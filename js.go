package wasa

import (
	"syscall/js"
)

func isJSValueUndefined(v js.Value) bool {
	return v.Type() == js.TypeUndefined
}

func isJSValueNull(v js.Value) bool {
	return v.Type() == js.TypeNull
}

func isJSValueValid(v js.Value) bool {
	return !isJSValueUndefined(v) && !isJSValueNull(v)
}

type jsNode interface {
	jsValue() js.Value
}

type jsElt struct {
	jElt js.Value
}

func newJSElt(v js.Value) jsElt {
	return jsElt{
		jElt: v,
	}
}

func undefinedJSElt() jsElt {
	return jsElt{
		jElt: js.Undefined(),
	}
}

func (e jsElt) is(v js.Value) bool {
	if !isJSValueValid(e.jElt) {
		return false
	}
	return e.jElt == v
}

func (e jsElt) isValid() bool {
	return isJSValueValid(e.jElt)
}

func (e jsElt) jsValue() js.Value {
	return e.jElt
}

func (e jsElt) call(method string, args ...interface{}) js.Value {
	return e.jElt.Call(method, args...)
}

func (e jsElt) get(names ...string) js.Value {
	curr := e.jElt
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

func (e jsElt) appendChild(n jsNode) {
	if !isJSValueValid(e.jElt) {
		return
	}
	e.jElt.Call("appendChild", n.jsValue())
}

func (e jsElt) remove() {
	if !isJSValueValid(e.jElt) {
		return
	}
	e.jElt.Call("remove")
}

func (e jsElt) replaceChild(which jsNode, with jsNode) {
	if !isJSValueValid(e.jElt) {
		return
	}
	e.jElt.Call("replaceChild", with.jsValue(), which.jsValue())
}

func (e jsElt) setInnerHTML(html string) {
	if !isJSValueValid(e.jElt) {
		return
	}
	e.jElt.Set("innerHTML", html)
}

func (e jsElt) setAttribute(key, val string) {
	if !isJSValueValid(e.jElt) {
		return
	}
	e.jElt.Call("setAttribute", key, val)
}

func (e jsElt) addEventListener(event string, cb js.Func) {
	if !isJSValueValid(e.jElt) {
		return
	}
	e.jElt.Call("addEventListener", event, cb)
}
