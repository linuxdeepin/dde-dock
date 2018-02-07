/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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
