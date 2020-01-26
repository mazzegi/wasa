package app

import "github.com/mazzegi/wasa"

type Footer struct {
	root *wasa.Elt
	doc  *wasa.Document
}

func NewFooter(doc *wasa.Document) *Footer {
	e := &Footer{
		root: wasa.NewElt("footer", wasa.Class("footer")),
	}
	e.setupUI()
	return e
}

func (e *Footer) Elt() *wasa.Elt {
	return e.root
}

func (e *Footer) setupUI() {

}
