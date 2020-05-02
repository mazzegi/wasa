package app

import (
	"math/rand"

	"github.com/mazzegi/wasa"
	"github.com/mazzegi/wasa/draw"
	"github.com/mazzegi/wasa/wlog"
)

type App struct {
	root   *wasa.Elt
	doc    *wasa.Document
	canvas *wasa.Canvas
}

func New() (*App, error) {
	doc, err := wasa.NewDocument("Canvas Example")
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
	a.doc.Run(a.root)
	wlog.Infof("app: run ... done")
}

func (a *App) setupUI() {
	a.root = wasa.NewElt("div", wasa.Class("root"))
	a.root.Append(wasa.NewElt(wasa.StyleTag, wasa.Data(appStyles)))

	canvasContainer := wasa.NewElt(wasa.DivTag, wasa.Class("canvas-container"))
	a.canvas = wasa.NewCanvas("canvas", "canvas")
	canvasContainer.Append(a.canvas.Elt())
	canvasContainer.LCC.On(wasa.Mounted, func() {
		a.canvas.InitCtx()
	})

	controlContainer := wasa.NewElt(wasa.DivTag, wasa.Class("control-container"))
	btn := wasa.NewElt(wasa.ButtonTag, wasa.Data("Start"))
	a.doc.Callback(wasa.ClickEvent, btn, func(evt *wasa.Event) {
		a.renderCanvas()
	})
	controlContainer.Append(btn)

	a.root.Append(
		canvasContainer,
		controlContainer,
	)

	// a.doc.AfterMount(func() {
	// 	a.canvas.InitCtx()
	// })
}

func (a *App) renderCanvas() {
	a.canvas.Clear()
	a.canvas.MoveTo(draw.P(10, 10))
	a.canvas.LineTo(draw.P(100+rand.Intn(50), 100+rand.Intn(50)))
	a.canvas.Stroke()
	a.canvas.Font("24px Arial")

	a.canvas.FillStyle("#009900")
	a.canvas.FillText("😋", draw.P(200+rand.Intn(50), 200+rand.Intn(50)))
	a.canvas.FillStyle("#990000")
	a.canvas.FillText("😞", draw.P(200+rand.Intn(50), 200+rand.Intn(50)))

}

var appStyles = `
html{
	height: 100%;
}
body{
	height: 100%;
    margin: 0;
	font-family: Arial, Helvetica, sans-serif; 
	font-size: 13.333px;   
}
.root{
	height: 100%;	
	margin: auto;
	width: 70%;
}
.canvas-container{
	margin-top: 80px;
	border: 1px solid lightgray;
}
.control-container{
	border: 1px solid lightgray;
}
.canvas{
	width: 100%;
	height: 600px;
}
`
