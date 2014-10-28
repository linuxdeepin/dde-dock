package utils

import (
	"os"
	"path"
	"pkg.linuxdeepin.com/lib/glib-2.0"
	"pkg.linuxdeepin.com/lib/utils"
)

func ConfigFilePath(name string) string {
	return path.Join(glib.GetUserConfigDir(), name)
}

func ConfigFile(name string, defaultFile string) (*glib.KeyFile, error) {
	file := glib.NewKeyFile()
	conf := ConfigFilePath(name)
	if !utils.IsFileExist(conf) {
		os.MkdirAll(path.Dir(conf), DirDefaultPerm)
		if defaultFile == "" {
			f, err := os.Create(conf)
			if err != nil {
				return nil, err
			}
			defer f.Close()
		} else {
			CopyFile(defaultFile, conf, CopyFileNotKeepSymlink)
		}
	}

	if ok, err := file.LoadFromFile(conf, glib.KeyFileFlagsNone); !ok {
		file.Free()
		return nil, err
	}
	return file, nil
}

func uniqueStringList(l []string) []string {
	m := make(map[string]bool, 0)
	for _, v := range l {
		m[v] = true
	}
	n := make([]string, 0)
	for k, _ := range m {
		n = append(n, k)
	}
	return n
}
