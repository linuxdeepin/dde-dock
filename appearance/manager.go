/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package appearance

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/user"
	"strconv"
	"sync"
	"time"

	"gir/gio-2.0"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.accounts"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.sessionmanager"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.wm"
	"pkg.deepin.io/dde/api/theme_thumb"
	"pkg.deepin.io/dde/daemon/appearance/background"
	"pkg.deepin.io/dde/daemon/appearance/fonts"
	"pkg.deepin.io/dde/daemon/appearance/subthemes"
	ddbus "pkg.deepin.io/dde/daemon/dbus"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/gsprop"
	"pkg.deepin.io/lib/dbusutil/proxy"
	"pkg.deepin.io/lib/fsnotify"
	"pkg.deepin.io/lib/strv"
	dutils "pkg.deepin.io/lib/utils"
)

// The supported types
const (
	TypeGtkTheme          = "gtk"
	TypeIconTheme         = "icon"
	TypeCursorTheme       = "cursor"
	TypeBackground        = "background"
	TypeGreeterBackground = "greeterbackground"
	TypeStandardFont      = "standardfont"
	TypeMonospaceFont     = "monospacefont"
	TypeFontSize          = "fontsize"
)

const (
	wrapBgSchema    = "com.deepin.wrap.gnome.desktop.background"
	gnomeBgSchema   = "org.gnome.desktop.background"
	gsKeyBackground = "picture-uri"

	appearanceSchema    = "com.deepin.dde.appearance"
	gsKeyGtkTheme       = "gtk-theme"
	gsKeyIconTheme      = "icon-theme"
	gsKeyCursorTheme    = "cursor-theme"
	gsKeyFontStandard   = "font-standard"
	gsKeyFontMonospace  = "font-monospace"
	gsKeyFontSize       = "font-size"
	gsKeyBackgroundURIs = "background-uris"

	defaultStandardFont   = "Noto Sans"
	defaultMonospaceFont  = "Noto Mono"
	defaultFontConfigFile = "/usr/share/deepin-default-settings/fontconfig.json"

	dbusServiceName = "com.deepin.daemon.Appearance"
	dbusPath        = "/com/deepin/daemon/Appearance"
	dbusInterface   = dbusServiceName
)

// Manager shows current themes and fonts settings, emit 'Changed' signal if modified
// if themes list changed will emit 'Refreshed' signal
type Manager struct {
	service *dbusutil.Service
	sigLoop *dbusutil.SignalLoop

	GtkTheme      gsprop.String
	IconTheme     gsprop.String
	CursorTheme   gsprop.String
	Background    gsprop.String
	StandardFont  gsprop.String
	MonospaceFont gsprop.String

	FontSize gsprop.Double `prop:"access:rw"`

	userObj   *accounts.User
	imageBlur *accounts.ImageBlur
	xSettings *sessionmanager.XSettings

	setting        *gio.Settings
	wrapBgSetting  *gio.Settings
	gnomeBgSetting *gio.Settings

	defaultFontConfig   DefaultFontConfig
	defaultFontConfigMu sync.Mutex

	watcher    *fsnotify.Watcher
	endWatcher chan struct{}

	currentDesktopBgs []string
	currentGreeterBg  string

	wm *wm.Wm

	signals *struct {
		// Theme setting changed
		Changed struct {
			type0 string
			value string
		}

		// Theme list refreshed
		Refreshed struct {
			type0 string
		}
	}

	methods *struct {
		Delete         func() `in:"type,name"`
		GetScaleFactor func() `out:"scale_factor"`
		List           func() `in:"type" out:"list"`
		Set            func() `in:"type,value"`
		SetScaleFactor func() `in:"scale_factor"`
		Show           func() `in:"type,names" out:"detail"`
		Thumbnail      func() `in:"type,name" out:"file"`
	}
}

// NewManager will create a 'Manager' object
func newManager(service *dbusutil.Service) *Manager {
	var m = new(Manager)
	m.service = service
	m.setting = gio.NewSettings(appearanceSchema)
	m.wrapBgSetting = gio.NewSettings(wrapBgSchema)

	m.GtkTheme.Bind(m.setting, gsKeyGtkTheme)
	m.IconTheme.Bind(m.setting, gsKeyIconTheme)
	m.CursorTheme.Bind(m.setting, gsKeyCursorTheme)
	m.StandardFont.Bind(m.setting, gsKeyFontStandard)
	m.MonospaceFont.Bind(m.setting, gsKeyFontMonospace)
	m.Background.Bind(m.wrapBgSetting, gsKeyBackground)
	m.FontSize.Bind(m.setting, gsKeyFontSize)

	m.gnomeBgSetting, _ = dutils.CheckAndNewGSettings(gnomeBgSchema)

	var err error
	m.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		logger.Warning("New file watcher failed:", err)
	} else {
		m.endWatcher = make(chan struct{})
	}

	return m
}

func (m *Manager) initCurrentBgs() {
	m.currentDesktopBgs = m.setting.GetStrv(gsKeyBackgroundURIs)
	greeterBg, err := m.userObj.GreeterBackground().Get(0)
	if err == nil {
		m.currentGreeterBg = greeterBg
	} else {
		logger.Warning(err)
	}
}

func (m *Manager) isBgInUse(file string) bool {
	if file == m.currentGreeterBg {
		return true
	}

	for _, bg := range m.currentDesktopBgs {
		if bg == file {
			return true
		}
	}
	return false
}

func (m *Manager) listBackground() background.Backgrounds {
	origin := background.ListBackground()
	result := make(background.Backgrounds, len(origin))

	for idx, bg := range origin {
		var deletable bool
		if bg.Deletable {
			// custom
			if !m.isBgInUse(bg.Id) {
				deletable = true
			}
		}
		result[idx] = &background.Background{
			Id:        bg.Id,
			Deletable: deletable,
		}
	}
	return result
}

func (m *Manager) destroy() {
	m.sigLoop.Stop()
	m.xSettings.RemoveHandler(proxy.RemoveAllHandlers)

	if m.setting != nil {
		m.setting.Unref()
		m.setting = nil
	}

	if m.wrapBgSetting != nil {
		m.wrapBgSetting.Unref()
		m.wrapBgSetting = nil
	}

	if m.gnomeBgSetting != nil {
		m.gnomeBgSetting.Unref()
		m.gnomeBgSetting = nil
	}

	if m.watcher != nil {
		close(m.endWatcher)
		m.watcher.Close()
		m.watcher = nil
	}

	m.endCursorChangedHandler()
}

// resetFonts reset StandardFont and MonospaceFont
func (m *Manager) resetFonts() {
	defaultStandardFont, defaultMonospaceFont := m.getDefaultFonts()
	logger.Debugf("getDefaultFonts standard: %q, mono: %q",
		defaultStandardFont, defaultMonospaceFont)
	if defaultStandardFont != m.StandardFont.Get() {
		m.StandardFont.Set(defaultStandardFont)
	}

	if defaultMonospaceFont != m.MonospaceFont.Get() {
		m.MonospaceFont.Set(defaultMonospaceFont)
	}

	err := fonts.SetFamily(defaultStandardFont, defaultMonospaceFont,
		m.FontSize.Get())
	if err != nil {
		logger.Debug("resetFonts fonts.SetFamily failed", err)
		return
	}
	m.checkFontConfVersion()
}

func (m *Manager) initUserObj(systemConn *dbus.Conn) {
	cur, err := user.Current()
	if err != nil {
		logger.Warning("failed to get current user:", err)
		return
	}
	m.userObj, err = ddbus.NewUserByUid(systemConn, cur.Uid)
	if err != nil {
		logger.Warning("failed to new user object", err)
		return
	}

	// sync desktop backgrounds
	userBackgrounds, err := m.userObj.DesktopBackgrounds().Get(0)
	if err != nil {
		logger.Warning(err)
		return
	}

	gsBackgrounds := m.setting.GetStrv(gsKeyBackgroundURIs)
	if !strv.Strv(userBackgrounds).Equal(gsBackgrounds) {
		m.setDesktopBackgrounds(gsBackgrounds)
	}
}

func (m *Manager) init() {
	background.SetCustomWallpaperDeleteCallback(func(file string) {
		logger.Debug("imageBlur delete", file)
		err := m.imageBlur.Delete(0, file)
		if err != nil {
			logger.Warning("imageBlur delete err:", err)
		}
	})

	sessionBus := m.service.Conn()
	m.wm = wm.NewWm(sessionBus)

	m.xSettings = sessionmanager.NewXSettings(sessionBus)
	theme_thumb.Init(m.getScaleFactor())

	m.sigLoop = dbusutil.NewSignalLoop(sessionBus, 10)
	m.xSettings.InitSignalExt(m.sigLoop, true)
	m.xSettings.ConnectSetScaleFactorDone(m.handleSetScaleFactorDone)
	m.sigLoop.Start()

	err := m.loadDefaultFontConfig(defaultFontConfigFile)
	if err != nil {
		logger.Warning("load default font config failed:", err)
	}

	m.doSetGtkTheme(m.GtkTheme.Get())
	m.doSetIconTheme(m.IconTheme.Get())
	m.doSetCursorTheme(m.CursorTheme.Get())

	// Init theme list
	time.AfterFunc(time.Second*10, func() {
		if !dutils.IsFileExist(fonts.DeepinFontConfig) {
			m.resetFonts()
		} else {
			m.correctFontName()
		}

		subthemes.ListGtkTheme()
		subthemes.ListIconTheme()
		subthemes.ListCursorTheme()
		background.ListBackground()
		fonts.GetFamilyTable()

		setDQtTheme(dQtFile, dQtSectionTheme,
			[]string{
				dQtKeyIcon,
				dQtKeyFont,
				dQtKeyMonoFont,
				dQtKeyFontSize},
			[]string{
				m.IconTheme.Get(),
				m.StandardFont.Get(),
				m.MonospaceFont.Get(),
				strconv.FormatFloat(m.FontSize.Get(), 'f', 1, 64)})
		err := saveDQtTheme(dQtFile)
		if err != nil {
			logger.Warning("Failed to save qt theme:", err)
			return
		}
	})

	systemConn, err := dbus.SystemBus()
	if err != nil {
		logger.Warning(err)
		return
	}

	m.initUserObj(systemConn)
	m.initCurrentBgs()
	m.imageBlur = accounts.NewImageBlur(systemConn)
}

func (m *Manager) correctFontName() {
	defaultStandardFont, defaultMonospaceFont := m.getDefaultFonts()

	var changed = false
	table := fonts.GetFamilyTable()
	stand := table.GetFamily(m.StandardFont.Get())
	if stand != nil {
		// for virtual font
		if stand.Id != m.StandardFont.Get() {
			changed = true
			m.StandardFont.Set(stand.Id)
		}
	} else {
		changed = true
		m.StandardFont.Set(defaultStandardFont)
	}

	mono := table.GetFamily(m.MonospaceFont.Get())
	if mono != nil {
		if mono.Id != m.MonospaceFont.Get() {
			changed = true
			m.MonospaceFont.Set(mono.Id)
		}
	} else {
		changed = true
		m.MonospaceFont.Set(defaultMonospaceFont)
	}

	if !changed && m.checkFontConfVersion() {
		return
	}

	err := fonts.SetFamily(m.StandardFont.Get(), m.MonospaceFont.Get(),
		m.FontSize.Get())
	if err != nil {
		logger.Debug("[correctFontName]-----------set font failed:", err)
		return
	}
}

func (m *Manager) doSetGtkTheme(value string) error {
	if !subthemes.IsGtkTheme(value) {
		return fmt.Errorf("invalid gtk theme '%v'", value)
	}

	return subthemes.SetGtkTheme(value)
}

func (m *Manager) doSetIconTheme(value string) error {
	if !subthemes.IsIconTheme(value) {
		return fmt.Errorf("invalid icon theme '%v'", value)
	}

	err := subthemes.SetIconTheme(value)
	if err != nil {
		return err
	}

	return m.writeDQtTheme(dQtKeyIcon, value)
}

func (m *Manager) doSetCursorTheme(value string) error {
	if !subthemes.IsCursorTheme(value) {
		return fmt.Errorf("invalid cursor theme '%v'", value)
	}

	return subthemes.SetCursorTheme(value)
}

func (m *Manager) doSetBackground(value string) error {
	logger.Debugf("call doSetBackground %q", value)
	if !background.IsBackgroundFile(value) {
		return errors.New("invalid background")
	}

	file, err := background.Prepare(value)
	if err != nil {
		return err
	}
	logger.Debug("prepare result:", file)
	uri := dutils.EncodeURI(file, dutils.SCHEME_FILE)
	m.wm.ChangeCurrentWorkspaceBackground(dbus.FlagNoAutoStart, uri)

	_, err = m.imageBlur.Get(0, file)
	if err != nil {
		logger.Warning("call imageBlur.Get err:", err)
	}
	return nil
}

func (m *Manager) doSetGreeterBackground(value string) error {
	if m.userObj == nil {
		return fmt.Errorf("create user object failed")
	}

	return m.userObj.SetGreeterBackground(0, value)
}

func (m *Manager) doSetStandardFont(value string) error {
	if !fonts.IsFontFamily(value) {
		return fmt.Errorf("invalid font family '%v'", value)
	}

	err := fonts.SetFamily(value, m.MonospaceFont.Get(), m.FontSize.Get())
	if err != nil {
		return err
	}
	return m.writeDQtTheme(dQtKeyFont, value)
}

func (m *Manager) doSetMonnospaceFont(value string) error {
	if !fonts.IsFontFamily(value) {
		return fmt.Errorf("invalid font family '%v'", value)
	}

	err := fonts.SetFamily(m.StandardFont.Get(), value, m.FontSize.Get())
	if err != nil {
		return err
	}
	return m.writeDQtTheme(dQtKeyMonoFont, value)
}

func (m *Manager) doSetFontSize(size float64) error {
	if !fonts.IsFontSizeValid(size) {
		logger.Debug("[doSetFontSize] invalid size:", size)
		return fmt.Errorf("invalid font size '%v'", size)
	}

	err := fonts.SetFamily(m.StandardFont.Get(), m.MonospaceFont.Get(), size)
	if err != nil {
		return err
	}

	return m.writeDQtTheme(dQtKeyFontSize, strconv.FormatFloat(size, 'f', 1, 64))
}

func (*Manager) doShow(ifc interface{}) (string, error) {
	if ifc == nil {
		return "", fmt.Errorf("not found target")
	}
	content, err := json.Marshal(ifc)
	return string(content), err
}

func (m *Manager) loadDefaultFontConfig(filename string) error {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var defaultFontConfig DefaultFontConfig
	if err := json.Unmarshal(contents, &defaultFontConfig); err != nil {
		return err
	}

	m.defaultFontConfigMu.Lock()
	m.defaultFontConfig = defaultFontConfig
	m.defaultFontConfigMu.Unlock()

	logger.Debugf("load default font config ok %#v", defaultFontConfig)
	return nil
}

func (m *Manager) getDefaultFonts() (standard string, monospace string) {
	m.defaultFontConfigMu.Lock()
	cfg := m.defaultFontConfig
	m.defaultFontConfigMu.Unlock()

	if cfg == nil {
		return defaultStandardFont, defaultMonospaceFont
	}
	return cfg.Get()
}

func (m *Manager) writeDQtTheme(key, value string) error {
	setDQtTheme(dQtFile, dQtSectionTheme,
		[]string{key}, []string{value})
	return saveDQtTheme(dQtFile)
}

func (m *Manager) setDesktopBackgrounds(val []string) {
	if m.userObj != nil {
		err := m.userObj.SetDesktopBackgrounds(0, val)
		if err != nil {
			logger.Warning("call userObj.SetDesktopBackgrounds err:", err)
		}
	}
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}
