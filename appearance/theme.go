/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
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

package appearance

import (
	"os"
	"path"
	. "pkg.deepin.io/dde-daemon/appearance/utils"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/glib-2.0"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	systemPreviewPath = "/usr/share/personalization/preview"
	userPreviewPath   = ".local/share/personalization/preview"
)

type Theme struct {
	Name        string // theme id
	DisplayName string
	GtkTheme    string
	IconTheme   string
	SoundTheme  string
	CursorTheme string
	Background  string
	FontName    string
	FontMono    string
	FontSize    int32
	Type        int32
	Preview     []string

	filePath   string
	objectPath string

	eventHandler func(string)
	watcher      *dutils.WatchProxy
}

func NewTheme(info PathInfo, handler func(string)) *Theme {
	if handler == nil {
		logger.Error("Event Handler nil")
		return nil
	}

	t := &Theme{}

	t.eventHandler = handler
	t.Type = info.FileFlag
	t.filePath = info.FilePath

	t.setPropsFromFile()
	t.setPropPreview(t.getPreviewList())
	t.objectPath = themeDBusPath + t.Name

	if t.Type == FileFlagUserOwned {
		t.watcher = dutils.NewWatchProxy()
		if t.watcher != nil {
			t.watcher.SetFileList(t.getDirList())
			t.watcher.SetEventHandler(t.handleEvent)
			go t.watcher.StartWatch()
		}
	}

	return t
}

func (t *Theme) Destroy() {
	if t.watcher != nil {
		t.watcher.EndWatch()
	}
	dbus.UnInstallObject(t)
}

func (t *Theme) getThemePreview(themeType string) string {
	var target string
	switch themeType {
	case "gtk":
		target = path.Join("WindowThemes", t.GtkTheme)
	case "icon":
		target = path.Join("IconThemes", t.IconTheme)
	case "cursor":
		target = path.Join("CursorThemes", t.CursorTheme)
	}

	filename := path.Join(os.Getenv("HOME"), userPreviewPath,
		target+"-preview.png")
	if dutils.IsFileExist(filename) {
		return dutils.EncodeURI(filename, dutils.SCHEME_FILE)
	}

	filename = path.Join(systemPreviewPath,
		target+"-preview.png")
	if dutils.IsFileExist(filename) {
		return dutils.EncodeURI(filename, dutils.SCHEME_FILE)
	}

	return ""
}

func (t *Theme) getPreviewList() []string {
	var list []string

	prev := t.getThemePreview("gtk")
	if len(prev) != 0 {
		list = append(list, prev)
	}
	prev = t.getThemePreview("icon")
	if len(prev) != 0 {
		list = append(list, prev)
	}
	prev = t.getThemePreview("cursor")
	if len(prev) != 0 {
		list = append(list, prev)
	}

	return list
}

func (t *Theme) getDirList() []string {
	return []string{path.Join(t.filePath, "theme.ini")}
}

func isThemeInfoSame(info1, info2 *Theme) bool {
	if info1 == nil || info2 == nil {
		return false
	}

	if info1.GtkTheme != info2.GtkTheme ||
		info1.IconTheme != info2.IconTheme ||
		info1.CursorTheme != info2.CursorTheme ||
		info1.SoundTheme != info2.SoundTheme ||
		info1.FontName != info2.FontName ||
		info1.FontMono != info2.FontMono ||
		info1.FontSize != info2.FontSize ||
		info1.Background != info2.Background {
		return false
	}

	return true
}

func getThemeInfoFromFile(config string) (Theme, error) {
	var info Theme
	kFile := glib.NewKeyFile()
	defer kFile.Free()

	_, err := kFile.LoadFromFile(config,
		glib.KeyFileFlagsKeepComments|
			glib.KeyFileFlagsKeepTranslations)
	if err != nil {
		return info, err
	}

	info.Name, err = kFile.GetString(groupKeyTheme, themeKeyId)
	if err != nil {
		return info, err
	}
	info.DisplayName, _ = kFile.GetLocaleString(groupKeyTheme,
		themeKeyName, "\x00")
	info.GtkTheme, _ = kFile.GetString(groupKeyComponent,
		themeKeyGtk)
	info.IconTheme, _ = kFile.GetString(groupKeyComponent,
		themeKeyIcon)
	info.SoundTheme, _ = kFile.GetString(groupKeyComponent,
		themeKeySound)
	info.CursorTheme, _ = kFile.GetString(groupKeyComponent,
		themeKeyCursor)
	info.FontName, _ = kFile.GetString(groupKeyComponent,
		themeKeyFontName)
	info.FontMono, _ = kFile.GetString(groupKeyComponent,
		themeKeyFontMono)
	info.Background, _ = kFile.GetString(groupKeyComponent,
		themeKeyBackground)
	interval, _ := kFile.GetInteger(groupKeyComponent,
		themeKeyFontSize)
	info.FontSize = int32(interval)

	return info, nil
}
