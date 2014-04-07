package main

import (
	nm "dbus/org/freedesktop/networkmanager"
	"dlib/dbus"
	"fmt"
)

type ConnectionSession struct {
	coreObjPath dbus.ObjectPath
	objPath     dbus.ObjectPath
	data        _ConnectionData
	connType    string

	CurrentUUID string

	AllowSave bool // TODO really needed?

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

func NewConnectionSessionByCreate(connType string) (session *ConnectionSession, err error) {
	if !isStringInArray(connType, supportedConnectionTypes) {
		err = fmt.Errorf("connection type is out of support: %s", connType)
		LOGGER.Error(err)
		return
	}

	session = doNewConnectionSession()
	session.CurrentUUID = newUUID()
	session.objPath = dbus.ObjectPath(fmt.Sprintf("/com/deepin/daemon/ConnectionSession/%s", randString(8)))

	// TODO
	session.connType = connType

	session.updatePropErrors()
	session.updatePropAvailableKeys()
	session.updatePropAllowSave(false)

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
	session.connType = getSettingConnectionType(session.data)

	session.updatePropErrors()
	session.updatePropAvailableKeys()
	session.updatePropAllowSave(false)

	// TODO
	LOGGER.Debug("NewConnectionSessionByOpen():", session.data)

	return
}

// Save save current connection session.
func (session *ConnectionSession) Save() bool {
	if !session.AllowSave {
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
	switch session.connType {
	case typeWired:
		pages = []string{
			pageGeneral,
			pageIPv4,
			pageIPv6,
		}
	case typeWireless:
		pages = []string{
			pageGeneral,
			// pageWifi, // TODO
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
		switch session.connType {
		case typeWired:
		case typeWireless:
			// TODO
			keys = getSettingWirelessSecurityAvailableKeys(session.data)
		}
	}
	if len(keys) == 0 {
		LOGGER.Warning("there is no avaiable keys for page", page)
	}
	return
}

//比如获得当前链接支持的加密方式 EAP字段: TLS、MD5、FAST、PEAP
//获得ip设置方式 : Manual、Link-Local Only、Automatic(DHCP)
//获得当前可用mac地址(这种字段是有几个可选值但用户也可用手动输入一个其他值)
// TODO
func (session *ConnectionSession) GetAvailableValues(page, key string) (values []string) {
	switch page {
	case pageGeneral:
	case pageIPv4:
		values, _ = getSettingIp4ConfigAvailableValues(key)
	case pageIPv6:
		values, _ = getSettingIp6ConfigAvailableValues(key)
	case pageSecurity: // TODO
		// switch session.connType {
		// case typeWired:
		// case typeWireless:
		// 	// TODO
		// }
	}
	return
}

func (session *ConnectionSession) GetKey(page, key string) (value string) {
	switch page {
	default:
		LOGGER.Error("GetKey: invalid page name", page)
	case pageGeneral:
		value = generalGetSettingConnectionKeyJSON(session.data, key)
	case pageEthernet:
		value = generalGetSettingWiredKeyJSON(session.data, key)
	case pageWifi:
		value = generalGetSettingWirelessKeyJSON(session.data, key)
	case pageIPv4:
		value = generalGetSettingIp4ConfigKeyJSON(session.data, key)
	case pageIPv6:
		value = generalGetSettingIp6ConfigKeyJSON(session.data, key)
	case pageSecurity: // TODO
		switch session.connType {
		case typeWired:
		case typeWireless:
			// switch method {
			// value = generalGetSettingWirelessSecurityKeyJSON(session.data, key)
			// value = generalGetSetting8021xKeyJSON(session.data, key)
			// }
		}
	}
	return
}

func (session *ConnectionSession) SetKey(page, key, value string) {
	switch page {
	default:
		LOGGER.Error("SetKey: invalid page name", page)
	case pageGeneral:
		generalSetSettingConnectionKeyJSON(session.data, key, value)
	case pageEthernet:
		generalSetSettingWiredKeyJSON(session.data, key, value)
	case pageWifi:
		generalSetSettingWirelessKeyJSON(session.data, key, value)
	case pageIPv4:
		generalSetSettingIp4ConfigKeyJSON(session.data, key, value)
	case pageIPv6:
		generalSetSettingIp6ConfigKeyJSON(session.data, key, value)
	case pageSecurity: // TODO
		switch session.connType {
		case typeWired:
		case typeWireless:
			// switch method {
			// generalSetSettingWirelessSecurityKeyJSON(session.data, key, value)
			// generalSetSetting8021xKeyJSON(session.data, key, value)
			// }
		}
	}

	session.updatePropErrors()
	session.updatePropAvailableKeys()
	if session.isErrorOccured() {
		session.updatePropAllowSave(false)
	} else {
		session.updatePropAllowSave(true)
	}

	return
}

// TODO CheckValues check target value if is correct.
// func (session *ConnectionSession) CheckValue(page, key, value string) (ok bool) {
// 	return
// }

// DebugConnectionTypes return all supported connection types, only for debugging.
// TODO move to manager
func (session *ConnectionSession) DebugListSupportedConnectionTypes() []string {
	return supportedConnectionTypes
}
