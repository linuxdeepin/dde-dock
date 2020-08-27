/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package power

import (
	"strings"
	"sync"
	"time"

	"path/filepath"

	"pkg.deepin.io/dde/api/powersupply/battery"
	gudev "pkg.deepin.io/gir/gudev-1.0"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

type Battery struct {
	service *dbusutil.Service
	exit    chan struct{}
	mutex   sync.Mutex

	gudevClient       *gudev.Client
	changedProperties []string

	PropsMu   sync.RWMutex
	SysfsPath string
	IsPresent bool

	Manufacturer string
	ModelName    string
	SerialNumber string
	Name         string
	Technology   string

	Energy           float64
	EnergyFull       float64
	EnergyFullDesign float64
	EnergyRate       float64

	Voltage     float64
	Percentage  float64
	Capacity    float64
	Status      battery.Status
	TimeToEmpty uint64
	TimeToFull  uint64
	UpdateTime  int64

	batteryHistory []float64

	refreshDone func()

	methods *struct {
		Debug func() `in:"cmd"`
	}
}

const (
	checkTimeSliceCount = 2
	checkTimeFloatRange = 120
)

func newBattery(manager *Manager, device *gudev.Device) *Battery {
	sysfsPath := device.GetSysfsPath()
	logger.Debugf("NewBattery sysfsPath: %q", sysfsPath)
	if manager == nil || manager.gudevClient == nil ||
		device == nil {
		return nil
	}
	bat := &Battery{
		service:     manager.service,
		gudevClient: manager.gudevClient,
		SysfsPath:   sysfsPath,
	}
	ok := bat.refresh(device)
	if !ok {
		return nil
	}
	bat.resetUpdateInterval(60 * time.Second)
	return bat
}

const (
	batteryDBusInterface = dbusInterface + ".Battery"
)

func (*Battery) GetInterfaceName() string {
	return batteryDBusInterface
}

func (bat *Battery) getObjPath() dbus.ObjectPath {
	return dbus.ObjectPath(dbusPath + "/battery_" + getValidName(filepath.Base(bat.SysfsPath)))
}

func getValidName(n string) string {
	// dbus objpath 0-9 a-z A-Z _
	n = strings.Replace(n, "-", "_x0", -1)
	n = strings.Replace(n, ".", "_x1", -1)
	n = strings.Replace(n, ":", "_x2", -1)
	return n
}

func checkTimeStabilized(s []uint64, t uint64) bool {
	if len(s) <= checkTimeSliceCount {
		// 记录少于 3
		return false
	}

	// 循环最后三次记录
	length := len(s)
	for i := 1; i <= checkTimeSliceCount; i++ {
		if t+checkTimeFloatRange <= s[length-i] {
			// 所处记录大于 t+允许浮动值
			return false
		}
		if t-checkTimeFloatRange >= s[length-i] {
			// 所处记录小于 t-允许浮动值
			return false
		}
	}

	return true
}

func (bat *Battery) setRefreshDoneCallback(fn func()) {
	bat.refreshDone = fn
}

func (bat *Battery) newDevice() *gudev.Device {
	return bat.gudevClient.QueryBySysfsPath(bat.SysfsPath)
}

func (bat *Battery) notifyChange(propNames ...string) {
	bat.changedProperties = append(bat.changedProperties, propNames...)
}

func (bat *Battery) refresh(dev *gudev.Device) (ok bool) {
	endDelay := bat.service.DelayEmitPropertyChanged(bat)
	batInfo := battery.GetBatteryInfo(dev)
	if batInfo == nil {
		return
	}

	setTimeToFull := true
	if batInfo.Status == battery.StatusCharging {
		/*
		 * bug: https://pms.uniontech.com/zentao/bug-view-15187.html
		 * \desc Because the time to full provide by Upower module is not stable
		 * when just started charging. There are now two way solve that:
		 * 1. take the average value of the times, but the frequency of refresh battry is 60 seconds.
		 * If add more time tickers may introducing new bugs because the refresh() method can used by
		 * other method. Just like: https://pms.uniontech.com/zentao/bug-view-37382.html
		 * 2. just show time to full on the next refresh, it will be 60s after charging.
		 * \warn Add frequency of refresh battry will take a lot of overhead in cpu
		 * refresh battry every 1s take 5% cpu performence on arm laptop
		 */
		if bat.Status == battery.StatusDischarging {
			logger.Debug("Just started charging")
			setTimeToFull = false
		}
	}

	bat._refresh(batInfo, setTimeToFull)
	if endDelay != nil {
		err := endDelay()
		if err != nil {
			logger.Warning(err)
		}
	}
	ok = true
	return
}

func (bat *Battery) _refresh(info *battery.BatteryInfo, setTimeToFull bool) {
	logger.Debug("Refresh", bat.Name)
	isPresent := true
	var updateTime int64
	if info == nil {
		isPresent = false
		info = &battery.BatteryInfo{}
	} else {
		now := time.Now()
		updateTime = now.Unix()
		logger.Debugf("now %v updateTime %v", now, updateTime)
	}

	logger.Debug("Name", info.Name)
	logger.Debug("Technology", info.Technology)
	logger.Debug("Manufacturer", info.Manufacturer)
	logger.Debug("ModelName", info.ModelName)
	logger.Debug("SerialNumber", info.SerialNumber)
	logger.Debugf("energy %v", info.Energy)
	logger.Debugf("energyFull %v", info.EnergyFull)
	logger.Debugf("EnergyFullDesign %v", info.EnergyFullDesign)
	logger.Debugf("EnergyRate %v", info.EnergyRate)
	logger.Debugf("voltage %v", info.Voltage)
	logger.Debugf("percentage %.4f%%", info.Percentage)
	logger.Debugf("capacity %.4f%%", info.Capacity)
	logger.Debug("status", info.Status)
	logger.Debugf("timeToEmpty %v (%vs), timeToFull %v (%vs)",
		time.Duration(info.TimeToEmpty)*time.Second,
		info.TimeToEmpty,
		time.Duration(info.TimeToFull)*time.Second,
		info.TimeToFull)

	/* lie to full */
	bat.appendToHistory(info.Percentage)
	if info.Percentage > 97.0 && bat.getHistoryLength() >= 10 && bat.calcHistoryVariance() < 0.3 {
		logger.Debugf("fake 100 : true percentage %.4f%% variance %.4f%%", info.Percentage, bat.calcHistoryVariance())
		info.Percentage = 100.0
		info.TimeToFull = 0
	}

	bat.PropsMu.Lock()
	bat.setPropIsPresent(isPresent)
	bat.setPropUpdateTime(updateTime)
	bat.setPropName(info.Name)
	bat.setPropTechnology(info.Technology)
	bat.setPropManufacturer(info.Manufacturer)
	bat.setPropModelName(info.ModelName)
	bat.setPropSerialNumber(info.SerialNumber)
	bat.setPropEnergy(info.Energy)
	bat.setPropEnergyFull(info.EnergyFull)
	bat.setPropEnergyFullDesign(info.EnergyFullDesign)
	bat.setPropEnergyRate(info.EnergyRate)
	bat.setPropVoltage(info.Voltage)
	bat.setPropPercentage(info.Percentage)
	bat.setPropCapacity(info.Capacity)
	bat.setPropStatus(info.Status)
	bat.setPropTimeToEmpty(info.TimeToEmpty)
	if setTimeToFull {
		bat.setPropTimeToFull(info.TimeToFull)
	} else {
		bat.setPropTimeToFull(0)
	}
	bat.PropsMu.Unlock()

	logger.Debugf("Refresh %v done", bat.Name)
	if bat.refreshDone != nil {
		bat.refreshDone()
	}
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
			case _, ok := <-ticker.C:
				if !ok {
					logger.Error("Invalid ticker event")
					return
				}

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
