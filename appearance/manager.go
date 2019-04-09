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
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.accounts"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.sessionmanager"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.wm"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.login1"
	"github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/randr"
	"pkg.deepin.io/dde/api/theme_thumb"
	"pkg.deepin.io/dde/daemon/appearance/background"
	"pkg.deepin.io/dde/daemon/appearance/fonts"
	"pkg.deepin.io/dde/daemon/appearance/subthemes"
	"pkg.deepin.io/dde/daemon/common/dsync"
	ddbus "pkg.deepin.io/dde/daemon/dbus"
	"pkg.deepin.io/dde/daemon/session/common"
	"pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/gsprop"
	"pkg.deepin.io/lib/dbusutil/proxy"
	"pkg.deepin.io/lib/fsnotify"
	"pkg.deepin.io/lib/strv"
	dutils "pkg.deepin.io/lib/utils"
	"pkg.deepin.io/lib/xdg/basedir"
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

	appearanceSchema        = "com.deepin.dde.appearance"
	xSettingsSchema         = "com.deepin.xsettings"
	gsKeyIndividualScaling  = "individual-scaling"
	gsKeyGtkTheme           = "gtk-theme"
	gsKeyIconTheme          = "icon-theme"
	gsKeyCursorTheme        = "cursor-theme"
	gsKeyFontStandard       = "font-standard"
	gsKeyFontMonospace      = "font-monospace"
	gsKeyFontSize           = "font-size"
	gsKeyBackgroundURIs     = "background-uris"
	gsKeyOpacity            = "opacity"
	gsKeyWallpaperSlideshow = "wallpaper-slideshow"

	wsPolicyLogin  = "login"
	wsPolicyWakeup = "wakeup"

	defaultIconTheme      = "deepin"
	defaultGtkTheme       = "deepin"
	defaultCursorTheme    = "deepin"
	defaultStandardFont   = "Noto Sans"
	defaultMonospaceFont  = "Noto Mono"
	defaultFontConfigFile = "/usr/share/deepin-default-settings/fontconfig.json"

	dbusServiceName = "com.deepin.daemon.Appearance"
	dbusPath        = "/com/deepin/daemon/Appearance"
	dbusInterface   = dbusServiceName
)

var wrConfigFile = filepath.Join(basedir.GetUserConfigDir(), "deepin/dde-daemon/appearance/wallpaper-slideshow.json")

// Manager shows current themes and fonts settings, emit 'Changed' signal if modified
// if themes list changed will emit 'Refreshed' signal
type Manager struct {
	service        *dbusutil.Service
	sessionSigLoop *dbusutil.SignalLoop
	sysSigLoop     *dbusutil.SignalLoop
	xConn          *x.Conn
	syncConfig     *dsync.Config

	GtkTheme      gsprop.String
	IconTheme     gsprop.String
	CursorTheme   gsprop.String
	Background    gsprop.String
	StandardFont  gsprop.String
	MonospaceFont gsprop.String
	Opacity       gsprop.Double `prop:"access:rw"`

	FontSize           gsprop.Double `prop:"access:rw"`
	WallpaperSlideShow gsprop.String `prop:"access:rw"`

	wsLoop      *WSLoop
	wsScheduler *WSScheduler

	userObj       *accounts.User
	imageBlur     *accounts.ImageBlur
	xSettings     *sessionmanager.XSettings
	login1Manager *login1.Manager

	setting        *gio.Settings
	xSettingsGs    *gio.Settings
	wrapBgSetting  *gio.Settings
	gnomeBgSetting *gio.Settings

	defaultFontConfig   DefaultFontConfig
	defaultFontConfigMu sync.Mutex

	watcher    *fsnotify.Watcher
	endWatcher chan struct{}

	desktopBgs []string
	greeterBg  string

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
		Delete                func() `in:"type,name"`
		GetScaleFactor        func() `out:"scale_factor"`
		List                  func() `in:"type" out:"list"`
		Set                   func() `in:"type,value"`
		SetScaleFactor        func() `in:"scale_factor"`
		Show                  func() `in:"type,names" out:"detail"`
		Thumbnail             func() `in:"type,name" out:"file"`
		SetScreenScaleFactors func() `in:"scaleFactors"`
		GetScreenScaleFactors func() `out:"scaleFactors"`
	}
}

// NewManager will create a 'Manager' object
func newManager(service *dbusutil.Service) *Manager {
	var m = new(Manager)
	m.service = service
	m.setting = gio.NewSettings(appearanceSchema)
	m.xSettingsGs = gio.NewSettings(xSettingsSchema)
	m.wrapBgSetting = gio.NewSettings(wrapBgSchema)

	m.GtkTheme.Bind(m.setting, gsKeyGtkTheme)
	m.IconTheme.Bind(m.setting, gsKeyIconTheme)
	m.CursorTheme.Bind(m.setting, gsKeyCursorTheme)
	m.StandardFont.Bind(m.setting, gsKeyFontStandard)
	m.MonospaceFont.Bind(m.setting, gsKeyFontMonospace)
	m.Background.Bind(m.wrapBgSetting, gsKeyBackground)
	m.FontSize.Bind(m.setting, gsKeyFontSize)
	m.Opacity.Bind(m.setting, gsKeyOpacity)
	m.WallpaperSlideShow.Bind(m.setting, gsKeyWallpaperSlideshow)

	m.wsLoop = newWSLoop()
	m.wsScheduler = newWSScheduler()

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
	m.desktopBgs = m.getBackgroundURIs()

	if m.userObj == nil {
		return
	}
	greeterBg, err := m.userObj.GreeterBackground().Get(0)
	if err == nil {
		m.greeterBg = greeterBg
	} else {
		logger.Warning(err)
	}
}

func (m *Manager) getBackgroundURIs() []string {
	return m.setting.GetStrv(gsKeyBackgroundURIs)
}

func (m *Manager) isBgInUse(file string) bool {
	if file == m.greeterBg {
		return true
	}

	for _, bg := range m.desktopBgs {
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
	m.sessionSigLoop.Stop()
	m.xSettings.RemoveHandler(proxy.RemoveAllHandlers)
	m.syncConfig.Destroy()

	m.sysSigLoop.Stop()
	m.login1Manager.RemoveHandler(proxy.RemoveAllHandlers)

	m.wsScheduler.stop()

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
		err := m.watcher.Close()
		if err != nil {
			logger.Warning(err)
		}
		m.watcher = nil
	}

	if m.xConn != nil {
		m.xConn.Close()
		m.xConn = nil
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

	err = common.ActivateSysDaemonService("com.deepin.daemon.Accounts")
	if err != nil {
		logger.Warning(err)
	}

	m.userObj, err = ddbus.NewUserByUid(systemConn, cur.Uid)
	if err != nil {
		logger.Warning("failed to new user object:", err)
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

func (m *Manager) init() error {
	background.SetCustomWallpaperDeleteCallback(func(file string) {
		logger.Debug("imageBlur delete", file)
		err := m.imageBlur.Delete(0, file)
		if err != nil {
			logger.Warning("imageBlur delete err:", err)
		}
	})

	sessionBus := m.service.Conn()
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	m.xConn, err = x.NewConn()
	if err != nil {
		return err
	}

	_, err = randr.QueryVersion(m.xConn, randr.MajorVersion,
		randr.MinorVersion).Reply(m.xConn)
	if err != nil {
		logger.Warning(err)
	}

	m.wm = wm.NewWm(sessionBus)
	m.imageBlur = accounts.NewImageBlur(systemBus)

	m.xSettings = sessionmanager.NewXSettings(sessionBus)
	theme_thumb.Init(m.getScaleFactor())

	m.sessionSigLoop = dbusutil.NewSignalLoop(sessionBus, 10)
	m.sessionSigLoop.Start()
	m.xSettings.InitSignalExt(m.sessionSigLoop, true)
	_, err = m.xSettings.ConnectSetScaleFactorStarted(handleSetScaleFactorStarted)
	if err != nil {
		logger.Warning(err)
	}
	_, err = m.xSettings.ConnectSetScaleFactorDone(handleSetScaleFactorDone)
	if err != nil {
		logger.Warning(err)
	}

	m.sysSigLoop = dbusutil.NewSignalLoop(systemBus, 10)
	m.sysSigLoop.Start()
	m.login1Manager = login1.NewManager(systemBus)
	m.login1Manager.InitSignalExt(m.sysSigLoop, true)
	m.initWallpaperSlideshow()

	err = m.loadDefaultFontConfig(defaultFontConfigFile)
	if err != nil {
		logger.Warning("load default font config failed:", err)
	}

	// set gtk theme
	gtkThemes := subthemes.ListGtkTheme()
	currentGtkTheme := m.GtkTheme.Get()
	if gtkThemes.Get(currentGtkTheme) == nil {
		m.GtkTheme.Set(defaultGtkTheme)
		currentGtkTheme = defaultGtkTheme
	}
	m.doSetGtkTheme(currentGtkTheme)

	// set icon theme
	iconThemes := subthemes.ListIconTheme()
	currentIconTheme := m.IconTheme.Get()
	if iconThemes.Get(currentIconTheme) == nil {
		m.IconTheme.Set(defaultIconTheme)
		currentIconTheme = defaultIconTheme
	}
	m.doSetIconTheme(currentIconTheme)

	// set cursor theme
	cursorThemes := subthemes.ListCursorTheme()
	currentCursorTheme := m.CursorTheme.Get()
	if cursorThemes.Get(currentCursorTheme) == nil {
		m.CursorTheme.Set(defaultCursorTheme)
		currentCursorTheme = defaultCursorTheme
	}
	m.doSetCursorTheme(currentCursorTheme)

	// Init theme list
	time.AfterFunc(time.Second*10, func() {
		if !dutils.IsFileExist(fonts.DeepinFontConfig) {
			m.resetFonts()
		} else {
			m.correctFontName()
		}

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

	m.initUserObj(systemBus)
	m.initCurrentBgs()
	m.syncConfig = dsync.NewConfig("appearance", &syncConfig{m: m}, m.sessionSigLoop, dbusPath, logger)
	return nil
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

	// set dde-kwin decoration theme
	var ddeKWinTheme string
	switch value {
	case "deepin":
		ddeKWinTheme = "light"
	case "deepin-dark":
		ddeKWinTheme = "dark"
	}
	if ddeKWinTheme != "" {
		err := m.wm.SetDecorationDeepinTheme(0, ddeKWinTheme)
		if err != nil {
			logger.Warning(err)
		}
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

func (m *Manager) doSetBackground(value string) (string, error) {
	logger.Debugf("call doSetBackground %q", value)
	if !background.IsBackgroundFile(value) {
		return "", errors.New("invalid background")
	}

	file, err := background.Prepare(value)
	if err != nil {
		return "", err
	}
	logger.Debug("prepare result:", file)
	uri := dutils.EncodeURI(file, dutils.SCHEME_FILE)
	m.wm.ChangeCurrentWorkspaceBackground(dbus.FlagNoAutoStart, uri)

	_, err = m.imageBlur.Get(0, file)
	if err != nil {
		logger.Warning("call imageBlur.Get err:", err)
	}
	return file, nil
}

func (m *Manager) doSetGreeterBackground(value string) error {
	value = dutils.EncodeURI(value, dutils.SCHEME_FILE)
	m.greeterBg = value
	if m.userObj == nil {
		return errors.New("user object is nil")
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

func (m *Manager) doSetMonospaceFont(value string) error {
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

func (m *Manager) saveWSConfig(t time.Time) error {
	cfg := &WSConfig{
		LastChange: t,
		Showed:     m.wsLoop.GetShowed(),
	}

	return cfg.save(wrConfigFile)
}

func (m *Manager) autoChangeBg(t time.Time) {
	logger.Debug("autoChangeBg", t)
	file := m.wsLoop.GetNext()
	if file == "" {
		logger.Warning("file is empty")
		return
	}

	_, err := m.doSetBackground(file)
	if err != nil {
		logger.Warning(err)
	}

	err = m.saveWSConfig(t)
	if err != nil {
		logger.Warning(err)
	}
}

func (m *Manager) initWallpaperSlideshow() {
	_, err := m.login1Manager.ConnectPrepareForSleep(func(before bool) {
		if !before {
			// after sleep
			if m.WallpaperSlideShow.Get() == wsPolicyWakeup {
				m.autoChangeBg(time.Now())
			}
		}
	})
	if err != nil {
		logger.Warning(err)
	}
	m.wsScheduler.fn = m.autoChangeBg

	policy := m.WallpaperSlideShow.Get()

	if isValidWSPolicy(policy) {
		m.loadWSConfig()
	}

	if policy == wsPolicyLogin {
		runDir, err := basedir.GetUserRuntimeDir(true)
		if err != nil {
			logger.Warning(err)
		} else {
			markFile := filepath.Join(runDir, "dde-daemon-wallpaper-slideshow-login")
			_, err = os.Stat(markFile)
			if os.IsNotExist(err) {
				m.autoChangeBg(time.Now())
				err = touchFile(markFile)
				if err != nil {
					logger.Warning(err)
				}
			} else if err != nil {
				logger.Warning(err)
			}
		}

	} else {
		nSec, err := strconv.ParseUint(policy, 10, 32)
		if err == nil {
			m.wsScheduler.updateInterval(time.Duration(nSec) * time.Second)
		}
	}
}

func (m *Manager) loadWSConfig() {
	cfg := loadWSConfigSafe(wrConfigFile)
	logger.Debug("loadWSConfig lastChange:", cfg.LastChange)

	m.wsScheduler.mu.Lock()
	m.wsScheduler.lastRun = cfg.LastChange
	m.wsScheduler.mu.Unlock()

	m.wsLoop.mu.Lock()
	for _, file := range cfg.Showed {
		m.wsLoop.showed[file] = struct{}{}
	}
	m.wsLoop.mu.Unlock()
}

func (m *Manager) updateWSPolicy(policy string) {
	if isValidWSPolicy(policy) {
		m.loadWSConfig()
	}
	nSec, err := strconv.ParseUint(policy, 10, 32)
	if err == nil {
		m.wsScheduler.updateInterval(time.Duration(nSec) * time.Second)
	} else {
		m.wsScheduler.stop()
	}
}

func touchFile(filename string) error {
	return ioutil.WriteFile(filename, nil, 0644)
}
