/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package appearance

import (
	"path"
	"pkg.deepin.io/dde/daemon/appearance/background"
	"pkg.deepin.io/lib/utils"
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
