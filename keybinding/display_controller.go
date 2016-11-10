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
	"errors"
	"fmt"
	. "pkg.deepin.io/dde/daemon/keybinding/shortcuts"
)

type DisplayController struct {
	disp     *display.Display
	blDaemon *backlight.Backlight
}

func NewDisplayController(blDaemon *backlight.Backlight) (*DisplayController, error) {
	c := new(DisplayController)
	c.blDaemon = blDaemon
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

func (c *DisplayController) getBacklightInfo() (int32, int32, error) {
	blDaemon := c.blDaemon
	if blDaemon == nil {
		return 0, 0, ErrIsNil{"DisplayController.blDaemon"}
	}

	list, _ := blDaemon.ListSysPath()
	if len(list) == 0 {
		return 0, 0, errors.New("no backlight found")
	}

	cur, err := blDaemon.GetBrightness(list[0])
	if err != nil {
		return 0, 0, err
	}

	max, err := blDaemon.GetMaxBrightness(list[0])
	if err != nil {
		return 0, 0, err
	}
	return cur, max, nil
}

func (c *DisplayController) changeBrightness(raised bool) error {
	if c.disp == nil {
		return ErrIsNil{"DisplayController.disp"}
	}

	outputs, err := c.disp.ListOutputNames()
	if err != nil {
		return err
	}

	brightnessMap := c.disp.Brightness.Get()
	var step float64 = 0.05
	if !raised {
		step = -step
	}

	var backlightValid bool = true
	cur, max, err := c.getBacklightInfo()
	if err != nil {
		backlightValid = false
	}
	for _, output := range outputs {
		v, ok := brightnessMap[output]
		if !ok {
			v = 1.0
		}
		var discrete float64
		supported, _ := c.disp.SupportedBacklight(output)
		if backlightValid && supported {
			// TODO: Some drivers will also set the brightness when the brightness up/down key is pressed
			v = float64(cur) / float64(max)
		}
		discrete = v + step
		if discrete > 1.0 {
			discrete = 1
		}
		if discrete < 0.02 {
			discrete = 0.02
		}
		logger.Debug("[changeBrightness] will set to:", output, discrete)
		if err := c.disp.SetBrightness(output, discrete); err != nil {
			logger.Warningf("changeBrightness set failed, output: %v, discrete: %v, err: %v", output, discrete, err)
		}
	}
	if err := c.disp.SaveBrightness(); err != nil {
		return fmt.Errorf("saveBrightness failed", err)
	}

	var osd = "BrightnessUp"
	if !raised {
		osd = "BrightnessDown"
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
