/*
 * Copyright (C) 2015 ~ 2017 Deepin Technology Co., Ltd.
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

package lastore

import (
	"fmt"
	"os"

	"pkg.deepin.io/lib/gettext"
)

type NotifyAction struct {
	Id       string
	Name     string
	Callback func()
}

func getAppStoreAppName() string {
	_, err := os.Stat("/usr/share/applications/deepin-app-store.desktop")
	if err == nil {
		return "deepin-app-store"
	}
	return "deepin-appstore"
}

const (
	notifyExpireTimeoutNever   = 0
	notifyExpireTimeoutDefault = -1
)

func (l *Lastore) sendNotify(icon string, msg string, actions []NotifyAction, expireTimeout int32, appName string) {
	logger.Infof("sendNotify icon: %s, msg: %q, actions: %v, timeout: %d",
		icon, msg, actions, expireTimeout)
	n := l.notifications
	var as []string
	for _, action := range actions {
		as = append(as, action.Id, action.Name)
	}

	if icon == "" {
		icon = "deepin-appstore"
	}
	id, err := n.Notify(0, appName, 0, icon, "",
		msg, as, nil, expireTimeout)
	if err != nil {
		logger.Warningf("Notify failed: %q: %v\n", msg, err)
		return
	}

	if len(actions) == 0 {
		return
	}

	hid, err := l.notifications.ConnectActionInvoked(
		func(id0 uint32, actionId string) {
			logger.Debugf("notification action invoked id: %d, actionId: %q",
				id0, actionId)
			if id != id0 {
				return
			}

			for _, action := range actions {
				if action.Id == actionId && action.Callback != nil {
					action.Callback()
					return
				}
			}
			logger.Warningf("not found action id %q in %v", actionId, actions)
		})
	if err != nil {
		logger.Warning(err)
		return
	}
	logger.Debugf("notifyIdHidMap[%d]=%d", id, hid)
	l.notifyIdHidMap[id] = hid
}

// NotifyInstall send desktop notify for install job
func (l *Lastore) notifyInstall(pkgId string, succeed bool, ac []NotifyAction) {
	var msg string
	if succeed {
		msg = fmt.Sprintf(gettext.Tr("%q installed successfully."), pkgId)
		l.sendNotify("package_install_succeed", msg, ac, notifyExpireTimeoutDefault, getAppStoreAppName())
	} else {
		msg = fmt.Sprintf(gettext.Tr("%q failed to install."), pkgId)
		l.sendNotify("package_install_failed", msg, ac, notifyExpireTimeoutDefault, getAppStoreAppName())
	}
}

func (l *Lastore) notifyRemove(pkgId string, succeed bool, ac []NotifyAction) {
	var msg string
	if succeed {
		msg = fmt.Sprintf(gettext.Tr("%q removed successfully."), pkgId)
	} else {
		msg = fmt.Sprintf(gettext.Tr("%q failed to remove."), pkgId)
	}
	l.sendNotify("deepin-appstore", msg, ac, notifyExpireTimeoutDefault, getAppStoreAppName())
}

//NotifyLowPower send notify for low power
func (l *Lastore) notifyLowPower() {
	msg := gettext.Tr("In order to prevent automatic shutdown, please plug in for normal update.")
	l.sendNotify("notification-battery_low", msg, nil, notifyExpireTimeoutDefault, getAppStoreAppName())
}

func (l *Lastore) notifyAutoClean() {
	msg := gettext.Tr("Package cache wiped")
	l.sendNotify("deepin-appstore", msg, nil, notifyExpireTimeoutDefault, "dde-control-center")
}

func (l *Lastore) notifySourceModified(actions []NotifyAction) {
	msg := gettext.Tr("Your system source has been modified, please restore to official source for your normal use")
	l.sendNotify("dialog-warning", msg, actions, notifyExpireTimeoutNever, getAppStoreAppName())
}

func (l *Lastore) notifyUpdateSource(actions []NotifyAction) {
	msg := gettext.Tr("New system edition available")
	l.sendNotify("dde-control-center", msg, actions, notifyExpireTimeoutNever, "dde-control-center")
}
