package keybind

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"strings"
	"unicode"
)

var KbdMappingNotifyCallback func()

// Initialize attaches the appropriate callbacks to make key bindings easier.
// i.e., update state of the world on a MappingNotify.
func Initialize(xu *xgbutil.XUtil) {
	// Listen to mapping notify events
	xevent.MappingNotifyFun(mappingNotifyHandler).Connect(xu, xevent.NoWindow)

	// Give us an initial mapping state...
	keyMap, modMap := MapsGet(xu)
	KeyMapSet(xu, keyMap)
	ModMapSet(xu, modMap)
}

func mappingNotifyHandler(xu *xgbutil.XUtil, ev xevent.MappingNotifyEvent) {
	keyMap, modMap := MapsGet(xu)
	KeyMapSet(xu, keyMap)
	ModMapSet(xu, modMap)

	if ev.Request == xproto.MappingKeyboard {
		if KbdMappingNotifyCallback != nil {
			KbdMappingNotifyCallback()
		}
	}
}

func StrToKeycodes(xu *xgbutil.XUtil, str string) ([]xproto.Keycode, error) {
	sym, ok := strToKeysym(str)
	if !ok {
		return nil, errors.New("fail to get keysym")
	}

	return keysymToKeycodes(xu, sym), nil
}

func keysymToKeycodes(xu *xgbutil.XUtil, keysym xproto.Keysym) []xproto.Keycode {
	min, max := minMaxKeycodeGet(xu)
	keyMap := KeyMapGet(xu)
	if keyMap == nil {
		panic("keybind.Initialize must be called before using the keybind " +
			"package.")
	}

	symsLen := int(keyMap.KeysymsPerKeycode)
	keycodes := make([]xproto.Keycode, symsLen)
	offset := 0
	for kc := int(min); kc <= int(max); kc++ {
		for c := 0; c < symsLen; c++ {
			if keysym == keyMap.Keysyms[offset+c] && keycodes[c] == 0 {
				keycodes[c] = xproto.Keycode(kc)
			}
		}
		offset += symsLen
	}
	return keycodes
}

// LookupString attempts to convert a (modifiers, keycode) to an english string.
// It essentially implements the rules described at http://goo.gl/qum9q
// Namely, the bulleted list that describes how key syms should be interpreted
// when various modifiers are pressed.
// Note that we ignore the logic that asks us to check if particular key codes
// are mapped to particular modifiers (i.e., "XK_Caps_Lock" to "Lock" modifier).
// We just check if the modifiers are activated. That's good enough for me.
// XXX: We ignore num lock stuff.
// XXX: We ignore MODE SWITCH stuff. (i.e., we don't use group 2 key syms.)
func LookupString(xu *xgbutil.XUtil, mods uint16,
	keycode xproto.Keycode) string {

	k1, k2, _, _ := interpretSymList(xu, keycode)

	shift := mods&xproto.ModMaskShift > 0
	lock := mods&xproto.ModMaskLock > 0
	switch {
	case !shift && !lock:
		return k1
	case !shift && lock:
		if len(k1) == 1 && unicode.IsLower(rune(k1[0])) {
			return k2
		} else {
			return k1
		}
	case shift && lock:
		if len(k2) == 1 && unicode.IsLower(rune(k2[0])) {
			return string(unicode.ToUpper(rune(k2[0])))
		} else {
			return k2
		}
	case shift:
		return k2
	}

	return ""
}

// interpretSymList interprets the keysym list for a particular keycode as
// described in the third and fourth paragraphs of http://goo.gl/qum9q
func interpretSymList(xu *xgbutil.XUtil, keycode xproto.Keycode) (
	k1 string, k2 string, k3 string, k4 string) {

	ks1 := KeysymGet(xu, keycode, 0)
	ks2 := KeysymGet(xu, keycode, 1)
	ks3 := KeysymGet(xu, keycode, 2)
	ks4 := KeysymGet(xu, keycode, 3)

	// follow the rules, third paragraph
	switch {
	case ks2 == 0 && ks3 == 0 && ks4 == 0:
		ks3 = ks1
	case ks3 == 0 && ks4 == 0:
		ks3 = ks1
		ks4 = ks2
	case ks4 == 0:
		ks4 = 0
	}

	// Now convert keysyms to strings, so we can do alphabetic shit.
	k1 = KeysymToStr(ks1)
	k2 = KeysymToStr(ks2)
	k3 = KeysymToStr(ks3)
	k4 = KeysymToStr(ks4)

	// follow the rules, fourth paragraph
	if k2 == "" {
		if len(k1) == 1 && unicode.IsLetter(rune(k1[0])) {
			k1 = string(unicode.ToLower(rune(k1[0])))
			k2 = string(unicode.ToUpper(rune(k1[0])))
		} else {
			k2 = k1
		}
	}
	if k4 == "" {
		if len(k3) == 1 && unicode.IsLetter(rune(k3[0])) {
			k3 = string(unicode.ToLower(rune(k3[0])))
			k4 = string(unicode.ToUpper(rune(k4[0])))
		} else {
			k4 = k3
		}
	}

	return
}

// KeysymGet is a shortcut alias for 'KeysymGetWithMap' using the current
// keymap stored in XUtil.
// keybind.Initialize MUST have been called before using this function.
func KeysymGet(xu *xgbutil.XUtil, keycode xproto.Keycode,
	column byte) xproto.Keysym {

	return KeysymGetWithMap(xu, KeyMapGet(xu), keycode, column)
}

// KeysymGetWithMap uses the given key map and finds a keysym associated
// with the given keycode in the current X environment.
func KeysymGetWithMap(xu *xgbutil.XUtil, keyMap *xgbutil.KeyboardMapping,
	keycode xproto.Keycode, column byte) xproto.Keysym {

	min, _ := minMaxKeycodeGet(xu)
	i := (int(keycode)-int(min))*int(keyMap.KeysymsPerKeycode) + int(column)

	return keyMap.Keysyms[i]
}

// KeysymToStr converts a keysym to a string if one is available.
// If one is found, KeysymToStr also checks the 'weirdKeysyms' map, which
// contains a map from multi-character strings to single character
// representations (i.e., 'braceleft' to '{').
// If no match is found initially, an empty string is returned.
func KeysymToStr(keysym xproto.Keysym) string {
	symStr, ok := strKeysyms[keysym]
	if !ok {
		return ""
	}

	shortSymStr, ok := weirdKeysyms[symStr]
	if ok {
		return string(shortSymStr)
	}

	return symStr
}

func strToKeysym(str string) (xproto.Keysym, bool) {
	// Do some fancy case stuff before we give up.
	sym, ok := keysyms[str]
	if !ok {
		sym, ok = keysyms[strings.Title(str)]
	}
	if !ok {
		sym, ok = keysyms[strings.ToLower(str)]
	}
	if !ok {
		sym, ok = keysyms[strings.ToUpper(str)]
	}

	return sym, ok
}

// A convenience function to grab the KeyboardMapping and ModifierMapping
// from X. We need to do this on startup (see Initialize) and whenever we
// get a MappingNotify event.
func MapsGet(xu *xgbutil.XUtil) (*xproto.GetKeyboardMappingReply,
	*xproto.GetModifierMappingReply) {

	min, max := minMaxKeycodeGet(xu)
	newKeymap, keyErr := xproto.GetKeyboardMapping(xu.Conn(), min,
		byte(max-min+1)).Reply()
	newModmap, modErr := xproto.GetModifierMapping(xu.Conn()).Reply()

	// If there are errors, we really need to panic. We just can't do
	// any key binding without a mapping from the server.
	if keyErr != nil {
		panic(fmt.Sprintf("COULD NOT GET KEYBOARD MAPPING: %v\n"+
			"THIS IS AN UNRECOVERABLE ERROR.\n",
			keyErr))
	}
	if modErr != nil {
		panic(fmt.Sprintf("COULD NOT GET MODIFIER MAPPING: %v\n"+
			"THIS IS AN UNRECOVERABLE ERROR.\n",
			keyErr))
	}

	return newKeymap, newModmap
}

// KeyMapGet accessor.
func KeyMapGet(xu *xgbutil.XUtil) *xgbutil.KeyboardMapping {
	return xu.Keymap
}

// KeyMapSet updates XUtil.keymap.
// This is exported for use in the keybind package. You probably shouldn't
// use this. (You may need to use this if you're rolling your own event loop,
// and still want to use the keybind package.)
func KeyMapSet(xu *xgbutil.XUtil, keyMapReply *xproto.GetKeyboardMappingReply) {
	xu.Keymap = &xgbutil.KeyboardMapping{keyMapReply}
}

// ModMapGet accessor.
func ModMapGet(xu *xgbutil.XUtil) *xgbutil.ModifierMapping {
	return xu.Modmap
}

// ModMapSet updates XUtil.modmap.
// This is exported for use in the keybind package. You probably shouldn't
// use this. (You may need to use this if you're rolling your own event loop,
// and still want to use the keybind package.)
func ModMapSet(xu *xgbutil.XUtil, modMapReply *xproto.GetModifierMappingReply) {
	xu.Modmap = &xgbutil.ModifierMapping{modMapReply}
}

// minMaxKeycodeGet a simple accessor to the X setup info to return the
// minimum and maximum keycodes. They are typically 8 and 255, respectively.
func minMaxKeycodeGet(xu *xgbutil.XUtil) (xproto.Keycode, xproto.Keycode) {
	return xu.Setup().MinKeycode, xu.Setup().MaxKeycode
}
