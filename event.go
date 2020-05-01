package wasa

import (
	"syscall/js"

	"github.com/mazzegi/wasa/errors"
)

type Event struct {
	doc       *Document
	eventType string
	this      js.Value
	jsEvent   js.Value
	args      []js.Value
}

func NewEvent(doc *Document, eventType string, this js.Value, vals []js.Value) (*Event, error) {
	if len(vals) == 0 {
		return nil, errors.Errorf("js-event without args")
	}
	return &Event{
		doc:       doc,
		eventType: eventType,
		this:      this,
		jsEvent:   vals[0],
		args:      vals[1:],
	}, nil
}

func (e *Event) Document() *Document {
	return e.doc
}

func (e *Event) Type() string {
	return e.eventType
}

func (e *Event) This() js.Value {
	return e.this
}

func (e *Event) Target() js.Value {
	return e.jsEvent.Get("target")
}

func (e *Event) TargetID() string {
	return e.jsEvent.Get("target").Get("id").String()
}

func (e *Event) JSEvent() js.Value {
	return e.jsEvent
}

func (e *Event) Args() []js.Value {
	return e.args
}
