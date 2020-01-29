// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mazzegi/wasa/devutil"
)

var (
	dist     = flag.String("dist", "dist", "dist directory to serve")
	wasmExec = flag.String("wasmexec", "../../wasm_exec.js", "path to wasm_exec.js")
	mainGo   = flag.String("maingo", "../test/main.go", "location of main.go file to build and serve")
)

func main() {
	flag.Parse()
	fmt.Printf("make ...\n")
	fmt.Printf("dist:     %s\n", *dist)
	fmt.Printf("wasmExec: %s\n", *wasmExec)
	fmt.Printf("main.go:  %s\n", *mainGo)
	err := devutil.Make(*dist, *wasmExec, *mainGo)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("make succeeded\n")
}
