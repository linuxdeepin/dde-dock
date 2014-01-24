package main

import (
	"dlib/dbus"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

const (
	_GRUB2_THEME_MANAGER_DEST = "com.deepin.daemon.Grub2"
	_GRUB2_THEME_MANAGER_PATH = "/com/deepin/daemon/Grub2/ThemeManager"
	_GRUB2_THEME_MANAGER_IFC  = "com.deepin.daemon.Grub2.ThemeManager"
)

func (tm *ThemeManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_GRUB2_THEME_MANAGER_DEST,
		_GRUB2_THEME_MANAGER_PATH,
		_GRUB2_THEME_MANAGER_IFC,
	}
}

// TODO
func (tm *ThemeManager) GetInstalledThemes() []string {
	themes := make([]string, 0)
	files, err := ioutil.ReadDir(_THEME_DIR)
	if err == nil {
		for _, f := range files {
			if f.IsDir() && tm.isThemeValid(f.Name()) {
				themes = append(themes, f.Name())
			}
		}
	}
	return themes
}

func (tm *ThemeManager) GetEnabledTheme() string {
	return tm.getThemeName(tm.getEnabledThemeMainFile())
}

func (tm *ThemeManager) EnableTheme(themeName string) bool {
	if tm.isThemeValid(themeName) {
		file, ok := tm.getThemeMainFile(themeName)
		if ok {
			tm.setEnabledThemeMainFile(file)
			return true
		}
	}
	return false
}

func (tm *ThemeManager) DisableTheme() {
	tm.enabledThemeMainFile = ""
}

func (tm *ThemeManager) InstallTheme(archive string) bool {
	themePathInZip, err := findFileInTarGz(archive, _THEME_MAIN_FILE)
	if err != nil {
		logError("install theme %s failed: %v", archive, err) // TODO
		return false
	}
	themePathPrefix := path.Dir(themePathInZip)
	unTarGz(archive, _THEME_DIR, themePathPrefix)
	return true
}

func (tm *ThemeManager) UninstallTheme(themeName string) bool {
	themePath, exists := tm.getThemePath(themeName)
	if !exists {
		return false
	}
	err := os.RemoveAll(themePath)
	if err != nil {
		return false
	}
	return true
}

func (tm *ThemeManager) IsThemeCustomizable(themeName string) bool {
	_, okTpl := tm.getThemeTplFile(themeName)
	_, okJson := tm.getThemeTplJsonFile(themeName)
	return okTpl && okJson
}

func (tm *ThemeManager) GetThemeLastCustomizeValues(themeName string) (background, itemColor, selectedItemColor string, ok bool) {
	tplJsonData, err := tm.getThemeTplJsonData(themeName)
	if err != nil {
		return "", "", "", false
	}

	background = tplJsonData.LastTplValue.Background
	itemColor = tplJsonData.LastTplValue.ItemColor
	selectedItemColor = tplJsonData.LastTplValue.SelectedItemColor
	return background, itemColor, selectedItemColor, true

	// return background's absolute path
	background = tm.getBgFileAbsPath(themeName, background)
	if !isFileExists(background) {
		logError("theme [%s]: background file is not exists", background) // TODO
	}
	ok = true
	return
}

func (tm *ThemeManager) GetThemeDefaultCustomizeValues(themeName string) (background, itemColor, selectedItemColor string, ok bool) {
	tplJsonData, err := tm.getThemeTplJsonData(themeName)
	if err != nil {
		return "", "", "", false
	}

	background = tplJsonData.DefaultTplValue.Background
	itemColor = tplJsonData.DefaultTplValue.ItemColor
	selectedItemColor = tplJsonData.DefaultTplValue.SelectedItemColor
	return background, itemColor, selectedItemColor, true

	// return background's absolute path
	background = tm.getBgFileAbsPath(themeName, background)
	if !isFileExists(background) {
		logError("theme [%s]: background file is not exists", background) // TODO
	}
	ok = true
	return
}

func (tm *ThemeManager) CustomTheme(themeName, background, itemColor, selectedItemColor string) bool {
	// copy background file to theme dir if need
	_, err := tm.copyBgFileToThemeDir(themeName, background)
	if err != nil {
		return false
	}
	background = tm.getNewBgFileName(background)

	tplJsonData, err := tm.getThemeTplJsonData(themeName)
	if err != nil {
		return false
	}
	tplJsonData.LastTplValue.Background = background
	tplJsonData.LastTplValue.ItemColor = itemColor
	tplJsonData.LastTplValue.SelectedItemColor = selectedItemColor

	// generate a new theme.txt from template
	tplFile, ok := tm.getThemeTplFile(themeName)
	if !ok {
		logError("theme [%s]: template file is not existed", themeName) // TODO
	}

	tplFileContent, err := ioutil.ReadFile(tplFile)
	if err != nil {
		logError(err.Error()) // TODO
		return false
	}

	themeFileContent, err := tm.getCustomizedThemeContent(tplFileContent, tplJsonData.LastTplValue)
	themeMainFile, _ := tm.getThemeMainFile(themeName)
	err = ioutil.WriteFile(themeMainFile, themeFileContent, 0644)
	if err != nil {
		logError(err.Error()) // TODO
		return false
	}

	// store the customized key-values to json file
	jsonContent, err := json.Marshal(tplFileContent)
	if err != nil {
		logError(err.Error()) // TODO
		return false
	}
	jsonFile, _ := tm.getThemeTplJsonFile(themeName)
	err = ioutil.WriteFile(jsonFile, jsonContent, 0644)
	if err != nil {
		logError(err.Error()) // TODO
		return false
	}

	return true
}
