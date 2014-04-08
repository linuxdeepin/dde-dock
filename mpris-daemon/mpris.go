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
        libmpris "dbus/org/mpris/mediaplayer2"
        "strings"
)

const (
        MPRIS_FILTER_KEY = "org.mpris.MediaPlayer2"
        MPRIS_PATH       = "/org/mpris/MediaPlayer2"
        SEEK_DISTANCE    = int64(10000000) // 10s
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
                if obj.PlaybackStatus.GetValue().(string) == "Playing" {
                        return obj
                }
        }

        return nil
}

func listenAudioSignal() {
        mediaKeyObj.ConnectAudioPlay(func(press bool) {
                obj := getActiveMprisClient()
                if obj == nil {
                        return
                }
                obj.Play()
        })

        mediaKeyObj.ConnectAudioPause(func(press bool) {
                obj := getActiveMprisClient()
                if obj == nil {
                        return
                }
                obj.Pause()
        })

        mediaKeyObj.ConnectAudioStop(func(press bool) {
                obj := getActiveMprisClient()
                if obj == nil {
                        return
                }
                obj.Stop()
        })

        mediaKeyObj.ConnectAudioPrevious(func(press bool) {
                obj := getActiveMprisClient()
                if obj == nil {
                        return
                }
                obj.Previous()
        })

        mediaKeyObj.ConnectAudioNext(func(press bool) {
                obj := getActiveMprisClient()
                if obj == nil {
                        return
                }
                obj.Next()
        })

        mediaKeyObj.ConnectAudioRewind(func(press bool) {
                obj := getActiveMprisClient()
                if obj == nil {
                        return
                }
                obj.Seek(obj.Position.GetValue().(int64) - SEEK_DISTANCE)
        })

        mediaKeyObj.ConnectAudioForward(func(press bool) {
                obj := getActiveMprisClient()
                if obj == nil {
                        return
                }
                obj.Seek(obj.Position.GetValue().(int64) + SEEK_DISTANCE)
        })

        mediaKeyObj.ConnectAudioRepeat(func(press bool) {
                obj := getActiveMprisClient()
                if obj == nil {
                        return
                }
                obj.Play()
        })
}
