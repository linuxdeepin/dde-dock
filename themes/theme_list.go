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
        "os"
        "os/user"
)

type PathInfo struct {
        path    string
        t       string  // 'system' or 'local'
}

const (
        THEMES_PATH             = "/usr/share/themes/"
        THEMES_LOCAL_PATH       = "/.themes/"
        ICONS_LOCAL_PATH        = "/.icons/"
        ICONS_PATH              = "/usr/share/icons/"
        PATH_TYPE_SYSTEM        = "system"
        PATH_TYPE_LOCAL         = "local"
        THUMB_BASE_PATH         = "/usr/share/deepin-personalization/"
        THUMB_THEME_PATH        = THUMB_BASE_PATH + "themes"
        THUMB_GTK_PATH          = THUMB_BASE_PATH + "gtk"
        THUMB_ICON_PATH         = THUMB_BASE_PATH + "icons"
        THUMB_CURSOR_PATH       = THUMB_BASE_PATH + "cursor"
        THUMB_LOCAL_BASE_PATH   = "/.deepin-personalization/"
        THUMB_LOCAL_THEME_PATH  = THUMB_LOCAL_BASE_PATH + "themes"
        THUMB_LOCAL_GTK_PATH    = THUMB_LOCAL_BASE_PATH + "gtk"
        THUMB_LOCAL_ICON_PATH   = THUMB_LOCAL_BASE_PATH + "icons"
        THUMB_LOCAL_CURSOR_PATH = THUMB_LOCAL_BASE_PATH + "cursor"
)

func getHomeDir() string {
        u, err := user.Current()
        if err != nil {
                logObject.Info("Get current user info failed:%v", err)
                panic(err)
        }
        return u.HomeDir
}

func getValidGtkThemes() []PathInfo {
        localDir := getHomeDir() + THEMES_LOCAL_PATH
        sysDirs := []PathInfo{PathInfo{THEMES_PATH, PATH_TYPE_SYSTEM}}
        localDirs := []PathInfo{PathInfo{localDir, PATH_TYPE_LOCAL}}
        conditions := []string{"gtk-2.0", "gtk-3.0", "metacity-1"}

        sysList := getValidThemes(sysDirs, conditions)
        localList := getValidThemes(localDirs, conditions)
        for _, l := range sysList {
                if isElementExist(l, localList) {
                        continue
                }
                localList = append(localList, l)
        }
        //logObject.Info("Valid Gtk Theme: %v", localList)

        return localList
}

func getValidIconThemes() []PathInfo {
        localDir := getHomeDir() + ICONS_LOCAL_PATH
        sysDirs := []PathInfo{PathInfo{ICONS_PATH, PATH_TYPE_SYSTEM}}
        localDirs := []PathInfo{PathInfo{localDir, PATH_TYPE_LOCAL}}
        conditions := []string{"index.theme"}

        sysList := getValidThemes(sysDirs, conditions)
        localList := getValidThemes(localDirs, conditions)
        for _, l := range sysList {
                if isElementExist(l, localList) {
                        continue
                }
                localList = append(localList, l)
        }
        //logObject.Info("Valid Icon Theme: %v", localList)

        return localList
}

func getValidCursorThemes() []PathInfo {
        localDir := getHomeDir() + ICONS_LOCAL_PATH
        sysDirs := []PathInfo{PathInfo{ICONS_PATH, PATH_TYPE_SYSTEM}}
        localDirs := []PathInfo{PathInfo{localDir, PATH_TYPE_LOCAL}}
        conditions := []string{"cursors"}

        sysList := getValidThemes(sysDirs, conditions)
        localList := getValidThemes(localDirs, conditions)
        for _, l := range sysList {
                if isElementExist(l, localList) {
                        continue
                }
                localList = append(localList, l)
        }
        //logObject.Info("Valid Cursor Theme: %v", localList)

        return localList
}

func getValidThemes(dirs []PathInfo, conditions []string) []PathInfo {
        valid := []PathInfo{}
        for _, dir := range dirs {
                f, err := os.Open(dir.path)
                if err != nil {
                        logObject.Info("Open '%s' failed: %s\n",
                                dir.path, err)
                        continue
                }
                defer f.Close()

                infos, err1 := f.Readdir(0)
                if err1 != nil {
                        logObject.Info("ReadDir '%s' failed: %s\n",
                                dir.path, err)
                        continue
                }

                for _, info := range infos {
                        if !info.IsDir() {
                                continue
                        }

                        if filterTheme(dir.path+info.Name(),
                                conditions) {
                                tmp := PathInfo{info.Name(), dir.t}
                                valid = append(valid, tmp)
                        }
                }
        }

        return valid
}

func filterTheme(dir string, conditions []string) bool {
        f, err := os.Open(dir)
        if err != nil {
                logObject.Info("Open '%s' failed: %s\n", dir, err)
                return false
        }
        defer f.Close()

        names, err1 := f.Readdirnames(0)
        if err1 != nil {
                logObject.Info("ReadDir '%s' failed: %s\n", dir, err)
                return false
        }

        cnt := 0
        for _, name := range names {
                for _, condition := range conditions {
                        if name == condition {
                                cnt++
                                break
                        }
                }
        }

        if cnt == len(conditions) {
                return true
        }

        return false
}

func getThemeThumbList() []PathInfo {
        sysDirs := []PathInfo{PathInfo{THUMB_THEME_PATH + "/", PATH_TYPE_SYSTEM}}
        localDir := getHomeDir() + THUMB_LOCAL_THEME_PATH + "/"
        localDirs := []PathInfo{PathInfo{localDir, PATH_TYPE_LOCAL}}
        conditions := []string{"theme.ini"}

        //sysList := getThumbList(sysDirs)
        //localList := getThumbList(localDirs)
        sysList := getValidThemes(sysDirs, conditions)
        localList := getValidThemes(localDirs, conditions)
        //logObject.Info("System Theme List: %v\n", sysList)
        //logObject.Info("Local Theme List: %v\n", localList)
        for _, l := range sysList {
                if isElementExist(l, localList) {
                        continue
                }
                localList = append(localList, l)
        }
        //logObject.Info("Gtk Thumb List:%v", localList)

        return localList
}

func getGtkThumbList() []PathInfo {
        sysDirs := []PathInfo{PathInfo{THUMB_GTK_PATH, PATH_TYPE_SYSTEM}}
        localDir := getHomeDir() + THUMB_LOCAL_GTK_PATH
        localDirs := []PathInfo{PathInfo{localDir, PATH_TYPE_LOCAL}}

        sysList := getThumbList(sysDirs)
        localList := getThumbList(localDirs)
        for _, l := range sysList {
                if isElementExist(l, localList) {
                        continue
                }
                localList = append(localList, l)
        }
        //logObject.Info("Gtk Thumb List:%v", localList)

        return localList
}

func getIconThumbList() []PathInfo {
        sysDirs := []PathInfo{PathInfo{THUMB_ICON_PATH, PATH_TYPE_SYSTEM}}
        localDir := getHomeDir() + THUMB_LOCAL_ICON_PATH
        localDirs := []PathInfo{PathInfo{localDir, PATH_TYPE_LOCAL}}

        sysList := getThumbList(sysDirs)
        localList := getThumbList(localDirs)
        for _, l := range sysList {
                if isElementExist(l, localList) {
                        continue
                }
                localList = append(localList, l)
        }
        //logObject.Info("Icon Thumb List:%v", localList)

        return localList
}

func getCursorThumbList() []PathInfo {
        sysDirs := []PathInfo{PathInfo{THUMB_CURSOR_PATH, PATH_TYPE_SYSTEM}}
        localDir := getHomeDir() + THUMB_LOCAL_CURSOR_PATH
        localDirs := []PathInfo{PathInfo{localDir, PATH_TYPE_LOCAL}}

        sysList := getThumbList(sysDirs)
        localList := getThumbList(localDirs)
        for _, l := range sysList {
                if isElementExist(l, localList) {
                        continue
                }
                localList = append(localList, l)
        }
        //logObject.Info("Cursor Thumb List:%v", localList)

        return localList
}

func getThumbList(dirs []PathInfo) []PathInfo {
        list := []PathInfo{}
        for _, dir := range dirs {
                f, err := os.Open(dir.path)
                if err != nil {
                        logObject.Info("Open '%s' failed: %v",
                                dir.path, err)
                        return list
                }

                infos, err1 := f.Readdir(0)
                if err1 != nil {
                        logObject.Info("ReadDir '%s' failed: %v",
                                dir.path, err)
                        return list
                }

                for _, info := range infos {
                        if info.IsDir() {
                                tmp := PathInfo{info.Name(), dir.t}
                                list = append(list, tmp)
                        }
                }
        }

        return list
}

func isElementExist(ele PathInfo, list []PathInfo) bool {
        for _, e := range list {
                if ele.path == e.path {
                        return true
                }
        }

        return false
}
