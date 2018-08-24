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

package network

import (
	nmdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.networkmanager"

	"fmt"

	"pkg.deepin.io/dde/daemon/network/nm"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
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
	nmAp    *nmdbus.AccessPoint
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
	ap.nmAp.InitSignalExt(m.sysSigLoop, true)
	ap.nmAp.AccessPoint().ConnectPropertiesChanged(func(properties map[string]dbus.Variant) {
		if !m.isAccessPointExists(apPath) {
			return
		}

		m.accessPointsLock.Lock()
		defer m.accessPointsLock.Unlock()
		ignoredBefore := ap.shouldBeIgnore()
		ap.updateProps()
		ignoredNow := ap.shouldBeIgnore()
		apJSON, _ := marshalJSON(ap)
		if ignoredNow == ignoredBefore {
			// ignored state not changed, only send properties changed
			// signal when not ignored
			if ignoredNow {
				logger.Debugf("access point(ignored) properties changed %#v", ap)
			} else {
				//logger.Debugf("access point properties changed %#v", ap)
				m.service.Emit(m, "AccessPointPropertiesChanged", string(devPath), apJSON)
			}
		} else {
			// ignored state changed, if became ignored now, send
			// removed signal or send added signal
			if ignoredNow {
				logger.Debugf("access point is ignored %#v", ap)
				m.service.Emit(m, "AccessPointRemoved", string(devPath), apJSON)
			} else {
				logger.Debugf("ignored access point available %#v", ap)
				m.service.Emit(m, "AccessPointAdded", string(devPath), apJSON)
			}
		}
	})

	if ap.shouldBeIgnore() {
		logger.Debugf("new access point is ignored %#v", ap)
	} else {
		apJSON, _ := marshalJSON(ap)
		m.service.Emit(m, "AccessPointAdded", string(devPath), apJSON)
	}

	return
}
func (m *Manager) destroyAccessPoint(ap *accessPoint) {
	// emit AccessPointRemoved signal
	apJSON, _ := marshalJSON(ap)
	m.service.Emit(m, "AccessPointRemoved", string(ap.devPath), apJSON)
	nmDestroyAccessPoint(ap.nmAp)
}
func (a *accessPoint) updateProps() {
	ssid, _ := a.nmAp.Ssid().Get(0)
	a.Ssid = decodeSsid(ssid)
	a.Secured = getApSecType(a.nmAp) != apSecNone
	a.SecuredInEap = getApSecType(a.nmAp) == apSecEap
	a.Strength, _ = a.nmAp.Strength().Get(0)
}

func getApSecType(ap *nmdbus.AccessPoint) apSecType {
	flags, _ := ap.Flags().Get(0)
	wpaFlags, _ := ap.WpaFlags().Get(0)
	rsnFlags, _ := ap.RsnFlags().Get(0)
	return doParseApSecType(flags, wpaFlags, rsnFlags)
}

func doParseApSecType(flags, wpaFlags, rsnFlags uint32) apSecType {
	r := apSecNone

	if (flags&nm.NM_802_11_AP_FLAGS_PRIVACY != 0) && (wpaFlags == nm.NM_802_11_AP_SEC_NONE) && (rsnFlags == nm.NM_802_11_AP_SEC_NONE) {
		r = apSecWep
	}
	if wpaFlags != nm.NM_802_11_AP_SEC_NONE {
		r = apSecPsk
	}
	if rsnFlags != nm.NM_802_11_AP_SEC_NONE {
		r = apSecPsk
	}
	if (wpaFlags&nm.NM_802_11_AP_SEC_KEY_MGMT_802_1X != 0) || (rsnFlags&nm.NM_802_11_AP_SEC_KEY_MGMT_802_1X != 0) {
		r = apSecEap
	}
	return r
}

// Check if current access point should be ignore in front-end. Hide
// the access point that strength less than 10 (not include 0 which
// should be caused by the network driver issue) and not activated.
func (a *accessPoint) shouldBeIgnore() bool {
	if a.Strength < 10 && a.Strength != 0 &&
		!manager.isAccessPointActivated(a.devPath, a.Ssid) {
		return true
	}
	return false
}

func (m *Manager) isAccessPointActivated(devPath dbus.ObjectPath, ssid string) bool {
	for _, path := range nmGetActiveConnections() {
		aconn := m.newActiveConnection(path)
		if aconn.typ == nm.NM_SETTING_WIRELESS_SETTING_NAME && isDBusPathInArray(devPath, aconn.Devices) {
			if ssid == string(nmGetWirelessConnectionSsidByUuid(aconn.Uuid)) {
				return true
			}
		}
	}
	return false
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
	//logger.Debug("add access point", devPath, apPath)
	m.accessPoints[devPath] = append(m.accessPoints[devPath], ap)
}

func (m *Manager) removeAccessPoint(devPath, apPath dbus.ObjectPath) {
	if !m.isAccessPointExists(apPath) {
		return
	}
	devPath, i := m.getAccessPointIndex(apPath)

	m.accessPointsLock.Lock()
	defer m.accessPointsLock.Unlock()
	//logger.Debug("remove access point", devPath, apPath)
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
func (m *Manager) GetAccessPoints(path dbus.ObjectPath) (apsJSON string, busErr *dbus.Error) {
	m.accessPointsLock.Lock()
	defer m.accessPointsLock.Unlock()
	filteredAccessPoints := make([]*accessPoint, 0)
	for _, aconn := range m.accessPoints[path] {
		if !aconn.shouldBeIgnore() {
			filteredAccessPoints = append(filteredAccessPoints, aconn)
		}
	}
	apsJSON, err := marshalJSON(filteredAccessPoints)
	busErr = dbusutil.ToError(err)
	return
}

func (m *Manager) ActivateAccessPoint(uuid string, apPath, devPath dbus.ObjectPath) (cpath dbus.ObjectPath, busErr *dbus.Error) {
	var err error
	cpath, err = m.activateAccessPoint(uuid, apPath, devPath)
	busErr = dbusutil.ToError(err)
	return
}

// ActivateAccessPoint add and activate connection for access point.
func (m *Manager) activateAccessPoint(uuid string, apPath, devPath dbus.ObjectPath) (cpath dbus.ObjectPath, err error) {
	logger.Debugf("ActivateAccessPoint: uuid=%s, apPath=%s, devPath=%s", uuid, apPath, devPath)
	defer logger.Debugf("ActivateAccessPoint end")

	if len(uuid) > 0 {
		cpath, err = m.activateConnection(uuid, devPath)
	} else {
		// if there is no connection for current access point, create one
		var nmAp *nmdbus.AccessPoint
		nmAp, err = nmNewAccessPoint(apPath)
		if err != nil {
			return
		}

		uuid = utils.GenUuid()
		ssid, err := nmAp.Ssid().Get(0)
		if err != nil {
			logger.Warning("failed to get Ap Ssid:", err)
		}

		data := newWirelessConnectionData(decodeSsid(ssid), uuid, ssid, getApSecType(nmAp))
		cpath, _, err = nmAddAndActivateConnection(data, devPath, true)
		if err != nil {
			logger.Warning("temp debug err is:", err)
		}
	}
	return
}

// CreateConnectionForAccessPoint will crate connection for target
// access point, it will set the right SSID and secret method
// automatically.
func (m *Manager) CreateConnectionForAccessPoint(apPath, devPath dbus.ObjectPath) (
	sessionPath dbus.ObjectPath, busErr *dbus.Error) {
	session, err := m.createConnectionForAccessPoint(apPath, devPath)
	if err != nil {
		return "/", dbusutil.ToError(err)
	}
	return session.sessionPath, nil
}

func (m *Manager) createConnectionForAccessPoint(apPath, devPath dbus.ObjectPath) (
	session *ConnectionSession, err error) {
	session, err = newConnectionSessionByCreate(connectionWireless, devPath, m.service)
	if err != nil {
		logger.Error(err)
		return
	}

	// setup access point data
	nmAp, err := nmNewAccessPoint(apPath)
	if err != nil {
		return
	}

	ssid, _ := nmAp.Ssid().Get(0)
	setSettingConnectionId(session.data, decodeSsid(ssid))
	setSettingWirelessSsid(session.data, ssid)
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
	session.setPropAllowDelete(false)

	// install dbus session
	m.addConnectionSession(session)
	return
}
