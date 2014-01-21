package main

import (
	"dlib/dbus"
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
	return tm.getThemeName(tm.enabledThemeMainFile)
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
	unTarGz(archive, _THEME_DIR, path.Dir(themePathInZip))
	return true
}

func (tm *ThemeManager) UninstallTheme(themeName string) bool {
	themePath, ok := tm.getThemePath(themeName)
	if !ok {
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

// TODO
func (tm *ThemeManager) GetThemeCustomizedValues(themeName string) (background, itemColor, selectedItemColor string) {
	return "", "", ""
}

// TODO
func (tm *ThemeManager) GetThemeCustomizedDefaultValues(themeName string) (background, itemColor, selectedItemColor string) {
	return "", "", ""
	// return true
}

// TODO
func (tm *ThemeManager) CustomTheme(themeName, background, itemColor, selectedItemColor string) bool {
	return true
}
