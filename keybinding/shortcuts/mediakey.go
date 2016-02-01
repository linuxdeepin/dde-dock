/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package shortcuts

func ListMediaShortcut() Shortcuts {
	s := newMediakeyGSetting()
	defer s.Unref()
	return doListShortcut(s, mediaIdNameMap(), KeyTypeMedia)
}

func resetMediaAccels() {
	s := newMediakeyGSetting()
	defer s.Unref()
	doResetAccels(s)
}

func disableMediaAccels(key string) {
	s := newMediakeyGSetting()
	defer s.Unref()
	doDisableAccles(s, key)
}

func addMediaAccel(key, accel string) {
	s := newMediakeyGSetting()
	defer s.Unref()
	doAddAccel(s, key, accel)
}

func delMediaAccel(key, accel string) {
	s := newMediakeyGSetting()
	defer s.Unref()
	doDelAccel(s, key, accel)
}

func mediaIdNameMap() map[string]string {
	var idMap = map[string]string{
		"calculator":  "Calculator",
		"eject":       "Eject",
		"email":       "Email client",
		"www":         "Web broswer",
		"media":       "Media player",
		"play":        "Play/Pause",
		"pause":       "Pause",
		"stop":        "Stop",
		"previous":    "Previous",
		"next":        "Next",
		"volume-mute": "Mute",
		"volume-down": "Volume down",
		"volume-up":   "Volume up",
	}
	return idMap
}
