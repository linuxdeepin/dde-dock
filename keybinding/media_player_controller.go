/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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

package keybinding

import (
	"errors"
	"strings"

	ofdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.dbus"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.login1"
	mpris2 "github.com/linuxdeepin/go-dbus-factory/org.mpris.mediaplayer2"
	. "pkg.deepin.io/dde/daemon/keybinding/shortcuts"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
)

const (
	senderTypeMpris = "org.mpris.MediaPlayer2"

	playerDelta int64 = 5000 * 1000 // 5s
)

type MediaPlayerController struct {
	conn         *dbus.Conn
	prevPlayer   string
	dbusDaemon   *ofdbus.DBus
	loginManager *login1.Manager
}

func NewMediaPlayerController(systemSigLoop *dbusutil.SignalLoop,
	sessionConn *dbus.Conn) *MediaPlayerController {

	c := new(MediaPlayerController)
	c.conn = sessionConn
	c.dbusDaemon = ofdbus.NewDBus(sessionConn)
	c.loginManager = login1.NewManager(systemSigLoop.Conn())
	c.loginManager.InitSignalExt(systemSigLoop, true)

	// pause all player before system sleep
	c.loginManager.ConnectPrepareForSleep(func(start bool) {
		if !start {
			return
		}
		c.pauseAllPlayer()
	})
	return c
}

func (c *MediaPlayerController) Destroy() {
	c.loginManager.RemoveHandler(proxy.RemoveAllHandlers)
}

func (c *MediaPlayerController) Name() string {
	return "Media Player"
}

func (c *MediaPlayerController) ExecCmd(cmd ActionCmd) error {
	player := c.getActiveMpris()
	if player == nil {
		return errors.New("no player found")
	}

	logger.Debug("[HandlerAction] active player dest name:", player.ServiceName_())
	switch cmd {
	case MediaPlayerPlay:
		return player.PlayPause(0)
	case MediaPlayerPause:
		return player.Pause(0)
	case MediaPlayerStop:
		return player.Stop(0)

	case MediaPlayerPrevious:
		if err := player.Previous(0); err != nil {
			return err
		}
		return player.Play(0)

	case MediaPlayerNext:
		if err := player.Next(0); err != nil {
			return err
		}
		return player.Play(0)

	case MediaPlayerRewind:
		pos, err := player.Position().Get(0)
		if err != nil {
			return err
		}

		var offset int64
		if pos-playerDelta > 0 {
			offset = -playerDelta
		}

		if err := player.Seek(0, offset); err != nil {
			return err
		}
		status, err := player.PlaybackStatus().Get(0)
		if err != nil {
			return err
		}
		if status != "Playing" {
			return player.PlayPause(0)
		}
	case MediaPlayerForword:
		if err := player.Seek(0, playerDelta); err != nil {
			return err
		}
		status, err := player.PlaybackStatus().Get(0)
		if err != nil {
			return err
		}
		if status != "Playing" {
			return player.PlayPause(0)
		}
	case MediaPlayerRepeat:
		return player.Play(0)

	default:
		return ErrInvalidActionCmd{cmd}
	}
	return nil
}

func (c *MediaPlayerController) pauseAllPlayer() {
	for _, sender := range c.getMprisSender() {
		player := mpris2.NewMediaPlayer(c.conn, sender)
		err := player.Pause(0)
		if err != nil {
			logger.Warning(err)
		}
	}
}

func (c *MediaPlayerController) getMprisSender() []string {
	if c.dbusDaemon == nil {
		return nil
	}
	var senders []string
	names, _ := c.dbusDaemon.ListNames(0)
	for _, name := range names {
		if strings.HasPrefix(name, senderTypeMpris) {
			senders = append(senders, name)
		}
	}

	return senders
}

func (c *MediaPlayerController) getActiveMpris() *mpris2.MediaPlayer {
	var senders = c.getMprisSender()

	length := len(senders)
	if length == 0 {
		return nil
	}

	for _, sender := range senders {
		player := mpris2.NewMediaPlayer(c.conn, sender)

		if length == 1 {
			return player
		}

		if length == 2 && strings.Contains(sender, "vlc") {
			return player
		}

		status, err := player.PlaybackStatus().Get(0)
		if err != nil {
			logger.Warning(err)
			continue
		}

		if status == "Playing" {
			c.prevPlayer = sender
			return player
		}

		if c.prevPlayer == sender {
			return player
		}
	}

	player := mpris2.NewMediaPlayer(c.conn, senders[0])
	return player
}
