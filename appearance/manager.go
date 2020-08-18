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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	dbus "github.com/godbus/dbus"
	accounts "github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.accounts"
	display "github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.display"
	imageeffect "github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.imageeffect"
	sessionmanager "github.com/linuxdeepin/go-dbus-factory/com.deepin.sessionmanager"
	wm "github.com/linuxdeepin/go-dbus-factory/com.deepin.wm"
	geoclue "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.geoclue2"
	login1 "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.login1"
	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/randr"
	"pkg.deepin.io/dde/api/theme_thumb"
	"pkg.deepin.io/dde/daemon/appearance/background"
	"pkg.deepin.io/dde/daemon/appearance/fonts"
	"pkg.deepin.io/dde/daemon/appearance/subthemes"
	"pkg.deepin.io/dde/daemon/common/dsync"
	ddbus "pkg.deepin.io/dde/daemon/dbus"
	"pkg.deepin.io/dde/daemon/session/common"
	gio "pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/gsprop"
	"pkg.deepin.io/lib/dbusutil/proxy"
	"pkg.deepin.io/lib/log"
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
	gsKeyGtkTheme           = "gtk-theme"
	gsKeyIconTheme          = "icon-theme"
	gsKeyCursorTheme        = "cursor-theme"
	gsKeyFontStandard       = "font-standard"
	gsKeyFontMonospace      = "font-monospace"
	gsKeyFontSize           = "font-size"
	gsKeyBackgroundURIs     = "background-uris"
	gsKeyOpacity            = "opacity"
	gsKeyWallpaperSlideshow = "wallpaper-slideshow"
	gsKeyWallpaperURIs      = "wallpaper-uris"
	gsKeyQtActiveColor      = "qt-active-color"

	propQtActiveColor = "QtActiveColor"

	wsPolicyLogin  = "login"
	wsPolicyWakeup = "wakeup"

	defaultIconTheme      = "bloom"
	defaultGtkTheme       = "deepin"
	autoGtkTheme          = "deepin-auto"
	defaultCursorTheme    = "bloom"
	defaultStandardFont   = "Noto Sans"
	defaultMonospaceFont  = "Noto Mono"
	defaultFontConfigFile = "/usr/share/deepin-default-settings/fontconfig.json"

	dbusServiceName = "com.deepin.daemon.Appearance"
	dbusPath        = "/com/deepin/daemon/Appearance"
	dbusInterface   = dbusServiceName
)

var wsConfigFile = filepath.Join(basedir.GetUserConfigDir(), "deepin/dde-daemon/appearance/wallpaper-slideshow.json")

// Manager shows current themes and fonts settings, emit 'Changed' signal if modified
// if themes list changed will emit 'Refreshed' signal
type Manager struct {
	service        *dbusutil.Service
	sessionSigLoop *dbusutil.SignalLoop
	sysSigLoop     *dbusutil.SignalLoop
	xConn          *x.Conn
	syncConfig     *dsync.Config
	bgSyncConfig   *dsync.Config

	GtkTheme           gsprop.String
	IconTheme          gsprop.String
	CursorTheme        gsprop.String
	Background         gsprop.String
	StandardFont       gsprop.String
	MonospaceFont      gsprop.String
	Opacity            gsprop.Double `prop:"access:rw"`
	FontSize           gsprop.Double `prop:"access:rw"`
	WallpaperSlideShow gsprop.String `prop:"access:rw"`
	WallpaperURIs      gsprop.String
	QtActiveColor      string `prop:"access:rw"`

	wsLoopMap      map[string]*WSLoop
	wsSchedulerMap map[string]*WSScheduler
	monitorMap     map[string]string

	userObj             *accounts.User
	imageBlur           *accounts.ImageBlur
	imageEffect         *imageeffect.ImageEffect
	xSettings           *sessionmanager.XSettings
	login1Manager       *login1.Manager
	geoclueClient       *geoclue.Client
	themeAutoTimer      *time.Timer
	display             *display.Display
	latitude            float64
	longitude           float64
	locationValid       bool
	detectSysClockTimer *time.Timer
	ts                  int64

	setting        *gio.Settings
	xSettingsGs    *gio.Settings
	wrapBgSetting  *gio.Settings
	gnomeBgSetting *gio.Settings

	defaultFontConfig   DefaultFontConfig
	defaultFontConfigMu sync.Mutex

	watcher    *fsnotify.Watcher
	endWatcher chan struct{}

	desktopBgs      []string
	greeterBg       string
	curMonitorSpace string
	wm              *wm.Wm

	//nolint
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

	//nolint
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
		SetMonitorBackground  func() `in:"monitorName,imageFile"`
		SetWallpaperSlideShow func() `in:"monitorName,wallpaperSlideShow"`
		GetWallpaperSlideShow func() `in:"monitorName" out:"slideShow"`
	}
}

type mapMonitorWorkspaceWSPolicy map[string]string
type mapMonitorWorkspaceWSConfig map[string]WSConfig
type mapMonitorWorkspaceWallpaperURIs map[string]string

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
	m.WallpaperURIs.Bind(m.setting, gsKeyWallpaperURIs)
	var err error
	m.QtActiveColor, err = m.getQtActiveColor()
	if err != nil {
		logger.Warning(err)
	}

	m.wsLoopMap = make(map[string]*WSLoop)
	m.wsSchedulerMap = make(map[string]*WSScheduler)

	m.gnomeBgSetting, _ = dutils.CheckAndNewGSettings(gnomeBgSchema)

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
	m.bgSyncConfig.Destroy()

	m.sysSigLoop.Stop()
	m.login1Manager.RemoveHandler(proxy.RemoveAllHandlers)
	for iSche := range m.wsSchedulerMap {
		m.wsSchedulerMap[iSche].stop()
	}

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

		logger.Debug("imageEffect delete", file)
		err = m.imageEffect.Delete(0, "all", file)
		if err != nil {
			logger.Warning("imageEffect delete err:", err)
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

	m.sessionSigLoop = dbusutil.NewSignalLoop(sessionBus, 10)
	m.sessionSigLoop.Start()

	m.wm = wm.NewWm(sessionBus)
	m.wm.InitSignalExt(m.sessionSigLoop, true)
	_, err = m.wm.ConnectWorkspaceCountChanged(m.handleWmWorkspaceCountChanged)
	if err != nil {
		logger.Warning(err)
	}
	m.imageBlur = accounts.NewImageBlur(systemBus)
	m.imageEffect = imageeffect.NewImageEffect(systemBus)

	m.xSettings = sessionmanager.NewXSettings(sessionBus)
	theme_thumb.Init(m.getScaleFactor())

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

	if currentGtkTheme == autoGtkTheme {
		m.updateThemeAuto(true)
	} else {
		if gtkThemes.Get(currentGtkTheme) == nil {
			m.GtkTheme.Set(defaultGtkTheme)
			err = m.doSetGtkTheme(defaultGtkTheme)
			if err != nil {
				logger.Warning("failed to set gtk theme:", err)
			}
		}
	}

	// set icon theme
	iconThemes := subthemes.ListIconTheme()
	currentIconTheme := m.IconTheme.Get()
	if iconThemes.Get(currentIconTheme) == nil {
		m.IconTheme.Set(defaultIconTheme)
		currentIconTheme = defaultIconTheme
	}
	err = m.doSetIconTheme(currentIconTheme)
	if err != nil {
		logger.Warning("failed to set icon theme:", err)
	}

	// set cursor theme
	cursorThemes := subthemes.ListCursorTheme()
	currentCursorTheme := m.CursorTheme.Get()
	if cursorThemes.Get(currentCursorTheme) == nil {
		m.CursorTheme.Set(defaultCursorTheme)
		currentCursorTheme = defaultCursorTheme
	}
	err = m.doSetCursorTheme(currentCursorTheme)
	if err != nil {
		logger.Warning("failed to set cursor theme:", err)
	}

	// Init theme list
	time.AfterFunc(time.Second*10, func() {
		if !dutils.IsFileExist(fonts.DeepinFontConfig) {
			m.resetFonts()
		} else {
			m.correctFontName()
		}

		fonts.GetFamilyTable()

		err = setDQtTheme(dQtFile, dQtSectionTheme,
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
		if err != nil {
			logger.Warning("failed to set deepin qt theme:", err)
		}
		err = saveDQtTheme(dQtFile)
		if err != nil {
			logger.Warning("Failed to save deepin qt theme:", err)
			return
		}
	})

	m.initUserObj(systemBus)
	m.initCurrentBgs()
	m.display = display.NewDisplay(sessionBus)
	m.display.InitSignalExt(m.sessionSigLoop, true)
	err = m.display.Monitors().ConnectChanged(func(hasValue bool, value []dbus.ObjectPath) {
		m.handleDisplayChanged(hasValue)
	})
	if err != nil {
		logger.Warning("failed to connect Monitors changed:", err)
	}
	err = m.display.Primary().ConnectChanged(func(hasValue bool, value string) {
		m.handleDisplayChanged(hasValue)
	})
	if err != nil {
		logger.Warning("failed to connect Primary changed:", err)
	}
	m.updateMonitorMap()
	m.syncConfig = dsync.NewConfig("appearance", &syncConfig{m: m}, m.sessionSigLoop, dbusPath, logger)
	m.bgSyncConfig = dsync.NewConfig("background", &backgroundSyncConfig{m: m}, m.sessionSigLoop,
		backgroundDBusPath, logger)
	return nil
}

func (m *Manager) handleDisplayChanged(hasValue bool) {
	if !hasValue {
		return
	}
	m.updateMonitorMap()
	err := m.doUpdateWallpaperURIs()
	if err != nil {
		logger.Warning("failed to update WallpaperURIs:", err)
	}
}

func (m *Manager) handleWmWorkspaceCountChanged(count int32) {
	logger.Debug("wm workspace count changed", count)
	bgs := m.setting.GetStrv(gsKeyBackgroundURIs)
	if len(bgs) < int(count) {
		allBgs := background.ListBackground()

		numAdded := int(count) - len(bgs)
		for i := 0; i < numAdded; i++ {
			idx := rand.Intn(len(allBgs))
			// Id is file url
			bgs = append(bgs, allBgs[idx].Id)
		}
		m.setting.SetStrv(gsKeyBackgroundURIs, bgs)
	} else if len(bgs) > int(count) {
		bgs = bgs[:int(count)]
		m.setting.SetStrv(gsKeyBackgroundURIs, bgs)
	}
	err := m.doUpdateWallpaperURIs()
	if err != nil {
		logger.Warning("failed to update WallpaperURIs:", err)
	}
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
	if value == autoGtkTheme {
		return nil
	}
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

func (m *Manager) doSetMonitorBackground(monitorName string, imageFile string) (string, error) {
	logger.Debugf("call doSetMonitorBackground monitor:%q file:%q", monitorName, imageFile)
	if !background.IsBackgroundFile(imageFile) {
		return "", errors.New("invalid background")
	}
	file, err := background.Prepare(imageFile)
	if err != nil {
		logger.Warning("failed to prepare:", err)
		return "", err
	}
	logger.Debug("prepare result:", file)
	uri := dutils.EncodeURI(file, dutils.SCHEME_FILE)
	err = m.wm.SetCurrentWorkspaceBackgroundForMonitor(0, uri, monitorName)
	if err != nil {
		return "", err
	}
	err = m.doUpdateWallpaperURIs()
	if err != nil {
		logger.Warning("failed to update WallpaperURIs:", err)
	}
	_, err = m.imageBlur.Get(0, file)
	if err != nil {
		logger.Warning("call imageBlur.Get err:", err)
	}
	go func() {
		outputFile, err := m.imageEffect.Get(0, "", file)
		if err != nil {
			logger.Warning("imageEffect Get err:", err)
		} else {
			logger.Warning("imageEffect Get outputFile:", outputFile)
		}
	}()
	return file, nil
}

func (m *Manager) updateMonitorMap() {
	monitorList, _ := m.display.ListOutputNames(0)
	primary, _ := m.display.Primary().Get(0)
	index := 0
	m.monitorMap = make(map[string]string)
	for _, item := range monitorList {
		if item == primary {
			m.monitorMap[item] = "Primary"
		} else {
			m.monitorMap[item] = "Subsidiary" + strconv.Itoa(index)
			index++
		}
	}
}

func (m *Manager) reverseMonitorMap() map[string]string {
	reverseMap := make(map[string]string)
	for k, v := range m.monitorMap {
		reverseMap[v] = k
	}
	return reverseMap
}

func (m *Manager) doUpdateWallpaperURIs() error {
	mapWallpaperURIs := make(mapMonitorWorkspaceWallpaperURIs)
	workspaceCount, _ := m.wm.WorkspaceCount(0)
	monitorList, _ := m.display.ListOutputNames(0)
	for _, monitor := range monitorList {
		for idx := int32(1); idx <= workspaceCount; idx++ {
			wallpaperURI, err := m.wm.GetWorkspaceBackgroundForMonitor(0, idx, monitor)
			if err != nil {
				logger.Warning("get wallpaperURI failed:", err)
				continue
			}
			key := m.monitorMap[monitor] + "&&" + strconv.Itoa(int(idx))
			mapWallpaperURIs[key] = wallpaperURI
		}
	}

	jsonString, err := doMarshalMonitorWorkspaceWallpaperURIs(mapWallpaperURIs)
	if err != nil {
		return err
	}
	m.WallpaperURIs.Set(jsonString)
	return nil
}

func doUnmarshalMonitorWorkspaceWallpaperURIs(jsonString string) (mapMonitorWorkspaceWallpaperURIs, error) {
	var cfg mapMonitorWorkspaceWallpaperURIs
	var byteMonitorWorkspaceWallpaperURIs = []byte(jsonString)
	err := json.Unmarshal(byteMonitorWorkspaceWallpaperURIs, &cfg)
	return cfg, err
}

func doMarshalMonitorWorkspaceWallpaperURIs(cfg mapMonitorWorkspaceWallpaperURIs) (string, error) {
	data, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func doUnmarshalWallpaperSlideshow(jsonString string) (mapMonitorWorkspaceWSPolicy, error) {
	var cfg mapMonitorWorkspaceWSPolicy
	var byteWallpaperSlideShow []byte = []byte(jsonString)
	err := json.Unmarshal(byteWallpaperSlideShow, &cfg)
	return cfg, err
}

func doMarshalWallpaperSlideshow(cfg mapMonitorWorkspaceWSPolicy) (string, error) {
	data, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func (m *Manager) doSetWallpaperSlideShow(monitorName string, wallpaperSlideShow string) error {
	idx, err := m.wm.GetCurrentWorkspace(0)
	if err != nil {
		logger.Warning("Get Current Workspace failure:", err)
		return err
	}
	cfg, err := doUnmarshalWallpaperSlideshow(m.WallpaperSlideShow.Get())
	if err != nil {
		logger.Warning("doUnmarshalWallpaperSlideshow Failed:", err)
	}
	if cfg == nil {
		cfg = make(mapMonitorWorkspaceWSPolicy)
	}

	key := monitorName + "&&" + strconv.Itoa(int(idx))
	cfg[key] = wallpaperSlideShow
	strAllWallpaperSlideShow, err := doMarshalWallpaperSlideshow(cfg)
	if err != nil {
		logger.Warning("Marshal Wallpaper Slideshow failure:", err)
	}
	m.WallpaperSlideShow.Set(strAllWallpaperSlideShow)
	m.curMonitorSpace = key
	return nil
}

func (m *Manager) doGetWallpaperSlideShow(monitorName string) (string, error) {
	idx, err := m.wm.GetCurrentWorkspace(0)
	if err != nil {
		logger.Warning("Get Current Workspace failure:", err)
		return "", err
	}
	cfg, err := doUnmarshalWallpaperSlideshow(m.WallpaperSlideShow.Get())
	if err != nil {
		return "", nil
	}
	key := monitorName + "&&" + strconv.Itoa(int(idx))
	wallpaperSlideShow := cfg[key]
	return wallpaperSlideShow, nil
}

func (m *Manager) doSetBackground(value string) (string, error) {
	logger.Debugf("call doSetBackground %q", value)
	if !background.IsBackgroundFile(value) {
		return "", errors.New("invalid background")
	}

	file, err := background.Prepare(value)
	if err != nil {
		logger.Warning("failed to prepare:", err)
		return "", err
	}
	logger.Debug("prepare result:", file)
	uri := dutils.EncodeURI(file, dutils.SCHEME_FILE)
	err = m.wm.ChangeCurrentWorkspaceBackground(0, uri)
	if err != nil {
		return "", err
	}

	_, err = m.imageBlur.Get(0, file)
	if err != nil {
		logger.Warning("call imageBlur.Get err:", err)
	}
	go func() {
		outputFile, err := m.imageEffect.Get(0, "", file)
		if err != nil {
			logger.Warning("imageEffect Get err:", err)
		} else {
			logger.Warning("imageEffect Get outputFile:", outputFile)
		}
	}()

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

	monoFont := m.MonospaceFont.Get()
	if !fonts.IsFontFamily(monoFont) {
		monoList := fonts.GetFamilyTable().ListMonospace()
		if len(monoList) == 0 {
			return fmt.Errorf("no valid mono font")
		}
		monoFont = monoList[0]
	}

	err := fonts.SetFamily(value, monoFont, m.FontSize.Get())
	if err != nil {
		return err
	}

	err = m.xSettings.SetString(0, "Qt/FontName", value)
	if err != nil {
		return err
	}

	return m.writeDQtTheme(dQtKeyFont, value)
}

func (m *Manager) doSetMonospaceFont(value string) error {
	if !fonts.IsFontFamily(value) {
		return fmt.Errorf("invalid font family '%v'", value)
	}

	standardFont := m.StandardFont.Get()
	if !fonts.IsFontFamily(standardFont) {
		standardList := fonts.GetFamilyTable().ListStandard()
		if len(standardList) == 0 {
			return fmt.Errorf("no valid standard font")
		}
		standardFont = standardList[0]
	}

	err := fonts.SetFamily(standardFont, value, m.FontSize.Get())
	if err != nil {
		return err
	}

	err = m.xSettings.SetString(0, "Qt/MonoFontName", value)
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

	err = m.xSettings.SetString(0, "Qt/FontPointSize",
		strconv.FormatFloat(size, 'f', -1, 64))
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
	err := setDQtTheme(dQtFile, dQtSectionTheme,
		[]string{key}, []string{value})
	if err != nil {
		logger.Warning("failed to set deepin qt theme:", err)
	}
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

func (m *Manager) saveWSConfig(monitorSpace string, t time.Time) error {
	cfg, _ := loadWSConfig(wsConfigFile)
	var tempCfg WSConfig
	tempCfg.LastChange = t
	if m.wsLoopMap[monitorSpace] != nil {
		tempCfg.Showed = m.wsLoopMap[monitorSpace].GetShowed()
	}
	if cfg == nil {
		cfg = make(mapMonitorWorkspaceWSConfig)
	}
	cfg[monitorSpace] = tempCfg
	return cfg.save(wsConfigFile)
}

func (m *Manager) autoChangeBg(monitorSpace string, t time.Time) {
	logger.Debug("autoChangeBg", monitorSpace, t)
	if m.wsLoopMap[monitorSpace] == nil {
		return
	}
	file := m.wsLoopMap[monitorSpace].GetNext()
	if file == "" {
		logger.Warning("file is empty")
		return
	}
	idx, err := m.wm.GetCurrentWorkspace(0)
	if err != nil {
		logger.Warning(err)
	}
	strIdx := strconv.Itoa(int(idx))
	splitter := strings.Index(monitorSpace, "&&")
	if splitter == -1 {
		logger.Warning("monitorSpace format error")
		return
	}
	if strIdx == monitorSpace[splitter+len("&&"):] {
		_, err := m.doSetMonitorBackground(monitorSpace[:splitter], file)
		if err != nil {
			logger.Warning("failed to set background:", err)
		}
	}
	err = m.saveWSConfig(monitorSpace, t)
	if err != nil {
		logger.Warning(err)
	}
}

func (m *Manager) initWallpaperSlideshow() {
	m.loadWSConfig()
	cfg, err := doUnmarshalWallpaperSlideshow(m.WallpaperSlideShow.Get())
	if err == nil {
		for monitorSpace, policy := range cfg {
			_, ok := m.wsSchedulerMap[monitorSpace]
			if !ok {
				m.wsSchedulerMap[monitorSpace] = newWSScheduler(m.autoChangeBg)
			}
			_, ok = m.wsLoopMap[monitorSpace]
			if !ok {
				m.wsLoopMap[monitorSpace] = newWSLoop()
			}
			if isValidWSPolicy(policy) {
				if policy == wsPolicyLogin {
					err := m.changeBgAfterLogin(monitorSpace)
					if err != nil {
						logger.Warning("failed to change background after login:", err)
					}
				} else {
					nSec, err := strconv.ParseUint(policy, 10, 32)
					if err == nil && m.wsSchedulerMap[monitorSpace] != nil {
						m.wsSchedulerMap[monitorSpace].updateInterval(monitorSpace, time.Duration(nSec)*time.Second)
					}
				}
			}
		}
	} else {
		logger.Debug("doUnmarshalWallpaperSlideshow err is ", err)
	}
}

func (m *Manager) changeBgAfterLogin(monitorSpace string) error {
	runDir, err := basedir.GetUserRuntimeDir(true)
	if err != nil {
		return err
	}

	currentSessionId, err := getSessionId("/proc/self/sessionid")
	if err != nil {
		return err
	}

	var needChangeBg bool
	markFile := filepath.Join(runDir, "dde-daemon-wallpaper-slideshow-login"+monitorSpace)
	sessionId, err := getSessionId(markFile)
	if err == nil {
		if sessionId != currentSessionId {
			needChangeBg = true
		}
	} else if os.IsNotExist(err) {
		needChangeBg = true
	} else if err != nil {
		return err
	}

	if needChangeBg {
		m.autoChangeBg(monitorSpace, time.Now())
		err = ioutil.WriteFile(markFile, []byte(currentSessionId), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func getSessionId(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(content)), nil
}

func (m *Manager) loadWSConfig() {
	cfg := loadWSConfigSafe(wsConfigFile)
	for monitorSpace := range cfg {
		_, ok := m.wsSchedulerMap[monitorSpace]
		if !ok {
			m.wsSchedulerMap[monitorSpace] = newWSScheduler(m.autoChangeBg)
		}
		m.wsSchedulerMap[monitorSpace].mu.Lock()
		m.wsSchedulerMap[monitorSpace].lastSetBg = cfg[monitorSpace].LastChange
		m.wsSchedulerMap[monitorSpace].mu.Unlock()

		_, ok = m.wsLoopMap[monitorSpace]
		if !ok {
			m.wsLoopMap[monitorSpace] = newWSLoop()
		}
		m.wsLoopMap[monitorSpace].mu.Lock()
		for _, file := range cfg[monitorSpace].Showed {
			m.wsLoopMap[monitorSpace].showed[file] = struct{}{}
		}
		m.wsLoopMap[monitorSpace].mu.Unlock()
	}
}

func (m *Manager) updateWSPolicy(policy string) {
	cfg, err := doUnmarshalWallpaperSlideshow(policy)
	m.loadWSConfig()
	if err == nil {
		for monitorSpace, policy := range cfg {
			_, ok := m.wsSchedulerMap[monitorSpace]
			if !ok {
				m.wsSchedulerMap[monitorSpace] = newWSScheduler(m.autoChangeBg)
			}
			_, ok = m.wsLoopMap[monitorSpace]
			if !ok {
				m.wsLoopMap[monitorSpace] = newWSLoop()
			}
			if m.curMonitorSpace == monitorSpace && isValidWSPolicy(policy) {
				nSec, err := strconv.ParseUint(policy, 10, 32)
				if err == nil {
					m.wsSchedulerMap[monitorSpace].lastSetBg = time.Now()
					m.wsSchedulerMap[monitorSpace].updateInterval(monitorSpace, time.Duration(nSec)*time.Second)
					err = m.saveWSConfig(monitorSpace, time.Now())
					if err != nil {
						logger.Warning(err)
					}
				} else {
					m.wsSchedulerMap[monitorSpace].stop()
				}
			}
		}
	}
}

func (m *Manager) enableDetectSysClock(enabled bool) {
	logger.Debug("enableDetectSysClock:", enabled)
	nSec := 60 // 1 min
	if logger.GetLogLevel() == log.LevelDebug {
		// debug mode: 10 s
		nSec = 10
	}
	interval := time.Duration(nSec) * time.Second
	if enabled {
		m.ts = time.Now().Unix()
		if m.detectSysClockTimer == nil {
			m.detectSysClockTimer = time.AfterFunc(interval, func() {
				nowTs := time.Now().Unix()
				d := nowTs - m.ts - int64(nSec)
				if !(-2 < d && d < 2) {
					m.handleSysClockChanged()
				}

				m.ts = time.Now().Unix()
				m.detectSysClockTimer.Reset(interval)
			})
		} else {
			m.detectSysClockTimer.Reset(interval)
		}
	} else {
		// disable
		if m.detectSysClockTimer != nil {
			m.detectSysClockTimer.Stop()
		}
	}
}

func (m *Manager) handleSysClockChanged() {
	logger.Debug("system clock changed")
	if m.locationValid {
		m.autoSetTheme(m.latitude, m.longitude)
		m.resetThemeAutoTimer()
	}
}

func (m *Manager) updateThemeAuto(enabled bool) {
	m.enableDetectSysClock(enabled)
	logger.Debug("updateThemeAuto:", enabled)
	if enabled {
		var err error
		if m.themeAutoTimer == nil {
			m.themeAutoTimer = time.AfterFunc(0, func() {
				if m.locationValid {
					m.autoSetTheme(m.latitude, m.longitude)

					time.AfterFunc(5*time.Second, func() {
						m.resetThemeAutoTimer()
					})
				}
			})
		} else {
			m.themeAutoTimer.Reset(0)
		}

		if m.geoclueClient == nil {
			m.geoclueClient, err = getGeoclueClient()
			if err != nil {
				logger.Warning("failed to get geoclue client:", err)
				return
			}

			m.geoclueClient.InitSignalExt(m.sysSigLoop, true)
			_, err = m.geoclueClient.ConnectLocationUpdated(
				func(old dbus.ObjectPath, newLoc dbus.ObjectPath) {
					sysBus, err := dbus.SystemBus()
					if err != nil {
						logger.Warning(err)
						return
					}
					loc, err := geoclue.NewLocation(sysBus, newLoc)
					if err != nil {
						logger.Warning(err)
						return
					}

					latitude, err := loc.Latitude().Get(0)
					if err != nil {
						logger.Warning("failed to get latitude:", err)
						return
					}

					longitude, err := loc.Longitude().Get(0)
					if err != nil {
						logger.Warning("failed to get longitude:", err)
						return
					}
					m.updateLocation(latitude, longitude)
				})
			if err != nil {
				logger.Warning(err)
			}
		}

		locPath, err := m.geoclueClient.Location().Get(0)
		if err == nil {
			if locPath != "/" {
				latitude, longitude, err := getLocation(locPath)
				if err == nil {
					m.updateLocation(latitude, longitude)
				} else {
					logger.Warning("failed to get location:", err)
				}
			} else {
				logger.Debug("wait location updated signal")
			}
		} else {
			logger.Warning("failed to get geoclue client location path:", err)
		}

		err = m.geoclueClient.Start(0)
		if err != nil {
			logger.Warning("failed to start geoclue client:", err)
			return
		}

	} else {
		// disable geoclue client
		if m.geoclueClient != nil {
			err := m.geoclueClient.Stop(0)
			if err != nil {
				logger.Warning("failed to stop geoclue client:", err)
			}

			m.geoclueClient.RemoveAllHandlers()
			m.geoclueClient = nil
		}
		m.latitude = 0
		m.longitude = 0
		m.locationValid = false
		if m.themeAutoTimer != nil {
			m.themeAutoTimer.Stop()
		}
	}
}

func (m *Manager) updateLocation(latitude, longitude float64) {
	m.latitude = latitude
	m.longitude = longitude
	m.locationValid = true
	logger.Debugf("update location, latitude: %v, longitude: %v",
		latitude, longitude)
	m.autoSetTheme(latitude, longitude)
	m.resetThemeAutoTimer()
}

func (m *Manager) resetThemeAutoTimer() {
	if m.themeAutoTimer == nil {
		logger.Debug("themeAutoTimer is nil")
		return
	}
	if !m.locationValid {
		logger.Debug("location is invalid")
		return
	}

	now := time.Now()
	changeTime, err := getThemeAutoChangeTime(now, m.latitude, m.longitude)
	if err != nil {
		logger.Warning("failed to get theme auto change time:", err)
		return
	}

	interval := changeTime.Sub(now)
	logger.Debug("change theme after:", interval)
	m.themeAutoTimer.Reset(interval)
}

func (m *Manager) autoSetTheme(latitude, longitude float64) {
	now := time.Now()
	if m.GtkTheme.Get() != autoGtkTheme {
		return
	}

	sunriseT, sunsetT, err := getSunriseSunset(now, latitude, longitude)
	if err != nil {
		logger.Warning(err)
		return
	}
	logger.Debugf("now: %v, sunrise: %v, sunset: %v",
		now, sunriseT, sunsetT)
	themeName := getThemeAutoName(isDaytime(now, sunriseT, sunsetT))
	logger.Debug("auto theme name:", themeName)

	currentTheme := m.GtkTheme.Get()
	if currentTheme != themeName {
		err = m.doSetGtkTheme(themeName)
		if err != nil {
			logger.Warning(err)
		}
	}
}

func (m *Manager) getQtActiveColor() (string, error) {
	str := m.xSettingsGs.GetString(gsKeyQtActiveColor)
	return xsColorToHexColor(str)
}

func xsColorToHexColor(str string) (string, error) {
	fields := strings.Split(str, ",")
	if len(fields) != 4 {
		return "", errors.New("length of fields is not 4")
	}

	var array [4]uint16
	for idx, field := range fields {
		v, err := strconv.ParseUint(field, 10, 16)
		if err != nil {
			return "", err
		}
		array[idx] = uint16(v)
	}

	var byteArr [4]byte
	for idx, value := range array {
		byteArr[idx] = byte((float64(value) / float64(math.MaxUint16)) * float64(math.MaxUint8))
	}
	return byteArrayToHexColor(byteArr), nil
}

func byteArrayToHexColor(p [4]byte) string {
	// p : [R G B A]
	if p[3] == 255 {
		return fmt.Sprintf("#%02X%02X%02X", p[0], p[1], p[2])
	}
	return fmt.Sprintf("#%02X%02X%02X%02X", p[0], p[1], p[2], p[3])
}

var hexColorReg = regexp.MustCompile(`^#([0-9A-F]{6}|[0-9A-F]{8})$`)

func parseHexColor(hexColor string) (array [4]byte, err error) {
	hexColor = strings.ToUpper(hexColor)
	match := hexColorReg.FindStringSubmatch(hexColor)
	if match == nil {
		err = errors.New("invalid hex color format")
		return
	}
	hexNums := string(match[1])
	count := 4
	if len(hexNums) == 6 {
		count = 3
		array[3] = 255
	}

	for i := 0; i < count; i++ {
		array[i], err = parseHexNum(hexNums[i*2 : i*2+2])
		if err != nil {
			return
		}
	}
	return
}

func parseHexNum(str string) (byte, error) {
	v, err := strconv.ParseUint(str, 16, 8)
	return byte(v), err
}

func (m *Manager) setQtActiveColor(hexColor string) error {
	xsColor, err := hexColorToXsColor(hexColor)
	if err != nil {
		return err
	}

	ok := m.xSettingsGs.SetString(gsKeyQtActiveColor, xsColor)
	if !ok {
		return errors.New("failed to save")
	}
	return nil
}

func hexColorToXsColor(hexColor string) (string, error) {
	byteArr, err := parseHexColor(hexColor)
	if err != nil {
		return "", err
	}
	var array [4]uint16
	for idx, value := range byteArr {
		array[idx] = uint16((float64(value) / float64(math.MaxUint8)) * float64(math.MaxUint16))
	}
	return fmt.Sprintf("%d,%d,%d,%d", array[0], array[1],
		array[2], array[3]), nil
}
