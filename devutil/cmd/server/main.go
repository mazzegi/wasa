package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mazzegi/wasa/devutil"
)

var (
	bind     = flag.String("bind", ":8080", "bind address")
	dist     = flag.String("dist", "dist", "dist directory to serve")
	wasmExec = flag.String("wasmexec", "../../wasm_exec.js", "path to wasm_exec.js")
	mainGo   = flag.String("maingo", "../test/main.go", "location of main.go file to build and serve")
)

func main() {
	flag.Parse()
	c := devutil.ServerConfig{
		Bind:     *bind,
		DistDir:  *dist,
		MainGo:   *mainGo,
		WasmExec: *wasmExec,
	}

	sigC := make(chan os.Signal)
	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	s, err := devutil.NewServer(c)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
	go s.ListenAndServe()

	<-sigC
	s.Close()
}
