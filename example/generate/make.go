// +build ignore

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mazzegi/wasa/devutil"
)

func main() {
	start := time.Now()
	fmt.Println("make ...")
	err := devutil.Make("", "dist")
	if err != nil {
		fmt.Println("ERROR: make:", err)
		os.Exit(1)
	}
	fmt.Println("make succeeded:", time.Since(start))
}
