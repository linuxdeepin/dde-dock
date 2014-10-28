package utils

import (
	"io/ioutil"
	"os"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/interfaces"
)

func SaveKeyFile(file RateConfigFileInterface, path string) error {
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
