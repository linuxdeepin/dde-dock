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

package bluetooth

import (
	"fmt"
	"os/exec"
	"strconv"
	"sync"
	"time"

	notifications "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.notifications"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	. "pkg.deepin.io/lib/gettext"
)

const (
	notifyIconBluetoothConnected     = "notification-bluetooth-connected"
	notifyIconBluetoothDisconnected  = "notification-bluetooth-disconnected"
	notifyIconBluetoothConnectFailed = "notification-bluetooth-error"
	// dialog use for show pinCode
	notifyDdeDialogPath = "/usr/lib/deepin-daemon/dde-bluetooth-dialog"
	// notification window stay time
	notifyTimerDuration = 5 * time.Second
)

const bluetoothDialog string = "dde-bluetooth-dialog"

var globalNotifications *notifications.Notifications
var globalNotifyId uint32
var globalNotifyMu sync.Mutex

func initNotifications() error {
	// init global notification timer instance
	globalTimerNotifier = GetTimerNotifyInstance()

	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	globalNotifications = notifications.NewNotifications(sessionBus)

	// monitor notification-close-signal
	sessionLoop := dbusutil.NewSignalLoop(sessionBus, 10)
	sessionLoop.Start()
	globalNotifications.InitSignalExt(sessionLoop, true)
	_, err = globalNotifications.ConnectActionInvoked(func(id uint32, actionKey string) {
		// has received signal, use id to compare with last globalNotifyId
		if id == globalNotifyId {
			// if it is the same, then send chan to instance chan to close window
			globalTimerNotifier.actionInvokedChan <- true
		}
	})
	if err != nil {
		logger.Warningf("listen action invoked failed,err:%v", err)
	}

	return nil
}

func notify(icon, summary, body string) {
	logger.Info("notify", icon, summary, body)

	globalNotifyMu.Lock()
	nid := globalNotifyId
	globalNotifyMu.Unlock()

	nid, err := globalNotifications.Notify(0, "dde-control-center", nid, icon,
		summary, body, nil, nil, -1)
	if err != nil {
		logger.Warning(err)
		return
	}
	globalNotifyMu.Lock()
	globalNotifyId = nid
	globalNotifyMu.Unlock()
}

func notifyConnected(alias string) {
	format := Tr("Connect %q successfully")
	notify(notifyIconBluetoothConnected, "", fmt.Sprintf(format, alias))
}

func notifyDisconnected(alias string) {
	format := Tr("%q disconnected")
	notify(notifyIconBluetoothDisconnected, "", fmt.Sprintf(format, alias))
}

func notifyConnectFailedHostDown(alias string) {
	format := Tr("Make sure %q is turned on and in range")
	notifyConnectFailedAux(alias, format)
}

func notifyConnectFailedAux(alias, format string) {
	notify(notifyIconBluetoothConnectFailed, Tr("Bluetooth connection failed"), fmt.Sprintf(format, alias))
}

// notify pc initiative connect to device
// so dont need to show notification window
func notifyInitiativeConnect(dev *device, pinCode string) error {
	if checkProcessExists(bluetoothDialog) {
		logger.Info("initiative already exist")
		return nil
	}

	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	//use command to open osd window to show pin code
	cmd := exec.Command(notifyDdeDialogPath, pinCode, string(dev.Path), timestamp)
	err := cmd.Start()
	if err != nil {
		logger.Infof("execute cmd command failed,err:%v", err)
		return err
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			logger.Warning(err)
		}
	}()

	return nil
}

// device passive connect to pc
// so need to show notification window
func notifyPassiveConnect(dev *device, pinCode string) error {
	format := Tr("Click here to connect to %q")
	summary := Tr("Add Bluetooth devices")
	body := fmt.Sprintf(format, dev.Name)
	globalNotifyMu.Lock()
	nid := globalNotifyId
	globalNotifyMu.Unlock()
	// check if bluetooth dialog is exist
	if checkProcessExists(bluetoothDialog) {
		logger.Info("Passive is not exist")
		return nil
	}
	var as = []string{"pair", Tr("Pair")}
	var timestamp = strconv.FormatInt(time.Now().UnixNano(), 10)
	cmd := notifyDdeDialogPath + "," + pinCode + "," + string(dev.Path) + "," + timestamp
	hints := map[string]dbus.Variant{"x-deepin-action-pair": dbus.MakeVariant(cmd)}

	// to make sure last notification has been closed
	err := globalNotifications.CloseNotification(0, nid)
	if err != nil {
		logger.Warningf("close last notification failed,err:%v", err)
	}

	// notify connect request to dde-control-center
	// set notify time out as -1, default time out is 5 seconds
	nid, err = globalNotifications.Notify(0, "dde-control-center", nid, notifyIconBluetoothConnected,
		summary, body, as, hints, -1)
	if err != nil {
		logger.Warningf("notify message failed,err:%v", err)
		return err
	}

	globalNotifyMu.Lock()
	globalNotifyId = nid
	globalNotifyMu.Unlock()

	// to avoid to show more than one window and fix notification time out incorrect
	// need to reset timer
	globalTimerNotifier.timeout.Reset(notifyTimerDuration)

	return nil
}

// global timer notifier
var globalTimerNotifier *timerNotify

// notify timer instance
// use chan bool instead of timer, but in case to fit new requirements of future flexibly, we keep element timer
type timerNotify struct {
	timeout           *time.Timer
	actionInvokedChan chan bool
}

// get timer instance
func GetTimerNotifyInstance() *timerNotify {
	// create a global timer notify object
	timerNotifier := &timerNotify{
		timeout:           time.NewTimer(notifyTimerDuration),
		actionInvokedChan: make(chan bool),
	}
	timerNotifier.timeout.Stop()
	return timerNotifier
}

// begin timer routine to monitor window click notification window or notification time out signal
func beginTimerNotify(notifyTimer *timerNotify) {
	for {
		select {
		case <-notifyTimer.timeout.C:
			// monitor time out signal
			logger.Info("user no response,close notify when time out")
			err := globalNotifications.CloseNotification(0, globalNotifyId)
			if err != nil {
				logger.Warningf("time out close notify icon failed,err:%v", err)
			}
		case <-notifyTimer.actionInvokedChan:
			// monitor click window signal
			logger.Info("user click notify,close notify")
			err := globalNotifications.CloseNotification(0, globalNotifyId)
			if err != nil {
				logger.Warningf("click event close notify icon failed,err:%v", err)
			}
			// if window is clicked, then stop timer
			notifyTimer.timeout.Stop()
		}
	}
}
