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

	"dbus/org/freedesktop/dbus"
	"dbus/org/freedesktop/login1"
	mpris2 "dbus/org/mpris/mediaplayer2"

	. "pkg.deepin.io/dde/daemon/keybinding/shortcuts"
)

const (
	senderTypeMpris = "org.mpris.MediaPlayer2"
	mprisPath       = "/org/mpris/MediaPlayer2"

	playerDelta int64 = 5000 * 1000 // 5s
)

type MediaPlayerController struct {
	prevPlayer   string
	dbusDaemon   *dbus.DBusDaemon
	loginManager *login1.Manager
}

func NewMediaPlayerController() (*MediaPlayerController, error) {
	c := new(MediaPlayerController)
	var err error
	c.dbusDaemon, err = dbus.NewDBusDaemon("org.freedesktop.DBus", "/")
	if err != nil {
		return nil, err
	}
	c.loginManager, err = login1.NewManager("org.freedesktop.login1", "/org/freedesktop/login1")
	if err != nil {
		return nil, err
	}

	// pause all player before system sleep
	c.loginManager.ConnectPrepareForSleep(func(actived bool) {
		if !actived {
			return
		}
		c.pauseAllPlayer()
	})
	return c, nil
}

func (c MediaPlayerController) Destroy() {
	if c.dbusDaemon != nil {
		dbus.DestroyDBusDaemon(c.dbusDaemon)
		c.dbusDaemon = nil
	}

	if c.loginManager != nil {
		login1.DestroyManager(c.loginManager)
		c.loginManager = nil
	}
}

func (c *MediaPlayerController) Name() string {
	return "Media Player"
}

func (c *MediaPlayerController) ExecCmd(cmd ActionCmd) error {
	player := c.getActiveMpris()
	if player == nil {
		return errors.New("no player found")
	}

	logger.Debug("[HandlerAction] active player dest name:", player.DestName)
	switch cmd {
	case MediaPlayerPlay:
		return player.PlayPause()
	case MediaPlayerPause:
		return player.Pause()
	case MediaPlayerStop:
		return player.Stop()

	case MediaPlayerPrevious:
		if err := player.Previous(); err != nil {
			return err
		}
		return player.Play()

	case MediaPlayerNext:
		if err := player.Next(); err != nil {
			return err
		}
		return player.Play()

	case MediaPlayerRewind:
		pos := player.Position.Get()
		var offset int64
		if pos-playerDelta > 0 {
			offset = -playerDelta
		}

		if err := player.Seek(offset); err != nil {
			return err
		}
		if player.PlaybackStatus.Get() != "Playing" {
			return player.PlayPause()
		}
	case MediaPlayerForword:
		if err := player.Seek(playerDelta); err != nil {
			return err
		}
		if player.PlaybackStatus.Get() != "Playing" {
			return player.PlayPause()
		}
	case MediaPlayerRepeat:
		return player.Play()

	default:
		return ErrInvalidActionCmd{cmd}
	}
	return nil
}

func (c *MediaPlayerController) pauseAllPlayer() {
	for _, sender := range c.getMprisSender() {
		player, err := mpris2.NewPlayer(sender, mprisPath)
		if err != nil {
			continue
		}
		player.Pause()
	}
}

func (c *MediaPlayerController) getMprisSender() []string {
	if c.dbusDaemon == nil {
		return nil
	}
	var senders []string
	names, _ := c.dbusDaemon.ListNames()
	for _, name := range names {
		if !strings.Contains(name, senderTypeMpris) {
			continue
		}

		senders = append(senders, name)
	}

	return senders
}

func (c *MediaPlayerController) getActiveMpris() *mpris2.Player {
	var senders = c.getMprisSender()

	length := len(senders)
	if length == 0 {
		return nil
	}

	for _, sender := range senders {
		player, err := mpris2.NewPlayer(sender, mprisPath)
		if err != nil {
			continue
		}

		if length == 1 {
			return player
		}

		if length == 2 && strings.Contains(sender, "vlc") {
			return player
		}

		if player.PlaybackStatus.Get() == "Playing" {
			c.prevPlayer = sender
			return player
		}

		if c.prevPlayer == sender {
			return player
		}
	}

	player, _ := mpris2.NewPlayer(senders[0], mprisPath)
	return player
}
