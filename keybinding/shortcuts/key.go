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

package shortcuts

import (
	"fmt"
	"strings"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/util/keysyms"
)

type Keycode x.Keycode
type Modifiers uint16

func (mods Modifiers) String() string {
	var keys []string
	if mods&x.ModMaskShift > 0 {
		keys = append(keys, "Shift")
	}
	if mods&x.ModMaskLock > 0 {
		keys = append(keys, "CapsLock")
	}
	if mods&x.ModMaskControl > 0 {
		keys = append(keys, "Control")
	}
	if mods&x.ModMask1 > 0 {
		keys = append(keys, "Alt")
	}
	if mods&x.ModMask2 > 0 {
		keys = append(keys, "NumLock")
	}
	if mods&x.ModMask3 > 0 {
		keys = append(keys, "Mod3")
	}
	if mods&x.ModMask4 > 0 {
		keys = append(keys, "Super")
	}
	if mods&x.ModMask5 > 0 {
		keys = append(keys, "Mod5")
	}
	return fmt.Sprintf("[%d|%s]", uint16(mods), strings.Join(keys, "-"))
}

type Key struct {
	Mods Modifiers
	Code Keycode
}

func (k Key) String() string {
	return fmt.Sprintf("Key<Mods=%s Code=%d>", k.Mods, k.Code)
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

func (k Key) ToAccel(keySymbols *keysyms.KeySymbols) ParsedAccel {
	str := keySymbols.LookupString(x.Keycode(k.Code), uint16(k.Mods))
	pa := ParsedAccel{
		Mods: k.Mods,
		Key:  str,
	}
	return pa.fix()
}

var IgnoreMods []uint16 = []uint16{
	0,
	x.ModMaskLock,              // Caps lock
	x.ModMask2,                 // Num lock
	x.ModMaskLock | x.ModMask2, // Caps and Num lock
}

func (k Key) Ungrab(conn *x.Conn) {
	rootWin := conn.GetDefaultScreen().Root
	for _, m := range IgnoreMods {
		x.UngrabKeyChecked(conn, x.Keycode(k.Code), rootWin, uint16(k.Mods)|m).Check(conn)
	}
}

func (k Key) Grab(conn *x.Conn) error {
	rootWin := conn.GetDefaultScreen().Root

	var err error
	for _, m := range IgnoreMods {
		err = x.GrabKeyChecked(conn, x.True, rootWin, uint16(k.Mods)|m, x.Keycode(k.Code),
			x.GrabModeAsync, x.GrabModeAsync).Check(conn)
		if err != nil {
			return err
		}
	}
	return nil
}
