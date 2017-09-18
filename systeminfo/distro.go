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
