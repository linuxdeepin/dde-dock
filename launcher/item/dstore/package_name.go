package dstore

import (
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
)

// GetPkgName returns package name of given desktop file.
func GetPkgName(soft DStore, path string) (string, error) {
	return soft.GetPkgNameFromPath(path)
}
