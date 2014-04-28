package main

import (
	"dlib/dbus"
	"fmt"
)

type FieldKeyErrors map[string]string
type PageKeyErrors map[string]FieldKeyErrors

type ConnectionSession struct {
	sessionPath    dbus.ObjectPath
	connPath       dbus.ObjectPath
	devPath        dbus.ObjectPath
	data           connectionData
	connectionType string

	CurrentUUID string

	AllowSave      bool // TODO really need?
	AvailablePages []string
	AvailableKeys  map[string][]string
	Errors         PageKeyErrors
}

//所有字段值都为string，后端自行转换为需要的值后提供给NM

func doNewConnectionSession(devPath dbus.ObjectPath, uuid string) (s *ConnectionSession) {
	s = &ConnectionSession{}
	s.sessionPath = dbus.ObjectPath(fmt.Sprintf("/com/deepin/daemon/ConnectionSession/%s", randString(8)))
	s.devPath = devPath
	s.CurrentUUID = uuid
	s.data = make(connectionData)
	s.AllowSave = false // TODO
	s.AvailablePages = make([]string, 0)
	s.AvailableKeys = make(map[string][]string)
	s.Errors = make(PageKeyErrors)
	return s
}

func NewConnectionSessionByCreate(connectionType string, devPath dbus.ObjectPath) (s *ConnectionSession, err error) {
	if !isStringInArray(connectionType, supportedConnectionTypes) {
		err = fmt.Errorf("connection type is out of support: %s", connectionType)
		logger.Error(err)
		return
	}

	s = doNewConnectionSession(devPath, newUUID())

	// TODO
	// new connection data, id is left here
	s.connectionType = connectionType
	switch s.connectionType {
	case typeWired:
		s.data = newWiredConnectionData("", s.CurrentUUID)
	case typeWireless:
		s.data = newWirelessConnectionData("", s.CurrentUUID, nil, apSecNone)
	case typePppoe:
		s.data = newPppoeConnectionData("", s.CurrentUUID)
	case typeVpnL2tp:
		s.data = newVpnL2tpConnectionData("", s.CurrentUUID)
	case typeVpnOpenconnect:
		s.data = newVpnOpenconnectConnectionData("", s.CurrentUUID)
	case typeVpnPptp:
		s.data = newVpnPptpConnectionData("", s.CurrentUUID)
	case typeVpnVpnc:
		s.data = newVpnVpncConnectionData("", s.CurrentUUID)
	case typeVpnOpenvpn:
		s.data = newVpnOpenvpnConnectionData("", s.CurrentUUID)
	}

	// s.updatePropAllowSave(false) // TODO
	s.updatePropAvailablePages()
	s.updatePropAvailableKeys()
	s.updatePropErrors()

	return
}

func NewConnectionSessionByOpen(uuid string, devPath dbus.ObjectPath) (s *ConnectionSession, err error) {
	connPath, err := nmGetConnectionByUuid(uuid)
	if err != nil {
		return
	}

	s = doNewConnectionSession(devPath, uuid)
	s.connPath = connPath

	// get connection data
	nmConn, err := nmNewSettingsConnection(connPath)
	if err != nil {
		return nil, err
	}
	s.data, err = nmConn.GetSettings()
	if err != nil {
		return nil, err
	}
	s.connectionType = generalGetConnectionType(s.data)

	// get secret data
	for _, secretFiled := range []string{fieldWirelessSecurity, field8021x} {
		if isSettingFieldExists(s.data, secretFiled) {
			wirelessSecrutiyData, err := nmConn.GetSecrets(fieldWirelessSecurity)
			if err == nil {
				for field, fieldData := range wirelessSecrutiyData {
					if !isSettingFieldExists(s.data, field) {
						addSettingField(s.data, field)
					}
					for key, value := range fieldData {
						s.data[field][key] = value
					}
				}
			}
		}
	}

	// s.updatePropAllowSave(false) // TODO
	s.updatePropAvailablePages()
	s.updatePropAvailableKeys()
	s.updatePropErrors()

	// TODO
	logger.Debug("NewConnectionSessionByOpen():", s.data)

	return
}

// Save save current connection s.
func (s *ConnectionSession) Save() bool {
	// if !s.AllowSave {
	// return false
	// }
	if s.isErrorOccured() {
		return false
	}

	// TODO what about the connection has been deleted?

	if len(s.connPath) > 0 {
		// update connection data and activate it
		nmConn, err := nmNewSettingsConnection(s.connPath)
		if err != nil {
			logger.Error(err)
			return false
		}
		err = nmConn.Update(s.data)
		if err != nil {
			logger.Error(err)
			return false
		}
		nmActivateConnection(s.connPath, s.devPath)
	} else {
		// create new connection and activate it
		nmAddAndActivateConnection(s.data, s.devPath)
	}

	dbus.UnInstallObject(s)
	return true
}

func (s *ConnectionSession) isErrorOccured() bool {
	for _, v := range s.Errors {
		if len(v) > 1 {
			return true
		}
	}
	return false
}

// Close cancel current connection s.
func (s *ConnectionSession) Close() {
	dbus.UnInstallObject(s)
}

// listPages return supported pages for target connection type.
func (s *ConnectionSession) listPages() (pages []string) {
	switch s.connectionType {
	case typeWired:
		pages = []string{
			pageGeneral,
			// pageEthernet, // TODO
			pageIPv4,
			pageIPv6,
			pageSecurity,
		}
	case typeWireless:
		pages = []string{
			pageGeneral,
			// pageWifi, // TODO need when setup adhoc
			pageIPv4,
			pageIPv6,
			pageSecurity,
		}
	case typePppoe:
		pages = []string{
			pageGeneral,
			// pageEthernet, // TODO
			pagePppoe,
			pagePpp,
			pageIPv4,
		}
	case typeVpnL2tp:
		pages = []string{
			pageGeneral,
			pageVpnL2tp,
			pageVpnL2tpPpp,
			pageVpnL2tpIpsec,
			pageIPv4,
		}
	case typeVpnOpenconnect:
		pages = []string{
			pageGeneral,
			pageVpnOpenconnect,
			pageIPv4,
			pageIPv6,
		}
	case typeVpnOpenvpn:
		pages = []string{
			pageGeneral,
			pageVpnOpenvpn,
			pageVpnOpenvpnAdvanced,
			pageVpnOpenvpnSecurity,
			pageVpnOpenvpnProxies,
			pageIPv4,
			pageIPv6,
		}
		// TODO make "Pages" as a dbus property
		// when connection type is static key, pageVpnOpenvpnTlsauth is not availabled
		if getSettingVpnOpenvpnKeyConnectionType(s.data) != NM_OPENVPN_CONTYPE_STATIC_KEY {
			pages = append(pages, pageVpnOpenvpnTlsauth)
		}
	case typeVpnPptp:
		pages = []string{
			pageGeneral,
			pageVpnPptp,
			pageVpnPptpPpp,
			pageIPv4,
		}
	case typeVpnVpnc:
		pages = []string{
			pageGeneral,
			pageVpnVpnc,
			pageVpnVpncAdvanced,
			pageIPv4,
		}
	}
	return
}

func (s *ConnectionSession) pageToFields(page string) (fields []string) {
	switch page {
	default:
		logger.Error("pageToFields: invalid page name", page)
	case pageGeneral:
		fields = []string{fieldConnection}
	case pageEthernet:
		fields = []string{fieldWired}
	case pageWifi:
		fields = []string{fieldWireless}
	case pageIPv4:
		fields = []string{fieldIpv4}
	case pageIPv6:
		fields = []string{fieldIpv6}
	case pageSecurity:
		switch s.connectionType {
		case typeWired:
			fields = []string{field8021x}
		case typeWireless:
			if isSettingFieldExists(s.data, field8021x) {
				fields = []string{fieldWirelessSecurity, field8021x}
			} else {
				fields = []string{fieldWirelessSecurity}
			}
		}
	case pagePppoe:
		fields = []string{fieldPppoe}
	case pagePpp:
		fields = []string{fieldPpp}
	case pageVpnL2tp:
		fields = []string{fieldVpnL2tp}
	case pageVpnL2tpPpp:
		fields = []string{fieldVpnL2tpPpp}
	case pageVpnL2tpIpsec:
		fields = []string{fieldVpnL2tpIpsec}
	case pageVpnOpenconnect:
		fields = []string{fieldVpnOpenconnect}
	case pageVpnOpenvpn:
		fields = []string{fieldVpnOpenvpn}
	case pageVpnOpenvpnAdvanced:
		fields = []string{fieldVpnOpenvpnAdvanced}
	case pageVpnOpenvpnSecurity:
		fields = []string{fieldVpnOpenvpnSecurity}
	case pageVpnOpenvpnTlsauth:
		fields = []string{fieldVpnOpenvpnTlsauth}
	case pageVpnOpenvpnProxies:
		fields = []string{fieldVpnOpenvpnProxies}
	case pageVpnPptp:
		fields = []string{fieldVpnPptp}
	case pageVpnPptpPpp:
		fields = []string{fieldVpnPptpPpp}
	case pageVpnVpnc:
		fields = []string{fieldVpnVpnc}
	case pageVpnVpncAdvanced:
		fields = []string{fieldVpnVpncAdvanced}
	}
	return
}

func (s *ConnectionSession) getFieldOfPageKey(page, key string) string {
	fields := s.pageToFields(page)
	for _, field := range fields {
		if generalIsKeyInSettingField(field, key) {
			return field
		}
	}
	logger.Errorf("get corresponding filed of key in page failed, page=%s, key=%s", page, key)
	return ""
}

// get valid keys of target page, show or hide some keys when special
// keys toggled
func (s *ConnectionSession) listKeys(page string) (keys []string) {
	fields := s.pageToFields(page)
	for _, field := range fields {
		// TODO
		// if isSettingFieldExists(s.data, field) {
		// }
		keys = appendStrArrayUnion(keys, generalGetSettingAvailableKeys(s.data, field)...)
	}
	if len(keys) == 0 {
		logger.Warning("there is no avaiable keys for page", page)
	}
	return
}

// GetAvailableValues get available values for target key.
// TODO
// func (s *ConnectionSession) GetAvailableValues(page, key string) (valuesJSON []map[string]string) {
func (s *ConnectionSession) GetAvailableValues(page, key string) (valuesJSON string) {
	// TODO
	var values []kvalue
	fields := s.pageToFields(page)
	for _, field := range fields {
		values = generalGetSettingAvailableValues(s.data, field, key)
		if len(values) > 0 {
			break
		}
	}
	// TODO
	// convert kvalue to kvalueJSON
	// for _, v := range values {
	// valueMap := make(map[string]string)
	// valueMap["Value"], _ = marshalJSON(v.Value)
	// valueMap["Text"] = v.Text
	// valuesJSON = append(valuesJSON, valueMap)

	// json, _ := marshalJSON(v.Value)
	// valuesJSON = append(valuesJSON, kvalueJSON{json, v.Text})

	// }
	valuesJSON, _ = marshalJSON(values)
	return
}

func (s *ConnectionSession) GetKey(page, key string) (value string) {
	field := s.getFieldOfPageKey(page, key)
	value = generalGetSettingKeyJSON(s.data, field, key)
	return
}

func (s *ConnectionSession) SetKey(page, key, value string) {
	field := s.getFieldOfPageKey(page, key)
	// logger.Debugf("SetKey(), page=%s, filed=%s, key=%s, value=%s", page, field, key, value) // TODO test
	generalSetSettingKeyJSON(s.data, field, key, value)

	s.updatePropAvailablePages()
	s.updatePropAvailableKeys()
	s.updatePropErrors()

	// TODO
	// if s.isErrorOccured() {
	// 	s.updatePropAllowSave(false)
	// } else {
	// 	s.updatePropAllowSave(true)
	// }

	return
}

// TODO remove CheckValues check target value if is correct.
// func (s *ConnectionSession) CheckValue(page, key, value string) (ok bool) {
// 	return
// }

func (s *ConnectionSession) DebugListKeyDetail() (info string) {
	for _, page := range s.listPages() {
		pageData, ok := s.AvailableKeys[page]
		if !ok {
			logger.Warning("no available keys for page", page)
			continue
		}
		for _, key := range pageData {
			field := s.getFieldOfPageKey(page, key)
			t := generalGetSettingKeyType(field, key)
			// TODO convert to value json
			values := generalGetSettingAvailableValues(s.data, field, key)
			info += fmt.Sprintf("%s->%s[%s]: %s (%s)\n", page, key, getKtypeDescription(t), s.GetKey(page, key), values)
		}
	}
	return
}

func (s *ConnectionSession) DebugGetConnectionData() connectionData {
	return s.data
}

func (s *ConnectionSession) DebugGetErrors() PageKeyErrors {
	return s.Errors
}
