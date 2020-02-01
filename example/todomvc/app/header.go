package app

import (
	"github.com/mazzegi/wasa"
	"github.com/mazzegi/wasa/example/todomvc/backend"
)

type Header struct {
	root    *wasa.Elt
	input   *wasa.Elt
	doc     *wasa.Document
	backend *backend.Backend
}

func NewHeader(doc *wasa.Document, backend *backend.Backend) *Header {
	h := &Header{
		root:    wasa.NewElt("header", wasa.Class("header")),
		doc:     doc,
		backend: backend,
	}
	h.setupUI()
	return h
}

func (e *Header) Elt() *wasa.Elt {
	return e.root
}

func (e *Header) setupUI() {
	e.root.Append(wasa.NewElt(wasa.H1Tag))
	e.input = wasa.NewElt(
		wasa.InputTag,
		wasa.Class("new-todo"),
		wasa.Attr("placeholder", "What needs to be done?"),
		wasa.Attr("autofocus", ""),
	)
	e.doc.Callback(wasa.KeyupEvent, e.input, func(evt *wasa.Event) {
		code := evt.JSEvent().Get("keyCode").Int()
		if code == 13 {
			e.backend.Add(e.input.Value())
		}
	})
	e.root.Append(e.input)
}

func (e *Header) render() {
	e.doc.Focus(e.input)
}
