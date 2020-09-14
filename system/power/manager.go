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
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	dbus "github.com/godbus/dbus"
	"pkg.deepin.io/dde/api/powersupply"
	"pkg.deepin.io/dde/api/powersupply/battery"
	gudev "pkg.deepin.io/gir/gudev-1.0"
	"pkg.deepin.io/lib/arch"
	"pkg.deepin.io/lib/dbusutil"
)

var noUEvent bool

func init() {
	if arch.Get() == arch.Sunway {
		noUEvent = true
	}
}

//go:generate dbusutil-gen -type Manager,Battery -import pkg.deepin.io/dde/api/powersupply/battery manager.go battery.go

// https://www.kernel.org/doc/Documentation/power/power_supply_class.txt
type Manager struct {
	service     *dbusutil.Service
	batteries   map[string]*Battery
	batteriesMu sync.Mutex
	ac          *AC
	gudevClient *gudev.Client

	// 电池是否电量低
	batteryLow bool
	// 初始化是否完成
	initDone bool

	// CPU操作接口
	cpus *CpuHandlers

	PropsMu      sync.RWMutex
	OnBattery    bool
	HasLidSwitch bool
	// battery display properties:
	HasBattery         bool
	BatteryPercentage  float64
	BatteryStatus      battery.Status
	BatteryTimeToEmpty uint64
	BatteryTimeToFull  uint64
	// 电池容量
	BatteryCapacity float64

	// 开启和关闭节能模式
	PowerSavingModeEnabled bool `prop:"access:rw"`

	// 自动切换节能模式，依据为是否插拔电源
	PowerSavingModeAuto bool `prop:"access:rw"`

	// 低电量时自动开启
	PowerSavingModeAutoWhenBatteryLow bool `prop:"access:rw"`

	// 开启节能模式时降低亮度的百分比值
	PowerSavingModeBrightnessDropPercent uint32 `prop:"access:rw"`

	// CPU频率调节模式，支持powersave和performance
	CpuGovernor string

	// CPU频率增强是否开启
	CpuBoost bool

	// 是否支持Boost
	IsBoostSupported bool

	// nolint
	methods *struct {
		GetBatteries   func() `out:"batteries"`
		Debug          func() `in:"cmd"`
		SetCpuGovernor func() `in:"governor"`
		SetCpuBoost    func() `in:"enabled"`
	}

	// nolint
	signals *struct {
		BatteryDisplayUpdate struct {
			timestamp int64
		}

		BatteryAdded struct {
			path dbus.ObjectPath
		}

		BatteryRemoved struct {
			path dbus.ObjectPath
		}

		LidClosed struct{}
		LidOpened struct{}
	}
}

func newManager(service *dbusutil.Service) (*Manager, error) {
	m := &Manager{
		service:           service,
		BatteryPercentage: 100,
		cpus:              NewCpuHandlers(),
	}
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

func newAC(manager *Manager, device *gudev.Device) *AC {
	sysfsPath := device.GetSysfsPath()
	return &AC{
		gudevClient: manager.gudevClient,
		sysfsPath:   sysfsPath,
	}
}

func (ac *AC) newDevice() *gudev.Device {
	return ac.gudevClient.QueryBySysfsPath(ac.sysfsPath)
}

func (m *Manager) refreshAC(ac *gudev.Device) { // 拔插电源时候触发
	online := ac.GetPropertyAsBoolean("POWER_SUPPLY_ONLINE")
	logger.Debug("ac online:", online)
	onBattery := !online

	m.PropsMu.Lock()
	m.setPropOnBattery(onBattery)
	m.PropsMu.Unlock()
	// 根据OnBattery的状态,修改节能模式
	m.updatePowerSavingMode()
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
		m.ac = newAC(m, ac)

		if noUEvent {
			go func() {
				c := time.Tick(2 * time.Second)
				for range c {
					err := m.RefreshMains()
					if err != nil {
						logger.Warning(err)
					}
				}
			}()
		}
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

	cfg := loadConfigSafe()
	// 将config.json中的配置完成初始化
	m.PowerSavingModeEnabled = cfg.PowerSavingModeEnabled                             // 开启和关闭节能模式
	m.PowerSavingModeAuto = cfg.PowerSavingModeAuto                                   // 自动切换节能模式，依据为是否插拔电源
	m.PowerSavingModeAutoWhenBatteryLow = cfg.PowerSavingModeAutoWhenBatteryLow       // 低电量时自动开启
	m.PowerSavingModeBrightnessDropPercent = cfg.PowerSavingModeBrightnessDropPercent // 开启节能模式时降低亮度的百分比值

	m.initAC(devices)
	m.initBatteries(devices)
	for _, dev := range devices {
		dev.Unref()
	}

	m.gudevClient.Connect("uevent", m.handleUEvent)
	m.initDone = true
	// init LMT config
	m.updatePowerSavingMode()

	var err error
	m.IsBoostSupported = m.cpus.IsBoostFileExist()
	m.CpuBoost, err = m.cpus.GetBoostEnabled()
	if err != nil {
		logger.Warning(err)
	}

	m.CpuGovernor, err = m.cpus.GetGovernor()
	if err != nil {
		logger.Warning(err)
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
				m.ac = newAC(m, device)
			} else if m.ac.sysfsPath != device.GetSysfsPath() {
				logger.Warning("found another AC", device.GetSysfsPath())
				return
			}

			// now m.ac != nil, and sysfsPath equal
			m.refreshAC(device)
			time.AfterFunc(1*time.Second, m.refreshBatteries)
			time.AfterFunc(3*time.Second, m.refreshBatteries)

		} else if powersupply.IsSystemBattery(device) {
			m.addAndExportBattery(device)
		}
	case "add":
		if powersupply.IsSystemBattery(device) {
			m.addAndExportBattery(device)
		}
		// ignore add mains

	case "remove":
		if powersupply.IsSystemBattery(device) {
			m.removeBattery(device)
		}
	}

}

func (m *Manager) initBatteries(devices []*gudev.Device) {
	m.batteries = make(map[string]*Battery)
	for _, dev := range devices {
		m.addBattery(dev)
	}
	logger.Debugf("initBatteries done %#v", m.batteries)
}

func (m *Manager) addAndExportBattery(dev *gudev.Device) {
	bat, added := m.addBattery(dev)
	if added {
		err := m.service.Export(bat.getObjPath(), bat)
		if err == nil {
			m.emitBatteryAdded(bat)
		} else {
			logger.Warning("failed to export battery:", err)
		}
	}
}

func (m *Manager) addBattery(dev *gudev.Device) (*Battery, bool) {
	logger.Debug("addBattery dev:", dev)
	if !powersupply.IsSystemBattery(dev) {
		return nil, false
	}

	sysfsPath := dev.GetSysfsPath()
	logger.Debug(sysfsPath)

	m.batteriesMu.Lock()
	bat, ok := m.batteries[sysfsPath]
	m.batteriesMu.Unlock()
	if ok {
		logger.Debugf("add battery failed , sysfsPath exists %q", sysfsPath)
		bat.Refresh()
		return bat, false
	}

	bat = newBattery(m, dev)
	if bat == nil {
		logger.Debugf("add batteries failed, sysfsPath %q, new battery failed", sysfsPath)
		return nil, false
	}

	m.batteriesMu.Lock()
	m.batteries[sysfsPath] = bat
	m.refreshBatteryDisplay()
	m.batteriesMu.Unlock()
	bat.setRefreshDoneCallback(m.refreshBatteryDisplay)
	return bat, true
}

// removeBattery remove the battery from Manager.batteries, and stop export it.
func (m *Manager) removeBattery(dev *gudev.Device) {
	sysfsPath := dev.GetSysfsPath()

	m.batteriesMu.Lock()
	bat, ok := m.batteries[sysfsPath]
	m.batteriesMu.Unlock()

	if ok {
		logger.Info("removeBattery", sysfsPath)
		m.batteriesMu.Lock()
		delete(m.batteries, sysfsPath)
		m.refreshBatteryDisplay()
		m.batteriesMu.Unlock()

		err := m.service.StopExport(bat)
		if err != nil {
			logger.Warning(err)
		}
		m.emitBatteryRemoved(bat)

		bat.destroy()
	} else {
		logger.Warning("removeBattery failed: invalid sysfsPath ", sysfsPath)
	}
}

func (m *Manager) emitBatteryAdded(bat *Battery) {
	err := m.service.Emit(m, "BatteryAdded", bat.getObjPath())
	if err != nil {
		logger.Warning(err)
	}
}

func (m *Manager) emitBatteryRemoved(bat *Battery) {
	err := m.service.Emit(m, "BatteryRemoved", bat.getObjPath())
	if err != nil {
		logger.Warning(err)
	}
}

func (m *Manager) destroy() {
	logger.Debug("destroy")
	m.batteriesMu.Lock()
	for _, bat := range m.batteries {
		bat.destroy()
	}
	m.batteries = nil
	m.batteriesMu.Unlock()

	if m.gudevClient != nil {
		m.gudevClient.Unref()
		m.gudevClient = nil
	}
}

const configFile = "/var/lib/dde-daemon/power/config.json"

type Config struct {
	PowerSavingModeEnabled               bool
	PowerSavingModeAuto                  bool
	PowerSavingModeAutoWhenBatteryLow    bool
	PowerSavingModeBrightnessDropPercent uint32
}

func loadConfig() (*Config, error) {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = json.Unmarshal(content, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func loadConfigSafe() *Config {
	cfg, err := loadConfig()
	if err != nil {
		// ignore not exist error
		if !os.IsNotExist(err) {
			logger.Warning(err)
		}
		return &Config{
			// default config
			PowerSavingModeAuto:                  true,
			PowerSavingModeEnabled:               false,
			PowerSavingModeAutoWhenBatteryLow:    false,
			PowerSavingModeBrightnessDropPercent: 20,
		}
	}
	// 新增字段后第一次启动时,缺少两个新增字段的json,导致亮度下降百分比字段默认为0,导致与默认值不符,需要处理
	// 低电量自动待机字段的默认值为false,不会导致错误影响
	// 正常情况下该字段范围为10-40,只有在该情况下会出现0的可能
	if cfg.PowerSavingModeBrightnessDropPercent == 0 {
		cfg.PowerSavingModeBrightnessDropPercent = 20
	}
	return cfg
}

func (m *Manager) saveConfig() error {
	logger.Debug("call saveConfig")

	var cfg Config
	m.PropsMu.RLock()
	cfg.PowerSavingModeAuto = m.PowerSavingModeAuto
	cfg.PowerSavingModeEnabled = m.PowerSavingModeEnabled
	cfg.PowerSavingModeAutoWhenBatteryLow = m.PowerSavingModeAutoWhenBatteryLow
	cfg.PowerSavingModeBrightnessDropPercent = m.PowerSavingModeBrightnessDropPercent
	m.PropsMu.RUnlock()

	dir := filepath.Dir(configFile)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	content, err := json.Marshal(&cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFile, content, 0644)
}
