package utils

import (
	"os"
	"pkg.deepin.io/lib/gio-2.0"
	"strings"
)

// dir default perm.
const (
	DirDefaultPerm os.FileMode = 0755
)

// CreateDesktopAppInfo is a helper function for creating GDesktopAppInfo object.
// if name is a path, gio.NewDesktopAppInfoFromFilename is used.
// otherwise, name must be desktop id and gio.NewDesktopAppInfo is used.
func CreateDesktopAppInfo(name string) *gio.DesktopAppInfo {
	if strings.ContainsRune(name, os.PathSeparator) {
		return gio.NewDesktopAppInfoFromFilename(name)
	} else {
		return gio.NewDesktopAppInfo(name)
	}
}
