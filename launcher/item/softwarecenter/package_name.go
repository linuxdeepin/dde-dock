package softwarecenter

import (
	"os/exec"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/interfaces"
	"strings"
)

func GetPkgName(soft SoftwareCenterInterface, path string) (string, error) {
	pkgName, err := soft.GetPkgNameFromPath(path)
	if err != nil {
		return getPkgNameFromCommandLine(path)
	}

	return pkgName, nil
}

func getPkgNameFromCommandLine(path string) (string, error) {
	cmd := exec.Command("dpkg", "-S", path)
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	content, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.Split(string(content), ":")[0], nil
}
