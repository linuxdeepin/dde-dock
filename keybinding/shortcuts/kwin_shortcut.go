package shortcuts

import (
	"errors"

	"github.com/linuxdeepin/go-dbus-factory/com.deepin.wm"
	"pkg.deepin.io/dde/daemon/keybinding/util"
)

type kWinShortcut struct {
	BaseShortcut
	wm *wm.Wm
}

func newKWinShortcut(id, name string, keystrokes []string, wm *wm.Wm) *kWinShortcut {
	return &kWinShortcut{
		BaseShortcut: BaseShortcut{
			Id:         id,
			Type:       ShortcutTypeWM,
			Name:       name,
			Keystrokes: ParseKeystrokes(keystrokes),
		},
		wm: wm,
	}
}

func (ks *kWinShortcut) ReloadKeystrokes() bool {
	oldVal := ks.GetKeystrokes()
	keystrokes, err := ks.wm.GetAccel(0, ks.Id)
	if err != nil {
		logger.Warning("failed to get accel for %s: %v", ks.Id, err)
		return false
	}
	newVal := ParseKeystrokes(keystrokes)
	ks.setKeystrokes(newVal)
	return !keystrokesEqual(oldVal, newVal)
}

func (ks *kWinShortcut) SaveKeystrokes() error {
	accelJson, err := util.MarshalJSON(util.KWinAccel{
		Id:         ks.Id,
		Keystrokes: ks.getKeystrokesStrv(),
	})
	if err != nil {
		return err
	}

	ok, err := ks.wm.SetAccel(0, accelJson)
	if !ok {
		return errors.New("wm.SetAccel failed, id: " + ks.Id)
	}
	return err
}

func (ks *kWinShortcut) ShouldEmitSignalChanged() bool {
	return true
}
