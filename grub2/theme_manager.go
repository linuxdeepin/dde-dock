package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"text/template"
)

const (
	_THEME_DIR           = "/boot/grub/themes"
	_THEME_MAIN_FILE     = "theme.txt"
	_THEME_TPL_FILE      = "theme.tpl"
	_THEME_TPL_JSON_FILE = "theme_tpl.json" // json stores the key-values for template file
)

var _THEME_TEMPLATOR = template.New("theme-templator")

type TplValues struct {
	Background, ItemColor, SelectedItemColor string
}
type TplJsonData struct {
	DefaultTplValue, LastTplValue TplValues
}

type ThemeManager struct {
	themes               []*Theme // TODO
	enabledThemeMainFile string

	ThemeNames []string // TODO
}

func NewThemeManager() *ThemeManager {
	tm := &ThemeManager{}
	tm.enabledThemeMainFile = ""
	return tm
}

func (tm *ThemeManager) load() {
	tm.themes = make([]*Theme, 0)
	files, err := ioutil.ReadDir(_THEME_DIR)
	if err == nil {
		for _, f := range files {
			if f.IsDir() && tm.isThemeValid(f.Name()) {
				theme := NewTheme(tm, f.Name())
				tm.themes = append(tm.themes, theme)
			}
		}
	}
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
	p, err := findFileInTarGz(archive, _THEME_MAIN_FILE)
	if err != nil {
		return false
	}

	// check theme path level in archive
	p = path.Clean(p)
	if getPathLevel(path.Base(p)) != 1 {
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

func (tm *ThemeManager) getThemeTplJsonFile(themeName string) (file string, existed bool) {
	file = path.Join(_THEME_DIR, themeName, _THEME_TPL_JSON_FILE)
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

func (tm *ThemeManager) getThemeTplJsonData(themeName string) (*TplJsonData, error) {
	jsonFile, ok := tm.getThemeTplJsonFile(themeName)
	if !ok {
		err := errors.New(fmt.Sprintf("theme [%s]: json file for template is not exists", jsonFile))
		logError(err.Error()) // TODO
		return nil, err
	}

	fileContent, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		logError(err.Error()) // TODO
		return nil, err
	}

	tplJsonData, err := tm.getTplJsonData(fileContent)
	if err != nil {
		return nil, err
	}
	return tplJsonData, nil
}

func (tm *ThemeManager) getTplJsonData(fileContent []byte) (*TplJsonData, error) {
	tplJsonData := &TplJsonData{}
	err := json.Unmarshal(fileContent, tplJsonData)
	if err != nil {
		logError(err.Error()) // TODO
		return nil, err
	}
	return tplJsonData, nil
}

func (tm *ThemeManager) getBgFileAbsPath(themeName, bgFileRelPath string) string {
	themPath, _ := tm.getThemePath(themeName)
	bgFileAbsPath := path.Join(themPath, bgFileRelPath)
	return bgFileAbsPath
}

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
