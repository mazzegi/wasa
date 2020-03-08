package app

import (
	"fmt"

	"github.com/mazzegi/wasa"
	"github.com/mazzegi/wasa/example/todomvc/backend"
	"github.com/mazzegi/wasa/wlog"
)

type Footer struct {
	root       *wasa.Elt
	doc        *wasa.Document
	backend    *backend.Backend
	cntElt     *wasa.Elt
	aAll       *wasa.Elt
	aActive    *wasa.Elt
	aCompleted *wasa.Elt
}

func NewFooter(doc *wasa.Document, backend *backend.Backend) *Footer {
	e := &Footer{
		root:    wasa.NewElt("footer", wasa.Class("footer")),
		doc:     doc,
		backend: backend,
	}
	e.setupUI()
	return e
}

func (e *Footer) Elt() *wasa.Elt {
	return e.root
}

func (e *Footer) setupUI() {
	e.cntElt = wasa.NewElt("span", wasa.Class("todo-count"), wasa.Data("0 item left"))

	ulFilterElt := wasa.NewElt("ul", wasa.Class("filters"))
	ulFilterElt.Append(wasa.NewElt(wasa.StyleTag, wasa.Data(`
	a:hover{
		cursor: pointer;
	}
	`)))

	liAll := wasa.NewElt("li")
	e.aAll = wasa.NewElt("a", wasa.Class("selected"), wasa.Data("All"))
	e.doc.Callback(wasa.ClickEvent, e.aAll, func(evt *wasa.Event) {
		wlog.Infof("filter:all:clicked")
		e.backend.ChangeFilter(backend.All)
	})
	liAll.Append(e.aAll)

	liActive := wasa.NewElt("li")
	e.aActive = wasa.NewElt("a", wasa.Data("Active"))
	e.doc.Callback(wasa.ClickEvent, e.aActive, func(evt *wasa.Event) {
		wlog.Infof("filter:active:clicked")
		e.backend.ChangeFilter(backend.Active)
	})
	liActive.Append(e.aActive)

	liCompleted := wasa.NewElt("li")
	e.aCompleted = wasa.NewElt("a", wasa.Data("Completed"))
	e.doc.Callback(wasa.ClickEvent, e.aCompleted, func(evt *wasa.Event) {
		wlog.Infof("filter:completed:clicked")
		e.backend.ChangeFilter(backend.Completed)
	})
	liCompleted.Append(e.aCompleted)
	ulFilterElt.Append(
		liAll,
		liActive,
		liCompleted,
	)

	clearCompletedElt := wasa.NewElt("button", wasa.Class("clear-completed"), wasa.Data("Clear completed"))
	e.doc.Callback(wasa.ClickEvent, clearCompletedElt, func(evt *wasa.Event) {
		wlog.Infof("clear completed")
		e.backend.DeleteCompleted()
	})

	e.root.Append(e.cntElt, ulFilterElt, clearCompletedElt)
	e.root.Hidden = true
}

func (e *Footer) render() {
	if e.backend.IsEmpty() {
		e.root.Hidden = true
	} else {
		e.root.Hidden = false
	}
	wasa.Data(fmt.Sprintf("%d item(s) left", e.backend.Count()))(e.cntElt)

	wasa.Class("")(e.aAll)
	wasa.Class("")(e.aActive)
	wasa.Class("")(e.aCompleted)
	switch e.backend.Filter() {
	case backend.All:
		wasa.Class("selected")(e.aAll)
	case backend.Active:
		wasa.Class("selected")(e.aActive)
	case backend.Completed:
		wasa.Class("selected")(e.aCompleted)
	}

	//e.cntElt.Invalidate()
}
