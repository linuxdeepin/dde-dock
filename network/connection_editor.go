package main

import (
	"fmt"
)

const (
	pageGeneral = "General"
	pageIPv4    = "IPv4"
	pageIPv6    = "IPv6"
	page8021x   = "802.1xSecurity"
)

var supportedConnectionTypes = []string{
	NM_SETTING_WIRED_SETTING_NAME,
	NM_SETTING_WIRELESS_SETTING_NAME,
}

type ConnectionEditor struct {
	data _ConnectionData

	connType string // TODO

	CurrentUUID string
	HasChanged  bool

	currentPage string

	//前端只显示此列表中的字段,会跟随当前正在编辑的值而改变
	// TODO more documentation
	CurrentFields []string
	//返回当前page下错误的字段和对应的错误原因
	CurrentErrors []string
}

func NewConnectionEditor() *ConnectionEditor {
	editor := &ConnectionEditor{}
	return editor
}

//所有字段值都为string，后端自行转换为需要的值后提供给NM

//新建一个Connection 返回uuid (此时这个Connection还未提交到NM)
//如果是支持的类型则设置CurrentUUID属性
// Create try to create a new connection, return empty string if error ocurred.
func (editor *ConnectionEditor) Create(connType string) (uuid string) {
	if !isStringInArray(connType, supportedConnectionTypes) {
		return ""
	}

	// TODO
	editor.data = make(_ConnectionData)
	editor.connType = connType

	uuid = newUUID()
	editor.updatePropCurrentUUID(uuid)
	editor.updatePropHasChanged(true)

	// TODO
	// editor.currentPage = editor.getDefaultPage(connType)
	// editor.updatePropCurrentFields()

	// TODO current errors

	return uuid
}

// // TODO get default page of target connection type
// func (editor *ConnectionEditor) getDefaultPage(connType string) (defpage string) {
// 	switch connType {
// 	case NM_SETTING_WIRED_SETTING_NAME:
// 		defpage = pageIPv4 // TODO
// 	case NM_SETTING_WIRELESS_SETTING_NAME:
// 		defpage = "default page" // TODO
// 	}
// 	return
// }

//保存当前Connection的修改。  不错任何处理如果Changed属性=false
func (editor *ConnectionEditor) Save() {
	// TODO
	if !editor.HasChanged {
		return
	}
}

//打开uuid指定的Connection 如果无法通过
//org.freedesktop.NetworkManager.Settings的GetConnectionByUuid得到结果
//则设置Error属性如果成功打开则设置CurrentUUID属性
func (editor *ConnectionEditor) Open(uuid string) (err error) {
	// TODO
	return
}

//根据CurrentUUID返回此Connection支持的设置页面
func (editor *ConnectionEditor) ListPages() (pages []string) {
	// TODO
	switch editor.connType {
	case NM_SETTING_WIRED_SETTING_NAME:
		pages = []string{
			pageGeneral,
			pageIPv4,
			pageIPv6,
			page8021x,
		}
	case NM_SETTING_WIRELESS_SETTING_NAME:
	}
	return
}

// get valid fields for target page
func (editor *ConnectionEditor) listFields(page string) (fields []string) {
	switch editor.connType {
	case NM_SETTING_WIRED_SETTING_NAME:
		switch editor.currentPage {
		case pageGeneral:
			fields = []string{"General_field1", "General_field2"}
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
			fields = []string{"IPv6_field1", "IPv6_field2"}
		case page8021x:
			fields = []string{"802.1xSecurity_field1", "802.1xSecurity_field2"}
		}
	case NM_SETTING_WIRELESS_SETTING_NAME:
	}
	return
}

//设置/获得字段的值都受这里设置page的影响。
func (editor *ConnectionEditor) SwitchPage(page string) {
	// TODO HasChanged
	editor.currentPage = page
	editor.updatePropCurrentFields()
}

//比如获得当前链接支持的加密方式 EAP字段: TLS、MD5、FAST、PEAP
//获得ip设置方式 : Manual、Link-Local Only、Automatic(DHCP)
//获得当前可用mac地址(这种字段是有几个可选值但用户也可用手动输入一个其他值)
func (editor *ConnectionEditor) GetAvailableValue(key string) (values []string) {
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
func (editor *ConnectionEditor) DebugListFields() []string {
	// TODO
	return editor.listFields(editor.currentPage)
}

// DebugConnectionTypes return all supported connection types, only for debugging.
func (editor *ConnectionEditor) DebugConnectionTypes() []string {
	return supportedConnectionTypes
}

//设置某个字段， 会影响CurrentFields属性，某些值会导致其他属性进入不可用状态
func (editor *ConnectionEditor) SetField(key, value string) (err error) {
	page, ok := editor.data[editor.currentPage]
	if !ok {
		err = fmt.Errorf("invalid page: data[%s]", editor.currentPage)
		return
	}

	varient, ok := page[key]
	if !ok {
		// field is not yet exits
		page[key] = wrapVarient(value)
		editor.HasChanged = true
		return
	}

	oldValue, err := unwrapVarient(varient)
	if err != nil {
		return
	}
	if oldValue == value {
		return
	}
	page[key] = wrapVarient(value)
	editor.HasChanged = true

	// TODO processing logic
	editor.updatePropCurrentFields()

	return
}

func (editor *ConnectionEditor) GetField(key string) (value string, err error) {
	page, ok := editor.data[editor.currentPage]
	if !ok {
		err = fmt.Errorf("invalid page: data[%s]", editor.currentPage)
		return
	}

	varient, ok := page[key]
	if !ok {
		err = fmt.Errorf("invalid field: data[%s][%s]", editor.currentPage, key)
		return
	}

	value, err = unwrapVarient(varient)
	return
}
