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
	"strconv"
	"sync"
	"time"

	"dbus/org/bluez"
	"pkg.deepin.io/lib/dbus"
)

const (
	dbusAgentDest = "."
	dbusAgentPath = dbusBluetoothPath + "/Agent"
	dbusAgentIfs  = "org.bluez.Agent1"
)

type authorize struct {
	dpath  dbus.ObjectPath
	key    string
	accept bool
}

type agent struct {
	agentManager *bluez.AgentManager1
	b            *Bluetooth
	rspChan      chan authorize

	lk            sync.Mutex
	requestDevice dbus.ObjectPath
}

func (a *agent) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusAgentDest,
		ObjectPath: dbusAgentPath,
		Interface:  dbusAgentIfs,
	}
}

/*****************************************************************************/

//Release method gets called when the service daemon unregisters the agent.
//An agent can use it to do cleanup tasks. There is no need to unregister the
//agent, because when this method gets called it has already been unregistered.
func (a *agent) Release() {
	logger.Info("Release()")
}

//RequestPinCode method gets called when the service daemon needs to get the passkey for an authentication.
//The return value should be a string of 1-16 characters length. The string can be alphanumeric.
//Possible errors: org.bluez.Error.Rejected
//                 org.bluez.Error.Canceled
func (a *agent) RequestPinCode(dpath dbus.ObjectPath) (pincode string, err error) {
	logger.Info("RequestPinCode()")
	//return utils.RandString(8), nil
	auth, err := a.emitRequest(dpath, "RequestPinCode")
	if nil != err {
		return "", err
	}
	return auth.key, err
}

//DisplayPinCode method gets called when the service daemon needs to display a pincode for an authentication.
//An empty reply should be returned. When the pincode needs no longer to be displayed, the Cancel method
//of the agent will be called. This is used during the pairing process of keyboards that don't support
//Bluetooth 2.1 Secure Simple Pairing, in contrast to DisplayPasskey which is used for those that do.
//This method will only ever be called once since older keyboards do not support typing notification.
//Note that the PIN will always be a 6-digit number, zero-padded to 6 digits. This is for harmony with
//the later specification.
//Possible errors: org.bluez.Error.Rejected
//				   org.bluez.Error.Canceled
func (a *agent) DisplayPinCode(dpath dbus.ObjectPath, pincode string) (err error) {
	logger.Info("DisplayPinCode()", pincode)
	dbus.Emit(a.b, "DisplayPinCode", dpath, pincode)
	return
}

//RequestPasskey method gets called when the service daemon needs to get the passkey for an authentication.
//The return value should be a numeric value between 0-999999.
//Possible errors: org.bluez.Error.Rejected
//				   org.bluez.Error.Canceled
func (a *agent) RequestPasskey(dpath dbus.ObjectPath) (passkey uint32, err error) {
	//passkey = rand.Uint32() % 999999
	logger.Info("RequestPasskey()")
	auth, err := a.emitRequest(dpath, "RequestPasskey")
	if nil != err {
		return 0, err
	}
	key, err := strconv.ParseUint(auth.key, 10, 32)
	passkey = uint32(key)
	return passkey, err
}

//DisplayPasskey method gets called when the service daemon needs to display a passkey for an authentication.
//The entered parameter indicates the number of already typed keys on the remote side.
//An empty reply should be returned. When the passkey needs no longer to be displayed, the Cancel method
//of the agent will be called.
//During the pairing process this method might be called multiple times to update the entered value.
//Note that the passkey will always be a 6-digit number, so the display should be zero-padded at the start if
//the value contains less than 6 digits.
func (a *agent) DisplayPasskey(dpath dbus.ObjectPath, passkey uint32, entered uint16) {
	logger.Info("DisplayPasskey()", passkey, entered)
	dbus.Emit(a.b, "DisplayPasskey", dpath, passkey, entered)
	return
}

//RequestConfirmation This method gets called when the service daemon needs to confirm a passkey for an authentication.
//To confirm the value it should return an empty reply or an error in case the passkey is invalid.
//Note that the passkey will always be a 6-digit number, so the display should be zero-padded at the start if
//the value contains less than 6 digits.
//Possible errors: org.bluez.Error.Rejected
//			       org.bluez.Error.Canceled
func (a *agent) RequestConfirmation(dpath dbus.ObjectPath, passkey uint32) (err error) {
	logger.Info("RequestConfirmation", dpath, passkey)
	key := fmt.Sprintf("%06d", passkey)
	_, err = a.emitRequest(dpath, "RequestConfirmation", key)
	return err
}

//RequestAuthorization method gets called to request the user to authorize an incoming pairing attempt which
//would in other circumstances trigger the just-works model.
//Possible errors: org.bluez.Error.Rejected
//				   org.bluez.Error.Canceled
func (a *agent) RequestAuthorization(dpath dbus.ObjectPath) (err error) {
	logger.Info("RequestAuthorization()")
	_, err = a.emitRequest(dpath, "RequestAuthorization")
	return err
}

//AuthorizeService method gets called when the service daemon needs to authorize a connection/service request.
//Possible errors: org.bluez.Error.Rejected
//				   org.bluez.Error.Canceled
func (a *agent) AuthorizeService(dpath dbus.ObjectPath, uuid string) (err error) {
	logger.Info("AuthorizeService()")
	// TODO: DO NOT forbiden device connect service
	//dbus.Emit(a.b, "AuthorizeService")
	//return a.emitRequest(dpath, uuid, "AuthorizeService")
	return nil
}

//Cancel method gets called to indicate that the agent request failed before a reply was returned.
func (a *agent) Cancel() {
	logger.Info("Cancel()")
	a.rspChan <- authorize{dpath: a.requestDevice, accept: false, key: ""}
}

/*****************************************************************************/

func newAgent() (a *agent) {
	a = &agent{
		rspChan: make(chan authorize, 1),
	}
	return
}

func (a *agent) init() (err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			a.destroy()
		}
	}()

	a.agentManager, err = bluez.NewAgentManager1(dbusBluezDest, dbusBluezPath)
	if nil != err {
		logger.Info("Get AgentManager failed: ", err)
		return err
	}
	err = a.agentManager.RegisterAgent(dbusAgentPath, "DisplayYesNo")
	if nil != err {
		logger.Info("Register agent failed: ", err)
		return err
	}
	err = a.agentManager.RequestDefaultAgent(dbusAgentPath)
	if nil != err {
		logger.Info("Set default agent failed: ", err)
		return err
	}
	return nil
}

func (a *agent) destroy() {
	a.agentManager.UnregisterAgent(dbusAgentPath)
	dbus.UnInstallObject(a)
}

func (a *agent) waitRespone() (auth authorize, err error) {
	// TODO: time is too short or long
	logger.Info("waitRespone")
	defer func() { a.requestDevice = "" }()
	t := time.NewTimer(60 * time.Second)
	for {
		select {
		case auth = <-a.rspChan:
			logger.Info("receive", auth)
			if !auth.accept {
				err = errBluezRejected
				logger.Warningf("emitRequest return with: %v", err)
				return
			}
			logger.Infof("emitRequest accept %v with %v", a.requestDevice, auth.key)
			return
		case <-t.C:
			logger.Info("timeout")
			err = errBluezCanceled
			logger.Warningf("emitRequest return with: %v", err)
			return
		}
	}
	logger.Error("Shoud not run here!!!")
	err = errBluezCanceled
	return
}

func (a *agent) emit(signal string, dpath dbus.ObjectPath, ins ...interface{}) (err error) {
	sins := [](interface{}){}
	sins = append(sins, dpath)
	sins = append(sins, ins...)
	return dbus.Emit(a.b, signal, sins...)
}

func (a *agent) emitRequest(dpath dbus.ObjectPath, signal string, ins ...interface{}) (auth authorize, err error) {
	logger.Info("emitRequest", dpath, signal, ins)

	a.lk.Lock()
	a.requestDevice = dpath
	a.lk.Unlock()

	_, err = a.b.getDevice(dpath)
	if nil != err {
		logger.Warningf("emitRequest can not find device: %v, %v", dpath, err)
		return auth, errBluezCanceled
	}

	logger.Debug("Send Signal for device: ", dpath, signal, ins)
	a.emit(signal, dpath, ins...)

	return a.waitRespone()
}

func (b *Bluetooth) feed(dpath dbus.ObjectPath, accept bool, key string) (err error) {
	_, err = b.getDevice(dpath)
	if nil != err {
		logger.Warningf("FeedRequest can not find device: %v, %v", dpath, err)
		return err
	}

	b.agent.lk.Lock()
	if b.agent.requestDevice != dpath {
		b.agent.lk.Unlock()
		logger.Warningf("FeedRequest can not find match device: %v, %v", b.agent.requestDevice, dpath)
		return errBluezCanceled
	}
	b.agent.lk.Unlock()

	b.agent.rspChan <- authorize{dpath: dpath, accept: accept, key: key}
	return nil
}

//Confirm should call when you receive RequestConfirmation signal
func (b *Bluetooth) Confirm(dpath dbus.ObjectPath, accept bool) (err error) {
	logger.Info("Confirm()")
	return b.feed(dpath, accept, "")
}

//FeedPinCode should call when you receive RequestPinCode signal, notice that accept must true
//if you accept connect request. If accept is false, pincode will be ignored.
func (b *Bluetooth) FeedPinCode(dpath dbus.ObjectPath, accept bool, pincode string) (err error) {
	logger.Info("FeedPinCode()")
	return b.feed(dpath, accept, pincode)
}

//FeedPasskey should call when you receive RequestPasskey signal, notice that accept must true
//if you accept connect request. If accept is false, passkey will be ignored.
//passkey must be range in 0~999999.
func (b *Bluetooth) FeedPasskey(dpath dbus.ObjectPath, accept bool, passkey uint32) (err error) {
	logger.Info("FeedPasskey()")
	return b.feed(dpath, accept, fmt.Sprintf("%06d", passkey))
}
