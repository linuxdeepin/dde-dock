package main

import (
// "text/template"
)

// const (
// 	_THEME_MAIN_FILE     = "theme.txt"
// 	_THEME_TPL_FILE      = "theme.tpl"
// 	_THEME_TPL_JSON_FILE = "theme_tpl.json" // json stores the key-values for template file
// )

// var _THEME_TEMPLATOR = template.New("theme-templator")

// type TplValues struct {
// 	Background, ItemColor, SelectedItemColor string
// }
// type TplJsonData struct {
// 	DefaultTplValue, LastTplValue TplValues
// }

type Theme struct {
	mainFile string
	tplFile  string
	jsonFile string

	Name              string
	Customizable      bool
	Background        string `access:"readwrite"`
	ItemColor         string `access:"readwrite"`
	SelectedItemColor string `access:"readwrite"`
}

func NewTheme(tm *ThemeManager, name string) *Theme {
	theme := &Theme{Name: name}
	theme.Customizable = tm.IsThemeCustomizable(name) // TODO
	theme.mainFile, _ = tm.getThemeMainFile(name)
	if path, ok := tm.getThemeTplFile(name); ok {
		theme.tplFile = path
	}
	if path, ok := tm.getThemeTplJsonFile(name); ok {
		theme.jsonFile = path
	}

	return theme
}

// TODO
func (theme *Theme) setBackground(background string) {
}

// TODO
func (theme *Theme) setItemColor(itemColor string) {
}

// TODO
func (theme *Theme) setSelectedItemColor(selectedItemColor string) {
}
