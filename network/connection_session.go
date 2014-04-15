package main

import (
	nm "dbus/org/freedesktop/networkmanager"
	"dlib/dbus"
	"fmt"
)

type ConnectionSession struct {
	coreObjPath    dbus.ObjectPath
	objPath        dbus.ObjectPath
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

func doNewConnectionSession() (session *ConnectionSession) {
	session = &ConnectionSession{}
	session.data = make(_ConnectionData)
	session.AllowSave = false // TODO
	session.AvailableKeys = make(map[string][]string)
	session.Errors = make(map[string]map[string]string)
	return session
}

func NewConnectionSessionByCreate(connectionType string) (session *ConnectionSession, err error) {
	if !isStringInArray(connectionType, supportedConnectionTypes) {
		err = fmt.Errorf("connection type is out of support: %s", connectionType)
		Logger.Error(err)
		return
	}

	session = doNewConnectionSession()
	session.CurrentUUID = newUUID()
	session.objPath = dbus.ObjectPath(fmt.Sprintf("/com/deepin/daemon/ConnectionSession/%s", randString(8)))

	// TODO
	// new connection data, id is left here
	session.connectionType = connectionType
	switch session.connectionType {
	case typeWired:
		session.data = newWiredConnectionData("", session.CurrentUUID)
	case typeWireless:
		session.data = newWirelessConnectionData("", session.CurrentUUID, nil, ApKeyNone)
	case typePppoe:
		session.data = newPppoeConnectionData("", session.CurrentUUID)
	}

	session.updatePropErrors()
	session.updatePropAvailableKeys()
	// session.updatePropAllowSave(false) // TODO

	return
}

func NewConnectionSessionByOpen(uuid string) (session *ConnectionSession, err error) {
	coreObjPath, err := NMSettings.GetConnectionByUuid(uuid)
	if err != nil {
		return
	}

	session = doNewConnectionSession()
	session.coreObjPath = coreObjPath
	session.CurrentUUID = uuid
	session.objPath = dbus.ObjectPath(fmt.Sprintf("/com/deepin/daemon/ConnectionSession/%s", randString(8)))

	// get connection data
	nmConn, err := nm.NewSettingsConnection(NMDest, coreObjPath)
	if err != nil {
		return nil, err
	}
	session.data, err = nmConn.GetSettings()
	if err != nil {
		return nil, err
	}
	session.connectionType = getSettingConnectionType(session.data)

	// get secret data
	for _, secretFiled := range []string{fieldWirelessSecurity, field8021x} {
		if isSettingFieldExists(session.data, secretFiled) {
			wirelessSecrutiyData, err := nmConn.GetSecrets(fieldWirelessSecurity)
			if err == nil {
				for field, fieldData := range wirelessSecrutiyData {
					if !isSettingFieldExists(session.data, field) {
						addSettingField(session.data, field)
					}
					for key, value := range fieldData {
						session.data[field][key] = value
					}
				}
			}
		}
	}

	session.updatePropErrors()
	session.updatePropAvailableKeys()
	// session.updatePropAllowSave(false) // TODO

	// TODO
	Logger.Debug("NewConnectionSessionByOpen():", session.data)

	return
}

// Save save current connection session.
func (session *ConnectionSession) Save() bool {
	// if !session.AllowSave {
	// return false
	// }
	if session.isErrorOccured() {
		return false
	}

	// TODO what about the connection has been deleted?

	// update connection data
	nmConn, err := nm.NewSettingsConnection(NMDest, session.coreObjPath)
	if err != nil {
		Logger.Error(err)
		return false
	}
	err = nmConn.Update(session.data)
	if err != nil {
		Logger.Error(err)
		return false
	}

	dbus.UnInstallObject(session)
	return true
}

func (session *ConnectionSession) isErrorOccured() bool {
	for _, v := range session.Errors {
		if len(v) > 1 {
			return true
		}
	}
	return false
}

// Close cancel current connection session.
func (session *ConnectionSession) Close() {
	dbus.UnInstallObject(session)
}

//根据CurrentUUID返回此Connection支持的设置页面
func (session *ConnectionSession) ListPages() (pages []string) {
	switch session.connectionType {
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

func (session *ConnectionSession) pageToFields(page string) (fields []string) {
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
		switch session.connectionType {
		case typeWired:
			fields = []string{field8021x}
		case typeWireless:
			if isSettingFieldExists(session.data, field8021x) {
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

func (session *ConnectionSession) getFieldOfPageKey(page, key string) string {
	fields := session.pageToFields(page)
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
func (session *ConnectionSession) listKeys(page string) (keys []string) {
	fields := session.pageToFields(page)
	for _, field := range fields {
		if isSettingFieldExists(session.data, field) {
			keys = appendStrArrayUnion(keys, generalGetSettingAvailableKeys(session.data, field)...)
		}
	}
	if len(keys) == 0 {
		Logger.Warning("there is no avaiable keys for page", page)
	}
	return
}

// GetAvailableValues get available values for target key.
func (session *ConnectionSession) GetAvailableValues(page, key string) (values []string) {
	fields := session.pageToFields(page)
	for _, field := range fields {
		values, _ = generalGetSettingAvailableValues(session.data, field, key)
		if len(values) > 0 {
			break
		}
	}
	return
}

func (session *ConnectionSession) GetKey(page, key string) (value string) {
	field := session.getFieldOfPageKey(page, key)
	value = generalGetSettingKeyJSON(session.data, field, key)
	return
}

func (session *ConnectionSession) SetKey(page, key, value string) {
	field := session.getFieldOfPageKey(page, key)
	generalSetSettingKeyJSON(session.data, field, key, value)
	session.updatePropErrors()
	session.updatePropAvailableKeys()
	// TODO
	// if session.isErrorOccured() {
	// 	session.updatePropAllowSave(false)
	// } else {
	// 	session.updatePropAllowSave(true)
	// }

	return
}

// TODO remove CheckValues check target value if is correct.
// func (session *ConnectionSession) CheckValue(page, key, value string) (ok bool) {
// 	return
// }

func (session *ConnectionSession) DebugListKeyDetail() (info string) {
	for _, page := range session.ListPages() {
		pageData, ok := session.AvailableKeys[page]
		if !ok {
			Logger.Warning("no available keys for page", page)
			continue
		}
		for _, key := range pageData {
			field := session.getFieldOfPageKey(page, key)
			t := generalGetSettingKeyType(field, key)
			values, _ := generalGetSettingAvailableValues(session.data, field, key)
			info += fmt.Sprintf("%s->%s[%s]: %s (%s)\n", page, key, getKtypeDescription(t), session.GetKey(page, key), values)
		}
	}
	return
}

func (session *ConnectionSession) DebugGetConnectionData() _ConnectionData {
	return session.data
}

func (session *ConnectionSession) DebugGetErrors() map[string]map[string]string {
	return session.Errors
}
