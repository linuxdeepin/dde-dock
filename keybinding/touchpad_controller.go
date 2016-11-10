/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package keybinding

import (
	"dbus/com/deepin/daemon/inputdevices"
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

func (c *TouchpadController) enable(val bool) {
	// check touchpad exist?
	exist := c.touchpad.Exist.Get()
	if !exist {
		return
	}

	c.touchpad.TPadEnable.Set(val)

	osd := "TouchpadOn"
	if !val {
		osd = "TouchpadOff"
	}
	showOSD(osd)
}

func (c *TouchpadController) toggle() {
	// check touchpad exist?
	exist := c.touchpad.Exist.Get()
	if !exist {
		return
	}

	c.touchpad.TPadEnable.Set(!c.touchpad.TPadEnable.Get())

	showOSD("TouchpadToggle")
}
