package app

import "github.com/mazzegi/wasa"

type Main struct {
	root     *wasa.Elt
	doc      *wasa.Document
	todoList *wasa.Elt
}

func NewMain(doc *wasa.Document) *Main {
	e := &Main{
		root: wasa.NewElt("section", wasa.Class("main")),
		doc:  doc,
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
	e.todoList = wasa.NewElt("ul", wasa.Class("todo-list"))

	e.root.Append(
		toggleAll,
		toggleAllLabel,
		e.todoList,
	)
	e.root.Hidden = true
}

func (e *Main) render(r *repo) {
	if r.isEmpty() {
		e.root.Hidden = true
	} else {
		e.root.Hidden = false
	}
	e.todoList.RemoveAll()
	for _, item := range r.items {
		e.renderItem(item)
	}
}

func (e *Main) renderItem(item item) {
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
	label := wasa.NewElt("label", wasa.Data(item.Text))
	delete := wasa.NewElt("button", wasa.Class("destroy"))
	view.Append(toggleCompleted, label, delete)

	e.todoList.Append(li)
}
