/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package appearance

import (
	"path"
	. "pkg.linuxdeepin.com/dde-daemon/appearance/factory"
	"pkg.linuxdeepin.com/dde-daemon/appearance/fonts"
	. "pkg.linuxdeepin.com/dde-daemon/appearance/utils"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/dbus/property"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"sort"
	"strings"
	"sync"
)

const (
	dbusSender = "com.deepin.daemon.ThemeManager"

	deepinGSKeyTheme   = "current-theme"
	deepinGSKeyPicture = "picture-uri"
	deepinGSKeySound   = "current-sound-theme"
	deepinGSKeyGreeter = "greeter-theme"

	defaultDThemeId = "Deepin"
)

type Manager struct {
	ThemeList        []string
	GtkThemeList     []string
	IconThemeList    []string
	CursorThemeList  []string
	SoundThemeList   []string
	BackgroundList   []string
	GreeterThemeList []string
	CurrentTheme     *property.GSettingsStringProperty `access:"readwrite"`
	GreeterTheme     *property.GSettingsStringProperty `access:"readwrite"`

	dtheme  FactoryInterface
	gtk     FactoryInterface
	icon    FactoryInterface
	cursor  FactoryInterface
	sound   FactoryInterface
	bg      FactoryInterface
	greeter FactoryInterface
	font    *fonts.FontManager

	themeObjMap map[string]*Theme

	settings      *gio.Settings
	wrapSetting   *gio.Settings
	gnomeSettings *gio.Settings

	lock    sync.Mutex
	wLocker sync.Mutex
}

func NewManager() *Manager {
	m := &Manager{}

	m.dtheme = NewFactory(ObjectTypeDeepinTheme, m.handleDThemeChanged)
	m.gtk = NewFactory(ObjectTypeGtk, m.setPropGtkThemeList)
	m.icon = NewFactory(ObjectTypeIcon, m.setPropIconThemeList)
	m.cursor = NewFactory(ObjectTypeCursor, m.setPropSoundThemeList)
	m.sound = NewFactory(ObjectTypeSound, m.setPropSoundThemeList)
	m.greeter = NewFactory(ObjectTypeGreeter, m.setPropGreeterThemeList)
	m.bg = NewFactory(ObjectTypeBackground, m.setPropBackgroundList)
	m.font = fonts.NewFontManager()

	m.setPropThemeList(getThemeObjectList(m.dtheme.GetNameStrList()))
	m.setPropGtkThemeList(m.gtk.GetNameStrList())
	m.setPropIconThemeList(m.icon.GetNameStrList())
	m.setPropCursorThemeList(m.cursor.GetNameStrList())
	m.setPropSoundThemeList(m.sound.GetNameStrList())
	m.setPropGreeterThemeList(m.greeter.GetNameStrList())
	m.setPropBackgroundList(m.bg.GetNameStrList())

	m.settings = NewGSettings("com.deepin.dde.personalization")
	m.CurrentTheme = property.NewGSettingsStringProperty(
		m, "CurrentTheme",
		m.settings, deepinGSKeyTheme)
	m.GreeterTheme = property.NewGSettingsStringProperty(
		m, "GreeterTheme",
		m.settings, deepinGSKeyGreeter)

	m.wrapSetting = NewGSettings("com.deepin.wrap.gnome.desktop.background")
	m.gnomeSettings = CheckAndNewGSettings("org.gnome.desktop.background")

	m.themeObjMap = make(map[string]*Theme)
	m.listenGSettings()
	m.rebuildThemes()

	return m
}

func (m *Manager) rebuildThemes() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.destroyThemes()

	var curThemeValid bool
	for _, name := range m.ThemeList {
		info, err := m.dtheme.GetInfoByName(path.Base(name))
		if err != nil {
			continue
		}
		t := NewTheme(info, m.applyTheme)
		err = dbus.InstallOnSession(t)
		if err != nil {
			logger.Warning("Install dbus failed:", info, err)
			continue
		}

		m.themeObjMap[info.BaseName] = t
		if info.BaseName == m.CurrentTheme.Get() {
			curThemeValid = true
			m.applyTheme(info.BaseName)
		}
	}

	if !curThemeValid {
		m.applyTheme(defaultDThemeId)
	}
}

func (m *Manager) destroyTheme(name string) {
	t, ok := m.themeObjMap[name]
	if !ok {
		return
	}

	t.Destroy()
	delete(m.themeObjMap, name)
}

func (m *Manager) destroyThemes() {
	for n, _ := range m.themeObjMap {
		m.destroyTheme(n)
	}

	m.themeObjMap = make(map[string]*Theme)
}

func (m *Manager) destroy() {
	m.destroyThemes()
	m.dtheme.Destroy()
	m.gtk.Destroy()
	m.icon.Destroy()
	m.cursor.Destroy()
	m.sound.Destroy()
	m.bg.Destroy()
	m.greeter.Destroy()
	m.font.Destroy()

	Unref(m.settings)
	Unref(m.gnomeSettings)

	// if dbus.InstallOnSession() not called, what will happen?
	// Now(2014.14.18) it is ok.
	dbus.UnInstallObject(m)
}

func sortNameByDeepin(list []string) []string {
	sort.Strings(list)
	deepinList := []string{}
	tmpList := []string{}

	for _, n := range list {
		t := strings.ToLower(n)
		if strings.Contains(t, "deepin") {
			deepinList = append(deepinList, n)
			continue
		}

		tmpList = append(tmpList, n)
	}

	if len(tmpList) > 0 {
		deepinList = append(deepinList, tmpList...)
	}

	return deepinList
}
