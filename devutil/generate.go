package devutil

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
	types map[string]struct{}
	names map[string]struct{}
}

func NewGenerator(opts ...Option) (*Generator, error) {
	g := &Generator{
		types: map[string]struct{}{},
		names: map[string]struct{}{},
	}
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
	fmt.Printf("*** process %q ****\n", inFile)
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

func (g *Generator) wasaType(n *html.Node, prefix string) string {
	for _, attr := range n.Attr {
		if attr.Key == "wasa-type" {
			return prefix + norm(attr.Val)
		}
	}
	return ""
}

func (g *Generator) wasaName(n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == "wasa-name" {
			return norm(attr.Val)
		}
	}
	return ""
}

func (g *Generator) uniqueType(ty string) string {
	contains := func(s string) bool {
		if _, ok := g.types[s]; ok {
			return true
		}
		return false
	}
	i := 0
	cand := ty
	for contains(cand) {
		i++
		cand = fmt.Sprintf("%s%d", ty, i)
	}
	return cand
}

func (g *Generator) uniqueName(name string) string {
	fmt.Printf("find unique name for %q\n", name)
	contains := func(s string) bool {
		if _, ok := g.names[s]; ok {
			return true
		}
		return false
	}
	i := 0
	cand := name
	for contains(cand) {
		i++
		cand = fmt.Sprintf("%s_%d", name, i)
	}
	fmt.Printf("return unique name for %q\n", cand)
	return cand
}

func (g *Generator) forceWasaType(n *html.Node, prefix string) string {
	if s := g.wasaType(n, prefix); s != "" {
		return s
	}
	ty := g.uniqueType(prefix + norm(n.Data))
	g.useType(ty)
	return ty
}

func (g *Generator) forceWasaName(n *html.Node) string {
	if s := g.wasaName(n); s != "" {
		return s
	}
	name := g.uniqueName(norm(n.Data))
	g.useName(name)
	return name
}

func (g *Generator) useType(ty string) {
	g.types[ty] = struct{}{}
}

func (g *Generator) useName(name string) {
	g.names[name] = struct{}{}
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

	prefix := ""
	for _, n := range nodes {
		if n.Type == html.ElementNode && n.Data == "wasa" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "prefix" && c.FirstChild != nil && c.FirstChild.Type == html.TextNode {
					prefix = c.FirstChild.Data
					fmt.Printf("assign prefix: %s\n", prefix)
				}
			}
		}
	}

	types := []*structType{}
	for _, n := range nodes {
		if mustSkip(n) {
			continue
		}
		if n.Type == html.ElementNode && n.Data == "wasa" {
			continue
		}
		fmt.Printf("process %s\n", n.Data)
		wsType := g.forceWasaType(n, prefix)
		st := newStructType(wsType, wsType)
		err := st.process(g, n, prefix)
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
		varStyle = fmt.Sprintf("style%s%s", prefix, strings.Title(norm(name)))
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

type yieldedElement struct {
	name string
	typ  string
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

func (t *structType) process(g *Generator, n *html.Node, prefix string) error {
	t.tag = n.Data
	t.attrs = n.Attr
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if mustSkip(c) {
			continue
		}
		wsName := g.forceWasaName(c)
		if c.Type == html.ElementNode && c.Data == "yield" {
			wsType := g.wasaType(c, prefix)
			if wsType == "" {
				fmt.Printf("WARN: skipping yield-node due to missing wasa-type\n")
				continue
			}
			t.childs = append(t.childs, &yieldedElement{
				name: wsName,
				typ:  wsType,
			})
			fmt.Printf("added yielded elt: %q/%q\n", wsName, wsType)
		} else if c.FirstChild == nil {
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
			wsType := g.forceWasaType(c, prefix)
			st := newStructType(wsType, wsName)
			err := st.process(g, c, prefix)
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
		case *yieldedElement:
			writeLine("    %s *%s", c.name, c.typ)
		}
	}
	writeLine("}\n")

	//get-element
	writeLine("func (e *%s) Elt() *wasa.Elt {", t.typeName)
	writeLine("    return e.root")
	writeLine("}\n")

	//append-element
	writeLine("func (e *%s) Append(elt ...*wasa.Elt) {", t.typeName)
	writeLine("    e.root.Append(elt...)")
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
		case *yieldedElement:
			writeLine("    e.%s = New%s()", child.name, child.typ)
			writeLine("    e.root.Append(e.%s.Elt())", child.name)
		}
	}
	writeLine("    return e")
	writeLine("}")

	subs = append(subs, code)
	return subs
}
