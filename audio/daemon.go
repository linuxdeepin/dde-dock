/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/pulse"
)

var (
	logger = log.NewLogger("daemon/audio")
	_audio *Audio
)

func init() {
	loader.Register(NewAudioDaemon(logger))
}

type Daemon struct {
	*loader.ModuleBase
}

func NewAudioDaemon(logger *log.Logger) *Daemon {
	var d = new(Daemon)
	d.ModuleBase = loader.NewModuleBase("audio", d, logger)
	return d
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func finalize() {
	_audio.destroy()
	_audio = nil
	logger.EndTracing()
}

func (*Daemon) Start() error {
	if _audio != nil {
		return nil
	}

	logger.BeginTracing()

	ctx := pulse.GetContext()
	if ctx == nil {
		logger.Error("Failed to connect pulseaudio server")
		logger.EndTracing()
		return nil
	}

	_audio = NewAudio(ctx)

	if err := dbus.InstallOnSession(_audio); err != nil {
		logger.Error("Failed InstallOnSession:", err)
		finalize()
		return err
	}

	initDefaultVolume(_audio)
	return nil
}

func (*Daemon) Stop() error {
	if _audio == nil {
		return nil
	}

	finalize()
	return nil
}
