// +build ignore

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mazzegi/wasa/devutil"
)

func main() {
	handler, err := devutil.NewMakeServeHandler("", "dist", "lib.wasm")
	if err != nil {
		fmt.Println("ERROR: init make-serve-handler:", err)
		os.Exit(1)
	}
	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)
}
