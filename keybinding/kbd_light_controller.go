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
	"dbus/com/deepin/daemon/helper/backlight"
	"errors"
	. "pkg.deepin.io/dde/daemon/keybinding/shortcuts"
	commonbl "pkg.deepin.io/lib/backlight/common"
	kbdbl "pkg.deepin.io/lib/backlight/keyboard"
)

const backlightTypeKeyboard = 2

type KbdLightController struct {
	backlightHelper *backlight.Backlight
}

func NewKbdLightController(backlightHelper *backlight.Backlight) *KbdLightController {
	return &KbdLightController{
		backlightHelper: backlightHelper,
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

func getKbdBlController() (*commonbl.Controller, error) {
	controllers, err := kbdbl.List()
	if err != nil {
		return nil, err
	}
	if len(controllers) > 0 {
		return controllers[0], nil
	}
	return nil, errors.New("not found keyboard backlight controller")
}

// Let the keyboard light brightness switch between 0 to max
func (c *KbdLightController) toggle() error {
	if c.backlightHelper == nil {
		return ErrIsNil{"KbdLightController.backlightHelper"}
	}

	controller, err := getKbdBlController()
	if err != nil {
		return err
	}
	value, err := controller.GetBrightness()
	if err != nil {
		return err
	}

	if value == 0 {
		value = controller.MaxBrightness
	} else {
		value = 0
	}
	logger.Debug("[KbdLightController.toggle] will set kbd backlight to:", value)
	return c.backlightHelper.SetBrightness(backlightTypeKeyboard, controller.Name, int32(value))
}

var kbdBacklightStep int = 0

func (c *KbdLightController) changeBrightness(raised bool) error {
	if c.backlightHelper == nil {
		return ErrIsNil{"KbdLightController.backlightHelper"}
	}

	controller, err := getKbdBlController()
	if err != nil {
		return err
	}
	value, err := controller.GetBrightness()
	if err != nil {
		return err
	}

	maxValue := controller.MaxBrightness

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
	logger.Debug("[KbdLightController.changeBrightness] kbd backlight info:", value, maxValue, kbdBacklightStep)
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

	logger.Debug("[KbdLightController.changeBrightness] will set kbd backlight to:", value)
	return c.backlightHelper.SetBrightness(backlightTypeKeyboard, controller.Name, int32(value))
}
