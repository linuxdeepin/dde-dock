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
	"encoding/binary"
	"fmt"
	"os"
	"sync"
	"time"

	"pkg.deepin.io/dde/daemon/network/nm"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/strv"
	"pkg.deepin.io/lib/utils"
)

type sectionErrors map[string]string
type sessionErrors map[string]sectionErrors

// ConnectionSession used to provide DBus session to edit the
// connections. With the interfaces such as SetKey, GetKey,
// AvailableKeys and GetAvailableValues, the front-end could show the
// related widgets automatically.
type ConnectionSession struct {
	service          *dbusutil.Service
	sessionPath      dbus.ObjectPath
	devPath          dbus.ObjectPath
	data             connectionData
	dataLocker       sync.RWMutex
	connectionExists bool

	PropsMu        sync.RWMutex
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
	signals *struct {
		ConnectionDataChanged struct{}
	}

	methods *struct {
		DebugGetErrors           func() `out:"errors"`
		DebugListKeyDetail       func() `out:"info"`
		GetAllKeys               func() `out:"keys"`
		GetAvailableValues       func() `in:"section,key" out:"valuesJSON"`
		GetKey                   func() `in:"section,key" out:"valueJSON"`
		GetKeyName               func() `in:"section,key" out:"name"`
		IsDefaultExpandedSection func() `in:"vsection" out:"result"`
		ListAvailableKeyDetail   func() `out:"detail"`
		Save                     func() `in:"activated" out:"ok"`
		SetKey                   func() `in:"section,key,valueJSON"`
		SetKeyFd                 func() `in:"section,key" out:"fd"`
	}
}

func doNewConnectionSession(devPath dbus.ObjectPath, uuid string,
	service *dbusutil.Service) (s *ConnectionSession) {
	s = &ConnectionSession{}
	s.service = service
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

func newConnectionSessionByCreate(connectionType string, devPath dbus.ObjectPath, service *dbusutil.Service) (s *ConnectionSession, err error) {
	if !isStringInArray(connectionType, supportedConnectionTypes) {
		err = fmt.Errorf("connection type is out of support: %s", connectionType)
		logger.Error(err)
		return
	}

	s = doNewConnectionSession(devPath, utils.GenUuid(), service)
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
		macAddress, err := nmGeneralGetDeviceHwAddr(devPath, true)
		if err == nil && macAddress != "" {
			setSettingWiredMacAddress(s.data, convertMacAddressToArrayByte(macAddress))
		}
	case connectionWireless:
		s.data = newWirelessConnectionData(id, s.Uuid, nil, apSecNone)
	case connectionWirelessAdhoc:
		s.data = newWirelessAdhocConnectionData(id, s.Uuid)
	case connectionWirelessHotspot:
		s.data = newWirelessHotspotConnectionData(id, s.Uuid)
		if devPath != "" && devPath != "/" {
			setSettingConnectionInterfaceName(s.data, nmGetDeviceInterface(devPath))
		}
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
		logicSetSettingVkVpnAutoconnect(s.data, false)
	}

	fillSectionCache(s.data)
	s.setProps()
	s.setPropAllowDelete(false)
	logger.Debugf("newConnectionSessionByCreate(): %#v", s.data)
	return
}

func newConnectionSessionByOpen(uuid string, devPath dbus.ObjectPath,
	service *dbusutil.Service) (s *ConnectionSession, err error) {

	logger.Debugf("newConnectionSessionByOpen uuid: %q, devPath: %q",
		uuid, devPath)
	connectionPath, err := nmGetConnectionByUuid(uuid)
	if err != nil {
		return
	}

	s = doNewConnectionSession(devPath, uuid, service)
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
	s.setPropAllowDelete(true)
	s.getSecrets()

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
	s.getSecretsFromNM()
}

//func (s *ConnectionSession) getSecretsFromKeyring() (ok bool) {
//	for _, section := range s.AvailableSections {
//		realSetting := getAliasSettingRealName(section)
//		if values, okNest := secretGetAll(s.Uuid, realSetting); okNest {
//			ok = true
//			secretsData := buildKeyringSecret(s.data, realSetting, values)
//			s.doGetSecrets(secretsData)
//		}
//	}
//	return
//}
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
			s.doGetSecretsFromNM(nm.NM_SETTING_802_1X_SETTING_NAME)
		}
	case connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
		if getSettingVk8021xEnable(s.data) {
			s.doGetSecretsFromNM(nm.NM_SETTING_802_1X_SETTING_NAME)
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
	case connectionVpnL2tp, connectionVpnOpenconnect, connectionVpnPptp, connectionVpnVpnc, connectionVpnOpenvpn, connectionVpnStrongswan:
		s.doGetSecretsFromNM(nm.NM_SETTING_VPN_SETTING_NAME)
	}
}

//func (s *ConnectionSession) getSecretsFromNM() {
//	switch getCustomConnectionType(s.data) {
//	case connectionWired:
//		if getSettingVk8021xEnable(s.data) {
//			// TODO 8021x secret
//			// s.doGetSecretsFromNM(nm.NM_SETTING_802_1X_SETTING_NAME)
//		}
//	case connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
//		if getSettingVk8021xEnable(s.data) {
//			// TODO 8021x secret
//			// s.doGetSecretsFromNM(nm.NM_SETTING_802_1X_SETTING_NAME)
//		} else {
//			s.doGetSecretsFromNM(nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME)
//		}
//	case connectionPppoe:
//		s.doGetSecretsFromNM(nm.NM_SETTING_PPPOE_SETTING_NAME)
//	case connectionMobileGsm:
//		// FIXME: if the connection owns no secret key, such as "US -> AT&T -> MEdia Net (phones)"
//		// it will popup password dialog when editing the connection.
//		// s.doGetSecretsFromNM(nm.NM_SETTING_GSM_SETTING_NAME)
//	case connectionMobileCdma:
//		// FIXME: same with connectionMobileGsm
//		// s.doGetSecretsFromNM(nm.NM_SETTING_CDMA_SETTING_NAME)
//	case connectionVpnL2tp, connectionVpnOpenconnect, connectionVpnPptp, connectionVpnVpnc, connectionVpnOpenvpn:
//		// ignore vpn secrets
//	}
//}
func (s *ConnectionSession) doGetSecretsFromNM(secretSection string) {
	if isSettingExists(s.data, secretSection) {
		if secretsData, err := nmGetConnectionSecrets(s.ConnectionPath, secretSection); err == nil {
			if isSettingExists(s.data, secretSection) {
				s.doGetSecrets(secretsData)
			}
		}
	}
}

//func (s *ConnectionSession) updateSecretsToKeyring() {
//	for sectionName, sectionData := range s.data {
//		for keyName, variant := range sectionData {
//			if isSecretKey(s.data, sectionName, keyName) {
//				if sectionName == nm.NM_SETTING_VPN_SETTING_NAME && keyName == nm.NM_SETTING_VPN_SECRETS {
//					// dispatch vpn secret keys specially
//					vpnSecrets := getSettingVpnSecrets(s.data)
//					for k, v := range vpnSecrets {
//						secretSet(s.Uuid, sectionName, k, v)
//					}
//				} else if value, ok := variant.Value().(string); ok {
//					secretSet(s.Uuid, sectionName, keyName, value)
//				}
//			}
//		}
//	}
//}

// Save save current connection session, and the 'activated' will special whether activate it,
// but if the connection non-exists, the 'activated' no effect
func (s *ConnectionSession) Save(activated bool) (bool, *dbus.Error) {
	ok, err := s.save(activated)
	return ok, dbusutil.ToError(err)
}

func (s *ConnectionSession) save(activated bool) (ok bool, err error) {
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

		correctConnectionData(s.data)
		err = nmConn.Update(0, s.data)
		if err != nil {
			logger.Error(err)
			return false, err
		}
		manager.secretAgent.saveSecrets(s.data, s.ConnectionPath)

		if activated {
			_, err := manager.activateConnection(s.Uuid, s.devPath)
			if err != nil {
				logger.Warning(err)
			}
		}

	} else {
		// create new connection and activate it
		connectionType := getCustomConnectionType(s.data)

		// keep ID same with SSID for wireless connections
		switch connectionType {
		case connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
			setSettingConnectionId(s.data, decodeSsid(getSettingWirelessSsid(s.data)))
		}

		// wired, wireless, wireless-adhoc, wireless-hotspot should auto activated if available
		if activated {
			_, _, err = nmAddAndActivateConnection(s.data, s.devPath, true)
		} else {
			_, err = nmAddConnection(s.data)
		}
		if err != nil {
			logger.Error("Failed to save non-exists connection:", err)
		}
		s.connectionExists = true
	}

	//s.updateSecretsToKeyring()

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
func (s *ConnectionSession) Close() *dbus.Error {
	s.dataLocker.Lock()
	defer s.dataLocker.Unlock()
	manager.removeConnectionSession(s)
	return nil
}

// GetAllKeys return all may used section key information in current session.
func (s *ConnectionSession) GetAllKeys() (infoJSON string, err *dbus.Error) {
	s.dataLocker.Lock()
	defer s.dataLocker.Unlock()
	allVsectionInfo := make([]VsectionInfo, 0)
	vsections := getAllVsections(s.data)
	for _, vsection := range vsections {
		if sectionInfo, ok := virtualSections[vsection]; ok {
			sectionInfo.fixExpanded(s.data)
			// for _, keyInfo := range sectionInfo.Keys {
			// keyInfo.fixReadonly(s.data)
			// }
			allVsectionInfo = append(allVsectionInfo, sectionInfo)
		} else {
			logger.Errorf("get virtaul section info failed: %s", allVsectionInfo)
		}
	}
	infoJSON, _ = marshalJSON(allVsectionInfo)
	return
}

// GetAvailableValues return available values marshaled by JSON for target key.
func (s *ConnectionSession) GetAvailableValues(section, key string) (valuesJSON string,
	err *dbus.Error) {
	s.dataLocker.RLock()
	defer s.dataLocker.RUnlock()
	var values []kvalue
	values = generalGetSettingAvailableValues(s.data, section, key)
	valuesJSON, _ = marshalJSON(values)
	return
}

func (s *ConnectionSession) getKey(section, key string) (valueJSON string) {
	s.dataLocker.RLock()
	defer s.dataLocker.RUnlock()
	valueJSON = generalGetSettingKeyJSON(s.data, section, key)
	return
}

// GetKey get target key value which marshaled by JSON.
func (s *ConnectionSession) GetKey(section, key string) (valueJSON string, err *dbus.Error) {
	valueJSON = s.getKey(section, key)
	return
}

// GetKeyName return the display name for special key.
func (s *ConnectionSession) GetKeyName(section, key string) (name string, busErr *dbus.Error) {
	s.dataLocker.RLock()
	defer s.dataLocker.RUnlock()
	name, err := getRelatedKeyName(s.data, section, key)
	busErr = dbusutil.ToError(err)
	return
}

// SetKey set target key with new value, the value should be marshaled by JSON.
func (s *ConnectionSession) SetKey(section, key, valueJSON string) *dbus.Error {
	s.setKey(section, key, valueJSON)
	return nil
}

func (s *ConnectionSession) setKey(section, key, valueJSON string) {
	logger.Debugf("SetKey section=%q, key=%q, len(valueJSON)=%d", section, key, len(valueJSON))
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
		// fix cloned address deleted not work
		if key == "cloned-mac-address" && valueJSON == "\"\"" &&
			(section == "802-3-ethernet" || section == "802-11-wireless") {
			generalSetSettingKeyJSON(s.data, section, "assigned-mac-address", valueJSON)
		}
	}
	s.updateErrorsWhenSettingKey(section, key, err)
	s.setProps()
	return
}

func (s *ConnectionSession) SetKeyFd(section, key string) (dbus.UnixFD, *dbus.Error) {
	fd, err := s.setKeyFd(section, key)
	return fd, dbusutil.ToError(err)
}

func (s *ConnectionSession) setKeyFd(section, key string) (dbus.UnixFD, error) {
	const deadline = 10 * time.Second
	r, w, err := os.Pipe()
	if err != nil {
		return 0, err
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		var length uint32
		err := binary.Read(r, binary.LittleEndian, &length)
		if err != nil {
			logger.Warning("SetKeyFd failed to read length:", err)
			return
		}

		logger.Debug("length:", length)
		if length > 1024 {
			logger.Warning("SetKeyFd length > 1024")
			return
		}
		buf := make([]byte, length)
		n, err := r.Read(buf)
		if err != nil {
			logger.Warning("SetKeyFd failed to read value:", err)
			return
		}
		secretValue := string(buf[:n])
		ch <- secretValue
	}()

	go func() {
		var secretValue string
		var ok bool

		select {
		case secretValue, ok = <-ch:
		case <-time.After(deadline):
			logger.Warning("SetKeyFd timeout")
		}
		r.Close()
		w.Close()
		if ok {
			s.setKey(section, key, secretValue)
		}
	}()

	return dbus.UnixFD(w.Fd()), nil
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
func (s *ConnectionSession) IsDefaultExpandedSection(vsection string) (bool, *dbus.Error) {
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

// DebugGetErrors get current errors.
func (s *ConnectionSession) DebugGetErrors() (sessionErrors, *dbus.Error) {
	return s.Errors, nil
}

// DebugListKeyDetail get all key deails, including all the available
// key values.
func (s *ConnectionSession) DebugListKeyDetail() (info string, err *dbus.Error) {
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
						getKtypeDesc(t), s.getKey(section, key), valuesJSON)
				} else {
					info += fmt.Sprintf("%s: %s[%s](%s): %s\n", vsection, section, key,
						getKtypeDesc(t), s.getKey(section, key))
				}
			}
		}
	}
	return
}

// ListAvailableKeyDetail get all available key details, include vitual sections
func (s *ConnectionSession) ListAvailableKeyDetail() (string, *dbus.Error) {
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
			info += fmt.Sprintf("\"Value\": %s", s.getKey(keyInfo.Section, keyInfo.Key))
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
	return info, nil
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
		if ok {
			if isInterfaceEmpty(tmp.Value()) {
				delete(value, "addresses")
			} else {
				addrs := interfaceToIpv6Addresses(tmp.Value())
				value["addresses"] = dbus.MakeVariant(addrs)
			}
		}

		tmp, ok = value["routes"]
		if ok {
			if isInterfaceEmpty(tmp.Value()) {
				delete(value, "routes")
			} else {
				routes := interfaceToIpv6Routes(tmp.Value())
				value["routes"] = dbus.MakeVariant(routes)
			}
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
