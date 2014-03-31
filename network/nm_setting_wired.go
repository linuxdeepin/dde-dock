package main

// https://developer.gnome.org/libnm-util/0.9/NMSettingWired.html
// https://developer.gnome.org/NetworkManager/unstable/ref-settings.html

// The setting's name; these names are defined by the specification
// and cannot be changed after the object has been created. Each
// setting class has a name, and all objects of that class share the
// same name.
const NM_SETTING_WIRED_SETTING_NAME = "802-3-ethernet"

const (
	// Specific port type to use if multiple the device supports
	// multiple attachment methods. One of 'tp' (Twisted Pair), 'aui'
	// (Attachment Unit Interface), 'bnc' (Thin Ethernet) or 'mii'
	// (Media Independent Interface. If the device supports only one
	// port type, this setting is ignored.
	NM_SETTING_WIRED_PORT = "port"

	// If non-zero, request that the device use only the specified
	// speed. In Mbit/s, ie 100 == 100Mbit/s.
	NM_SETTING_WIRED_SPEED = "speed"

	// If specified, request that the device only use the specified
	// duplex mode. Either 'half' or 'full'.
	NM_SETTING_WIRED_DUPLEX = "duplex"

	// If TRUE, allow auto-negotiation of port speed and duplex
	// mode. If FALSE, do not allow auto-negotiation,in which case the
	// 'speed' and 'duplex' properties should be set.
	NM_SETTING_WIRED_AUTO_NEGOTIATE = "auto-negotiate"

	// If specified, this connection will only apply to the ethernet
	// device whose permanent MAC address matches. This property does
	// not change the MAC address of the device (i.e. MAC spoofing).
	NM_SETTING_WIRED_MAC_ADDRESS = "mac-address"

	// If specified, request that the device use this MAC address
	// instead of its permanent MAC address. This is known as MAC
	// cloning or spoofing.
	NM_SETTING_WIRED_CLONED_MAC_ADDRESS = "cloned-mac-address"

	// If specified, this connection will never apply to the ethernet
	// device whose permanent MAC address matches an address in the
	// list. Each MAC address is in the standard hex-digits-and-colons
	// notation (00:11:22:33:44:55).
	NM_SETTING_WIRED_MAC_ADDRESS_BLACKLIST = "mac-address-blacklist"

	// If non-zero, only transmit packets of the specified size or
	// smaller, breaking larger packets up into multiple Ethernet
	// frames.
	NM_SETTING_WIRED_MTU = "mtu"

	// Identifies specific subchannels that this network device uses
	// for communcation with z/VM or s390 host. Like the 'mac-address'
	// property for non-z/VM devices, this property can be used to
	// ensure this connection only applies to the network device that
	// uses these subchannels. The list should contain exactly 3
	// strings, and each string may only be composed of hexadecimal
	// characters and the period (.) character.
	NM_SETTING_WIRED_S390_SUBCHANNELS = "s390-subchannels"

	// s390 network device type; one of 'qeth', 'lcs', or 'ctc',
	// representing the different types of virtual network devices
	// available on s390 systems.
	NM_SETTING_WIRED_S390_NETTYPE = "s390-nettype"

	// Dictionary of key/value pairs of s390-specific device
	// options. Both keys and values must be strings. Allowed keys
	// include 'portno', 'layer2', 'portname', 'protocol', among
	// others.
	NM_SETTING_WIRED_S390_OPTIONS = "s390-options"
)

func newWireedConnectionData(id, uuid string) (data _ConnectionData) {
	data = make(_ConnectionData)

	addConnectionDataField(data, fieldConnection)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, typeWired)

	// TODO
	addConnectionDataField(data, fieldWired)
	setSettingWiredDuplex(data, "full")

	addConnectionDataField(data, fieldIPv4)
	setSettingIp4ConfigMethod(data, NM_SETTING_IP4_CONFIG_METHOD_AUTO)

	addConnectionDataField(data, fieldIPv6)
	setSettingIp6ConfigMethod(data, NM_SETTING_IP6_CONFIG_METHOD_AUTO)

	return
}

// TODO Get available keys
func getSettingWiredAvailableKeys(data _ConnectionData) (keys []string) {
	keys = []string{
		NM_SETTING_WIRED_MAC_ADDRESS,
		NM_SETTING_WIRED_MTU,
	}
	return
}

// TODO Check whether the values are correct
func checkSettingWiredValues(data _ConnectionData) (errs map[string]string) {
	errs = make(map[string]string)
	return
}

// TODO use logic setter
func generalSetSettingWiredKeyJSON(data _ConnectionData, key, value string) {
	switch key {
	default:
		LOGGER.Error("generalSetSettingWiredKey: invalide key", key)
	case NM_SETTING_WIRED_PORT:
		setSettingWiredPortJSON(data, value)
	case NM_SETTING_WIRED_SPEED:
		setSettingWiredSpeedJSON(data, value)
	case NM_SETTING_WIRED_DUPLEX:
		setSettingWiredDuplexJSON(data, value)
	case NM_SETTING_WIRED_AUTO_NEGOTIATE:
		setSettingWiredAutoNegotiateJSON(data, value)
	case NM_SETTING_WIRED_MAC_ADDRESS:
		setSettingWiredMacAddressJSON(data, value)
	case NM_SETTING_WIRED_CLONED_MAC_ADDRESS:
		setSettingWiredClonedMacAddressJSON(data, value)
	case NM_SETTING_WIRED_MAC_ADDRESS_BLACKLIST:
		setSettingWiredMacAddressBlacklistJSON(data, value)
	case NM_SETTING_WIRED_MTU:
		setSettingWiredMtuJSON(data, value)
	case NM_SETTING_WIRED_S390_SUBCHANNELS:
		setSettingWiredS390SubchannelsJSON(data, value)
	case NM_SETTING_WIRED_S390_NETTYPE:
		setSettingWiredS390NettypeJSON(data, value)
	case NM_SETTING_WIRED_S390_OPTIONS:
		setSettingWiredS390OptionsJSON(data, value)
	}
	return
}
