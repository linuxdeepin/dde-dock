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

package main

import (
        libkeybind "dbus/com/deepin/daemon/keybinding"
        libdbus "dbus/org/freedesktop/dbus"
        libmpris "dbus/org/mpris/mediaplayer2"
        "dlib/gio-2.0"
        "os/exec"
        "strings"
)

const (
        MPRIS_FILTER_KEY = "org.mpris.MediaPlayer2"
        MPRIS_PATH       = "/org/mpris/MediaPlayer2"
        SEEK_DISTANCE    = int64(5000000) // 5s

        MIME_TYPE_BROWSER = "x-scheme-handler/http"
        MIME_TYPE_EMAIL   = "x-scheme-handler/mailto"
        CALCULATOR_CMD    = "/usr/bin/gnome-calculator"
)

var (
        dbusObj     *libdbus.DBusDaemon
        mediaKeyObj *libkeybind.MediaKey
)

func getMprisClients() ([]string, bool) {
        list := []string{}
        names, err := dbusObj.ListNames()
        if err != nil {
                logObj.Info("List DBus Sender Names: ", err)
                return list, false
        }

        for _, name := range names {
                if strings.Contains(name, MPRIS_FILTER_KEY) {
                        list = append(list, name)
                }
        }

        return list, true
}

func getActiveMprisClient() *libmpris.Player {
        list, ok := getMprisClients()
        if !ok {
                return nil
        }

        for _, dest := range list {
                obj, err := libmpris.NewPlayer(dest, MPRIS_PATH)
                if err != nil {
                        logObj.Warningf("New mpris player failed: '%v'' for sender: '%s'", err, dest)
                        continue
                }
                if len(list) == 1 {
                        return obj
                }
                if obj.PlaybackStatus.GetValue().(string) == "Playing" {
                        return obj
                }
        }

        return nil
}

func listenAudioSignal() {
        mediaKeyObj.ConnectAudioPlay(func(press bool) {
                if press {
                        return
                }

                obj := getActiveMprisClient()
                if obj == nil {
                        return
                }
                //obj.Play()
                obj.PlayPause()
        })

        mediaKeyObj.ConnectAudioPause(func(press bool) {
                if press {
                        return
                }

                obj := getActiveMprisClient()
                if obj == nil {
                        return
                }
                obj.Pause()
        })

        mediaKeyObj.ConnectAudioStop(func(press bool) {
                if press {
                        return
                }

                obj := getActiveMprisClient()
                if obj == nil {
                        return
                }
                obj.Stop()
        })

        mediaKeyObj.ConnectAudioPrevious(func(press bool) {
                if press {
                        return
                }

                obj := getActiveMprisClient()
                if obj == nil {
                        return
                }
                obj.Previous()
        })

        mediaKeyObj.ConnectAudioNext(func(press bool) {
                if press {
                        return
                }

                obj := getActiveMprisClient()
                if obj == nil {
                        return
                }
                obj.Next()
        })

        mediaKeyObj.ConnectAudioRewind(func(press bool) {
                if press {
                        return
                }

                obj := getActiveMprisClient()
                if obj == nil {
                        return
                }
                pos := obj.Position.GetValue().(int64)
                //println("Current Position: ", pos)
                nextPos := pos - SEEK_DISTANCE
                if nextPos < 0 {
                        nextPos = 0
                } else {
                        nextPos = 0 - SEEK_DISTANCE
                }
                //println("Rewind Position: ", nextPos)
                obj.Seek(nextPos)
                if obj.PlaybackStatus.GetValue().(string) != "Playing" {
                        obj.PlayPause()
                }
        })

        mediaKeyObj.ConnectAudioForward(func(press bool) {
                if press {
                        return
                }

                obj := getActiveMprisClient()
                if obj == nil {
                        return
                }
                //pos := obj.Position.GetValue().(int64)
                //println("Current Position: ", pos)
                //nextPos := pos + SEEK_DISTANCE
                //println("Forward Position: ", nextPos)
                obj.Seek(SEEK_DISTANCE)
                if obj.PlaybackStatus.GetValue().(string) != "Playing" {
                        obj.PlayPause()
                }
        })

        mediaKeyObj.ConnectAudioRepeat(func(press bool) {
                if press {
                        return
                }

                obj := getActiveMprisClient()
                if obj == nil {
                        return
                }
                obj.Play()
        })

        mediaKeyObj.ConnectLaunchEmail(func(press bool) {
                if press {
                        return
                }

                if cmd, ok := getCommandByMimeType(MIME_TYPE_EMAIL); ok {
                        go exec.Command(cmd).Run()
                }
        })

        mediaKeyObj.ConnectLaunchBrowser(func(press bool) {
                if press {
                        return
                }

                if cmd, ok := getCommandByMimeType(MIME_TYPE_BROWSER); ok {
                        go exec.Command(cmd).Run()
                }
        })

        mediaKeyObj.ConnectLaunchCalculator(func(press bool) {
                if press {
                        return
                }

                go exec.Command(CALCULATOR_CMD).Run()
        })
}

func startMprisDaemon() {
        var err error

        dbusObj, err = libdbus.NewDBusDaemon("org.freedesktop.DBus",
                "/")
        if err != nil {
                logObj.Info("New DBusDaemon Failed: ", err)
                panic(err)
        }

        mediaKeyObj, err = libkeybind.NewMediaKey(
                "com.deepin.daemon.KeyBinding",
                "/com/deepin/daemon/MediaKey")
        if err != nil {
                logObj.Info("New MediaKey Object Failed: ", err)
                panic(err)
        }

        listenAudioSignal()
}

func getCommandByMimeType(mimeType string) (string, bool) {
        if appInfo := gio.AppInfoGetDefaultForType(mimeType, false); appInfo != nil {
                return appInfo.GetExecutable(), true
        }

        return "", false
}
