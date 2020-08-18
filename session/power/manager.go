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
	"os"
	"sync"
	"syscall"

	dbus "github.com/godbus/dbus"
	"pkg.deepin.io/dde/daemon/common/dsync"
	"pkg.deepin.io/dde/daemon/session/common"
	gio "pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/gsprop"
)

//go:generate dbusutil-gen -type Manager manager.go
type Manager struct {
	service              *dbusutil.Service
	sessionSigLoop       *dbusutil.SignalLoop
	systemSigLoop        *dbusutil.SignalLoop
	syncConfig           *dsync.Config
	helper               *Helper
	settings             *gio.Settings
	warnLevelCountTicker *countTicker
	warnLevelConfig      *WarnLevelConfigManager
	submodules           map[string]submodule
	inhibitor            *sleepInhibitor
	inhibitFd            dbus.UnixFD

	PropsMu sync.RWMutex
	// 是否有盖子，一般笔记本电脑才有
	LidIsPresent bool
	// 是否使用电池, 接通电源时为 false, 使用电池时为 true
	OnBattery bool
	//是否使用Wayland
	UseWayland bool
	// 警告级别
	WarnLevel WarnLevel

	// 是否有环境光传感器
	HasAmbientLightSensor bool

	// dbusutil-gen: ignore-below
	// 电池是否可用，是否存在
	BatteryIsPresent map[string]bool
	// 电池电量百分比
	BatteryPercentage map[string]float64
	// 电池状态
	BatteryState map[string]uint32

	// 接通电源时，不做任何操作，到显示屏保的时间
	LinePowerScreensaverDelay gsprop.Int `prop:"access:rw"`
	// 接通电源时，不做任何操作，到关闭屏幕的时间
	LinePowerScreenBlackDelay gsprop.Int `prop:"access:rw"`
	// 接通电源时，不做任何操作，到睡眠的时间
	LinePowerSleepDelay gsprop.Int `prop:"access:rw"`

	// 使用电池时，不做任何操作，到显示屏保的时间
	BatteryScreensaverDelay gsprop.Int `prop:"access:rw"`
	// 使用电池时，不做任何操作，到关闭屏幕的时间
	BatteryScreenBlackDelay gsprop.Int `prop:"access:rw"`
	// 使用电池时，不做任何操作，到睡眠的时间
	BatterySleepDelay gsprop.Int `prop:"access:rw"`

	// 关闭屏幕前是否锁定
	ScreenBlackLock gsprop.Bool `prop:"access:rw"`
	// 睡眠前是否锁定
	SleepLock gsprop.Bool `prop:"access:rw"`

	// 废弃
	LidClosedSleep gsprop.Bool `prop:"access:rw"`

	// 接通电源时，笔记本电脑盖上盖子 待机（默认选择）、睡眠、关闭显示器、无任何操作
	LinePowerLidClosedAction gsprop.Enum `prop:"access:rw"`

	// 接通电源时，按下电源按钮 关机（默认选择）、待机、睡眠、关闭显示器、无任何操作
	LinePowerPressPowerBtnAction gsprop.Enum `prop:"access:rw"` // keybinding中监听power按键事件,获取gsettings的值

	// 使用电池时，笔记本电脑盖上盖子 待机（默认选择）、睡眠、关闭显示器、无任何操作
	BatteryLidClosedAction gsprop.Enum `prop:"access:rw"`

	// 使用电池时，按下电源按钮 关机（默认选择）、待机、睡眠、关闭显示器、无任何操作
	BatteryPressPowerBtnAction gsprop.Enum `prop:"access:rw"` // keybinding中监听power按键事件,获取gsettings的值

	// 接通电源时，不做任何操作，到自动锁屏的时间
	LinePowerLockDelay gsprop.Int `prop:"access:rw"`
	// 使用电池时，不做任何操作，到自动锁屏的时间
	BatteryLockDelay gsprop.Int `prop:"access:rw"`

	// 打开电量通知
	LowPowerNotifyEnable gsprop.Bool `prop:"access:rw"` // 开启后默认当电池仅剩余达到电量水平低时（默认15%）发出系统通知“电池电量低，请连接电源”；
	// 当电池仅剩余为设置低电量时（默认5%），发出系统通知“电池电量耗尽”，进入待机模式；

	// 电池低电量通知百分比
	LowPowerNotifyThreshold gsprop.Int `prop:"access:rw"` // 设置电量低提醒的阈值，可设置范围10%-25%，默认为20%

	// 自动待机电量百分比
	LowPowerAutoSleepThreshold gsprop.Int `prop:"access:rw"` // 设置电池电量进入待机模式（s3）的阈值，可设置范围为1%-9%，默认为5%（范围待确定）

	savingModeBrightnessDropPercent gsprop.Int // 用来接收和保存来自system power中降低的屏幕亮度值

	AmbientLightAdjustBrightness gsprop.Bool `prop:"access:rw"`
	ambientLightClaimed          bool
	lightLevelUnit               string
	lidSwitchState               uint
	sessionActive                bool

	// if prepare suspend, ignore idle off
	prepareSuspend       int
	prepareSuspendLocker sync.Mutex
}

func newManager(service *dbusutil.Service) (*Manager, error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	m := new(Manager)
	m.service = service
	sessionBus := service.Conn()
	m.sessionSigLoop = dbusutil.NewSignalLoop(sessionBus, 10)
	m.systemSigLoop = dbusutil.NewSignalLoop(systemBus, 10)
	m.inhibitFd = -1
	m.prepareSuspend = suspendStateUnknown

	m.syncConfig = dsync.NewConfig("power", &syncConfig{m: m}, m.sessionSigLoop, dbusPath, logger)

	helper, err := newHelper(systemBus, sessionBus)
	if err != nil {
		return nil, err
	}
	m.helper = helper

	m.settings = gio.NewSettings(gsSchemaPower)
	m.warnLevelConfig = NewWarnLevelConfigManager(m.settings)

	m.LinePowerScreensaverDelay.Bind(m.settings, settingKeyLinePowerScreensaverDelay)
	m.LinePowerScreenBlackDelay.Bind(m.settings, settingKeyLinePowerScreenBlackDelay)
	m.LinePowerSleepDelay.Bind(m.settings, settingKeyLinePowerSleepDelay)
	m.LinePowerLockDelay.Bind(m.settings, settingKeyLinePowerLockDelay)
	m.BatteryScreensaverDelay.Bind(m.settings, settingKeyBatteryScreensaverDelay)
	m.BatteryScreenBlackDelay.Bind(m.settings, settingKeyBatteryScreenBlackDelay)
	m.BatterySleepDelay.Bind(m.settings, settingKeyBatterySleepDelay)
	m.BatteryLockDelay.Bind(m.settings, settingKeyBatteryLockDelay)
	m.ScreenBlackLock.Bind(m.settings, settingKeyScreenBlackLock)
	m.SleepLock.Bind(m.settings, settingKeySleepLock)

	m.LinePowerLidClosedAction.Bind(m.settings, settingKeyLinePowerLidClosedAction)
	m.LinePowerPressPowerBtnAction.Bind(m.settings, settingKeyLinePowerPressPowerBtnAction)
	m.BatteryLidClosedAction.Bind(m.settings, settingKeyBatteryLidClosedAction)
	m.BatteryPressPowerBtnAction.Bind(m.settings, settingKeyBatteryPressPowerBtnAction)
	m.LowPowerNotifyEnable.Bind(m.settings, settingKeyLowPowerNotifyEnable)
	m.LowPowerNotifyThreshold.Bind(m.settings, settingKeyLowPowerNotifyThreshold)
	m.LowPowerAutoSleepThreshold.Bind(m.settings, settingKeyLowPowerAutoSleepThreshold)
	m.savingModeBrightnessDropPercent.Bind(m.settings, settingKeyBrightnessDropPercent)
	m.initGSettingsConnectChanged()
	m.AmbientLightAdjustBrightness.Bind(m.settings,
		settingKeyAmbientLightAdjuestBrightness)

	power := m.helper.Power
	err = common.ActivateSysDaemonService(power.ServiceName_())
	if err != nil {
		logger.Warning(err)
	}

	m.LidIsPresent, err = power.HasLidSwitch().Get(0)
	if err != nil {
		logger.Warning(err)
	}

	m.OnBattery, err = power.OnBattery().Get(0)
	if err != nil {
		logger.Warning(err)
	}

	logger.Info("LidIsPresent", m.LidIsPresent)
	m.HasAmbientLightSensor, _ = helper.SensorProxy.HasAmbientLight().Get(0)
	logger.Debug("HasAmbientLightSensor:", m.HasAmbientLightSensor)
	if m.HasAmbientLightSensor {
		m.lightLevelUnit, _ = helper.SensorProxy.LightLevelUnit().Get(0)
	}

	m.sessionActive, _ = helper.SessionWatcher.IsActive().Get(0)

	// init battery display
	m.BatteryIsPresent = make(map[string]bool)
	m.BatteryPercentage = make(map[string]float64)
	m.BatteryState = make(map[string]uint32)

	return m, nil
}

func (m *Manager) init() {
	m.claimOrReleaseAmbientLight()
	m.sessionSigLoop.Start()
	m.systemSigLoop.Start()

	if len(os.Getenv("WAYLAND_DISPLAY")) != 0 {
		m.UseWayland = true
	} else {
		m.UseWayland = false
	}

	m.helper.initSignalExt(m.systemSigLoop, m.sessionSigLoop)

	// init sleep inhibitor
	m.inhibitor = newSleepInhibitor(m.helper.LoginManager)
	m.inhibitor.OnBeforeSuspend = m.handleBeforeSuspend
	m.inhibitor.OnWakeup = m.handleWakeup
	err := m.inhibitor.block()
	if err != nil {
		logger.Warning(err)
	}

	m.handleBatteryDisplayUpdate()
	power := m.helper.Power
	_, err = power.ConnectBatteryDisplayUpdate(func(timestamp int64) {
		logger.Debug("BatteryDisplayUpdate", timestamp)
		m.handleBatteryDisplayUpdate()
	})
	if err != nil {
		logger.Warning(err)
	}

	err = m.helper.SensorProxy.LightLevel().ConnectChanged(func(hasValue bool, value float64) {
		if !hasValue {
			return
		}
		m.handleLightLevelChanged(value)
	})
	if err != nil {
		logger.Warning(err)
	}

	_, err = m.helper.SysDBusDaemon.ConnectNameOwnerChanged(
		func(name string, oldOwner string, newOwner string) {
			serviceName := m.helper.SensorProxy.ServiceName_()
			if name == serviceName && newOwner != "" {
				logger.Debug("sensorProxy restarted")
				hasSensor, _ := m.helper.SensorProxy.HasAmbientLight().Get(0)
				var lightLevelUnit string
				if hasSensor {
					lightLevelUnit, _ = m.helper.SensorProxy.LightLevelUnit().Get(0)
				}

				m.PropsMu.Lock()
				m.setPropHasAmbientLightSensor(hasSensor)
				m.ambientLightClaimed = false
				m.lightLevelUnit = lightLevelUnit
				m.PropsMu.Unlock()

				m.claimOrReleaseAmbientLight()
			}
		})
	if err != nil {
		logger.Warning(err)
	}

	err = m.helper.SessionWatcher.IsActive().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}

		m.PropsMu.Lock()
		m.sessionActive = value
		m.PropsMu.Unlock()

		logger.Debug("session active changed to:", value)
		m.claimOrReleaseAmbientLight()
	})
	if err != nil {
		logger.Warning(err)
	}

	m.warnLevelConfig.setChangeCallback(m.handleBatteryDisplayUpdate)

	m.initPowerModule()

	m.initOnBatteryChangedHandler()
	m.initSubmodules()
	m.startSubmodules()
	m.inhibitLogind()
}

func (m *Manager) initPowerModule() {
	init := m.settings.GetBoolean(settingKeyPowerModuleInitialized)
	if !init {
		// TODO: 也许有更好的判断台式机的方法
		power := m.helper.Power
		hasBattery, err := power.HasBattery().Get(0)
		if err != nil {
			logger.Warning(err)
		} else {
			if !hasBattery {
				// 无电池，判断为台式机, 设置待机为 从不
				m.LinePowerSleepDelay.Set(0)
				m.BatterySleepDelay.Set(0)
			}
		}
		m.settings.SetBoolean(settingKeyPowerModuleInitialized, true)
	}
}

func (m *Manager) isX11SessionActive() (bool, error) {
	return m.helper.SessionWatcher.IsX11SessionActive(0)
}

func (m *Manager) destroy() {
	m.destroySubmodules()
	m.releaseAmbientLight()
	m.permitLogind()

	if m.helper != nil {
		m.helper.Destroy()
		m.helper = nil
	}

	if m.inhibitor != nil {
		err := m.inhibitor.unblock()
		if err != nil {
			logger.Warning(err)
		}
		m.inhibitor = nil
	}

	m.systemSigLoop.Stop()
	m.sessionSigLoop.Stop()
	m.syncConfig.Destroy()
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}

func (m *Manager) Reset() *dbus.Error {
	logger.Debug("Reset settings")

	var settingKeys = []string{
		settingKeyLinePowerScreenBlackDelay,
		settingKeyLinePowerSleepDelay,
		settingKeyLinePowerLockDelay,
		settingKeyLinePowerLidClosedAction,
		settingKeyLinePowerPressPowerBtnAction,

		settingKeyBatteryScreenBlackDelay,
		settingKeyBatterySleepDelay,
		settingKeyBatteryLockDelay,
		settingKeyBatteryLidClosedAction,
		settingKeyBatteryPressPowerBtnAction,

		settingKeyScreenBlackLock,
		settingKeySleepLock,
		settingKeyPowerButtonPressedExec,

		settingKeyLowPowerNotifyEnable,
		settingKeyLowPowerNotifyThreshold,
		settingKeyLowPowerAutoSleepThreshold,
		settingKeyBrightnessDropPercent,
	}
	for _, key := range settingKeys {
		logger.Debug("reset setting", key)
		m.settings.Reset(key)
	}
	return nil
}

func (m *Manager) inhibitLogind() {
	fd, err := m.helper.LoginManager.Inhibit(0,
		"handle-power-key:handle-lid-switch", dbusServiceName,
		"handling key press and lid switch close", "block")
	logger.Debug("inhibitLogind fd:", fd)
	if err != nil {
		logger.Warning(err)
		return
	}
	m.inhibitFd = fd
}

func (m *Manager) permitLogind() {
	if m.inhibitFd != -1 {
		err := syscall.Close(int(m.inhibitFd))
		if err != nil {
			logger.Warning("failed to close inhibitFd:", err)
		}
		m.inhibitFd = -1
	}
}
