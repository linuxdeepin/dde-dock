/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package keybinding

type Mediakey struct {
	AudioMute         func(bool)
	AudioUp           func(bool)
	AudioDown         func(bool)
	BrightnessUp      func(bool)
	BrightnessDown    func(bool)
	KbdBrightnessUp   func(bool)
	KbdBrightnessDown func(bool)
	CapsLockOn        func(bool)
	CapsLockOff       func(bool)
	NumLockOn         func(bool)
	NumLockOff        func(bool)
	SwitchMonitors    func(bool)
	TouchpadOn        func(bool)
	TouchpadOff       func(bool)
	TouchpadToggle    func(bool)
	PowerOff          func(bool)
	PowerSleep        func(bool)
	PowerSuspend      func(bool)
	SwitchLayout      func(bool)
	AudioPlay         func(bool)
	AudioPause        func(bool)
	AudioStop         func(bool)
	AudioPrevious     func(bool)
	AudioNext         func(bool)
	AudioRewind       func(bool)
	AudioForward      func(bool)
	AudioRepeat       func(bool)
	AudioMedia        func(bool)
	LaunchEmail       func(bool)
	LaunchBrowser     func(bool)
	LaunchCalculator  func(bool)
	Eject             func(bool)
}
