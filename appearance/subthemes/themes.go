/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package subthemes

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/linuxdeepin/go-dbus-factory/com.deepin.api.cursorhelper"
	"pkg.deepin.io/dde/api/theme_thumb"
	"pkg.deepin.io/dde/api/themes"
	"pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/dbus1"
)

const (
	appearanceSchema  = "com.deepin.dde.appearance"
	gsKeyExcludedIcon = "excluded-icon-themes"
)

type Theme struct {
	Id   string
	Path string

	Deletable bool
}
type Themes []*Theme

var (
	cacheGtkThemes    Themes
	cacheIconThemes   Themes
	cacheCursorThemes Themes

	home = os.Getenv("HOME")
)

func RefreshGtkThemes() {
	cacheGtkThemes = getThemes(themes.ListGtkTheme())
}

func RefreshIconThemes() {
	infos := getThemes(themes.ListIconTheme())
	s := gio.NewSettings(appearanceSchema)
	defer s.Unref()
	blacklist := s.GetStrv(gsKeyExcludedIcon)

	var ret Themes
	for _, info := range infos {
		if isItemInList(info.Id, blacklist) {
			continue
		}
		ret = append(ret, info)
	}
	cacheIconThemes = ret
}

func RefreshCursorThemes() {
	cacheCursorThemes = getThemes(themes.ListCursorTheme())
}

func ListGtkTheme() Themes {
	if len(cacheGtkThemes) == 0 {
		RefreshGtkThemes()
	}
	return cacheGtkThemes
}

func ListIconTheme() Themes {
	if len(cacheIconThemes) == 0 {
		RefreshIconThemes()
	}
	return cacheIconThemes
}

func ListCursorTheme() Themes {
	if len(cacheCursorThemes) == 0 {
		RefreshCursorThemes()
	}
	return cacheCursorThemes
}

func IsGtkTheme(id string) bool {
	return themes.IsThemeInList(id, themes.ListGtkTheme())
}

func IsIconTheme(id string) bool {
	return themes.IsThemeInList(id, themes.ListIconTheme())
}

func IsCursorTheme(id string) bool {
	return themes.IsThemeInList(id, themes.ListCursorTheme())
}

func SetGtkTheme(id string) error {
	return themes.SetGtkTheme(id)
}

func SetIconTheme(id string) error {
	return themes.SetIconTheme(id)
}

func SetCursorTheme(id string) error {
	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	helper := cursorhelper.NewCursorHelper(sessionBus)
	return helper.Set(0, id)
}

func GetGtkThumbnail(id string) (string, error) {
	info := ListGtkTheme().Get(id)
	if info == nil {
		return "", fmt.Errorf("not found %q", id)
	}

	descFile := path.Join(info.Path, "index.theme")
	return theme_thumb.GetGtk(id, descFile)
}

func GetIconThumbnail(id string) (string, error) {
	info := ListIconTheme().Get(id)
	if info == nil {
		return "", fmt.Errorf("not found %q", id)
	}

	descFile := path.Join(info.Path, "index.theme")
	return theme_thumb.GetIcon(id, descFile)
}

func GetCursorThumbnail(id string) (string, error) {
	info := ListCursorTheme().Get(id)
	if info == nil {
		return "", fmt.Errorf("not found %q", id)
	}
	descFile := path.Join(info.Path, "cursor.theme")
	return theme_thumb.GetCursor(id, descFile)
}

func (infos Themes) GetIds() []string {
	var ids []string
	for _, info := range infos {
		ids = append(ids, info.Id)
	}
	return ids
}

func (infos Themes) Get(id string) *Theme {
	for _, info := range infos {
		if id == info.Id {
			return info
		}
	}
	return nil
}

func (infos Themes) ListGet(ids []string) Themes {
	var ret Themes
	for _, id := range ids {
		info := infos.Get(id)
		if info == nil {
			continue
		}
		ret = append(ret, info)
	}
	return ret
}

func (infos Themes) Delete(id string) error {
	info := infos.Get(id)
	if info == nil {
		return fmt.Errorf("not found %q", id)
	}
	return info.Delete()
}

func (info *Theme) Delete() error {
	if !info.Deletable {
		return fmt.Errorf("permission denied")
	}
	return os.RemoveAll(info.Path)
}

func getThemes(files []string) Themes {
	var infos Themes
	for _, v := range files {
		infos = append(infos, &Theme{
			Id:        path.Base(v),
			Path:      v,
			Deletable: isDeletable(v),
		})
	}
	return infos
}

func isDeletable(file string) bool {
	return strings.Contains(file, home)
}

func isItemInList(item string, list []string) bool {
	for _, v := range list {
		if item == v {
			return true
		}
	}
	return false
}
