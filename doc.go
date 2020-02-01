package wasa

import (
	"log"
	"syscall/js"
	"time"

	"github.com/pkg/errors"
)

type Document struct {
	jsElt
	glb     js.Value
	jDoc    js.Value
	body    jsElt
	events  map[string]struct{}
	root    *Elt
	focus   *Elt
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
		jsElt:   newJSElt(jDoc),
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
		d.appendChild(body)
	} else {
		jBody := bodySlice.Index(0)
		d.body = newJSElt(jBody)
	}
	return nil
}

func (d *Document) CreateElementNode(tag string) (jsElt, error) {
	v := d.jDoc.Call("createElement", tag)
	if !isJSValueValid(v) {
		return undefinedJSElt(), errors.Errorf("doc-createElement returned invalid js-value (%s)", v.Type().String())
	}
	return newJSElt(v), nil
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

func (d *Document) Focus(elt *Elt) {
	d.focus = elt
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
	log.Printf("doc: register-event (%s)", eventType)
	d.addEventListener(eventType, js.FuncOf(func(this js.Value, vals []js.Value) interface{} {
		evt, err := NewEvent(d, eventType, this, vals)
		if err != nil {
			return err
		}
		log.Printf("doc:on: (%s) -> (%s)", eventType, evt.TargetID())
		start := time.Now()
		if elt, ok := d.root.findByTarget(evt.Target()); ok {
			if cb, ok := elt.findCallback(eventType); ok {
				go func() {
					log.Printf("doc:on: (%s) -> (%s). found in (%s)", eventType, evt.TargetID(), time.Since(start))
					defer d.signalRender()
					cb(evt)
				}()
			}
		}

		return nil
	}))
	d.events[eventType] = struct{}{}
}

func (d *Document) signalRender() {
	d.renderC <- struct{}{}
}

func (d *Document) Run(root *Elt) {
	log.Printf("doc: enter render loop ...")
	d.root = root
	d.root.mount(d, d.body)
	for {
		select {
		case <-d.renderC:
			log.Printf("render ...")
			d.render(d.root, d.body)
			if d.focus != nil {
				d.focus.Call("focus")
				d.focus = nil
			}
		}
	}
}

func (d *Document) render(e *Elt, parent jsElt) {
	if e.modified {
		e.mount(d, parent)
		return
	}
	for _, c := range e.Childs {
		d.render(c, e.jsElt)
	}
}
