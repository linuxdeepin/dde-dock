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
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/pulse"
)

var (
	logger = log.NewLogger("daemon/audio")
)

func init() {
	loader.Register(NewAudioDaemon(logger))
}

type Daemon struct {
	*loader.ModuleBase
	audio *Audio
}

func NewAudioDaemon(logger *log.Logger) *Daemon {
	var d = new(Daemon)
	d.ModuleBase = loader.NewModuleBase("audio", d, logger)
	return d
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Start() error {
	if d.audio != nil {
		return nil
	}

	service := loader.GetService()
	ctx := pulse.GetContext()
	if ctx == nil {
		logger.Error("failed to connect pulseaudio server")
		return nil
	}

	d.audio = newAudio(ctx, service)

	err := service.Export(dbusPath, d.audio)
	if err != nil {
		return err
	}
	err = service.RequestName(dbusServiceName)
	if err != nil {
		return err
	}

	go d.audio.handleEvent()
	go d.audio.handleStateChanged()
	initDefaultVolume(d.audio)
	return nil
}

func (d *Daemon) Stop() error {
	if d.audio == nil {
		return nil
	}

	d.audio.destroy()

	service := loader.GetService()
	err := service.StopExport(d.audio)
	if err != nil {
		logger.Warning(err)
	}

	err = service.ReleaseName(dbusServiceName)
	if err != nil {
		logger.Warning(err)
	}

	d.audio = nil
	return nil
}
