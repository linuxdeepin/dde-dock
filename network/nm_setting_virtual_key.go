package network

import . "dlib/gettext"
import "fmt"

// TODO remove
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

// vpn
const NM_SETTING_VK_VPN_TYPE = "vk-vpn-type"

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

const (
	vkTypeWrapper       = "wrapper"
	vkTypeEnableWrapper = "enable-wrapper"

	// control other sections or keys, no related key, and the related
	// section always is a virtual section, such as "vk-vpn-type", for
	// there is no real related section, the key's name must be unique
	vkTypeController = "controller"
)

type vkeyInfo struct {
	Name           string
	Type           ktype
	VkType         string // could be "wrapper", "enable-wrapper", "control"
	RelatedSection string
	RelatedKeys    []string
	Available      bool // check if is used by front-end
	Optional       bool // if key is optional(such as child key gateway of ip address), will ignore error for it
}

func getVkeyInfo(section, vkey string) (info vkeyInfo, ok bool) {
	for _, vk := range virtualKeys {
		if vk.RelatedSection == section && vk.Name == vkey {
			info = vk
			ok = true
			return
		}
	}
	logger.Errorf("invalid virtual key, section=%s, vkey=%s", section, vkey)
	ok = false
	return
}

func isVirtualKey(section, key string) bool {
	// TODO
	if isStringInArray(key, getVkeysOfSection(section)) {
		return true
	}
	return false
}

// get all virtual keys in target section
func getVkeysOfSection(section string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.RelatedSection == section {
			vks = append(vks, vk.Name)
		}
	}
	// logger.Debug("getVkeysOfSection: filed:", section, vks) // TODO test
	return
}

func getSettingVkeyType(section, key string) (t ktype) {
	t = ktypeUnknown
	for _, vk := range virtualKeys {
		if vk.RelatedSection == section && vk.Name == key {
			t = vk.Type
		}
	}
	return
}

func generalGetSettingVkeyAvailableValues(data connectionData, section, key string) (values []kvalue) {
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
					kvalue{"none", Tr("None")},
					kvalue{"wep", Tr("WEP 40/128-bit Key")},
					kvalue{"wpa-psk", Tr("WPA & WPA2 Personal")},
					kvalue{"wpa-eap", Tr("WPA & WPA2 Enterprise")},
				}
			} else {
				values = []kvalue{
					kvalue{"none", Tr("None")},
					kvalue{"wep", Tr("WEP 40/128-bit Key")},
					kvalue{"wpa-psk", Tr("WPA & WPA2 Personal")},
				}
			}
		}
	case sectionPppoe:
	case sectionPpp:
	case sectionVpnL2tpPpp:
		switch key {
		case NM_SETTING_VK_VPN_L2TP_MPPE_SECURITY:
			values = []kvalue{
				kvalue{"default", Tr("All Available (default)")},
				kvalue{"128-bit", Tr("128-bit (most secure)")},
				kvalue{"40-bit", Tr("40-bit (less secure)")},
			}
		}
	case sectionVpnPptpPpp:
		switch key {
		case NM_SETTING_VK_VPN_PPTP_MPPE_SECURITY:
			values = []kvalue{
				kvalue{"default", Tr("All Available (default)")},
				kvalue{"128-bit", Tr("128-bit (most secure)")},
				kvalue{"40-bit", Tr("40-bit (less secure)")},
			}
		}
	case sectionVpnVpncAdvanced:
		switch key {
		case NM_SETTING_VK_VPN_VPNC_KEY_ENCRYPTION_METHOD:
			values = []kvalue{
				kvalue{"secure", Tr("Secure (default)")},
				kvalue{"weak", Tr("Weak")},
				kvalue{"none", Tr("None")},
			}
		}
	}

	if len(values) == 0 {
		logger.Warningf("there is no available values for virtual key, %s->%s", section, key)
	}
	return
}

// general function to append available keys, will dispatch virtual keys specially
func appendAvailableKeys(data connectionData, keys []string, section, key string) (newKeys []string) {
	newKeys = appendStrArrayUnique(keys)
	relatedVks := getRelatedAvailableVkeys(section, key)
	if len(relatedVks) > 0 {
		for _, vk := range relatedVks {
			// if is enable wrapper virtual key, both virtual key and
			// real key will be appended
			if isEnableWrapperVkey(section, vk) {
				if isSettingKeyExists(data, section, key) {
					newKeys = appendStrArrayUnique(newKeys, key)
				}
			}
		}
		newKeys = appendStrArrayUnique(newKeys, relatedVks...)
	} else {
		newKeys = appendStrArrayUnique(newKeys, key)
	}
	return
}

func getRelatedAvailableVkeys(section, key string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.RelatedSection == section && isStringInArray(key, vk.RelatedKeys) && vk.Available {
			vks = append(vks, vk.Name)
		}
	}
	return
}

// get related virtual keys of target key
func getRelatedVkeys(section, key string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.RelatedSection == section && isStringInArray(key, vk.RelatedKeys) {
			vks = append(vks, vk.Name)
		}
	}
	return
}

func isWrapperVkey(section, vkey string) bool {
	vkInfo, ok := getVkeyInfo(section, vkey)
	if !ok {
		return false
	}
	if vkInfo.VkType == vkTypeWrapper {
		return true
	}
	return false
}

func isEnableWrapperVkey(section, vkey string) bool {
	vkInfo, ok := getVkeyInfo(section, vkey)
	if !ok {
		return false
	}
	if vkInfo.VkType == vkTypeEnableWrapper {
		return true
	}
	return false
}

func isControlVkey(section, vkey string) bool {
	vkInfo, ok := getVkeyInfo(section, vkey)
	if !ok {
		return false
	}
	if vkInfo.VkType == vkTypeController {
		return true
	}
	return false
}

func isOptionalVkey(section, vkey string) (optional bool) {
	for _, vk := range virtualKeys {
		if vk.RelatedSection == section && vk.Name == vkey {
			optional = vk.Optional
		}
	}
	return
}

// Virtual key with none related section
func getSettingVkMobileServiceType(data connectionData) (serviceType string) {
	if isSettingSectionExists(data, NM_SETTING_GSM_SETTING_NAME) {
		serviceType = mobileServiceGsm
	} else if isSettingSectionExists(data, NM_SETTING_CDMA_SETTING_NAME) {
		serviceType = mobileServiceCdma
	} else {
		logger.Error("get mobile service type failed, neither gsm section nor cdma section")
	}
	return
}
func logicSetSettingVkMobileServiceType(data connectionData, serviceType string) (err error) {
	switch serviceType {
	case mobileServiceGsm:
		removeSettingSection(data, sectionCdma)
		initSettingSectionGsm(data)
	case mobileServiceCdma:
		removeSettingSection(data, sectionGsm)
		initSettingSectionCdma(data)
	default:
		err = fmt.Errorf("invalid mobile service type", serviceType)
	}
	return
}

func getSettingVkVpnType(data connectionData) (vpnType string) {
	vpnType = getCustomConnectionType(data)
	return
}
func logicSetSettingVkVpnType(data connectionData, vpnType string) (err error) {
	removeSettingSection(data, sectionVpn)
	switch vpnType {
	case connectionVpnL2tp:
		initSettingSectionVpnL2tp(data)
	case connectionVpnPptp:
		initSettingSectionVpnPptp(data)
	case connectionVpnOpenconnect:
		initSettingSectionVpnOpenconnect(data)
	case connectionVpnOpenvpn:
		initSettingSectionVpnOpenvpn(data)
	case connectionVpnVpnc:
		initSettingSectionVpnVpnc(data)
	}
	return
}
