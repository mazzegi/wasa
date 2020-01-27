package app

type item struct {
	Text      string
	Completed bool
}

type repo struct {
	items []item
}

func newRepo() *repo {
	r := &repo{
		items: []item{},
	}
	r.items = append(r.items, item{
		Text:      "Buy some juicy stuff",
		Completed: false,
	}, item{
		Text:      "Date the Simpsons",
		Completed: false,
	})

	return r
}

func (r *repo) isEmpty() bool {
	return len(r.items) == 0
}
