package network

import (
	"dlib/dbus"
	"fmt"
)

type sectionErrors map[string]string
type sessionErrors map[string]sectionErrors

type ConnectionSession struct {
	sessionPath dbus.ObjectPath
	devPath     dbus.ObjectPath
	data        connectionData

	ConnectionPath dbus.ObjectPath
	Uuid           string
	Type           string

	AvailableSections []string
	AvailableKeys     map[string][]string // TODO collection of available keys in sections(not virtual section)

	Errors           sessionErrors
	settingKeyErrors sessionErrors
}

func doNewConnectionSession(devPath dbus.ObjectPath, uuid string) (s *ConnectionSession) {
	s = &ConnectionSession{}
	s.sessionPath = dbus.ObjectPath(fmt.Sprintf("/com/deepin/daemon/ConnectionSession/%s", randString(8)))
	s.devPath = devPath
	s.Uuid = uuid
	s.data = make(connectionData)
	s.AvailableSections = make([]string, 0)
	s.AvailableKeys = make(map[string][]string)
	s.Errors = make(sessionErrors)
	s.settingKeyErrors = make(sessionErrors)
	return s
}

func newConnectionSessionByCreate(connectionType string, devPath dbus.ObjectPath) (s *ConnectionSession, err error) {
	if !isStringInArray(connectionType, supportedConnectionTypes) {
		err = fmt.Errorf("connection type is out of support: %s", connectionType)
		logger.Error(err)
		return
	}

	s = doNewConnectionSession(devPath, newUUID())

	s.Type = connectionType
	id := genConnectionId(s.Type)
	switch s.Type {
	case connectionWired:
		s.data = newWiredConnectionData(id, s.Uuid)
	case connectionWireless:
		s.data = newWirelessConnectionData(id, s.Uuid, nil, apSecNone)
	case connectionWirelessAdhoc:
		s.data = newWirelessAdhocConnectionData(id, s.Uuid)
	case connectionWirelessHotspot:
		s.data = newWirelessHotspotConnectionData(id, s.Uuid)
	case connectionPppoe:
		s.data = newPppoeConnectionData(id, s.Uuid)
	case connectionMobileGsm:
		s.data = newMobileConnectionData(id, s.Uuid, mobileServiceGsm)
	case connectionMobileCdma:
		s.data = newMobileConnectionData(id, s.Uuid, mobileServiceCdma)
	case connectionVpnL2tp:
		s.data = newVpnL2tpConnectionData(id, s.Uuid)
	case connectionVpnOpenconnect:
		s.data = newVpnOpenconnectConnectionData(id, s.Uuid)
	case connectionVpnPptp:
		s.data = newVpnPptpConnectionData(id, s.Uuid)
	case connectionVpnVpnc:
		s.data = newVpnVpncConnectionData(id, s.Uuid)
	case connectionVpnOpenvpn:
		s.data = newVpnOpenvpnConnectionData(id, s.Uuid)
	}

	s.updatePropConnectionType()
	s.updatePropAvailableSections()
	s.updatePropAvailableKeys()
	s.updatePropErrors()

	logger.Debug("newConnectionSessionByCreate():", s.data)
	return
}

func newConnectionSessionByOpen(uuid string, devPath dbus.ObjectPath) (s *ConnectionSession, err error) {
	connectionPath, err := nmGetConnectionByUuid(uuid)
	if err != nil {
		return
	}

	s = doNewConnectionSession(devPath, uuid)
	s.ConnectionPath = connectionPath

	// get connection data
	s.data, err = nmGetConnectionData(s.ConnectionPath)
	if err != nil {
		return nil, err
	}
	s.Type = getCustomConnectionType(s.data)

	s.fixValues()

	// execute asynchronous to avoid front-end blocked if
	// NeedSecrets() signal send to it
	go s.getSecrets()

	s.updatePropConnectionType()
	s.updatePropAvailableSections()
	s.updatePropAvailableKeys()
	s.updatePropErrors()

	logger.Debug("NewConnectionSessionByOpen():", s.data)
	return
}

func (s *ConnectionSession) fixValues() {
	// append missing sectionIpv6
	if !isSettingSectionExists(s.data, sectionIpv6) && isStringInArray(sectionIpv6, getAvailableSections(s.data)) {
		initSettingSectionIpv6(s.data)
	}

	// vpn plugin data and secret
	if getSettingConnectionType(s.data) == NM_SETTING_VPN_SETTING_NAME {
		if !isSettingVpnDataExists(s.data) {
			setSettingVpnData(s.data, make(map[string]string))
		}
		if !isSettingVpnSecretsExists(s.data) {
			setSettingVpnSecrets(s.data, make(map[string]string))
		}
	}

	// append missing sectionWired for pppoe
	if s.Type == connectionPppoe {
		if !isSettingSectionExists(s.data, sectionWired) {
			initSettingSectionWired(s.data)
		}
	}

	// TODO fix secret flags
	// if isSettingVpnOpenvpnKeyCertpassFlagsExists(s.data) && getSettingVpnOpenvpnKeyCertpassFlags(s.data) == 1 {
	// setSettingVpnOpenvpnKeyCertpassFlags(s.data, NM_OPENVPN_SECRET_FLAG_SAVE)
	// }
}

func (s *ConnectionSession) getSecrets() {
	// get secret data
	switch s.Type {
	case connectionWired:
		if getSettingVk8021xEnable(s.data) {
			s.doGetSecrets(section8021x)
		}
	case connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
		if getSettingVk8021xEnable(s.data) {
			s.doGetSecrets(section8021x)
		} else {
			s.doGetSecrets(sectionWirelessSecurity)
		}
	case connectionPppoe:
		s.doGetSecrets(sectionPppoe)
	case connectionMobileGsm:
		s.doGetSecrets(sectionGsm)
	case connectionMobileCdma:
		s.doGetSecrets(sectionCdma)
	case connectionVpnL2tp, connectionVpnOpenconnect, connectionVpnPptp, connectionVpnVpnc, connectionVpnOpenvpn:
		// TODO
	}
}

func (s *ConnectionSession) doGetSecrets(secretField string) {
	if isSettingSectionExists(s.data, secretField) {
		secrets, err := nmGetConnectionSecrets(s.ConnectionPath, secretField)
		if err == nil {
			for section, sectionData := range secrets {
				if !isSettingSectionExists(s.data, section) {
					addSettingSection(s.data, section)
				}
				for key, value := range sectionData {
					s.data[section][key] = value
				}
			}
		}
	}
}

// Save save current connection s.
func (s *ConnectionSession) Save() bool {
	// TODO what about the connection has been deleted?

	if s.isErrorOccured() {
		logger.Debug("Errors occured when saving:", s.Errors)
		return false
	}

	if getSettingConnectionReadOnly(s.data) {
		logger.Debug("read only connection, don't allowed to save")
		return false
	}

	if len(s.ConnectionPath) > 0 {
		// update connection data and activate it
		nmConn, err := nmNewSettingsConnection(s.ConnectionPath)
		if err != nil {
			logger.Error(err)
			return false
		}
		err = nmConn.Update(s.data)
		if err != nil {
			logger.Error(err)
			return false
		}
		nmActivateConnection(s.ConnectionPath, s.devPath)
	} else {
		// create new connection and activate it
		// TODO vpn ad-hoc hotspot
		if s.Type == connectionWired || s.Type == connectionWireless {
			nmAddAndActivateConnection(s.data, s.devPath)
		} else {
			nmAddConnection(s.data)
		}
	}

	removeConnectionSession(s)
	return true
}

func (s *ConnectionSession) isErrorOccured() bool {
	for _, v := range s.Errors {
		if len(v) > 0 {
			return true
		}
	}
	return false
}

// Close cancel current connection.
func (s *ConnectionSession) Close() {
	removeConnectionSession(s)
}

// GetAvailableValues return available values marshaled by json for target key.
func (s *ConnectionSession) GetAvailableValues(vsection, key string) (valuesJSON string) {
	var values []kvalue
	sections := getRelatedSectionsOfVsection(s.data, vsection)
	for _, section := range sections {
		values = generalGetSettingAvailableValues(s.data, section, key)
		if len(values) > 0 {
			break
		}
	}
	valuesJSON, _ = marshalJSON(values)
	return
}

func (s *ConnectionSession) GetKey(vsection, key string) (value string) {
	section := getSectionOfKeyInVsection(s.data, vsection, key)
	value = generalGetSettingKeyJSON(s.data, section, key)
	return
}

func (s *ConnectionSession) SetKey(vsection, key, value string) {
	section := getSectionOfKeyInVsection(s.data, vsection, key)
	err := generalSetSettingKeyJSON(s.data, section, key, value)
	// logger.Debugf("SetKey(), %v, vsection=%s, filed=%s, key=%s, value=%s", err == nil, vsection, section, key, value) // TODO test
	s.updateErrorsWhenSettingKey(vsection, key, err)

	s.updatePropAvailableSections()
	s.updatePropAvailableKeys()
	s.updatePropErrors()

	return
}

func (s *ConnectionSession) updateErrorsWhenSettingKey(vsection, key string, err error) {
	if err == nil {
		// delete key error if exists
		sectionErrors, ok := s.settingKeyErrors[vsection]
		if ok {
			_, ok := sectionErrors[key]
			if ok {
				delete(sectionErrors, key)
			}
		}
	} else {
		// append key error
		sectionErrorsData, ok := s.settingKeyErrors[vsection]
		if !ok {
			sectionErrorsData = make(sectionErrors)
			s.settingKeyErrors[vsection] = sectionErrorsData
		}
		sectionErrorsData[key] = err.Error()
	}
}

// Debug functions
func (s *ConnectionSession) DebugGetConnectionData() connectionData {
	return s.data
}
func (s *ConnectionSession) DebugGetErrors() sessionErrors {
	return s.Errors
}
func (s *ConnectionSession) DebugListKeyDetail() (info string) {
	for _, vsection := range getAvailableVsections(s.data) {
		vsectionData, ok := s.AvailableKeys[vsection]
		if !ok {
			logger.Warning("no available keys for vsection", vsection)
			continue
		}
		for _, key := range vsectionData {
			section := getSectionOfKeyInVsection(s.data, vsection, key)
			t := generalGetSettingKeyType(section, key)
			values := generalGetSettingAvailableValues(s.data, section, key)
			valuesJSON, _ := marshalJSON(values)
			info += fmt.Sprintf("%s->%s[%s]: %s (%s)\n", vsection, key, getKtypeDescription(t), s.GetKey(vsection, key), valuesJSON)
		}
	}
	return
}
