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

package inputdevices

import (
	dutils "pkg.deepin.io/lib/utils"
)

const (
	wmKbdSchemaID    = "com.deepin.wrap.gnome.desktop.peripherals.keyboard"
	wmKbdKeyRepeat   = "repeat"
	wmKbdKeyDelay    = "delay"
	wmKbdKeyInterval = "repeat-interval"

	wmTPadSchemaID         = "com.deepin.wrap.gnome.desktop.peripherals.touchpad"
	wmTPadKeyEdgeScroll    = "edge-scrolling-enabled"
	wmTPadKeyNaturalScroll = "natural-scroll"
	wmTPadKeyTapClick      = "tap-to-click"
	// enum: mouse, left, right
	wmTPadKeyLeftHanded = "left-handed"

	wmMouseSchemaID = "com.deepin.wrap.gnome.desktop.peripherals.mouse"
)

func setWMKeyboardRepeat(repeat bool, delay, interval uint32) {
	setting, err := dutils.CheckAndNewGSettings(wmKbdSchemaID)
	if err != nil {
		logger.Warning("Failed to new wm keyboard settings")
		return
	}
	defer setting.Unref()

	if setting.GetBoolean(wmKbdKeyRepeat) != repeat {
		setting.SetBoolean(wmKbdKeyRepeat, repeat)
	}
	if setting.GetUint(wmKbdKeyDelay) != delay {
		setting.SetUint(wmKbdKeyDelay, delay)
	}
	if setting.GetUint(wmKbdKeyInterval) != interval {
		setting.SetUint(wmKbdKeyInterval, interval)
	}
}

func setWMTPadBoolKey(key string, value bool) {
	setting, err := dutils.CheckAndNewGSettings(wmTPadSchemaID)
	if err != nil {
		logger.Warning("Failed to new wm touchpad settings")
		return
	}
	defer setting.Unref()

	switch key {
	case wmTPadKeyEdgeScroll, wmTPadKeyNaturalScroll, wmTPadKeyTapClick:
		if v := setting.GetBoolean(key); v == value {
			return
		}
		setting.SetBoolean(key, value)
		//case wmTPadKeyLeftHanded:
		//v := setting.GetString(key)
		//var tmp = "right"
		//if value {
		//tmp = "left"
		//}
		//if v == tmp {
		//return
		//}
		//setting.SetString(key, tmp)
	}
}

func setWMMouseBoolKey(key string, value bool) {
	return
	setting, err := dutils.CheckAndNewGSettings(wmMouseSchemaID)
	if err != nil {
		logger.Warning("Failed to new wm mouse settings")
		return
	}
	defer setting.Unref()

	switch key {
	case wmTPadKeyNaturalScroll, wmTPadKeyLeftHanded:
		if v := setting.GetBoolean(key); v == value {
			return
		}
		setting.SetBoolean(key, value)
	}
}
