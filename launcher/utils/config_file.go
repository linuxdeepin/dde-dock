package utils

// #cgo pkg-config: glib-2.0
// #include <glib.h>
import "C"
import (
	"os"
	"path"
	"pkg.deepin.io/lib/glib-2.0"
	"pkg.deepin.io/lib/utils"
)

// GReloadUserSpecialDirsCache reloads user special dirs cache.
func GReloadUserSpecialDirsCache() {
	C.g_reload_user_special_dirs_cache()
}

// ConfigFilePath returns path in user's config dir.
func ConfigFilePath(name string) string {
	return path.Join(glib.GetUserConfigDir(), name)
}

// ConfigFile open the given keyfile, this file will be created if not existed.
func ConfigFile(name string) (*glib.KeyFile, error) {
	file := glib.NewKeyFile()
	conf := ConfigFilePath(name)
	if !utils.IsFileExist(conf) {
		os.MkdirAll(path.Dir(conf), DirDefaultPerm)
		f, err := os.Create(conf)
		if err != nil {
			return nil, err
		}
		defer f.Close()
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
	var n []string
	for k := range m {
		n = append(n, k)
	}
	return n
}
