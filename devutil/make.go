package devutil

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Manifest struct {
	MainGo    string `json:"main"`
	AssetsDir string `json:"assets,omitempty"`
	CSSDir    string `json:"css,omitempty"`
}

func loadManifest(dir string) (Manifest, error) {
	file := filepath.Join(dir, "build.manifest.json")
	f, err := os.Open(file)
	if err != nil {
		return Manifest{}, errors.Wrapf(err, "open (%s)", file)
	}
	defer f.Close()

	var m Manifest
	err = json.NewDecoder(f).Decode(&m)
	if err != nil {
		return Manifest{}, errors.Wrapf(err, "decode (%s)", file)
	}
	return m, nil
}

func Make(srcDir, distDir string) error {
	//scan manifest
	manifest, err := loadManifest(srcDir)
	if err != nil {
		return errors.Wrap(err, "load manifest")
	}

	//
	err = os.MkdirAll(distDir, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "create dist-dir (%s)", distDir)
	}

	//copy wasm_exec.js and index.html
	indexSource := filepath.Join(srcDir, "index.html")
	wasmExecSource := filepath.Join(srcDir, "wasm_exec.js")
	indexDest := filepath.Join(distDir, "index.html")
	wasmExecDest := filepath.Join(distDir, "wasm_exec.js")

	err = CopyFile(indexSource, indexDest)
	if err != nil {
		return errors.Wrapf(err, "copy (%s)", indexSource)
	}
	err = CopyFile(wasmExecSource, wasmExecDest)
	if err != nil {
		return errors.Wrapf(err, "copy (%s)", wasmExecSource)
	}

	if manifest.AssetsDir != "" {
		assetsSource := filepath.Join(srcDir, manifest.AssetsDir)
		assetsDest := filepath.Join(distDir, manifest.AssetsDir)
		err = CopyDirectory(assetsSource, assetsDest)
		if err != nil {
			return errors.Wrapf(err, "copy (%s)", assetsSource)
		}
	}

	if manifest.CSSDir != "" {
		cssSource := filepath.Join(srcDir, manifest.CSSDir)
		cssDest := filepath.Join(distDir, manifest.CSSDir)
		err = CopyDirectory(cssSource, cssDest)
		if err != nil {
			return errors.Wrapf(err, "copy (%s)", cssSource)
		}
	}

	//build
	lib := filepath.Join(distDir, "lib.wasm")
	err = BuildWASM(manifest.MainGo, lib)
	if err != nil {
		return errors.Wrapf(err, "build (%s)", manifest.MainGo)
	}

	return nil
}
