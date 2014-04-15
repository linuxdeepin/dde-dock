package main

import (
	"dlib/dbus"
	"fmt"
)

type ConnectionSession struct {
	sessionPath    dbus.ObjectPath
	connPath       dbus.ObjectPath
	devPath        dbus.ObjectPath
	data           _ConnectionData
	connectionType string

	CurrentUUID string

	AllowSave bool // TODO really need?

	// 前端只显示此列表中的字段, 会跟随当前正在编辑的值而改变
	// TODO more documentation
	AvailableKeys map[string][]string

	// 返回所有 page 下错误的字段和对应的错误原因
	Errors map[string]map[string]string
}

//所有字段值都为string，后端自行转换为需要的值后提供给NM

func doNewConnectionSession(devPath dbus.ObjectPath, uuid string) (s *ConnectionSession) {
	s = &ConnectionSession{}
	s.sessionPath = dbus.ObjectPath(fmt.Sprintf("/com/deepin/daemon/ConnectionSession/%s", randString(8)))
	s.devPath = devPath
	s.CurrentUUID = uuid
	s.data = make(_ConnectionData)
	s.AllowSave = false // TODO
	s.AvailableKeys = make(map[string][]string)
	s.Errors = make(map[string]map[string]string)
	return s
}

func NewConnectionSessionByCreate(connectionType string, devPath dbus.ObjectPath) (s *ConnectionSession, err error) {
	if !isStringInArray(connectionType, supportedConnectionTypes) {
		err = fmt.Errorf("connection type is out of support: %s", connectionType)
		Logger.Error(err)
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
		s.data = newWirelessConnectionData("", s.CurrentUUID, nil, ApSecNone)
	case typePppoe:
		s.data = newPppoeConnectionData("", s.CurrentUUID)
	}

	s.updatePropErrors()
	s.updatePropAvailableKeys()
	// s.updatePropAllowSave(false) // TODO

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
	s.connectionType = getSettingConnectionType(s.data)

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

	s.updatePropErrors()
	s.updatePropAvailableKeys()
	// s.updatePropAllowSave(false) // TODO

	// TODO
	Logger.Debug("NewConnectionSessionByOpen():", s.data)

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
			Logger.Error(err)
			return false
		}
		err = nmConn.Update(s.data)
		if err != nil {
			Logger.Error(err)
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

//根据CurrentUUID返回此Connection支持的设置页面
func (s *ConnectionSession) ListPages() (pages []string) {
	switch s.connectionType {
	case typeWired:
		pages = []string{
			pageGeneral,
			pageIPv4,
			pageIPv6,
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
			pagePppoe,
			pageIPv4,
			// pagePpp,			// TODO if need
		}
	}
	return
}

func (s *ConnectionSession) pageToFields(page string) (fields []string) {
	switch page {
	default:
		Logger.Error("pageToFields: invalid page name", page)
	case pageGeneral:
		fields = []string{fieldConnection}
	case pageEthernet:
		fields = []string{fieldWired}
	case pageWifi:
		fields = []string{fieldWireless}
	case pageIPv4:
		fields = []string{fieldIPv4}
	case pageIPv6:
		fields = []string{fieldIPv6}
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
	Logger.Errorf("get corresponding filed of key in page failed, page=%s, key=%s", page, key)
	return ""
}

// get valid keys of target page, show or hide some keys when special
// keys toggled
func (s *ConnectionSession) listKeys(page string) (keys []string) {
	fields := s.pageToFields(page)
	for _, field := range fields {
		if isSettingFieldExists(s.data, field) {
			keys = appendStrArrayUnion(keys, generalGetSettingAvailableKeys(s.data, field)...)
		}
	}
	if len(keys) == 0 {
		Logger.Warning("there is no avaiable keys for page", page)
	}
	return
}

// GetAvailableValues get available values for target key.
func (s *ConnectionSession) GetAvailableValues(page, key string) (values []string) {
	fields := s.pageToFields(page)
	for _, field := range fields {
		values, _ = generalGetSettingAvailableValues(s.data, field, key)
		if len(values) > 0 {
			break
		}
	}
	return
}

func (s *ConnectionSession) GetKey(page, key string) (value string) {
	field := s.getFieldOfPageKey(page, key)
	value = generalGetSettingKeyJSON(s.data, field, key)
	return
}

func (s *ConnectionSession) SetKey(page, key, value string) {
	field := s.getFieldOfPageKey(page, key)
	generalSetSettingKeyJSON(s.data, field, key, value)
	s.updatePropErrors()
	s.updatePropAvailableKeys()
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
	for _, page := range s.ListPages() {
		pageData, ok := s.AvailableKeys[page]
		if !ok {
			Logger.Warning("no available keys for page", page)
			continue
		}
		for _, key := range pageData {
			field := s.getFieldOfPageKey(page, key)
			t := generalGetSettingKeyType(field, key)
			values, _ := generalGetSettingAvailableValues(s.data, field, key)
			info += fmt.Sprintf("%s->%s[%s]: %s (%s)\n", page, key, getKtypeDescription(t), s.GetKey(page, key), values)
		}
	}
	return
}

func (s *ConnectionSession) DebugGetConnectionData() _ConnectionData {
	return s.data
}

func (s *ConnectionSession) DebugGetErrors() map[string]map[string]string {
	return s.Errors
}
