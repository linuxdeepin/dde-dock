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

type Manager struct {
        ThemeList    []string `access:"readwrite"`
        CurrentTheme string   `access:"readwrite"`
}

/*
   Return only contains thumbnail theme
   First, get all local themes
   Then, determine whether contains thumbnail
*/
func (op *Manager) GetGtkThemeList() []string {
        return []string{}
}

/*
   Return only contains thumbnail theme
   First, get all local themes
   Then, determine whether contains thumbnail
*/
func (op *Manager) GetIconNameList() []string {
        return []string{}
}

/*
   Return only contains thumbnail theme
   First, get all local themes
   Then, determine whether contains thumbnail
*/
func (op *Manager) GetCursorNameList() []string {
        return []string{}
}

// Has not yet been determined
func (op *Manager) GetFontNameList() []string {
        return []string{}
}

/*
   Unlimited
   Return all bg.
*/
func (op *Manager) GetBackgroundList() []string {
        return []string{}
}

func newManager() *Manager {
        m := &Manager{}

        m.setPropName("ThemeList")
        m.setPropName("CurrentTheme")

        return m
}
