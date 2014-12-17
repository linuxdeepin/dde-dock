package utils

import (
	"os"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"strings"
)

const (
	DirDefaultPerm os.FileMode = 0755
)

func CreateDesktopAppInfo(name string) *gio.DesktopAppInfo {
	if strings.ContainsRune(name, os.PathSeparator) {
		return gio.NewDesktopAppInfoFromFilename(name)
	} else {
		return gio.NewDesktopAppInfo(name)
	}
}
