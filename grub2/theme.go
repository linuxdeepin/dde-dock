/*
 * Copyright (C) 2013 ~ 2017 Deepin Technology Co., Ltd.
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

package grub2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"text/template"

	graphic "pkg.deepin.io/lib/gdkpixbuf"
	"pkg.deepin.io/lib/utils"
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
)

func SetDefaultThemeDir(dir string) {
	themeDir = dir
	themeMainFile = themeDir + "/theme.txt"
	themeTplFile = themeDir + "/theme.tpl"
	themeJSONFile = themeDir + "/theme_tpl.json"
	themeBgOrigSrcFile = themeDir + "/background_origin_source"
	themeBgSrcFile = themeDir + "/background_source"
	themeBgFile = themeDir + "/background.png"
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
	grub2       *Grub2
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

	ItemColor         string
	SelectedItemColor string

	// Signal:
	BackgroundChanged func()
}

// NewTheme create Theme object.
func NewTheme(grub2 *Grub2) *Theme {
	theme := &Theme{}
	theme.grub2 = grub2
	theme.themeDir = themeDir
	theme.mainFile = themeMainFile
	theme.tplFile = themeTplFile
	theme.jsonFile = themeJSONFile
	theme.bgSrcFile = themeBgSrcFile
	theme.bgFile = themeBgFile

	return theme
}

func newTplJSONData() (d *TplJSONData) {
	d = &TplJSONData{}
	return
}

func (theme *Theme) getScreenWidthHeight() (w, h uint16) {
	var err error
	w, h, err = theme.grub2.getScreenWidthHeight()
	if err != nil {
		return 1024, 768
	}
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
	theme.setPropItemColor(theme.tplJSONData.CurrentScheme.ItemColor)
	theme.setPropSelectedItemColor(theme.tplJSONData.CurrentScheme.SelectedItemColor)
}

// reset to default configuration
func (theme *Theme) reset() {
	theme.tplJSONData.CurrentScheme = theme.tplJSONData.DarkScheme
	theme.setPropItemColor(theme.tplJSONData.CurrentScheme.ItemColor)
	theme.setPropSelectedItemColor(theme.tplJSONData.CurrentScheme.SelectedItemColor)
	theme.setCustomTheme()

	// reset theme background
	go func() {
		theme.updateLock.Lock()
		defer theme.updateLock.Unlock()
		theme.setPropUpdating(true)
		resetThemeBackground()
		screenWidth, screenHeight := theme.getScreenWidthHeight()
		generateThemeBackground(screenWidth, screenHeight)
		theme.setPropUpdating(false)
		theme.emitSignalBackgroundChanged()
	}()
}

// Fix issue that the theme background will keep default size as
// 1024x768 if updating grub-themes-deepin package lonely
func (theme *Theme) regenerateBackgroundIfNeed() {
	logger.Debug("check if need regenerate theme background")
	wantWidth, wantHeight := theme.getScreenWidthHeight()
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
		generateThemeBackground(wantWidth, wantHeight)
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

func (theme *Theme) setCustomTheme() {
	logger.Debugf("set custom theme: %v", theme.tplJSONData.CurrentScheme)

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

	err = writeThemeMainFile(themeFileContent)
	if err != nil {
		logger.Warning(err)
		return
	}

	// store the customized key-values to json file
	jsonContent, err := json.Marshal(theme.tplJSONData)
	if err != nil {
		return
	}

	writeThemeTplFile(jsonContent)
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

func delta(v1, v2 float64) float64 {
	if v1 > v2 {
		return v1 - v2
	}
	return v2 - v1
}

// write file content to "/boot/grub/themes/deepin/theme.txt".
func writeThemeMainFile(data []byte) error {
	return ioutil.WriteFile(themeMainFile, data, 0664)
}

// write file content to "/boot/grub/themes/deepin/theme_tpl.json".
func writeThemeTplFile(data []byte) error {
	return ioutil.WriteFile(themeJSONFile, data, 0664)
}

// link background_origin_source to background_source
func resetThemeBackground() error {
	os.Remove(themeBgSrcFile)
	return os.Symlink(themeBgOrigSrcFile, themeBgSrcFile)
}

// generate the background for deepin grub2
// theme depends on screen resolution.
func generateThemeBackground(screenWidth, screenHeight uint16) (err error) {
	imgWidth, imgHeight, err := graphic.GetImageSize(themeBgSrcFile)
	if err != nil {
		return err
	}
	logger.Infof("source background size %dx%d", imgWidth, imgHeight)
	logger.Infof("background size %dx%d", screenWidth, screenHeight)
	return graphic.ScaleImagePrefer(themeBgSrcFile, themeBgFile, int(screenWidth), int(screenHeight), graphic.GDK_INTERP_HYPER, graphic.FormatPng)
}

func (theme *Theme) doSetBackgroundSourceFile(imageFile string) bool {
	theme.updateLock.Lock()
	defer theme.updateLock.Unlock()
	theme.setPropUpdating(true)
	screenWidth, screenHeight := theme.getScreenWidthHeight()
	setThemeBackgroundSourceFile(imageFile, screenWidth, screenHeight)
	theme.setPropUpdating(false)
	theme.emitSignalBackgroundChanged()

	// set item color through background's dominant color
	_, _, v, _ := graphic.GetDominantColorOfImage(theme.bgSrcFile)
	if v < 0.5 {
		// background is dark
		theme.tplJSONData.CurrentScheme = theme.tplJSONData.DarkScheme
		logger.Info("background is dark, use the dark theme scheme")
	} else {
		// background is bright
		theme.tplJSONData.CurrentScheme = theme.tplJSONData.BrightScheme
		logger.Info("background is bright, so use the bright theme scheme")
	}
	theme.setPropItemColor(theme.tplJSONData.CurrentScheme.ItemColor)
	theme.setPropSelectedItemColor(theme.tplJSONData.CurrentScheme.SelectedItemColor)
	theme.setCustomTheme()

	logger.Info("update background sucess")
	return true
}

// setup a new background source file
// for deepin grub2 theme, and then generate the background depends on
// screen resolution.
func setThemeBackgroundSourceFile(imageFile string, screenWidth, screenHeight uint16) (err error) {
	// if background source file is a symlink, just delete it
	if utils.IsSymlink(themeBgSrcFile) {
		os.Remove(themeBgSrcFile)
	}

	// backup background source file
	err = utils.CopyFile(imageFile, themeBgSrcFile)
	if err != nil {
		return err
	}

	// generate a new background
	return generateThemeBackground(screenWidth, screenHeight)
}
