package wasa

import (
	"syscall/js"

	"github.com/mazzegi/wasa/draw"
)

type Canvas struct {
	elt *Elt
	ctx js.Value
}

func NewCanvas(class string, id string) *Canvas {
	elt := NewElt("canvas", Class(class), ID(id))
	c := &Canvas{
		elt: elt,
	}
	return c
}

func (c *Canvas) InitCtx() {
	c.ctx = c.elt.Call("getContext", "2d")
	w := c.elt.Get("clientWidth").Int()
	h := c.elt.Get("clientHeight").Int()
	c.elt.jsElt.jElt.Set("width", w)
	c.elt.jsElt.jElt.Set("height", h)
}

func (c *Canvas) Clear() {
	w := c.elt.Get("width").Int()
	h := c.elt.Get("height").Int()
	c.ctx.Call("beginPath")
	c.ctx.Call("clearRect", 0, 0, w, h)

}

func (c *Canvas) Elt() *Elt {
	return c.elt
}

func (c *Canvas) MoveTo(p draw.Pt) {
	c.ctx.Call("moveTo", p.X, p.Y)
}

func (c *Canvas) LineTo(p draw.Pt) {
	c.ctx.Call("lineTo", p.X, p.Y)
}

func (c *Canvas) Stroke() {
	c.ctx.Call("stroke")
}

func (c *Canvas) Font(fs string) {
	c.ctx.Set("font", fs)
}

func (c *Canvas) FillStyle(fs string) {
	c.ctx.Set("fillStyle", fs)
}

func (c *Canvas) FillText(s string, p draw.Pt) {
	c.ctx.Call("fillText", s, p.X, p.Y)
}

func (c *Canvas) StrokeText(s string, p draw.Pt) {
	c.ctx.Call("strokeText", s, p.X, p.Y)
}
