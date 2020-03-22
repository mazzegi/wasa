package app

import (
	"github.com/mazzegi/wasa"
	"github.com/mazzegi/wasa/wlog"
)

type App struct {
	comp *Component
	doc  *wasa.Document
}

func New() (*App, error) {
	doc, err := wasa.NewDocument("Generate Example")
	if err != nil {
		return nil, err
	}
	a := &App{
		doc: doc,
	}

	a.setupUI()
	return a, nil
}

func (a *App) Run() {
	wlog.Infof("app: run ...")
	a.doc.Run(a.comp.Elt())
	wlog.Infof("app: run ... done")
}

func (a *App) setupUI() {
	a.comp = NewComponent()
}
