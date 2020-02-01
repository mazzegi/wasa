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

	cmd := exec.Command("go", "build", "-o", out, mainGo)
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
