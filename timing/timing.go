package timing

import (
	"time"

	"github.com/mazzegi/wasa/wlog"
)

type Timing struct {
	start time.Time
	topic string
}

func New(topic string) *Timing {
	return &Timing{
		start: time.Now(),
		topic: topic,
	}
}

func (t *Timing) Reset() {
	t.start = time.Now()
}

func (t *Timing) Log(label string) {
	wlog.Infof("[timing] [%s] <%s>:[%s] (since %s)", t.topic, label, time.Now().Round(1*time.Microsecond).Format(time.RFC3339Nano), time.Since(t.start))
}
