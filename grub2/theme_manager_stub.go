package main

import (
	"dlib/dbus"
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

func (tm *ThemeManager) OnPropertiesChanged(name string, oldv interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logError("%v", err) // TODO
		}
	}()
	switch name {
	case "EnabledTheme":
		if len(tm.EnabledTheme) == 0 {
			tm.setEnabledThemeMainFile("")
		} else {
			_, theme := tm.getTheme(tm.EnabledTheme)
			if theme == nil {
				panic(newError("theme not found: %s", tm.EnabledTheme))
			}
			tm.setEnabledThemeMainFile(theme.mainFile)
		}
	}
}

func (tm *ThemeManager) InstallTheme(archive string) bool {
	themePathInZip, err := findFileInTarGz(archive, _THEME_MAIN_FILE)
	if err != nil {
		logError("install theme %s failed: %v", archive, err) // TODO
		return false
	}
	themePathPrefix := path.Dir(themePathInZip)
	unTarGz(archive, _THEME_DIR, themePathPrefix)

	tm.load()
	// TODO append theme object to list
	// theme, err := NewTheme(tm, themeName)
	// if err == nil {
	// tm.themes = append(tm.themes, theme)
	// }
	// tm.makeThemeNames()

	logInfo("install theme success: %s", archive)
	return true
}

// TODO
func (tm *ThemeManager) UninstallTheme(themeName string) bool {
	_, theme := tm.getTheme(themeName)
	err := os.RemoveAll(theme.themePath)
	if err != nil {
		return false
	}

	tm.load()
	// TODO delete theme object from list
	// i, theme := tm.getTheme(themeName)
	// copy(tm.themes[i:], tm.themes[i+1:])
	// tm.themes[len(tm.themes)-1] = nil
	// tm.themes = tm.themes[:len(tm.themes)-1]
	// dbus.UnInstallObject(theme)
	// tm.makeThemeNames()

	logInfo("uninstall theme success: %s", themeName)
	return true
}

func (tm *ThemeManager) GetThemeId(themeName string) int32 {
	for _, t := range tm.themes {
		if t.Name == themeName {
			return t.id
		}
	}
	return -1
}
