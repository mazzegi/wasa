package app

import (
	"log"

	"github.com/mazzegi/wasa"
)

type App struct {
	root   *wasa.Elt
	doc    *wasa.Document
	repo   *repo
	header *Header
	main   *Main
	footer *Footer
}

func New() (*App, error) {
	doc, err := wasa.NewDocument("WASA / TodoMVC")
	if err != nil {
		return nil, err
	}
	a := &App{
		root: wasa.NewElt("section", wasa.Class("todoapp")),
		doc:  doc,
		repo: newRepo(),
	}

	a.root.Append(wasa.NewElt(wasa.StyleTag, wasa.Data(baseCSS)))
	a.root.Append(wasa.NewElt(wasa.StyleTag, wasa.Data(indexCSS)))

	a.setupUI()
	a.render()
	return a, nil
}

func (a *App) Run() {
	log.Printf("app: run ...")
	a.doc.Run(a.root)
	log.Printf("app: run ... done")
}

func (a *App) setupUI() {
	a.header = NewHeader(a.doc)
	a.main = NewMain(a.doc)
	a.footer = NewFooter(a.doc)
	a.root.Append(a.header.Elt(), a.main.Elt(), a.footer.Elt())
}

func (a *App) render() {
	a.header.render(a.repo)
	a.main.render(a.repo)
	a.footer.render(a.repo)
	a.root.Invalidate()
}
