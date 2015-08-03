package appearance

import (
	"encoding/json"
	"fmt"
	"pkg.deepin.io/dde/daemon/appearance/background"
	"pkg.deepin.io/dde/daemon/appearance/fonts"
	"pkg.deepin.io/dde/daemon/appearance/subthemes"
	"pkg.deepin.io/lib/gio-2.0"
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
}

func NewManager() *Manager {
	var m = new(Manager)
	m.setting = gio.NewSettings(appearanceSchema)
	m.Theme = m.setting.GetString(gsKeyTheme)
	m.FontSize = m.setting.GetInt(gsKeyFontSize)

	m.init()

	return m
}

func (m *Manager) init() {
	dt := m.getCurrentDTheme()
	subthemes.SetGtkTheme(dt.Gtk.Id)
	subthemes.SetIconTheme(dt.Icon.Id)
	subthemes.SetCursorTheme(dt.Cursor.Id)
}

func (m *Manager) destroy() {
	m.setting.Unref()
}

func (m *Manager) doSetDTheme(id string) error {
	dt := m.getCurrentDTheme()
	if dt.Id == id {
		return nil
	}

	err := SetDTheme(id)
	if err != nil {
		return err
	}
	m.setting.SetString(gsKeyTheme, id)
	m.Theme = m.setting.GetString(gsKeyTheme)

	return m.doSetFontSize(ListDTheme().Get(id).FontSize)
}

func (m *Manager) doSetGtkTheme(value string) error {
	dt := m.getCurrentDTheme()
	if dt.Gtk.Id == value {
		return nil
	}

	if !subthemes.IsGtkTheme(value) {
		return fmt.Errorf("Invalid gtk theme '%v'", value)
	}

	return m.setDThemeByComponent(&dthemeComponent{
		Gtk:           value,
		Icon:          dt.Icon.Id,
		Cursor:        dt.Cursor.Id,
		Background:    dt.Background.URI,
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

	return m.setDThemeByComponent(&dthemeComponent{
		Gtk:           dt.Gtk.Id,
		Icon:          value,
		Cursor:        dt.Cursor.Id,
		Background:    dt.Background.URI,
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

	return m.setDThemeByComponent(&dthemeComponent{
		Gtk:           dt.Gtk.Id,
		Icon:          dt.Icon.Id,
		Cursor:        value,
		Background:    dt.Background.URI,
		StandardFont:  dt.StandardFont.Id,
		MonospaceFont: dt.MonospaceFont.Id,
	})
}

func (m *Manager) doSetBackground(value string) error {
	dt := m.getCurrentDTheme()
	if dt.Background.URI == value {
		return nil
	}

	if !background.IsBackgroundFile(value) {
		return fmt.Errorf("Invalid background file '%v'", value)
	}

	return m.setDThemeByComponent(&dthemeComponent{
		Gtk:           dt.Gtk.Id,
		Icon:          dt.Icon.Id,
		Cursor:        dt.Cursor.Id,
		Background:    value,
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

	return m.setDThemeByComponent(&dthemeComponent{
		Gtk:           dt.Gtk.Id,
		Icon:          dt.Icon.Id,
		Cursor:        dt.Cursor.Id,
		Background:    dt.Background.URI,
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

	return m.setDThemeByComponent(&dthemeComponent{
		Gtk:           dt.Gtk.Id,
		Icon:          dt.Icon.Id,
		Cursor:        dt.Cursor.Id,
		Background:    dt.Background.URI,
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

	err := fonts.SetSize(size)
	if err != nil {
		return err
	}

	m.setting.SetInt(gsKeyFontSize, m.FontSize)
	m.FontSize = m.setting.GetInt(gsKeyFontSize)
	return nil
}

func (m *Manager) getCurrentDTheme() *DTheme {
	id := m.setting.GetString(gsKeyTheme)
	dt := ListDTheme().Get(id)
	if dt != nil {
		return dt
	}

	m.doSetDTheme(dthemeDefaultId)
	return ListDTheme().Get(dthemeDefaultId)
}

func (m *Manager) setDThemeByComponent(component *dthemeComponent) error {
	id := ListDTheme().findDThemeId(component)
	if len(id) != 0 {
		return m.doSetDTheme(id)
	}

	err := doWriteCustomDTheme(component)
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
