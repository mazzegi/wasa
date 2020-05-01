package app

import (
    "github.com/mazzegi/wasa"
)

var styleDetailsComponent = `
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

type DetailsType1 struct {
    root *wasa.Elt
}

func (e *DetailsType1) Elt() *wasa.Elt {
    return e.root
}

func NewDetailsType1() *DetailsType1 {
    e := &DetailsType1{}
    e.root = wasa.NewElt("div")
    e.root.Append(wasa.NewElt(wasa.StyleTag, wasa.Data(styleDetailsComponent)))
    return e
}

type DetailsDiv struct {
    root *wasa.Elt
    Span *wasa.Elt
}

func (e *DetailsDiv) Elt() *wasa.Elt {
    return e.root
}

func NewDetailsDiv() *DetailsDiv {
    e := &DetailsDiv{}
    e.root = wasa.NewElt("div")

    e.Span = wasa.NewElt("span")
    wasa.Data("Hello World")(e.Span)
    e.root.Append(e.Span)
    return e
}

type DetailsComplex struct {
    root *wasa.Elt
    List *wasa.Elt
}

func (e *DetailsComplex) Elt() *wasa.Elt {
    return e.root
}

func NewDetailsComplex() *DetailsComplex {
    e := &DetailsComplex{}
    e.root = wasa.NewElt("div")
    wasa.Attr("class", "complex")(e.root)

    e.List = wasa.NewElt("ul")
    e.root.Append(e.List)
    return e
}

type DetailsComponent struct {
    root *wasa.Elt
    Label *wasa.Elt
    Button *wasa.Elt
    Div *DetailsDiv
    Complex *DetailsComplex
    Clockwise *DetailsType1
}

func (e *DetailsComponent) Elt() *wasa.Elt {
    return e.root
}

func NewDetailsComponent() *DetailsComponent {
    e := &DetailsComponent{}
    e.root = wasa.NewElt("div")
    e.root.Append(wasa.NewElt(wasa.StyleTag, wasa.Data(styleDetailsComponent)))
    wasa.Attr("class", "component")(e.root)

    e.Label = wasa.NewElt("label")
    wasa.Data("Label-For Nothing")(e.Label)
    e.root.Append(e.Label)

    e.Button = wasa.NewElt("button")
    wasa.Data("Button")(e.Button)
    e.root.Append(e.Button)

    e.Div = NewDetailsDiv()
    e.root.Append(e.Div.Elt())

    e.Complex = NewDetailsComplex()
    e.root.Append(e.Complex.Elt())

    e.Clockwise = NewDetailsType1()
    e.root.Append(e.Clockwise.Elt())
    return e
}

