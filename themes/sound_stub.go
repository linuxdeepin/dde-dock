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
// #include <stdlib.h>
// #include "canberra_wrapper.h"
import "C"
import "unsafe"

import (
        "dlib/dbus"
)

const (
        SOUND_DEST = "com.deepin.daemon.Themes"
        SOUND_PATH = "/com/deepin/daemon/Sound"
        SOUND_IFC  = "com.deepin.daemon.Sound"
)

const (
        SOUND_THEME_PATH      = "/usr/share/sounds/"
        SOUND_THEME_MAIN_FILE = "index.theme"
)

type Sound struct{}

func (s *Sound) GetDBusInfo() dbus.DBusInfo {
        return dbus.DBusInfo{
                SOUND_DEST,
                SOUND_PATH,
                SOUND_IFC,
        }
}

func (s *Sound) PlaySystemSound(event string) (err error) {
        return s.PlayThemeSound(s.getCurrentSoundTheme(), event)
}

func (s *Sound) getCurrentSoundTheme() string {
        currentSoundTheme := "LinuxDeepin" // default theme name
        if objManager != nil {
                if objTheme := objManager.getThemeObject(objManager.CurrentTheme); objTheme != nil {
                        currentSoundTheme = objTheme.SoundThemeName
                }
        }
        return currentSoundTheme
}

func (s *Sound) PlayThemeSound(theme, event string) (err error) {
        go func() {
                ctheme := C.CString(theme)
                defer C.free(unsafe.Pointer(ctheme))
                cevent := C.CString(event)
                defer C.free(unsafe.Pointer(cevent))
                ret := C.canberra_play_system_sound(ctheme, cevent)
                if ret != 0 {
                        logObject.Error("play system sound failed: theme=%s, event=%s, %s",
                                theme, event, C.GoString(C.ca_strerror(ret)))
                }
        }()
        return
}

func (s *Sound) PlaySoundFile(file string) (err error) {
        go func() {
                cfile := C.CString(file)
                defer C.free(unsafe.Pointer(cfile))
                ret := C.canberra_play_sound_file(cfile)
                if ret != 0 {
                        logObject.Error("play sound file failed: %s, %s", file, C.GoString(C.ca_strerror(ret)))
                }
        }()
        return
}
