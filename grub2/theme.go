package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"
	"text/template"
)

const (
	_THEME_PATH        = "/boot/grub/themes/deepin" // TODO
	_THEME_MAIN_FILE   = _THEME_PATH + "/theme.txt"
	_THEME_TPL_FILE    = _THEME_PATH + "/theme.tpl"
	_THEME_JSON_FILE   = _THEME_PATH + "/theme_tpl.json" // json stores the key-values for template file
	_THEME_BG_SRC_FILE = _THEME_PATH + "/background_source.png"
	_THEME_BG_FILE     = _THEME_PATH + "/background.png"
)

var _THEME_TEMPLATOR = template.New("theme-templator")

// TODO
type TplValues struct {
	// Background, ItemColor, SelectedItemColor string
	ItemColor, SelectedItemColor string
}
type TplJsonData struct {
	// DefaultTplValue, LastTplValue TplValues
	ItemColor, SelectedItemColor string
}

type Theme struct {
	themePath string
	mainFile  string
	tplFile   string
	jsonFile  string
	bgSrcFile string // TODO
	bgFile    string // TODO
	// relBgFile string // TODO remove relative background file path

	Background string `access:"read"` // absolute background file path
	// BackgroundSource  string `access:"read"` // TODO
	ItemColor         string `access:"readwrite"`
	SelectedItemColor string `access:"readwrite"`
}

func NewTheme() *Theme {
	theme := &Theme{}
	theme.themePath = _THEME_PATH
	theme.mainFile = _THEME_MAIN_FILE
	theme.tplFile = _THEME_TPL_FILE
	theme.jsonFile = _THEME_JSON_FILE
	theme.bgSrcFile = _THEME_BG_SRC_FILE
	theme.bgFile = _THEME_BG_FILE

	tplJsonData, err := theme.getThemeTplJsonData()
	if err != nil {
		panic(err) // TODO
	}

	// theme.relBgFile = tplJsonData.LastTplValue.Background // TODO
	theme.makeBackground()
	theme.ItemColor = tplJsonData.ItemColor
	theme.SelectedItemColor = tplJsonData.SelectedItemColor

	return theme
}

// TODO Update variable 'Background'
func (theme *Theme) makeBackground() {
	theme.Background = theme.bgFile
	// theme.Background = theme.getBgFileAbsPath(theme.relBgFile)
	if !isFileExists(theme.Background) {
		logError("theme: background file is not exists, %s", theme.Background)
	}
}

// TODO remove
// func (theme *Theme) setBackground(background string) {
// 	// copy background file to theme dir if need
// 	theme.copyBgFileToThemeDir(background)
// 	theme.relBgFile = theme.getNewBgFileName(background)
// 	// theme.customTheme()			// TODO
// }

func (theme *Theme) setItemColor(itemColor string) {
	// theme.customTheme()			// TODO
}

func (theme *Theme) setSelectedItemColor(selectedItemColor string) {
	// theme.customTheme()			// TODO
}

// TODO remove
func (theme *Theme) copyBgFileToThemeDir(imageFile string) (newBgFile string, err error) {
	bgFileName := theme.getNewBgFileName(imageFile)
	newBgFile = theme.getBgFileAbsPath(bgFileName)
	_, err = copyFile(newBgFile, imageFile)
	if err != nil {
		logError(err.Error())
	}
	return
}

// TODO remove
func (theme *Theme) getBgFileAbsPath(bgFileRelPath string) string {
	bgFileAbsPath := path.Join(theme.themePath, bgFileRelPath)
	return bgFileAbsPath
}

// TODO remove
func (theme *Theme) getNewBgFileName(imageFile string) string {
	fileName := path.Base(imageFile)
	i := strings.LastIndex(fileName, ".")
	fileExt := fileName[i:]
	return "background" + fileExt
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
	tplJsonData, err := theme.getThemeTplJsonData()
	if err != nil {
		panic(err)
	}
	// tplJsonData.Background = theme.relBgFile // TODO
	tplJsonData.ItemColor = theme.ItemColor
	tplJsonData.SelectedItemColor = theme.SelectedItemColor

	// generate a new theme.txt from template
	tplFileContent, err := ioutil.ReadFile(theme.tplFile)
	if err != nil {
		logError(err.Error())
		panic(err) // TODO
	}
	themeFileContent, err := theme.getCustomizedThemeContent(tplFileContent, tplJsonData)
	err = ioutil.WriteFile(theme.mainFile, themeFileContent, 0644)
	if err != nil {
		logError(err.Error())
		panic(err)
	}

	// store the customized key-values to json file
	jsonContent, err := json.Marshal(tplFileContent)
	if err != nil {
		logError(err.Error())
		panic(err)
	}
	err = ioutil.WriteFile(theme.jsonFile, jsonContent, 0644)
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
