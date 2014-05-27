package main

import (
	"dlib/dbus"
	"fmt"
)

// TODO rename
type sectionErrors map[string]string
type sessionErrors map[string]sectionErrors

type ConnectionSession struct {
	sessionPath dbus.ObjectPath
	devPath     dbus.ObjectPath
	data        connectionData

	ConnectionPath dbus.ObjectPath
	Uuid           string
	Type           string

	AvailablePages []string
	AvailableKeys  map[string][]string
	Errors         sessionErrors
	errorsSetKey   sessionErrors
}

//所有字段值都为string，后端自行转换为需要的值后提供给NM

func doNewConnectionSession(devPath dbus.ObjectPath, uuid string) (s *ConnectionSession) {
	s = &ConnectionSession{}
	s.sessionPath = dbus.ObjectPath(fmt.Sprintf("/com/deepin/daemon/ConnectionSession/%s", randString(8)))
	s.devPath = devPath
	s.Uuid = uuid
	s.data = make(connectionData)
	s.AvailablePages = make([]string, 0)
	s.AvailableKeys = make(map[string][]string)
	s.Errors = make(sessionErrors)
	s.errorsSetKey = make(sessionErrors)
	return s
}

func NewConnectionSessionByCreate(connectionType string, devPath dbus.ObjectPath) (s *ConnectionSession, err error) {
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
	s.updatePropAvailablePages()
	s.updatePropAvailableKeys()
	s.updatePropErrors()

	return
}

func NewConnectionSessionByOpen(uuid string, devPath dbus.ObjectPath) (s *ConnectionSession, err error) {
	connectionPath, err := nmGetConnectionByUuid(uuid)
	if err != nil {
		return
	}

	s = doNewConnectionSession(devPath, uuid)
	s.ConnectionPath = connectionPath

	// get connection data
	nmConn, err := nmNewSettingsConnection(connectionPath)
	if err != nil {
		return nil, err
	}
	s.data, err = nmConn.GetSettings()
	if err != nil {
		return nil, err
	}
	s.Type = getCustomConnectinoType(s.data)

	s.fixValues()

	// get secret data
	// TODO sectionVpnSecurity
	for _, secretFiled := range []string{sectionWirelessSecurity, section8021x, sectionGsm, sectionCdma} {
		if isSettingSectionExists(s.data, secretFiled) {
			wirelessSecrutiyData, err := nmConn.GetSecrets(sectionWirelessSecurity)
			if err == nil {
				for section, sectionData := range wirelessSecrutiyData {
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

	s.updatePropConnectionType()
	s.updatePropAvailablePages()
	s.updatePropAvailableKeys()
	s.updatePropErrors()

	// TODO
	logger.Debug("NewConnectionSessionByOpen():", s.data)

	return
}

func (s *ConnectionSession) fixValues() {
	// append missing sectionIpv6
	if !isSettingSectionExists(s.data, sectionIpv6) && isStringInArray(sectionIpv6, s.listSections()) {
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

// Save save current connection s.
func (s *ConnectionSession) Save() bool {
	if s.isErrorOccured() {
		logger.Debug("Errors occured when saving:", s.Errors)
		return false
	}

	if getSettingConnectionReadOnly(s.data) {
		logger.Debug("read only connection, don't allowed to save")
		return false
	}
	// TODO what about the connection has been deleted?

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

	dbus.UnInstallObject(s)
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

// Close cancel current connection s.
func (s *ConnectionSession) Close() {
	dbus.UnInstallObject(s)
}

// listSections return all pages related sections
func (s *ConnectionSession) listSections() (sections []string) {
	for _, page := range s.listPages() {
		sections = appendStrArrayUnion(sections, s.pageToSections(page)...)
	}
	return
}

// listPages return supported pages for target connection type.
func (s *ConnectionSession) listPages() (pages []string) {
	switch s.Type {
	case connectionWired:
		pages = []string{
			pageGeneral,
			pageEthernet,
			pageIPv4,
			pageIPv6,
			pageSecurity,
		}
	case connectionWireless:
		pages = []string{
			pageGeneral,
			pageWifi,
			pageIPv4,
			pageIPv6,
			pageSecurity,
		}
	case connectionWirelessAdhoc:
		pages = []string{
			pageGeneral,
			pageWifi,
			pageIPv4,
			pageIPv6,
			pageSecurity,
		}
	case connectionWirelessHotspot:
		pages = []string{
			pageGeneral,
			pageWifi,
			pageIPv4,
			pageIPv6,
			pageSecurity,
		}
	case connectionPppoe:
		pages = []string{
			pageGeneral,
			pageEthernet,
			pagePppoe,
			pagePpp,
			pageIPv4,
		}
	case connectionVpnL2tp:
		pages = []string{
			pageGeneral,
			pageVpnL2tp,
			pageVpnL2tpPpp,
			pageVpnL2tpIpsec,
			pageIPv4,
		}
	case connectionVpnOpenconnect:
		pages = []string{
			pageGeneral,
			pageVpnOpenconnect,
			pageIPv4,
			pageIPv6,
		}
	case connectionVpnOpenvpn:
		pages = []string{
			pageGeneral,
			pageVpnOpenvpn,
			pageVpnOpenvpnAdvanced,
			pageVpnOpenvpnSecurity,
			pageVpnOpenvpnProxies,
			pageIPv4,
			pageIPv6,
		}
		// when connection connection is static key, pageVpnOpenvpnTlsauth is not available
		if getSettingVpnOpenvpnKeyConnectionType(s.data) != NM_OPENVPN_CONTYPE_STATIC_KEY {
			pages = append(pages, pageVpnOpenvpnTlsauth)
		}
	case connectionVpnPptp:
		pages = []string{
			pageGeneral,
			pageVpnPptp,
			pageVpnPptpPpp,
			pageIPv4,
		}
	case connectionVpnVpnc:
		pages = []string{
			pageGeneral,
			pageVpnVpnc,
			pageVpnVpncAdvanced,
			pageIPv4,
		}
	case connectionMobileGsm:
		pages = []string{
			pageGeneral,
			pageMobile,
			pagePpp,
			pageIPv4,
		}
	case connectionMobileCdma:
		pages = []string{
			pageGeneral,
			pageMobileCdma,
			pagePpp,
			pageIPv4,
		}
	}
	return
}

func (s *ConnectionSession) pageToSections(page string) (sections []string) {
	switch page {
	default:
		logger.Error("pageToSections: invalid page name", page)
	case pageGeneral:
		sections = []string{sectionConnection}
	case pageMobile:
		sections = []string{sectionGsm}
	case pageMobileCdma:
		sections = []string{sectionCdma}
	case pageEthernet:
		sections = []string{sectionWired}
	case pageWifi:
		sections = []string{sectionWireless}
	case pageIPv4:
		sections = []string{sectionIpv4}
	case pageIPv6:
		sections = []string{sectionIpv6}
	case pageSecurity:
		switch s.Type {
		case connectionWired:
			sections = []string{section8021x}
		case connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
			if isSettingSectionExists(s.data, section8021x) {
				sections = []string{sectionWirelessSecurity, section8021x}
			} else {
				sections = []string{sectionWirelessSecurity}
			}
		}
	case pagePppoe:
		sections = []string{sectionPppoe}
	case pagePpp:
		sections = []string{sectionPpp}
	case pageVpnL2tp:
		sections = []string{sectionVpnL2tp}
	case pageVpnL2tpPpp:
		sections = []string{sectionVpnL2tpPpp}
	case pageVpnL2tpIpsec:
		sections = []string{sectionVpnL2tpIpsec}
	case pageVpnOpenconnect:
		sections = []string{sectionVpnOpenconnect}
	case pageVpnOpenvpn:
		sections = []string{sectionVpnOpenvpn}
	case pageVpnOpenvpnAdvanced:
		sections = []string{sectionVpnOpenvpnAdvanced}
	case pageVpnOpenvpnSecurity:
		sections = []string{sectionVpnOpenvpnSecurity}
	case pageVpnOpenvpnTlsauth:
		sections = []string{sectionVpnOpenvpnTlsauth}
	case pageVpnOpenvpnProxies:
		sections = []string{sectionVpnOpenvpnProxies}
	case pageVpnPptp:
		sections = []string{sectionVpnPptp}
	case pageVpnPptpPpp:
		sections = []string{sectionVpnPptpPpp}
	case pageVpnVpnc:
		sections = []string{sectionVpnVpnc}
	case pageVpnVpncAdvanced:
		sections = []string{sectionVpnVpncAdvanced}
	}
	return
}

func (s *ConnectionSession) getSectionOfPageKey(page, key string) string {
	sections := s.pageToSections(page)
	for _, section := range sections {
		if generalIsKeyInSettingSection(section, key) {
			return section
		}
	}
	logger.Errorf("get corresponding filed of key in page failed, page=%s, key=%s", page, key)
	return ""
}

// get valid keys of target page, show or hide some keys when special
// keys toggled
func (s *ConnectionSession) listKeys(page string) (keys []string) {
	sections := s.pageToSections(page)
	for _, section := range sections {
		// TODO
		// if isSettingSectionExists(s.data, section) {
		// }
		keys = appendStrArrayUnion(keys, generalGetSettingAvailableKeys(s.data, section)...)
	}
	if len(keys) == 0 {
		logger.Warning("there is no avaiable keys for page", page)
	}
	return
}

// GetAvailableValues return available values marshaled by json for target key.
func (s *ConnectionSession) GetAvailableValues(page, key string) (valuesJSON string) {
	var values []kvalue
	sections := s.pageToSections(page) // TODO
	for _, section := range sections {
		values = generalGetSettingAvailableValues(s.data, section, key)
		if len(values) > 0 {
			break
		}
	}
	valuesJSON, _ = marshalJSON(values)
	return
}

func (s *ConnectionSession) GetKey(page, key string) (value string) {
	section := s.getSectionOfPageKey(page, key)
	value = generalGetSettingKeyJSON(s.data, section, key)
	return
}

func (s *ConnectionSession) SetKey(page, key, value string) {
	section := s.getSectionOfPageKey(page, key)
	err := generalSetSettingKeyJSON(s.data, section, key, value)
	// logger.Debugf("SetKey(), %v, page=%s, filed=%s, key=%s, value=%s", err == nil, page, section, key, value) // TODO test
	s.updateErrorsWhenSettingKey(page, key, err)

	s.updatePropAvailablePages()
	s.updatePropAvailableKeys()
	s.updatePropErrors()

	return
}

func (s *ConnectionSession) updateErrorsWhenSettingKey(page, key string, err error) {
	if err == nil {
		// delete key error if exists
		sectionErrors, ok := s.errorsSetKey[page]
		if ok {
			_, ok := sectionErrors[key]
			if ok {
				delete(sectionErrors, key)
			}
		}
	} else {
		// append key error
		sectionErrorsData, ok := s.errorsSetKey[page]
		if !ok {
			sectionErrorsData = make(sectionErrors)
			s.errorsSetKey[page] = sectionErrorsData
		}
		sectionErrorsData[key] = err.Error()
	}
}

func (s *ConnectionSession) DebugListKeyDetail() (info string) {
	for _, page := range s.listPages() {
		pageData, ok := s.AvailableKeys[page]
		if !ok {
			logger.Warning("no available keys for page", page)
			continue
		}
		for _, key := range pageData {
			section := s.getSectionOfPageKey(page, key)
			t := generalGetSettingKeyType(section, key)
			// TODO convert to value json
			values := generalGetSettingAvailableValues(s.data, section, key)
			info += fmt.Sprintf("%s->%s[%s]: %s (%s)\n", page, key, getKtypeDescription(t), s.GetKey(page, key), values)
		}
	}
	return
}

func (s *ConnectionSession) DebugGetConnectionData() connectionData {
	return s.data
}

func (s *ConnectionSession) DebugGetErrors() sessionErrors {
	return s.Errors
}
