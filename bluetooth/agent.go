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

	bluez "github.com/linuxdeepin/go-dbus-factory/org.bluez"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	agentDBusPath      = dbusPath + "/Agent"
	agentDBusInterface = "org.bluez.Agent1"
)

type authorize struct {
	path   dbus.ObjectPath
	key    string
	accept bool
}

type agent struct {
	service      *dbusutil.Service
	bluezManager *bluez.Manager

	b       *Bluetooth
	rspChan chan authorize

	mu            sync.Mutex
	requestDevice dbus.ObjectPath

	methods *struct {
		RequestPinCode       func() `in:"device" out:"pinCode"`
		DisplayPinCode       func() `in:"device,pinCode"`
		RequestPasskey       func() `in:"device" out:"passkey"`
		DisplayPasskey       func() `in:"device,passkey,entered"`
		RequestConfirmation  func() `in:"device,passkey"`
		RequestAuthorization func() `in:"device"`
		AuthorizeService     func() `in:"device,uuid"`
	}
}

func (*agent) GetInterfaceName() string {
	return agentDBusInterface
}

/*****************************************************************************/

//Release method gets called when the service daemon unregisters the agent.
//An agent can use it to do cleanup tasks. There is no need to unregister the
//agent, because when this method gets called it has already been unregistered.
func (a *agent) Release() *dbus.Error {
	logger.Info("Release()")
	return nil
}

//RequestPinCode method gets called when the service daemon needs to get the passkey for an authentication.
//The return value should be a string of 1-16 characters length. The string can be alphanumeric.
//Possible errors: org.bluez.Error.Rejected
//                 org.bluez.Error.Canceled
func (a *agent) RequestPinCode(dpath dbus.ObjectPath) (pincode string, busErr *dbus.Error) {
	logger.Info("RequestPinCode()")

	d, err := a.b.getDevice(dpath)
	if err != nil {
		logger.Warning(err)
		return "", dbusutil.ToError(err)
	}
	d.agentWorkStart()
	defer d.agentWorkEnd()

	//return utils.RandString(8), nil
	auth, err := a.emitRequest(dpath, "RequestPinCode")
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	return auth.key, nil
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
func (a *agent) DisplayPinCode(devPath dbus.ObjectPath, pinCode string) (err *dbus.Error) {
	logger.Info("DisplayPinCode()", pinCode)
	a.b.service.Emit(a.b, "DisplayPinCode", devPath, pinCode)
	return
}

//RequestPasskey method gets called when the service daemon needs to get the passkey for an authentication.
//The return value should be a numeric value between 0-999999.
//Possible errors: org.bluez.Error.Rejected
//				   org.bluez.Error.Canceled
func (a *agent) RequestPasskey(dpath dbus.ObjectPath) (passkey uint32, busErr *dbus.Error) {
	logger.Info("RequestPasskey()")

	d, err := a.b.getDevice(dpath)
	if err != nil {
		logger.Warning(err)
		return 0, dbusutil.ToError(err)
	}
	d.agentWorkStart()
	defer d.agentWorkEnd()

	//passkey = rand.Uint32() % 999999
	auth, err := a.emitRequest(dpath, "RequestPasskey")
	if err != nil {
		return 0, dbusutil.ToError(err)
	}

	key, err := strconv.ParseUint(auth.key, 10, 32)
	if err != nil {
		return 0, dbusutil.ToError(err)
	}
	passkey = uint32(key)
	return passkey, nil
}

//DisplayPasskey method gets called when the service daemon needs to display a passkey for an authentication.
//The entered parameter indicates the number of already typed keys on the remote side.
//An empty reply should be returned. When the passkey needs no longer to be displayed, the Cancel method
//of the agent will be called.
//During the pairing process this method might be called multiple times to update the entered value.
//Note that the passkey will always be a 6-digit number, so the display should be zero-padded at the start if
//the value contains less than 6 digits.
func (a *agent) DisplayPasskey(dpath dbus.ObjectPath, passkey uint32,
	entered uint16) *dbus.Error {

	logger.Info("DisplayPasskey()", passkey, entered)
	err := a.b.service.Emit(a.b, "DisplayPasskey", dpath, passkey, uint32(entered))
	if err != nil {
		logger.Warning("Failed to emit signal 'DisplayPasskey':", err, dpath, passkey, entered)
	}
	return nil
}

//RequestConfirmation This method gets called when the service daemon needs to confirm a passkey for an authentication.
//To confirm the value it should return an empty reply or an error in case the passkey is invalid.
//Note that the passkey will always be a 6-digit number, so the display should be zero-padded at the start if
//the value contains less than 6 digits.
//Possible errors: org.bluez.Error.Rejected
//			       org.bluez.Error.Canceled
func (a *agent) RequestConfirmation(dpath dbus.ObjectPath, passkey uint32) *dbus.Error {
	logger.Info("RequestConfirmation", dpath, passkey)

	d, err := a.b.getDevice(dpath)
	if err != nil {
		logger.Warning(err)
		return dbusutil.ToError(err)
	}
	d.agentWorkStart()
	defer d.agentWorkEnd()

	key := fmt.Sprintf("%06d", passkey)
	_, err = a.emitRequest(dpath, "RequestConfirmation", key)
	return dbusutil.ToError(err)
}

//RequestAuthorization method gets called to request the user to authorize an incoming pairing attempt which
//would in other circumstances trigger the just-works model.
//Possible errors: org.bluez.Error.Rejected
//				   org.bluez.Error.Canceled
func (a *agent) RequestAuthorization(dpath dbus.ObjectPath) *dbus.Error {
	logger.Info("RequestAuthorization()")

	d, err := a.b.getDevice(dpath)
	if err != nil {
		logger.Warning(err)
		return dbusutil.ToError(err)
	}
	d.agentWorkStart()
	defer d.agentWorkEnd()

	_, err = a.emitRequest(dpath, "RequestAuthorization")
	return dbusutil.ToError(err)
}

//AuthorizeService method gets called when the service daemon needs to authorize a connection/service request.
//Possible errors: org.bluez.Error.Rejected
//				   org.bluez.Error.Canceled
func (a *agent) AuthorizeService(dpath dbus.ObjectPath, uuid string) *dbus.Error {
	logger.Info("AuthorizeService()")
	// TODO: DO NOT forbiden device connect service
	//dbus.Emit(a.b, "AuthorizeService")
	//return a.emitRequest(dpath, uuid, "AuthorizeService")
	return nil
}

//Cancel method gets called to indicate that the agent request failed before a reply was returned.
func (a *agent) Cancel() *dbus.Error {
	logger.Info("Cancel()")
	a.rspChan <- authorize{path: a.requestDevice, accept: false, key: ""}
	a.emitCancelled()
	return nil
}

/*****************************************************************************/

func newAgent(service *dbusutil.Service) (a *agent) {
	a = &agent{
		service: service,
		rspChan: make(chan authorize),
	}
	return
}

func (a *agent) init() {
	sysBus := a.service.Conn()
	a.bluezManager = bluez.NewManager(sysBus)
	a.registerDefaultAgent()
	return
}

func (a *agent) registerDefaultAgent() {
	// register agent
	err := a.bluezManager.RegisterAgent(0, agentDBusPath, "DisplayYesNo")
	if err != nil {
		logger.Warning("failed to register agent:", err)
		return
	}

	// request default agent
	err = a.bluezManager.RequestDefaultAgent(0, agentDBusPath)
	if err != nil {
		logger.Warning("failed to become the default agent:", err)
		err = a.bluezManager.UnregisterAgent(0, agentDBusPath)
		if err != nil {
			logger.Warning(err)
		}
		return
	}
}

func (a *agent) destroy() {
	err := a.bluezManager.UnregisterAgent(0, agentDBusPath)
	if err != nil {
		logger.Warning(err)
	}

	err = a.service.StopExport(a)
	if err != nil {
		logger.Warning(err)
	}
}

func (a *agent) waitResponse() (auth authorize, err error) {
	logger.Info("waitResponse")

	defer func() {
		a.mu.Lock()
		a.requestDevice = ""
		a.mu.Unlock()
	}()

	t := time.NewTimer(60 * time.Second)
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
		a.emitCancelled()
		return
	}
}

func (a *agent) emit(signal string, devPath dbus.ObjectPath, args ...interface{}) (err error) {
	var args0 []interface{}
	args0 = append(args0, devPath)
	args0 = append(args0, args...)
	return a.b.service.Emit(a.b, signal, args0...)
}

func (a *agent) emitCancelled() {
	a.mu.Lock()
	devPath := a.requestDevice
	a.mu.Unlock()

	if devPath == "" {
		logger.Warning("failed to emitCancelled, devPath is empty")
		return
	}
	err := a.b.service.Emit(a.b, "Cancelled", devPath)
	if err != nil {
		logger.Warning(err)
	}
}

func (a *agent) emitRequest(devPath dbus.ObjectPath, signal string, args ...interface{}) (auth authorize, err error) {
	logger.Info("emitRequest", devPath, signal, args)

	a.mu.Lock()
	a.requestDevice = devPath
	a.mu.Unlock()

	_, err = a.b.getDevice(devPath)
	if nil != err {
		logger.Warningf("emitRequest can not find device: %v, %v", devPath, err)
		return auth, errBluezCanceled
	}

	logger.Debug("Send Signal for device: ", devPath, signal, args)
	a.emit(signal, devPath, args...)

	return a.waitResponse()
}
