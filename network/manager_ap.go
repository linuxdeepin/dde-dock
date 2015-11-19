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

package network

import (
	nm "dbus/org/freedesktop/networkmanager"
	"fmt"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/utils"
)

type apSecType uint32

const (
	apSecNone apSecType = iota
	apSecWep
	apSecPsk
	apSecEap
)

type accessPoint struct {
	nmAp    *nm.AccessPoint
	devPath dbus.ObjectPath

	Ssid         string
	Secured      bool
	SecuredInEap bool
	Strength     uint8
	Path         dbus.ObjectPath
}

func (m *Manager) newAccessPoint(devPath, apPath dbus.ObjectPath) (ap *accessPoint, err error) {
	nmAp, err := nmNewAccessPoint(apPath)
	if err != nil {
		return
	}

	ap = &accessPoint{
		nmAp:    nmAp,
		devPath: devPath,
		Path:    apPath,
	}
	ap.updateProps()
	if len(ap.Ssid) == 0 {
		err = fmt.Errorf("ignore hidden access point")
		return
	}

	// connect property changed signals
	ap.nmAp.ConnectPropertiesChanged(func(properties map[string]dbus.Variant) {
		if !m.isAccessPointExists(apPath) {
			return
		}

		m.accessPointsLock.Lock()
		defer m.accessPointsLock.Unlock()
		ap.updateProps()
		logger.Debugf("access point properties changed %#v", ap)
		apJSON, _ := marshalJSON(ap)
		dbus.Emit(m, "AccessPointPropertiesChanged", string(devPath), apJSON)
	})

	apJSON, _ := marshalJSON(ap)
	dbus.Emit(m, "AccessPointAdded", string(devPath), apJSON)

	return
}
func (m *Manager) destroyAccessPoint(ap *accessPoint) {
	// emit AccessPointRemoved signal
	apJSON, _ := marshalJSON(ap)
	dbus.Emit(m, "AccessPointRemoved", string(ap.devPath), apJSON)
	nmDestroyAccessPoint(ap.nmAp)
}
func (a *accessPoint) updateProps() {
	a.Ssid = string(a.nmAp.Ssid.Get())
	a.Secured = getApSecType(a.nmAp) != apSecNone
	a.SecuredInEap = getApSecType(a.nmAp) == apSecEap
	a.Strength = a.nmAp.Strength.Get()
}
func getApSecType(ap *nm.AccessPoint) apSecType {
	return doParseApSecType(ap.Flags.Get(), ap.WpaFlags.Get(), ap.RsnFlags.Get())
}
func doParseApSecType(flags, wpaFlags, rsnFlags uint32) apSecType {
	r := apSecNone

	if (flags&NM_802_11_AP_FLAGS_PRIVACY != 0) && (wpaFlags == NM_802_11_AP_SEC_NONE) && (rsnFlags == NM_802_11_AP_SEC_NONE) {
		r = apSecWep
	}
	if wpaFlags != NM_802_11_AP_SEC_NONE {
		r = apSecPsk
	}
	if rsnFlags != NM_802_11_AP_SEC_NONE {
		r = apSecPsk
	}
	if (wpaFlags&NM_802_11_AP_SEC_KEY_MGMT_802_1X != 0) || (rsnFlags&NM_802_11_AP_SEC_KEY_MGMT_802_1X != 0) {
		r = apSecEap
	}
	return r
}

func (m *Manager) clearAccessPoints() {
	m.accessPointsLock.Lock()
	defer m.accessPointsLock.Unlock()
	for _, aps := range m.accessPoints {
		for _, ap := range aps {
			m.destroyAccessPoint(ap)
		}
	}
	m.accessPoints = make(map[dbus.ObjectPath][]*accessPoint)
}

func (m *Manager) addAccessPoint(devPath, apPath dbus.ObjectPath) {
	if m.isAccessPointExists(apPath) {
		return
	}

	m.accessPointsLock.Lock()
	defer m.accessPointsLock.Unlock()
	ap, err := m.newAccessPoint(devPath, apPath)
	if err != nil {
		return
	}
	logger.Debug("add access point", devPath, apPath)
	m.accessPoints[devPath] = append(m.accessPoints[devPath], ap)
}

func (m *Manager) removeAccessPoint(devPath, apPath dbus.ObjectPath) {
	if !m.isAccessPointExists(apPath) {
		return
	}
	devPath, i := m.getAccessPointIndex(apPath)

	m.accessPointsLock.Lock()
	defer m.accessPointsLock.Unlock()
	logger.Info("remove access point", devPath, apPath)
	m.accessPoints[devPath] = m.doRemoveAccessPoint(m.accessPoints[devPath], i)
}
func (m *Manager) doRemoveAccessPoint(aps []*accessPoint, i int) []*accessPoint {
	m.destroyAccessPoint(aps[i])
	copy(aps[i:], aps[i+1:])
	aps[len(aps)-1] = nil
	aps = aps[:len(aps)-1]
	return aps
}

func (m *Manager) getAccessPoint(apPath dbus.ObjectPath) (ap *accessPoint) {
	devPath, i := m.getAccessPointIndex(apPath)
	if i < 0 {
		logger.Warning("access point not found", apPath)
		return
	}

	m.accessPointsLock.Lock()
	defer m.accessPointsLock.Unlock()
	ap = m.accessPoints[devPath][i]
	return
}
func (m *Manager) isAccessPointExists(apPath dbus.ObjectPath) bool {
	_, i := m.getAccessPointIndex(apPath)
	if i >= 0 {
		return true
	}
	return false
}
func (m *Manager) getAccessPointIndex(apPath dbus.ObjectPath) (devPath dbus.ObjectPath, index int) {
	m.accessPointsLock.Lock()
	defer m.accessPointsLock.Unlock()
	for d, aps := range m.accessPoints {
		for i, ap := range aps {
			if ap.Path == apPath {
				return d, i
			}
		}
	}
	return "", -1
}

func (m *Manager) isSsidExists(devPath dbus.ObjectPath, ssid string) bool {
	for _, ap := range m.accessPoints[devPath] {
		if ap.Ssid == ssid {
			return true
		}
	}
	return false
}

// GetAccessPoints return all access points object which marshaled by json.
func (m *Manager) GetAccessPoints(path dbus.ObjectPath) (apsJSON string, err error) {
	m.accessPointsLock.Lock()
	defer m.accessPointsLock.Unlock()
	apsJSON, err = marshalJSON(m.accessPoints[path])
	return
}

// ActivateAccessPoint add and activate connection for access point.
func (m *Manager) ActivateAccessPoint(uuid string, apPath, devPath dbus.ObjectPath) (cpath dbus.ObjectPath, err error) {
	logger.Debugf("ActivateAccessPoint: uuid=%s, apPath=%s, devPath=%s", uuid, apPath, devPath)
	defer logger.Debugf("ActivateAccessPoint end") // TODO test

	if len(uuid) > 0 {
		cpath, err = m.ActivateConnection(uuid, devPath)
	} else {
		// if there is no connection for current access point, create one
		var nmAp *nm.AccessPoint
		nmAp, err = nmNewAccessPoint(apPath)
		if err != nil {
			return
		}
		defer nmDestroyAccessPoint(nmAp)

		uuid = utils.GenUuid()
		data := newWirelessConnectionData(string(nmAp.Ssid.Get()), uuid, []byte(nmAp.Ssid.Get()), getApSecType(nmAp))
		cpath, _, err = nmAddAndActivateConnection(data, devPath)
	}
	return
}

func (m *Manager) CreateConnectionForAccessPoint(apPath, devPath dbus.ObjectPath) (session *ConnectionSession, err error) {
	session, err = newConnectionSessionByCreate(connectionWireless, devPath)
	if err != nil {
		logger.Error(err)
		return
	}

	// setup access point data
	nmAp, err := nmNewAccessPoint(apPath)
	if err != nil {
		return
	}
	defer nmDestroyAccessPoint(nmAp)

	setSettingConnectionId(session.data, string(nmAp.Ssid.Get()))
	setSettingWirelessSsid(session.data, []byte(nmAp.Ssid.Get()))
	secType := getApSecType(nmAp)
	switch secType {
	case apSecNone:
		logicSetSettingVkWirelessSecurityKeyMgmt(session.data, "none")
	case apSecWep:
		logicSetSettingVkWirelessSecurityKeyMgmt(session.data, "wep")
	case apSecPsk:
		logicSetSettingVkWirelessSecurityKeyMgmt(session.data, "wpa-psk")
	case apSecEap:
		logicSetSettingVkWirelessSecurityKeyMgmt(session.data, "wpa-eap")
	}
	session.setProps()

	// install dbus session
	m.addConnectionSession(session)
	return
}
