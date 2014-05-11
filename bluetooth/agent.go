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
func (a *Agent) Release() {}
func (a *Agent) RequestPinCode(dpath dbus.ObjectPath) (pincode string) {
	// TODO
	return
}
func (a *Agent) DisplayPinCode(dpath dbus.ObjectPath, pincode string)                 {}
func (a *Agent) DisplayPasskey(dpath dbus.ObjectPath, passkey uint32, entered uint16) {}
func (a *Agent) RequestConfirmation(dpath dbus.ObjectPath, passkey uint32)            {}
func (a *Agent) RequestAuthorization(dpath dbus.ObjectPath)                           {}
func (a *Agent) AuthorizeService(dpath dbus.ObjectPath, uuid string) {
	// TODO
}
func (a *Agent) Cancel() {}

// TODO
func (b *Bluetooth) FeedPinCode(pincode string) (err error) {
	return
}

// TODO
func (b *Bluetooth) FeedAuthorizeService(pincode string) (err error) {
	return
}
