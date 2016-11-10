/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package shortcuts

import (
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"strings"
)

type Keycode xproto.Keycode
type Modifiers uint16

func (mods Modifiers) String() string {
	var keys []string
	if mods&xproto.ModMaskShift > 0 {
		keys = append(keys, "Shift")
	}
	if mods&xproto.ModMaskLock > 0 {
		keys = append(keys, "CapsLock")
	}
	if mods&xproto.ModMaskControl > 0 {
		keys = append(keys, "Control")
	}
	if mods&xproto.ModMask1 > 0 {
		keys = append(keys, "Alt")
	}
	if mods&xproto.ModMask2 > 0 {
		keys = append(keys, "NumLock")
	}
	if mods&xproto.ModMask3 > 0 {
		keys = append(keys, "Mod3")
	}
	if mods&xproto.ModMask4 > 0 {
		keys = append(keys, "Super")
	}
	if mods&xproto.ModMask5 > 0 {
		keys = append(keys, "Mod5")
	}
	return fmt.Sprintf("[%d|%s]", uint16(mods), strings.Join(keys, "-"))
}

type Key struct {
	Mods Modifiers
	Code Keycode
}

func keysymToWeird(sym string) string {
	switch sym {
	case "-":
		return "minus"
	case "=":
		return "equal"
	case "\\":
		return "backslash"
	case "?":
		return "question"
	case "!":
		return "exclam"
	case "#":
		return "numbersign"
	case ";":
		return "semicolon"
	case "'":
		return "apostrophe"
	case "<":
		return "less"
	case ".":
		return "period"
	case "/":
		return "slash"
	case "(":
		return "parenleft"
	case "[":
		return "bracketleft"
	case ")":
		return "parenright"
	case "]":
		return "bracketright"
	case "\"":
		return "quotedbl"
	case " ":
		return "space"
	case "$":
		return "dollar"
	case "+":
		return "plus"
	case "*":
		return "asterisk"
	case "_":
		return "underscore"
	case "|":
		return "bar"
	case "`":
		return "grave"
	case "@":
		return "at"
	case "%":
		return "percent"
	case ">":
		return "greater"
	case "^":
		return "asciicircum"
	case "{":
		return "braceleft"
	case ":":
		return "colon"
	case ",":
		return "comma"
	case "~":
		return "asciitilde"
	case "&":
		return "ampersand"
	case "}":
		return "braceright"
	}

	return sym
}

func (k Key) ToAccel(xu *xgbutil.XUtil) ParsedAccel {
	str := keybind.LookupString(xu, uint16(k.Mods), xproto.Keycode(k.Code))
	pa := ParsedAccel{
		Mods: k.Mods,
		Key:  str,
	}
	return pa.fix()
}

func (k Key) Ungrab(xu *xgbutil.XUtil) {
	keybind.Ungrab(xu, xu.RootWin(), uint16(k.Mods), xproto.Keycode(k.Code))
}

type Keys []Key

func (keys Keys) Grab(xu *xgbutil.XUtil) error {
	rootWin := xu.RootWin()
	grabedKeys := make([]Key, len(keys))
	for _, key := range keys {
		// grab a key
		code := xproto.Keycode(key.Code)
		mods := uint16(key.Mods)

		err := keybind.GrabChecked(xu, rootWin, mods, code)
		logger.Debug("keybind.GrabChecked", key.Mods, code)
		if err != nil {
			for _, gk := range grabedKeys {
				keybind.Ungrab(xu, rootWin, uint16(gk.Mods), xproto.Keycode(gk.Code))
			}
			return err
		}
		grabedKeys = append(grabedKeys, key)
	}
	return nil
}

func (keys Keys) Ungrab(xu *xgbutil.XUtil) {
	for _, key := range keys {
		key.Ungrab(xu)
	}
}
