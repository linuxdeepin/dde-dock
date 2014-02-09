package main

import (
	"dlib/dbus"
	"io/ioutil"
	"path"
	"strings"
)

const (
	_THEME_DIR           = "/boot/grub/themes"
	_THEME_MAIN_FILE     = "theme.txt"
	_THEME_TPL_FILE      = "theme.tpl"
	_THEME_TPL_JSON_FILE = "theme_tpl.json" // json stores the key-values for template file
)

type ThemeManager struct {
	themes               []*Theme
	enabledThemeMainFile string // TODO

	ThemeNames   []string
	EnabledTheme string // TODO
}

func NewThemeManager() *ThemeManager {
	tm := &ThemeManager{}
	tm.enabledThemeMainFile = ""
	return tm
}

func (tm *ThemeManager) load() {
	// TODO clear themes load last time, uninstall dbus interface
	for _, t := range tm.themes {
		dbus.UnInstallObject(t)
	}
	tm.themes = make([]*Theme, 0)

	themeId = 0
	files, err := ioutil.ReadDir(_THEME_DIR)
	if err == nil {
		for _, f := range files {
			if f.IsDir() && tm.isThemeValid(f.Name()) {
				theme, err := NewTheme(tm, f.Name())
				if err == nil {
					tm.themes = append(tm.themes, theme)

					err := dbus.InstallOnSystem(theme)
					if err != nil {
						panic(err)
					}

				}
				logInfo("found theme: %s", theme.Name)
			}
		}
	}
	tm.makeThemeNames()
}

// Update variable 'ThemeNames'
func (tm *ThemeManager) makeThemeNames() {
	tm.ThemeNames = make([]string, 0)
	for _, t := range tm.themes {
		tm.ThemeNames = append(tm.ThemeNames, t.Name)
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

func (tm *ThemeManager) isThemeCustomizable(themeName string) bool {
	_, okTpl := tm.getThemeTplFile(themeName)
	_, okJson := tm.getThemeTplJsonFile(themeName)
	return okTpl && okJson
}

// TODO remove
func (tm *ThemeManager) getThemeName(themeMainFile string) string {
	if len(themeMainFile) == 0 {
		return ""
	}
	return path.Base(path.Dir(themeMainFile))
}

// TODO remove
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

func (tm *ThemeManager) getTheme(themeName string) (int, *Theme) {
	for i, t := range tm.themes {
		if t.Name == themeName {
			return i, t
		}
	}
	return -1, nil
}
