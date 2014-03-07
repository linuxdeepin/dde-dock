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
        ThemeList    []string
        CurrentTheme string  `access:"readwrite"`
        GtkThemeList []string
        //GtkBasePath     string
        IconThemeList []string
        //IconBasePath    string
        CursorThemeList []string
        //CursorBasePath  string
        FontThemeList []string
        pathNameMap   map[string]PathInfo
}

func (op *Manager) GetGtkBasePath(name string) string {
        list := getGtkThemeList()

        t := getThemeType(name, list)
        if t == PATH_TYPE_SYSTEM {
                return THUMB_GTK_PATH + "/" + name
        } else if t == PATH_TYPE_LOCAL {
                return THUMB_LOCAL_GTK_PATH + "/" + name
        }

        return ""
}

func (op *Manager) GetIconBasePath(name string) string {
        list := getIconThemeList()

        t := getThemeType(name, list)
        if t == PATH_TYPE_SYSTEM {
                return THUMB_ICON_PATH + "/" + name
        } else if t == PATH_TYPE_LOCAL {
                return THUMB_LOCAL_ICON_PATH + "/" + name
        }

        return ""
}

func (op *Manager) GetCursorBasePath(name string) string {
        list := getCursorThemeList()

        t := getThemeType(name, list)
        if t == PATH_TYPE_SYSTEM {
                return THUMB_CURSOR_PATH + "/" + name
        } else if t == PATH_TYPE_LOCAL {
                return THUMB_LOCAL_CURSOR_PATH + "/" + name
        }

        return ""
}

func getThemeType(name string, list []PathInfo) string {
        for _, l := range list {
                if name == l.path {
                        return l.t
                }
        }

        return ""
}

func (op *Manager) SetTheme(gtk, icon, cursor, font string) string {
        return ""
}

func getThemeList() []PathInfo {
        return getThemeThumbList()
}

/*
   Return only contains thumbnail theme
   First, get all local themes
   Then, determine whether contains thumbnail
*/
func getGtkThemeList() []PathInfo {
        valid := getValidGtkThemes()
        thumb := getGtkThumbList()

        list := []PathInfo{}
        for _, v := range valid {
                if isElementExist(v, thumb) {
                        list = append(list, v)
                }
        }

        return list
}

/*
   Return only contains thumbnail theme
   First, get all local themes
   Then, determine whether contains thumbnail
*/
func getIconThemeList() []PathInfo {
        valid := getValidIconThemes()
        thumb := getIconThumbList()

        list := []PathInfo{}
        for _, v := range valid {
                if isElementExist(v, thumb) {
                        list = append(list, v)
                }
        }

        return list
}

/*
   Return only contains thumbnail theme
   First, get all local themes
   Then, determine whether contains thumbnail
*/
func getCursorThemeList() []PathInfo {
        valid := getValidCursorThemes()
        thumb := getCursorThumbList()

        list := []PathInfo{}
        for _, v := range valid {
                if isElementExist(v, thumb) {
                        list = append(list, v)
                }
        }

        return list
}

// Has not yet been determined
func getFontNameList() []string {
        return []string{}
}

/*
   Unlimited
   Return all bg.
*/
func getBackgroundList() []string {
        return []string{}
}

func newManager() *Manager {
        m := &Manager{}

        m.pathNameMap = make(map[string]PathInfo)
        m.setPropName("ThemeList")
        m.setPropName("CurrentTheme")
        m.setPropName("GtkThemeList")
        m.setPropName("IconThemeList")
        m.setPropName("CursorThemeList")

        return m
}
