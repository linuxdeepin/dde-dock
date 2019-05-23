package screensaver

import (
	"encoding/json"
	"os"
	"path/filepath"

	"pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/keyfile"
	"pkg.deepin.io/lib/xdg/basedir"
)

const (
	dScreenSaverPath        = "/com/deepin/daemon/ScreenSaver"
	dScreenSaverServiceName = "com.deepin.daemon.ScreenSaver"

	gsSchemaPower                  = "com.deepin.dde.power"
	gsKeyBatteryScreensaverDelay   = "battery-screensaver-delay"
	gsKeyLinePowerScreensaverDelay = "line-power-screensaver-delay"

	sectionGeneral       = "General"
	keyCurrent           = "currentScreenSaver"
	keyLockScreenAtAwake = "lockScreenAtAwake"
)

var dScreensaverConfigFile = filepath.Join(basedir.GetUserConfigDir(), "deepin/deepin-screensaver.conf")

type syncConfig struct {
}

func (sc *syncConfig) Get() (interface{}, error) {
	gs := gio.NewSettings(gsSchemaPower)
	defer gs.Unref()
	var v syncData
	v.BatteryDelay = int(gs.GetInt(gsKeyBatteryScreensaverDelay))
	v.LinePowerDelay = int(gs.GetInt(gsKeyLinePowerScreensaverDelay))

	kf := keyfile.NewKeyFile()
	err := kf.LoadFromFile(dScreensaverConfigFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return &v, nil
	}
	v.Current, _ = kf.GetString(sectionGeneral, keyCurrent)
	v.LockScreenAtAwake, _ = kf.GetBool(sectionGeneral, keyLockScreenAtAwake)
	return &v, nil
}

func (sc *syncConfig) Set(data []byte) error {
	var v syncData
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	gs := gio.NewSettings(gsSchemaPower)
	defer gs.Unref()
	gs.SetInt(gsKeyBatteryScreensaverDelay, int32(v.BatteryDelay))
	gs.SetInt(gsKeyLinePowerScreensaverDelay, int32(v.LinePowerDelay))

	kf := keyfile.NewKeyFile()
	err = kf.LoadFromFile(dScreensaverConfigFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	kf.SetString(sectionGeneral, keyCurrent, v.Current)
	kf.SetBool(sectionGeneral, keyLockScreenAtAwake, v.LockScreenAtAwake)
	return kf.SaveToFile(dScreensaverConfigFile)
}

// version: 1.0
type syncData struct {
	Version           string `json:"version"`
	BatteryDelay      int    `json:"battery_delay"`
	LinePowerDelay    int    `json:"line_power_delay"`
	LockScreenAtAwake bool   `json:"lock_screen_at_awake"`
	Current           string `json:"current"`
}
