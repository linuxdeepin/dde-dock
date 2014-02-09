package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"
	"text/template"
)

var _THEME_TEMPLATOR = template.New("theme-templator")

type TplValues struct {
	Background, ItemColor, SelectedItemColor string
}
type TplJsonData struct {
	DefaultTplValue, LastTplValue TplValues
}

var themeId int32 = 0

type Theme struct {
	id        int32
	tm        *ThemeManager
	themePath string
	mainFile  string
	tplFile   string
	jsonFile  string
	relBgFile string // relative background file path

	Name              string
	Customizable      bool
	Background        string `access:"readwrite"` // absolute background file path
	ItemColor         string `access:"readwrite"`
	SelectedItemColor string `access:"readwrite"`
}

func NewTheme(tm *ThemeManager, name string) (*Theme, error) {
	theme := &Theme{}
	theme.id = themeId
	themeId++
	theme.tm = tm
	theme.Name = name
	theme.themePath, _ = tm.getThemePath(name)
	theme.mainFile, _ = tm.getThemeMainFile(name)
	if path, ok := tm.getThemeTplFile(name); ok {
		theme.tplFile = path
	}
	if path, ok := tm.getThemeJsonFile(name); ok {
		theme.jsonFile = path
	}

	// TODO
	theme.Customizable = tm.isThemeCustomizable(name)
	if theme.Customizable {
		tplJsonData, err := theme.getThemeTplJsonData()
		if err != nil {
			return nil, err
		}

		theme.relBgFile = tplJsonData.LastTplValue.Background
		theme.makeAbsBgFile()
		theme.ItemColor = tplJsonData.LastTplValue.ItemColor
		theme.SelectedItemColor = tplJsonData.LastTplValue.SelectedItemColor
	}

	return theme, nil
}

// TODO
func (theme *Theme) setBackground(background string) {
	// copy background file to theme dir if need
	theme.copyBgFileToThemeDir(background)
	theme.relBgFile = theme.getNewBgFileName(background)
	theme.customTheme()
}

func (theme *Theme) setItemColor(itemColor string) {
	theme.customTheme()
}

func (theme *Theme) setSelectedItemColor(selectedItemColor string) {
	theme.customTheme()
}

// TODO Update variable 'Background' which means the absolute background file path
func (theme *Theme) makeAbsBgFile() {
	theme.Background = theme.getBgFileAbsPath(theme.relBgFile)
	if !isFileExists(theme.Background) {
		logError("theme [%s]: background file is not exists, %s", theme.Name, theme.Background) // TODO
	}
}

func (theme *Theme) copyBgFileToThemeDir(imageFile string) (newBgFile string, err error) {
	bgFileName := theme.getNewBgFileName(imageFile)
	newBgFile = theme.getBgFileAbsPath(bgFileName)
	_, err = copyFile(newBgFile, imageFile)
	if err != nil {
		logError(err.Error()) // TODO
	}
	return
}

func (theme *Theme) getBgFileAbsPath(bgFileRelPath string) string {
	bgFileAbsPath := path.Join(theme.themePath, bgFileRelPath)
	return bgFileAbsPath
}

func (theme *Theme) getNewBgFileName(imageFile string) string {
	fileName := path.Base(imageFile)
	i := strings.LastIndex(fileName, ".")
	fileExt := fileName[i:]
	return "background" + fileExt
}

// TODO
func (theme *Theme) getThemeTplJsonData() (*TplJsonData, error) {
	fileContent, err := ioutil.ReadFile(theme.jsonFile)
	if err != nil {
		logError(err.Error()) // TODO
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
		logError(err.Error()) // TODO
		return nil, err
	}
	return tplJsonData, nil
}

func (theme *Theme) customTheme() {
	tplJsonData, err := theme.getThemeTplJsonData()
	if err != nil {
		panic(err)
	}
	tplJsonData.LastTplValue.Background = theme.relBgFile
	tplJsonData.LastTplValue.ItemColor = theme.ItemColor
	tplJsonData.LastTplValue.SelectedItemColor = theme.SelectedItemColor

	// generate a new theme.txt from template
	tplFileContent, err := ioutil.ReadFile(theme.tplFile)
	if err != nil {
		logError(err.Error()) // TODO
		panic(err)
	}
	themeFileContent, err := theme.getCustomizedThemeContent(tplFileContent, tplJsonData.LastTplValue)
	err = ioutil.WriteFile(theme.mainFile, themeFileContent, 0644)
	if err != nil {
		logError(err.Error()) // TODO
		panic(err)
	}

	// store the customized key-values to json file
	jsonContent, err := json.Marshal(tplFileContent)
	if err != nil {
		logError(err.Error()) // TODO
		panic(err)
	}
	err = ioutil.WriteFile(theme.jsonFile, jsonContent, 0644)
	if err != nil {
		logError(err.Error()) // TODO
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
