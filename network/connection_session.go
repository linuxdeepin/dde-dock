/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

import (
	"fmt"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/utils"
	"time"
)

type sectionErrors map[string]string
type sessionErrors map[string]sectionErrors

type ConnectionSession struct {
	sessionPath      dbus.ObjectPath
	devPath          dbus.ObjectPath
	data             connectionData
	connectionExists bool

	ConnectionPath dbus.ObjectPath
	Uuid           string
	Type           string // customized connection types, e.g. connectionMobileGsm

	AllowDelete           bool
	AllowEditConnectionId bool

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
	s.connectionExists = false

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

	// disable vpn autoconnect default
	if isVpnConnection(s.data) {
		manager.config.addVpnConfig(s.Uuid)
		logicSetSettingVkVpnAutoconnect(s.data, false)
	}

	fillSectionCache(s.data)
	s.setProps()
	logger.Debugf("newConnectionSessionByCreate(): %#v", s.data)
	return
}

func newConnectionSessionByOpen(uuid string, devPath dbus.ObjectPath) (s *ConnectionSession, err error) {
	connectionPath, err := nmGetConnectionByUuid(uuid)
	if err != nil {
		return
	}

	s = doNewConnectionSession(devPath, uuid)
	s.connectionExists = true
	s.ConnectionPath = connectionPath

	// get connection data
	s.data, err = nmGetConnectionData(s.ConnectionPath)
	if err != nil {
		return nil, err
	}

	s.fixValues()

	// execute asynchronous to avoid front-end block if
	// NeedSecrets() signal emitted
	chSecret := make(chan int)
	go func() {
		s.getSecrets()
		chSecret <- 0
	}()
	select {
	case <-time.After(500 * time.Millisecond):
	case <-chSecret:
	}

	fillSectionCache(s.data)
	s.setProps()
	logger.Debugf("NewConnectionSessionByOpen(): %#v", s.data)
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

	// do not use s.Type here for that it may be not initialized
	switch getCustomConnectionType(s.data) {
	case connectionPppoe:
		// append missing sectionWired for pppoe
		if !isSettingSectionExists(s.data, sectionWired) {
			initSettingSectionWired(s.data)
		}
	case connectionMobileGsm, connectionMobileCdma:
		addSettingSection(s.data, sectionPpp)
		logicSetSettingVkPppEnableLcpEcho(s.data, true)
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
		// FIXME: if the connection owns no secret key, such as "US -> AT&T -> MEdia Net (phones)"
		// it will popup password dialog when editing the connection.
		// s.doGetSecrets(sectionGsm)
	case connectionMobileCdma:
		// FIXME: same with connectionMobileGsm
		// s.doGetSecrets(sectionCdma)
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

	logger.Debugf("Save connection: %#v", s.data)

	if s.isErrorOccured() {
		logger.Info("Save Errors:", s.Errors)
		logger.Info("Save settingKeyErrors:", s.settingKeyErrors)
		return false, nil
	}

	if getSettingConnectionReadOnly(s.data) {
		err = fmt.Errorf("read only connection, don't allowed to save")
		logger.Debug(err)
		return false, err
	}

	refileSectionCache(s.data)

	if s.connectionExists {
		// update connection data and activate it
		nmConn, err := nmNewSettingsConnection(s.ConnectionPath)
		if err != nil {
			return false, err
		}
		defer nmDestroySettingsConnection(nmConn)

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

		// keep ID same with SSID for wireless connections
		switch connectionType {
		case connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
			setSettingConnectionId(s.data, string(getSettingWirelessSsid(s.data)))
		}

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
	// clean up vpn config if abort the new connection
	if !s.connectionExists {
		if isVpnConnection(s.data) {
			manager.config.removeVpnConfig(s.Uuid)
		}
	}

	manager.removeConnectionSession(s)
}

// GetAllKeys return all may used section key information in current session.
func (s *ConnectionSession) GetAllKeys() (infoJSON string) {
	allVsectionInfo := make([]VsectionInfo, 0)
	vsections := getAllVsections(s.data)
	for _, vsection := range vsections {
		if sectionInfo, ok := virtualSections[vsection]; ok {
			sectionInfo.fixExpanded(s.data)
			for _, keyInfo := range sectionInfo.Keys {
				keyInfo.fixReadonly(s.data)
			}
			allVsectionInfo = append(allVsectionInfo, sectionInfo)
		} else {
			logger.Errorf("get virtaul section info failed: %s", allVsectionInfo)
		}
	}
	infoJSON, _ = marshalJSON(allVsectionInfo)
	return
}

// GetAvailableValues return available values marshaled by json for target key.
func (s *ConnectionSession) GetAvailableValues(section, key string) (valuesJSON string) {
	var values []kvalue
	values = generalGetSettingAvailableValues(s.data, section, key)
	valuesJSON, _ = marshalJSON(values)
	return
}

func (s *ConnectionSession) GetKey(section, key string) (valueJSON string) {
	valueJSON = generalGetSettingKeyJSON(s.data, section, key)
	return
}

func (s *ConnectionSession) SetKey(section, key, valueJSON string) {
	logger.Debugf("SetKey(), section=%s, key=%s, valueJSON=%s", section, key, valueJSON)
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
	for _, vsection := range s.AvailableVirtualSections {
		for _, section := range getAvailableSectionsOfVsection(s.data, vsection) {
			sectionKeys, ok := s.AvailableKeys[section]
			if !ok {
				logger.Warning("no available keys for section", section)
				continue
			}
			for _, key := range sectionKeys {
				// TODO: remove?
				// section := getSectionOfKeyInVsection(s.data, vsection, key)
				t := generalGetSettingKeyType(section, key)
				if values := generalGetSettingAvailableValues(s.data, section, key); len(values) > 0 {
					valuesJSON, _ := marshalJSON(values)
					info += fmt.Sprintf("%s: %s[%s](%s): %s (%s)\n", vsection, section, key,
						getKtypeDesc(t), s.GetKey(section, key), valuesJSON)
				} else {
					info += fmt.Sprintf("%s: %s[%s](%s): %s\n", vsection, section, key,
						getKtypeDesc(t), s.GetKey(section, key))
				}
			}
		}
	}
	return
}
