/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
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

package appearance

import (
	"path"
	"pkg.deepin.io/dde/daemon/appearance/background"
	"pkg.deepin.io/lib/utils"
	"time"
)

const (
	deepinBackgroundPrefix = "file:///usr/share/wallpapers/deepin/"
)

func (m *Manager) initBackground() {
	bg := correctBackgroundPath(m.Background.Get())
	if len(bg) == 0 {
		m.wrapBgSetting.Reset(gsKeyBackground)
	} else if bg != m.Background.Get() {
		m.wrapBgSetting.SetString(gsKeyBackground, bg)
	} else if !checkBlurredBackgroundExists(bg) {
		// If the corresponding blurred image doesn't exist, set the background
		// again to trigger the blur process.

		// This function is executed before gsettings signals are correctly
		// connected, so we need some time.
		time.AfterFunc(5*time.Second, func() {
			m.wrapBgSetting.SetString(gsKeyBackground, bg)
		})
	}
}

// correctBackgroundPath the bg path has changed because of deleting bg from deepin-artwork-themes
func correctBackgroundPath(bg string) string {
	if utils.IsFileExist(bg) {
		return bg
	}

	uri := deepinBackgroundPrefix + path.Base(bg)
	if background.ListBackground().Get(uri) != nil {
		return uri
	}
	return ""
}

func checkBlurredBackgroundExists(srcURI string) bool {
	id, _ := utils.SumStrMd5(utils.DecodeURI(srcURI))
	blurredImagePath := "/var/cache/image-blur/" + id + path.Ext(srcURI)

	return utils.IsFileExist(blurredImagePath)
}
