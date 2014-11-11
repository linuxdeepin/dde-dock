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
	"fmt"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/utils"
	"time"
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

	AvailableVirtualSections []string
	AvailableSections        []string

	// collection of available keys in each section(not virtual section)
	AvailableKeys map[string][]string

	Errors           sessionErrors
	settingKeyErrors sessionErrors

	// signal
	ConnectionDataChanged func()
}

func doNewConnectionSession(devPath dbus.ObjectPath, uuid string) (s *ConnectionSession) {
	s = &ConnectionSession{}
	s.sessionPath = dbus.ObjectPath(fmt.Sprintf("/com/deepin/daemon/ConnectionSession/%s", utils.RandString(8)))
	s.devPath = devPath
	s.Uuid = uuid
	s.data = make(connectionData)
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

	s = doNewConnectionSession(devPath, utils.GenUuid())

	// expand wrapper connection type
	id := genConnectionId(connectionType)
	switch connectionType {
	case connectionMobile:
		connectionType = connectionMobileGsm
	case connectionVpn:
		connectionType = connectionVpnL2tp
	}
	switch connectionType {
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
		s.data = newMobileConnectionData(id, s.Uuid, connectionMobileGsm)
	case connectionMobileCdma:
		s.data = newMobileConnectionData(id, s.Uuid, connectionMobileCdma)
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

	s.setProps()
	logger.Infof("newConnectionSessionByCreate(): %#v", s.data)
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

	s.setProps()
	logger.Infof("NewConnectionSessionByOpen(): %#v", s.data)
	return
}

func (s *ConnectionSession) fixValues() {
	// append missing sectionIpv6
	if !isSettingSectionExists(s.data, sectionIpv6) && isStringInArray(sectionIpv6, getAvailableSections(s.data)) {
		initSettingSectionIpv6(s.data)
	}

	// fix ipv6 addresses and routes data structure, interface{}
	if isSettingIp6ConfigAddressesExists(s.data) {
		setSettingIp6ConfigAddresses(s.data, getSettingIp6ConfigAddresses(s.data))
	}
	if isSettingIp6ConfigRoutesExists(s.data) {
		setSettingIp6ConfigRoutes(s.data, getSettingIp6ConfigRoutes(s.data))
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
	if getCustomConnectionType(s.data) == connectionPppoe {
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
	switch getCustomConnectionType(s.data) {
	case connectionWired:
		if getSettingVk8021xEnable(s.data) {
			// TODO 8021x secret
			// s.doGetSecrets(section8021x)
		}
	case connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
		if getSettingVk8021xEnable(s.data) {
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
func (s *ConnectionSession) Save() (ok bool, err error) {
	// TODO what about the connection has been deleted?

	logger.Infof("Save connection: %#v", s.data)

	if s.isErrorOccured() {
		logger.Debug("Errors:", s.Errors)
		logger.Debug("settingKeyErrors:", s.settingKeyErrors)
		return false, nil
	}

	if getSettingConnectionReadOnly(s.data) {
		err = fmt.Errorf("read only connection, don't allowed to save")
		logger.Debug(err)
		return false, err
	}

	if len(s.ConnectionPath) > 0 {
		// update connection data and activate it
		nmConn, err := nmNewSettingsConnection(s.ConnectionPath)
		if err != nil {
			return false, err
		}
		err = nmConn.Update(s.data)
		if err != nil {
			logger.Error(err)
			return false, err
		}
		manager.ActivateConnection(s.Uuid, s.devPath)
	} else {
		// create new connection and activate it
		// TODO vpn ad-hoc hotspot
		connectionType := getCustomConnectionType(s.data)
		switch connectionType {
		case connectionWired, connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
			nmAddAndActivateConnection(s.data, s.devPath)
		default:
			nmAddConnection(s.data)
		}
	}

	manager.removeConnectionSession(s)
	return true, nil
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
	manager.removeConnectionSession(s)
}

// GetAvailableValues return available values marshaled by json for target key.
func (s *ConnectionSession) GetAvailableValues(section, key string) (valuesJSON string) {
	var values []kvalue
	values = generalGetSettingAvailableValues(s.data, section, key)
	valuesJSON, _ = marshalJSON(values)
	return
}

func (s *ConnectionSession) GetKey(section, key string) (valueJSON string) {
	// logger.Debugf("GetKey(), section=%s, key=%s", section, key) // TODO test
	valueJSON = generalGetSettingKeyJSON(s.data, section, key)
	// logger.Debugf("GetKey(), section=%s, key=%s, valueJSON=%s", section, key, valueJSON) // TODO test
	return
}

func (s *ConnectionSession) SetKey(section, key, valueJSON string) {
	logger.Debugf("SetKey(), section=%s, key=%s, valueJSON=%s", section, key, valueJSON) // TODO test
	err := generalSetSettingKeyJSON(s.data, section, key, valueJSON)
	s.updateErrorsWhenSettingKey(section, key, err)
	s.setProps()
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

func (s *ConnectionSession) IsDefaultExpandedSection(vsection string) (bool, error) {
	switch s.Type {
	case connectionWired:
		switch vsection {
		case vsectionIpv4:
			return true, nil
		}
	case connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
		switch vsection {
		case vsectionSecurity:
			return true, nil
		}
	case connectionPppoe:
		switch vsection {
		case vsectionPppoe:
			return true, nil
		}
	case connectionMobileGsm, connectionMobileCdma:
		switch vsection {
		case vsectionMobile:
			return true, nil
		}
	case connectionVpnL2tp, connectionVpnOpenconnect, connectionVpnPptp, connectionVpnVpnc, connectionVpnOpenvpn:
		switch vsection {
		case vsectionVpn:
			return true, nil
		}
	default:
		switch vsection {
		case vsectionIpv4:
			return true, nil
		}
	}
	return false, nil
}

// Debug functions
func (s *ConnectionSession) DebugGetConnectionData() connectionData {
	return s.data
}
func (s *ConnectionSession) DebugGetErrors() sessionErrors {
	return s.Errors
}
func (s *ConnectionSession) DebugListKeyDetail() (info string) {
	// TODO
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
