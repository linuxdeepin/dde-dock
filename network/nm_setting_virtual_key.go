package main

import (
	"dlib"
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

// wired
const NM_SETTING_VK_WIRED_ENABLE_MTU = "vk-enable-mtu"

// wireless
const NM_SETTING_VK_WIRELESS_ENABLE_MTU = "vk-enable-mtu"

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
	NM_SETTING_VK_VPN_OPENVPN_KEY_ENABLE_PORT                 = "vk-enable-port"
	NM_SETTING_VK_VPN_OPENVPN_KEY_ENABLE_RENEG_SECONDS        = "vk-enable-reneg-seconds"
	NM_SETTING_VK_VPN_OPENVPN_KEY_ENABLE_TUNNEL_MTU           = "vk-enable-tunnel-mtu"
	NM_SETTING_VK_VPN_OPENVPN_KEY_ENABLE_FRAGMENT_SIZE        = "vk-enable-fragment-size"
	NM_SETTING_VK_VPN_OPENVPN_KEY_ENABLE_STATIC_KEY_DIRECTION = "vk-static-key-direction"
	NM_SETTING_VK_VPN_OPENVPN_KEY_ENABLE_TA_DIR               = "vk-ta-dir"
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

func generalGetSettingVkAvailableValues(data connectionData, field, key string) (values []kvalue) {
	switch field {
	case field8021x:
		switch key {
		case NM_SETTING_VK_802_1X_EAP:
			values = getSetting8021xAvailableValues(data, NM_SETTING_802_1X_EAP)
		}
	case fieldConnection:
	case fieldIpv4:
	case fieldIpv6:
	case fieldWired:
	case fieldWireless:
	case fieldWirelessSecurity:
		switch key {
		case NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT:
			if getSettingWirelessMode(data) == NM_SETTING_WIRELESS_MODE_INFRA {
				values = []kvalue{
					kvalue{"none", dlib.Tr("None")},
					kvalue{"wep", dlib.Tr("WEP 40/128-bit Key")},
					kvalue{"wpa-psk", dlib.Tr("WPA & WPA2 Personal")},
					kvalue{"wpa-eap", dlib.Tr("WPA & WPA2 Enterprise")},
				}
			} else {
				values = []kvalue{
					kvalue{"none", dlib.Tr("None")},
					kvalue{"wep", dlib.Tr("WEP 40/128-bit Key")},
					kvalue{"wpa-psk", dlib.Tr("WPA & WPA2 Personal")},
				}
			}
		}
	case fieldPppoe:
	case fieldPpp:
	case fieldVpnVpncAdvanced:
		switch key {
		case NM_SETTING_VK_VPN_VPNC_KEY_ENCRYPTION_METHOD:
			values = []kvalue{
				kvalue{"secure", dlib.Tr("Secure (default)")},
				kvalue{"weak", dlib.Tr("Weak")},
				kvalue{"none", dlib.Tr("None")},
			}
		}
	}
	if len(values) == 0 {
		logger.Warningf("there is no available values for virtual key, %s->%s", field, key)
	}
	return
}

func appendAvailableKeys(data connectionData, keys []string, field, key string) (newKeys []string) {
	newKeys = appendStrArrayUnion(keys)
	relatedVks := getRelatedAvailableVirtualKeys(field, key)
	if len(relatedVks) > 0 {
		for _, vk := range relatedVks {
			// if is enable wrapper virtual key, both virtual key and
			// real key will be appended
			if isEnableWrapperVirtualKey(field, vk) {
				if isSettingKeyExists(data, field, key) {
					newKeys = appendStrArrayUnion(newKeys, key)
				}
			}
		}
		newKeys = appendStrArrayUnion(newKeys, relatedVks...)
	} else {
		newKeys = appendStrArrayUnion(newKeys, key)
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
