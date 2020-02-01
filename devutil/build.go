package devutil

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

func BuildWASM(mainGo string, out string) error {
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

	cmd := exec.Command("go", "build", "-o", out, mainFile)
	cmd.Env = append(os.Environ(),
		"GOARCH=wasm",
		"GOOS=js",
	)
	stdErr, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "go build: %s", string(stdErr))
	}
	return nil
}
