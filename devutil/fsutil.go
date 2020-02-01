package devutil

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

func ExistsFolder(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	if !fi.IsDir() {
		return false
	}
	return true
}

func ExistsFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return false
	}
	return true
}

func CopyFile(source string, target string) error {
	err := os.MkdirAll(path.Dir(target), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "mkdirall (%s)", path.Dir(target))
	}
	sF, err := os.Open(source)
	if err != nil {
		return errors.Wrapf(err, "open source (%s)", source)
	}
	defer sF.Close()
	tF, err := os.Create(target)
	if err != nil {
		return errors.Wrapf(err, "create target (%s)", target)
	}
	defer tF.Close()
	_, err = io.Copy(tF, sF)
	if err != nil {
		return errors.Wrap(err, "copy")
	}
	//copy permissions/mode
	fiS, err := os.Stat(source)
	if err != nil {
		return errors.Wrapf(err, "stat (%s)", source)
	}
	err = os.Chmod(target, fiS.Mode())
	return err
}

func CopyDirectory(scrDir, dest string) error {
	entries, err := ioutil.ReadDir(scrDir)
	if err != nil {
		return errors.Wrapf(err, "read-dir (%s)", scrDir)
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(scrDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return errors.Wrapf(err, "stat (%s)", sourcePath)
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := CreateDirIfNotExists(destPath, 0755); err != nil {
				return err
			}
			if err := CopyDirectory(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if err := CopyFile(sourcePath, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateDirIfNotExists(dir string, perm os.FileMode) error {
	if ExistsFolder(dir) {
		return nil
	}
	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}
	return nil
}
