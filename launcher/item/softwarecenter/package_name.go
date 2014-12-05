package softwarecenter

import (
	. "pkg.linuxdeepin.com/dde-daemon/launcher/interfaces"
)

func GetPkgName(soft SoftwareCenterInterface, path string) (string, error) {
	return soft.GetPkgNameFromPath(path)
}
