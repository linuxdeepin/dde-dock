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
	"dbus/com/deepin/daemon/display"
	"dbus/com/deepin/daemon/helper/backlight"
	ddbus "pkg.deepin.io/dde/daemon/dbus"
	. "pkg.deepin.io/dde/daemon/keybinding/shortcuts"
)

type DisplayController struct {
	disp            *display.Display
	backlightHelper *backlight.Backlight
}

func NewDisplayController(backlightHelper *backlight.Backlight) (*DisplayController, error) {
	c := new(DisplayController)
	c.backlightHelper = backlightHelper
	var err error
	c.disp, err = display.NewDisplay("com.deepin.daemon.Display",
		"/com/deepin/daemon/Display")
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (*DisplayController) Name() string {
	return "Display"
}

func (c *DisplayController) ExecCmd(cmd ActionCmd) error {
	switch cmd {
	case DisplayModeSwitch:
		showOSD("SwitchMonitors")
		return nil

	case MonitorBrightnessUp:
		return c.changeBrightness(true)

	case MonitorBrightnessDown:
		return c.changeBrightness(false)

	default:
		return ErrInvalidActionCmd{cmd}
	}
}

func (c *DisplayController) changeBrightness(raised bool) error {
	if c.disp == nil || !ddbus.IsSessionBusActivated(c.disp.DestName) {
		return ErrIsNil{"DisplayController.disp"}
	}

	var osd = "BrightnessUp"
	if !raised {
		osd = "BrightnessDown"
	}

	err := c.disp.ChangeBrightness(raised)
	if err != nil {
		return err
	}

	showOSD(osd)
	return nil
}

func (c *DisplayController) Destroy() {
	if c.disp != nil {
		display.DestroyDisplay(c.disp)
		c.disp = nil
	}
}
