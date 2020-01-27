package app

import (
	"github.com/mazzegi/wasa"
	"github.com/mazzegi/wasa/example/todomvc/backend"
	"log"
)

type Main struct {
	root     *wasa.Elt
	doc      *wasa.Document
	backend  *backend.Backend
	todoList *wasa.Elt
}

func NewMain(doc *wasa.Document, backend *backend.Backend) *Main {
	e := &Main{
		root:    wasa.NewElt("section", wasa.Class("main")),
		doc:     doc,
		backend: backend,
	}
	e.setupUI()
	return e
}

func (e *Main) Elt() *wasa.Elt {
	return e.root
}

func (e *Main) setupUI() {
	toggleAll := wasa.NewElt("input", wasa.Attr("type", "checkpoint"), wasa.Class("toggle-all"), wasa.ID("toggle-all"))
	toggleAllLabel := wasa.NewElt("label", wasa.Attr("for", "toggle-all"), wasa.Data("Mark all as complete"))
	e.doc.Callback(wasa.ClickEvent, toggleAll, func(evt *wasa.Event) {
		e.backend.ToggleAll()
	})

	e.todoList = wasa.NewElt("ul", wasa.Class("todo-list"))
	e.root.Append(
		toggleAll,
		toggleAllLabel,
		e.todoList,
	)
	e.root.Hidden = true
}

func (e *Main) render() {
	if e.backend.IsEmpty() {
		e.root.Hidden = true
	} else {
		e.root.Hidden = false
	}
	e.todoList.RemoveAll()
	e.backend.Each(e.renderItem)
	//e.todoList.Invalidate()
}

func (e *Main) renderItem(item *backend.Item) {
	li := wasa.NewElt("li")
	if item.Completed {
		wasa.Class("completed")(li)
	}
	view := wasa.NewElt("div", wasa.Class("view"))
	li.Append(view)

	toggleCompleted := wasa.NewElt("input", wasa.Class("toggle"), wasa.Type("checkbox"))
	if item.Completed {
		wasa.Attr("checked", "")(toggleCompleted)
	}
	e.doc.Callback(wasa.ClickEvent, toggleCompleted, func(evt *wasa.Event) {
		e.backend.ToggleComplete(item.ID)
		log.Printf("toggle-completed (%d)", item.ID)
	})

	label := wasa.NewElt("label", wasa.Data(item.Text))
	deleteBtn := wasa.NewElt("button", wasa.Class("destroy"))
	e.doc.Callback(wasa.ClickEvent, deleteBtn, func(evt *wasa.Event) {
		e.backend.Delete(item.ID)
		log.Printf("delete (%d)", item.ID)
	})

	view.Append(toggleCompleted, label, deleteBtn)
	e.todoList.Append(li)
}
