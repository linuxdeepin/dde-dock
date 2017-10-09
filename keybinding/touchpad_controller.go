/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
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
	"dbus/com/deepin/daemon/inputdevices"

	ddbus "pkg.deepin.io/dde/daemon/dbus"
	. "pkg.deepin.io/dde/daemon/keybinding/shortcuts"
)

type TouchpadController struct {
	touchpad *inputdevices.TouchPad
}

func NewTouchpadController() (*TouchpadController, error) {
	c := new(TouchpadController)
	var err error
	c.touchpad, err = inputdevices.NewTouchPad("com.deepin.daemon.InputDevices", "/com/deepin/daemon/InputDevice/TouchPad")
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *TouchpadController) Destroy() {
	if c.touchpad != nil {
		inputdevices.DestroyTouchPad(c.touchpad)
		c.touchpad = nil
	}
}

func (*TouchpadController) Name() string {
	return "Touchpad"
}

func (c *TouchpadController) ExecCmd(cmd ActionCmd) error {
	switch cmd {
	case TouchpadToggle:
		c.toggle()
	case TouchpadOn:
		c.enable(true)
	case TouchpadOff:
		c.enable(false)
	default:
		return ErrInvalidActionCmd{cmd}
	}
	return nil
}

func (c *TouchpadController) enable(val bool) error {
	if c.touchpad == nil || !ddbus.IsSessionBusActivated(c.touchpad.DestName) {
		return ErrIsNil{"TouchpadController.touchpad"}
	}
	// check touchpad exist?
	exist := c.touchpad.Exist.Get()
	if !exist {
		return nil
	}

	c.touchpad.TPadEnable.Set(val)

	osd := "TouchpadOn"
	if !val {
		osd = "TouchpadOff"
	}
	showOSD(osd)
	return nil
}

func (c *TouchpadController) toggle() error {
	if c.touchpad == nil || !ddbus.IsSessionBusActivated(c.touchpad.DestName) {
		return ErrIsNil{"TouchpadController.touchpad"}
	}
	// check touchpad exist?
	exist := c.touchpad.Exist.Get()
	if !exist {
		return nil
	}

	c.touchpad.TPadEnable.Set(!c.touchpad.TPadEnable.Get())

	showOSD("TouchpadToggle")
	return nil
}
