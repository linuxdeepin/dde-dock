/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package audio

import (
	libdbus "dbus/org/freedesktop/dbus"
	mpris2 "dbus/org/mpris/mediaplayer2"
	"encoding/json"
	"pkg.deepin.io/dde/api/soundutils"
	"pkg.deepin.io/lib/pulse"
	"strings"
)

func isVolumeValid(v float64) bool {
	if v < 0 || v > pulse.VolumeUIMax {
		return false
	}
	return true
}

func playFeedback() {
	playFeedbackWithDevice("")
}

func playFeedbackWithDevice(device string) {
	soundutils.PlaySystemSound(soundutils.EventVolumeChanged, device, false)
}

func toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}

const (
	mprisPlayerDestPrefix = "org.mpris.MediaPlayer2"
	mprisPlayerObjPath    = "/org/mpris/MediaPlayer2"
)

func getMprisPlayers() ([]string, error) {
	dbusDaemon, err := libdbus.NewDBusDaemon("org.freedesktop.DBus", "/")
	if err != nil {
		return nil, err
	}
	defer libdbus.DestroyDBusDaemon(dbusDaemon)

	var playerDests []string
	names, err := dbusDaemon.ListNames()
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		if strings.HasPrefix(name, mprisPlayerDestPrefix) {
			// is mpris player
			playerDests = append(playerDests, name)
		}
	}
	return playerDests, nil
}

func pauseAllPlayers() {
	playerDests, err := getMprisPlayers()
	if err != nil {
		logger.Warning("getMprisPlayers failed:", err)
		return
	}

	logger.Debug("pause all players")
	for _, playerDest := range playerDests {
		player, err := mpris2.NewPlayer(playerDest, mprisPlayerObjPath)
		if err != nil {
			continue
		}
		defer mpris2.DestroyPlayer(player)
		player.Pause()
	}
}
