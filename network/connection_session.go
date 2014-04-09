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
		LOGGER.Error(err)
		return
	}

	session = doNewConnectionSession()
	session.CurrentUUID = newUUID()
	session.objPath = dbus.ObjectPath(fmt.Sprintf("/com/deepin/daemon/ConnectionSession/%s", randString(8)))

	// TODO
	session.connectionType = connectionType
	switch session.connectionType {
	case typeWired:
		session.data = newWireedConnectionData("", session.CurrentUUID)
	case typeWireless:
		session.data = newWirelessConnectionData("", session.CurrentUUID, nil, ApKeyNone)
	}

	session.updatePropErrors()
	session.updatePropAvailableKeys()
	// session.updatePropAllowSave(false) // TODO

	return
}

func NewConnectionSessionByOpen(uuid string) (session *ConnectionSession, err error) {
	coreObjPath, err := _NMSettings.GetConnectionByUuid(uuid)
	if err != nil {
		err = fmt.Errorf("counld not find connection with uuid equal %s", uuid)
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

	session.updatePropErrors()
	session.updatePropAvailableKeys()
	// session.updatePropAllowSave(false) // TODO

	// TODO
	LOGGER.Debug("NewConnectionSessionByOpen():", session.data)

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
		LOGGER.Error(err)
		return false
	}
	err = nmConn.Update(session.data)
	if err != nil {
		LOGGER.Error(err)
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
	}
	return
}

// get valid keys of target page, show or hide some keys when special
// keys toggled
func (session *ConnectionSession) listKeys(page string) (keys []string) {
	switch page {
	case pageGeneral:
		keys = getSettingConnectionAvailableKeys(session.data)
	case pageEthernet:
		keys = getSettingWiredAvailableKeys(session.data)
	case pageWifi:
		keys = getSettingWirelessAvailableKeys(session.data)
	case pageIPv4:
		keys = getSettingIp4ConfigAvailableKeys(session.data)
	case pageIPv6:
		keys = getSettingIp6ConfigAvailableKeys(session.data)
	case pageSecurity: // TODO
		switch session.connectionType {
		case typeWired:
			keys = getSetting8021xAvailableKeys(session.data)
		case typeWireless:
			securityField := getSettingWirelessSec(session.data)
			switch securityField {
			case fieldWirelessSecurity:
				keys = getSettingWirelessSecurityAvailableKeys(session.data)
			case field8021x:
				keys = getSetting8021xAvailableKeys(session.data)
			}
		}
	}
	if len(keys) == 0 {
		LOGGER.Warning("there is no avaiable keys for page", page)
	}
	return
}

// GetAvailableValues get available values for target key.
func (session *ConnectionSession) GetAvailableValues(page, key string) (values []string) {
	switch page {
	case pageGeneral:
		values, _ = getSettingConnectionAvailableValues(key)
	case pageIPv4:
		values, _ = getSettingIp4ConfigAvailableValues(key)
	case pageIPv6:
		values, _ = getSettingIp6ConfigAvailableValues(key)
	case pageSecurity: // TODO
		switch session.connectionType {
		case typeWired:
			values, _ = getSetting8021xAvailableValues(key)
		case typeWireless:
			securityField := getSettingWirelessSec(session.data)
			switch securityField {
			case fieldWirelessSecurity:
				values, _ = getSettingWirelessSecurityAvailableValues(key)
			case field8021x:
				values, _ = getSetting8021xAvailableValues(key)
			}
		}
	}
	return
}

func (session *ConnectionSession) pageToField(page string) (field string) {
	switch page {
	default:
		LOGGER.Error("pageToField: invalid page name", page)
	case pageGeneral:
		field = fieldConnection
	case pageEthernet:
		field = fieldWired
	case pageWifi:
		field = fieldWireless
	case pageIPv4:
		field = fieldIPv4
	case pageIPv6:
		field = fieldIPv6
	case pageSecurity: // TODO
		switch session.connectionType {
		case typeWired:
			field = field8021x
		case typeWireless:
			securityField := getSettingWirelessSec(session.data)
			switch securityField {
			case fieldWirelessSecurity:
				field = fieldWirelessSecurity
			case field8021x:
				field = field8021x
			}
		}
	}
	return
}

func (session *ConnectionSession) GetKey(page, key string) (value string) {
	field := session.pageToField(page)
	if isVirtualKey(field, key) {
		value = generalGetVirtualKeyJSON(session.data, field, key)
	} else {
		value = generalGetKeyJSON(session.data, field, key)
	}
	return
}

func (session *ConnectionSession) SetKey(page, key, value string) {
	field := session.pageToField(page)
	if isVirtualKey(field, key) {
		generalSetVirtualKeyJSON(session.data, field, key, value)
	} else {
		generalSetKeyJSON(session.data, field, key, value)
	}
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

// DebugConnectionTypes return all supported connection types, only for debugging.
// TODO move to manager
func (session *ConnectionSession) DebugListSupportedConnectionTypes() []string {
	return supportedConnectionTypes
}
