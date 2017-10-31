/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"gir/gio-2.0"
	"pkg.deepin.io/lib/strv"
	"pkg.deepin.io/lib/xdg/basedir"
)

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

type Family struct {
	Id   string
	Name string

	Styles []string
	//Files  []string
}
type Families []*Family

func ListStandardFamily() Families {
	return ListFont().ListStandard().convertToFamilies()
}

func ListMonospaceFamily() Families {
	return ListFont().ListMonospace().convertToFamilies()
}

func ListAllFamily() Families {
	return ListFont().convertToFamilies()
}

func IsFontFamily(value string) bool {
	if isVirtualFont(value) {
		return true
	}

	info := ListAllFamily().Get(value)
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

	families := ListAllFamily()
	standInfo := families.Get(standard)
	if standInfo == nil {
		return fmt.Errorf("Invalid standard id '%s'", standard)
	}
	standard += " " + standInfo.preferredStyle()
	monoInfo := families.Get(monospace)
	if monoInfo == nil {
		return fmt.Errorf("Invalid monospace id '%s'", monospace)
	}
	monospace += " " + monoInfo.preferredStyle()

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

func (infos Families) GetIds() []string {
	var ids []string
	for _, info := range infos {
		ids = append(ids, info.Id)
	}
	return ids
}

func (infos Families) Get(id string) *Family {
	if id == "" {
		return nil
	}
	if isVirtualFont(id) {
		id = fcFontMatch(id)
	}

	for _, info := range infos {
		if info.Id == id {
			return info
		}
	}
	return nil
}

func (infos Families) add(info *Family) Families {
	v := infos.Get(info.Id)
	if v == nil {
		infos = append(infos, info)
		return infos
	}

	v.Styles = compositeList(v.Styles, info.Styles)
	//v.Files = compositeList(v.Files, info.Files)
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

func compositeList(l1, l2 []string) []string {
	for _, v := range l2 {
		if isItemInList(v, l1) {
			continue
		}
		l1 = append(l1, v)
	}
	return l1
}

func isItemInList(item string, list []string) bool {
	for _, v := range list {
		if item == v {
			return true
		}
	}
	return false
}

func writeFontConfig(content, file string) error {
	err := os.MkdirAll(path.Dir(file), 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, []byte(content), 0644)
}

// If set pixelsize, wps-office-wps will not show some text.
//
//func configContent(standard, mono string, pixel float64) string {
func configContent(standard, mono string) string {
	return fmt.Sprintf(`<?xml version="1.0"?>
<!DOCTYPE fontconfig SYSTEM "fonts.dtd">
<fontconfig>
    <match target="pattern">
        <test qual="any" name="family">
            <string>serif</string>
        </test>
        <edit name="family" mode="assign" binding="strong">
            <string>%s</string>
            <string>%s</string>
        </edit>
        <edit name="style" mode="assign" binding="strong">
            <string>Regular</string>
            <string>normal</string>
            <string>Standard</string>
            <string>Normale</string>
            <string>Medium</string>
            <string>Italic</string>
            <string>Black</string>
            <string>Light</string>
            <string>Bold</string>
            <string>BoldItalic</string>
            <string>DemiLight</string>
            <string>Thin</string>
       </edit>
    </match>

    <match target="pattern">
        <test qual="any" name="family">
            <string>sans-serif</string>
        </test>
        <edit name="family" mode="assign" binding="strong">
            <string>%s</string>
            <string>%s</string>
        </edit>
        <edit name="style" mode="assign" binding="strong">
            <string>Regular</string>
            <string>normal</string>
            <string>Standard</string>
            <string>Normale</string>
            <string>Medium</string>
            <string>Italic</string>
            <string>Black</string>
            <string>Light</string>
            <string>Bold</string>
            <string>BoldItalic</string>
            <string>DemiLight</string>
            <string>Thin</string>
       </edit>
    </match>

    <match target="pattern">
        <test qual="any" name="family">
            <string>monospace</string>
        </test>
        <edit name="family" mode="assign" binding="strong">
            <string>%s</string>
            <string>%s</string>
        </edit>
        <edit name="style" mode="assign" binding="strong">
            <string>Regular</string>
            <string>normal</string>
            <string>Standard</string>
            <string>Normale</string>
            <string>Medium</string>
            <string>Italic</string>
            <string>Black</string>
            <string>Light</string>
            <string>Bold</string>
            <string>BoldItalic</string>
            <string>DemiLight</string>
            <string>Thin</string>
       </edit>
    </match>

    <match target="font">
        <edit name="hinting"><bool>true</bool></edit>
        <edit name="autohint"><bool>false</bool></edit>
        <edit name="hintstyle"><const>hintfull</const></edit>
        <edit name="rgba"><const>rgb</const></edit>
        <edit name="lcdfilter"><const>lcddefault</const></edit>
        <edit name="embeddedbitmap"><bool>false</bool></edit>
        <edit name="embolden"><bool>false</bool></edit>
    </match>

</fontconfig>`, standard, fallbackStandard,
		standard, fallbackStandard,
		mono, fallbackMonospace)
}
