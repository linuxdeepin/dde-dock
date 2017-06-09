package keybind

import (
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

// GrabChecked Grabs a key with mods on a particular window.
// This is the same as Grab, except that it issue a checked request.
// Which means that an error could be returned and handled on the spot.
// (Checked requests are slower than unchecked requests.)
// This will also grab all combinations of modifiers found in xevent.IgnoreMods.
func GrabChecked(xu *xgbutil.XUtil, win xproto.Window,
	mods uint16, key xproto.Keycode) error {

	var err error
	for _, m := range xevent.IgnoreMods {
		err = xproto.GrabKeyChecked(xu.Conn(), true, win, mods|m, key,
			xproto.GrabModeAsync, xproto.GrabModeAsync).Check()
		if err != nil {
			return err
		}
	}
	return nil
}

// Ungrab undoes Grab. It will handle all combinations od modifiers found
// in xevent.IgnoreMods.
func Ungrab(xu *xgbutil.XUtil, win xproto.Window,
	mods uint16, key xproto.Keycode) {

	for _, m := range xevent.IgnoreMods {
		xproto.UngrabKeyChecked(xu.Conn(), key, win, mods|m).Check()
	}
}

// GrabKeyboard grabs the entire keyboard.
// Returns whether GrabStatus is successful and an error if one is reported by
// XGB. It is possible to not get an error and the grab to be unsuccessful.
// The purpose of 'win' is that after a grab is successful, ALL Key*Events will
// be sent to that window. Make sure you have a callback attached :-)
func GrabKeyboard(xu *xgbutil.XUtil, win xproto.Window) error {
	reply, err := xproto.GrabKeyboard(xu.Conn(), false, win, 0,
		xproto.GrabModeAsync, xproto.GrabModeAsync).Reply()
	if err != nil {
		return fmt.Errorf("GrabKeyboard: Error grabbing keyboard on "+
			"window '%x': %s", win, err)
	}

	switch reply.Status {
	case xproto.GrabStatusSuccess:
		// all is well
	case xproto.GrabStatusAlreadyGrabbed:
		return fmt.Errorf("GrabKeyboard: Could not grab keyboard. " +
			"Status: AlreadyGrabbed.")
	case xproto.GrabStatusInvalidTime:
		return fmt.Errorf("GrabKeyboard: Could not grab keyboard. " +
			"Status: InvalidTime.")
	case xproto.GrabStatusNotViewable:
		return fmt.Errorf("GrabKeyboard: Could not grab keyboard. " +
			"Status: NotViewable.")
	case xproto.GrabStatusFrozen:
		return fmt.Errorf("GrabKeyboard: Could not grab keyboard. " +
			"Status: Frozen.")
	}
	return nil
}

// UngrabKeyboard undoes GrabKeyboard.
func UngrabKeyboard(xu *xgbutil.XUtil) {
	xproto.UngrabKeyboard(xu.Conn(), 0)
}
