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

package screenedge

import (
	libarea "dbus/com/deepin/api/xmousearea"
	libdsp "dbus/com/deepin/daemon/display"
	"dbus/com/deepin/dde/launcher"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/gio-2.0"
	"pkg.deepin.io/lib/log"
	"sync"
)

var (
	dspObj    *libdsp.Display
	areaObj   *libarea.XMouseArea
	launchObj *launcher.Launcher
	logger    = log.NewLogger("daemon/screenedge")

	mutex         = new(sync.Mutex)
	edgeActionMap = make(map[string]string)
)
var zoneSettings = func() func() *gio.Settings {
	var initZoneSettings sync.Once
	var _zoneSettings *gio.Settings

	return func() *gio.Settings {
		initZoneSettings.Do(func() {
			_zoneSettings = gio.NewSettings("com.deepin.dde.zone")
		})
		return _zoneSettings
	}
}()

// Enable desktop edge zone detected
//
// 是否启用桌面边缘热区功能
func (op *Manager) EnableZoneDetected(enable bool) {
	if enable {
		unregisterZoneArea()
		registerZoneArea()
	} else {
		unregisterZoneArea()
	}
}

// Set left-top edge action
func (op *Manager) SetTopLeft(value string) {
	mutex.Lock()
	defer mutex.Unlock()
	zoneSettings().SetString("left-up", value)
}

// Get left-top edge action
func (op *Manager) TopLeftAction() string {
	return zoneSettings().GetString("left-up")
}

// Set left-bottom edge action
func (op *Manager) SetBottomLeft(value string) {
	mutex.Lock()
	defer mutex.Unlock()
	zoneSettings().SetString("left-down", value)
}

// Get left-bottom edge action
func (op *Manager) BottomLeftAction() string {
	return zoneSettings().GetString("left-down")
}

// Set right-top edge action
func (op *Manager) SetTopRight(value string) {
	mutex.Lock()
	defer mutex.Unlock()
	zoneSettings().SetString("right-up", value)
}

// Get right-top edge action
func (op *Manager) TopRightAction() string {
	return zoneSettings().GetString("right-up")
}

// Set right-bottom edge action
func (op *Manager) SetBottomRight(value string) {
	mutex.Lock()
	defer mutex.Unlock()
	zoneSettings().SetString("right-down", value)
}

// Get right-bottom edge action
func (op *Manager) BottomRightAction() string {
	return zoneSettings().GetString("right-down")
}

func initDBusIFC() error {
	var err error
	dspObj, err = libdsp.NewDisplay("com.deepin.daemon.Display",
		"/com/deepin/daemon/Display")
	if err != nil {
		finalizeDBusIFC()
		return err
	}

	areaObj, err = libarea.NewXMouseArea("com.deepin.api.XMouseArea",
		"/com/deepin/api/XMouseArea")
	if err != nil {
		finalizeDBusIFC()
		return err
	}

	launchObj, err = launcher.NewLauncher("com.deepin.dde.launcher",
		"/com/deepin/dde/launcher")
	if err != nil {
		finalizeDBusIFC()
		return err
	}

	return nil
}

func finalizeDBusIFC() {
	if dspObj != nil {
		libdsp.DestroyDisplay(dspObj)
		dspObj = nil
	}

	if areaObj != nil {
		libarea.DestroyXMouseArea(areaObj)
		areaObj = nil
	}

	if launchObj != nil {
		launcher.DestroyLauncher(launchObj)
		launchObj = nil
	}
}

var _m *Manager

func finalize() {
	finalizeDBusIFC()
	_m.destroy()
	_m = nil
}

type Daemon struct {
	*loader.ModuleBase
}

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("screenedge", daemon, logger)
	return daemon
}

func (d *Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Stop() error {
	if _m == nil {
		return nil
	}

	finalize()
	logger.EndTracing()
	return nil
}

func (m *Manager) initEdgeActionMap() {
	mutex.Lock()
	defer mutex.Unlock()
	edgeActionMap[leftTopEdge] = zoneSettings().GetString("left-up")
	edgeActionMap[leftBottomEdge] = zoneSettings().GetString("left-down")
	edgeActionMap[rightTopEdge] = zoneSettings().GetString("right-up")
	edgeActionMap[rightBottomEdge] = zoneSettings().GetString("right-down")
}

func (d *Daemon) Start() error {
	if _m != nil {
		return nil
	}

	logger.BeginTracing()

	err := initDBusIFC()
	if err != nil {
		logger.Error("Create dbus interface failed:", err)
		logger.EndTracing()
		return err
	}

	_m = newManager()
	err = dbus.InstallOnSession(_m)
	if err != nil {
		logger.Error("Install Zone Session Failed: ", err)
		finalize()
		return err
	}

	_m.initEdgeActionMap()
	handleSettingsChanged()

	return nil
}
