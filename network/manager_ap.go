package main

import (
	nm "dbus/org/freedesktop/networkmanager"
	"dlib/dbus"
)

type apSecType uint32

const (
	apSecNone apSecType = iota
	apSecWep
	apSecPsk
	apSecEap
)

type accessPoint struct {
	Ssid         string
	Secured      bool
	SecuredInEap bool
	Strength     uint8
	Path         dbus.ObjectPath
	nmAp         *nm.AccessPoint
}

func newAccessPoint(apPath dbus.ObjectPath) (ap accessPoint, err error) {
	nmAp, err := nmNewAccessPoint(apPath)
	if err != nil {
		return
	}

	ap = accessPoint{
		Ssid:         string(nmAp.Ssid.Get()),
		Secured:      getApSecType(nmAp) != apSecNone,
		SecuredInEap: getApSecType(nmAp) == apSecEap,
		Strength:     calcApStrength(nmAp.Strength.Get()),
		Path:         nmAp.Path,
		nmAp:         nmAp,
	}
	return
}

func calcApStrength(s uint8) uint8 {
	switch {
	case s <= 10:
		return 0
	case s <= 25:
		return 25
	case s <= 50:
		return 50
	case s <= 75:
		return 75
	case s <= 100:
		return 100
	}
	return 0
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

func (m *Manager) addAccessPoint(devPath, apPath dbus.ObjectPath) {
	if m.isAccessPointExists(devPath, apPath) {
		return
	}
	ap, err := newAccessPoint(apPath)
	if err != nil {
		return
	}
	if len(ap.Ssid) == 0 {
		// ignore hidden access point
		return
	}

	// connect property, access point strength
	ap.nmAp.Strength.ConnectChanged(func() {
		// firstly, check if the access point is still exists to ignore
		// dbus error when getting property
		if m.isAccessPointExists(devPath, apPath) {
			ap.Strength = calcApStrength(ap.nmAp.Strength.Get())
			if m.AccessPointPropertiesChanged != nil {
				apJSON, _ := marshalJSON(ap)
				// logger.Debug(string(devPath), apJSON) // TODO test
				m.AccessPointPropertiesChanged(string(devPath), apJSON)
			}
			m.updatePropAccessPoints()
		}
	})

	// emit signal
	if m.AccessPointAdded != nil {
		apJSON, _ := marshalJSON(ap)
		m.AccessPointAdded(string(devPath), apJSON)
	}
	m.accessPoints[devPath] = append(m.accessPoints[devPath], &ap)
	m.updatePropAccessPoints()
}
func (m *Manager) removeAccessPoint(devPath, apPath dbus.ObjectPath) {
	// emit signal
	if m.AccessPointRemoved != nil {
		apJSON, _ := marshalJSON(accessPoint{Path: apPath})
		m.AccessPointRemoved(string(devPath), apJSON)
	}
	m.doRemoveAccessPoint(devPath, apPath)
	m.updatePropAccessPoints()
}
func (m *Manager) doRemoveAccessPoint(devPath, apPath dbus.ObjectPath) {
	i := m.getAccessPointIndex(devPath, apPath)
	if i < 0 {
		return
	}
	aps := m.accessPoints[devPath]
	copy(aps[i:], aps[i+1:])
	aps[len(aps)-1] = nil
	aps = aps[:len(aps)-1]
	m.accessPoints[devPath] = aps
}
func (m *Manager) isAccessPointExists(devPath, apPath dbus.ObjectPath) bool {
	if m.getAccessPointIndex(devPath, apPath) >= 0 {
		return true
	}
	return false
}
func (m *Manager) getAccessPointIndex(devPath, apPath dbus.ObjectPath) int {
	for i, ap := range m.accessPoints[devPath] {
		if ap.Path == apPath {
			return i
		}
	}
	return -1
}

// GetAccessPoints return all access points object which marshaled by json.
func (m *Manager) GetAccessPoints(path dbus.ObjectPath) (apsJSON string, err error) {
	aps, err := m.doGetAccessPoints(path)
	if err != nil {
		return
	}
	apsJSON, err = marshalJSON(aps)
	return
}
func (m *Manager) doGetAccessPoints(devPath dbus.ObjectPath) (aps []accessPoint, err error) {
	apPaths := nmGetAccessPoints(devPath)
	for _, path := range apPaths {
		ap, err := newAccessPoint(path)
		if err != nil {
			continue
		}
		if len(ap.Ssid) == 0 {
			// ignore hidden access point
			continue
		}
		aps = append(aps, ap)
	}
	return
}

// TODO remove
// GetAccessPointProperty return access point object which marshaled by json.
func (m *Manager) getAccessPointProperty(apPath dbus.ObjectPath) (apJSON string, err error) {
	ap, err := newAccessPoint(apPath)
	if err != nil {
		return
	}
	apJSON, err = marshalJSON(ap)
	return
}

// ActivateAccessPoint add and activate connection for access point.
func (m *Manager) ActivateAccessPoint(apPath, devPath dbus.ObjectPath) (uuid string, err error) {
	logger.Debugf("ActivateAccessPoint: apPath=%s, devPath=%s", apPath, devPath)
	// if there is no connection for current access point, create one
	ap, err := nmNewAccessPoint(apPath)
	if err != nil {
		return
	}
	cpath, ok := nmGetWirelessConnection(ap.Ssid.Get(), devPath)
	if ok {
		logger.Debug("activate access point", cpath) // TODO test
		uuid = nmGetConnectionUuid(cpath)
		_, err = nmActivateConnection(cpath, devPath)
	} else {
		logger.Debug("add and activate access point", cpath) // TODO test
		uuid = newUUID()
		data := newWirelessConnectionData(string(ap.Ssid.Get()), uuid, []byte(ap.Ssid.Get()), getApSecType(ap))
		_, _, err = nmAddAndActivateConnection(data, devPath)
	}
	return
}

func (m *Manager) CreateConnectionForAccessPoint(apPath, devPath dbus.ObjectPath) (session *ConnectionSession, err error) {
	session, err = NewConnectionSessionByCreate(typeWireless, devPath)
	if err != nil {
		logger.Error(err)
		return
	}

	// setup access point data
	ap, err := nmNewAccessPoint(apPath)
	if err != nil {
		return
	}
	setSettingConnectionId(session.data, string(ap.Ssid.Get()))
	setSettingWirelessSsid(session.data, []byte(ap.Ssid.Get()))
	secType := getApSecType(ap)
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

	// install dbus session
	err = dbus.InstallOnSession(session)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

// TODO remove
func (m *Manager) editConnectionForAccessPoint(apPath, devPath dbus.ObjectPath) (session *ConnectionSession, err error) {
	// session, err = NewConnectionSessionByOpen(uuid, devPath)
	// if err != nil {
	// 	logger.Error(err)
	// 	return
	// }

	// // install dbus session
	// err = dbus.InstallOnSession(session)
	// if err != nil {
	// 	logger.Error(err)
	// 	return
	// }

	return
}

// TODO remove
// GetConnectionUuidByAccessPoint return the connection's uuid of access point, return empty if none.
func (m *Manager) getConnectionUuidByAccessPoint(apPath dbus.ObjectPath) (uuid string, err error) {
	ap, err := nmNewAccessPoint(apPath)
	if err != nil {
		return
	}

	// TODO check wifi hw addr
	cpath, ok := nmGetWirelessConnection(ap.Ssid.Get(), "")
	if !ok {
		return
	}

	uuid = nmGetConnectionUuid(cpath)

	logger.Debugf("GetConnectionUuidByAccessPoint: apPath=%s, uuid=%s", apPath, uuid) // TODO test
	return
}
