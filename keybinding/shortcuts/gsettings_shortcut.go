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

package shortcuts

import (
	"gir/gio-2.0"
)

type GSettingsShortcut struct {
	BaseShortcut
	gsettings *gio.Settings
}

func NewGSettingsShortcut(gsettings *gio.Settings, id string, type0 int32,
	keystrokes []string, name string) *GSettingsShortcut {
	gs := &GSettingsShortcut{
		BaseShortcut: BaseShortcut{
			Id:         id,
			Type:       type0,
			Keystrokes: ParseKeystrokes(keystrokes),
			Name:       name,
		},
		gsettings: gsettings,
	}

	return gs
}

func (gs *GSettingsShortcut) SaveKeystrokes() error {
	keystrokesStrv := make([]string, 0, len(gs.Keystrokes))
	for _, ks := range gs.Keystrokes {
		keystrokesStrv = append(keystrokesStrv, ks.String())
	}
	gs.gsettings.SetStrv(gs.Id, keystrokesStrv)
	logger.Debugf("GSettingsShortcut.SaveKeystrokes id: %v, keystrokes: %v", gs.Id, keystrokesStrv)
	return nil
}

func keystrokesEqual(s1 []*Keystroke, s2 []*Keystroke) bool {
	l1 := len(s1)
	l2 := len(s2)

	if l1 != l2 {
		return false
	}

	for i := 0; i < l1; i++ {
		ks1 := s1[i]
		ks2 := s2[i]

		if ks1.String() != ks2.String() {
			return false
		}
	}
	return true
}

func (gs *GSettingsShortcut) ReloadKeystrokes() bool {
	oldVal := gs.GetKeystrokes()
	id := gs.GetId()
	newVal := ParseKeystrokes(gs.gsettings.GetStrv(id))
	gs.setKeystrokes(newVal)
	return !keystrokesEqual(oldVal, newVal)
}
