//handle LidSwitch, PowerButton and Battery status event.
package main

import "os/exec"
import "dbus/com/deepin/sessionmanager"
import "dbus/org/freedesktop/upower"
import "dbus/org/freedesktop/login1"
import "dbus/com/deepin/daemon/keybinding"
import "fmt"

const (
	//sync with com.deepin.daemon.power.schemas
	ActionBlank       int32 = 0
	ActionSuspend           = 1
	ActionShutdown          = 2
	ActionHibernate         = 3
	ActionInteractive       = 4
	ActionNothing           = 5
	ActionLogout            = 6
)

func doLock() {
	if m, err := sessionmanager.NewSessionManager("com.deepin.SessionManager", "/com/deepin/SessionManager"); err != nil {
		LOGGER.Warning("can't build SessionManager Object:", err)
	} else {
		if err = m.RequestLock(); err != nil {
			LOGGER.Warning("Lock failed:", err)
		}
		sessionmanager.DestroySessionManager(m)
	}

}

func doShowLowpower() {
	go exec.Command("/usr/lib/deepin-daemon/lowpower").Run()
}
func doCloseLowpower() {
	go exec.Command("killall", "lowpower").Run()
}

func doShutDown() {
	if m, err := sessionmanager.NewSessionManager("com.deepin.SessionManager", "/com/deepin/SessionManager"); err != nil {
		LOGGER.Warning("can't build SessionManager Object:", err)
	} else {
		if err = m.RequestShutdown(); err != nil {
			LOGGER.Warning("Shutdown failed:", err)
		}
		sessionmanager.DestroySessionManager(m)
	}
}

func doSuspend() {
	if m, err := sessionmanager.NewSessionManager("com.deepin.SessionManager", "/com/deepin/SessionManager"); err != nil {
		LOGGER.Warning("can't build SessionManager Object:", err)
	} else {
		if err = m.RequestSuspend(); err != nil {
			LOGGER.Warning("Suspend failed:", err)
		}
		fmt.Println("RequestSuspend...", err)
		sessionmanager.DestroySessionManager(m)
	}
}

func doLogout() {
	if m, err := sessionmanager.NewSessionManager("com.deepin.SessionManager", "/com/deepin/SessionManager"); err != nil {
		LOGGER.Warning("can't build SessionManager Object:", err)
	} else {
		if err = m.Logout(); err != nil {
			LOGGER.Warning("ShutDown failed:", err)
		}
		sessionmanager.DestroySessionManager(m)
	}
}

func doShutDownInteractive() {
	go exec.Command("/usr/lib/deepin-daemon/dshutdown").Run()
}

func (up *Power) handlePowerButton() {
	fmt.Println("HandlePowerButton:")
	switch up.PowerButtonAction.Get() {
	case ActionInteractive:
		doShutDownInteractive()
	case ActionShutdown:
		doShutDown()
	case ActionNothing:
	default:
		LOGGER.Warning("invalid LidSwitchAction:", up.LidClosedAction)
	}
}

func (up *Power) handleCloseLidSwitch() {
	switch up.LidClosedAction.Get() {
	case ActionInteractive:
		doShutDownInteractive()
	case ActionSuspend:
		doSuspend()
	case ActionShutdown:
		doShutDown()
	case ActionNothing:
	default:
		LOGGER.Warning("invalid LidSwitchAction:", up.LidClosedAction.Get())
	}
}

func (p *Power) initEventHandle() {
	up, err := upower.NewUpower(UPOWER_BUS_NAME, "/org/freedesktop/UPower")
	if err != nil {
		LOGGER.Error("Can't build org.freedesktop.UPower:", err)
	} else {
		up.ConnectChanged(func() {
			currentLidCloed := up.LidIsClosed.Get()
			fmt.Println("LidState:", currentLidCloed)
			if p.lidIsClosed != currentLidCloed {
				p.lidIsClosed = currentLidCloed
				if currentLidCloed {
					p.handleCloseLidSwitch()
				}
			}
			p.lidIsClosed = currentLidCloed

		})
	}

	mediaKey, err := keybinding.NewMediaKey("com.deepin.daemon.KeyBinding", "/com/deepin/daemon/MediaKey")
	if err != nil {
		LOGGER.Error("Can't build com.deepin.daemon.KeyBinding:", err)
	} else {
		mediaKey.ConnectPowerOff(func(press bool) {
			//prevent mediaKey be destroyed
			mediaKey.DestName = mediaKey.DestName

			if !press {
				p.handlePowerButton()
			}
		})
	}

	login, err := login1.NewManager("org.freedesktop.login1", "/org/freedesktop/login1")
	if err != nil {
		LOGGER.Error("Can't build org.freedesktop.login1:", err)
	} else {
		login.ConnectPrepareForSleep(func(before bool) {
			fmt.Println("Sleep change...", before)
			if before {
				if p.coreSettings.GetBoolean("lock-enabled") {
					doLock()
				}
			} else {
				p.handleBatteryPercentage()
				if p.lowBatteryStatus == lowBatteryStatusAction {
					doShowLowpower()
				}
			}
		})
	}
}
