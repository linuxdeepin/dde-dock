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
        "github.com/howeyc/fsnotify"
        "os"
        "regexp"
)

var (
        watchReset = make(chan bool)
        fsWatchMap = make(map[string]*fsnotify.Watcher)
)

func (op *Manager) listenThemeDir(dir string) {
        //logObject.Infof("Listen Dir: %s\n", dir)
        watcher, err := fsnotify.NewWatcher()
        if err != nil {
                logObject.Infof("Create new watch failed: %v", err)
                return
        }

        if ok := objUtil.IsFileExist(dir); !ok {
                err = os.MkdirAll(dir, 0755)
                if err != nil {
                        logObject.Infof("Make dir '%s' failed: %v", dir, err)
                        return
                }
        }

        err = watcher.Watch(dir)
        if err != nil {
                logObject.Infof("Watch '%s' failed: %v", dir, err)
                return
        }
        fsWatchMap[dir] = watcher

        go func() {
                defer func() {
                        if watcher != nil {
                                watcher.Close()
                        }
                        delete(fsWatchMap, dir)
                }()

                for {
                        select {
                        case ev := <-watcher.Event:
                                if ev == nil {
                                        break
                                }
                                if ok, _ := regexp.MatchString(`\.swa?px?$`, ev.Name); ok {
                                        break
                                }
                                op.updateAllProps()
                                //if ev.IsDelete() || ev.IsCreate() {
                                //watchReset <- true
                                //}
                        case err := <-watcher.Error:
                                if err != nil {
                                        logObject.Warningf("Watch Error: %v", err)
                                }
                        }
                }
        }()
}

func (op *Manager) recursionListenDir(dir string) {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Warning("Recover Error: ", err)
                        return
                }
        }()

        if ok := objUtil.IsFileExist(dir); !ok {
                err := os.MkdirAll(dir, 0755)
                if err != nil {
                        logObject.Infof("Make dir '%s' failed: %v", dir, err)
                        return
                }
        }
        op.listenThemeDir(dir)

        f, err := os.Open(dir)
        if err != nil {
                logObject.Warningf("Open '%s' failed: %v", dir, err)
                return
        }

        finfos, err1 := f.Readdir(0)
        if err1 != nil {
                logObject.Warningf("ReadDir '%s' failed: %v", dir, err1)
                return
        }

        for _, finfo := range finfos {
                if finfo.IsDir() {
                        op.recursionListenDir(dir + "/" + finfo.Name())
                }
        }
}

func (op *Manager) startListenDirs() {
        homeDir := getHomeDir()

        op.listenThemeDir(THEMES_PATH)
        op.listenThemeDir(homeDir + THEMES_LOCAL_PATH)

        op.listenThemeDir(ICONS_PATH)
        op.listenThemeDir(homeDir + ICONS_LOCAL_PATH)

        op.recursionListenDir(THUMB_BASE_PATH)
        op.recursionListenDir(homeDir + THUMB_LOCAL_BASE_PATH)
        op.recursionListenDir(SOUND_THEME_PATH)
}

func (op *Manager) resetListenDirs() {
        for {
                select {
                case <-watchReset:
                        for k, v := range fsWatchMap {
                                if v != nil {
                                        v.Close()
                                }
                                delete(fsWatchMap, k)
                        }
                        op.startListenDirs()
                }
        }
}
