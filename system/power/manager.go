/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

import (
	"errors"
	"gir/gudev-1.0"
	"pkg.deepin.io/dde/api/powersupply"
	"pkg.deepin.io/dde/api/powersupply/battery"
	"pkg.deepin.io/lib/dbus"
	"sync"
	"time"
)

// https://www.kernel.org/doc/Documentation/power/power_supply_class.txt
type Manager struct {
	OnBattery   bool
	batteries   map[string]*Battery
	ac          *AC
	gudevClient *gudev.Client
	mutex       sync.Mutex

	// battery display properties:
	HasBattery         bool
	BatteryPercentage  float64
	BatteryStatus      battery.Status
	BatteryTimeToEmpty uint64
	BatteryTimeToFull  uint64

	HasLidSwitch bool

	// Signals:
	BatteryDisplayUpdate func(timestamp int64)
	BatteryAdded         func(objpath string)
	BatteryRemoved       func(objpath string)
	LidClosed            func()
	LidOpened            func()
}

func NewManager() (*Manager, error) {
	m := &Manager{}
	err := m.init()
	if err != nil {
		m.destroy()
		return nil, err
	}
	return m, nil
}

type AC struct {
	gudevClient *gudev.Client
	sysfsPath   string
}

func NewAC(manager *Manager, device *gudev.Device) *AC {
	sysfsPath := device.GetSysfsPath()
	return &AC{
		gudevClient: manager.gudevClient,
		sysfsPath:   sysfsPath,
	}
}

func (ac *AC) newDevice() *gudev.Device {
	return ac.gudevClient.QueryBySysfsPath(ac.sysfsPath)
}

func (m *Manager) refreshAC(ac *gudev.Device) {
	online := ac.GetPropertyAsBoolean("POWER_SUPPLY_ONLINE")
	logger.Debug("ac online:", online)
	m.setPropOnBattery(!online)
}

func (m *Manager) initAC(devices []*gudev.Device) {
	var ac *gudev.Device
	for _, dev := range devices {
		if powersupply.IsMains(dev) {
			ac = dev
			break
		}
	}
	if ac != nil {
		m.refreshAC(ac)
		m.ac = NewAC(m, ac)
	}
}

func (m *Manager) init() error {
	subsystems := []string{"power_supply", "input"}
	m.gudevClient = gudev.NewClient(subsystems)
	if m.gudevClient == nil {
		return errors.New("gudevClient is nil")
	}

	m.initLidSwitch()
	devices := powersupply.GetDevices(m.gudevClient)
	m.initAC(devices)
	m.initBatteries(devices)
	for _, dev := range devices {
		dev.Unref()
	}

	m.gudevClient.Connect("uevent", m.handleUEvent)

	err := dbus.InstallOnSystem(m)
	if err != nil {
		logger.Warning("Install manager failed:", err)
		return err
	}

	return nil
}

func (m *Manager) handleUEvent(client *gudev.Client, action string, device *gudev.Device) {
	logger.Debug("on uevent action:", action)
	defer device.Unref()

	switch action {
	case "change":
		if powersupply.IsMains(device) {
			if m.ac == nil {
				m.ac = NewAC(m, device)
			} else if m.ac.sysfsPath != device.GetSysfsPath() {
				logger.Warning("found another AC", device.GetSysfsPath())
				return
			}

			// now m.ac != nil, and sysfsPath equal
			m.refreshAC(device)
			time.AfterFunc(1*time.Second, m.RefreshBatteries)
			time.AfterFunc(3*time.Second, m.RefreshBatteries)
		} else if powersupply.IsSystemBattery(device) {
			m.addBattery(device)
		}
	case "add":
		if powersupply.IsSystemBattery(device) {
			m.addBattery(device)
		}
		// ignore add mains

	case "remove":
		m.removeBattery(device)
	}

}

func (m *Manager) initBatteries(devices []*gudev.Device) {
	m.batteries = make(map[string]*Battery)
	for _, dev := range devices {
		m.addBattery(dev)
	}
	logger.Debugf("initBatteries done %#v", m.batteries)
}

func (m *Manager) addBattery(dev *gudev.Device) (*Battery, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	logger.Debug("addBattery dev:", dev)
	if !powersupply.IsSystemBattery(dev) {
		return nil, false
	}

	sysfsPath := dev.GetSysfsPath()
	logger.Debug(sysfsPath)
	bat0, ok := m.batteries[sysfsPath]
	if ok {
		logger.Debugf("add battery failed , sysfsPath exists %q", sysfsPath)
		bat0.Refresh()
		return bat0, false
	}

	bat := NewBattery(m, dev)
	if bat == nil {
		logger.Warning("add batteries failed, sysfsPath %q, new batttery failed", sysfsPath)
		return nil, false
	}
	err := dbus.InstallOnSystem(bat)
	if err != nil {
		logger.Warning("Install battery failed:", err)
		bat.destroy()
		return nil, false
	}

	m.batteries[sysfsPath] = bat
	m.refreshBatteryDisplay()
	bat.setRefreshDoneCallback(m.refreshBatteryDisplay)
	// signal BatteryAdded
	logger.Debug("before emit BatteryAdded")
	dbus.Emit(m, "BatteryAdded", string(bat.GetDBusInfo().ObjectPath))
	logger.Debug("after emit BatteryAdded")
	return bat, true
}

func (m *Manager) removeBattery(dev *gudev.Device) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	sysfsPath := dev.GetSysfsPath()
	bat, ok := m.batteries[sysfsPath]
	if ok {
		logger.Info("removeBattery", sysfsPath)
		dbus.UnInstallObject(bat)
		bat.destroy()
		delete(m.batteries, sysfsPath)
		m.refreshBatteryDisplay()
		// signal BatteryRemoved
		dbus.Emit(m, "BatteryRemoved", string(bat.GetDBusInfo().ObjectPath))
	} else {
		logger.Warning("removeBattery failed: invalid sysfsPath ", sysfsPath)
	}
}

func (m *Manager) destroy() {
	logger.Debug("destroy")
	for _, bat := range m.batteries {
		bat.destroy()
	}
	m.batteries = nil

	if m.gudevClient != nil {
		m.gudevClient.Unref()
		m.gudevClient = nil
	}
}
