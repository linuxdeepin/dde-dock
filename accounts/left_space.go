/**
 * Copyright (C) 2017 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package accounts

import (
	"os"
	"pkg.deepin.io/lib/utils"
)

const (
	// 50MB
	filesystemMinFreeSize = 50 * 1024 * 1024
)

// checkLeftSpace Check disk left space, if no space left, then remove '~/.cache/deepin'
func (u *User) checkLeftSpace() {
	info, err := utils.QueryFilesytemInfo(u.HomeDir)
	if err != nil {
		logger.Warning("Failed to get filesystem info:", err, u.HomeDir)
		return
	}
	logger.Debugf("--------User '%s' left space: %#v", u.UserName, info)
	if info.AvailSize > filesystemMinFreeSize {
		return
	}

	logger.Debug("No space left, will remove deepin cache")
	u.removeCache()
}

func (u *User) removeCache() {
	var file = u.HomeDir + "/.cache"
	logger.Debug("-------Will remove:", file)
	err := os.RemoveAll(file)
	if err != nil {
		logger.Warning("Failed to remove cache:", err)
	}
}
