/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package grub2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	graphic "pkg.deepin.io/lib/gdkpixbuf"
	"pkg.deepin.io/lib/utils"
	"sync"
	"text/template"
)

const DefaultThemeDir = "/boot/grub/themes/deepin"

var (
	themeDir           = DefaultThemeDir
	themeMainFile      = themeDir + "/theme.txt"
	themeTplFile       = themeDir + "/theme.tpl"
	themeJSONFile      = themeDir + "/theme_tpl.json"
	themeBgOrigSrcFile = themeDir + "/background_origin_source"
	themeBgSrcFile     = themeDir + "/background_source"
	themeBgFile        = themeDir + "/background.png"
	themeBgThumbFile   = themeDir + "/background_thumb.png"
)

func SetDefaultThemeDir(dir string) {
	themeDir = dir
	themeMainFile = themeDir + "/theme.txt"
	themeTplFile = themeDir + "/theme.tpl"
	themeJSONFile = themeDir + "/theme_tpl.json"
	themeBgOrigSrcFile = themeDir + "/background_origin_source"
	themeBgSrcFile = themeDir + "/background_source"
	themeBgFile = themeDir + "/background.png"
	themeBgThumbFile = themeDir + "/background_thumb.png"
}

// TplJSONData read JSON data from
// "/boot/grub/themes/deepin/theme_tpl.json" which stores the
// key-values for template file.
type TplJSONData struct {
	BrightScheme, DarkScheme, CurrentScheme ThemeScheme
}

// ThemeScheme stores scheme data which be used when customing deepin grub2 theme.
type ThemeScheme struct {
	ItemColor, SelectedItemColor, TerminalBox, MenuPixmapStyle, ScrollbarThumb string
}

// Theme is a dbus object which provide properties and methods to
// setup deepin grub2 theme.
type Theme struct {
	themeDir    string
	mainFile    string
	tplFile     string
	jsonFile    string
	bgSrcFile   string
	bgFile      string
	bgThumbFile string
	tplJSONData *TplJSONData

	Updating   bool
	updateLock sync.Mutex

	Background        string // background thumbnail, always equal with bgThumbFile, used by front-end
	ItemColor         string `access:"readwrite"`
	SelectedItemColor string `access:"readwrite"`
}

// NewTheme create Theme object.
func NewTheme() *Theme {
	theme := &Theme{}
	theme.themeDir = themeDir
	theme.mainFile = themeMainFile
	theme.tplFile = themeTplFile
	theme.jsonFile = themeJSONFile
	theme.bgSrcFile = themeBgSrcFile
	theme.bgFile = themeBgFile
	theme.bgThumbFile = themeBgThumbFile

	return theme
}

func newTplJSONData() (d *TplJSONData) {
	d = &TplJSONData{}
	return
}

func (theme *Theme) initTheme() {
	var err error
	theme.tplJSONData, err = theme.getThemeTplJSON()
	if err != nil {
		logger.Error(err)
		theme.tplJSONData = newTplJSONData()
	}

	// init properties
	theme.setPropBackground(theme.bgThumbFile)
	theme.setPropItemColor(theme.tplJSONData.CurrentScheme.ItemColor)
	theme.setPropSelectedItemColor(theme.tplJSONData.CurrentScheme.SelectedItemColor)
}

// reset to default configuration
func (theme *Theme) reset() {
	theme.tplJSONData.CurrentScheme = theme.tplJSONData.DarkScheme
	theme.setPropItemColor(theme.tplJSONData.CurrentScheme.ItemColor)
	theme.setPropSelectedItemColor(theme.tplJSONData.CurrentScheme.SelectedItemColor)
	theme.customTheme()

	// reset theme background
	go func() {
		theme.updateLock.Lock()
		defer theme.updateLock.Unlock()
		theme.setPropUpdating(true)
		grub2extDoResetThemeBackground()
		screenWidth, screenHeight := parseCurrentGfxmode()
		grub2extDoGenerateThemeBackground(screenWidth, screenHeight)
		theme.setPropBackground(theme.bgThumbFile)
		theme.setPropUpdating(false)
	}()
}

// Fix issue that the theme background will keep default size as
// 1024x768 if updating grub-themes-deepin package lonely
func (theme *Theme) regenerateBackgroundIfNeed() {
	logger.Debug("check if need regenerate theme background")
	wantWidth, wantHeight := parseCurrentGfxmode()
	bgw, bgh, _ := graphic.GetImageSize(theme.bgFile)
	srcbgw, srcbgh, _ := graphic.GetImageSize(theme.bgSrcFile)
	needGenerate := false
	logger.Debugf("expected size: %dx%d, source background: %dx%d, background: %dx%d",
		wantWidth, wantHeight, srcbgw, srcbgh, bgw, bgh)
	if srcbgw >= int(wantWidth) && srcbgh >= int(wantHeight) {
		// if source background is bigger than expected size, the size
		// of background should equal with it
		if delta(float64(bgw), float64(wantWidth)) > 5 ||
			delta(float64(bgh), float64(wantHeight)) > 5 {
			needGenerate = true
		}
	} else {
		// if source background is smaller than expected size, the
		// scale of backgound should equle with it
		scalebg := float64(bgw) / float64(bgh)
		scaleScreen := float64(wantWidth) / float64(wantHeight)
		if delta(scalebg, scaleScreen) > 0.1 {
			needGenerate = true
		}
	}

	// regenerate theme background if thumbnail not exists
	if !utils.IsFileExist(theme.bgThumbFile) {
		needGenerate = true
	}

	if needGenerate {
		grub2extDoGenerateThemeBackground(wantWidth, wantHeight)
		theme.setPropBackground(theme.bgThumbFile)
		logger.Info("update theme background success")
	}
}

func (theme *Theme) getThemeTplJSON() (*TplJSONData, error) {
	fileContent, err := ioutil.ReadFile(theme.jsonFile)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	tplJSONData, err := theme.getTplJSONData(fileContent)
	if err != nil {
		return nil, err
	}
	logger.Debugf("theme template json data: %v", tplJSONData)
	return tplJSONData, nil
}

func (theme *Theme) getTplJSONData(fileContent []byte) (*TplJSONData, error) {
	tplJSONData := &TplJSONData{}
	err := json.Unmarshal(fileContent, tplJSONData)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return tplJSONData, nil
}

func (theme *Theme) customTheme() {
	logger.Debugf("custom theme: %v", theme.tplJSONData.CurrentScheme)

	// generate a new theme.txt from template
	tplFileContent, err := ioutil.ReadFile(theme.tplFile)
	if err != nil {
		logger.Error(err)
		return
	}
	themeFileContent, err := theme.getCustomizedThemeContent(tplFileContent, theme.tplJSONData.CurrentScheme)
	if err != nil {
		logger.Error(err)
		return
	}
	if len(themeFileContent) == 0 {
		logger.Error("theme content is empty")
	}

	grub2extDoWriteThemeMainFile(string(themeFileContent))

	// store the customized key-values to json file
	jsonContent, err := json.Marshal(theme.tplJSONData)
	if err != nil {
		return
	}

	grub2extDoWriteThemeTplFile(string(jsonContent))
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
