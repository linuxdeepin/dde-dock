package network

import (
	"dlib/dbus"
	"fmt"
	"time"
)

type sectionErrors map[string]string
type sessionErrors map[string]sectionErrors

type ConnectionSession struct {
	sessionPath dbus.ObjectPath
	devPath     dbus.ObjectPath

	ConnectionPath dbus.ObjectPath
	Uuid           string
	Type           string
	Data           connectionData

	AvailableVirtualSections []string
	AvailableSections        []string

	// collection of available keys in each section(not virtual section)
	AvailableKeys map[string][]string

	Errors           sessionErrors
	settingKeyErrors sessionErrors
}

func doNewConnectionSession(devPath dbus.ObjectPath, uuid string) (s *ConnectionSession) {
	s = &ConnectionSession{}
	s.sessionPath = dbus.ObjectPath(fmt.Sprintf("/com/deepin/daemon/ConnectionSession/%s", randString(8)))
	s.devPath = devPath
	s.Uuid = uuid
	s.Data = make(connectionData)
	s.AvailableVirtualSections = make([]string, 0)
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
		s.Data = newWiredConnectionData(id, s.Uuid)
	case connectionWireless:
		s.Data = newWirelessConnectionData(id, s.Uuid, nil, apSecNone)
	case connectionWirelessAdhoc:
		s.Data = newWirelessAdhocConnectionData(id, s.Uuid)
	case connectionWirelessHotspot:
		s.Data = newWirelessHotspotConnectionData(id, s.Uuid)
	case connectionPppoe:
		s.Data = newPppoeConnectionData(id, s.Uuid)
	case connectionMobileGsm:
		s.Data = newMobileConnectionData(id, s.Uuid, mobileServiceGsm)
	case connectionMobileCdma:
		s.Data = newMobileConnectionData(id, s.Uuid, mobileServiceCdma)
	case connectionVpnL2tp:
		s.Data = newVpnL2tpConnectionData(id, s.Uuid)
	case connectionVpnOpenconnect:
		s.Data = newVpnOpenconnectConnectionData(id, s.Uuid)
	case connectionVpnPptp:
		s.Data = newVpnPptpConnectionData(id, s.Uuid)
	case connectionVpnVpnc:
		s.Data = newVpnVpncConnectionData(id, s.Uuid)
	case connectionVpnOpenvpn:
		s.Data = newVpnOpenvpnConnectionData(id, s.Uuid)
	}

	s.updatePropConnectionType()
	s.updatePropAvailableVirtualSections()
	s.updatePropAvailableSections()
	s.updatePropAvailableKeys()
	s.updatePropErrors()

	logger.Debug("newConnectionSessionByCreate():", s.Data)
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
	s.Data, err = nmGetConnectionData(s.ConnectionPath)
	if err != nil {
		return nil, err
	}
	s.Type = getCustomConnectionType(s.Data)

	s.fixValues()

	// execute asynchronous to avoid front-end block if
	// NeedSecrets() signal emit
	chSecret := make(chan int)
	go func() {
		s.getSecrets()
		chSecret <- 0
	}()
	select {
	case <-time.After(500 * time.Millisecond):
	case <-chSecret:
	}

	s.updatePropConnectionType()
	s.updatePropAvailableVirtualSections()
	s.updatePropAvailableSections()
	s.updatePropAvailableKeys()
	s.updatePropErrors()

	logger.Debug("NewConnectionSessionByOpen():", s.Data)
	return
}

func (s *ConnectionSession) fixValues() {
	// append missing sectionIpv6
	if !isSettingSectionExists(s.Data, sectionIpv6) && isStringInArray(sectionIpv6, getAvailableSections(s.Data)) {
		initSettingSectionIpv6(s.Data)
	}

	// vpn plugin data and secret
	if getSettingConnectionType(s.Data) == NM_SETTING_VPN_SETTING_NAME {
		if !isSettingVpnDataExists(s.Data) {
			setSettingVpnData(s.Data, make(map[string]string))
		}
		if !isSettingVpnSecretsExists(s.Data) {
			setSettingVpnSecrets(s.Data, make(map[string]string))
		}
	}

	// append missing sectionWired for pppoe
	if s.Type == connectionPppoe {
		if !isSettingSectionExists(s.Data, sectionWired) {
			initSettingSectionWired(s.Data)
		}
	}

	// TODO fix secret flags
	// if isSettingVpnOpenvpnKeyCertpassFlagsExists(s.Data) && getSettingVpnOpenvpnKeyCertpassFlags(s.Data) == 1 {
	// setSettingVpnOpenvpnKeyCertpassFlags(s.Data, NM_OPENVPN_SECRET_FLAG_SAVE)
	// }
}

func (s *ConnectionSession) getSecrets() {
	// get secret data
	switch s.Type {
	case connectionWired:
		if getSettingVk8021xEnable(s.Data) {
			// TODO 8021x secret
			// s.doGetSecrets(section8021x)
		}
	case connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
		if getSettingVk8021xEnable(s.Data) {
			// TODO 8021x secret
			// s.doGetSecrets(section8021x)
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
	if isSettingSectionExists(s.Data, secretField) {
		secrets, err := nmGetConnectionSecrets(s.ConnectionPath, secretField)
		if err == nil {
			for section, sectionData := range secrets {
				if !isSettingSectionExists(s.Data, section) {
					addSettingSection(s.Data, section)
				}
				for key, value := range sectionData {
					s.Data[section][key] = value
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

	if getSettingConnectionReadOnly(s.Data) {
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
		err = nmConn.Update(s.Data)
		if err != nil {
			logger.Error(err)
			return false
		}
		nmActivateConnection(s.ConnectionPath, s.devPath)
	} else {
		// create new connection and activate it
		// TODO vpn ad-hoc hotspot
		if s.Type == connectionWired || s.Type == connectionWireless {
			nmAddAndActivateConnection(s.Data, s.devPath)
		} else {
			nmAddConnection(s.Data)
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
func (s *ConnectionSession) GetAvailableValues(section, key string) (valuesJSON string) {
	// TODO remove
	// var values []kvalue
	// sections := getRelatedSectionsOfVsection(s.Data, vsection)
	// for _, section := range sections {
	// 	values = generalGetSettingAvailableValues(s.Data, section, key)
	// 	if len(values) > 0 {
	// 		break
	// 	}
	// }
	values := generalGetSettingAvailableValues(s.Data, section, key)
	valuesJSON, _ = marshalJSON(values)
	return
}

func (s *ConnectionSession) GetKey(section, key string) (value string) {
	// section := getSectionOfKeyInVsection(s.Data, vsection, key)
	value = generalGetSettingKeyJSON(s.Data, section, key)
	return
}

func (s *ConnectionSession) SetKey(section, key, value string) {
	// section := getSectionOfKeyInVsection(s.Data, vsection, key)
	err := generalSetSettingKeyJSON(s.Data, section, key, value)
	// logger.Debugf("SetKey(), %v, vsection=%s, filed=%s, key=%s, value=%s", err == nil, vsection, section, key, value) // TODO test
	s.updateErrorsWhenSettingKey(section, key, err)

	s.updatePropAvailableVirtualSections()
	s.updatePropAvailableSections()
	s.updatePropAvailableKeys()
	s.updatePropErrors()

	return
}

func (s *ConnectionSession) updateErrorsWhenSettingKey(section, key string, err error) {
	if err == nil {
		// delete key error if exists
		sectionErrors, ok := s.settingKeyErrors[section]
		if ok {
			_, ok := sectionErrors[key]
			if ok {
				delete(sectionErrors, key)
			}
		}
	} else {
		// append key error
		sectionErrorsData, ok := s.settingKeyErrors[section]
		if !ok {
			sectionErrorsData = make(sectionErrors)
			s.settingKeyErrors[section] = sectionErrorsData
		}
		sectionErrorsData[key] = err.Error()
	}
}

// Debug functions
func (s *ConnectionSession) DebugGetConnectionData() connectionData {
	return s.Data
}
func (s *ConnectionSession) DebugGetErrors() sessionErrors {
	return s.Errors
}
func (s *ConnectionSession) DebugListKeyDetail() (info string) {
	for _, vsection := range getAvailableVsections(s.Data) {
		vsectionData, ok := s.AvailableKeys[vsection]
		if !ok {
			logger.Warning("no available keys for vsection", vsection)
			continue
		}
		for _, key := range vsectionData {
			section := getSectionOfKeyInVsection(s.Data, vsection, key)
			t := generalGetSettingKeyType(section, key)
			values := generalGetSettingAvailableValues(s.Data, section, key)
			valuesJSON, _ := marshalJSON(values)
			info += fmt.Sprintf("%s->%s[%s]: %s (%s)\n", vsection, key, getKtypeDescription(t), s.GetKey(vsection, key), valuesJSON)
		}
	}
	return
}
