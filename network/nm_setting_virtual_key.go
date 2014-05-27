package main

import (
	"dlib"
)

// If there is none related section for virtual key, it means that the
// virtual key used to control multiple sections, such as change
// connection type, and the key's name must be unique.
const NM_SETTING_VK_NONE_RELATED_FIELD = "<none>"

// For a virtual key with none related key, it is often used to
// control multiple keys in same section.
const NM_SETTING_VK_NONE_RELATED_KEY = "<none>"

// Virtual key names

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

// mobile
const NM_SETTING_VK_MOBILE_SERVICE_TYPE = "vk-mobile-service-type"

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
const NM_SETTING_VK_VPN_L2TP_REQUIRE_MPPE = "vk-require-mppe"
const NM_SETTING_VK_VPN_L2TP_MPPE_SECURITY = "vk-mppe-security"
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
const NM_SETTING_VK_VPN_PPTP_REQUIRE_MPPE = "vk-require-mppe"
const NM_SETTING_VK_VPN_PPTP_MPPE_SECURITY = "vk-mppe-security"
const NM_SETTING_VK_VPN_PPTP_ENABLE_LCP_ECHO = "vk-enable-lcp-echo"

// vpn-vpnc
const (
	NM_SETTING_VK_VPN_VPNC_KEY_HYBRID_AUTHMODE   = "vk-hybrid-authmode"
	NM_SETTING_VK_VPN_VPNC_KEY_ENCRYPTION_METHOD = "vk-encryption-method"
	NM_SETTING_VK_VPN_VPNC_KEY_DISABLE_DPD       = "vk-disable-dpd"
)

// wireless security
const NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT = "vk-key-mgmt"

// virtualKey stores virtual key info for each sections.
type virtualKey struct {
	Name          string
	Type          ktype
	RelatedSection  string
	RelatedKey    string
	EnableWrapper bool // check if the virtual key is a wrapper just to enable target key
	Available     bool // check if is used by front-end
	Optional      bool // if key is optional, will ignore error for it
}

func getVirtualKeyInfo(section, vkey string) (vkInfo virtualKey, ok bool) {
	for _, vk := range virtualKeys {
		if vk.RelatedSection == section && vk.Name == vkey {
			vkInfo = vk
			ok = true
			return
		}
	}
	logger.Errorf("invalid virtual key, section=%s, vkey=%s", section, vkey)
	ok = false
	return
}

func isVirtualKey(section, key string) bool {
	if isStringInArray(key, getVirtualKeysOfSection(section)) {
		return true
	}
	return false
}

func getVirtualKeysOfSection(section string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.RelatedSection == section {
			vks = append(vks, vk.Name)
		}
	}
	// logger.Debug("getVirtualKeysOfSection: filed:", section, vks) // TODO test
	return
}

func getSettingVkKeyType(section, key string) (t ktype) {
	t = ktypeUnknown
	for _, vk := range virtualKeys {
		if vk.RelatedSection == section && vk.Name == key {
			t = vk.Type
		}
	}
	return
}

func generalGetSettingVkAvailableValues(data connectionData, section, key string) (values []kvalue) {
	switch section {
	case section8021x:
		switch key {
		case NM_SETTING_VK_802_1X_EAP:
			values = getSetting8021xAvailableValues(data, NM_SETTING_802_1X_EAP)
		}
	case sectionConnection:
	case sectionIpv4:
	case sectionIpv6:
	case sectionWired:
	case sectionWireless:
	case sectionWirelessSecurity:
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
	case sectionPppoe:
	case sectionPpp:
	case sectionVpnL2tpPpp:
		switch key {
		case NM_SETTING_VK_VPN_L2TP_MPPE_SECURITY:
			values = []kvalue{
				kvalue{"default", dlib.Tr("All Available (default)")},
				kvalue{"128-bit", dlib.Tr("128-bit (most secure)")},
				kvalue{"40-bit", dlib.Tr("40-bit (less secure)")},
			}
		}
	case sectionVpnPptpPpp:
		switch key {
		case NM_SETTING_VK_VPN_PPTP_MPPE_SECURITY:
			values = []kvalue{
				kvalue{"default", dlib.Tr("All Available (default)")},
				kvalue{"128-bit", dlib.Tr("128-bit (most secure)")},
				kvalue{"40-bit", dlib.Tr("40-bit (less secure)")},
			}
		}
	case sectionVpnVpncAdvanced:
		switch key {
		case NM_SETTING_VK_VPN_VPNC_KEY_ENCRYPTION_METHOD:
			values = []kvalue{
				kvalue{"secure", dlib.Tr("Secure (default)")},
				kvalue{"weak", dlib.Tr("Weak")},
				kvalue{"none", dlib.Tr("None")},
			}
		}
	}

	// dispatch virtual keys that with none related sections
	if len(values) == 0 {
		switch key {
		case NM_SETTING_VK_MOBILE_SERVICE_TYPE:
			values = []kvalue{
				kvalue{"gsm", dlib.Tr("GSM (GPRS, EDGE, UMTS, HSPA)")},
				kvalue{"cdma", dlib.Tr("CDMA (1xRTT, EVDO)")},
			}
		}
	}

	if len(values) == 0 {
		logger.Warningf("there is no available values for virtual key, %s->%s", section, key)
	}
	return
}

func appendAvailableKeys(data connectionData, keys []string, section, key string) (newKeys []string) {
	newKeys = appendStrArrayUnion(keys)
	relatedVks := getRelatedAvailableVirtualKeys(section, key)
	if len(relatedVks) > 0 {
		for _, vk := range relatedVks {
			// if is enable wrapper virtual key, both virtual key and
			// real key will be appended
			if isEnableWrapperVirtualKey(section, vk) {
				if isSettingKeyExists(data, section, key) {
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

func getRelatedAvailableVirtualKeys(section, key string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.RelatedSection == section && vk.RelatedKey == key && vk.Available {
			vks = append(vks, vk.Name)
		}
	}
	return
}

// get related virtual key(s) for target key
func getRelatedVirtualKeys(section, key string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.RelatedSection == section && vk.RelatedKey == key {
			vks = append(vks, vk.Name)
		}
	}
	return
}

func isEnableWrapperVirtualKey(section, vkey string) bool {
	vkInfo, ok := getVirtualKeyInfo(section, vkey)
	if !ok {
		return false
	}
	return vkInfo.EnableWrapper
}

func isOptionalChildVirtualKeys(section, vkey string) (optional bool) {
	for _, vk := range virtualKeys {
		if vk.RelatedSection == section && vk.Name == vkey {
			optional = vk.Optional
		}
	}
	return
}
