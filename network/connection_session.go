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
	"pkg.linuxdeepin.com/lib/dbus"
	"fmt"
	"time"
)

type sectionErrors map[string]string
type sessionErrors map[string]sectionErrors

type ConnectionSession struct {
	sessionPath dbus.ObjectPath
	devPath     dbus.ObjectPath

	Data           connectionData
	ConnectionPath dbus.ObjectPath
	Uuid           string
	Type           string

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

	s = doNewConnectionSession(devPath, genUuid())

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
		s.Data = newMobileConnectionData(id, s.Uuid, connectionMobileGsm)
	case connectionMobileCdma:
		s.Data = newMobileConnectionData(id, s.Uuid, connectionMobileCdma)
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

	s.updateProps()
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

	s.updateProps()
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
	if getCustomConnectionType(s.Data) == connectionPppoe {
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
	switch getCustomConnectionType(s.Data) {
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
func (s *ConnectionSession) Save() (ok bool, err error) {
	// TODO what about the connection has been deleted?

	if s.isErrorOccured() {
		logger.Debug("Errors occured when saving:", s.Errors)
		return false, nil
	}

	if getSettingConnectionReadOnly(s.Data) {
		err = fmt.Errorf("read only connection, don't allowed to save")
		logger.Debug(err)
		return false, err
	}

	if len(s.ConnectionPath) > 0 {
		// update connection data and activate it
		nmConn, err := nmNewSettingsConnection(s.ConnectionPath)
		if err != nil {
			logger.Error(err)
			return false, err
		}
		err = nmConn.Update(s.Data)
		if err != nil {
			logger.Error(err)
			return false, err
		}
		nmActivateConnection(s.ConnectionPath, s.devPath)
	} else {
		// create new connection and activate it
		// TODO vpn ad-hoc hotspot
		connectionType := getCustomConnectionType(s.Data)
		switch connectionType {
		case connectionWired, connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
			nmAddAndActivateConnection(s.Data, s.devPath)
		default:
			nmAddConnection(s.Data)
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
	values = generalGetSettingAvailableValues(s.Data, section, key)
	valuesJSON, _ = marshalJSON(values)
	return
}

func (s *ConnectionSession) GetKey(section, key string) (valueJSON string) {
	// section := getSectionOfKeyInVsection(s.Data, vsection, key)
	valueJSON = generalGetSettingKeyJSON(s.Data, section, key)
	return
}

func (s *ConnectionSession) SetKey(section, key, valueJSON string) {
	// logger.Debugf("SetKey(), section=%s, key=%s, valueJSON=%s", section, key, valueJSON) // TODO test
	err := generalSetSettingKeyJSON(s.Data, section, key, valueJSON)
	s.updateErrorsWhenSettingKey(section, key, err)
	s.updateProps()
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

	// ignore errors that not available
	for errSection, _ := range s.settingKeyErrors {
		if !isStringInArray(errSection, s.AvailableSections) {
			delete(s.settingKeyErrors, errSection)
		}
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
	// TODO
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
