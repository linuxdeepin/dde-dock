package appearance

import (
	"encoding/json"
	"fmt"
	"path"

	"pkg.deepin.io/dde/daemon/appearance/background"
	"pkg.deepin.io/dde/daemon/appearance/dtheme"
	"pkg.deepin.io/dde/daemon/appearance/fonts"
	"pkg.deepin.io/dde/daemon/appearance/subthemes"
	"pkg.deepin.io/lib/gio-2.0"
	"pkg.deepin.io/lib/glib-2.0"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	TypeDTheme        string = "dtheme"
	TypeGtkTheme             = "gtk"
	TypeIconTheme            = "icon"
	TypeCursorTheme          = "cursor"
	TypeBackground           = "background"
	TypeStandardFont         = "standardfont"
	TypeMonospaceFont        = "monospacefont"
	TypeFontSize             = "fontsize"
)

const (
	dthemeDefaultId = "Deepin"
	dthemeCustomId  = "Custom"

	wrapBgSchema    = "com.deepin.wrap.gnome.desktop.background"
	gnomeBgSchema   = "org.gnome.desktop.background"
	gsKeyBackground = "picture-uri"

	appearanceSchema = "com.deepin.dde.appearance"
	gsKeyTheme       = "theme"
	gsKeyFontSize    = "font-size"
)

type Manager struct {
	// Current desktop theme
	Theme string
	// Current desktop font size
	FontSize int32

	// Theme changed signal
	// ty, name
	Changed func(string, string)

	setting *gio.Settings

	wrapBgSetting  *gio.Settings
	gnomeBgSetting *gio.Settings
}

func NewManager() *Manager {
	var m = new(Manager)
	m.setting = gio.NewSettings(appearanceSchema)
	m.setPropTheme(m.setting.GetString(gsKeyTheme))
	m.setPropFontSize(m.setting.GetInt(gsKeyFontSize))

	m.wrapBgSetting, _ = dutils.CheckAndNewGSettings(wrapBgSchema)
	m.gnomeBgSetting, _ = dutils.CheckAndNewGSettings(gnomeBgSchema)

	m.init()

	return m
}

func (m *Manager) destroy() {
	m.setting.Unref()

	if m.wrapBgSetting != nil {
		m.wrapBgSetting.Unref()
	}

	if m.gnomeBgSetting != nil {
		m.gnomeBgSetting.Unref()
	}
}

func (m *Manager) init() {
	var file = path.Join(glib.GetUserConfigDir(), "fontconfig", "fonts.conf")
	if dutils.IsFileExist(file) {
		return
	}

	dt := m.getCurrentDTheme()
	if dt == nil {
		logger.Error("Not found valid dtheme")
		return
	}

	err := fonts.SetFamily(dt.StandardFont.Id, dt.MonospaceFont.Id, dt.FontSize)
	if err != nil {
		logger.Debug("[init]----------- font failed:", err)
		return
	}
}

func (m *Manager) doSetDTheme(id string) error {
	err := dtheme.SetDTheme(id)
	if err != nil {
		return err
	}

	if m.Theme == id {
		return nil
	}

	m.setPropTheme(id)
	m.setting.SetString(gsKeyTheme, id)

	return m.doSetFontSize(dtheme.ListDTheme().Get(id).FontSize)
}

func (m *Manager) doSetGtkTheme(value string) error {
	dt := m.getCurrentDTheme()
	if dt.Gtk.Id == value {
		return nil
	}

	if !subthemes.IsGtkTheme(value) {
		return fmt.Errorf("Invalid gtk theme '%v'", value)
	}

	subthemes.SetGtkTheme(value)
	return m.setDThemeByComponent(&dtheme.ThemeComponent{
		Gtk:           value,
		Icon:          dt.Icon.Id,
		Cursor:        dt.Cursor.Id,
		Background:    dt.Background.Id,
		StandardFont:  dt.StandardFont.Id,
		MonospaceFont: dt.MonospaceFont.Id,
	})
}

func (m *Manager) doSetIconTheme(value string) error {
	dt := m.getCurrentDTheme()
	if dt.Icon.Id == value {
		return nil
	}

	if !subthemes.IsIconTheme(value) {
		return fmt.Errorf("Invalid icon theme '%v'", value)
	}

	subthemes.SetIconTheme(value)
	return m.setDThemeByComponent(&dtheme.ThemeComponent{
		Gtk:           dt.Gtk.Id,
		Icon:          value,
		Cursor:        dt.Cursor.Id,
		Background:    dt.Background.Id,
		StandardFont:  dt.StandardFont.Id,
		MonospaceFont: dt.MonospaceFont.Id,
	})
}

func (m *Manager) doSetCursorTheme(value string) error {
	dt := m.getCurrentDTheme()
	if dt.Cursor.Id == value {
		return nil
	}

	if !subthemes.IsCursorTheme(value) {
		return fmt.Errorf("Invalid cursor theme '%v'", value)
	}

	subthemes.SetCursorTheme(value)
	return m.setDThemeByComponent(&dtheme.ThemeComponent{
		Gtk:           dt.Gtk.Id,
		Icon:          dt.Icon.Id,
		Cursor:        value,
		Background:    dt.Background.Id,
		StandardFont:  dt.StandardFont.Id,
		MonospaceFont: dt.MonospaceFont.Id,
	})
}

func (m *Manager) doSetBackground(value string) error {
	dt := m.getCurrentDTheme()
	if dt.Background.Id == value {
		return nil
	}

	if !background.IsBackgroundFile(value) {
		return fmt.Errorf("Invalid background file '%v'", value)
	}

	uri, err := background.ListBackground().Set(value)
	if err != nil {
		return err
	}
	return m.setDThemeByComponent(&dtheme.ThemeComponent{
		Gtk:           dt.Gtk.Id,
		Icon:          dt.Icon.Id,
		Cursor:        dt.Cursor.Id,
		Background:    uri,
		StandardFont:  dt.StandardFont.Id,
		MonospaceFont: dt.MonospaceFont.Id,
	})
}

func (m *Manager) doSetStandardFont(value string) error {
	dt := m.getCurrentDTheme()
	if dt.StandardFont.Id == value {
		return nil
	}

	if !fonts.IsFontFamily(value) {
		return fmt.Errorf("Invalid font family '%v'", value)
	}

	//fonts.SetFamily(value, dt.MonospaceFont.Id, m.FontSize)
	return m.setDThemeByComponent(&dtheme.ThemeComponent{
		Gtk:           dt.Gtk.Id,
		Icon:          dt.Icon.Id,
		Cursor:        dt.Cursor.Id,
		Background:    dt.Background.Id,
		StandardFont:  value,
		MonospaceFont: dt.MonospaceFont.Id,
	})
}

func (m *Manager) doSetMonnospaceFont(value string) error {
	dt := m.getCurrentDTheme()
	if dt.MonospaceFont.Id == value {
		return nil
	}

	if !fonts.IsFontFamily(value) {
		return fmt.Errorf("Invalid font family '%v'", value)
	}

	//fonts.SetFamily(dt.StandardFont.Id, value, m.FontSize)
	return m.setDThemeByComponent(&dtheme.ThemeComponent{
		Gtk:           dt.Gtk.Id,
		Icon:          dt.Icon.Id,
		Cursor:        dt.Cursor.Id,
		Background:    dt.Background.Id,
		StandardFont:  dt.StandardFont.Id,
		MonospaceFont: value,
	})
}

func (m *Manager) doSetFontSize(size int32) error {
	if m.FontSize == size {
		return nil
	}

	if !fonts.IsFontSizeValid(size) {
		return fmt.Errorf("Invalid font size '%v'", size)
	}

	m.setPropFontSize(size)
	m.setting.SetInt(gsKeyFontSize, size)

	if size == fonts.GetFontSize() {
		return nil
	}
	dt := m.getCurrentDTheme()
	if dt == nil {
		return fmt.Errorf("Not found valid dtheme")
	}

	return fonts.SetFamily(dt.StandardFont.Id, dt.MonospaceFont.Id, dt.FontSize)
}

func (m *Manager) getCurrentDTheme() *dtheme.DTheme {
	id := m.setting.GetString(gsKeyTheme)
	dt := dtheme.ListDTheme().Get(id)
	if dt != nil {
		return dt
	}

	m.doSetDTheme(dthemeDefaultId)
	return dtheme.ListDTheme().Get(dthemeDefaultId)
}

func (m *Manager) setDThemeByComponent(component *dtheme.ThemeComponent) error {
	id := dtheme.ListDTheme().FindDThemeId(component)
	if len(id) != 0 {
		return m.doSetDTheme(id)
	}

	err := dtheme.WriteCustomTheme(component)
	if err != nil {
		return err
	}
	return m.doSetDTheme(dthemeCustomId)
}

func (*Manager) doShow(ifc interface{}) (string, error) {
	if ifc == nil {
		return "", fmt.Errorf("Not found target")
	}
	content, err := json.Marshal(ifc)
	return string(content), err
}
