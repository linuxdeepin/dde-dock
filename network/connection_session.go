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
	"pkg.deepin.io/dde/daemon/network/nm"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/strv"
	"pkg.deepin.io/lib/utils"
	"sync"
	"time"
)

type sectionErrors map[string]string
type sessionErrors map[string]sectionErrors

// ConnectionSession used to provide DBus session to edit the
// connections. With the interfaces such as SetKey, GetKey,
// AvailableKeys and GetAvailableValues, the front-end could show the
// related widgets automatically.
type ConnectionSession struct {
	sessionPath      dbus.ObjectPath
	devPath          dbus.ObjectPath
	data             connectionData
	dataLocker       sync.RWMutex
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

	// used by font-end to update widget value that with proeprty
	// "alwaysUpdate", which should only update value when visible
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
		// set mac address
		macAddress, err := nmGeneralGetDeviceHwAddr(devPath)
		if err == nil && macAddress != "" {
			setSettingWiredMacAddress(s.data, convertMacAddressToArrayByte(macAddress))
		}
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
	case connectionVpnStrongswan:
		s.data = newVpnStrongswanConnectionData(id, s.Uuid)
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
	fillSectionCache(s.data)
	s.setProps()

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

	logger.Debugf("NewConnectionSessionByOpen(): %#v", s.data)
	return
}

func (s *ConnectionSession) fixValues() {
	// append missing nm.NM_SETTING_IP6_CONFIG_SETTING_NAME
	if !isSettingExists(s.data, nm.NM_SETTING_IP6_CONFIG_SETTING_NAME) && isStringInArray(nm.NM_SETTING_IP6_CONFIG_SETTING_NAME, getAvailableSections(s.data)) {
		initSettingSectionIpv6(s.data)
	}

	// fix ipv6 addresses and routes data structure, interface{}
	if isSettingIP6ConfigAddressesExists(s.data) {
		setSettingIP6ConfigAddresses(s.data, getSettingIP6ConfigAddresses(s.data))
	}
	if isSettingIP6ConfigRoutesExists(s.data) {
		setSettingIP6ConfigRoutes(s.data, getSettingIP6ConfigRoutes(s.data))
	}

	// remove address-data and gateway fields in IP4/IP6 section to keep
	// compatible with NetworkManager 1.0+
	if isSettingKeyExists(s.data, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, "address-data") {
		removeSettingKey(s.data, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, "address-data")
	}
	if isSettingKeyExists(s.data, nm.NM_SETTING_IP6_CONFIG_SETTING_NAME, "address-data") {
		removeSettingKey(s.data, nm.NM_SETTING_IP6_CONFIG_SETTING_NAME, "address-data")
	}
	if isSettingKeyExists(s.data, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, "gateway") {
		removeSettingKey(s.data, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, "gateway")
	}
	if isSettingKeyExists(s.data, nm.NM_SETTING_IP6_CONFIG_SETTING_NAME, "gateway") {
		removeSettingKey(s.data, nm.NM_SETTING_IP6_CONFIG_SETTING_NAME, "gateway")
	}

	// vpn plugin data and secret
	if getSettingConnectionType(s.data) == nm.NM_SETTING_VPN_SETTING_NAME {
		if !isSettingVpnDataExists(s.data) {
			setSettingVpnData(s.data, make(map[string]string))
		}
		if !isSettingVpnSecretsExists(s.data) {
			setSettingVpnSecrets(s.data, make(map[string]string))
		}
		switch getCustomConnectionType(s.data) {
		case connectionVpnStrongswan:
			// fix vpn strongswan password flags
			setSettingVpnStrongswanKeyPasswordFlags(s.data, nm.NM_SETTING_SECRET_FLAG_NONE)
		}
	}

	// do not use s.Type here for that it may be not initialized
	switch getCustomConnectionType(s.data) {
	case connectionPppoe:
		// append missing nm.NM_SETTING_WIRED_SETTING_NAME for pppoe
		if !isSettingExists(s.data, nm.NM_SETTING_WIRED_SETTING_NAME) {
			initSettingSectionWired(s.data)
		}
	case connectionMobileGsm, connectionMobileCdma:
		addSetting(s.data, nm.NM_SETTING_PPP_SETTING_NAME)
		logicSetSettingVkPppEnableLcpEcho(s.data, true)
	}

	// TODO fix secret flags
	// if isSettingVpnOpenvpnKeyCertpassFlagsExists(s.data) && getSettingVpnOpenvpnKeyCertpassFlags(s.data) == 1 {
	// setSettingVpnOpenvpnKeyCertpassFlags(s.data, nm.NM_OPENVPN_SECRET_FLAG_SAVE)
	// }
}

func (s *ConnectionSession) getSecrets() {
	if !s.getSecretsFromKeyring() {
		logger.Info("get secrets from keyring failed, try network-manager configuration again")
		s.getSecretsFromNM()
	}
}

func (s *ConnectionSession) getSecretsFromKeyring() (ok bool) {
	for _, section := range s.AvailableSections {
		realSetting := getAliasSettingRealName(section)
		if values, okNest := secretGetAll(s.Uuid, realSetting); okNest {
			ok = true
			secretsData := buildKeyringSecret(s.data, realSetting, values)
			s.doGetSecrets(secretsData)
		}
	}
	return
}
func (s *ConnectionSession) doGetSecrets(secretsData connectionData) {
	for section, sectionData := range secretsData {
		if !isSettingExists(s.data, section) {
			addSetting(s.data, section)
		}
		for key, value := range sectionData {
			s.data[section][key] = value
		}
	}
}

func (s *ConnectionSession) getSecretsFromNM() {
	switch getCustomConnectionType(s.data) {
	case connectionWired:
		if getSettingVk8021xEnable(s.data) {
			// TODO 8021x secret
			// s.doGetSecretsFromNM(nm.NM_SETTING_802_1X_SETTING_NAME)
		}
	case connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
		if getSettingVk8021xEnable(s.data) {
			// TODO 8021x secret
			// s.doGetSecretsFromNM(nm.NM_SETTING_802_1X_SETTING_NAME)
		} else {
			s.doGetSecretsFromNM(nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME)
		}
	case connectionPppoe:
		s.doGetSecretsFromNM(nm.NM_SETTING_PPPOE_SETTING_NAME)
	case connectionMobileGsm:
		// FIXME: if the connection owns no secret key, such as "US -> AT&T -> MEdia Net (phones)"
		// it will popup password dialog when editing the connection.
		// s.doGetSecretsFromNM(nm.NM_SETTING_GSM_SETTING_NAME)
	case connectionMobileCdma:
		// FIXME: same with connectionMobileGsm
		// s.doGetSecretsFromNM(nm.NM_SETTING_CDMA_SETTING_NAME)
	case connectionVpnL2tp, connectionVpnOpenconnect, connectionVpnPptp, connectionVpnVpnc, connectionVpnOpenvpn:
		// ignore vpn secrets
	}
}
func (s *ConnectionSession) doGetSecretsFromNM(secretSection string) {
	if isSettingExists(s.data, secretSection) {
		if secretsData, err := nmGetConnectionSecrets(s.ConnectionPath, secretSection); err == nil {
			if isSettingExists(s.data, secretSection) {
				s.doGetSecrets(secretsData)
			}
		}
	}
}

func (s *ConnectionSession) updateSecretsToKeyring() {
	for sectionName, sectionData := range s.data {
		for keyName, variant := range sectionData {
			if isSecretKey(s.data, sectionName, keyName) {
				if sectionName == nm.NM_SETTING_VPN_SETTING_NAME && keyName == nm.NM_SETTING_VPN_SECRETS {
					// dispatch vpn secret keys specially
					vpnSecrets := getSettingVpnSecrets(s.data)
					for k, v := range vpnSecrets {
						secretSet(s.Uuid, sectionName, k, v)
					}
				} else if value, ok := variant.Value().(string); ok {
					secretSet(s.Uuid, sectionName, keyName, value)
				}
			}
		}
	}
}

// Save save current connection session, and the 'activated' will special whether activate it,
// but if the connection non-exists, the 'activated' no effect
func (s *ConnectionSession) Save(activated bool) (ok bool, err error) {
	logger.Debugf("Save connection: %#v", s.data)
	s.dataLocker.Lock()
	defer s.dataLocker.Unlock()

	if s.isErrorOccurred() {
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

		correctConnectionData(s.data)
		err = nmConn.Update(s.data)
		if err != nil {
			logger.Error(err)
			return false, err
		}
		if !activated {
			err = nmConn.Save()
		} else {
			_, err = manager.ActivateConnection(s.Uuid, s.devPath)
		}
		if err != nil {
			logger.Error("Failed to save exists connection:", err)
		}
	} else {
		// create new connection and activate it
		connectionType := getCustomConnectionType(s.data)

		// keep ID same with SSID for wireless connections
		switch connectionType {
		case connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
			setSettingConnectionId(s.data, string(getSettingWirelessSsid(s.data)))
		}

		switch connectionType {
		case connectionWired, connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
			_, _, err = nmAddAndActivateConnection(s.data, s.devPath, true)
		default:
			_, err = nmAddConnection(s.data)
		}
		if err != nil {
			logger.Error("Failed to save non-exists connection:", err)
		}
	}

	s.updateSecretsToKeyring()

	manager.removeConnectionSession(s)
	return true, nil
}

func (s *ConnectionSession) isErrorOccurred() bool {
	for _, v := range s.Errors {
		if len(v) > 0 {
			return true
		}
	}
	return false
}

// Close cancel current connection.
func (s *ConnectionSession) Close() {
	s.dataLocker.Lock()
	defer s.dataLocker.Unlock()
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
	s.dataLocker.Lock()
	defer s.dataLocker.Unlock()
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

// GetAvailableValues return available values marshaled by JSON for target key.
func (s *ConnectionSession) GetAvailableValues(section, key string) (valuesJSON string) {
	s.dataLocker.RLock()
	defer s.dataLocker.RUnlock()
	var values []kvalue
	values = generalGetSettingAvailableValues(s.data, section, key)
	valuesJSON, _ = marshalJSON(values)
	return
}

// GetKey get target key value which marshaled by JSON.
func (s *ConnectionSession) GetKey(section, key string) (valueJSON string) {
	s.dataLocker.RLock()
	defer s.dataLocker.RUnlock()
	valueJSON = generalGetSettingKeyJSON(s.data, section, key)
	return
}

// GetKeyName return the display name for special key.
func (s *ConnectionSession) GetKeyName(section, key string) (name string, err error) {
	s.dataLocker.RLock()
	defer s.dataLocker.RUnlock()
	name, err = getRelatedKeyName(s.data, section, key)
	return
}

// SetKey set target key with new value, the value should be marshaled by JSON.
func (s *ConnectionSession) SetKey(section, key, valueJSON string) {
	logger.Debugf("SetKey(), section=%s, key=%s, valueJSON=%s", section, key, valueJSON)
	s.dataLocker.Lock()
	defer s.dataLocker.Unlock()
	err := generalSetSettingKeyJSON(s.data, section, key, valueJSON)
	if err == nil {
		// set addresses mask when method in manual
		if section == "ipv4" && key == "method" && valueJSON == "\"manual\"" {
			addresses := interfaceToArrayArrayUint32(generalGetSettingDefaultValue(section, "addresses"))
			logicSetSettingVkIp4ConfigAddressesMask(s.data, convertIpv4PrefixToNetMask(addresses[0][1]))
			logicSetSettingVkIp4ConfigAddressesAddress(s.data, "")
		}
	}
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

// IsDefaultExpandedSection check if target virtual section should be
// expanded default.
func (s *ConnectionSession) IsDefaultExpandedSection(vsection string) (bool, error) {
	switch s.Type {
	case connectionWired:
		switch vsection {
		case nm.NM_SETTING_VS_IPV4:
			return true, nil
		}
	case connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
		switch vsection {
		case nm.NM_SETTING_VS_SECURITY:
			return true, nil
		}
	case connectionPppoe:
		switch vsection {
		case nm.NM_SETTING_VS_PPPOE:
			return true, nil
		}
	case connectionMobileGsm, connectionMobileCdma:
		switch vsection {
		case nm.NM_SETTING_VS_MOBILE:
			return true, nil
		}
	case connectionVpnL2tp, connectionVpnOpenconnect, connectionVpnPptp, connectionVpnVpnc, connectionVpnOpenvpn:
		switch vsection {
		case nm.NM_SETTING_VS_VPN:
			return true, nil
		}
	default:
		switch vsection {
		case nm.NM_SETTING_VS_IPV4:
			return true, nil
		}
	}
	return false, nil
}

// Debug functions

// DebugGetConnectionData get current connection data.
func (s *ConnectionSession) DebugGetConnectionData() connectionData {
	return s.data
}

// DebugGetErrors get current errors.
func (s *ConnectionSession) DebugGetErrors() sessionErrors {
	return s.Errors
}

// DebugListKeyDetail get all key deails, including all the available
// key values.
func (s *ConnectionSession) DebugListKeyDetail() (info string) {
	for _, vsection := range s.AvailableVirtualSections {
		for _, section := range getAvailableSectionsOfVsection(s.data, vsection) {
			sectionKeys, ok := s.AvailableKeys[section]
			if !ok {
				logger.Warning("no available keys for section", section)
				continue
			}
			for _, key := range sectionKeys {
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

// ListAvailableKeyDetail get all available key details, include vitual sections
func (s *ConnectionSession) ListAvailableKeyDetail() string {
	var (
		firstKey      bool   = true
		firstVSection bool   = true
		info          string = "["

		allVSections = getAllVsections(s.data)
	)
	for _, vsection := range allVSections {
		if !strv.Strv(s.AvailableVirtualSections).Contains(vsection) {
			continue
		}

		sectionInfo, ok := virtualSections[vsection]
		if !ok {
			continue
		}

		if !firstVSection {
			info += ", "
		}
		firstKey = true
		info += fmt.Sprintf("{\"VirtualSection\": \"%s\", \"Keys\":[", vsection)
		sections := getAvailableSectionsOfVsection(s.data, vsection)
		for _, keyInfo := range sectionInfo.Keys {
			if !strv.Strv(sections).Contains(keyInfo.Section) {
				continue
			}

			sectionKeys, ok := s.AvailableKeys[keyInfo.Section]
			if !ok {
				continue
			}

			if len(sectionKeys) == 0 || !strv.Strv(sectionKeys).Contains(keyInfo.Key) {
				continue
			}

			if !firstKey {
				info += ","
			}

			info += fmt.Sprintf("{\"Section\": \"%s\", \"Key\": \"%s\",", keyInfo.Section, keyInfo.Key)
			info += fmt.Sprintf("\"Value\": %s", s.GetKey(keyInfo.Section, keyInfo.Key))
			values := generalGetSettingAvailableValues(s.data, keyInfo.Section, keyInfo.Key)
			if len(values) > 0 {
				valuesJSON, _ := marshalJSON(values)
				info += fmt.Sprintf(",\"Values\": %s", valuesJSON)
			}

			info += "}"
			firstKey = false
		}
		info += "]}"
		firstVSection = false
	}

	info += "]"
	return info
}

func correctConnectionData(data connectionData) {
	correctIPv6DataType(data)
	correctIgnoreAutoDNS(data)
}

func correctIPv6DataType(data connectionData) {
	for section, value := range data {
		if section != "ipv6" {
			continue
		}

		tmp, ok := value["addresses"]
		if ok && !isInterfaceEmpty(tmp.Value()) {
			addrs := interfaceToIpv6Addresses(tmp.Value())
			value["addresses"] = dbus.MakeVariant(addrs)
		}

		tmp, ok = value["routes"]
		if ok && !isInterfaceEmpty(tmp.Value()) {
			routes := interfaceToIpv6Routes(tmp.Value())
			value["routes"] = dbus.MakeVariant(routes)
		}
	}
}

func correctIgnoreAutoDNS(data connectionData) {
	for section, value := range data {
		if section != "ipv4" && section != "ipv6" {
			continue
		}

		// if dns specialed, and method is 'auto', set ignore-auto-dns to true
		method, ok := value["method"]
		logger.Debug("[correctIgnoreAutoDNS] method:", method.String())
		if !ok || method.Value().(string) != "auto" {
			continue
		}

		if section == "ipv4" {
			dns := getSettingIP4ConfigDns(data)
			logger.Debug("[correctIgnoreAutoDNS] ipv4 dns:", dns)
			if len(dns) == 0 {
				removeSettingIP4ConfigIgnoreAutoDns(data)
				continue
			}
			setSettingIP4ConfigIgnoreAutoDns(data, true)
		} else if section == "ipv6" {
			dns := getSettingIP6ConfigDns(data)
			logger.Debug("[correctIgnoreAutoDNS] ipv6 dns:", dns)
			if len(dns) == 0 {
				removeSettingIP6ConfigIgnoreAutoDns(data)
				continue
			}
			setSettingIP6ConfigIgnoreAutoDns(data, true)
		}
	}
}
