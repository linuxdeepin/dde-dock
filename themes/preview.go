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
        return gtkPreviewPath(name)
}

func (op *PreviewPath) IconPath(name string) string {
        return iconPreviewPath(name)
}

func (op *PreviewPath) CursorPath(name string) string {
        return cursorPreviewPath(name)
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
                path = PREVIEW_GTK_PATH
        } else if t == PATH_TYPE_LOCAL {
                homeDir := getHomeDir()
                path = homeDir + PREVIEW_LOCAL_GTK_PATH
        }

        path += "/" + name + "/preview.png"
        if ok := objUtil.IsFileExist(path); !ok {
                return ""
        }

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
                path = PREVIEW_ICON_PATH
        } else if t == PATH_TYPE_LOCAL {
                homeDir := getHomeDir()
                path = homeDir + PREVIEW_LOCAL_ICON_PATH
        }

        path += "/" + name + "/preview.png"
        if ok := objUtil.IsFileExist(path); !ok {
                return ""
        }

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
                path = PREVIEW_CURSOR_PATH
        } else if t == PATH_TYPE_LOCAL {
                homeDir := getHomeDir()
                path = homeDir + PREVIEW_LOCAL_CURSOR_PATH
        }

        path += "/" + name + "/preview.png"
        if ok := objUtil.IsFileExist(path); !ok {
                return ""
        }

        return path
}
