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
	"bytes"
	"errors"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"pkg.deepin.io/dde/daemon/keybinding/keybind"
	"strconv"
	"strings"
)

type Accel struct {
	Parsed    ParsedAccel
	GrabedKey Key
	Shortcut  Shortcut
}

// ParsedAccel
// field Mods ignore mod2(Num_Lock) and lock(Caps_Lock)
type ParsedAccel struct {
	Mods Modifiers
	Key  string
}

func (ap1 ParsedAccel) Equal(xu *xgbutil.XUtil, ap2 ParsedAccel) bool {
	logger.Debug(ap1, " equal? ", ap2)
	if ap1.Mods != ap2.Mods {
		logger.Debug("Mods no equal, return false")
		return false
	}
	// ap1.Mods == ap2.Mods
	if ap1.Key == ap2.Key {
		logger.Debug("Key equal, return true")
		return true
	}

	// ap1.Key != ap2.Key
	codes1, err := keybind.StrToKeycodes(xu, ap1.Key)
	if err != nil {
		return false
	}
	codes2, err := keybind.StrToKeycodes(xu, ap2.Key)
	if err != nil {
		return false
	}

	keycodesEq := isKeycodesEqual(codes1, codes2)
	logger.Debug("keycodesEq:", keycodesEq)
	return keycodesEq
}

func isKeycodesEqual(list1, list2 []xproto.Keycode) bool {
	logger.Debug("isKeycodesEqual:", list1, list2)
	l1 := len(list1)
	l2 := len(list2)
	if l1 != l2 {
		return false
	}

	for i, v := range list1 {
		if v != list2[i] {
			return false
		}
	}

	return true
}

func (pa ParsedAccel) MarshalJSON() ([]byte, error) {
	str := pa.String()
	quoted := strconv.Quote(str)
	return []byte(quoted), nil
}

func GetKeyFirstCode(xu *xgbutil.XUtil, str string) (xproto.Keycode, error) {
	codes, err := keybind.StrToKeycodes(xu, str)
	if err != nil {
		return 0, err
	}
	var code xproto.Keycode
	for _, kc := range codes {
		if kc != 0 {
			code = kc
			break
		}
	}
	if code == 0 {
		return 0, errors.New("not found keycode")
	}
	logger.Debugf("GetKeyFirstCode str %q codes: %v code: %d", str, codes, code)
	return code, nil
}

func (pa ParsedAccel) QueryKey(xu *xgbutil.XUtil) (Key, error) {
	code, err := GetKeyFirstCode(xu, pa.Key)
	if err != nil {
		return Key{}, err
	}
	return Key{
		Mods: pa.Mods,
		Code: Keycode(code),
	}, nil
}

func splitStandardAccel(accel string) ([]string, error) {
	if accel == "" {
		return nil, nil
	}

	var keys []string
	reader := strings.NewReader(accel)
	for {
		ch, err := reader.ReadByte()
		if err != nil {
			// eof
			break
		}

		switch ch {
		case '<':
			// read byte is not '>' , fill buf key
			// read byte is '>' push key.String() to keys
			var key bytes.Buffer
		Loop:
			for {
				ch, err := reader.ReadByte()
				if err != nil {
					// eof
					return nil, errors.New("> not found")
				}
				switch ch {
				case '>':
					break Loop
				case '<':
					return nil, errors.New("unexpect < found")
				default:
					key.WriteByte(ch)
				}
			}
			if key.Len() > 0 {
				keys = append(keys, key.String())
			} else {
				return nil, errors.New("empty modifier found")
			}
		default:
			reader.UnreadByte()
			var key bytes.Buffer
			// read rest bytes
			for {
				ch, err := reader.ReadByte()
				if err != nil {
					break
				}
				switch ch {
				case '<', '>':
					return nil, errors.New("unexpect < or > found")
				default:
					key.WriteByte(ch)
				}
			}
			keys = append(keys, key.String())
		}
	}
	return keys, nil
}

// <Super>L mods (mod4) key L
// <Super>% mods (mod4, shift) key %
// <Control><Alt>T mods (control,mod1) key T
// <Control><shift>T mods(control,shift) key T
// <Control>> mods(control) key >
// <Control>< invalid
// Super< invalid
// <Super> mods() key Super
// Print mods() key Print
// <Control>Print mods(Control) key Print
// check ParsedAccel.Key valid later
func ParseStandardAccel(accel string) (ParsedAccel, error) {
	parts, err := splitStandardAccel(accel)
	if err != nil {
		return ParsedAccel{}, err
	}
	switch len(parts) {
	case 0:
		return ParsedAccel{}, errors.New("empty parts")
	case 1:
		return ParsedAccel{Key: parts[0]}, nil
	}

	key := parts[len(parts)-1]
	// check key valid

	var mods Modifiers
	for _, part := range parts[:len(parts)-1] {
		switch strings.ToLower(part) {
		case "shift":
			mods |= xproto.ModMaskShift
		case "control":
			mods |= xproto.ModMaskControl
		case "alt":
			mods |= xproto.ModMask1
		case "super":
			mods |= xproto.ModMask4
		default:
			return ParsedAccel{}, errors.New("unexpect mod " + part)
		}
	}

	return ParsedAccel{
		Mods: mods,
		Key:  key,
	}, nil
}

func ParseStandardAccels(accelStrv []string) []ParsedAccel {
	parsedAccels := make([]ParsedAccel, 0, len(accelStrv))
	for _, accel := range accelStrv {
		parsed, err := ParseStandardAccel(accel)
		if err == nil {
			parsedAccels = append(parsedAccels, parsed)
		}
		// TODO else warning
	}
	return parsedAccels
}

// get standard key sequece
func (pa ParsedAccel) String() string {
	var keys []string
	mods := pa.Mods
	if mods&xproto.ModMaskShift > 0 {
		keys = append(keys, "<Shift>")
	}
	if mods&xproto.ModMaskControl > 0 {
		keys = append(keys, "<Control>")
	}
	if mods&xproto.ModMask1 > 0 {
		keys = append(keys, "<Alt>")
	}
	if mods&xproto.ModMask4 > 0 {
		keys = append(keys, "<Super>")
	}

	keys = append(keys, pa.Key)
	return strings.Join(keys, "")
}

func isGoodSingleKey(key string) bool {
	// single key
	switch key {
	case "f1", "f2", "f3", "f4", "f5", "f6",
		"f7", "f8", "f9", "f10", "f11", "f12",
		"print", "backspace", "delete", "super_l", "super_r":
		return true
	default:
		if strings.HasPrefix(key, "xf86") {
			return true
		}
		return false
	}
}

func (pa ParsedAccel) IsGood() bool {
	keyLower := strings.ToLower(pa.Key)
	if pa.Mods == 0 {
		return isGoodSingleKey(keyLower)
	}

	// pa.Mod > 0
	if pa.Mods&^xproto.ModMaskShift == 0 {
		// mods is <Shift>
		// TODO
		return isGoodSingleKey(keyLower)
	}

	switch keyLower {
	case "shift_r", "shift_l",
		"alt_r", "alt_l",
		"meta_r", "meta_l",
		"super_r", "super_l",
		"hyper_r", "hyper_l",
		"control_r", "control_l":
		return false
	}

	return true
}

// char is a-z
func isLowerAlpha(char byte) bool {
	if 'a' <= char && char <= 'z' {
		return true
	}
	return false
}

func (pa ParsedAccel) fix() ParsedAccel {
	logger.Debug("before fix", pa)
	var keyLower = strings.ToLower(pa.Key)
	var key string
	switch keyLower {
	case "l1":
		key = "F11"
	case "l2":
		key = "F12"
	case "kb_tab", "iso_left_tab":
		key = "Tab"
	default:
		key = keysymToWeird(pa.Key)
	}

	if len(key) == 1 && isLowerAlpha(key[0]) {
		key = strings.ToUpper(key)
	}

	if pa.Mods > 0 && pa.Mods&^xproto.ModMask4 == 0 {
		// pa is <Super>Super_L or <Super>Super_R
		if keyLower == "super_l" || keyLower == "super_r" {
			pa.Mods = 0
		}
	}
	pa.Key = key
	logger.Debug("after fix", pa)
	return pa
}
