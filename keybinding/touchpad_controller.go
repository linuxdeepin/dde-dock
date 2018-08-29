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
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.inputdevices"
	. "pkg.deepin.io/dde/daemon/keybinding/shortcuts"
	"pkg.deepin.io/lib/dbus1"
)

type TouchPadController struct {
	touchPad *inputdevices.TouchPad
}

func NewTouchPadController(sessionConn *dbus.Conn) *TouchPadController {
	c := new(TouchPadController)
	c.touchPad = inputdevices.NewTouchPad(sessionConn)
	return c
}

func (*TouchPadController) Name() string {
	return "TouchPad"
}

func (c *TouchPadController) ExecCmd(cmd ActionCmd) error {
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

func (c *TouchPadController) enable(val bool) error {
	exist, err := c.touchPad.Exist().Get(0)
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}

	err = c.touchPad.TPadEnable().Set(0, val)
	if err != nil {
		return err
	}

	osd := "TouchpadOn"
	if !val {
		osd = "TouchpadOff"
	}
	showOSD(osd)
	return nil
}

func (c *TouchPadController) toggle() error {
	// check touchpad exist?
	exist, err := c.touchPad.Exist().Get(0)
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}

	if globalConfig.HandleTouchPadToggle {
		enabled, err := c.touchPad.TPadEnable().Get(0)
		if err != nil {
			return err
		}
		err = c.touchPad.TPadEnable().Set(0, !enabled)
		if err != nil {
			return err
		}
	}

	showOSD("TouchpadToggle")
	return nil
}
