package main

import (
	"bytes"
	"encoding/json"
	"path"
	"strings"
	"text/template"
)

const (
	_THEME_DIR       = "/boot/grub/themes"
	_THEME_MAIN_FILE = "theme.txt"
	_THEME_TPL_FILE  = "theme.tpl"

	// json stores the key-values for template file
	_THEME_TPL_JSON_LAST    = "theme_tpl_last.json"
	_THEME_TPL_JSON_DEFAULT = "theme_tpl_default.json"

	_THEME_TPL_KEY_BACKGROUND          = "Background"
	_THEME_TPL_KEY_ITEM_COLOR          = "ItemColor"
	_THEME_TPL_KEY_SELECTED_ITEM_COLOR = "SelectedItemColor"
)

var _THEME_TEMPLATOR = template.New("theme-templator")

type ThemeManager struct {
	enabledThemeMainFile string
}

func NewThemeManager() *ThemeManager {
	tm := &ThemeManager{}
	tm.enabledThemeMainFile = ""
	return tm
}

func (tm *ThemeManager) setEnabledThemeMainFile(file string) {
	// if the theme.txt file is not under theme dir(/boot/grub/themes), ignore it
	if strings.HasPrefix(file, _THEME_DIR) {
		tm.enabledThemeMainFile = file
	} else {
		tm.enabledThemeMainFile = ""
	}
}

func (tm *ThemeManager) getEnabledThemeMainFile() string {
	return tm.enabledThemeMainFile
}

func (tm *ThemeManager) isThemeValid(themeName string) bool {
	_, okPath := tm.getThemePath(themeName)
	_, okMainFile := tm.getThemeMainFile(themeName)
	if okPath && okMainFile {
		return true
	} else {
		return false
	}
}

func (tm *ThemeManager) isThemeArchiveValid(archive string) bool {
	_, err := findFileInTarGz(archive, _THEME_MAIN_FILE)
	if err != nil {
		return false
	}
	return true
}

func (tm *ThemeManager) getThemeName(themeMainFile string) string {
	if len(themeMainFile) == 0 {
		return ""
	}
	return path.Base(path.Dir(themeMainFile))
}

func (tm *ThemeManager) getThemePath(themeName string) (themePath string, existed bool) {
	themePath = path.Join(_THEME_DIR, themeName)
	existed = isFileExists(themePath)
	return
}

func (tm *ThemeManager) getThemeMainFile(themeName string) (file string, existed bool) {
	file = path.Join(_THEME_DIR, themeName, _THEME_MAIN_FILE)
	existed = isFileExists(file)
	return
}

func (tm *ThemeManager) getThemeTplFile(themeName string) (file string, existed bool) {
	file = path.Join(_THEME_DIR, themeName, _THEME_TPL_FILE)
	existed = isFileExists(file)
	return
}

func (tm *ThemeManager) getThemeTplDefaultJsonFile(themeName string) (file string, existed bool) {
	file = path.Join(_THEME_DIR, themeName, _THEME_TPL_JSON_DEFAULT)
	existed = isFileExists(file)
	return
}

func (tm *ThemeManager) getThemeTplLastJsonFile(themeName string) (file string, existed bool) {
	file = path.Join(_THEME_DIR, themeName, _THEME_TPL_JSON_LAST)
	existed = isFileExists(file)
	return
}

func (tm *ThemeManager) getCustomizedThemeContent(fileContent []byte, tplData interface{}) ([]byte, error) {
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

func (tm *ThemeManager) getValuesInJson(fileContent []byte) (background, itemColor, selectedItemColor string, ok bool) {
	tplData := make(map[string]string)
	err := json.Unmarshal(fileContent, &tplData)
	if err != nil {
		logError(err.Error()) // TODO
		return "", "", "", false
	}
	background = tplData[_THEME_TPL_KEY_BACKGROUND]
	itemColor = tplData[_THEME_TPL_KEY_ITEM_COLOR]
	selectedItemColor = tplData[_THEME_TPL_KEY_SELECTED_ITEM_COLOR]
	return background, itemColor, selectedItemColor, true
}

func (tm *ThemeManager) getBgFileAbsPath(themeName, bgFileRelPath string) string {
	themPath, _ := tm.getThemePath(themeName)
	bgFileAbsPath := path.Join(themPath, bgFileRelPath)
	return bgFileAbsPath
}

// TODO
// func (tm *ThemeManager) getBgFileRelPath(themeName, bgFileAbsPath string) string {
// 	themPath, _ := tm.getThemePath(themeName)
// 	i := strings.Index(bgFileAbsPath, themPath)
// 	if i >= 0 {
// 		bgFileRelPath := bgFileAbsPath[i:]
// 	}
// 	// existed = isFileExists(bgFileAbsPath)
// 	return bgFileRelPath
// }

func (tm *ThemeManager) copyBgFileToThemeDir(themeName, imageFile string) (newBgFile string, err error) {
	bgFileName := tm.getNewBgFileName(imageFile)
	newBgFile = tm.getBgFileAbsPath(themeName, bgFileName)
	_, err = copyFile(newBgFile, imageFile)
	if err != nil {
		logError(err.Error()) // TODO
	}
	return
}

func (tm *ThemeManager) getNewBgFileName(imageFile string) string {
	fileName := path.Base(imageFile)
	i := strings.LastIndex(fileName, ".")
	fileExt := fileName[i:]
	return "background" + fileExt
}
