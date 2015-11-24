package utils

import (
	"io/ioutil"
	"os"
	"pkg.deepin.io/lib/glib-2.0"
)

// SaveKeyFile saves key file.
func SaveKeyFile(file *glib.KeyFile, path string) error {
	_, content, err := file.ToData()
	if err != nil {
		return err
	}

	stat, err := os.Lstat(path)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, []byte(content), stat.Mode())
	if err != nil {
		return err
	}
	return nil
}
