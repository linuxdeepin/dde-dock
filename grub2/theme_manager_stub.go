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
	_, ok := tm.getThemeTplFile(themeName)
	return ok
}

func (tm *ThemeManager) GetThemeLastCustomizeValues(themeName string) (background, itemColor, selectedItemColor string, ok bool) {
	jsonFile, ok := tm.getThemeTplLastJsonFile(themeName)
	if !ok {
		return tm.GetThemeDefaultCustomizeValues(themeName)
	}

	fileContent, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		logError(err.Error()) // TODO
		return "", "", "", false
	}
	background, itemColor, selectedItemColor, ok = tm.getValuesInJson(fileContent)
	background = tm.getBgFileAbsPath(themeName, background)
	if !isFileExists(background) {
		logError("theme [%s]: background file not existed", background)
	}
	return
}

func (tm *ThemeManager) GetThemeDefaultCustomizeValues(themeName string) (background, itemColor, selectedItemColor string, ok bool) {
	jsonFile, ok := tm.getThemeTplDefaultJsonFile(themeName)
	if !ok {
		logError("theme [%s]: default json data file is not existed", themeName) // TODO
		return "", "", "", false
	}

	fileContent, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		logError(err.Error()) // TODO
		return "", "", "", false
	}
	background, itemColor, selectedItemColor, ok = tm.getValuesInJson(fileContent)
	background = tm.getBgFileAbsPath(themeName, background)
	if !isFileExists(background) {
		logError("theme [%s]: background file not existed", background)
	}
	return
}

func (tm *ThemeManager) CustomTheme(themeName, background, itemColor, selectedItemColor string) bool {
	_, err := tm.copyBgFileToThemeDir(themeName, background)
	if err != nil {
		return false
	}
	background = tm.getNewBgFileName(background)

	tplData := make(map[string]string)
	tplData[_THEME_TPL_KEY_BACKGROUND] = background
	tplData[_THEME_TPL_KEY_ITEM_COLOR] = itemColor
	tplData[_THEME_TPL_KEY_SELECTED_ITEM_COLOR] = selectedItemColor

	tplFile, ok := tm.getThemeTplFile(themeName)
	if !ok {
		logError("theme [%s]: template file is not existed", themeName) // TODO
	}

	tplFileContent, err := ioutil.ReadFile(tplFile)
	if err != nil {
		logError(err.Error()) // TODO
		return false
	}

	themeFileContent, err := tm.getCustomizedThemeContent(tplFileContent, tplData)
	themeMainFile, _ := tm.getThemeMainFile(themeName)
	err = ioutil.WriteFile(themeMainFile, themeFileContent, 0644)
	if err != nil {
		logError(err.Error()) // TODO
		return false
	}

	// store the customized key-values to json file
	jsonContent, err := json.Marshal(tplData)
	if err != nil {
		logError(err.Error()) // TODO
		return false
	}
	lastJsonFile, _ := tm.getThemeTplLastJsonFile(themeName)
	err = ioutil.WriteFile(lastJsonFile, jsonContent, 0644)
	if err != nil {
		logError(err.Error()) // TODO
		return false
	}

	return true
}
