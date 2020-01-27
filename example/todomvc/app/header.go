package app

import (
	"log"

	"github.com/mazzegi/wasa"
)

type Header struct {
	root *wasa.Elt
	doc  *wasa.Document
}

func NewHeader(doc *wasa.Document) *Header {
	h := &Header{
		root: wasa.NewElt("header", wasa.Class("header")),
		doc:  doc,
	}
	h.setupUI()
	return h
}

func (e *Header) Elt() *wasa.Elt {
	return e.root
}

func (e *Header) setupUI() {
	e.root.Append(wasa.NewElt(wasa.H1Tag))
	inElt := wasa.NewElt(
		wasa.InputTag,
		wasa.Class("new-todo"),
		wasa.Attr("placeholder", "What needs to be done?"),
		wasa.Attr("autofocus", ""),
	)
	e.doc.Callback(wasa.KeyupEvent, inElt, func(e *wasa.Event) {
		log.Printf("input: keyup (%s)", e.JSEvent().Get("keyCode").String())
	})
	e.root.Append(inElt)
}

func (e *Header) render(r *repo) {

}
