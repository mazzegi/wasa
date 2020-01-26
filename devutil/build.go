package devutil

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

const wasmFile = "lib.wasm"

func buildWASM(mainGo string, out string) error {
	mainGo, _ = filepath.Abs(mainGo)
	out, _ = filepath.Abs(out)

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(wd)
	mainDir := filepath.Dir(mainGo)
	mainFile := filepath.Base(mainGo)
	os.Chdir(mainDir)

	start := time.Now()
	cmd := exec.Command("go", "build", "-o", out, mainFile)
	cmd.Env = append(os.Environ(),
		"GOARCH=wasm",
		"GOOS=js",
	)
	stdErr, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "go build: %s", string(stdErr))
	}
	fmt.Printf("build (%s) -> (%s) done in (%s)\n", mainGo, out, time.Since(start))
	return nil
}
