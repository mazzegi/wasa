package app

import (
	"github.com/mazzegi/wasa"
	"github.com/mazzegi/wasa/example/todomvc/backend"
	"github.com/mazzegi/wasa/wlog"
)

type App struct {
	root    *wasa.Elt
	doc     *wasa.Document
	backend *backend.Backend
	header  *Header
	main    *Main
	footer  *Footer
}

func New(be *backend.Backend) (*App, error) {
	doc, err := wasa.NewDocument("WASA / TodoMVC")
	if err != nil {
		return nil, err
	}
	a := &App{
		root:    wasa.NewElt("section", wasa.Class("todoapp")),
		doc:     doc,
		backend: be,
	}
	a.backend.Subscribe(func() {
		a.render()
	})

	api := doc.GetGlobal("wasaenv", "api").String()
	wlog.Infof("api=(%s)", api)
	loc := doc.Location()
	wlog.Infof("loc=(%s)", loc)

	a.setupUI()
	a.render()
	return a, nil
}

func (a *App) Run() {
	wlog.Infof("app: run ...")
	a.doc.Run(a.root)
	wlog.Infof("app: run ... done")
}

func (a *App) setupUI() {
	a.header = NewHeader(a.doc, a.backend)
	a.main = NewMain(a.doc, a.backend)
	a.footer = NewFooter(a.doc, a.backend)
	a.root.Append(a.header.Elt(), a.main.Elt(), a.footer.Elt())
}

func (a *App) render() {
	wlog.Infof("app: render ...")
	a.header.render()
	a.main.render()
	a.footer.render()
	a.root.Invalidate()
	wlog.Infof("app: render ... done")
}
