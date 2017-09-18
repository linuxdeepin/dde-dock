/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package systeminfo

import (
	"fmt"
	"pkg.deepin.io/lib/keyfile"
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
	kfile := keyfile.NewKeyFile()
	if err := kfile.LoadFromFile(file); err != nil {
		return "", err
	}

	version, err := kfile.GetString(versionGroupRelease,
		versionKeyVersion)
	if err != nil {
		return "", err
	}

	ty, _ := kfile.GetLocaleString(versionGroupRelease,
		versionKeyType, "")
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
