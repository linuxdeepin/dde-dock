package main

import (
	"bytes"
	"dbus/com/deepin/api/image"
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
	_THEME_TEMPLATOR                   = template.New("theme-templator")
	_UPDATE_THEME_BACKGROUND_ID uint32 = 0
	dimg                        *image.Image
)

func init() {
	var err error
	dimg, err = image.NewImage("/com/deepin/api/Image")
	if err != nil {
		panic(err)
	}
}

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

	var err error
	theme.tplJsonData, err = theme.getThemeTplJsonData()
	if err != nil {
		panic(err) // TODO
	}

	// init properties
	theme.Background = theme.bgFile
	theme.ItemColor = theme.tplJsonData.CurrentScheme.ItemColor
	theme.SelectedItemColor = theme.tplJsonData.CurrentScheme.SelectedItemColor

	return theme
}

func (theme *Theme) setItemColor(itemColor string) {
	theme.tplJsonData.CurrentScheme.ItemColor = itemColor
	theme.customTheme()
}

func (theme *Theme) setSelectedItemColor(selectedItemColor string) {
	theme.tplJsonData.CurrentScheme.SelectedItemColor = selectedItemColor
	theme.customTheme()
}

func (theme *Theme) getThemeTplJsonData() (*TplJsonData, error) {
	fileContent, err := ioutil.ReadFile(theme.jsonFile)
	if err != nil {
		logError(err.Error())
		return nil, err
	}

	tplJsonData, err := theme.getTplJsonData(fileContent)
	if err != nil {
		return nil, err
	}
	logInfo("theme template json data: %v", tplJsonData)
	return tplJsonData, nil
}

func (theme *Theme) getTplJsonData(fileContent []byte) (*TplJsonData, error) {
	tplJsonData := &TplJsonData{}
	err := json.Unmarshal(fileContent, tplJsonData)
	if err != nil {
		logError(err.Error())
		return nil, err
	}
	return tplJsonData, nil
}

func (theme *Theme) customTheme() {
	logInfo("custom theme: %v", theme.tplJsonData.CurrentScheme)

	// generate a new theme.txt from template
	tplFileContent, err := ioutil.ReadFile(theme.tplFile)
	if err != nil {
		logError(err.Error())
		panic(err) // TODO
	}
	themeFileContent, err := theme.getCustomizedThemeContent(tplFileContent, theme.tplJsonData.CurrentScheme)
	err = ioutil.WriteFile(theme.mainFile, themeFileContent, 0664)
	if err != nil {
		logError(err.Error())
		panic(err)
	}

	// store the customized key-values to json file
	jsonContent, err := json.Marshal(theme.tplJsonData)
	if err != nil {
		logError(err.Error())
		panic(err)
	}
	err = ioutil.WriteFile(theme.jsonFile, jsonContent, 0664)
	if err != nil {
		logError(err.Error())
		panic(err)
	}
}

func (theme *Theme) getCustomizedThemeContent(fileContent []byte, tplData interface{}) ([]byte, error) {
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

// Generate background to fit the screen resolution.
func (theme *Theme) generateBackground() {
	screenWidth, screenHeight := getPrimaryScreenBestResolution()
	imgWidth, imgHeight, err := dimg.GetImageSize(theme.bgSrcFile)
	if err != nil {
		panic(err)
	}
	logInfo("source background size %dx%d", imgWidth, imgHeight)

	w, h := getImgClipSizeByResolution(screenWidth, screenHeight, imgWidth, imgHeight)
	logInfo("background size %dx%d", w, h)
	err = dimg.ClipPNG(theme.bgSrcFile, theme.bgFile, 0, 0, w, h)
	if err != nil {
		panic(err)
	}
}
