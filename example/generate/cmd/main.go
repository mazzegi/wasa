package main

import (
	"github.com/mazzegi/wasa/example/generate/app"
	"github.com/mazzegi/wasa/wlog"
)

func main() {
	wlog.InstallConsoleLogger()
	a, err := app.New()
	if err != nil {
		panic(err)
	}
	a.Run()
}
