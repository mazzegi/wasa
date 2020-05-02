package wasa

import (
	"net/url"
	"syscall/js"

	"github.com/mazzegi/wasa/errors"
	"github.com/mazzegi/wasa/wlog"
)

type Document struct {
	js.Value
	glb     js.Value
	jDoc    js.Value
	body    js.Value
	events  map[string]struct{}
	root    *Elt
	renderC chan struct{}
}

func NewDocument(title string) (*Document, error) {
	glb := js.Global()
	if !isJSValueValid(glb) {
		return nil, errors.Errorf("js-global is not valid (%s)", glb.Type().String())
	}
	jDoc := glb.Get("document")
	if !isJSValueValid(jDoc) {
		return nil, errors.Errorf("js-document is not valid (%s)", jDoc.Type().String())
	}
	doc := &Document{
		Value:   jDoc,
		glb:     glb,
		jDoc:    jDoc,
		events:  map[string]struct{}{},
		renderC: make(chan struct{}),
	}
	err := doc.createBodyIfNotExists()
	if err != nil {
		return nil, errors.Wrap(err, "create-body")
	}
	jDoc.Set("title", title)
	return doc, nil
}

func (d *Document) createBodyIfNotExists() error {
	bodySlice := d.jDoc.Call("getElementsByTagName", "body")
	if bodySlice.Length() == 0 {
		body, err := d.CreateElementNode("body")
		if err != nil {
			return err
		}
		d.body = body
		d.Call("appendChild", body)
	} else {
		jBody := bodySlice.Index(0)
		d.body = jBody
	}
	return nil
}

func (d *Document) CreateElementNode(tag string) (js.Value, error) {
	v := d.jDoc.Call("createElement", tag)
	if !isJSValueValid(v) {
		return js.Undefined(), errors.Errorf("doc-createElement returned invalid js-value (%s)", v.Type().String())
	}
	return v, nil
}

func (d *Document) GetGlobal(names ...string) js.Value {
	curr := d.glb
	for _, name := range names {
		curr = curr.Get(name)
		if !isJSValueValid(curr) {
			return curr
		}
	}
	return curr
}

func (d *Document) Location() *url.URL {
	raw := d.GetGlobal("window", "location", "href").String()
	url, _ := url.Parse(raw)
	return url
}

//Callbacks
func (d *Document) Callback(eventType string, elt *Elt, cb ElementCallback) {
	d.registerEvent(eventType)
	elt.callback(eventType, cb)
}

func (d *Document) registerEvent(eventType string) {
	if _, contains := d.events[eventType]; contains {
		return
	}
	wlog.Infof("doc: register-event (%s)", eventType)
	d.Call("addEventListener", eventType, js.FuncOf(func(this js.Value, vals []js.Value) interface{} {
		evt, err := NewEvent(d, eventType, this, vals)
		if err != nil {
			return err
		}
		_, stack, ok := d.root.stackToTarget(evt.Target())
		if !ok || len(stack) == 0 {
			return nil
		}
		for i := len(stack) - 1; i >= 0; i-- {
			if cb, ok := stack[i].findCallback(eventType); ok {
				go func() {
					defer d.SignalRender()
					cb(evt)
				}()
				return nil
			}
		}
		return nil
	}))
	d.events[eventType] = struct{}{}
}

func (d *Document) SignalRender() {
	d.renderC <- struct{}{}
}

func (d *Document) Run(root *Elt) {
	d.root = root
	d.root.mount(d, d.body)
	for {
		select {
		case <-d.renderC:
			d.render(d.root, d.body)
			d.root.rendered()
		}
	}
}

func (d *Document) render(e *Elt, parent js.Value) {
	if e.modified {
		e.mount(d, parent)
		return
	}
	for _, c := range e.Childs {
		d.render(c, e.Value)
	}
}
