/**
 * Copyright (c) 2011 ~ 2015 Deepin, Inc.
 *               2013 ~ 2015 jouyouyun
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

package keybinding

type Mediakey struct {
	AudioMute        func(bool)
	AudioUp          func(bool)
	AudioDown        func(bool)
	BrightnessUp     func(bool)
	BrightnessDown   func(bool)
	CapsLockOn       func(bool)
	CapsLockOff      func(bool)
	NumLockOn        func(bool)
	NumLockOff       func(bool)
	SwitchMonitors   func(bool)
	TouchpadOn       func(bool)
	TouchpadOff      func(bool)
	TouchpadToggle   func(bool)
	PowerOff         func(bool)
	PowerSleep       func(bool)
	SwitchLayout     func(bool)
	AudioPlay        func(bool)
	AudioPause       func(bool)
	AudioStop        func(bool)
	AudioPrevious    func(bool)
	AudioNext        func(bool)
	AudioRewind      func(bool)
	AudioForward     func(bool)
	AudioRepeat      func(bool)
	LaunchEmail      func(bool)
	LaunchBrowser    func(bool)
	LaunchCalculator func(bool)
}
