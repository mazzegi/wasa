package app

import (
	"log"

	"github.com/mazzegi/wasa"
)

type Footer struct {
	root *wasa.Elt
	doc  *wasa.Document
}

func NewFooter(doc *wasa.Document) *Footer {
	e := &Footer{
		root: wasa.NewElt("footer", wasa.Class("footer")),
		doc:  doc,
	}
	e.setupUI()
	return e
}

func (e *Footer) Elt() *wasa.Elt {
	return e.root
}

func (e *Footer) setupUI() {
	cntElt := wasa.NewElt("span", wasa.Class("todo-count"), wasa.Data("0 item left"))

	ulFilterElt := wasa.NewElt("ul", wasa.Class("filters"))

	liAll := wasa.NewElt("li")
	aAll := wasa.NewElt("a", wasa.Class("selected"), wasa.Data("All"))
	e.doc.Callback(wasa.ClickEvent, aAll, func(e *wasa.Event) {
		log.Printf("filter:all:clicked")
	})
	liAll.Append(aAll)

	liActive := wasa.NewElt("li")
	aActive := wasa.NewElt("a", wasa.Data("Active"))
	e.doc.Callback(wasa.ClickEvent, aActive, func(e *wasa.Event) {
		log.Printf("filter:active:clicked")
	})
	liActive.Append(aActive)

	liCompleted := wasa.NewElt("li")
	aCompleted := wasa.NewElt("a", wasa.Data("Completed"))
	e.doc.Callback(wasa.ClickEvent, aCompleted, func(e *wasa.Event) {
		log.Printf("filter:completed:clicked")
	})
	liCompleted.Append(aCompleted)
	ulFilterElt.Append(
		liAll,
		liActive,
		liCompleted,
	)

	clearCompletedElt := wasa.NewElt("button", wasa.Class("clear-completed"), wasa.Data("Clear completed"))
	e.doc.Callback(wasa.ClickEvent, clearCompletedElt, func(e *wasa.Event) {
		log.Printf("clear completed")
	})

	e.root.Append(cntElt, ulFilterElt, clearCompletedElt)
	e.root.Hidden = true
}

func (e *Footer) render(r *repo) {
	if r.isEmpty() {
		e.root.Hidden = true
	} else {
		e.root.Hidden = false
	}
}
