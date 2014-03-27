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
	currentUUID string // TODO hide property
	currentPage string

	HasChanged bool

	//前端只显示此列表中的字段,会跟随当前正在编辑的值而改变
	// TODO more documentation
	CurrentFields []string
	//返回当前page下错误的字段和对应的错误原因
	CurrentErrors []string
}

//所有字段值都为string，后端自行转换为需要的值后提供给NM

func NewConnectionSessionByCreate(connType string) (session *ConnectionSession, err error) {
	if !isStringInArray(connType, supportedConnectionTypes) {
		err = fmt.Errorf("connection type is out of support: %s", connType)
		LOGGER.Error(err)
		return
	}

	session = &ConnectionSession{}
	session.currentUUID = newUUID()
	session.objPath = dbus.ObjectPath(fmt.Sprintf("/com/deepin/daemon/ConnectionSession/%s", randString(8)))

	// TODO
	session.data = make(_ConnectionData)
	session.connType = connType

	// session.updatePropCurrentUUID(uuid)
	// session.updatePropHasChanged(true)

	// TODO
	// session.currentPage = session.getDefaultPage(connType)
	// session.updatePropCurrentFields()

	// TODO current errors
	return
}

func NewConnectionSessionByOpen(uuid string) (session *ConnectionSession, err error) {
	coreObjPath, ok := _Manager.getConnectionPathByUUID(uuid)
	if !ok {
		err = fmt.Errorf("counld not find connection with uuid equal %s", uuid)
		return
	}

	session = &ConnectionSession{}
	session.coreObjPath = coreObjPath
	session.currentUUID = uuid
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

	return
}

// Save save current connection session.
func (session *ConnectionSession) Save() {
	// TODO
	if !session.HasChanged {
		dbus.UnInstallObject(session)
		return
	}

	// TODO error fields

	dbus.UnInstallObject(session)
}

// Cancel cancel current connection session.
func (session *ConnectionSession) Cancel() {
	dbus.UnInstallObject(session)
}

//根据CurrentUUID返回此Connection支持的设置页面
func (session *ConnectionSession) ListPages() (pages []string) {
	// TODO
	switch session.connType {
	case typeWired:
		pages = []string{
			pageGeneral,
			pageEthernet,
			pageIPv4,
			pageIPv6,
		}
	case typeWireless:
		pages = []string{
			pageGeneral,
			pageWifi,
			pageIPv4,
			pageIPv6,
			pageSecurity,
		}
	}
	return
}

// get valid fields for target page
func (session *ConnectionSession) listFields(page string) (fields []string) {
	switch session.connType {
	case typeWired:
		switch page {
		case pageGeneral:
			fields = []string{
				NM_SETTING_CONNECTION_ID,
				NM_SETTING_CONNECTION_UUID,
				NM_SETTING_CONNECTION_TYPE,
				NM_SETTING_CONNECTION_AUTOCONNECT,
				NM_SETTING_CONNECTION_TIMESTAMP,
				NM_SETTING_CONNECTION_READ_ONLY,
				NM_SETTING_CONNECTION_PERMISSIONS,
				NM_SETTING_CONNECTION_ZONE,
				NM_SETTING_CONNECTION_MASTER,
				NM_SETTING_CONNECTION_SLAVE_TYPE,
				NM_SETTING_CONNECTION_SECONDARIES,
			}
		case pageEthernet:
			fields = []string{
				NM_SETTING_WIRED_PORT,
				NM_SETTING_WIRED_SPEED,
				NM_SETTING_WIRED_DUPLEX,
				NM_SETTING_WIRED_AUTO_NEGOTIATE,
				NM_SETTING_WIRED_MAC_ADDRESS,
				NM_SETTING_WIRED_CLONED_MAC_ADDRESS,
				NM_SETTING_WIRED_MAC_ADDRESS_BLACKLIST,
				NM_SETTING_WIRED_MTU,
				NM_SETTING_WIRED_S390_SUBCHANNELS,
				NM_SETTING_WIRED_S390_NETTYPE,
				NM_SETTING_WIRED_S390_OPTIONS,
			}
		case pageIPv4:
			fields = []string{
				NM_SETTING_IP4_CONFIG_METHOD,
				NM_SETTING_IP4_CONFIG_DNS,
				NM_SETTING_IP4_CONFIG_DNS_SEARCH,
				NM_SETTING_IP4_CONFIG_ADDRESSES,
				NM_SETTING_IP4_CONFIG_ROUTES,
				NM_SETTING_IP4_CONFIG_IGNORE_AUTO_ROUTES,
				NM_SETTING_IP4_CONFIG_IGNORE_AUTO_DNS,
				NM_SETTING_IP4_CONFIG_DHCP_CLIENT_ID,
				NM_SETTING_IP4_CONFIG_DHCP_SEND_HOSTNAME,
				NM_SETTING_IP4_CONFIG_DHCP_HOSTNAME,
				NM_SETTING_IP4_CONFIG_NEVER_DEFAULT,
				NM_SETTING_IP4_CONFIG_MAY_FAIL,
			}
		case pageIPv6:
			fields = []string{
				NM_SETTING_IP6_CONFIG_METHOD,
				NM_SETTING_IP6_CONFIG_DNS,
				NM_SETTING_IP6_CONFIG_DNS_SEARCH,
				NM_SETTING_IP6_CONFIG_ADDRESSES,
				NM_SETTING_IP6_CONFIG_ROUTES,
				NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES,
				NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS,
				NM_SETTING_IP6_CONFIG_NEVER_DEFAULT,
				NM_SETTING_IP6_CONFIG_MAY_FAIL,
				NM_SETTING_IP6_CONFIG_IP6_PRIVACY,
				NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME,
			}
		case pageSecurity:
			fields = []string{
				NM_SETTING_WIRELESS_SECURITY_KEY_MGMT,
				NM_SETTING_WIRELESS_SECURITY_WEP_TX_KEYIDX,
				NM_SETTING_WIRELESS_SECURITY_AUTH_ALG,
				NM_SETTING_WIRELESS_SECURITY_PROTO,
				NM_SETTING_WIRELESS_SECURITY_PAIRWISE,
				NM_SETTING_WIRELESS_SECURITY_GROUP,
				NM_SETTING_WIRELESS_SECURITY_LEAP_USERNAME,
				NM_SETTING_WIRELESS_SECURITY_WEP_KEY0,
				NM_SETTING_WIRELESS_SECURITY_WEP_KEY1,
				NM_SETTING_WIRELESS_SECURITY_WEP_KEY2,
				NM_SETTING_WIRELESS_SECURITY_WEP_KEY3,
				NM_SETTING_WIRELESS_SECURITY_WEP_KEY_FLAGS,
				NM_SETTING_WIRELESS_SECURITY_WEP_KEY_TYPE,
				NM_SETTING_WIRELESS_SECURITY_PSK,
				NM_SETTING_WIRELESS_SECURITY_PSK_FLAGS,
				NM_SETTING_WIRELESS_SECURITY_LEAP_PASSWORD,
				NM_SETTING_WIRELESS_SECURITY_LEAP_PASSWORD_FLAGS,
			}
		}
	case typeWireless:
		// TODO
	}
	return
}

//设置/获得字段的值都受这里设置page的影响。
func (session *ConnectionSession) SwitchPage(page string) {
	// TODO HasChanged
	session.currentPage = page
	session.updatePropCurrentFields()
}

//比如获得当前链接支持的加密方式 EAP字段: TLS、MD5、FAST、PEAP
//获得ip设置方式 : Manual、Link-Local Only、Automatic(DHCP)
//获得当前可用mac地址(这种字段是有几个可选值但用户也可用手动输入一个其他值)
func (session *ConnectionSession) GetAvailableValue(key string) (values []string) {
	// TODO
	switch key {
	case NM_SETTING_IP4_CONFIG_METHOD:
		values = []string{
			NM_SETTING_IP4_CONFIG_METHOD_AUTO,
			NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL,
			NM_SETTING_IP4_CONFIG_METHOD_MANUAL,
			NM_SETTING_IP4_CONFIG_METHOD_SHARED,
		}
	case NM_SETTING_IP4_CONFIG_DNS:
		values = []string{}
	}
	return
}

//仅仅调试使用，返回某个页面支持的字段。 因为字段如何安排(位置、我们是否要提供这个字段)是由前端决定的。
//*****在设计前端显示内容的时候和这个返回值关联很大*****
// DebugListFields return all fields of current page, only for debugging.
func (session *ConnectionSession) DebugListFields(page string) []string {
	// TODO
	return session.listFields(session.currentPage)
}

// DebugConnectionTypes return all supported connection types, only for debugging.
// TODO move to manager
func (session *ConnectionSession) DebugGetSupportedConnectionTypes() []string {
	return supportedConnectionTypes
}

// TODO
// func (session *ConnectionSession) DebugGetKeyType(key string) ktypeDescription {
// 	return ktypeDescriptions[0]
// }

// TODO panic error
// func (*Manager) DebugListKeyTypes() (uint32, string) {
// 	LOGGER.Debug(ktypeDescriptions)
// 	return ktypeDescriptions[0].t, ktypeDescriptions[0].desc
// }

//设置某个字段， 会影响CurrentFields属性，某些值会导致其他属性进入不可用状态
// TODO SetField
func (session *ConnectionSession) SetKey(key, value string) {
	switch session.connType {
	case typeWired:
		switch session.currentPage {
		case pageGeneral:
			switch key {
			case NM_SETTING_CONNECTION_ID:
				setSettingConnectionId(session.data, value)
			case NM_SETTING_CONNECTION_UUID:
				setSettingConnectionUuid(session.data, value)
			case NM_SETTING_CONNECTION_TYPE:
				setSettingConnectionType(session.data, value)
			case NM_SETTING_CONNECTION_AUTOCONNECT:
				setSettingConnectionAutoconnect(session.data, value)
			case NM_SETTING_CONNECTION_TIMESTAMP:
				setSettingConnectionTimestamp(session.data, value)
			case NM_SETTING_CONNECTION_READ_ONLY:
				setSettingConnectionReadOnly(session.data, value)
			case NM_SETTING_CONNECTION_PERMISSIONS:
				setSettingConnectionPermissions(session.data, value)
			case NM_SETTING_CONNECTION_ZONE:
				setSettingConnectionZone(session.data, value)
			case NM_SETTING_CONNECTION_MASTER:
				setSettingConnectionMaster(session.data, value)
			case NM_SETTING_CONNECTION_SLAVE_TYPE:
				setSettingConnectionSlaveType(session.data, value)
			case NM_SETTING_CONNECTION_SECONDARIES:
				setSettingConnectionSecondaries(session.data, value)
			}
		case pageEthernet:
			switch key {
			case NM_SETTING_WIRED_PORT:
				setSettingWiredPort(session.data, value)
			case NM_SETTING_WIRED_SPEED:
				setSettingWiredSpeed(session.data, value)
			case NM_SETTING_WIRED_DUPLEX:
				setSettingWiredDuplex(session.data, value)
			case NM_SETTING_WIRED_AUTO_NEGOTIATE:
				setSettingWiredAutoNegotiate(session.data, value)
			case NM_SETTING_WIRED_MAC_ADDRESS:
				setSettingWiredMacAddress(session.data, value)
			case NM_SETTING_WIRED_CLONED_MAC_ADDRESS:
				setSettingWiredClonedMacAddress(session.data, value)
			case NM_SETTING_WIRED_MAC_ADDRESS_BLACKLIST:
				setSettingWiredMacAddressBlacklist(session.data, value)
			case NM_SETTING_WIRED_MTU:
				setSettingWiredMtu(session.data, value)
			case NM_SETTING_WIRED_S390_SUBCHANNELS:
				setSettingWiredS390Subchannels(session.data, value)
			case NM_SETTING_WIRED_S390_NETTYPE:
				setSettingWiredS390Nettype(session.data, value)
			case NM_SETTING_WIRED_S390_OPTIONS:
				setSettingWiredS390Options(session.data, value)
			}
		case pageIPv4:
			switch key {
			case NM_SETTING_IP4_CONFIG_METHOD:
				setSettingIp4ConfigMethod(session.data, value)
			case NM_SETTING_IP4_CONFIG_DNS:
				setSettingIp4ConfigDns(session.data, value)
			case NM_SETTING_IP4_CONFIG_DNS_SEARCH:
				setSettingIp4ConfigDnsSearch(session.data, value)
			case NM_SETTING_IP4_CONFIG_ADDRESSES:
				setSettingIp4ConfigAddresses(session.data, value)
			case NM_SETTING_IP4_CONFIG_ROUTES:
				setSettingIp4ConfigRoutes(session.data, value)
			case NM_SETTING_IP4_CONFIG_IGNORE_AUTO_ROUTES:
				setSettingIp4ConfigIgnoreAutoRoutes(session.data, value)
			case NM_SETTING_IP4_CONFIG_IGNORE_AUTO_DNS:
				setSettingIp4ConfigIgnoreAutoDns(session.data, value)
			case NM_SETTING_IP4_CONFIG_DHCP_CLIENT_ID:
				setSettingIp4ConfigDhcpClientId(session.data, value)
			case NM_SETTING_IP4_CONFIG_DHCP_SEND_HOSTNAME:
				setSettingIp4ConfigDhcpSendHostname(session.data, value)
			case NM_SETTING_IP4_CONFIG_DHCP_HOSTNAME:
				setSettingIp4ConfigDhcpHostname(session.data, value)
			case NM_SETTING_IP4_CONFIG_NEVER_DEFAULT:
				setSettingIp4ConfigNeverDefault(session.data, value)
			case NM_SETTING_IP4_CONFIG_MAY_FAIL:
				setSettingIp4ConfigMayFail(session.data, value)
			}
		case pageIPv6:
			// TODO
			switch key {
			case NM_SETTING_IP6_CONFIG_METHOD:
			case NM_SETTING_IP6_CONFIG_DNS:
			case NM_SETTING_IP6_CONFIG_DNS_SEARCH:
			case NM_SETTING_IP6_CONFIG_ADDRESSES:
			case NM_SETTING_IP6_CONFIG_ROUTES:
			case NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES:
			case NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS:
			case NM_SETTING_IP6_CONFIG_NEVER_DEFAULT:
			case NM_SETTING_IP6_CONFIG_MAY_FAIL:
			case NM_SETTING_IP6_CONFIG_IP6_PRIVACY:
			case NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME:
			}
		case pageSecurity:
			// TODO
			switch key {
			case NM_SETTING_WIRELESS_SECURITY_KEY_MGMT:
			case NM_SETTING_WIRELESS_SECURITY_WEP_TX_KEYIDX:
			case NM_SETTING_WIRELESS_SECURITY_AUTH_ALG:
			case NM_SETTING_WIRELESS_SECURITY_PROTO:
			case NM_SETTING_WIRELESS_SECURITY_PAIRWISE:
			case NM_SETTING_WIRELESS_SECURITY_GROUP:
			case NM_SETTING_WIRELESS_SECURITY_LEAP_USERNAME:
			case NM_SETTING_WIRELESS_SECURITY_WEP_KEY0:
			case NM_SETTING_WIRELESS_SECURITY_WEP_KEY1:
			case NM_SETTING_WIRELESS_SECURITY_WEP_KEY2:
			case NM_SETTING_WIRELESS_SECURITY_WEP_KEY3:
			case NM_SETTING_WIRELESS_SECURITY_WEP_KEY_FLAGS:
			case NM_SETTING_WIRELESS_SECURITY_WEP_KEY_TYPE:
			case NM_SETTING_WIRELESS_SECURITY_PSK:
			case NM_SETTING_WIRELESS_SECURITY_PSK_FLAGS:
			case NM_SETTING_WIRELESS_SECURITY_LEAP_PASSWORD:
			case NM_SETTING_WIRELESS_SECURITY_LEAP_PASSWORD_FLAGS:
			}
		}
	case typeWireless:
		// TODO
	}
	return
}

// TODO
func (session *ConnectionSession) GetKey(key string) (value string) {
	switch session.connType {
	case typeWired:
		switch session.currentPage {
		case pageGeneral:
			switch key {
			case NM_SETTING_CONNECTION_ID:
				value = getSettingConnectionId(session.data)
			case NM_SETTING_CONNECTION_UUID:
				value = getSettingConnectionUuid(session.data)
			case NM_SETTING_CONNECTION_TYPE:
				value = getSettingConnectionType(session.data)
			case NM_SETTING_CONNECTION_AUTOCONNECT:
				value = getSettingConnectionAutoconnect(session.data)
			case NM_SETTING_CONNECTION_TIMESTAMP:
				value = getSettingConnectionTimestamp(session.data)
			case NM_SETTING_CONNECTION_READ_ONLY:
				value = getSettingConnectionReadOnly(session.data)
			case NM_SETTING_CONNECTION_PERMISSIONS:
				value = getSettingConnectionPermissions(session.data)
			case NM_SETTING_CONNECTION_ZONE:
				value = getSettingConnectionZone(session.data)
			case NM_SETTING_CONNECTION_MASTER:
				value = getSettingConnectionMaster(session.data)
			case NM_SETTING_CONNECTION_SLAVE_TYPE:
				value = getSettingConnectionSlaveType(session.data)
			case NM_SETTING_CONNECTION_SECONDARIES:
				value = getSettingConnectionSecondaries(session.data)
			}
		case pageEthernet:
			switch key {
			case NM_SETTING_WIRED_PORT:
				value = getSettingWiredPort(session.data)
			case NM_SETTING_WIRED_SPEED:
				value = getSettingWiredSpeed(session.data)
			case NM_SETTING_WIRED_DUPLEX:
				value = getSettingWiredDuplex(session.data)
			case NM_SETTING_WIRED_AUTO_NEGOTIATE:
				value = getSettingWiredAutoNegotiate(session.data)
			case NM_SETTING_WIRED_MAC_ADDRESS:
				value = getSettingWiredMacAddress(session.data)
			case NM_SETTING_WIRED_CLONED_MAC_ADDRESS:
				value = getSettingWiredClonedMacAddress(session.data)
			case NM_SETTING_WIRED_MAC_ADDRESS_BLACKLIST:
				value = getSettingWiredMacAddressBlacklist(session.data)
			case NM_SETTING_WIRED_MTU:
				value = getSettingWiredMtu(session.data)
			case NM_SETTING_WIRED_S390_SUBCHANNELS:
				value = getSettingWiredS390Subchannels(session.data)
			case NM_SETTING_WIRED_S390_NETTYPE:
				value = getSettingWiredS390Nettype(session.data)
			case NM_SETTING_WIRED_S390_OPTIONS:
				value = getSettingWiredS390Options(session.data)
			}
		case pageIPv4:
			switch key {
			case NM_SETTING_IP4_CONFIG_METHOD:
				value = getSettingIp4ConfigMethod(session.data)
			case NM_SETTING_IP4_CONFIG_DNS:
				value = getSettingIp4ConfigDns(session.data)
			case NM_SETTING_IP4_CONFIG_DNS_SEARCH:
				value = getSettingIp4ConfigDnsSearch(session.data)
			case NM_SETTING_IP4_CONFIG_ADDRESSES:
				value = getSettingIp4ConfigAddresses(session.data)
			case NM_SETTING_IP4_CONFIG_ROUTES:
				value = getSettingIp4ConfigRoutes(session.data)
			case NM_SETTING_IP4_CONFIG_IGNORE_AUTO_ROUTES:
				value = getSettingIp4ConfigIgnoreAutoRoutes(session.data)
			case NM_SETTING_IP4_CONFIG_IGNORE_AUTO_DNS:
				value = getSettingIp4ConfigIgnoreAutoDns(session.data)
			case NM_SETTING_IP4_CONFIG_DHCP_CLIENT_ID:
				value = getSettingIp4ConfigDhcpClientId(session.data)
			case NM_SETTING_IP4_CONFIG_DHCP_SEND_HOSTNAME:
				value = getSettingIp4ConfigDhcpSendHostname(session.data)
			case NM_SETTING_IP4_CONFIG_DHCP_HOSTNAME:
				value = getSettingIp4ConfigDhcpHostname(session.data)
			case NM_SETTING_IP4_CONFIG_NEVER_DEFAULT:
				value = getSettingIp4ConfigNeverDefault(session.data)
			case NM_SETTING_IP4_CONFIG_MAY_FAIL:
				value = getSettingIp4ConfigMayFail(session.data)
			}
		case pageIPv6:
			// TODO
			switch key {
			case NM_SETTING_IP6_CONFIG_METHOD:
			case NM_SETTING_IP6_CONFIG_DNS:
			case NM_SETTING_IP6_CONFIG_DNS_SEARCH:
			case NM_SETTING_IP6_CONFIG_ADDRESSES:
			case NM_SETTING_IP6_CONFIG_ROUTES:
			case NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES:
			case NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS:
			case NM_SETTING_IP6_CONFIG_NEVER_DEFAULT:
			case NM_SETTING_IP6_CONFIG_MAY_FAIL:
			case NM_SETTING_IP6_CONFIG_IP6_PRIVACY:
			case NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME:
			}
		case pageSecurity:
			// TODO
			switch key {
			case NM_SETTING_WIRELESS_SECURITY_KEY_MGMT:
			case NM_SETTING_WIRELESS_SECURITY_WEP_TX_KEYIDX:
			case NM_SETTING_WIRELESS_SECURITY_AUTH_ALG:
			case NM_SETTING_WIRELESS_SECURITY_PROTO:
			case NM_SETTING_WIRELESS_SECURITY_PAIRWISE:
			case NM_SETTING_WIRELESS_SECURITY_GROUP:
			case NM_SETTING_WIRELESS_SECURITY_LEAP_USERNAME:
			case NM_SETTING_WIRELESS_SECURITY_WEP_KEY0:
			case NM_SETTING_WIRELESS_SECURITY_WEP_KEY1:
			case NM_SETTING_WIRELESS_SECURITY_WEP_KEY2:
			case NM_SETTING_WIRELESS_SECURITY_WEP_KEY3:
			case NM_SETTING_WIRELESS_SECURITY_WEP_KEY_FLAGS:
			case NM_SETTING_WIRELESS_SECURITY_WEP_KEY_TYPE:
			case NM_SETTING_WIRELESS_SECURITY_PSK:
			case NM_SETTING_WIRELESS_SECURITY_PSK_FLAGS:
			case NM_SETTING_WIRELESS_SECURITY_LEAP_PASSWORD:
			case NM_SETTING_WIRELESS_SECURITY_LEAP_PASSWORD_FLAGS:
			}
		}
	case typeWireless:
		// TODO
	}
	return
}
