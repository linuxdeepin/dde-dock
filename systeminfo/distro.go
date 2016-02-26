/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package systeminfo

import (
	"fmt"
)

const (
	distroFileLSB    = "/etc/lsb-release"

	distroIdKeyLSB   = "DISTRIB_ID"
	distroDescKeyLSB = "DISTRIB_DESCRIPTION"
	distroVerKeyLSB  = "DISTRIB_RELEASE"
	distroKeyDelim   = "="
)

func getDistro() (string, string, string, error) {
	distroId, distroDesc, distroVer, err := getDistroFromLSB(distroFileLSB)
	if err == nil {
		return distroId, distroDesc, distroVer, nil
	}

	return "", "", "", err
}

func getDistroFromLSB(file string) (string, string, string, error) {
	ret, err := parseInfoFile(file, distroKeyDelim)
	if err != nil {
		return "", "", "", err
	}

	distroId, ok := ret[distroIdKeyLSB]
	if !ok {
		return "", "", "", fmt.Errorf("Cannot find the key '%s'", distroIdKeyLSB)
	}

	distroDesc, ok := ret[distroDescKeyLSB]
	if !ok {
		return "", "", "", fmt.Errorf("Cannot find the key '%s'", distroDescKeyLSB)
	}

	if distroDesc[0] == '"' && distroDesc[len(distroDesc) - 1] == '"' {
		distroDesc = distroDesc[1:len(distroDesc) - 1]
	}

	distroVer, ok := ret[distroVerKeyLSB]
	if !ok {
		return "", "", "", fmt.Errorf("Cannot find the key '%s'", distroVerKeyLSB)
	}

	return distroId, distroDesc, distroVer, nil
}
