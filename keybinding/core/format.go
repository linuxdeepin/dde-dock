/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package core

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/keybind"
	"strings"
)

const (
	accelDelim = "-"
)

// FormatAccel format the accel as '<Control><Alt>T'
//
// accel: the input accel
// ret0: the formated accel
func FormatAccel(accel string) string {
	if len(accel) == 0 {
		return ""
	}

	return marshalKeys(unmarshalAccel(accel))
}

// FormatKeyEvent lookup keyevent to accel
func FormatKeyEvent(state uint16, detail int) (string, error) {
	xu, err := Initialize()
	if err != nil {
		return "", err
	}

	mod := keybind.ModifierString(state)
	key := keybind.LookupString(xu, state, xproto.Keycode(detail))
	if detail == 65 {
		key = "space"
	}

	if len(mod) != 0 {
		key = mod + "-" + key
	}

	return FormatAccel(key), nil
}

// IsAccelEqual check 's1', 's2' whether equal
func IsAccelEqual(s1, s2 string) bool {
	if len(s1) == 0 || len(s2) == 0 {
		return false
	}

	if s1 == s2 {
		return true
	}

	news1 := formatAccelToXGB(s1)
	news2 := formatAccelToXGB(s2)
	if news1 == news2 {
		return true
	}

	xu, err := Initialize()
	if err != nil {
		return false
	}

	mod1, keys1, _ := keybind.ParseString(xu, news1)
	mod2, keys2, _ := keybind.ParseString(xu, news2)
	if mod1 != mod2 {
		return false
	}

	return isKeycodesEqual(keys1, keys2)
}

// IsKeyMatch check 's' and 'mod, keycode' whether equal
func IsKeyMatch(s string, mod uint16, keycode int) bool {
	if len(s) == 0 {
		return false
	}

	tmp, err := FormatKeyEvent(mod, keycode)
	if err != nil {
		return false
	}

	return IsAccelEqual(s, tmp)
}

// formatAccelToXGB format the accel as 'control-alt-t'
func formatAccelToXGB(accel string) string {
	if len(accel) == 0 {
		return ""
	}

	return marshalKeysToXGB(unmarshalAccel(accel))
}

// To format '<Control><Alt>T'
func marshalKeys(keys []string) string {
	var (
		accel  string
		length = len(keys)
	)

	// ['control', 'shift', '%'] should be ['control','%']
	if isShiftKey(keys[length-1]) {
		keys, _ = delItemFromList("shift", keys)
		length = len(keys)
	}

	for i, k := range keys {
		if i == length-1 {
			break
		}

		accel += "<" + strings.Title(xgbModToKey(k)) + ">"
	}
	accel += strings.Title(keys[length-1])
	return accel
}

// To format 'control-alt-t'
func marshalKeysToXGB(keys []string) string {
	if len(keys) == 1 {
		return keys[0]
	}

	var (
		accel  string
		length = len(keys)
	)

	// ['control', '%'] should be ['control', 'shift', '%']
	if isShiftKey(keys[length-1]) {
		keys, _ = addItemToList("shift", keys)
		length = len(keys)
	}

	for i, k := range keys {
		if i == length-1 {
			break
		}
		accel += strings.ToLower(keyToXGBMod(k)) + accelDelim
	}
	accel += keysymToWeird(strings.ToLower(keys[length-1]))
	return accel
}

func unmarshalAccel(accel string) []string {
	var keys []string

	if hasFormated(accel) {
		keys = unmarshalFormatedAccel(accel)
	} else {
		keys = unmarshalXGBAccel(accel)
	}

	return uniqueKeys(filterSpecialKey(keys))
}

func unmarshalFormatedAccel(accel string) []string {
	var (
		keys  []string
		key   string
		start int
		end   int
		match bool
	)
	for i, ch := range accel {
		if ch == '<' {
			match = true
			start = i
			continue
		}

		if ch == '>' && match {
			end = i
			match = false

			var tmp string
			for j := start + 1; j < end; j++ {
				tmp += string(accel[j])
			}
			keys = append(keys, tmp)
			continue
		}

		if !match {
			key += string(accel[i])
		}
	}
	keys = append(keys, key)

	return keys
}

func unmarshalXGBAccel(accel string) []string {
	var (
		idx int = -1

		keys []string
		key  string
	)
	for i, ch := range accel {
		if ch == '-' {
			if (idx + 1) == i {
				idx = -1
				keys = append(keys, "-")
			} else {
				keys = append(keys, key)
				idx = i
				key = ""
			}
			continue
		}

		key += string(accel[i])
	}
	if len(key) > 0 {
		keys = append(keys, key)
	}

	return keys
}

func uniqueKeys(keys []string) []string {
	if len(keys) <= 1 {
		return keys
	}

	var ret []string
	for _, k := range keys {
		if strings.ToLower(k) == "primary" {
			k = "control"
		}

		if isItemInList(k, ret) {
			continue
		}

		ret = append(ret, k)
	}
	return ret
}

func filterSpecialKey(keys []string) []string {
	if len(keys) == 0 {
		return nil
	}

	var filterKeys = []string{
		"lock", "caps_lock",
		"mod2", "num_lock",
	}

	var ret []string
	length := len(keys)
	for i, key := range keys {
		// if length == 1 && key == 'caps_lock'
		if i == length-1 {
			break
		}

		if isItemInList(strings.ToLower(key), filterKeys) {
			continue
		}

		ret = append(ret, key)
	}
	ret = append(ret, keys[length-1])

	return ret
}

func hasFormated(accel string) bool {
	if accel[0] == '<' {
		return true
	}

	return false
}

func keyToXGBMod(key string) string {
	switch strings.ToLower(key) {
	case "caps_lock":
		return "lock"
	case "alt":
		return "mod1"
	case "meta":
		return "mod1"
	case "num_lock":
		return "mod2"
	case "super":
		return "mod4"
	case "hyper":
		return "mod4"
	}

	return key
}

func xgbModToKey(mod string) string {
	switch strings.ToLower(mod) {
	case "mod1":
		return "alt"
	case "mod2":
		return "num_lock"
	case "mod4":
		return "super"
	case "lock":
		return "caps_lock"
	}

	return mod
}

func isShiftKey(key string) bool {
	var keys = []string{
		"~", "!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "+",
		"{", "}", "|", ":", "\"", "<", ">", "?",
	}

	return isItemInList(key, keys)
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

func isItemInList(item string, items []string) bool {
	for _, v := range items {
		if v == item {
			return true
		}
	}

	return false
}

func addItemToList(item string, list []string) ([]string, bool) {
	if isItemInList(item, list) {
		return list, false
	}

	var ret []string
	length := len(list)
	ret = append(ret, list[:length-1]...)
	ret = append(ret, item)
	ret = append(ret, list[length-1])

	return ret, true
}

func delItemFromList(item string, list []string) ([]string, bool) {
	var (
		deleted bool
		ret     []string
	)
	for _, k := range list {
		if item == k {
			deleted = true
			continue
		}

		ret = append(ret, k)
	}
	return ret, deleted
}
