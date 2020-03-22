package app

import (
	"github.com/mazzegi/wasa"
)

var styleComponent = `
html{
    height: 100%;
}
body{
    height: 100%;
    margin: 0;
    font-family: Arial, Helvetica, sans-serif; 
    font-size: 13.333px;   
}
.root{
    height: 100%;	
    margin: auto;
    width: 70%;
}
`

type Type1 struct {
	root *wasa.Elt
}

func (e *Type1) Elt() *wasa.Elt {
	return e.root
}

func NewType1() *Type1 {
	e := &Type1{}
	e.root = wasa.NewElt("div")
	e.root.Append(wasa.NewElt(wasa.StyleTag, wasa.Data(styleComponent)))
	return e
}

type Complex struct {
	root *wasa.Elt
	List *wasa.Elt
}

func (e *Complex) Elt() *wasa.Elt {
	return e.root
}

func NewComplex() *Complex {
	e := &Complex{}
	e.root = wasa.NewElt("div")
	wasa.Attr("class", "complex")(e.root)

	e.List = wasa.NewElt("ul")
	e.root.Append(e.List)
	return e
}

type Component struct {
	root      *wasa.Elt
	Label     *wasa.Elt
	Button    *wasa.Elt
	Complex   *Complex
	Clockwise *wasa.Elt
}

func (e *Component) Elt() *wasa.Elt {
	return e.root
}

func NewComponent() *Component {
	e := &Component{}
	e.root = wasa.NewElt("div")
	e.root.Append(wasa.NewElt(wasa.StyleTag, wasa.Data(styleComponent)))
	wasa.Attr("class", "component")(e.root)

	e.Label = wasa.NewElt("label")
	wasa.Data("Label-For")(e.Label)
	e.root.Append(e.Label)

	e.Button = wasa.NewElt("button")
	wasa.Data("Button")(e.Button)
	e.root.Append(e.Button)

	e.Complex = NewComplex()
	e.root.Append(e.Complex.Elt())

	e.Clockwise = wasa.NewElt("yield")
	e.root.Append(e.Clockwise)
	return e
}
