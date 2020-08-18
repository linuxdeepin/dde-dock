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
	"errors"
	"fmt"

	dbus "github.com/godbus/dbus"
	nmdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.networkmanager"
	"pkg.deepin.io/dde/daemon/network/nm"
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

func (v apSecType) String() string {
	switch v {
	case apSecNone:
		return "none"
	case apSecWep:
		return "wep"
	case apSecPsk:
		return "wpa-psk"
	case apSecEap:
		return "wpa-eap"
	default:
		return fmt.Sprintf("<invalid apSecType %d>", v)
	}
}

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
	_, err = ap.nmAp.AccessPoint().ConnectPropertiesChanged(func(properties map[string]dbus.Variant) {
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
				err = m.service.Emit(m, "AccessPointPropertiesChanged", string(devPath), apJSON)
			}
		} else {
			// ignored state changed, if became ignored now, send
			// removed signal or send added signal
			if ignoredNow {
				logger.Debugf("access point is ignored %#v", ap)
				err = m.service.Emit(m, "AccessPointRemoved", string(devPath), apJSON)
			} else {
				logger.Debugf("ignored access point available %#v", ap)
				err = m.service.Emit(m, "AccessPointAdded", string(devPath), apJSON)
			}
		}
		if err != nil {
			logger.Warning("failed to emit signal:", err)
		}
	})
	if err != nil {
		logger.Warning("failed to monitor changing properties of AccessPoint", err)
	}

	if ap.shouldBeIgnore() {
		logger.Debugf("new access point is ignored %#v", ap)
	} else {
		apJSON, _ := marshalJSON(ap)
		err = m.service.Emit(m, "AccessPointAdded", string(devPath), apJSON)
		if err != nil {
			logger.Warning("failed to emit signal:", err)
		}
	}

	return
}

func (m *Manager) destroyAccessPoint(ap *accessPoint) {
	// emit AccessPointRemoved signal
	apJSON, _ := marshalJSON(ap)
	err := m.service.Emit(m, "AccessPointRemoved", string(ap.devPath), apJSON)
	if err != nil {
		logger.Warning("failed to emit signal:", err)
	}
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
	_, i := m.getAccessPointIndex(apPath)

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

func (m *Manager) isAccessPointExists(apPath dbus.ObjectPath) bool {
	_, i := m.getAccessPointIndex(apPath)
	return i >= 0
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

// GetAccessPoints return all access points object which marshaled by json.
func (m *Manager) GetAccessPoints(path dbus.ObjectPath) (apsJSON string, busErr *dbus.Error) {
	m.accessPointsLock.Lock()
	defer m.accessPointsLock.Unlock()
	accessPoints := m.accessPoints[path]
	filteredAccessPoints := make([]*accessPoint, 0, len(m.accessPoints))
	for _, ap := range accessPoints {
		if !ap.shouldBeIgnore() {
			filteredAccessPoints = append(filteredAccessPoints, ap)
		}
	}
	apsJSON, err := marshalJSON(filteredAccessPoints)
	busErr = dbusutil.ToError(err)
	return
}

func (m *Manager) ActivateAccessPoint(uuid string, apPath, devPath dbus.ObjectPath) (dbus.ObjectPath,
	*dbus.Error) {
	var err error
	cpath, err := m.activateAccessPoint(uuid, apPath, devPath)
	if err != nil {
		logger.Warning("failed to activate access point:", err)
		return "/", dbusutil.ToError(err)
	}
	return cpath, nil
}

func fixApSecTypeChange(uuid string, secType apSecType) (needUserEdit bool, err error) {
	var cpath dbus.ObjectPath
	cpath, err = nmGetConnectionByUuid(uuid)
	if err != nil {
		return
	}

	var conn *nmdbus.ConnectionSettings
	conn, err = nmNewSettingsConnection(cpath)
	if err != nil {
		return
	}
	var connData connectionData
	connData, err = conn.GetSettings(0)
	if err != nil {
		return
	}

	secTypeOld, err := getApSecTypeFromConnData(connData)
	if err != nil {
		logger.Warning("failed to get apSecType from connData")
		return false, nil
	}

	if secTypeOld == secType {
		return
	}
	logger.Debug("apSecType change to", secType)

	switch secType {
	case apSecNone:
		err = logicSetSettingVkWirelessSecurityKeyMgmt(connData, "none")
	case apSecWep:
		err = logicSetSettingVkWirelessSecurityKeyMgmt(connData, "wep")
	case apSecPsk:
		err = logicSetSettingVkWirelessSecurityKeyMgmt(connData, "wpa-psk")
	case apSecEap:
		needUserEdit = true
		return
	}
	if err != nil {
		logger.Debug("failed to set VKWirelessSecutiryKeyMgmt")
		return
	}

	// fix ipv6 addresses and routes data structure, interface{}
	if isSettingIP6ConfigAddressesExists(connData) {
		setSettingIP6ConfigAddresses(connData, getSettingIP6ConfigAddresses(connData))
	}
	if isSettingIP6ConfigRoutesExists(connData) {
		setSettingIP6ConfigRoutes(connData, getSettingIP6ConfigRoutes(connData))
	}

	err = conn.Update(0, connData)
	return
}

// ActivateAccessPoint add and activate connection for access point.
func (m *Manager) activateAccessPoint(uuid string, apPath, devPath dbus.ObjectPath) (cpath dbus.ObjectPath, err error) {
	logger.Debugf("ActivateAccessPoint: uuid=%s, apPath=%s, devPath=%s", uuid, apPath, devPath)

	cpath = "/"
	var nmAp *nmdbus.AccessPoint
	nmAp, err = nmNewAccessPoint(apPath)
	if err != nil {
		return
	}
	secType := getApSecType(nmAp)
	if uuid != "" {
		var needUserEdit bool
		needUserEdit, err = fixApSecTypeChange(uuid, secType)
		if err != nil {
			return
		}
		if needUserEdit {
			err = errors.New("need user edit")
			return
		}
		cpath, err = m.activateConnection(uuid, devPath)
	} else {
		// if there is no connection for current access point, create one
		uuid = utils.GenUuid()
		var ssid []byte
		ssid, err = nmAp.Ssid().Get(0)
		if err != nil {
			logger.Warning("failed to get Ap Ssid:", err)
			return
		}

		data := newWirelessConnectionData(decodeSsid(ssid), uuid, ssid, secType)
		cpath, _, err = nmAddAndActivateConnection(data, devPath, true)
		if err != nil {
			return
		}
	}
	return
}
