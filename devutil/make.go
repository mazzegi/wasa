package devutil

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

func Make(dist, wasmExec, mainGo string) error {
	err := os.MkdirAll(dist, os.ModePerm)
	if err != nil {
		return err
	}
	bs, err := exec.Command("cp", wasmExec, dist).CombinedOutput()
	if err != nil {
		return errors.Errorf("cp wasm_exec.js: %s", string(bs))
	}
	fmt.Printf("copied wasm_exec.js (%s) -> (%s)\n", wasmExec, dist)

	fIndex, err := os.Create(filepath.Join(dist, "index.html"))
	if err != nil {
		return errors.Wrap(err, "create index.html")
	}
	defer fIndex.Close()

	t, err := template.New("index.html").Parse(indexHTMLTemplate)
	if err != nil {
		return errors.Wrap(err, "parse index.html template")
	}
	data := struct {
	}{}
	err = t.Execute(fIndex, data)
	if err != nil {
		return errors.Wrap(err, "execute index.html template")
	}
	fmt.Printf("created index.html\n")

	wasmPath := filepath.Join(dist, wasmFile)
	err = buildWASM(mainGo, wasmPath)
	if err != nil {
		return errors.Wrap(err, "build wasm")
	}
	fmt.Printf("build wasm\n")

	return nil
}
