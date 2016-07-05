/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
