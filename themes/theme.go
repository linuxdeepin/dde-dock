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
        "strings"
)

type Theme struct {
        Name           string
        Type           string  //system or local theme
        GtkTheme       string
        IconTheme      string
        GtkCursorTheme string
        GtkFontName    string
        BackgroundFile string
        PreviewPath    string
        ThumbnailPath  string
        basePath       string
        path           string
}

func newTheme(path string, info PathInfo) *Theme {
        m := &Theme{}

        m.path = path
        m.Name = info.path
        m.Type = strings.ToLower(info.t)

        if m.Type == "system" {
                m.basePath = THUMB_THEME_PATH + "/" + m.Name
        } else if m.Type == "local" {
                homeDir := getHomeDir()
                m.basePath = homeDir + THUMB_LOCAL_THEME_PATH + "/" + m.Name
        }
        m.PreviewPath = m.basePath + "/preview.png"
        m.ThumbnailPath = m.basePath + "/thumbnail.png"

        m.updateThemeInfo()

        return m
}
