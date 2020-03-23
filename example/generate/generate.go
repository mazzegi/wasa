// +build ignore

package main

import (
	"log"

	"github.com/mazzegi/wasa/devutil"
)

func main() {
	g, err := devutil.NewGenerator()
	if err != nil {
		log.Fatalf("new-generator: %v", err)
	}
	err = g.ProcessDir("app")
	if err != nil {
		log.Fatalf("process-dir: %v", err)
	}
}
