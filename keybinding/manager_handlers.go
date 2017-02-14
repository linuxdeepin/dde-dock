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
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	. "pkg.deepin.io/dde/daemon/keybinding/shortcuts"
)

func (m *Manager) initHandlers() {
	m.handlers = make([]KeyEventFunc, ActionTypeCount)
	logger.Debug("initHandlers", len(m.handlers))

	m.handlers[ActionTypeNonOp] = func(ev *KeyEvent) {
		logger.Debug("non-Op do nothing")
	}

	m.handlers[ActionTypeExecCmd] = func(ev *KeyEvent) {
		action := ev.Shortcut.GetAction()
		arg, ok := action.Arg.(*ActionExecCmdArg)
		if !ok {
			logger.Warning(ErrTypeAssertionFail)
			return
		}

		go func() {
			err := execCmd(arg.Cmd)
			if err != nil {
				logger.Warning("execCmd error:", err)
			}
		}()
	}

	m.handlers[ActionTypeShowNumLockOSD] = func(ev *KeyEvent) {
		save := m.keyboardSetting.GetBoolean(gsKeySaveNumLockState)
		if ev.Mods&xproto.ModMask2 > 0 {
			if save {
				m.NumLockState.Set(int32(NumLockOff))
			}
			showOSD("NumLockOff")
		} else {
			if save {
				m.NumLockState.Set(int32(NumLockOn))
			}
			showOSD("NumLockOn")
		}

	}

	m.handlers[ActionTypeShowCapsLockOSD] = func(ev *KeyEvent) {
		if !canShowCapsOSD() {
			return
		}
		if ev.Mods&xproto.ModMaskLock > 0 {
			// lowercase alphabet
			showOSD("CapsLockOff")
		} else {
			// uppercase alphabet
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
			err := execCmd(queryCommandByMime(mimeType))
			if err != nil {
				logger.Warning("execCmd error:", err)
			}
		}()
	}

	m.handlers[ActionTypeAudioCtrl] = buildHandlerFromController(m.audioController)
	m.handlers[ActionTypeMediaPlayerCtrl] = buildHandlerFromController(m.mediaPlayerController)
	m.handlers[ActionTypeDisplayCtrl] = buildHandlerFromController(m.displayController)
	m.handlers[ActionTypeKbdLightCtrl] = buildHandlerFromController(m.kbdLightController)
	m.handlers[ActionTypeTouchpadCtrl] = buildHandlerFromController(m.touchpadController)

	m.handlers[ActionTypeSystemShutdown] = func(ev *KeyEvent) {
		systemShutdown()
	}

	m.handlers[ActionTypeSystemSuspend] = func(ev *KeyEvent) {
		systemSuspend()
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
