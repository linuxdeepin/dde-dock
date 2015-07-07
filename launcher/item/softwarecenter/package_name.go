package softwarecenter

import (
	. "pkg.deepin.io/dde-daemon/launcher/interfaces"
)

func GetPkgName(soft SoftwareCenterInterface, path string) (string, error) {
	return soft.GetPkgNameFromPath(path)
}
