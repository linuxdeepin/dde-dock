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
	"fmt"
	"time"

	sys_network "github.com/linuxdeepin/go-dbus-factory/com.deepin.system.network"
	. "pkg.deepin.io/dde/daemon/keybinding/shortcuts"
	dbus "pkg.deepin.io/lib/dbus1"
)

func (m *Manager) shouldShowCapsLockOSD() bool {
	return m.gsKeyboard.GetBoolean(gsKeyShowCapsLockOSD)
}

func (m *Manager) initHandlers() {
	m.handlers = make([]KeyEventFunc, ActionTypeCount)
	logger.Debug("initHandlers", len(m.handlers))

	m.handlers[ActionTypeNonOp] = func(ev *KeyEvent) {
		logger.Debug("non-Op do nothing")
	}

	m.handlers[ActionTypeExecCmd] = func(ev *KeyEvent) {
		// prevent shortcuts such as switch window managers from being
		// triggered twice by mistake.
		const minExecCmdInterval = 600 * time.Millisecond
		type0 := ev.Shortcut.GetType()
		id := ev.Shortcut.GetId()
		if type0 == ShortcutTypeSystem && id == "wm-switcher" {
			now := time.Now()
			duration := now.Sub(m.lastExecCmdTime)
			logger.Debug("handle ActionTypeExecCmd duration:", duration)
			if 0 < duration && duration < minExecCmdInterval {
				logger.Debug("handle ActionTypeExecCmd ignore key event")
				return
			}
			m.lastExecCmdTime = now
		}

		action := ev.Shortcut.GetAction()
		arg, ok := action.Arg.(*ActionExecCmdArg)
		if !ok {
			logger.Warning(ErrTypeAssertionFail)
			return
		}

		go func() {
			err := m.execCmd(arg.Cmd, true)
			if err != nil {
				logger.Warning("execCmd error:", err)
			}
		}()
	}

	m.handlers[ActionTypeShowNumLockOSD] = func(ev *KeyEvent) {
		state, err := queryNumLockState(m.conn)
		if err != nil {
			logger.Warning(err)
			return
		}
		save := m.gsKeyboard.GetBoolean(gsKeySaveNumLockState)
		switch state {
		case NumLockOn:
			if save {
				m.NumLockState.Set(int32(NumLockOn))
			}
			showOSD("NumLockOn")
		case NumLockOff:
			if save {
				m.NumLockState.Set(int32(NumLockOff))
			}
			showOSD("NumLockOff")
		}
	}

	m.handlers[ActionTypeShowCapsLockOSD] = func(ev *KeyEvent) {
		if !m.shouldShowCapsLockOSD() {
			return
		}

		state, err := queryCapsLockState(m.conn)
		if err != nil {
			logger.Warning(err)
			return
		}

		switch state {
		case CapsLockOff:
			showOSD("CapsLockOff")
		case CapsLockOn:
			showOSD("CapsLockOn")
		}
	}

	m.handlers[ActionTypeOpenMimeType] = func(ev *KeyEvent) {
		action := ev.Shortcut.GetAction()
		mimeType, ok := action.Arg.(string)
		if !ok {
			logger.Warning(ErrTypeAssertionFail)
			return
		}

		go func() {
			err := m.execCmd(queryCommandByMime(mimeType), true)
			if err != nil {
				logger.Warning("execCmd error:", err)
			}
		}()
	}

	m.handlers[ActionTypeDesktopFile] = func(ev *KeyEvent) {
		action := ev.Shortcut.GetAction()

		go func() {
			err := m.runDesktopFile(action.Arg.(string))
			if err != nil {
				logger.Warning("runDesktopFile error:", err)
			}
		}()
	}

	m.handlers[ActionTypeAudioCtrl] = buildHandlerFromController(m.audioController)
	m.handlers[ActionTypeMediaPlayerCtrl] = buildHandlerFromController(m.mediaPlayerController)
	m.handlers[ActionTypeDisplayCtrl] = buildHandlerFromController(m.displayController)
	m.handlers[ActionTypeKbdLightCtrl] = buildHandlerFromController(m.kbdLightController)
	m.handlers[ActionTypeTouchpadCtrl] = buildHandlerFromController(m.touchPadController)
	m.handlers[ActionTypeToggleWireless] = func(ev *KeyEvent) {
		if m.gsMediaKey.GetBoolean(gsKeyUpperLayerWLAN) {
			sysBus, err := dbus.SystemBus()
			if err != nil {
				logger.Warning(err)
				return
			}
			sysNetwork := sys_network.NewNetwork(sysBus)
			enabled, err := sysNetwork.ToggleWirelessEnabled(0)
			if err != nil {
				logger.Warning("failed to toggle wireless enabled:", err)
				return
			}
			if enabled {
				showOSD("WLANOn")
			} else {
				showOSD("WLANOff")
			}

		} else {
			state, err := getRfkillWlanState()
			if err != nil {
				logger.Warning(err)
				return
			}
			if state == 0 {
				showOSD("WLANOff")
			} else {
				showOSD("WLANOn")
			}
		}
	}

	m.handlers[ActionTypeSystemShutdown] = func(ev *KeyEvent) {
		cmd := getPowerButtonPressedExec()

		go func() {
			err := m.execCmd(cmd, false)
			if err != nil {
				logger.Warning("execCmd error:", err)
			}
		}()
	}

	m.handlers[ActionTypeSystemSuspend] = func(ev *KeyEvent) {
		systemSuspend()
	}

	m.handlers[ActionTypeSystemLogOff] = func(ev *KeyEvent) {
		systemLogout()
	}

	m.handlers[ActionTypeSystemAway] = func(ev *KeyEvent) {
		systemAway()
	}

	// handle Switch Kbd Layout
	m.handlers[ActionTypeSwitchKbdLayout] = func(ev *KeyEvent) {
		logger.Debug("Switch Kbd Layout state", m.switchKbdLayoutState)
		flags := m.ShortcutSwitchLayout.Get()
		action := ev.Shortcut.GetAction()
		arg, ok := action.Arg.(uint32)
		if !ok {
			logger.Warning(ErrTypeAssertionFail)
			return
		}

		if arg&flags == 0 {
			return
		}

		switch m.switchKbdLayoutState {
		case SKLStateNone:
			m.switchKbdLayoutState = SKLStateWait
			go m.sklWait()

		case SKLStateWait:
			m.switchKbdLayoutState = SKLStateOSDShown
			m.terminateSKLWait()
			showOSD("SwitchLayout")

		case SKLStateOSDShown:
			showOSD("SwitchLayout")
		}
	}

	m.handlers[ActionTypeShowControlCenter] = func(ev *KeyEvent) {
		err := m.execCmd("dbus-send --session --dest=com.deepin.dde.ControlCenter  --print-reply /com/deepin/dde/ControlCenter com.deepin.dde.ControlCenter.Show",
			false)
		if err != nil {
			logger.Warning("failed to show control center:", err)
		}
	}

	m.shortcutManager.SetAllModKeysReleasedCallback(func() {
		switch m.switchKbdLayoutState {
		case SKLStateWait:
			showOSD("DirectSwitchLayout")
			m.terminateSKLWait()
		case SKLStateOSDShown:
			showOSD("SwitchLayoutDone")
		case SKLStateNone:
			return
		}
		m.switchKbdLayoutState = SKLStateNone
	})
}

func (m *Manager) sklWait() {
	defer func() {
		logger.Debug("sklWait end")
		m.sklWaitQuit = nil
	}()

	m.sklWaitQuit = make(chan int)
	timer := time.NewTimer(350 * time.Millisecond)
	select {
	case <-m.sklWaitQuit:
		return
	case _, ok := <-timer.C:
		if !ok {
			logger.Error("Invalid ticker event")
			return
		}

		logger.Debug("timer fired")
		if m.switchKbdLayoutState == SKLStateWait {
			m.switchKbdLayoutState = SKLStateOSDShown
			showOSD("SwitchLayout")
		}
	}
}

func (m *Manager) terminateSKLWait() {
	if m.sklWaitQuit != nil {
		close(m.sklWaitQuit)
	}
}

type Controller interface {
	ExecCmd(cmd ActionCmd) error
	Name() string
}

func buildHandlerFromController(c Controller) KeyEventFunc {
	return func(ev *KeyEvent) {
		if c == nil {
			logger.Warning("controller is nil")
			return
		}
		name := c.Name()

		action := ev.Shortcut.GetAction()
		cmd, ok := action.Arg.(ActionCmd)
		if !ok {
			logger.Warning(ErrTypeAssertionFail)
			return
		}
		logger.Debugf("%v Controller exec cmd %v", name, cmd)
		if err := c.ExecCmd(cmd); err != nil {
			logger.Warning(name, "Controller exec cmd err:", err)
		}
	}
}

type ErrInvalidActionCmd struct {
	Cmd ActionCmd
}

func (err ErrInvalidActionCmd) Error() string {
	return fmt.Sprintf("invalid action cmd %v", err.Cmd)
}

type ErrIsNil struct {
	Name string
}

func (err ErrIsNil) Error() string {
	return fmt.Sprintf("%s is nil", err.Name)
}
