//handle LidSwitch, PowerButton and Battery status event.
package power

import "os/exec"
import "dbus/com/deepin/sessionmanager"
import "time"
import "dbus/com/deepin/daemon/display"
import "syscall"

const (
	//sync with com.deepin.daemon.power.schemas
	//
	// 按下电源键和合上笔记本盖时支持的操作
	//
	// 关闭显示器
	ActionBlank int32 = 0
	// 挂起
	ActionSuspend = 1
	// 关机
	ActionShutdown = 2
	// 休眠
	ActionHibernate = 3
	// 询问
	ActionInteractive = 4
	// 无
	ActionNothing = 5
	// 注销
	ActionLogout = 6
)

func doLock() {
	if m, err := sessionmanager.NewSessionManager("com.deepin.SessionManager", "/com/deepin/SessionManager"); err != nil {
		logger.Warning("can't build SessionManager Object:", err)
	} else {
		if err = m.RequestLock(); err != nil {
			logger.Warning("Lock failed:", err)
		}
		sessionmanager.DestroySessionManager(m)
	}

}

func doShowLowpower() {
	go exec.Command("/usr/lib/deepin-daemon/dde-lowpower").Run()
}
func doCloseLowpower() {
	go exec.Command("killall", "dde-lowpower").Run()
}

func doShutDown() {
	if m, err := sessionmanager.NewSessionManager("com.deepin.SessionManager", "/com/deepin/SessionManager"); err != nil {
		logger.Warning("can't build SessionManager Object:", err)
	} else {
		if err = m.RequestShutdown(); err != nil {
			logger.Warning("Shutdown failed:", err)
		}
		sessionmanager.DestroySessionManager(m)
	}
}

func doSuspend() {
	if m, err := sessionmanager.NewSessionManager("com.deepin.SessionManager", "/com/deepin/SessionManager"); err != nil {
		logger.Warning("can't build SessionManager Object:", err)
	} else {
		if err = m.RequestSuspend(); err != nil {
			logger.Warning("Suspend failed:", err)
		}
		logger.Debug("RequestSuspend...", err)
		sessionmanager.DestroySessionManager(m)
	}
}

func doLogout() {
	if m, err := sessionmanager.NewSessionManager("com.deepin.SessionManager", "/com/deepin/SessionManager"); err != nil {
		logger.Warning("can't build SessionManager Object:", err)
	} else {
		if err = m.Logout(); err != nil {
			logger.Warning("ShutDown failed:", err)
		}
		sessionmanager.DestroySessionManager(m)
	}
}

func doShutDownInteractive() {
	go exec.Command("dde-shutdown").Run()
}

func (up *Power) handlePowerButton() {
	switch up.PowerButtonAction.Get() {
	case ActionInteractive:
		doShutDownInteractive()
	case ActionShutdown:
		doShutDown()
	case ActionSuspend:
		doSuspend()
	case ActionNothing:
	default:
		logger.Warning("invalid LidSwitchAction:", up.LidClosedAction)
	}
}

func (up *Power) handleLidSwitch(opened bool) {
	if opened {
		logger.Info("Lid opened...")
		//TODO: DPMS ON
	} else {
		logger.Info("Lid closed...")
		//TODO: DPMS OFF
		switch up.LidClosedAction.Get() {
		case ActionInteractive:
			doShutDownInteractive()
		case ActionSuspend:
			if isMultihead() && !up.coreSettings.GetBoolean("lid-close-suspend-with-external-monitor") {
				logger.Info("Prevent suspend when lidclosed because another monitor connected")
				return
			}
			doSuspend()
		case ActionShutdown:
			doShutDown()
		case ActionNothing:
		default:
			logger.Warning("invalid LidSwitchAction:", up.LidClosedAction.Get())
		}
	}
}

func isMultihead() bool {
	if dp, err := display.NewDisplay("com.deepin.daemon.Display", "/com/deepin/daemon/Display"); err != nil {
		logger.Error("Can't build com.deepin.daemon.Display Object:", err)
		return false
	} else {
		paths := dp.Monitors.Get()
		if len(paths) > 1 {
			return true
		} else if len(paths) == 1 {
			if m, err := display.NewMonitor("com.deepin.daemon.Display", paths[0]); err != nil {
				return false
			} else if m.IsComposited.Get() {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

func (p *Power) initEventHandle() {
	if upower != nil {
		upower.LidIsClosed.ConnectChanged(func() {
			currentLidClosed := upower.LidIsClosed.Get()
			if p.lidIsClosed != currentLidClosed {
				p.lidIsClosed = currentLidClosed
				p.handleLidSwitch(!currentLidClosed)
			}
			p.lidIsClosed = currentLidClosed
		})
	}

	if mediaKey != nil {
		mediaKey.ConnectPowerOff(func(press bool) {
			if !press {
				p.handlePowerButton()
			}
		})
	}

	if login1 != nil {
		var blockSleep, unblockSleep = func() (func(), func()) {
			var blockFD = -1
			return func() {
					if blockFD == -1 {
						fd, err := login1.Inhibit("sleep", "lock screen", "run screenlock..", "delay")
						blockFD = int(fd)
						if err != nil {
							logger.Warning("inbhibit login1.sleep failed", err)
						}
					}
				}, func() {
					if blockFD >= 0 {
						err := syscall.Close(blockFD)
						if err != nil {
							logger.Warning("error when close fd:", err)
						}
						blockFD = -1
					}
				}
		}()

		blockSleep()

		login1.ConnectPrepareForSleep(func(before bool) {
			if before {
				unblockSleep()
				return
			}

			// Wakeup
			time.AfterFunc(time.Second*1, func() {
				p.screensaver.SimulateUserActivity()
			})

			playSound("wakeup")
			p.handleBatteryPercentage()
			if p.coreSettings.GetBoolean(settingKeyLockEnabled) {
				now := time.Now()
				doLock()
				logger.Debug("screenlock ready time:", time.Now().Sub(now))
			}

			blockSleep()
		})
	}
}
