package main

import (
	"os/user"
)

// Virtual key names.

const NM_SETTING_VK_NONE_RELATED_KEY = "<none>"

// 802-1x
const (
	NM_SETTING_VK_802_1X_ENABLE      = "vk-enable"
	NM_SETTING_VK_802_1X_EAP         = "vk-eap"
	NM_SETTING_VK_802_1X_PAC_FILE    = "vk-pac-file"
	NM_SETTING_VK_802_1X_CA_CERT     = "vk-ca-cert"
	NM_SETTING_VK_802_1X_CLIENT_CERT = "vk-client-cert"
	NM_SETTING_VK_802_1X_PRIVATE_KEY = "vk-private-key"
)

// connection
const NM_SETTING_VK_CONNECTION_NO_PERMISSION = "vk-no-permission"

// ipv4
const (
	NM_SETTING_VK_IP4_CONFIG_DNS               = "vk-dns"
	NM_SETTING_VK_IP4_CONFIG_ADDRESSES_ADDRESS = "vk-addresses-address"
	NM_SETTING_VK_IP4_CONFIG_ADDRESSES_MASK    = "vk-addresses-mask"
	NM_SETTING_VK_IP4_CONFIG_ADDRESSES_GATEWAY = "vk-addresses-gateway"
	NM_SETTING_VK_IP4_CONFIG_ROUTES_ADDRESS    = "vk-routes-address"
	NM_SETTING_VK_IP4_CONFIG_ROUTES_MASK       = "vk-routes-mask"
	NM_SETTING_VK_IP4_CONFIG_ROUTES_NEXTHOP    = "vk-routes-nexthop"
	NM_SETTING_VK_IP4_CONFIG_ROUTES_METRIC     = "vk-routes-metric"
)

// ipv6
const (
	NM_SETTING_VK_IP6_CONFIG_DNS               = "vk-dns"
	NM_SETTING_VK_IP6_CONFIG_ADDRESSES_ADDRESS = "vk-addresses-address"
	NM_SETTING_VK_IP6_CONFIG_ADDRESSES_PREFIX  = "vk-addresses-prefix"
	NM_SETTING_VK_IP6_CONFIG_ADDRESSES_GATEWAY = "vk-addresses-gateway"
	NM_SETTING_VK_IP6_CONFIG_ROUTES_ADDRESS    = "vk-routes-address"
	NM_SETTING_VK_IP6_CONFIG_ROUTES_PREFIX     = "vk-routes-prefix"
	NM_SETTING_VK_IP6_CONFIG_ROUTES_NEXTHOP    = "vk-routes-nexthop"
	NM_SETTING_VK_IP6_CONFIG_ROUTES_METRIC     = "vk-routes-metric"
)

// ppp
const NM_SETTING_VK_PPP_ENABLE_LCP_ECHO = "vk-enable-lcp-echo"

// vpn-l2tp
const NM_SETTING_VK_VPN_L2TP_ENABLE_LCP_ECHO = "vk-enable-lcp-echo"

// vpn-openvpn
const (
	NM_SETTING_VK_VPN_OPENVPN_KEY_ENABLE_PORT          = "vk-enable-port"
	NM_SETTING_VK_VPN_OPENVPN_KEY_ENABLE_RENEG_SECONDS = "vk-enable-reneg-seconds"
	NM_SETTING_VK_VPN_OPENVPN_KEY_ENABLE_TUNNEL_MTU    = "vk-enable-tunnel-mtu"
	NM_SETTING_VK_VPN_OPENVPN_KEY_ENABLE_FRAGMENT_SIZE = "vk-enable-fragment-size"
)

// vpn-pptp
const NM_SETTING_VK_VPN_PPTP_ENABLE_LCP_ECHO = "vk-enable-lcp-echo"

// vpn-vpnc
const (
	NM_SETTING_VK_VPN_VPNC_KEY_HYBRID_AUTHMODE   = "vk-hybrid-authmode"
	NM_SETTING_VK_VPN_VPNC_KEY_ENCRYPTION_METHOD = "vk-encryption-method"
	NM_SETTING_VK_VPN_VPNC_KEY_DISABLE_DPD       = "vk-disable-dpd"
)

// wireless security
const NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT = "vk-key-mgmt"

// virtualKey stores virtual key info for each fields.
type virtualKey struct {
	Name          string
	Type          ktype
	RelatedField  string
	RelatedKey    string
	EnableWrapper bool // check if the virtual key is a wrapper just to enable target key
	Available     bool // check if is used by front-end
	Optional      bool // if key is optional, will ignore error for it
}

func getVirtualKeyInfo(field, vkey string) (vkInfo virtualKey, ok bool) {
	for _, vk := range virtualKeys {
		if vk.RelatedField == field && vk.Name == vkey {
			vkInfo = vk
			ok = true
			return
		}
	}
	logger.Errorf("invalid virtual key, field=%s, vkey=%s", field, vkey)
	ok = false
	return
}

func isVirtualKey(field, key string) bool {
	if isStringInArray(key, getVirtualKeysOfField(field)) {
		return true
	}
	return false
}

func getVirtualKeysOfField(field string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.RelatedField == field {
			vks = append(vks, vk.Name)
		}
	}
	// logger.Debug("getVirtualKeysOfField: filed:", field, vks) // TODO test
	return
}

func getSettingVkKeyType(field, key string) (t ktype) {
	t = ktypeUnknown
	for _, vk := range virtualKeys {
		if vk.RelatedField == field && vk.Name == key {
			t = vk.Type
		}
	}
	return
}

func generalGetSettingVkAvailableValues(data connectionData, field, key string) (values []string) {
	switch field {
	case field8021x:
		switch key {
		case NM_SETTING_VK_802_1X_EAP:
			values, _ = getSetting8021xAvailableValues(data, NM_SETTING_802_1X_EAP)
		}
	case fieldConnection:
	case fieldIpv4:
	case fieldIpv6:
	case fieldWired:
	case fieldWireless:
	case fieldWirelessSecurity:
		switch key {
		case NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT:
			values = []string{"none", "wep", "wpa-psk", "wpa-eap"}
		}
	case fieldPppoe:
	case fieldPpp:
	case fieldVpnVpncAdvanced:
		switch key {
		case NM_SETTING_VK_VPN_VPNC_KEY_ENCRYPTION_METHOD:
			values = []string{"secure", "weak", "none"}
		}
	}
	if len(values) == 0 {
		logger.Warningf("there is no available values for virtual key, %s->%s", field, key)
	}
	return
}

func appendAvailableKeys(keys []string, field, key string) (appendKeys []string) {
	relatedVks := getRelatedAvailableVirtualKeys(field, key)
	if len(relatedVks) > 0 {
		appendKeys = appendStrArrayUnion(keys, relatedVks...)
	} else {
		appendKeys = appendStrArrayUnion(keys, key)
	}
	return
}

func getRelatedAvailableVirtualKeys(field, key string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.RelatedField == field && vk.RelatedKey == key && vk.Available {
			vks = append(vks, vk.Name)
		}
	}
	return
}

// get related virtual key(s) for target key
func getRelatedVirtualKeys(field, key string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.RelatedField == field && vk.RelatedKey == key {
			vks = append(vks, vk.Name)
		}
	}
	return
}

func isEnableWrapperVirtualKey(field, vkey string) bool {
	vkInfo, ok := getVirtualKeyInfo(field, vkey)
	if !ok {
		return false
	}
	return vkInfo.EnableWrapper
}

func isOptionalChildVirtualKeys(field, vkey string) (optional bool) {
	for _, vk := range virtualKeys {
		if vk.RelatedField == field && vk.Name == vkey {
			optional = vk.Optional
		}
	}
	return
}

func doGetOrNewSettingIp4ConfigAddresses(data connectionData) (addresses [][]uint32) {
	addresses = getSettingIp4ConfigAddresses(data)
	if len(addresses) == 0 {
		addresses = make([][]uint32, 1)
		addresses[0] = make([]uint32, 3)
	}
	if len(addresses[0]) != 3 {
		addresses[0] = make([]uint32, 3)
	}
	return
}

func doIsSettingIp4ConfigAddressesEmpty(data connectionData) bool {
	addresses := getSettingIp4ConfigAddresses(data)
	if len(addresses) == 0 {
		return true
	}
	if len(addresses[0]) != 3 {
		return true
	}
	return false
}

// Getter
func getSettingVk8021xEnable(data connectionData) (value bool) {
	if isSettingFieldExists(data, field8021x) {
		return true
	}
	return false
}
func getSettingVk8021xEap(data connectionData) (eap string) {
	eaps := getSetting8021xEap(data)
	if len(eaps) == 0 {
		logger.Error("eap value is empty")
		return
	}
	eap = eaps[0]
	return
}
func getSettingVk8021xPacFile(data connectionData) (value string) {
	pacFile := getSetting8021xPacFile(data)
	if len(pacFile) > 0 {
		value = toUriPath(pacFile)
	}
	return
}
func getSettingVk8021xCaCert(data connectionData) (value string) {
	caCert := getSetting8021xCaCert(data)
	value = byteArrayToStrPath(caCert)
	return
}
func getSettingVk8021xClientCert(data connectionData) (value string) {
	clientCert := getSetting8021xClientCert(data)
	value = byteArrayToStrPath(clientCert)
	return
}
func getSettingVk8021xPrivateKey(data connectionData) (value string) {
	privateKey := getSetting8021xPrivateKey(data)
	value = byteArrayToStrPath(privateKey)
	return
}
func getSettingVkConnectionNoPermission(data connectionData) (value bool) {
	permission := getSettingConnectionPermissions(data)
	if len(permission) > 0 {
		return false
	}
	return true
}
func getSettingVkIp4ConfigDns(data connectionData) (value string) {
	dns := getSettingIp4ConfigDns(data)
	if len(dns) == 0 {
		return
	}
	value = convertIpv4AddressToString(dns[0])
	return
}
func getSettingVkIp4ConfigAddressesAddress(data connectionData) (value string) {
	if doIsSettingIp4ConfigAddressesEmpty(data) {
		return
	}
	addresses := getSettingIp4ConfigAddresses(data)
	value = convertIpv4AddressToString(addresses[0][0])
	return
}
func getSettingVkIp4ConfigAddressesMask(data connectionData) (value string) {
	if doIsSettingIp4ConfigAddressesEmpty(data) {
		return
	}
	addresses := getSettingIp4ConfigAddresses(data)
	value = convertIpv4PrefixToNetMask(addresses[0][1])
	return
}
func getSettingVkIp4ConfigAddressesGateway(data connectionData) (value string) {
	if doIsSettingIp4ConfigAddressesEmpty(data) {
		return
	}
	addresses := getSettingIp4ConfigAddresses(data)
	value = convertIpv4AddressToStringNoZero(addresses[0][2])
	return
}
func getSettingVkIp4ConfigRoutesAddress(data connectionData) (value string) {
	// TODO
	// value := getSettingIp4ConfigRoutesAddress(data)
	return
}
func getSettingVkIp4ConfigRoutesMask(data connectionData) (value string) {
	// TODO
	// value := getSettingIp4ConfigRoutesMask(data)
	return
}
func getSettingVkIp4ConfigRoutesNexthop(data connectionData) (value string) {
	// TODO
	// value := getSettingIp4ConfigRoutesNexthop(data)
	return
}
func getSettingVkIp4ConfigRoutesMetric(data connectionData) (value string) {
	// TODO
	// value := getSettingIp4ConfigRoutesMetric(data)
	return
}
func getSettingVkIp6ConfigDns(data connectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigDns(data)
	return
}
func getSettingVkIp6ConfigAddressesAddress(data connectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigAddressesAddress(data)
	return
}
func getSettingVkIp6ConfigAddressesPrefix(data connectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigAddressesPrefix(data)
	return
}
func getSettingVkIp6ConfigAddressesGateway(data connectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigAddressesGateway(data)
	return
}
func getSettingVkIp6ConfigRoutesAddress(data connectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigRoutesAddress(data)
	return
}
func getSettingVkIp6ConfigRoutesPrefix(data connectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigRoutesPrefix(data)
	return
}
func getSettingVkIp6ConfigRoutesNexthop(data connectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigRoutesNexthop(data)
	return
}
func getSettingVkIp6ConfigRoutesMetric(data connectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigRoutesMetric(data)
	return
}
func getSettingVkPppEnableLcpEcho(data connectionData) (value bool) {
	if isSettingPppLcpEchoFailureExists(data) && isSettingPppLcpEchoIntervalExists(data) {
		return true
	}
	return false
}
func getSettingVkVpnL2tpEnableLcpEcho(data connectionData) (value bool) {
	if isSettingVpnL2tpKeyLcpEchoFailureExists(data) && isSettingVpnL2tpKeyLcpEchoIntervalExists(data) {
		return true
	}
	return false
}
func getSettingVkVpnPptpEnableLcpEcho(data connectionData) (value bool) {
	if isSettingVpnPptpKeyLcpEchoFailureExists(data) && isSettingVpnPptpKeyLcpEchoIntervalExists(data) {
		return true
	}
	return false
}
func getSettingVkVpnVpncKeyHybridAuthmode(data connectionData) (value bool) {
	if isSettingVpnVpncKeyAuthmodeExists(data) {
		return true
	}
	return false
}
func getSettingVkVpnVpncKeyEncryptionMethod(data connectionData) (value string) {
	if getSettingVpnVpncKeySingleDes(data) {
		return "weak"
	} else if getSettingVpnVpncKeyNoEncryption(data) {
		return "none"
	}
	return "secure"
}
func getSettingVkVpnVpncKeyDisableDpd(data connectionData) (value bool) {
	if isSettingVpnVpncKeyDpdIdleTimeoutExists(data) && getSettingVpnVpncKeyDpdIdleTimeout(data) == 0 {
		return true
	}
	return false
}
func getSettingVkWirelessSecurityKeyMgmt(data connectionData) (value string) {
	if !isSettingFieldExists(data, fieldWirelessSecurity) {
		value = "none"
		return
	}
	keyMgmt := getSettingWirelessSecurityKeyMgmt(data)
	switch keyMgmt {
	case "none":
		value = "wep"
	case "wpa-psk":
		value = "wpa-psk"
	case "wpa-eap":
		value = "wpa-eap"
	}
	return
}

// Logic setter, all virtual keys has a logic setter
func logicSetSettingVk8021xEnable(data connectionData, value bool) {
	if value {
		addSettingField(data, field8021x)
		logicSetSettingVk8021xEap(data, "tls")
	} else {
		removeSettingField(data, field8021x)
	}
}
func logicSetSettingVk8021xEap(data connectionData, value string) {
	logicSetSetting8021xEap(data, []string{value})
}
func logicSetSettingVk8021xPacFile(data connectionData, value string) {
	setSetting8021xPacFile(data, toLocalPath(value))
}
func logicSetSettingVk8021xCaCert(data connectionData, value string) {
	setSetting8021xCaCert(data, strToByteArrayPath(value))
}
func logicSetSettingVk8021xClientCert(data connectionData, value string) {
	setSetting8021xClientCert(data, strToByteArrayPath(value))
}
func logicSetSettingVk8021xPrivateKey(data connectionData, value string) {
	setSetting8021xPrivateKey(data, strToByteArrayPath(value))
}
func logicSetSettingVkConnectionNoPermission(data connectionData, value bool) {
	if value {
		removeSettingConnectionPermissions(data)
	} else {
		currentUser, err := user.Current()
		if err != nil {
			logger.Error(err)
			return
		}
		permission := "user:" + currentUser.Username + ":"
		setSettingConnectionPermissions(data, []string{permission})
	}
}
func logicSetSettingVkIp4ConfigDns(data connectionData, value string) {
	dns := getSettingIp4ConfigDns(data)
	if len(dns) == 0 {
		dns = make([]uint32, 1)
	}
	dns[0] = convertIpv4AddressToUint32(value)
	if dns[0] != 0 {
		setSettingIp4ConfigDns(data, dns)
	} else {
		removeSettingIp4ConfigDns(data)
	}
}
func logicSetSettingVkIp4ConfigAddressesAddress(data connectionData, value string) {
	addresses := doGetOrNewSettingIp4ConfigAddresses(data)
	addr := addresses[0]
	addr[0] = convertIpv4AddressToUint32(value)
	if !isUint32ArrayEmpty(addr) {
		setSettingIp4ConfigAddresses(data, addresses)
	} else {
		removeSettingIp4ConfigAddresses(data)
	}
}
func logicSetSettingVkIp4ConfigAddressesMask(data connectionData, value string) {
	addresses := doGetOrNewSettingIp4ConfigAddresses(data)
	addr := addresses[0]
	addr[1] = convertIpv4NetMaskToPrefix(value)
	if !isUint32ArrayEmpty(addr) {
		setSettingIp4ConfigAddresses(data, addresses)
	} else {
		removeSettingIp4ConfigAddresses(data)
	}
}
func logicSetSettingVkIp4ConfigAddressesGateway(data connectionData, value string) {
	if len(value) == 0 {
		value = ipv4Zero
	}
	addresses := doGetOrNewSettingIp4ConfigAddresses(data)
	addr := addresses[0]
	addr[2] = convertIpv4AddressToUint32(value)
	if !isUint32ArrayEmpty(addr) {
		setSettingIp4ConfigAddresses(data, addresses)
	} else {
		removeSettingIp4ConfigAddresses(data)
	}
}
func logicSetSettingVkIp4ConfigRoutesAddress(data connectionData, value string) {
	// TODO
	// setSettingIp4ConfigRoutesAddressJSON(data)
}
func logicSetSettingVkIp4ConfigRoutesMask(data connectionData, value string) {
	// TODO
	// setSettingIp4ConfigRoutesMaskJSON(data)
}
func logicSetSettingVkIp4ConfigRoutesNexthop(data connectionData, value string) {
	// TODO
	// setSettingIp4ConfigRoutesNexthopJSON(data)
}
func logicSetSettingVkIp4ConfigRoutesMetric(data connectionData, value string) {
	// TODO
	// setSettingIp4ConfigRoutesMetricJSON(data)
}
func logicSetSettingVkIp6ConfigDns(data connectionData, value string) {
	// TODO
	// setSettingIp6ConfigDnsJSON(data)
}
func logicSetSettingVkIp6ConfigAddressesAddress(data connectionData, value string) {
	// TODO
	// setSettingIp6ConfigAddressesAddressJSON(data)
}
func logicSetSettingVkIp6ConfigAddressesPrefix(data connectionData, value uint32) {
	// TODO
	// setSettingIp6ConfigAddressesPrefixJSON(data)
}
func logicSetSettingVkIp6ConfigAddressesGateway(data connectionData, value string) {
	// TODO
	// setSettingIp6ConfigAddressesGatewayJSON(data)
}
func logicSetSettingVkIp6ConfigRoutesAddress(data connectionData, value string) {
	// TODO
	// setSettingIp6ConfigRoutesAddressJSON(data)
}
func logicSetSettingVkIp6ConfigRoutesPrefix(data connectionData, value uint32) {
	// TODO
	// setSettingIp6ConfigRoutesPrefixJSON(data)
}
func logicSetSettingVkIp6ConfigRoutesNexthop(data connectionData, value string) {
	// TODO
	// setSettingIp6ConfigRoutesNexthopJSON(data)
}
func logicSetSettingVkIp6ConfigRoutesMetric(data connectionData, value uint32) {
	// TODO
	// setSettingIp6ConfigRoutesMetricJSON(data)
}
func logicSetSettingVkPppEnableLcpEcho(data connectionData, value bool) {
	if value {
		setSettingPppLcpEchoFailure(data, 5)
		setSettingPppLcpEchoInterval(data, 30)
	} else {
		removeSettingPppLcpEchoFailure(data)
		removeSettingPppLcpEchoInterval(data)
	}
}
func logicSetSettingVkVpnL2tpEnableLcpEcho(data connectionData, value bool) {
	if value {
		setSettingVpnL2tpKeyLcpEchoFailure(data, 5)
		setSettingVpnL2tpKeyLcpEchoInterval(data, 30)
	} else {
		removeSettingVpnL2tpKeyLcpEchoFailure(data)
		removeSettingVpnL2tpKeyLcpEchoInterval(data)
	}
}
func logicSetSettingVkVpnPptpEnableLcpEcho(data connectionData, value bool) {
	if value {
		setSettingVpnPptpKeyLcpEchoFailure(data, 5)
		setSettingVpnPptpKeyLcpEchoInterval(data, 30)
	} else {
		removeSettingVpnPptpKeyLcpEchoFailure(data)
		removeSettingVpnPptpKeyLcpEchoInterval(data)
	}
}
func logicSetSettingVkVpnVpncKeyHybridAuthmode(data connectionData, value bool) {
	if value {
		setSettingVpnVpncKeyAuthmode(data, "hybrid")
	} else {
		removeSettingVpnVpncKeyAuthmode(data)
	}
}
func logicSetSettingVkVpnVpncKeyEncryptionMethod(data connectionData, value string) {
	removeSettingVpnVpncKeySingleDes(data)
	removeSettingVpnVpncKeyNoEncryption(data)
	switch value {
	case "secure":
	case "weak":
		setSettingVpnVpncKeySingleDes(data, true)
	case "none":
		setSettingVpnVpncKeyNoEncryption(data, true)
	}
}
func logicSetSettingVkVpnVpncKeyDisableDpd(data connectionData, value bool) {
	if value {
		setSettingVpnVpncKeyDpdIdleTimeout(data, 0)
	} else {
		removeSettingVpnVpncKeyDpdIdleTimeout(data)
	}
}
func logicSetSettingVkWirelessSecurityKeyMgmt(data connectionData, value string) {
	switch value {
	default:
		logger.Error("invalid value", value)
	case "none":
		removeSettingField(data, fieldWirelessSecurity)
		removeSettingField(data, field8021x)
	case "wep":
		addSettingField(data, fieldWirelessSecurity)
		removeSettingField(data, field8021x)
		setSettingWirelessSecurityKeyMgmt(data, "none")
		setSettingWirelessSecurityAuthAlg(data, "open")
		setSettingWirelessSecurityWepKeyType(data, 1)
	case "wpa-psk":
		addSettingField(data, fieldWirelessSecurity)
		removeSettingField(data, field8021x)
		setSettingWirelessSecurityKeyMgmt(data, "wpa-psk")
	case "wpa-eap":
		addSettingField(data, fieldWirelessSecurity)
		addSettingField(data, field8021x)
		setSettingWirelessSecurityKeyMgmt(data, "wpa-eap")
		logicSetSetting8021xEap(data, []string{"tls"})
		// TODO
	}
}
