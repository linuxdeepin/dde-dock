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

package wacom

import (
	"fmt"
	"os/exec"
)

/**
 * xsetwacom set deviceName Rotate half/none
 * default: none
 */
func (w *Wacom) rotationAngle(leftHanded bool) {
	for _, info := range w.DeviceList {
		cmdline := "xsetwacom set \"" + info.Name + "\" Rotate "
		if leftHanded {
			cmdline += "half"
		} else {
			cmdline += "none"
		}

		err := doCommand(cmdline)
		if err != nil {
			w.debugInfo("Exec '%s' failed: %v", cmdline, err)
		}
	}
}

/**
 * xsetwacom set deviceName mode Relative/Absolute
 * default: Absolute for  stylus,  eraser  and  tablet  PC  touch;
 *          Relative for cursor and tablet touch.
 */
func (w *Wacom) cursorMode(cursorMode bool) {
	for _, info := range w.DeviceList {
		cmdline := "xsetwacom set \"" + info.Name + "\" mode "
		if cursorMode {
			cmdline += "Relative"
		} else {
			cmdline += "Absolute"
		}

		err := doCommand(cmdline)
		if err != nil {
			w.debugInfo("Exec '%s' failed: %v", cmdline, err)
		}
	}
}

/**
 * xsetwacom set deviceName Button 3 3/"KP_Page_Up"
 * default: 3
 */
func (w *Wacom) keyUpAction(action string) {
	value, ok := descActionMap[action]
	if !ok {
		return
	}

	for _, info := range w.DeviceList {
		cmdline := "xsetwacom set \"" + info.Name + "\" Button 3 " + value
		err := doCommand(cmdline)
		if err != nil {
			w.debugInfo("Exec '%s' failed: %v", cmdline, err)
		}
	}
}

/**
 * xsetwacom set deviceName Button 2 2/"KP_Page_Down"
 * default: 2
 */
func (w *Wacom) keyDownAction(action string) {
	value, ok := descActionMap[action]
	if !ok {
		return
	}

	for _, info := range w.DeviceList {
		cmdline := "xsetwacom set \"" + info.Name + "\" Button 2 " + value
		err := doCommand(cmdline)
		if err != nil {
			w.debugInfo("Exec '%s' failed: %v", cmdline, err)
		}
	}
}

/**
 * xsetwacom set deviceName Suppress 0-100
 * default: 2
 */
func (w *Wacom) doubleDelta(delta uint32) {
	for _, info := range w.DeviceList {
		cmdline := fmt.Sprintf("xsetwacom set \"%s\" Suppress %v", info.Name, delta)
		err := doCommand(cmdline)
		if err != nil {
			w.debugInfo("Exec '%s' failed: %v", cmdline, err)
		}
	}
}

/**
 * xsetwacom set deviceName Threshold 0-2047
 * default: 27
 */
func (w *Wacom) pressureSensitive(pressure uint32) {
	if pressure > 2047 {
		return
	}

	for _, info := range w.DeviceList {
		cmdline := fmt.Sprintf("xsetwacom set \"%s\" Threshold %v", info.Name, pressure)
		err := doCommand(cmdline)
		if err != nil {
			w.debugInfo("Exec '%s' failed: %v", cmdline, err)
		}
	}
}

func doCommand(cmdline string) (err error) {
	if len(cmdline) < 1 {
		return fmt.Errorf("doCommand args is nil")
	}

	err = exec.Command("/bin/sh", "-c", cmdline).Run()

	return
}
