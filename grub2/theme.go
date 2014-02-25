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
	"dlib/graph"
	"encoding/json"
	"io/ioutil"
	"text/template"
)

const (
	_THEME_PATH        = "/boot/grub/themes/deepin"
	_THEME_MAIN_FILE   = _THEME_PATH + "/theme.txt"
	_THEME_TPL_FILE    = _THEME_PATH + "/theme.tpl"
	_THEME_JSON_FILE   = _THEME_PATH + "/theme_tpl.json" // json stores the key-values for template file
	_THEME_BG_SRC_FILE = _THEME_PATH + "/background_source"
	_THEME_BG_FILE     = _THEME_PATH + "/background.png"
)

var (
	_UPDATE_THEME_BACKGROUND_ID uint32 = 0
)

type ThemeScheme struct {
	ItemColor, SelectedItemColor, TerminalBox, MenuPixmapStyle, ScrollbarThumb string
}

type TplJsonData struct {
	BrightScheme, DarkScheme, CurrentScheme ThemeScheme
}

type Theme struct {
	themePath   string
	mainFile    string
	tplFile     string
	jsonFile    string
	bgSrcFile   string
	bgFile      string
	tplJsonData *TplJsonData

	Background        string `access:"read"` // absolute background file path
	ItemColor         string `access:"readwrite"`
	SelectedItemColor string `access:"readwrite"`

	BackgroundUpdated func(uint32, bool)
}

func NewTheme() *Theme {
	theme := &Theme{}
	theme.themePath = _THEME_PATH
	theme.mainFile = _THEME_MAIN_FILE
	theme.tplFile = _THEME_TPL_FILE
	theme.jsonFile = _THEME_JSON_FILE
	theme.bgSrcFile = _THEME_BG_SRC_FILE
	theme.bgFile = _THEME_BG_FILE

	return theme
}

func (theme *Theme) load() {
	var err error
	theme.tplJsonData, err = theme.getThemeTplJsonData()
	if err != nil {
		panic(err)
	}

	// init properties
	theme.Background = theme.bgFile
	theme.ItemColor = theme.tplJsonData.CurrentScheme.ItemColor
	theme.SelectedItemColor = theme.tplJsonData.CurrentScheme.SelectedItemColor
	dbus.NotifyChange(theme, "Background")
	dbus.NotifyChange(theme, "ItemColor")
	dbus.NotifyChange(theme, "SelectedItemColor")
}

func (theme *Theme) setItemColor(itemColor string) {
	if len(itemColor) == 0 {
		// set a default value to avoid empty string
		itemColor = theme.tplJsonData.DarkScheme.ItemColor
	}
	theme.tplJsonData.CurrentScheme.ItemColor = itemColor
	dbus.NotifyChange(theme, "ItemColor")
	theme.customTheme()
}

func (theme *Theme) setSelectedItemColor(selectedItemColor string) {
	if len(selectedItemColor) == 0 {
		// set a default value to avoid empty string
		selectedItemColor = theme.tplJsonData.DarkScheme.SelectedItemColor
	}
	theme.tplJsonData.CurrentScheme.SelectedItemColor = selectedItemColor
	dbus.NotifyChange(theme, "SelectedItemColor")
	theme.customTheme()
}

func (theme *Theme) getThemeTplJsonData() (*TplJsonData, error) {
	fileContent, err := ioutil.ReadFile(theme.jsonFile)
	if err != nil {
		_LOGGER.Error(err.Error())
		return nil, err
	}

	tplJsonData, err := theme.getTplJsonData(fileContent)
	if err != nil {
		return nil, err
	}
	_LOGGER.Info("theme template json data: %v", tplJsonData)
	return tplJsonData, nil
}

func (theme *Theme) getTplJsonData(fileContent []byte) (*TplJsonData, error) {
	tplJsonData := &TplJsonData{}
	err := json.Unmarshal(fileContent, tplJsonData)
	if err != nil {
		_LOGGER.Error(err.Error())
		return nil, err
	}
	return tplJsonData, nil
}

// TODO split
func (theme *Theme) customTheme() {
	_LOGGER.Info("custom theme: %v", theme.tplJsonData.CurrentScheme)

	// generate a new theme.txt from template
	tplFileContent, err := ioutil.ReadFile(theme.tplFile)
	if err != nil {
		_LOGGER.Error(err.Error())
		return
	}
	themeFileContent, err := theme.getCustomizedThemeContent(tplFileContent, theme.tplJsonData.CurrentScheme)
	if err != nil {
		_LOGGER.Error(err.Error())
		return
	}
	if len(themeFileContent) == 0 {
		_LOGGER.Error("theme content is empty")
	}
	err = ioutil.WriteFile(theme.mainFile, themeFileContent, 0664)
	if err != nil {
		_LOGGER.Error(err.Error())
		return
	}

	// store the customized key-values to json file
	jsonContent, err := json.Marshal(theme.tplJsonData)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(theme.jsonFile, jsonContent, 0664)
	if err != nil {
		_LOGGER.Error(err.Error())
		return
	}
}

func (theme *Theme) getCustomizedThemeContent(fileContent []byte, tplData interface{}) ([]byte, error) {
	_THEME_TEMPLATOR := template.New("theme-templator")
	tpl, err := _THEME_TEMPLATOR.Parse(string(fileContent))
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

// TODO split
// Generate background to fit the screen resolution.
func (theme *Theme) generateBackground() {
	screenWidth, screenHeight := getPrimaryScreenBestResolution()
	imgWidth, imgHeight, err := graph.GetImageSize(theme.bgSrcFile)
	if err != nil {
		_LOGGER.Error(err.Error())
		return
	}
	_LOGGER.Info("source background size %dx%d", imgWidth, imgHeight)

	w, h := getImgClipSizeByResolution(screenWidth, screenHeight, imgWidth, imgHeight)
	_LOGGER.Info("background size %dx%d", w, h)
	err = graph.ClipPNG(theme.bgSrcFile, theme.bgFile, 0, 0, w, h)
	if err != nil {
		_LOGGER.Error(err.Error())
		return
	}
	dbus.NotifyChange(theme, "Background")
}
