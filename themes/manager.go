/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
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

package main

import (
        "dlib/glib-2.0"
        "os"
)

type Manager struct {
        ThemeList    []string
        CurrentTheme string
        GtkThemeList []string
        //GtkBasePath     string // TODO
        IconThemeList []string
        //IconBasePath    string
        CursorThemeList []string
        //CursorBasePath  string
        FontThemeList  []string
        BackgroundList []string
        SoundThemeList []string
        pathNameMap    map[string]PathInfo
}

func (op *Manager) GetPathViaName(name string) (string, bool) {
        if path, ok := themeNamePathMap[name]; ok {
                return path, true
        }

        return "", false
}

func (op *Manager) SetCurrentTheme(name string) bool {
        if _, ok := themeNamePathMap[name]; !ok {
                return false
        }

        if name != op.CurrentTheme {
                if obj := op.getThemeObject(name); obj != nil {
                        obj.setThemeViaXSettings()
                }
                personSettings.SetString(GKEY_CURRENT_THEME,
                        name)
                op.CurrentTheme = name
        }

        return true
}

func (op *Manager) SetGtkTheme(name string) (string, bool) {
        if len(name) <= 0 {
                return op.CurrentTheme, false
        }

        defer setGtkThemeViaXSettings(name)
        if obj := op.getThemeObject(op.CurrentTheme); obj != nil {
                v := op.setTheme(name, obj.IconTheme, obj.CursorTheme,
                        obj.FontName, obj.BackgroundFile, obj.SoundTheme)
                op.updateGSettingsKey(GKEY_CURRENT_THEME, v)
                return v, true
        }

        return op.CurrentTheme, false
}

func (op *Manager) SetIconTheme(name string) (string, bool) {
        if len(name) <= 0 {
                return op.CurrentTheme, false
        }

        defer setIconThemeViaXSettings(name)
        if obj := op.getThemeObject(op.CurrentTheme); obj != nil {
                v := op.setTheme(obj.GtkTheme, name,
                        obj.CursorTheme, obj.FontName, obj.BackgroundFile, obj.SoundTheme)
                op.updateGSettingsKey(GKEY_CURRENT_THEME, v)
                return v, true
        }

        return op.CurrentTheme, false
}

func (op *Manager) SetCursorTheme(name string) (string, bool) {
        if len(name) <= 0 {
                return op.CurrentTheme, false
        }

        defer setCursorThemeViaXSettings(name)
        if obj := op.getThemeObject(op.CurrentTheme); obj != nil {
                v := op.setTheme(obj.GtkTheme, obj.IconTheme,
                        name, obj.FontName, obj.BackgroundFile, obj.SoundTheme)
                op.updateGSettingsKey(GKEY_CURRENT_THEME, v)
                return v, true
        }

        return op.CurrentTheme, false
}

func (op *Manager) SetFontName(name string) (string, bool) {
        if len(name) <= 0 {
                return op.CurrentTheme, false
        }

        defer setGtkFontThemeViaXSettings(name)
        if obj := op.getThemeObject(op.CurrentTheme); obj != nil {
                v := op.setTheme(obj.GtkTheme, obj.IconTheme,
                        obj.CursorTheme, name, obj.BackgroundFile, obj.SoundTheme)
                op.updateGSettingsKey(GKEY_CURRENT_THEME, v)
                return v, true
        }

        return op.CurrentTheme, false
}

func (op *Manager) SetBackgroundFile(name string) (string, bool) {
        if len(name) <= 0 {
                return op.CurrentTheme, false
        }

        name, _, _ = objUtil.PathToFileURI(name)
        if ok, _ := objUtil.IsFileExist(name); !ok {
                return op.CurrentTheme, false
        }

        if !isElementExist(name, op.BackgroundList) {
                // Copy name to Custom dir
                logObject.Infof("Copy '%s' To Custom", name)
                dir := getHomeDir() + THUMB_LOCAL_THEME_PATH + "/Custom/" + PERSON_BG_DIR_NAME
                if ok, _ := objUtil.IsFileExist(dir); !ok {
                        if err := os.MkdirAll(dir, 0755); err != nil {
                                return op.CurrentTheme, false
                        }
                }
                src, _, _ := objUtil.URIToPath(name)
                baseName, _, _ := objUtil.GetBaseName(src)
                path := dir + "/" + baseName
                if ok, _ := objUtil.CopyFile(src, path); !ok {
                        return op.CurrentTheme, false
                }
                name = path
        }

        if obj := op.getThemeObject(op.CurrentTheme); obj != nil {
                v := op.setTheme(obj.GtkTheme, obj.IconTheme,
                        obj.CursorTheme, obj.FontName, name, obj.SoundTheme)
                op.updateGSettingsKey(GKEY_CURRENT_THEME, v)
                return v, true
        }

        // sync value to gsettings
        op.updateGSettingsKey(GKEY_CURRENT_BACKGROUND, name)

        return op.CurrentTheme, false
}

func (op *Manager) SetSoundTheme(name string) (string, bool) {
        if len(name) <= 0 {
                return op.CurrentTheme, false
        }

        if obj := op.getThemeObject(op.CurrentTheme); obj != nil {
                v := op.setTheme(obj.GtkTheme, obj.IconTheme,
                        obj.CursorTheme, obj.FontName, obj.BackgroundFile, name)
                op.updateGSettingsKey(GKEY_CURRENT_THEME, v)
                return v, true
        }

        // sync value to gsettings
        op.updateGSettingsKey(GKEY_CURRENT_SOUND_THEME, name)

        return op.CurrentTheme, false
}

func (op *Manager) setTheme(gtk, icon, cursor, gtkFont, bg, sound string) string {
        for _, path := range op.ThemeList {
                name, ok := isThemeExist(gtk, icon, cursor, gtkFont, bg, sound, path)
                if !ok {
                        continue
                } else {
                        return name
                }
        }

        createTheme("Custom", gtk, icon, cursor, gtkFont, bg, sound)
        op.updateAllProps()
        updateThemeObj(op.pathNameMap)

        return "Custom"
}

func getThemeList() []PathInfo {
        return getThemeThumbList()
}

/*
   Return only contains thumbnail theme
   First, get all local themes
   Then, determine whether contains thumbnail
*/
func getGtkThemeList() []PathInfo {
        valid := getValidGtkThemes()
        thumb := getGtkThumbList()

        list := []PathInfo{}
        for _, v := range valid {
                if isPathInfoInArray(v, thumb) && !isPathInfoInArray(v, list) {
                        list = append(list, v)
                }
        }

        return list
}

/*
   Return only contains thumbnail theme
   First, get all local themes
   Then, determine whether contains thumbnail
*/
func getIconThemeList() []PathInfo {
        valid := getValidIconThemes()
        thumb := getIconThumbList()

        list := []PathInfo{}
        for _, v := range valid {
                if isPathInfoInArray(v, thumb) && !isPathInfoInArray(v, list) {
                        list = append(list, v)
                }
        }

        return list
}

/*
   Return only contains thumbnail theme
   First, get all local themes
   Then, determine whether contains thumbnail
*/
func getCursorThemeList() []PathInfo {
        valid := getValidCursorThemes()
        thumb := getCursorThumbList()

        list := []PathInfo{}
        for _, v := range valid {
                if isPathInfoInArray(v, thumb) && !isPathInfoInArray(v, list) {
                        list = append(list, v)
                }
        }

        return list
}

// Has not yet been determined
func getFontNameList() []string {
        return []string{}
}

/*
   Return all sound theme names.
*/
func getSoundThemeList() []string {
        valid := getValidSoundThemes()
        list := []string{}
        for _, l := range valid {
                list = append(list, l.path)
        }
        return list
}

func isThemeExist(gtk, icon, cursor, gtkFont, bg, sound, path string) (string, bool) {
        obj, ok := themeObjMap[path]
        if !ok {
                return "", false
        }

        if gtk != obj.GtkTheme || icon != obj.IconTheme ||
                cursor != obj.CursorTheme || gtkFont != obj.FontName ||
                obj.BackgroundFile != bg || obj.SoundTheme != sound {
                return "", false
        }

        return obj.Name, true
}

func createTheme(name, gtk, icon, cursor, gtkFont, bg, sound string) bool {
        homeDir := getHomeDir()
        path := homeDir + THUMB_LOCAL_THEME_PATH + "/" + name
        logObject.Infof("Theme Dir:%s", path)
        if ok, _ := objUtil.IsFileExist(path); !ok {
                logObject.Infof("Create Theme Dir: %s", path)
                err := os.MkdirAll(path, 0755)
                if err != nil {
                        logObject.Infof("Mkdir '%s' failed: %v", path, err)
                        return false
                }
        }

        filename := path + "/" + "theme.ini"
        logObject.Infof("Theme Config File:%s", filename)
        if ok, _ := objUtil.IsFileExist(filename); !ok {
                logObject.Infof("Create Theme Config File: %s", filename)
                f, err := os.Create(filename)
                if err != nil {
                        logObject.Infof("Create '%s' failed: %v",
                                filename, err)
                        return false
                }
                f.Close()
        }

        mutex.Lock()
        defer mutex.Unlock()
        keyFile := glib.NewKeyFile()
        defer keyFile.Free()
        ok, err := keyFile.LoadFromFile(filename, glib.KeyFileFlagsKeepComments)
        if !ok {
                logObject.Warningf("LoadKeyFile '%s' failed", filename)
                return false
        }

        keyFile.SetString(THEME_GROUP_THEME, THEME_KEY_NAME, name)
        keyFile.SetString(THEME_GROUP_COMPONENT, THEME_KEY_GTK, gtk)
        keyFile.SetString(THEME_GROUP_COMPONENT, THEME_KEY_ICONS, icon)
        keyFile.SetString(THEME_GROUP_COMPONENT, THEME_KEY_CURSOR, cursor)
        keyFile.SetString(THEME_GROUP_COMPONENT, THEME_KEY_GTK_FONT, gtkFont)
        keyFile.SetString(THEME_GROUP_COMPONENT, THEME_KEY_BG, bg)
        keyFile.SetString(THEME_GROUP_COMPONENT, THEME_KEY_SOUND, sound)

        _, contents, err1 := keyFile.ToData()
        if err1 != nil {
                logObject.Warningf("KeyFile '%s' ToData failed: %s",
                        filename, err)
                return false
        }

        writeDatasToKeyFile(contents, filename)

        return true
}

func writeDatasToKeyFile(contents, filename string) {
        if len(filename) <= 0 {
                return
        }

        f, err := os.Create(filename + "~")
        if err != nil {
                logObject.Warningf("OpenFile '%s' failed: %v",
                        filename+"~", err)
                return
        }
        defer f.Close()

        if _, err = f.WriteString(contents); err != nil {
                logObject.Warningf("Write in '%s' failed: %v",
                        filename+"~", err)
                return
        }
        f.Sync()
        os.Rename(filename+"~", filename)
}

func newManager() *Manager {
        m := &Manager{}

        // TODO similar to updateAllProps()
        m.pathNameMap = make(map[string]PathInfo)
        m.setPropName("ThemeList")
        m.setPropName("GtkThemeList")
        m.setPropName("IconThemeList")
        m.setPropName("CursorThemeList")
        m.setPropName("SoundThemeList")
        m.setPropName("BackgroundList")

        // the following properties should be configure at end for their values
        // depends on other property
        m.setPropName("CurrentTheme")

        m.listenSettingsChanged()
        homeDir := getHomeDir()

        m.listenThemeDir(THEMES_PATH)
        m.listenThemeDir(homeDir + THEMES_LOCAL_PATH)

        m.listenThemeDir(ICONS_PATH)
        m.listenThemeDir(homeDir + ICONS_LOCAL_PATH)

        m.listenThemeDir(THUMB_BASE_PATH)
        m.listenThemeDir(homeDir + THUMB_LOCAL_BASE_PATH)

        m.listenThemeDir(SOUND_THEME_PATH)

        m.listenBackgroundDir(BACKGROUND_DEFAULT_DIR)

        return m
}

func isElementExist(ele string, list []string) bool {
        for _, l := range list {
                if ele == l {
                        return true
                }
        }

        return false
}
