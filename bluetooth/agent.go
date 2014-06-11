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
