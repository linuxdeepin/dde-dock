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
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.audio"
	. "pkg.deepin.io/dde/daemon/keybinding/shortcuts"
	"pkg.deepin.io/lib/dbus1"
)

const (
	volumeMin = 0
	volumeMax = 1.5
)

type AudioController struct {
	conn        *dbus.Conn
	audioDaemon *audio.Audio
}

func NewAudioController(sessionConn *dbus.Conn) *AudioController {
	return &AudioController{
		conn:        sessionConn,
		audioDaemon: audio.NewAudio(sessionConn),
	}
}

func (*AudioController) Name() string {
	return "Audio"
}

func (c *AudioController) ExecCmd(cmd ActionCmd) error {
	switch cmd {
	case AudioSinkMuteToggle:
		return c.toggleSinkMute()

	case AudioSinkVolumeUp:
		return c.changeSinkVolume(true)

	case AudioSinkVolumeDown:
		return c.changeSinkVolume(false)

	case AudioSourceMuteToggle:
		return c.toggleSourceMute()

	default:
		return ErrInvalidActionCmd{cmd}
	}
}

func (c *AudioController) toggleSinkMute() error {
	sink, err := c.getDefaultSink()
	if err != nil {
		return err
	}

	mute, err := sink.Mute().Get(0)
	if err != nil {
		return err
	}

	err = sink.SetMute(0, !mute)
	if err != nil {
		return err
	}
	showOSD("AudioMute")
	return nil
}

func (c *AudioController) toggleSourceMute() error {
	source, err := c.getDefaultSource()
	if err != nil {
		return err
	}

	mute, err := source.Mute().Get(0)
	if err != nil {
		return err
	}
	mute = !mute
	err = source.SetMute(0, mute)
	if err != nil {
		return err
	}

	var osd string
	if mute {
		osd = "AudioMicMuteOn"
	} else {
		osd = "AudioMicMuteOff"
	}
	showOSD(osd)
	return nil
}

func (c *AudioController) changeSinkVolume(raised bool) error {
	sink, err := c.getDefaultSink()
	if err != nil {
		return err
	}

	osd := "AudioUp"
	v, err := sink.Volume().Get(0)
	if err != nil {
		return err
	}

	var step = 0.05
	if !raised {
		step = -step
		osd = "AudioDown"
	}

	v += step
	if v < volumeMin {
		v = volumeMin
	} else if v > volumeMax {
		v = volumeMax
	}

	logger.Debug("[changeSinkVolume] will set volume to:", v)
	mute, err := sink.Mute().Get(0)
	if err != nil {
		return err
	}

	if mute {
		err = sink.SetMute(0, false)
		if err != nil {
			logger.Warning(err)
		}
	}

	err = sink.SetVolume(0, v, true)
	if err != nil {
		return err
	}
	showOSD(osd)
	return nil
}

func (c *AudioController) getDefaultSink() (*audio.Sink, error) {
	sinkPath, err := c.audioDaemon.DefaultSink().Get(0)
	if err != nil {
		return nil, err
	}

	sink, err := audio.NewSink(c.conn, sinkPath)
	if err != nil {
		return nil, err
	}
	return sink, nil
}

func (c *AudioController) getDefaultSource() (*audio.Source, error) {
	sourcePath, err := c.audioDaemon.DefaultSource().Get(0)
	if err != nil {
		return nil, err
	}
	source, err := audio.NewSource(c.conn, sourcePath)
	if err != nil {
		return nil, err
	}

	return source, nil
}
