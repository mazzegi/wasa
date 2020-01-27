package app

import (
	"github.com/mazzegi/wasa"
	"github.com/mazzegi/wasa/example/todomvc/backend"
)

type Header struct {
	root    *wasa.Elt
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
	inElt := wasa.NewElt(
		wasa.InputTag,
		wasa.Class("new-todo"),
		wasa.Attr("placeholder", "What needs to be done?"),
		wasa.Attr("autofocus", ""),
	)
	e.doc.Callback(wasa.KeyupEvent, inElt, func(evt *wasa.Event) {
		code := evt.JSEvent().Get("keyCode").Int()
		if code == 13 {
			e.backend.Add(inElt.Value())
		}
	})
	e.root.Append(inElt)
}

func (e *Header) render() {

}
