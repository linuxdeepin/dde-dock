package systeminfo

import (
	"fmt"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	versionFileDeepin = "/etc/deepin-version"
	versionFileLSB    = "/etc/lsb-release"

	versionGroupRelease = "Release"
	versionKeyVersion   = "Version"
	versionKeyType      = "Type"

	versionGroupAddition = "Addition"
	versionKeyMilestone  = "Milestone"

	versionKeyLSB   = "DISTRIB_RELEASE"
	versionKeyDelim = "="
)

func getVersion() (string, error) {
	version, err := getVersionFromDeepin(versionFileDeepin)
	if err == nil {
		return version, nil
	}

	version, err = getVersionFromLSB(versionFileLSB)
	if err == nil {
		return version, nil
	}

	return "", err
}

func getVersionFromDeepin(file string) (string, error) {
	kfile, err := dutils.NewKeyFileFromFile(file)
	if err != nil {
		return "", err
	}
	defer kfile.Free()

	version, err := kfile.GetString(versionGroupRelease,
		versionKeyVersion)
	if err != nil {
		return "", err
	}

	ty, _ := kfile.GetLocaleString(versionGroupRelease,
		versionKeyType, "\x00")
	if len(ty) != 0 {
		version = version + " " + ty
	}

	milestone, _ := kfile.GetString(versionGroupAddition,
		versionKeyMilestone)
	if len(milestone) != 0 {
		version = version + " " + milestone
	}

	return version, nil
}

func getVersionFromLSB(file string) (string, error) {
	ret, err := parseInfoFile(file, versionKeyDelim)
	if err != nil {
		return "", err
	}

	value, ok := ret[versionKeyLSB]
	if !ok {
		return "", fmt.Errorf("Can not find the key '%s'", versionKeyLSB)
	}

	return value, nil
}
