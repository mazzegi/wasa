package gen

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

type Option func(*Generator) error

type Generator struct {
}

func NewGenerator(opts ...Option) (*Generator, error) {
	g := &Generator{}
	for _, opt := range opts {
		err := opt(g)
		if err != nil {
			return nil, err
		}
	}
	return g, nil
}

func (g *Generator) ProcessDir(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.wasa.html"))
	if err != nil {
		return errors.Wrap(err, "glob")
	}
	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), ".wasa.html")
		out := filepath.Join(dir, fmt.Sprintf("wasagen_%s.go", name))
		err := g.ProcessFile(file, out, name, filepath.Base(dir))
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) ProcessFile(inFile, outFile string, name, pkg string) error {
	fIn, err := os.Open(inFile)
	if err != nil {
		return errors.Wrapf(err, "open (%s)", inFile)
	}
	defer fIn.Close()
	fOut, err := os.Create(outFile)
	if err != nil {
		return errors.Wrapf(err, "create (%s)", outFile)
	}
	defer fOut.Close()
	return g.Process(fIn, fOut, name, pkg)
}

func mustSkip(n *html.Node) bool {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "html", "head", "body", "script", "style":
			return true
		default:
			return false
		}
	}
	return true
}

func wasaType(n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == "wasa-type" {
			return norm(attr.Val)
		}
	}
	return ""
}

func wasaName(n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == "wasa-name" {
			return norm(attr.Val)
		}
	}
	return ""
}

func norm(s string) string {
	s = strings.Trim(s, " \r\n\t")
	s = strings.Title(s)
	s = strings.NewReplacer(
		"-", "_",
	).Replace(s)
	return s
}

func style(nodes []*html.Node) string {
	for _, n := range nodes {
		if n.Type == html.ElementNode && n.Data == "style" {
			if n.FirstChild.Type == html.TextNode {
				data := strings.Trim(n.FirstChild.Data, " \r\n\t")
				if data != "" {
					return data
				}
			}
		}
	}
	return ""
}

func (g *Generator) Process(in io.Reader, out io.Writer, name, pkg string) error {
	ctxNode := &html.Node{
		Type: html.ElementNode,
	}
	nodes, err := html.ParseFragment(in, ctxNode)
	if err != nil {
		return errors.Wrap(err, "parse html")
	}

	types := []*structType{}
	for _, n := range nodes {
		if mustSkip(n) {
			continue
		}
		fmt.Printf("process %s\n", n.Data)
		wsType := wasaType(n)
		if wsType == "" {
			fmt.Printf("WARN: skipping node due to missing wasa-type\n")
			continue
		}
		st := newStructType(wsType, wsType)
		err := st.process(n)
		if err != nil {
			return errors.Wrapf(err, "process type (%s)", wsType)
		}
		types = append(types, st)
	}

	fmt.Fprintf(out, "package %s\n\n", pkg)
	fmt.Fprintf(out, "import (\n    \"github.com/mazzegi/wasa\"\n)\n\n")

	style := style(nodes)
	varStyle := ""
	if style != "" {
		varStyle = fmt.Sprintf("style%s", strings.Title(name))
		code := fmt.Sprintf("var %s = `\n", varStyle)
		code += style + "\n`\n\n"
		fmt.Fprint(out, code)
	}

	for _, ty := range types {
		subs := ty.generateCode(varStyle)
		for _, sub := range subs {
			fmt.Fprint(out, sub)
			fmt.Fprint(out, "\n")
		}
	}

	return nil
}

type simpleElement struct {
	name  string
	tag   string
	attrs []html.Attribute
	data  string
}

type structType struct {
	typeName string
	name     string
	tag      string
	style    string
	attrs    []html.Attribute
	childs   []interface{}
}

func newStructType(typeName string, name string) *structType {
	return &structType{
		typeName: typeName,
		name:     name,
	}
}

func (t *structType) process(n *html.Node) error {
	t.tag = n.Data
	t.attrs = n.Attr
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if mustSkip(c) {
			continue
		}
		wsName := wasaName(c)
		if wsName == "" {
			fmt.Printf("WARN: skipping node due to missing wasa-name\n")
			continue
		}
		if c.FirstChild == nil {
			t.childs = append(t.childs, &simpleElement{
				name:  wsName,
				tag:   c.Data,
				attrs: c.Attr,
			})
		} else if c.FirstChild.Type == html.TextNode && c.FirstChild.NextSibling == nil {
			//just one text-node
			data := strings.Trim(c.FirstChild.Data, " \r\n\t")
			t.childs = append(t.childs, &simpleElement{
				name:  wsName,
				tag:   c.Data,
				attrs: c.Attr,
				data:  data,
			})
		} else {
			//first child is a non-text element - start new struct type
			wsType := wasaType(c)
			if wsType == "" {
				fmt.Printf("WARN: skipping node due to missing wasa-type\n")
				continue
			}

			st := newStructType(wsType, wsName)
			err := st.process(c)
			if err != nil {
				return errors.Wrapf(err, "process (%s)", wsType)
			}
			t.childs = append(t.childs, st)
		}
	}
	return nil
}

func (t *structType) generateCode(varStyle string) []string {
	code := ""
	writeLine := func(s string, args ...interface{}) {
		code += fmt.Sprintf(s+"\n", args...)
	}

	subs := []string{}
	writeLine("type %s struct {", t.typeName)
	writeLine("    root *wasa.Elt")
	for _, c := range t.childs {
		switch c := c.(type) {
		case *simpleElement:
			writeLine("    %s *wasa.Elt", c.name)
		case *structType:
			writeLine("    %s *%s", c.name, c.typeName)
			subs = append(subs, c.generateCode("")...)
		}
	}
	writeLine("}\n")

	//get-element
	writeLine("func (e *%s) Elt() *wasa.Elt {", t.typeName)
	writeLine("    return e.root")
	writeLine("}\n")

	//constructor
	writeLine("func New%s() *%s {", t.typeName, t.typeName)
	writeLine("    e := &%s{}", t.typeName)
	writeLine("    e.root = wasa.NewElt(%q)", t.tag)
	if varStyle != "" {
		writeLine("    e.root.Append(wasa.NewElt(wasa.StyleTag, wasa.Data(%s)))", varStyle)
	}
	for _, attr := range t.attrs {
		if !strings.HasPrefix(attr.Key, "wasa") {
			writeLine("    wasa.Attr(%q, %q)(e.root)", attr.Key, attr.Val)
		}
	}
	for _, child := range t.childs {
		writeLine("")
		switch child := child.(type) {
		case *simpleElement:
			writeLine("    e.%s = wasa.NewElt(%q)", child.name, child.tag)
			for _, attr := range child.attrs {
				if !strings.HasPrefix(attr.Key, "wasa") {
					writeLine("    wasa.Attr(%q, %q)(e.%s)", attr.Key, attr.Val, child.name)
				}
			}
			if child.data != "" {
				writeLine("    wasa.Data(%q)(e.%s)", child.data, child.name)
			}
			writeLine("    e.root.Append(e.%s)", child.name)
		case *structType:
			writeLine("    e.%s = New%s()", child.name, child.typeName)
			writeLine("    e.root.Append(e.%s.Elt())", child.name)
		}
	}
	writeLine("    return e")
	writeLine("}")

	subs = append(subs, code)
	return subs
}
