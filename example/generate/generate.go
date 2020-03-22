// +build ignore

package main

import (
	"log"

	"github.com/mazzegi/wasa/gen"
)

func main() {
	g, err := gen.NewGenerator()
	if err != nil {
		log.Fatalf("new-generator: %v", err)
	}
	err = g.ProcessDir("app")
	if err != nil {
		log.Fatalf("process-dir: %v", err)
	}
}
