/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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
	"errors"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/util/wm/ewmh"
	"pkg.deepin.io/lib/gsettings"
	"pkg.deepin.io/lib/procfs"
)

const submodulePSP = "PowerSavePlan"

func init() {
	submoduleList = append(submoduleList, newPowerSavePlan)
}

type powerSavePlan struct {
	manager            *Manager
	screenSaverTimeout int32
	metaTasks          metaTasks
	tasks              delayedTasks
	// key output name, value old brightness
	oldBrightnessTable map[string]float64
	mu                 sync.Mutex
	screensaverRunning bool

	atomNetWMStateFullscreen    x.Atom
	atomNetWMStateFocused       x.Atom
	fullscreenWorkaroundAppList []string
	enableSaveModeTime          time.Time
	setBrightnessTime           time.Time
}

func newPowerSavePlan(manager *Manager) (string, submodule, error) {
	p := new(powerSavePlan)
	p.manager = manager

	conn := manager.helper.xConn
	var err error
	p.atomNetWMStateFullscreen, err = conn.GetAtom("_NET_WM_STATE_FULLSCREEN")
	if err != nil {
		return submodulePSP, nil, err
	}
	p.atomNetWMStateFocused, err = conn.GetAtom("_NET_WM_STATE_FOCUSED")
	if err != nil {
		return submodulePSP, nil, err
	}

	p.fullscreenWorkaroundAppList = manager.settings.GetStrv(
		"fullscreen-workaround-app-list")
	return submodulePSP, p, nil
}

// 监听 GSettings 值改变, 更新节电计划
func (psp *powerSavePlan) initSettingsChangedHandler() {
	m := psp.manager
	gsettings.ConnectChanged(gsSchemaPower, "*", func(key string) {
		logger.Debug("setting changed", key)
		switch key {
		case settingKeyLinePowerScreensaverDelay,
			settingKeyLinePowerScreenBlackDelay,
			settingKeyLinePowerLockDelay,
			settingKeyLinePowerSleepDelay:
			if !m.OnBattery {
				logger.Debug("Change OnLinePower plan")
				psp.OnLinePower()
			}

		case settingKeyBatteryScreensaverDelay,
			settingKeyBatteryScreenBlackDelay,
			settingKeyBatteryLockDelay,
			settingKeyBatterySleepDelay:
			if m.OnBattery {
				logger.Debug("Change OnBattery plan")
				psp.OnBattery()
			}

		case settingKeyAmbientLightAdjuestBrightness:
			psp.manager.claimOrReleaseAmbientLight()
		}
	})
}

func (psp *powerSavePlan) OnBattery() {
	logger.Debug("Use OnBattery plan")
	m := psp.manager
	psp.Update(m.BatteryScreensaverDelay.Get(), m.BatteryLockDelay.Get(),
		m.BatteryScreenBlackDelay.Get(),
		m.BatterySleepDelay.Get())
}

func (psp *powerSavePlan) OnLinePower() {
	logger.Debug("Use OnLinePower plan")
	m := psp.manager
	psp.Update(m.LinePowerScreensaverDelay.Get(), m.LinePowerLockDelay.Get(),
		m.LinePowerScreenBlackDelay.Get(),
		m.LinePowerSleepDelay.Get())
}

func (psp *powerSavePlan) Reset() {
	m := psp.manager
	logger.Debug("OnBattery:", m.OnBattery)
	if m.OnBattery {
		psp.OnBattery()
	} else {
		psp.OnLinePower()
	}
}

func (psp *powerSavePlan) Start() error {
	psp.Reset()
	psp.initSettingsChangedHandler()

	helper := psp.manager.helper
	power := helper.Power
	display := helper.Display
	screenSaver := helper.ScreenSaver
	//OnBattery changed will effect current PowerSavePlan
	err := power.OnBattery().ConnectChanged(func(hasValue bool, value bool) {
		psp.Reset()
	})
	if err != nil {
		logger.Warning("failed to connectChanged OnBattery:", err)
	}
	err = power.PowerSavingModeEnabled().ConnectChanged(psp.handlePowerSavingModeChanged)
	if err != nil {
		logger.Warning("failed to connectChanged PowerSavingModeEnabled:", err)
	}
	err = power.PowerSavingModeBrightnessDropPercent().ConnectChanged(psp.handlePowerSavingModeBrightnessDropPercentChanged) // 监听自动降低亮度的属性的改变
	if err != nil {
		logger.Warning("failed to connectChanged PowerSavingModeBrightnessDropPercent:", err)
	}
	_, err = screenSaver.ConnectIdleOn(psp.HandleIdleOn)
	if err != nil {
		logger.Warning("failed to ConnectIdleOn:", err)
	}
	_, err = screenSaver.ConnectIdleOff(psp.HandleIdleOff)
	if err != nil {
		logger.Warning("failed to ConnectIdleOff:", err)
	}
	err = display.Brightness().ConnectChanged(psp.saveSetBrightnessTime)
	if err != nil {
		logger.Warning("failed to connectChanged Brightness:", err)
	}
	return nil
}

// 节能模式降低亮度的比例,并降低亮度
func (psp *powerSavePlan) handlePowerSavingModeBrightnessDropPercentChanged(hasValue bool, lowerValue uint32) {
	if !hasValue {
		return
	}
	logger.Debug("power saving mode lower brightness changed to", lowerValue)
	newLowerBrightnessScale := float64(lowerValue)
	psp.manager.PropsMu.RLock()
	hasLightSensor := psp.manager.HasAmbientLightSensor
	psp.manager.PropsMu.RUnlock()

	if hasLightSensor && psp.manager.AmbientLightAdjustBrightness.Get() {
		return
	}
	oldLowerBrightnessScale := float64(psp.manager.savingModeBrightnessDropPercent.Get()) // 保存之前的亮度下降值
	psp.manager.savingModeBrightnessDropPercent.Set(int32(lowerValue))
	savingModeEnable, err := psp.manager.helper.Power.PowerSavingModeEnabled().Get(0)
	if err != nil {
		logger.Error("get current power savingMode state error : ", err)
	}

	brightnessTable, err := psp.manager.helper.Display.GetBrightness(0)
	if err != nil {
		logger.Warning(err)
		return
	}

	if savingModeEnable {
		// adjust brightness by lowerBrightnessScale
		var unSavingModeBrightness float64
		// 判断亮度修改是手动调节亮度还是调节了节能选项
		brightnessChangedByManual := psp.isBrightnessChangedByManual()
		lowerBrightnessScale := 1 - newLowerBrightnessScale/100
		for key, value := range brightnessTable { // 反求未节能时的亮度
			if !brightnessChangedByManual {
				unSavingModeBrightness = value / (1 - oldLowerBrightnessScale/100)
				if unSavingModeBrightness > 1 {
					unSavingModeBrightness = 1
				}
			} else { // 在人为直接修改亮度之后,再使用节能选项调节,则不会按照比例反求亮度
				unSavingModeBrightness = value
			}

			value = unSavingModeBrightness * lowerBrightnessScale
			if value < 0.1 {
				value = 0.1
			}
			brightnessTable[key] = value
			psp.mu.Lock()
			psp.enableSaveModeTime = time.Now()
			psp.mu.Unlock()
		}
	} else {
		//else中(非节能状态下的调节)不需要做响应,需要降低亮度的预设值在之前已经保存了
		return
	}
	psp.manager.setAndSaveDisplayBrightness(brightnessTable)
}

//节能模式变化后的亮度修改
func (psp *powerSavePlan) handlePowerSavingModeChanged(hasValue bool, enabled bool) {
	const (
		multiLevelAdjustmentScale     = 0.2 // 分级调节时，默认按照20%亮度调节值调整亮度，并在亮度显示的设置分级中，归到所在分级
		multiLevelAdjustmentThreshold = 100 // 分级调节判断阈值，最大亮度值小于该值且不为0时，调节方式为分级调节
	)
	if !hasValue {
		return
	}
	logger.Debug("power saving mode enabled changed to", enabled)

	psp.manager.PropsMu.RLock()
	hasLightSensor := psp.manager.HasAmbientLightSensor
	psp.manager.PropsMu.RUnlock()

	if hasLightSensor && psp.manager.AmbientLightAdjustBrightness.Get() {
		return
	}

	brightnessTable, err := psp.manager.helper.Display.GetBrightness(0)
	if err != nil {
		logger.Warning(err)
		return
	}
	maxBacklightBrightness, err := psp.manager.helper.Display.MaxBacklightBrightness().Get(0)
	if err != nil {
		logger.Warning(err)
	}
	// 判断亮度调节方式是分级调节还是百分比滑动：最大亮度小于100且最大亮度不为0时，为分级调节
	isMultiLevelAdjustment := maxBacklightBrightness < multiLevelAdjustmentThreshold && maxBacklightBrightness != 0
	// 判断亮度修改是手动调节亮度还是调节了节能选项
	brightnessChangedByManual := psp.isBrightnessChangedByManual()
	lowerBrightnessScale := 1 - float64(psp.manager.savingModeBrightnessDropPercent.Get())/100
	oneStepValue := 1 / float64(maxBacklightBrightness)
	numSteps := math.Round(float64(maxBacklightBrightness) * multiLevelAdjustmentScale)
	if enabled {
		// reduce brightness when enabled saveMode
		for key, value := range brightnessTable {
			if isMultiLevelAdjustment {
				// 分级调节,减去需要降低的亮度
				value -= oneStepValue * numSteps
				if value < oneStepValue {
					value = oneStepValue
				}
			} else {
				// 非分级调节
				value *= lowerBrightnessScale
			}
			if value < 0.1 {
				value = 0.1
			}
			brightnessTable[key] = value
		}
		psp.mu.Lock()
		psp.enableSaveModeTime = time.Now()
		psp.mu.Unlock()
	} else {
		if !brightnessChangedByManual {
			logger.Debug("not manual adjust brightness")
			// increase brightness when disabled saveMode
			for key, value := range brightnessTable {
				if isMultiLevelAdjustment {
					// 分级调节,加上需要提升的亮度
					value += oneStepValue * numSteps
				} else {
					// 非分级调节
					value /= lowerBrightnessScale
				}
				if value > 1 {
					value = 1
				}
				brightnessTable[key] = value
			}
		} else {
			logger.Debug("manual adjust brightness")
			return
			// brightnessTable不变
		}
	}
	psp.manager.setAndSaveDisplayBrightness(brightnessTable)
}

// 取消之前的任务
func (psp *powerSavePlan) interruptTasks() {
	psp.tasks.CancelAll()
	psp.tasks.Wait(10*time.Millisecond, 200)
	psp.tasks = nil
}

func (psp *powerSavePlan) Destroy() {
	psp.interruptTasks()
}

func (psp *powerSavePlan) addTaskNoLock(t *delayedTask) {
	psp.tasks = append(psp.tasks, t)
}

func (psp *powerSavePlan) addTask(t *delayedTask) {
	psp.mu.Lock()
	psp.addTaskNoLock(t)
	psp.mu.Unlock()
}

func (psp *powerSavePlan) isBrightnessChangedByManual() bool {
	const interval = 2 * time.Second
	psp.mu.Lock()
	defer psp.mu.Unlock()

	return psp.setBrightnessTime.Sub(psp.enableSaveModeTime) > interval
}

type metaTask struct {
	delay     int32
	realDelay time.Duration
	name      string
	fn        func()
}

type metaTasks []metaTask

func (mts metaTasks) min() int32 {
	if len(mts) == 0 {
		return 0
	}

	min := mts[0].delay
	for _, t := range mts[1:] {
		if t.delay < min {
			min = t.delay
		}
	}
	return min
}

func (mts metaTasks) setRealDelay(min int32) {
	if min == 0 {
		return
	}
	for idx := range mts {
		t := &mts[idx]
		nSecs := t.delay - min
		if nSecs == 0 {
			t.realDelay = 1 * time.Millisecond
		} else {
			t.realDelay = time.Second * time.Duration(nSecs)
		}
	}
}

func (psp *powerSavePlan) Update(screenSaverStartDelay, lockDelay,
	screenBlackDelay, sleepDelay int32) {
	psp.mu.Lock()
	defer psp.mu.Unlock()

	psp.interruptTasks()
	logger.Debugf("update(screenSaverStartDelay=%vs, lockDelay=%vs,"+
		" screenBlackDelay=%vs, sleepDelay=%vs)",
		screenSaverStartDelay, lockDelay, screenBlackDelay, sleepDelay)

	tasks := make(metaTasks, 0, 4)
	if screenSaverStartDelay > 0 {
		tasks = append(tasks, metaTask{
			name:  "screenSaverStart",
			delay: screenSaverStartDelay,
			fn:    psp.startScreensaver,
		})
	}

	if lockDelay > 0 {
		tasks = append(tasks, metaTask{
			name:  "lock",
			delay: lockDelay,
			fn:    psp.lock,
		})
	}

	if screenBlackDelay > 0 {
		tasks = append(tasks, metaTask{
			name:  "screenBlack",
			delay: screenBlackDelay,
			fn:    psp.screenBlack,
		})
	}

	if sleepDelay > 0 {
		tasks = append(tasks, metaTask{
			name:  "sleep",
			delay: sleepDelay,
			fn:    psp.makeSystemSleep,
		})
	}

	min := tasks.min()
	tasks.setRealDelay(min)
	err := psp.setScreenSaverTimeout(min)
	if err != nil {
		logger.Warning("failed to set screen saver timeout:", err)
	}

	psp.metaTasks = tasks
}

func (psp *powerSavePlan) setScreenSaverTimeout(seconds int32) error {
	psp.screenSaverTimeout = seconds
	logger.Debugf("set ScreenSaver timeout to %d", seconds)
	err := psp.manager.helper.ScreenSaver.SetTimeout(0, uint32(seconds), 0, false)
	if err != nil {
		logger.Warningf("failed to set ScreenSaver timeout %d: %v", seconds, err)
	}
	return err
}

func (psp *powerSavePlan) saveCurrentBrightness() error {
	if psp.oldBrightnessTable == nil {
		var err error
		psp.oldBrightnessTable, err = psp.manager.helper.Display.Brightness().Get(0)
		if err != nil {
			return err
		}
		logger.Info("saveCurrentBrightness", psp.oldBrightnessTable)
		return nil
	}

	return errors.New("oldBrightnessTable is not nil")
}

func (psp *powerSavePlan) resetBrightness() {
	if psp.oldBrightnessTable != nil {
		logger.Debug("Reset all outputs brightness")
		psp.manager.setDisplayBrightness(psp.oldBrightnessTable)
		psp.oldBrightnessTable = nil
	}
}

func (psp *powerSavePlan) startScreensaver() {
	if os.Getenv("DESKTOP_CAN_SCREENSAVER") == "N" {
		logger.Info("do not start screensaver, env DESKTOP_CAN_SCREENSAVER == N")
		return
	}

	startScreensaver()
	psp.screensaverRunning = true
}

func (psp *powerSavePlan) stopScreensaver() {
	if !psp.screensaverRunning {
		return
	}
	stopScreensaver()
	psp.screensaverRunning = false
}

func (psp *powerSavePlan) makeSystemSleep() {
	logger.Info("sleep")
	psp.stopScreensaver()
	//psp.manager.setDPMSModeOn()
	//psp.resetBrightness()
	psp.manager.doSuspend()
}

func (psp *powerSavePlan) lock() {
	psp.manager.doLock(true)
}

// 降低显示器亮度，最终关闭显示器
func (psp *powerSavePlan) screenBlack() {
	manager := psp.manager
	logger.Info("Start screen black")

	adjustBrightnessEnabled := manager.settings.GetBoolean(settingKeyAdjustBrightnessEnabled)

	if adjustBrightnessEnabled {
		err := psp.saveCurrentBrightness()
		if err != nil {
			adjustBrightnessEnabled = false
			logger.Warning(err)
		} else {
			// half black
			brightnessTable := make(map[string]float64)
			brightnessRatio := 0.5
			logger.Debug("brightnessRatio:", brightnessRatio)
			for output, oldBrightness := range psp.oldBrightnessTable {
				brightnessTable[output] = oldBrightness * brightnessRatio
			}
			manager.setDisplayBrightness(brightnessTable)
		}
	} else {
		logger.Debug("adjust brightness disabled")
	}

	// full black
	const fullBlackTime = 5000 * time.Millisecond
	taskF := newDelayedTask("screenFullBlack", fullBlackTime, func() {
		psp.stopScreensaver()
		logger.Info("Screen full black")
		if manager.ScreenBlackLock.Get() {
			manager.lockWaitShow(5*time.Second, true)
		}

		if adjustBrightnessEnabled {
			// set min brightness for all outputs
			brightnessTable := make(map[string]float64)
			for output := range psp.oldBrightnessTable {
				brightnessTable[output] = 0.02
			}
			manager.setDisplayBrightness(brightnessTable)
		}
		manager.setDPMSModeOff()

	})
	psp.addTask(taskF)
}

func (psp *powerSavePlan) shouldPreventIdle() (bool, error) {
	conn := psp.manager.helper.xConn
	activeWin, err := ewmh.GetActiveWindow(conn).Reply(conn)
	if err != nil {
		return false, err
	}

	isFullscreenAndFocused, err := psp.isWindowFullScreenAndFocused(activeWin)
	if err != nil {
		return false, err
	}

	if !isFullscreenAndFocused {
		return false, nil
	}

	pid, err := ewmh.GetWMPid(conn, activeWin).Reply(conn)
	if err != nil {
		return false, err
	}

	p := procfs.Process(pid)
	cmdline, err := p.Cmdline()
	if err != nil {
		return false, err
	}

	for _, arg := range cmdline {
		for _, app := range psp.fullscreenWorkaroundAppList {
			if strings.Contains(arg, app) {
				logger.Debugf("match %q", app)
				return true, nil
			}
		}
	}
	return false, nil
}

// 开始 Idle
func (psp *powerSavePlan) HandleIdleOn() {
	psp.mu.Lock()
	defer psp.mu.Unlock()

	if psp.manager.shouldIgnoreIdleOn() {
		logger.Info("HandleIdleOn : IGNORE =========")
		return
	}

	if isActive, err := psp.manager.isX11SessionActive(); err != nil {
		logger.Warning(err)
		return
	} else if !isActive {
		logger.Info("X11 session is inactive, don't HandleIdleOn")
		return
	}

	// check window
	preventIdle, err := psp.shouldPreventIdle()
	if err != nil {
		logger.Warning(err)
	}
	if preventIdle {
		logger.Debug("prevent idle")
		err := psp.manager.helper.ScreenSaver.SimulateUserActivity(0)
		if err != nil {
			logger.Warning(err)
		}
		return
	}

	logger.Info("HandleIdleOn")

	for _, t := range psp.metaTasks {
		logger.Debugf("do %s after %v", t.name, t.realDelay)
		task := newDelayedTask(t.name, t.realDelay, t.fn)
		psp.addTaskNoLock(task)
	}

	_, err = os.Stat("/etc/deepin/no_suspend")
	if err == nil {
		if psp.manager.ScreenBlackLock.Get() {
			//m.setDPMSModeOn()
			//m.lockWaitShow(4 * time.Second)
			psp.manager.doLock(true)
			time.Sleep(time.Millisecond * 500)
		}
	}
}

// 结束 Idle
func (psp *powerSavePlan) HandleIdleOff() {
	psp.mu.Lock()
	defer psp.mu.Unlock()

	if psp.manager.shouldIgnoreIdleOff() {
		psp.manager.setPrepareSuspend(suspendStateFinish)
		logger.Info("HandleIdleOff : IGNORE =========")
		return
	}

	psp.manager.setPrepareSuspend(suspendStateFinish)
	logger.Info("HandleIdleOff")
	psp.interruptTasks()
	psp.manager.setDPMSModeOn()
	psp.resetBrightness()

	_, err := os.Stat("/etc/deepin/no_suspend")
	if err == nil {
		if psp.manager.ScreenBlackLock.Get() {
			//m.setDPMSModeOn()
			//m.lockWaitShow(4 * time.Second)
			psp.manager.doLock(false)
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func (psp *powerSavePlan) isWindowFullScreenAndFocused(xid x.Window) (bool, error) {
	conn := psp.manager.helper.xConn
	states, err := ewmh.GetWMState(conn, xid).Reply(conn)
	if err != nil {
		return false, err
	}
	found := 0
	for _, s := range states {
		if s == psp.atomNetWMStateFullscreen {
			found++
		}
		//后端的代码之前是基于deepin-wm这个窗口适配，现在换成了kwin,state中没有 focus 这个属性了
		//else if s == psp.atomNetWMStateFocused {
		//	found++
		//}
		if found == 1 {
			return true, nil
		}
	}
	return false, nil
}

func (psp *powerSavePlan) saveSetBrightnessTime(value bool, value2 map[string]float64) {
	if !value {
		return
	}
	psp.mu.Lock()
	psp.setBrightnessTime = time.Now()
	psp.mu.Unlock()
}
