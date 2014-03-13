/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

// #cgo pkg-config: glib-2.0 libcanberra
// #include "canberra_wrapper.h"
import "C"

import (
	"dlib/dbus"
)

const (
	SOUND_DEST = "com.deepin.daemon.Themes"
	SOUND_PATH = "/com/deepin/daemon/Sound"
	SOUND_IFC  = "com.deepin.daemon.Sound"
)

func (s *Sound) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		SOUND_DEST,
		SOUND_PATH,
		SOUND_IFC,
	}
}

// FIXME PlaySystemSound() maybe should be place in dde-api
func (s *Sound) PlaySystemSound(event string) (err error) {
	currentTheme := "LinuxDeepin" // TODO
	return s.PlayThemeSystemSound(currentTheme, event)
}

func (s *Sound) PlayThemeSystemSound(theme, event string) (err error) {
	go func() {
		ret := C.canberra_play_system_sound(C.CString(theme), C.CString(event))
		if ret != 0 {
			logObject.Error("play system sound failed: theme=%s, event=%s, %s",
				theme, event, C.GoString(C.ca_strerror(ret)))
		}
	}()
	return
}

// TODO remove
func (s *Sound) PlaySoundFile(file string) (err error) {
	go func() {
		ret := C.canberra_play_sound_file(C.CString(file))
		if ret != 0 {
			logObject.Error("play sound file failed: %s, %s\n", file, C.GoString(C.ca_strerror(ret)))
		}
	}()
	return
}
