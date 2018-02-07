/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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

package accounts

import (
	dutils "pkg.deepin.io/lib/utils"
)

const (
	versionFile = "/etc/deepin-version"
)

func getDeepinReleaseType() string {
	keyFile, err := dutils.NewKeyFileFromFile(versionFile)
	if err != nil {
		return ""
	}
	defer keyFile.Free()
	releaseType, err := keyFile.GetString("Release", "Type")
	return releaseType
}
