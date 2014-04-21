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
        "regexp"
)

var (
        watchReset   = make(chan bool)
        prevEventStr string
)

func (op *Manager) watchDirs(dirs []string) {
        for _, dir := range dirs {
                if ok := objUtil.IsFileExist(dir); !ok {
                        err := os.MkdirAll(dir, 0755)
                        if err != nil {
                                logObject.Infof("Make dir '%s' failed: %v", dir, err)
                                continue
                        }
                }
                watcher.Watch(dir)
        }

        for {
                select {
                case ev := <-watcher.Event:
                        if ev == nil {
                                break
                        }
                        if ok, _ := regexp.MatchString(`\.swa?px?$`, ev.Name); ok {
                                break
                        }
                        //if ok, _ := regexp.MatchString(`theme.ini$`, ev.Name); ok {
                        //break
                        //}
                        op.updateAllProps()
                        if ev.IsDelete() || ev.IsCreate() {
                                curEventStr := ev.String()
                                if prevEventStr == curEventStr {
                                        break
                                }
                                watchReset <- true
                                return
                        }
                }
        }
}

func getAllDirName(dir string) []string {
        f, err := os.Open(dir)
        if err != nil {
                logObject.Infof("Open '%s' failed: %v\n", dir, err)
                return []string{}
        }
        defer f.Close()

        finfos, err1 := f.Readdir(0)
        if err1 != nil {
                logObject.Infof("Readdir '%s' failed: %v\n", dir, err1)
                return []string{}
        }

        dirs := []string{}
        dirs = append(dirs, dir)
        for _, info := range finfos {
                if info == nil || !info.IsDir() {
                        continue
                }

                tmp := getAllDirName(dir + "/" + info.Name())
                dirs = append(dirs, tmp...)
        }

        return dirs
}

func (op *Manager) startListenDirs() {
        homeDir := getHomeDir()
        dirs := []string{}

        dirs = append(dirs, THEMES_PATH)
        dirs = append(dirs, homeDir+THEMES_LOCAL_PATH)

        dirs = append(dirs, ICONS_PATH)
        dirs = append(dirs, homeDir+ICONS_LOCAL_PATH)

        dirs = append(dirs, getAllDirName(THUMB_BASE_PATH)...)
        dirs = append(dirs, getAllDirName(homeDir+THUMB_LOCAL_BASE_PATH)...)
        dirs = append(dirs, getAllDirName(SOUND_THEME_PATH)...)

        //logObject.Info("Watch Dirs: ", dirs)
        go op.watchDirs(dirs)
}

func (op *Manager) resetListenDirs() {
        for {
                select {
                case <-watchReset:
                        op.startListenDirs()
                }
        }
}
