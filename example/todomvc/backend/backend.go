package backend

import "sync"

const (
	All       = "all"
	Active    = "active"
	Completed = "completed"
)

type Item struct {
	ID        int
	Text      string
	Completed bool
}

type Backend struct {
	sync.RWMutex
	items      []*Item
	nextID     int
	subs       []func()
	currFilter string
}

func New() *Backend {
	b := &Backend{
		items:      []*Item{},
		nextID:     1,
		currFilter: All,
	}
	return b
}

func (b *Backend) Subscribe(cb func()) {
	b.Lock()
	defer b.Unlock()
	b.subs = append(b.subs, cb)
}

func (b *Backend) publish() {
	for _, cb := range b.subs {
		cb()
	}
}

func (b *Backend) Filter() string {
	b.RLock()
	defer b.RUnlock()
	return b.currFilter
}

func (b *Backend) ChangeFilter(f string) {
	defer b.publish()
	b.Lock()
	defer b.Unlock()
	b.currFilter = f
}

func (b *Backend) Add(text string) {
	defer b.publish()
	b.Lock()
	defer b.Unlock()
	b.items = append(b.items, &Item{
		ID:        b.nextID,
		Text:      text,
		Completed: false,
	})
	b.nextID++

}

func (b *Backend) IsEmpty() bool {
	b.RLock()
	defer b.RUnlock()
	return len(b.items) == 0
}

func (b *Backend) Count() int {
	b.RLock()
	defer b.RUnlock()
	return len(b.items)
}

func (b *Backend) Each(fnc func(i *Item)) {
	b.Lock()
	defer b.Unlock()
	for _, i := range b.items {
		inc := true
		switch {
		case b.currFilter == Active:
			inc = !i.Completed
		case b.currFilter == Completed:
			inc = i.Completed
		}
		if inc {
			fnc(i)
		}
	}
}

func (b *Backend) Complete(id int) {
	defer b.publish()
	b.Lock()
	defer b.Unlock()
	for _, i := range b.items {
		if i.ID == id {
			i.Completed = true
			return
		}
	}
}

func (b *Backend) ToggleComplete(id int) {
	defer b.publish()
	b.Lock()
	defer b.Unlock()
	for _, i := range b.items {
		if i.ID == id {
			i.Completed = !i.Completed
			return
		}
	}
}

func (b *Backend) ToggleAll() {
	defer b.publish()
	b.Lock()
	defer b.Unlock()
	for _, i := range b.items {
		i.Completed = !i.Completed
	}
}

func (b *Backend) DeleteCompleted() {
	defer b.publish()
	b.Lock()
	defer b.Unlock()
	newItems := []*Item{}
	for _, i := range b.items {
		if !i.Completed {
			newItems = append(newItems, i)
		}
	}
	b.items = newItems
}

func (b *Backend) Delete(id int) {
	defer b.publish()
	b.Lock()
	defer b.Unlock()
	for idx, i := range b.items {
		if i.ID == id {
			b.items = append(b.items[:idx], b.items[idx+1:]...)
			return
		}
	}
}
