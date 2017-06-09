package keybinding

import (
	"errors"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgb/xtest"
	"github.com/BurntSushi/xgbutil"
	"pkg.deepin.io/dde/daemon/keybinding/shortcuts"
)

type NumLockState uint

const (
	NumLockOff NumLockState = iota
	NumLockOn
	NumLockUnknown
)

type CapsLockState uint

const (
	CapsLockOff CapsLockState = iota
	CapsLockOn
	CapsLockUnknown
)

func queryNumLockState(xu *xgbutil.XUtil) (NumLockState, error) {
	queryPointerReply, err := xproto.QueryPointer(xu.Conn(), xu.RootWin()).Reply()
	if err != nil {
		return NumLockUnknown, err
	}
	logger.Debugf("query pointer reply %#v", queryPointerReply)
	on := queryPointerReply.Mask&xproto.ModMask2 != 0
	if on {
		return NumLockOn, nil
	} else {
		return NumLockOff, nil
	}
}

func queryCapsLockState(xu *xgbutil.XUtil) (CapsLockState, error) {
	queryPointerReply, err := xproto.QueryPointer(xu.Conn(), xu.RootWin()).Reply()
	if err != nil {
		return CapsLockUnknown, err
	}
	logger.Debugf("query pointer reply %#v", queryPointerReply)
	on := queryPointerReply.Mask&xproto.ModMaskLock != 0
	if on {
		return CapsLockOn, nil
	} else {
		return CapsLockOff, nil
	}
}

func setNumLockState(xu *xgbutil.XUtil, state NumLockState) error {
	if !(state == NumLockOff || state == NumLockOn) {
		return errors.New("invalid numlock state")
	}

	state0, err := queryNumLockState(xu)
	if err != nil {
		return err
	}

	if state0 != state {
		return changeNumLockState(xu)
	}
	return nil
}

func changeNumLockState(xu *xgbutil.XUtil) (err error) {
	// get Num_Lock keycode
	code, err := shortcuts.GetKeyFirstCode(xu, "Num_Lock")
	if err != nil {
		return err
	}
	numLockKeycode := byte(code)
	logger.Debug("numLockKeycode is", numLockKeycode)

	x := xu.Conn()
	root := xu.RootWin()

	// fake key press
	err = xtest.FakeInputChecked(x, xproto.KeyPress, numLockKeycode, xproto.TimeCurrentTime, root, 0, 0, 0).Check()
	if err != nil {
		return err
	}
	// fake key release
	err = xtest.FakeInputChecked(x, xproto.KeyRelease, numLockKeycode, xproto.TimeCurrentTime, root, 0, 0, 0).Check()
	if err != nil {
		return err
	}
	return nil
}
