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
        "dlib/dbus"
)

type PreviewPath struct{}

const (
        PREVIEW_IFC = "com.deepin.daemon.PreviewPath"
)

func (op *PreviewPath) GtkPath(name string) string {
        obj := objManager.getThemeObject(name)
        return gtkPreviewPath(obj.GtkTheme)
}

func (op *PreviewPath) IconPath(name string) string {
        obj := objManager.getThemeObject(name)
        return iconPreviewPath(obj.IconTheme)
}

func (op *PreviewPath) CursorPath(name string) string {
        obj := objManager.getThemeObject(name)
        return cursorPreviewPath(obj.GtkCursorTheme)
}

func (op *PreviewPath) GetDBusInfo() dbus.DBusInfo {
        return dbus.DBusInfo{
                MANAGER_DEST,
                MANAGER_PATH,
                PREVIEW_IFC,
        }
}

func gtkPreviewPath(name string) string {
        path := ""

        list := getGtkThemeList()
        t := getThemeType(name, list)
        if len(t) <= 0 {
                return ""
        }
        if t == PATH_TYPE_SYSTEM {
                path = THUMB_GTK_PATH
        } else if t == PATH_TYPE_LOCAL {
                homeDir := getHomeDir()
                path = homeDir + THUMB_LOCAL_GTK_PATH
        }

        path += "/" + name + "/preview.png"

        return path
}

func iconPreviewPath(name string) string {
        path := ""

        list := getIconThemeList()
        t := getThemeType(name, list)
        if len(t) <= 0 {
                return ""
        }
        if t == PATH_TYPE_SYSTEM {
                path = THUMB_ICON_PATH
        } else if t == PATH_TYPE_LOCAL {
                homeDir := getHomeDir()
                path = homeDir + THUMB_LOCAL_ICON_PATH
        }

        path += "/" + name + "/preview.png"

        return path
}

func cursorPreviewPath(name string) string {
        path := ""

        list := getCursorThemeList()
        t := getThemeType(name, list)
        if len(t) <= 0 {
                return ""
        }
        if t == PATH_TYPE_SYSTEM {
                path = THUMB_CURSOR_PATH
        } else if t == PATH_TYPE_LOCAL {
                homeDir := getHomeDir()
                path = homeDir + THUMB_LOCAL_CURSOR_PATH
        }

        path += "/" + name + "/preview.png"

        return path
}
