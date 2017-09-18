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

package inputdevices

type devicePathInfo struct {
	Path string
	Type string
}
type devicePathInfos []*devicePathInfo

type Manager struct {
	Infos devicePathInfos

	kbd        *Keyboard
	mouse      *Mouse
	trackPoint *TrackPoint
	tpad       *Touchpad
	wacom      *Wacom
}

func NewManager() *Manager {
	var m = new(Manager)

	m.Infos = devicePathInfos{
		&devicePathInfo{
			Path: "com.deepin.daemon.InputDevice.Keyboard",
			Type: "keyboard",
		},
		&devicePathInfo{
			Path: "com.deepin.daemon.InputDevice.Mouse",
			Type: "mouse",
		},
		&devicePathInfo{
			Path: "com.deepin.daemon.InputDevice.TrackPoint",
			Type: "trackpoint",
		},
		&devicePathInfo{
			Path: "com.deepin.daemon.InputDevice.TouchPad",
			Type: "touchpad",
		},
	}

	m.kbd = getKeyboard()
	m.wacom = getWacom()
	m.tpad = getTouchpad()
	m.mouse = getMouse()
	m.trackPoint = getTrackPoint()

	return m
}

func (m *Manager) init() {
	m.kbd.init()
	m.kbd.handleGSettings()
	m.wacom.init()
	m.wacom.handleGSettings()
	m.tpad.init()
	m.tpad.handleGSettings()
	m.mouse.init()
	m.mouse.handleGSettings()
	m.trackPoint.init()
	m.trackPoint.handleGSettings()
}
