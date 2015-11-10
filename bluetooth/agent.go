/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package bluetooth

import (
	"dbus/org/bluez"
	"pkg.deepin.io/lib/dbus"
)

const (
	dbusAgentDest = "."
	dbusAgentPath = dbusBluetoothPath + "/Agent"
	dbusAgentIfs  = "org.bluez.Agent1"
)

type Agent struct {
	agentManager *bluez.AgentManager1
}

func (a *Agent) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusAgentDest,
		ObjectPath: dbusAgentPath,
		Interface:  dbusAgentIfs,
	}
}
func newAgent() (agent *Agent) {
	// TODO
	agent = &Agent{}
	return
}

func (a *Agent) init() (err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			a.destroy()
		}
	}()

	a.agentManager, err = bluez.NewAgentManager1(dbusBluezDest, dbusBluezPath)
	if nil != err {
		logger.Info("get agentmanager failed: ", err)
		return err
	}
	err = a.agentManager.RegisterAgent(dbusAgentPath, "DisplayYesNo")
	if nil != err {
		logger.Info("register agent failed: ", err)
		return err
	}
	err = a.agentManager.RequestDefaultAgent(dbusAgentPath)
	if nil != err {
		logger.Info("set defaulet agent failed: ", err)
		return err
	}
	return nil
}

func (a *Agent) destroy() {
	a.agentManager.UnregisterAgent(dbusAgentPath)
	dbus.UnInstallObject(a)
}

// TODO
func (a *Agent) Release() {
	// TODO
	logger.Info("Release()")
}

func (a *Agent) RequestPinCode(dpath dbus.ObjectPath) (pincode string) {
	// TODO
	logger.Info("RequestPinCode()")
	return
}

func (a *Agent) DisplayPinCode(dpath dbus.ObjectPath, pincode string) {
	logger.Info("DisplayPinCode()")
}

func (a *Agent) DisplayPasskey(dpath dbus.ObjectPath, passkey uint32, entered uint16) {
	logger.Info("DisplayPasskey()")
}

func (a *Agent) RequestConfirmation(dpath dbus.ObjectPath, passkey uint32) {
	logger.Info("RequestConfirmation", dpath, passkey)
}

func (a *Agent) RequestAuthorization(dpath dbus.ObjectPath) {
	// TODO
	logger.Info("RequestAuthorization()")
}

func (a *Agent) AuthorizeService(dpath dbus.ObjectPath, uuid string) {
	// TODO
	logger.Info("AuthorizeService()")
}

func (a *Agent) Cancel() {
	// TODO
	logger.Info("Cancel()")
}

// TODO
func (b *Bluetooth) FeedPinCode(pincode string) (err error) {
	logger.Info("FeedPinCode()")
	return
}

// TODO
func (b *Bluetooth) FeedAuthorizeService(pincode string) (err error) {
	logger.Info("FeedAuthorizeService()")
	return
}
