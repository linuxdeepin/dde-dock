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

package fonts

import (
	"crypto/md5"
	"fmt"
	"gir/gio-2.0"
	"path"
	"pkg.deepin.io/lib/strv"
	"pkg.deepin.io/lib/xdg/basedir"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Family struct {
	Id   string
	Name string

	Styles []string

	Monospace bool
	Show      bool
}

type FamilyHashTable map[string]*Family

const (
	fallbackStandard  = "Noto Sans"
	fallbackMonospace = "Noto Mono"
	defaultDPI        = 96

	xsettingsSchema = "com.deepin.xsettings"
	gsKeyFontName   = "gtk-font-name"
)

var (
	locker    sync.Mutex
	xsSetting = gio.NewSettings(xsettingsSchema)

	DeepinFontConfig = path.Join(basedir.GetUserConfigDir(), "fontconfig", "conf.d", "99-deepin.conf")
)

var stylePriorityList = []string{
	"Regular",
	"normal",
	"Standard",
	"Normale",
	"Medium",
	"Italic",
	"Black",
	"Light",
	"Bold",
	"BoldItalic",
	"DemiLight",
	"Thin",
}

func IsFontFamily(value string) bool {
	if isVirtualFont(value) {
		return true
	}

	info := GetFamilyTable().GetFamily(value)
	if info != nil {
		return true
	}
	return false
}

func IsFontSizeValid(size float64) bool {
	if size >= 7.0 && size <= 22.0 {
		return true
	}
	return false
}

func SetFamily(standard, monospace string, size float64) error {
	locker.Lock()
	defer locker.Unlock()

	if isVirtualFont(standard) {
		standard = fcFontMatch(standard)
	}
	if isVirtualFont(monospace) {
		monospace = fcFontMatch(monospace)
	}

	table := GetFamilyTable()
	standInfo := table.GetFamily(standard)
	if standInfo == nil {
		return fmt.Errorf("Invalid standard id '%s'", standard)
	}
	// standard += " " + standInfo.preferredStyle()
	monoInfo := table.GetFamily(monospace)
	if monoInfo == nil {
		return fmt.Errorf("Invalid monospace id '%s'", monospace)
	}
	// monospace += " " + monoInfo.preferredStyle()

	// fc-match can not real time update
	/*
		curStand := fcFontMatch("sans-serif")
		curMono := fcFontMatch("monospace")
		if (standInfo.Id == curStand || standInfo.Name == curStand) &&
			(monoInfo.Id == curMono || monoInfo.Name == curMono) {
			return nil
		}
	*/

	err := writeFontConfig(configContent(standInfo.Id, monoInfo.Id), DeepinFontConfig)
	if err != nil {
		return err
	}
	return setFontByXSettings(standard, size)
}

func GetFontSize() float64 {
	return getFontSize(xsSetting)
}

func (table FamilyHashTable) ListMonospace() []string {
	var ids []string
	for _, info := range table {
		if !info.Monospace {
			continue
		}
		ids = append(ids, info.Id)
	}
	sort.Strings(ids)
	return ids
}

func (table FamilyHashTable) ListStandard() []string {
	var ids []string
	for _, info := range table {
		if info.Monospace || !info.Show {
			continue
		}
		ids = append(ids, info.Id)
	}
	sort.Strings(ids)
	return ids
}

func (table FamilyHashTable) Get(key string) *Family {
	info, _ := table[key]
	return info
}

func (table FamilyHashTable) GetFamily(id string) *Family {
	info, ok := table[sumStrHash(id)]
	if !ok {
		return nil
	}
	return info
}

func (table FamilyHashTable) GetFamilies(ids []string) []*Family {
	var infos []*Family
	for _, id := range ids {
		info, ok := table[sumStrHash(id)]
		if !ok {
			continue
		}
		infos = append(infos, info)
	}
	return infos
}

func (info Family) preferredStyle() string {
	styles := strv.Strv(info.Styles)
	for _, v := range stylePriorityList {
		if styles.Contains(v) {
			return v
		}
	}
	return ""
}

func setFontByXSettings(name string, size float64) error {
	if size == -1 {
		size = getFontSize(xsSetting)
	}
	v := fmt.Sprintf("%s %v", name, size)
	if v == xsSetting.GetString(gsKeyFontName) {
		return nil
	}

	xsSetting.SetString(gsKeyFontName, v)
	return nil
}

func getFontSize(setting *gio.Settings) float64 {
	value := setting.GetString(gsKeyFontName)
	if len(value) == 0 {
		return 0
	}

	array := strings.Split(value, " ")
	size, _ := strconv.ParseFloat(array[len(array)-1], 64)
	return size
}

func isVirtualFont(name string) bool {
	switch name {
	case "monospace", "mono", "sans-serif", "sans", "serif":
		return true
	}
	return false
}

func sumStrHash(v string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(v)))
}
