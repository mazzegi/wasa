package app

import (
	"log"

	"github.com/mazzegi/wasa"
)

type App struct {
	root *wasa.Elt
	doc  *wasa.Document
}

func New() (*App, error) {
	doc, err := wasa.NewDocument("WASA / TodoMVC")
	if err != nil {
		return nil, err
	}
	a := &App{
		root: wasa.NewElt("section", wasa.Class("todoapp")),
		doc:  doc,
	}

	a.root.Append(wasa.NewElt(wasa.StyleTag, wasa.Data(baseCSS)))
	a.root.Append(wasa.NewElt(wasa.StyleTag, wasa.Data(indexCSS)))

	a.setupUI()
	return a, nil
}

func (a *App) Run() {
	log.Printf("app: run ...")
	a.doc.Run(a.root)
	log.Printf("app: run ... done")
}

func (a *App) setupUI() {
	header := NewHeader(a.doc)
	main := NewMain(a.doc)
	footer := NewFooter(a.doc)
	a.root.Append(header.Elt(), main.Elt(), footer.Elt())
}
