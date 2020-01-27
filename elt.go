package wasa

import (
	"syscall/js"

	"github.com/pkg/errors"
)

type ElementCallback func(e *Event)

type Attrs map[string]string

type Elts []*Elt

type Elt struct {
	modified  bool
	jsElt     jsElt
	Tag       string
	Attrs     Attrs
	Childs    Elts
	Data      string
	Hidden    bool
	Callbacks map[string]ElementCallback
}

func (e *Elt) Invalidate() {
	e.modified = true
	for _, c := range e.Childs {
		c.Invalidate()
	}
}

func (e *Elt) accept() {
	e.modified = false
	for _, c := range e.Childs {
		c.accept()
	}
}

func (e *Elt) mount(doc *Document, parent jsElt) error {
	gsxElt := e.jsElt
	eNode, err := e.createJSElt(doc)
	if err != nil {
		return errors.Wrap(err, "create-element-node")
	}
	if !gsxElt.isValid() {
		parent.appendChild(eNode)
	} else {
		parent.replaceChild(gsxElt, eNode)
		gsxElt.remove()
	}
	e.accept()
	return nil
}

func (e *Elt) createJSElt(doc *Document) (jsElt, error) {
	eNode, err := doc.CreateElementNode(e.Tag)
	if err != nil {
		return undefinedJSElt(), errors.Wrap(err, "create element")
	}
	for k, v := range e.Attrs {
		eNode.setAttribute(k, v)
	}
	if e.Data != "" {
		eNode.setInnerHTML(e.Data)
	}
	for _, c := range e.Childs {
		if c.Hidden {
			continue
		}
		cNode, err := c.createJSElt(doc)
		if err != nil {
			return undefinedJSElt(), errors.Wrap(err, "create child element node")
		}
		eNode.appendChild(cNode)
	}
	e.jsElt = eNode
	return eNode, nil
}

func (e *Elt) Append(elts ...*Elt) {
	e.Childs = append(e.Childs, elts...)
}

func (e *Elt) RemoveAll() {
	for _, c := range e.Childs {
		c.RemoveAll()
		c.jsElt.remove()
	}
	e.Childs = Elts{}
}

func (e *Elt) Remove(re *Elt) {
	for i, c := range e.Childs {
		if c.jsElt.is(re.jsElt.jElt) {
			c.RemoveAll()
			c.jsElt.remove()
			e.Childs = append(e.Childs[:i], e.Childs[i+1:]...)
			return
		}
	}
}

func (e *Elt) AddAttr(k, v string) {
	if e.Attrs == nil {
		e.Attrs = Attrs{}
	}
	e.Attrs[k] = v
}

func (e *Elt) callback(event string, cb ElementCallback) {
	if e.Callbacks == nil {
		e.Callbacks = map[string]ElementCallback{}
	}
	e.Callbacks[event] = cb
}

func (e *Elt) findCallback(event string) (ElementCallback, bool) {
	if e.Callbacks == nil {
		return nil, false
	}
	cb, ok := e.Callbacks[event]
	return cb, ok
}

// // some access helpers
func (e *Elt) Value() string {
	return e.jsElt.jElt.Get("value").String()
}

func (e *Elt) findByTarget(target js.Value) (*Elt, bool) {
	if e.jsElt.is(target) {
		return e, true
	}
	for _, c := range e.Childs {
		if fc, ok := c.findByTarget(target); ok {
			return fc, true
		}
	}
	return nil, false
}

///
type EltMod func(e *Elt)

func NewElt(tag string, mods ...EltMod) *Elt {
	e := &Elt{
		Tag:       tag,
		Attrs:     Attrs{},
		Callbacks: map[string]ElementCallback{},
	}
	for _, mod := range mods {
		mod(e)
	}
	return e
}

func Attr(k, v string) EltMod {
	return func(e *Elt) {
		e.AddAttr(k, v)
	}
}

func ID(s string) EltMod {
	return func(e *Elt) {
		e.AddAttr("id", s)
	}
}

func Class(s string) EltMod {
	return func(e *Elt) {
		e.AddAttr("class", s)
	}
}

func Style(s string) EltMod {
	return func(e *Elt) {
		e.AddAttr("style", s)
	}
}

func Type(s string) EltMod {
	return func(e *Elt) {
		e.AddAttr("type", s)
	}
}

func Value(s string) EltMod {
	return func(e *Elt) {
		e.AddAttr("value", s)
	}
}

func Data(s string) EltMod {
	return func(e *Elt) {
		e.Data = s
	}
}

func Hidden(h bool) EltMod {
	return func(e *Elt) {
		e.Hidden = h
	}
}
