package main

import (
	"dlib/dbus"
)

type Agent struct {
	b *Bluetooth
}

func newAgent() (agent *Agent) {
	// TODO
	agent = &Agent{}
	return
}

// TODO
func (a *Agent) Release() {
	// TODO
	logger.Debug("Release()")
}

func (a *Agent) RequestPinCode(dpath dbus.ObjectPath) (pincode string) {
	// TODO
	logger.Debug("RequestPinCode()")
	return
}

func (a *Agent) DisplayPinCode(dpath dbus.ObjectPath, pincode string) {
	logger.Debug("DisplayPinCode()")
}

func (a *Agent) DisplayPasskey(dpath dbus.ObjectPath, passkey uint32, entered uint16) {
	logger.Debug("DisplayPasskey()")
}

func (a *Agent) RequestConfirmation(dpath dbus.ObjectPath, passkey uint32) {
	logger.Debug("RequestConfirmation()")
}

func (a *Agent) RequestAuthorization(dpath dbus.ObjectPath) {
	// TODO
	logger.Debug("RequestAuthorization()")
}

func (a *Agent) AuthorizeService(dpath dbus.ObjectPath, uuid string) {
	// TODO
	logger.Debug("AuthorizeService()")
}

func (a *Agent) Cancel() {
	// TODO
	logger.Debug("Cancel()")
}

// TODO
func (b *Bluetooth) FeedPinCode(pincode string) (err error) {
	return
}

// TODO
func (b *Bluetooth) FeedAuthorizeService(pincode string) (err error) {
	return
}
