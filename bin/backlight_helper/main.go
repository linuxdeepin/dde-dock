/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import (
	"fmt"
	"os"
	"pkg.deepin.io/dde/daemon/backlight"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"sync"
	"time"
)

const (
	dbusDest = "com.deepin.daemon.helper.Backlight"
	dbusPath = "/com/deepin/daemon/helper/Backlight"
	dbusIFC  = "com.deepin.daemon.helper.Backlight"
)

var logger = log.NewLogger("backlight_helper")

type Manager struct {
	lcdInfos  backlight.SyspathInfos
	kbdInfos  backlight.SyspathInfos
	lcdLocker sync.Mutex
	kbdLocker sync.Mutex
}

// ListSysPath return all the backlight device syspath
func (m *Manager) ListSysPath() []string {
	var list []string
	for _, info := range m.lcdInfos {
		list = append(list, info.Path)
	}
	return list
}

// GetSysPathByType return the special type's syspath
// The type range: raw, platform, firmware
func (m *Manager) GetSysPathByType(ty string) (string, error) {
	for _, info := range m.lcdInfos {
		if info.Type == ty {
			return info.Path, nil
		}
	}
	return "", fmt.Errorf("Invalid backlight type: %s", ty)
}

// GetBrightness return the special syspath's brightness
func (m *Manager) GetBrightness(sysPath string) (int32, error) {
	m.lcdLocker.Lock()
	defer m.lcdLocker.Unlock()

	info, err := m.lcdInfos.Get(sysPath)
	if err != nil {
		return 0, err
	}
	return info.GetBrightness()
}

func (m *Manager) GetKbdBrightness() (int32, error) {
	m.kbdLocker.Lock()
	defer m.kbdLocker.Unlock()

	if len(m.kbdInfos) == 0 {
		return 0, fmt.Errorf("Unsupported keyboard backlight")
	}
	return m.kbdInfos[0].GetBrightness()
}

// GetBrightness return the special syspath's max brightness
func (m *Manager) GetMaxBrightness(sysPath string) (int32, error) {
	info, err := m.lcdInfos.Get(sysPath)
	if err != nil {
		return 0, err
	}
	return info.MaxBrightness, nil
}

func (m *Manager) GetKbdMaxBrightness() (int32, error) {
	if len(m.kbdInfos) == 0 {
		return 0, fmt.Errorf("Unsupported keyboard backlight")
	}
	return m.kbdInfos[0].MaxBrightness, nil
}

// SetBrightness set the special syspath's brightness
func (m *Manager) SetBrightness(sysPath string, value int32) error {
	m.lcdLocker.Lock()
	defer m.lcdLocker.Unlock()

	info, err := m.lcdInfos.Get(sysPath)
	if err != nil {
		return err
	}
	return info.SetBrightness(value)
}

func (m *Manager) SetKbdBrightness(value int32) error {
	m.kbdLocker.Lock()
	defer m.kbdLocker.Unlock()

	if len(m.kbdInfos) == 0 {
		return fmt.Errorf("Unsupported keyboard backlight")
	}
	return m.kbdInfos[0].SetBrightness(value)
}

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func main() {
	m := &Manager{
		lcdInfos: backlight.ListLCDBacklight(),
		kbdInfos: backlight.ListKbdBacklight(),
	}

	err := dbus.InstallOnSystem(m)
	if err != nil {
		logger.Error("Install session bus failed:", err)
		return
	}
	dbus.SetAutoDestroyHandler(time.Second*3, nil)
	dbus.DealWithUnhandledMessage()
	err = dbus.Wait()
	if err != nil {
		logger.Error("Lost dbus connection:", err)
		os.Exit(-1)
	}
	os.Exit(0)
}
