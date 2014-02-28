/**
 * Copyright (c) 2013 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
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
	"bytes"
	"dlib/dbus"
	"encoding/json"
	"io/ioutil"
	"text/template"
)

const (
	themePath      = "/boot/grub/themes/deepin"
	themeMainFile  = themePath + "/theme.txt"
	themeTplFile   = themePath + "/theme.tpl"
	themeJSONFile  = themePath + "/theme_tpl.json"
	themeBgSrcFile = themePath + "/background_source"
	themeBgFile    = themePath + "/background.png"
)

var updateThemeBackgroundID uint32 // mark the asynchronous operation's ID when setup background.

// ThemeScheme stores scheme data which be used when customing deepin grub2 theme.
type ThemeScheme struct {
	ItemColor, SelectedItemColor, TerminalBox, MenuPixmapStyle, ScrollbarThumb string
}

// TplJSONData read JSON data from
// "/boot/grub/themes/deepin/theme_tpl.json" which stores the
// key-values for template file.
type TplJSONData struct {
	BrightScheme, DarkScheme, CurrentScheme ThemeScheme
}

// Theme is a dbus object which provide properties and methods to
// setup deepin grub2 theme.
type Theme struct {
	themePath   string
	mainFile    string
	tplFile     string
	jsonFile    string
	bgSrcFile   string
	bgFile      string
	tplJSONData *TplJSONData

	Background        string `access:"read"` // absolute background file path
	ItemColor         string `access:"readwrite"`
	SelectedItemColor string `access:"readwrite"`

	BackgroundUpdated func(uint32, bool)
}

// NewTheme create Theme object.
func NewTheme() *Theme {
	theme := &Theme{}
	theme.themePath = themePath
	theme.mainFile = themeMainFile
	theme.tplFile = themeTplFile
	theme.jsonFile = themeJSONFile
	theme.bgSrcFile = themeBgSrcFile
	theme.bgFile = themeBgFile

	return theme
}

func (theme *Theme) load() {
	var err error
	theme.tplJSONData, err = theme.getThemeTplJSON()
	if err != nil {
		panic(err)
	}

	// init properties
	theme.Background = theme.bgFile
	theme.ItemColor = theme.tplJSONData.CurrentScheme.ItemColor
	theme.SelectedItemColor = theme.tplJSONData.CurrentScheme.SelectedItemColor
	dbus.NotifyChange(theme, "Background")
	dbus.NotifyChange(theme, "ItemColor")
	dbus.NotifyChange(theme, "SelectedItemColor")
}

func (theme *Theme) setItemColor(itemColor string) {
	if len(itemColor) == 0 {
		// set a default value to avoid empty string
		itemColor = theme.tplJSONData.DarkScheme.ItemColor
	}
	theme.tplJSONData.CurrentScheme.ItemColor = itemColor
	dbus.NotifyChange(theme, "ItemColor")
	theme.customTheme()
}

func (theme *Theme) setSelectedItemColor(selectedItemColor string) {
	if len(selectedItemColor) == 0 {
		// set a default value to avoid empty string
		selectedItemColor = theme.tplJSONData.DarkScheme.SelectedItemColor
	}
	theme.tplJSONData.CurrentScheme.SelectedItemColor = selectedItemColor
	dbus.NotifyChange(theme, "SelectedItemColor")
	theme.customTheme()
}

func (theme *Theme) getThemeTplJSON() (*TplJSONData, error) {
	fileContent, err := ioutil.ReadFile(theme.jsonFile)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	tplJSONData, err := theme.getTplJSONData(fileContent)
	if err != nil {
		return nil, err
	}
	logger.Debug("theme template json data: %v", tplJSONData)
	return tplJSONData, nil
}

func (theme *Theme) getTplJSONData(fileContent []byte) (*TplJSONData, error) {
	tplJSONData := &TplJSONData{}
	err := json.Unmarshal(fileContent, tplJSONData)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return tplJSONData, nil
}

func (theme *Theme) customTheme() {
	logger.Debug("custom theme: %v", theme.tplJSONData.CurrentScheme)

	// generate a new theme.txt from template
	tplFileContent, err := ioutil.ReadFile(theme.tplFile)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	themeFileContent, err := theme.getCustomizedThemeContent(tplFileContent, theme.tplJSONData.CurrentScheme)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if len(themeFileContent) == 0 {
		logger.Error("theme content is empty")
	}

	grub2ext.DoCustomTheme(string(themeFileContent))

	// store the customized key-values to json file
	jsonContent, err := json.Marshal(theme.tplJSONData)
	if err != nil {
		return
	}

	grub2ext.DoWriteThemeJson(string(jsonContent))
}

func (theme *Theme) getCustomizedThemeContent(fileContent []byte, tplData interface{}) ([]byte, error) {
	templator := template.New("theme-templator")
	tpl, err := templator.Parse(string(fileContent))
	if err != nil {
		return []byte(""), err
	}

	buf := bytes.NewBufferString("")
	err = tpl.Execute(buf, tplData)
	if err != nil {
		return []byte(""), err
	}
	return buf.Bytes(), nil
}
