package wasa

import (
	"net/url"
	"syscall/js"
	"time"

	"github.com/mazzegi/wasa/errors"
	"github.com/mazzegi/wasa/timing"
	"github.com/mazzegi/wasa/wlog"
)

type Document struct {
	jsElt
	glb         js.Value
	jDoc        js.Value
	body        jsElt
	events      map[string]struct{}
	root        *Elt
	focus       *Elt
	renderC     chan struct{}
	afterRender []func()
	afterMount  []func()
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

func (d *Document) Location() *url.URL {
	raw := d.GetGlobal("window", "location", "href").String()
	url, _ := url.Parse(raw)
	return url
}

func (d *Document) Focus(elt *Elt) {
	d.focus = elt
}

func (d *Document) BodyDimensions() (int, int) {
	w := d.body.get("clientWidth").Float()
	h := d.body.get("clientHeight").Float()
	return int(w), int(h)
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
	d.addEventListener(eventType, js.FuncOf(func(this js.Value, vals []js.Value) interface{} {
		evt, err := NewEvent(d, eventType, this, vals)
		if err != nil {
			return err
		}
		wlog.Infof("doc:on: (%s) -> (%s)", eventType, evt.TargetID())
		wlog.Raw("target:", evt.Target())
		start := time.Now()
		if _, stack, ok := d.root.findByTarget(evt.Target()); ok {
			if len(stack) > 0 {
				for i := len(stack) - 1; i >= 0; i-- {
					if cb, ok := stack[i].findCallback(eventType); ok {
						go func() {
							wlog.Infof("doc:on: (%s) -> (%s). found in (%s)", eventType, evt.TargetID(), time.Since(start))
							defer d.SignalRender()
							cb(evt)
						}()
						return nil
					}
				}
			}
		}
		return nil
	}))
	d.events[eventType] = struct{}{}
}

func (d *Document) AfterRender(cb func()) {
	d.afterRender = append(d.afterRender, cb)
}

func (d *Document) AfterMount(cb func()) {
	d.afterMount = append(d.afterMount, cb)
}

func (d *Document) SignalRender() {
	d.renderC <- struct{}{}
}

func (d *Document) Run(root *Elt) {
	wlog.Infof("doc: enter render loop ...")
	d.root = root
	d.root.mount(d, d.body)
	for _, cb := range d.afterMount {
		cb()
	}
	for _, cb := range d.afterRender {
		cb()
	}
	for {
		select {
		case <-d.renderC:
			t := timing.New("doc-render")
			t.Log("render")
			d.render(d.root, d.body)
			t.Log("after-render")
			if d.focus != nil {
				d.focus.Call("focus")
				t.Log("after-focus")
				d.focus = nil
			}
			for _, cb := range d.afterRender {
				cb()
			}
			t.Log("after-render-callbacks")
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
