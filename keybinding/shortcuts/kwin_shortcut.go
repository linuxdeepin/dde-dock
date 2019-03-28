package shortcuts

import (
	"encoding/json"
	"errors"

	"github.com/linuxdeepin/go-dbus-factory/com.deepin.wm"
)

type kwinShortcut struct {
	BaseShortcut
	wm *wm.Wm
}

func NewKwinShortcut(id, name string, keystrokes []string, wm *wm.Wm) *kwinShortcut {
	return &kwinShortcut{
		BaseShortcut: BaseShortcut{
			Id:         id,
			Type:       ShortcutTypeWM,
			Name:       name,
			Keystrokes: ParseKeystrokes(keystrokes),
		},
		wm: wm,
	}
}

func (ks *kwinShortcut) ReloadKeystrokes() bool {
	return false
}

func (ks *kwinShortcut) SaveKeystrokes() error {
	data, err := json.Marshal(kwinAccel{
		Id:         ks.Id,
		Keystrokes: ks.getKeystrokesStrv(),
	})
	if err != nil {
		return err
	}

	ok, err := ks.wm.SetAccel(0, string(data))
	if !ok {
		return errors.New("wm.SetAccel failed")
	}
	return err
}

func (ks *kwinShortcut) ShouldEmitSignalChanged() bool {
	return true
}
