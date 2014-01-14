/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
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

import (
	"testing"
)

func TestTheme(t *testing.T) {
	InitVariable()
	ReadThemeDir(THEME_DIR)
	m := NewManager()
	if m == nil {
		t.Error("New Manager Failed!")
		return
	}

	gtk := m.GtkTheme.Get()
	for _, v := range m.EnableGtkTheme {
		m.GtkTheme.Set(v.Name)
	}
	m.GtkTheme.Set(gtk)

	icon := m.IconTheme.Get()
	for _, v := range m.EnableIconTheme {
		m.IconTheme.Set(v.Name)
	}
	m.IconTheme.Set(icon)

	cursor := m.CursorTheme.Get()
	for _, v := range m.EnableCursorTheme {
		m.CursorTheme.Set(v.Name)
	}
	m.CursorTheme.Set(cursor)

        font := infaceSettings.GetString(SCHEMA_KEY_FONT)
        ok := infaceSettings.SetString(SCHEMA_KEY_FONT, "DejaVu Sans Mono 11")
        if !ok {
                t.Error("Set Font for gsd interface error!")
        }
        ok = infaceSettings.SetString(SCHEMA_KEY_FONT, font)
        if !ok {
                t.Error("Set Font for gsd interface error!")
        }
}
