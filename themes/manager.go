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

package themes

import (
	"dlib/dbus"
	"dlib/dbus/property"
	"dlib/glib-2.0"
	"github.com/howeyc/fsnotify"
	"os"
	"path"
)

type Manager struct {
	ThemeList         []string
	GtkThemeList      []string
	IconThemeList     []string
	CursorThemeList   []string
	SoundThemeList    []string
	BackgroundList    []string
	CurrentTheme      *property.GSettingsStringProperty `access:"readwrite"`
	currentSound      *property.GSettingsStringProperty
	currentBackground *property.GSettingsStringProperty
	themeObjMap       map[string]*Theme

	watcher    *fsnotify.Watcher
	quitFlag   chan bool
	bgWatcher  *fsnotify.Watcher
	bgQuitFlag chan bool
}

var _manager *Manager

func GetManager() *Manager {
	if _manager == nil {
		_manager = newManager()
	}

	return _manager
}

func (obj *Manager) isThemeExit(gtk, icon, sound, cursor, bg string, fontSize int32) (string, bool) {
	for _, t := range obj.themeObjMap {
		if gtk == t.GtkTheme && icon == t.IconTheme &&
			sound == t.SoundTheme && cursor == t.CursorTheme &&
			bg == t.Background && fontSize == t.FontSize {
			return t.Name, true
		}
	}

	return "", false
}

func (obj *Manager) modifyTheme(name, gtk, icon, sound, cursor, bg string, fontSize int32) bool {
	filename := ""
	t, ok := obj.themeObjMap[name]
	if !ok {
		if str, ok := obj.mkdirTheme(name); !ok {
			return false
		} else {
			filename = path.Join(str, "theme.ini")
		}
	} else {
		filename = path.Join(t.filePath, "theme.ini")
	}

	kf := glib.NewKeyFile()
	defer kf.Free()
	_, err := kf.LoadFromFile(filename, glib.KeyFileFlagsKeepComments)
	kf.SetString(THEME_GROUP_COMPONENT, THEME_KEY_GTK, gtk)
	kf.SetString(THEME_GROUP_COMPONENT, THEME_KEY_ICON, icon)
	kf.SetString(THEME_GROUP_COMPONENT, THEME_KEY_SOUND, sound)
	kf.SetString(THEME_GROUP_COMPONENT, THEME_KEY_CURSOR, cursor)
	kf.SetString(THEME_GROUP_COMPONENT, THEME_KEY_BACKGROUND, bg)
	kf.SetInteger(THEME_GROUP_COMPONENT, THEME_KEY_FONT_SIZE, int(fontSize))

	_, contents, err := kf.ToData()
	if err != nil {
		Logger.Error("Convert Keyfile ToData Failed:", err)
		return false
	}

	if !writeStringToKeyFile(filename, contents) {
		return false
	}

	return true
}

func (obj *Manager) mkdirTheme(name string) (string, bool) {
	if len(name) < 1 {
		return "", false
	}

	homeDir, _ := objUtil.GetHomeDir()
	filePath := path.Join(homeDir, PERSON_LOCAL_THEME_PATH, name)

	Logger.Debugf("%s path: %s", name, filePath)
	os.MkdirAll(filePath, 0755)

	return filePath, true
}

func (obj *Manager) rebuildThemes() {
	obj.destroyAllTheme()

	for _, t := range obj.ThemeList {
		name := path.Base(t)
		tList := getDThemeList()
		for _, l := range tList {
			if name == l.Name {
				user := newTheme(l)
				obj.themeObjMap[user.Name] = user
				dbus.InstallOnSession(user)
				if obj.CurrentTheme.GetValue().(string) == user.Name {
					user.setAllThemes()
				}
				break
			}
		}
	}
}

func (obj *Manager) destroyTheme(theme string) {
	if obj, ok := obj.themeObjMap[theme]; ok {
		if obj.Type == THEME_TYPE_LOCAL {
			obj.quitFlag <- true
			obj.watcher = nil
		}
		dbus.UnInstallObject(obj)
	}
}

func (obj *Manager) destroyAllTheme() {
	if len(obj.themeObjMap) < 1 {
		return
	}

	for n, _ := range obj.themeObjMap {
		obj.destroyTheme(n)
	}

	obj.themeObjMap = make(map[string]*Theme)
}

func (obj *Manager) getDThemeStrList() []string {
	list := []string{}
	tList := getDThemeList()

	for _, l := range tList {
		list = append(list, THEME_PATH+l.Name)
	}

	return list
}

func (obj *Manager) getGtkStrList() []string {
	list := []string{}
	tList := getGtkThemeList()

	for _, l := range tList {
		list = append(list, l.Name)
	}

	return list
}

func (obj *Manager) getIconStrList() []string {
	list := []string{}
	tList := getIconThemeList()

	for _, l := range tList {
		list = append(list, l.Name)
	}

	return list
}

func (obj *Manager) getSoundStrList() []string {
	list := []string{}
	tList := getSoundThemeList()

	for _, l := range tList {
		list = append(list, l.Name)
	}

	return list
}

func (obj *Manager) getCursorStrList() []string {
	list := []string{}
	tList := getCursorThemeList()

	for _, l := range tList {
		list = append(list, l.Name)
	}

	return list
}

func (obj *Manager) getBgStrList() []string {
	list := []string{}
	tList := getBackgroundList()

	for _, l := range tList {
		list = append(list, l.Path)
	}

	return list
}

func newManager() *Manager {
	m := &Manager{}

	m.setPropThemeList(m.getDThemeStrList())
	m.setPropGtkThemeList(m.getGtkStrList())
	m.setPropIconThemeList(m.getIconStrList())
	m.setPropSoundThemeList(m.getSoundStrList())
	m.setPropCursorThemeList(m.getCursorStrList())
	m.setPropBackgroundList(m.getBgStrList())

	m.CurrentTheme = property.NewGSettingsStringProperty(
		m, "CurrentTheme",
		themeSettings, GS_KEY_CURRENT_THEME)
	m.currentSound = property.NewGSettingsStringProperty(
		m, "currentSound",
		themeSettings, GS_KEY_CURRENT_SOUND)
	m.currentBackground = property.NewGSettingsStringProperty(
		m, "currentBackground",
		themeSettings, GS_KEY_CURRENT_BG)

	m.themeObjMap = make(map[string]*Theme)
	m.rebuildThemes()

	var err error
	m.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		Logger.Errorf("New Watcher Failed: %v", err)
		panic(err)
	}

	m.quitFlag = make(chan bool)

	m.bgWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		Logger.Errorf("New Watcher Failed: %v", err)
		panic(err)
	}
	m.bgQuitFlag = make(chan bool)

	m.listenGSettings()
	m.startWatch()
	go m.handleEvent()
	m.startBgWatch()
	go m.handleBgEvent()

	return m
}
