package keybinding

import (
	"strings"
)

var keysymWeirdsMap = map[string]string{
	"-":  "minus",
	"=":  "equal",
	"\\": "backslash",
	"?":  "question",
	"!":  "exclam",
	"#":  "numbersign",
	";":  "semicolon",
	"'":  "apostrophe",
	"<":  "less",
	".":  "period",
	"/":  "slash",
	"(":  "parenleft",
	"[":  "bracketleft",
	")":  "parenright",
	"]":  "bracketright",
	"\"": "quotedbl",
	" ":  "space",
	"$":  "dollar",
	"+":  "plus",
	"*":  "asterisk",
	"_":  "underscore",
	"|":  "bar",
	"`":  "grave",
	"@":  "at",
	"%":  "percent",
	">":  "greater",
	"^":  "asciicircum",
	"{":  "braceleft",
	":":  "colon",
	",":  "comma",
	"~":  "asciitilde",
	"&":  "ampersand",
	"}":  "braceright",
}

func convertKeysym2Weird(shortcut string) string {
	if len(shortcut) < 1 {
		return ""
	}

	array := strings.Split(shortcut, "--")
	strs := strings.Split(array[0], ACCEL_DELIM)

	tmp := ""
	for i, k := range strs {
		if v, ok := keysymWeirdsMap[k]; ok {
			k = v
		}

		if i != 0 {
			tmp += ACCEL_DELIM + k
		} else {
			tmp = k
		}
	}

	if len(array) > 1 {
		tmp += ACCEL_DELIM + keysymWeirdsMap[ACCEL_DELIM]
	}

	return tmp
}
