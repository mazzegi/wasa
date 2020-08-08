package wasa

import (
	"syscall/js"

	"github.com/mazzegi/wasa/errors"
	"github.com/mazzegi/wasa/wlog"
)

func isJSValueValid(v js.Value) bool {
	ty := v.Type()
	return ty != js.TypeUndefined && ty != js.TypeNull
}

type ElementCallback func(e *Event)

type Attrs map[string]string

type Elts []*Elt

type LCEvent string

const (
	Mounted   LCEvent = "mounted"
	Unmounted LCEvent = "unmounted"
	Rendered  LCEvent = "rendered"
)

type LCC struct {
	events map[LCEvent]func()
}

func (lcc *LCC) On(evt LCEvent, cb func()) {
	if lcc.events == nil {
		lcc.events = map[LCEvent]func(){}
	}
	lcc.events[evt] = cb
}

func (lcc *LCC) callback(evt LCEvent) {
	if cb, ok := lcc.events[evt]; ok {
		cb()
	}
}

type Elt struct {
	js.Value
	modified  bool
	hidden    bool
	Tag       string
	Attrs     Attrs
	Childs    Elts
	Data      string
	Callbacks map[string]ElementCallback
	key       string
	LCC       LCC
}

func (e *Elt) Key() string {
	return e.key
}

func (e *Elt) IsHidden() bool {
	return e.hidden
}

func (e *Elt) IsVisible() bool {
	return !e.hidden
}

func (e *Elt) Hide() {
	e.hidden = true
	for _, c := range e.Childs {
		c.Hide()
	}
}

func (e *Elt) Show() {
	e.hidden = false
	for _, c := range e.Childs {
		c.Show()
	}
}

func (e *Elt) mounted() {
	if e.hidden {
		return
	}
	e.LCC.callback(Mounted)
	for _, c := range e.Childs {
		c.mounted()
	}
}

func (e *Elt) unmounted() {
	if e.hidden {
		return
	}
	e.LCC.callback(Mounted)
	for _, c := range e.Childs {
		c.unmounted()
	}
}

func (e *Elt) rendered() {
	if e.hidden {
		return
	}
	e.LCC.callback(Rendered)
	for _, c := range e.Childs {
		c.rendered()
	}
}

func (e *Elt) Invalidate() {
	e.modified = true
	//e.LCC.callback(Unmounted)
	for _, c := range e.Childs {
		if c != nil {
			c.Invalidate()
		} else {
			wlog.Errorf("child is nil (this = %s, %v)", e.Tag, e.Attrs)
		}
	}
}

func (e *Elt) accept() {
	e.modified = false
	for _, c := range e.Childs {
		if c != nil {
			c.accept()
		} else {
			wlog.Errorf("child is nil (this = %s, %v)", e.Tag, e.Attrs)
		}
	}
}

func (e *Elt) mount(doc *Document, parent js.Value) error {
	oldV := e.Value
	newV, err := e.createJSElt(doc)
	if err != nil {
		return errors.Wrap(err, "create-element-node")
	}
	if !isJSValueValid(oldV) {
		parent.Call("appendChild", newV)
	} else {
		parent.Call("replaceChild", newV, oldV)
		oldV.Call("remove")
	}
	e.accept()
	e.mounted()
	return nil
}

func (e *Elt) createJSElt(doc *Document) (js.Value, error) {
	eNode, err := doc.CreateElementNode(e.Tag)
	if err != nil {
		return js.Undefined(), errors.Wrap(err, "create element")
	}
	for k, v := range e.Attrs {
		eNode.Call("setAttribute", k, v)
	}
	if e.Data != "" {
		eNode.Set("innerHTML", e.Data)
	}
	for _, c := range e.Childs {
		if c.hidden {
			continue
		}
		cNode, err := c.createJSElt(doc)
		if err != nil {
			return js.Undefined(), errors.Wrap(err, "create child element node")
		}
		eNode.Call("appendChild", cNode)
	}
	e.Value = eNode
	return eNode, nil
}

func (e *Elt) Append(elts ...*Elt) {
	e.Childs = append(e.Childs, elts...)
}

func (e *Elt) RemoveAll() {
	for _, c := range e.Childs {
		c.RemoveAll()
		c.Call("remove")
	}
	e.Childs = Elts{}
}

func (e *Elt) Remove(re *Elt) {
	for i, c := range e.Childs {
		if c.Value.Equal(re.Value) {
			c.RemoveAll()
			c.Call("remove")
			e.Childs = append(e.Childs[:i], e.Childs[i+1:]...)
			return
		}
	}
}

func (e *Elt) Replace(re *Elt, ne *Elt) {
	for i, c := range e.Childs {
		if c.Value.Equal(re.Value) {
			c.RemoveAll()
			c.Call("remove")
			e.Childs[i] = ne
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

// some access helpers
func (e *Elt) GetPath(names ...string) js.Value {
	curr := e.Value
	if !isJSValueValid(curr) {
		return curr
	}
	for _, name := range names {
		curr = curr.Get(name)
		if !isJSValueValid(curr) {
			return curr
		}
	}
	return curr
}

func (e *Elt) GetValue() string {
	return e.Get("value").String()
}

func (e *Elt) Is(target js.Value) bool {
	return e.Value.Equal(target)
}

func (e *Elt) stackToTarget(target js.Value) (match *Elt, stack []*Elt, found bool) {
	stack = []*Elt{e}
	if e.Is(target) {
		return e, stack, true
	}
	for _, c := range e.Childs {
		if fc, cstack, ok := c.stackToTarget(target); ok {
			stack = append(stack, cstack...)
			return fc, stack, true
		}
	}
	return nil, stack, false
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

func Key(k string) EltMod {
	return func(e *Elt) {
		e.key = k
	}
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
		e.hidden = h
	}
}
