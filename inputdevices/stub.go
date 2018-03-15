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

package inputdevices

const (
	dbusServiceName = "com.deepin.daemon.InputDevices"
	dbusPath        = "/com/deepin/daemon/InputDevices"
	dbusInterface   = dbusServiceName

	kbdDBusPath      = "/com/deepin/daemon/InputDevice/Keyboard"
	kbdDBusInterface = "com.deepin.daemon.InputDevice.Keyboard"

	mouseDBusPath           = "/com/deepin/daemon/InputDevice/Mouse"
	mouseDBusInterface      = "com.deepin.daemon.InputDevice.Mouse"
	trackPointDBusInterface = "com.deepin.daemon.InputDevice.TrackPoint"

	touchPadDBusPath      = "/com/deepin/daemon/InputDevice/TouchPad"
	touchPadDBusInterface = "com.deepin.daemon.InputDevice.TouchPad"

	wacomDBusPath      = "/com/deepin/daemon/InputDevice/Wacom"
	wacomDBusInterface = "com.deepin.daemon.InputDevice.Wacom"
)

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}

func (*Keyboard) GetInterfaceName() string {
	return kbdDBusInterface
}

func (*Mouse) GetInterfaceName() string {
	return mouseDBusInterface
}

func (*TrackPoint) GetInterfaceName() string {
	return trackPointDBusInterface
}

func (*Touchpad) GetInterfaceName() string {
	return touchPadDBusInterface
}

func (*Wacom) GetInterfaceName() string {
	return wacomDBusInterface
}
