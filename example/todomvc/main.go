package main

import (
	"github.com/mazzegi/wasa/example/todomvc/app"
	"github.com/mazzegi/wasa/example/todomvc/backend"
)

func main() {
	be := backend.New()
	a, err := app.New(be)
	if err != nil {
		panic(err)
	}
	a.Run()
}
