package app

import "github.com/mazzegi/wasa"

type Main struct {
	root *wasa.Elt
	doc  *wasa.Document
}

func NewMain(doc *wasa.Document) *Main {
	e := &Main{
		root: wasa.NewElt("section", wasa.Class("main")),
	}
	e.setupUI()
	return e
}

func (e *Main) Elt() *wasa.Elt {
	return e.root
}

func (e *Main) setupUI() {

}
