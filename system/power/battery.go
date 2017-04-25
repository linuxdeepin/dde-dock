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
	"gir/gudev-1.0"
	"pkg.deepin.io/dde/api/powersupply/battery"
	"pkg.deepin.io/lib/dbus"
	"sync"
	"time"
)

type Battery struct {
	exit  chan struct{}
	mutex sync.Mutex

	gudevClient       *gudev.Client
	changedProperties []string
	SysfsPath         string
	IsPresent         bool

	Name             string
	Technology       string
	Manufacturer     string
	ModelName        string
	SerialNumber     string
	EnergyFullDesign uint64
	VoltageMinDesign uint64

	EnergyFull  uint64
	EnergyNow   uint64
	PowerNow    uint64
	VoltageNow  uint64
	Percentage  float64
	Capacity    float64
	Status      battery.Status
	TimeToEmpty uint64
	TimeToFull  uint64
	UpdateTime  int64

	refreshDone func()
}

func NewBattery(manager *Manager, device *gudev.Device) *Battery {
	sysfsPath := device.GetSysfsPath()
	logger.Debugf("NewBattery sysfsPath: %q", sysfsPath)
	if manager == nil || manager.gudevClient == nil ||
		device == nil {
		return nil
	}
	bat := &Battery{
		gudevClient: manager.gudevClient,
		SysfsPath:   sysfsPath,
	}
	bat.refresh(device)
	bat.resetUpdateInterval(60 * time.Second)
	return bat
}

func (bat *Battery) setRefreshDoneCallback(fn func()) {
	bat.refreshDone = fn
}

func (bat *Battery) newDevice() *gudev.Device {
	return bat.gudevClient.QueryBySysfsPath(bat.SysfsPath)
}

func (bat *Battery) notifyChangeStart() {
	if bat.changedProperties != nil &&
		len(bat.changedProperties) > 0 {
		logger.Warning("some properties change notify may lose")
	}
	bat.changedProperties = make([]string, 0, 5)
}

func (bat *Battery) notifyChangeEnd() {
	logger.Debugf("changed props len: %v , %v",
		len(bat.changedProperties), bat.changedProperties)
	for _, propName := range bat.changedProperties {
		dbus.NotifyChange(bat, propName)
	}
	bat.changedProperties = nil
}

func (bat *Battery) notifyChange(propNames ...string) {
	bat.changedProperties = append(bat.changedProperties, propNames...)
}

func (bat *Battery) resetValues() {
	logger.Debug("resetValues")
	bat.setPropIsPresent(false)
	bat.setPropUpdateTime(0)

	bat.setPropName("")
	bat.setPropTechnology("")
	bat.setPropManufacturer("")
	bat.setPropModelName("")
	bat.setPropSerialNumber("")
	bat.setPropEnergyFullDesign(0)
	bat.setPropVoltageMinDesign(0)

	bat.setPropEnergyNow(0)
	bat.setPropEnergyFull(0)
	bat.setPropVoltageNow(0)
	bat.setPropPowerNow(0)
	bat.setPropPercentage(0)
	bat.setPropCapacity(0)
	bat.setPropStatus(battery.StatusUnknown)
	bat.setPropTimeToEmpty(0)
	bat.setPropTimeToFull(0)
}

func (bat *Battery) refresh(dev *gudev.Device) {
	bat.notifyChangeStart()
	batInfo := battery.GetBatteryInfo(dev)
	bat._refresh(batInfo)
	bat.notifyChangeEnd()
}

func (bat *Battery) _refresh(info *battery.BatteryInfo) {
	logger.Debug("Refresh", bat.Name)

	defer func() {
		logger.Debugf("Refresh %v done", bat.Name)
		if bat.refreshDone != nil {
			bat.refreshDone()
		}
	}()

	if info == nil {
		bat.resetValues()
		return
	}
	bat.setPropIsPresent(true)
	now := time.Now()
	updateTime := now.Unix()
	logger.Debugf("now %v updateTime %v", now, updateTime)
	bat.setPropUpdateTime(updateTime)

	logger.Debug("Name", info.Name)
	bat.setPropName(info.Name)

	logger.Debug("Technology", info.Technology)
	bat.setPropTechnology(info.Technology)

	logger.Debug("Manufacturer", info.Manufacturer)
	bat.setPropManufacturer(info.Manufacturer)

	logger.Debug("ModelName", info.ModelName)
	bat.setPropModelName(info.ModelName)

	logger.Debug("SerialNumber", info.SerialNumber)
	bat.setPropSerialNumber(info.SerialNumber)

	logger.Debugf("EnergyFullDesign %vµAh", info.EnergyFullDesign)
	bat.setPropEnergyFullDesign(info.EnergyFullDesign)

	logger.Debugf("VoltageMinDesign %vµV", info.VoltageMinDesign)
	bat.setPropVoltageMinDesign(info.VoltageMinDesign)

	logger.Debugf("energyNow %vµAh", info.EnergyNow)
	bat.setPropEnergyNow(info.EnergyNow)

	logger.Debugf("energyFull %vµAh", info.EnergyFull)
	if info.EnergyFull > info.EnergyFullDesign {
		logger.Warningf("EnergyFull (%v) is greater then EnergyFullDesign (%v)",
			info.EnergyFull, info.EnergyFullDesign)
	}
	bat.setPropEnergyFull(info.EnergyFull)

	logger.Debugf("powerNow %vµA", info.PowerNow)
	bat.setPropPowerNow(info.PowerNow)

	logger.Debugf("voltageNow %vµV", info.VoltageNow)
	bat.setPropVoltageNow(info.VoltageNow)

	logger.Debugf("percentage %.1f%%", info.Percentage)
	bat.setPropPercentage(info.Percentage)

	logger.Debugf("capacity %.1f%%", info.Capacity)
	bat.setPropCapacity(info.Capacity)

	logger.Debug("status", info.Status)
	bat.setPropStatus(info.Status)

	logger.Debugf("timeToEmpty %v (%vs), timeToFull %v (%vs)",
		time.Duration(info.TimeToEmpty)*time.Second,
		info.TimeToEmpty,
		time.Duration(info.TimeToFull)*time.Second,
		info.TimeToFull)

	bat.setPropTimeToEmpty(info.TimeToEmpty)
	bat.setPropTimeToFull(info.TimeToFull)
}

func (bat *Battery) Refresh() {
	dev := bat.newDevice()
	if dev != nil {
		bat.refresh(dev)
		dev.Unref()
	} else {
		logger.Warningf("Refresh %v failed", bat.Name)
	}
}

func (bat *Battery) startLoopUpdate(d time.Duration) chan struct{} {
	done := make(chan struct{}, 1)
	go func() {
		ticker := time.NewTicker(d)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				bat.Refresh()
			case <-done:
				return
			}
		}
	}()
	return done
}

func (bat *Battery) resetUpdateInterval(d time.Duration) {
	if bat.exit != nil {
		close(bat.exit)
	}
	bat.exit = bat.startLoopUpdate(d)
}

func (bat *Battery) destroy() {
	if bat.exit != nil {
		close(bat.exit)
		bat.exit = nil
	}
}

// 仅用于调试
func (bat *Battery) Debug(cmd string) {
	dev := bat.newDevice()
	if dev != nil {
		defer dev.Unref()

		switch cmd {
		case "reset-update-interval1":
			bat.resetUpdateInterval(1 * time.Second)
		case "reset-update-interval3":
			bat.resetUpdateInterval(3 * time.Second)
		default:
			logger.Info("Command no support")
		}
	}
}
