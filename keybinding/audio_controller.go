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
	"dbus/com/deepin/daemon/audio"

	ddbus "pkg.deepin.io/dde/daemon/dbus"
	. "pkg.deepin.io/dde/daemon/keybinding/shortcuts"
)

const (
	audioDaemonDest    = "com.deepin.daemon.Audio"
	audioDaemonObjPath = "/com/deepin/daemon/Audio"
)

type AudioController struct {
	audioDaemon *audio.Audio
}

func NewAudioController() (*AudioController, error) {
	c := &AudioController{}
	var err error
	// c.audioDaemon must not be nil
	c.audioDaemon, err = audio.NewAudio(audioDaemonDest, audioDaemonObjPath)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *AudioController) Destroy() {
	if c.audioDaemon != nil {
		audio.DestroyAudio(c.audioDaemon)
		c.audioDaemon = nil
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
	sink.SetMute(!sink.Mute.Get())
	showOSD("AudioMute")
	return nil
}

func (c *AudioController) toggleSourceMute() error {
	source, err := c.getDefaultSource()
	if err != nil {
		return err
	}
	source.SetMute(!source.Mute.Get())
	// TODO: here we can show osd
	return nil
}

func (c *AudioController) changeSinkVolume(raised bool) error {
	sink, err := c.getDefaultSink()
	if err != nil {
		return err
	}

	osd := "AudioUp"
	v := sink.Volume.Get()
	var step float64 = 0.05
	if !raised {
		step = -step
		osd = "AudioDown"
	}

	logger.Debug("[changeSinkVolume] old sink info:", sink.Name.Get(), v)
	v += step
	if v < 0 {
		v = 0
	} else if v > 1 {
		v = 1.0
	}

	logger.Debug("[changeSinkVolume] will set volume to:", v)
	if sink.Mute.Get() {
		sink.SetMute(false)
	}
	sink.SetVolume(v, true)
	showOSD(osd)
	return nil
}

func (c *AudioController) getDefaultSink() (*audio.AudioSink, error) {
	if c.audioDaemon == nil || !ddbus.IsSessionBusActivated(c.audioDaemon.DestName) {
		return nil, ErrIsNil{"AudioController.audioDaemon"}
	}
	sinkPath := c.audioDaemon.DefaultSink.Get()
	sink, err := audio.NewAudioSink(audioDaemonDest, sinkPath)
	if err != nil {
		return nil, err
	}

	return sink, nil
}

func (c *AudioController) getDefaultSource() (*audio.AudioSource, error) {
	if c.audioDaemon == nil || !ddbus.IsSessionBusActivated(c.audioDaemon.DestName) {
		return nil, ErrIsNil{"AudioController.audioDaemon"}
	}
	sourcePath := c.audioDaemon.DefaultSource.Get()
	source, err := audio.NewAudioSource(audioDaemonDest, sourcePath)
	if err != nil {
		return nil, err
	}

	return source, nil
}
