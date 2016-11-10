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
	"dbus/com/deepin/daemon/helper/backlight"
	. "pkg.deepin.io/dde/daemon/keybinding/shortcuts"
)

type KbdLightController struct {
	blDaemon *backlight.Backlight
}

func NewKbdLightController(blDaemon *backlight.Backlight) *KbdLightController {
	return &KbdLightController{
		blDaemon: blDaemon,
	}
}

func (c *KbdLightController) Name() string {
	return "Kbd Light"
}

func (c *KbdLightController) ExecCmd(cmd ActionCmd) error {
	switch cmd {
	case KbdLightBrightnessUp:
		return c.changeBrightness(true)
	case KbdLightBrightnessDown:
		return c.changeBrightness(false)
	case KbdLightToggle:
		return c.toggle()
	default:
		return ErrInvalidActionCmd{cmd}
	}
}

// Let the keyboard light brightness switch between 0 to max
func (c *KbdLightController) toggle() error {
	if c.blDaemon == nil {
		return ErrIsNil{"KbdLightController.blDaemon"}
	}

	value, err := c.blDaemon.GetKbdBrightness()
	if err != nil {
		return err
	}

	maxValue, err := c.blDaemon.GetKbdMaxBrightness()
	if err != nil {
		return err
	}
	if value == 0 {
		value = maxValue
	} else {
		value = 0
	}
	logger.Debug("[changeKbdBrightness] will set kbd backlight to:", value)
	return c.blDaemon.SetKbdBrightness(value)
}

var kbdBacklightStep int32 = 0

func (c *KbdLightController) changeBrightness(raised bool) error {
	if c.blDaemon == nil {
		return ErrIsNil{"KbdLightController.blDaemon"}
	}

	value, err := c.blDaemon.GetKbdBrightness()
	if err != nil {
		return err
	}

	maxValue, err := c.blDaemon.GetKbdMaxBrightness()
	if err != nil {
		return err
	}

	// step: (max < 10?1:max/10)
	if kbdBacklightStep == 0 {
		tmp := maxValue / 10
		if tmp == 0 {
			tmp = 1
		}
		// 4舍5入
		if float64(maxValue)/10 < float64(tmp)+0.5 {
			kbdBacklightStep = tmp
		} else {
			kbdBacklightStep = tmp + 1
		}
	}
	logger.Debug("[changeKbdBrightness] pld kbd backlight info:", value, maxValue, kbdBacklightStep)
	if raised {
		value += kbdBacklightStep
	} else {
		value -= kbdBacklightStep
	}

	if value < 0 {
		value = 0
	} else if value > maxValue {
		value = maxValue
	}

	logger.Debug("[changeKbdBrightness] will set kbd backlight to:", value)
	return c.blDaemon.SetKbdBrightness(value)
}
