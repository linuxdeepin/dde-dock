package lastore

import (
	"encoding/json"
	"errors"
	"pkg.deepin.io/lib/dbus1"
)

type syncConfig struct {
	l *Lastore
}

type syncData struct {
	Version             string `json:"version"`
	AutoCheckUpdates    bool   `json:"auto_check_updates"`
	AutoClean           bool   `json:"auto_clean"`
	AutoDownloadUpdates bool   `json:"auto_download_updates"`
	SmartMirrorEnabled  bool   `json:"smart_mirror_enabled"`
	SourceCheckEnabled  bool   `json:"source_check_enabled"`
}

const (
	syncVersion = "1.0"

	smartMirrorService = "com.deepin.lastore.Smartmirror"
	smartMirrorPath    = "/com/deepin/lastore/Smartmirror"
	smartMirrorIFC     = smartMirrorService
)

func (sc *syncConfig) Get() (interface{}, error) {
	var info syncData
	info.Version = syncVersion
	info.AutoCheckUpdates, _ = sc.l.core.AutoCheckUpdates().Get(0)
	info.AutoClean, _ = sc.l.core.AutoClean().Get(0)
	info.AutoDownloadUpdates, _ = sc.l.core.AutoDownloadUpdates().Get(0)
	info.SmartMirrorEnabled, _ = smartMirrorEnabledGet()
	info.SourceCheckEnabled = sc.l.SourceCheckEnabled
	return &info, nil
}

func (sc *syncConfig) Set(data []byte) error {
	var info syncData
	err := json.Unmarshal(data, &info)
	if err != nil {
		return err
	}
	sc.l.core.AutoCheckUpdates().Set(0, info.AutoCheckUpdates)
	sc.l.core.AutoClean().Set(0, info.AutoClean)
	sc.l.core.AutoDownloadUpdates().Set(0, info.AutoDownloadUpdates)
	smartMirrorEnabledSet(info.SmartMirrorEnabled)
	sc.l.SetSourceCheckEnabled(info.SourceCheckEnabled)
	return nil
}

func smartMirrorEnabledSet(enabled bool) error {
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	return conn.Object(smartMirrorService, smartMirrorPath).Call(
		smartMirrorIFC+".SetEnable", 0, enabled).Store()
}

func smartMirrorEnabledGet() (bool, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return false, err
	}
	var variant dbus.Variant
	err = conn.Object(smartMirrorService, smartMirrorPath).Call(
		"org.freedesktop.DBus.Properties.Get", 0, smartMirrorIFC, "Enable").Store(&variant)
	if err != nil {
		return false, err
	}

	if variant.Signature().String() != "b" {
		return false, errors.New("Not excepted value type")
	}
	return variant.Value().(bool), nil
}
