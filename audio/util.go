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

package audio

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math"
	"strings"
	"unicode"

	ofdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.dbus"
	mpris2 "github.com/linuxdeepin/go-dbus-factory/org.mpris.mediaplayer2"
	"pkg.deepin.io/dde/api/soundutils"
	"pkg.deepin.io/lib/dbus1"
	//"pkg.deepin.io/lib/pulse"
)

func isVolumeValid(v float64) bool {
	if v < 0 || v > gMaxUIVolume {
		return false
	}
	return true
}

func playFeedback() {
	playFeedbackWithDevice("")
}

func playFeedbackWithDevice(device string) {
	go soundutils.PlaySystemSound(soundutils.EventAudioVolumeChanged, device)
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
)

func getMprisPlayers(sessionConn *dbus.Conn) ([]string, error) {
	var playerNames []string
	dbusDaemon := ofdbus.NewDBus(sessionConn)
	names, err := dbusDaemon.ListNames(0)
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		if strings.HasPrefix(name, mprisPlayerDestPrefix) {
			// is mpris player
			playerNames = append(playerNames, name)
		}
	}
	return playerNames, nil
}

func pauseAllPlayers() {
	sessionConn, err := dbus.SessionBus()
	if err != nil {
		return
	}
	playerNames, err := getMprisPlayers(sessionConn)
	if err != nil {
		logger.Warning("getMprisPlayers failed:", err)
		return
	}

	logger.Debug("pause all players")
	for _, playerName := range playerNames {
		player := mpris2.NewMediaPlayer(sessionConn, playerName)
		err := player.Pause(0)
		if err != nil {
			logger.Warningf("failed to pause player %s: %v", playerName, err)
		}
	}
}

// 四舍五入
func floatPrecision(f float64) float64 {
	// 精确到小数点后2位
	pow10N := math.Pow10(2)
	return math.Trunc((f+0.5/pow10N)*pow10N) / pow10N
	// return math.Trunc((f)*pow10N) / pow10N
}

const defaultPaFile = "/etc/pulse/default.pa"

func loadDefaultPaConfig(filename string) (cfg defaultPaConfig) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Warning(err)
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Bytes()
		line = bytes.TrimLeftFunc(line, unicode.IsSpace)
		if bytes.HasPrefix(line, []byte{'#'}) {
			continue
		}

		if bytes.Contains(line, []byte("set-default-sink")) {
			cfg.setDefaultSink = true
		}
		if bytes.Contains(line, []byte("set-default-source")) {
			cfg.setDefaultSource = true
		}
	}
	err = scanner.Err()
	if err != nil {
		logger.Warning(err)
	}
	return
}

type defaultPaConfig struct {
	setDefaultSource bool
	setDefaultSink   bool
}
