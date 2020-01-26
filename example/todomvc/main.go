package main

import "github.com/mazzegi/wasa/example/todomvc/app"

func main() {
	a, err := app.New()
	if err != nil {
		panic(err)
	}
	a.Run()
}
