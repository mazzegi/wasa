package wasa

import (
	"log"
	"syscall/js"
	"time"

	"github.com/pkg/errors"
)

type Document struct {
	jsElt
	jDoc    js.Value
	body    jsElt
	events  map[string]struct{}
	root    *Elt
	renderC chan struct{}
}

func NewDocument() (*Document, error) {
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
		jDoc:    jDoc,
		events:  map[string]struct{}{},
		renderC: make(chan struct{}),
	}
	err := doc.createBodyIfNotExists()
	if err != nil {
		return nil, errors.Wrap(err, "create-body")
	}
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

//Callbacks
func (d *Document) Callback(event string, elt *Elt, cb ElementCallback) {
	d.registerEvent(event)
	elt.callback(event, cb)
}

func (d *Document) registerEvent(event string) {
	if _, contains := d.events[event]; contains {
		return
	}
	log.Printf("doc: register-event (%s)", event)
	d.addEventListener(event, js.FuncOf(func(this js.Value, vals []js.Value) interface{} {
		if len(vals) == 0 {
			return errors.Errorf("callback without args")
		}
		target := vals[0].Get("target")
		if !isJSValueValid(target) {
			return errors.Errorf("invalid target")
		}
		log.Printf("doc:on: (%s) -> (%s)", event, target.Get("id").String())
		start := time.Now()
		if elt, ok := d.root.findByTarget(target); ok {
			if cb, ok := elt.findCallback(event); ok {
				go func() {
					log.Printf("doc:on: (%s) -> (%s). found in (%s)", event, target.Get("id").String(), time.Since(start))
					defer d.signalRender()
					cb(d, target, vals[1:])
				}()
			}
		}

		return nil
	}))
	d.events[event] = struct{}{}
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
