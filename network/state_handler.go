package main

import "dbus/org/freedesktop/notifications"
import "fmt"

const (
	NM_DEVICE_STATE_REASON_UNKNOWN                        = 0
	NM_DEVICE_STATE_REASON_NONE                           = 1
	NM_DEVICE_STATE_REASON_NOW_MANAGED                    = 2
	NM_DEVICE_STATE_REASON_NOW_UNMANAGED                  = 3
	NM_DEVICE_STATE_REASON_CONFIG_FAILED                  = 4
	NM_DEVICE_STATE_REASON_CONFIG_UNAVAILABLE             = 5
	NM_DEVICE_STATE_REASON_CONFIG_EXPIRED                 = 6
	NM_DEVICE_STATE_REASON_NO_SECRETS                     = 7
	NM_DEVICE_STATE_REASON_SUPPLICANT_DISCONNECT          = 8
	NM_DEVICE_STATE_REASON_SUPPLICANT_CONFIG_FAILED       = 9
	NM_DEVICE_STATE_REASON_SUPPLICANT_FAILED              = 10
	NM_DEVICE_STATE_REASON_SUPPLICANT_TIMEOUT             = 11
	NM_DEVICE_STATE_REASON_PPP_START_FAILED               = 12
	NM_DEVICE_STATE_REASON_PPP_DISCONNECT                 = 13
	NM_DEVICE_STATE_REASON_PPP_FAILED                     = 14
	NM_DEVICE_STATE_REASON_DHCP_START_FAILED              = 15
	NM_DEVICE_STATE_REASON_DHCP_ERROR                     = 16
	NM_DEVICE_STATE_REASON_DHCP_FAILED                    = 17
	NM_DEVICE_STATE_REASON_SHARED_START_FAILED            = 18
	NM_DEVICE_STATE_REASON_SHARED_FAILED                  = 19
	NM_DEVICE_STATE_REASON_AUTOIP_START_FAILED            = 20
	NM_DEVICE_STATE_REASON_AUTOIP_ERROR                   = 21
	NM_DEVICE_STATE_REASON_AUTOIP_FAILED                  = 22
	NM_DEVICE_STATE_REASON_MODEM_BUSY                     = 23
	NM_DEVICE_STATE_REASON_MODEM_NO_DIAL_TONE             = 24
	NM_DEVICE_STATE_REASON_MODEM_NO_CARRIER               = 25
	NM_DEVICE_STATE_REASON_MODEM_DIAL_TIMEOUT             = 26
	NM_DEVICE_STATE_REASON_MODEM_DIAL_FAILED              = 27
	NM_DEVICE_STATE_REASON_MODEM_INIT_FAILED              = 28
	NM_DEVICE_STATE_REASON_GSM_APN_FAILED                 = 29
	NM_DEVICE_STATE_REASON_GSM_REGISTRATION_NOT_SEARCHING = 30
	NM_DEVICE_STATE_REASON_GSM_REGISTRATION_DENIED        = 31
	NM_DEVICE_STATE_REASON_GSM_REGISTRATION_TIMEOUT       = 32
	NM_DEVICE_STATE_REASON_GSM_REGISTRATION_FAILED        = 33
	NM_DEVICE_STATE_REASON_GSM_PIN_CHECK_FAILED           = 34
	NM_DEVICE_STATE_REASON_FIRMWARE_MISSING               = 35
	NM_DEVICE_STATE_REASON_REMOVED                        = 36
	NM_DEVICE_STATE_REASON_SLEEPING                       = 37
	NM_DEVICE_STATE_REASON_CONNECTION_REMOVED             = 38
	NM_DEVICE_STATE_REASON_USER_REQUESTED                 = 39
	NM_DEVICE_STATE_REASON_CARRIER                        = 40
	NM_DEVICE_STATE_REASON_CONNECTION_ASSUMED             = 41
	NM_DEVICE_STATE_REASON_SUPPLICANT_AVAILABLE           = 42
	NM_DEVICE_STATE_REASON_MODEM_NOT_FOUND                = 43
	NM_DEVICE_STATE_REASON_BT_FAILED                      = 44
	NM_DEVICE_STATE_REASON_GSM_SIM_NOT_INSERTED           = 45
	NM_DEVICE_STATE_REASON_GSM_SIM_PIN_REQUIRED           = 46
	NM_DEVICE_STATE_REASON_GSM_SIM_PUK_REQUIRED           = 47
	NM_DEVICE_STATE_REASON_GSM_SIM_WRONG                  = 48
	NM_DEVICE_STATE_REASON_INFINIBAND_MODE                = 49
	NM_DEVICE_STATE_REASON_DEPENDENCY_FAILED              = 50
	NM_DEVICE_STATE_REASON_BR2684_FAILED                  = 51
	NM_DEVICE_STATE_REASON_MODEM_MANAGER_UNAVAILABLE      = 52
	NM_DEVICE_STATE_REASON_SSID_NOT_FOUND                 = 53
	NM_DEVICE_STATE_REASON_SECONDARY_CONNECTION_FAILED    = 54
)

var DEVICEErrorTable = map[uint32]string{
	NM_DEVICE_STATE_REASON_UNKNOWN:                        "The reason for the device state change is unknown.",
	NM_DEVICE_STATE_REASON_NONE:                           "The state change is normal.",
	NM_DEVICE_STATE_REASON_NOW_MANAGED:                    "The device is now managed.",
	NM_DEVICE_STATE_REASON_NOW_UNMANAGED:                  "The device is no longer managed.",
	NM_DEVICE_STATE_REASON_CONFIG_FAILED:                  "The device could not be readied for configuration.",
	NM_DEVICE_STATE_REASON_CONFIG_UNAVAILABLE:             "IP configuration could not be reserved (no available address, timeout, etc).",
	NM_DEVICE_STATE_REASON_CONFIG_EXPIRED:                 "The IP configuration is no longer valid.",
	NM_DEVICE_STATE_REASON_NO_SECRETS:                     "Secrets were required, but not provided.",
	NM_DEVICE_STATE_REASON_SUPPLICANT_DISCONNECT:          "The 802.1X supplicant disconnected from the access point or authentication server.",
	NM_DEVICE_STATE_REASON_SUPPLICANT_CONFIG_FAILED:       "Configuration of the 802.1X supplicant failed.",
	NM_DEVICE_STATE_REASON_SUPPLICANT_FAILED:              "The 802.1X supplicant quit or failed unexpectedly.",
	NM_DEVICE_STATE_REASON_SUPPLICANT_TIMEOUT:             "The 802.1X supplicant took too long to authenticate.",
	NM_DEVICE_STATE_REASON_PPP_START_FAILED:               "The PPP service failed to start within the allowed time.",
	NM_DEVICE_STATE_REASON_PPP_DISCONNECT:                 "The PPP service disconnected unexpectedly.",
	NM_DEVICE_STATE_REASON_PPP_FAILED:                     "The PPP service quit or failed unexpectedly.",
	NM_DEVICE_STATE_REASON_DHCP_START_FAILED:              "The DHCP service failed to start within the allowed time.",
	NM_DEVICE_STATE_REASON_DHCP_ERROR:                     "The DHCP service reported an unexpected error.",
	NM_DEVICE_STATE_REASON_DHCP_FAILED:                    "The DHCP service quit or failed unexpectedly.",
	NM_DEVICE_STATE_REASON_SHARED_START_FAILED:            "The shared connection service failed to start.",
	NM_DEVICE_STATE_REASON_SHARED_FAILED:                  "The shared connection service quit or failed unexpectedly.",
	NM_DEVICE_STATE_REASON_AUTOIP_START_FAILED:            "The AutoIP service failed to start.",
	NM_DEVICE_STATE_REASON_AUTOIP_ERROR:                   "The AutoIP service reported an unexpected error.",
	NM_DEVICE_STATE_REASON_AUTOIP_FAILED:                  "The AutoIP service quit or failed unexpectedly.",
	NM_DEVICE_STATE_REASON_MODEM_BUSY:                     "Dialing failed because the line was busy.",
	NM_DEVICE_STATE_REASON_MODEM_NO_DIAL_TONE:             "Dialing failed because there was no dial tone.",
	NM_DEVICE_STATE_REASON_MODEM_NO_CARRIER:               "Dialing failed because there was carrier.",
	NM_DEVICE_STATE_REASON_MODEM_DIAL_TIMEOUT:             "Dialing timed out.",
	NM_DEVICE_STATE_REASON_MODEM_DIAL_FAILED:              "Dialing failed.",
	NM_DEVICE_STATE_REASON_MODEM_INIT_FAILED:              "Modem initialization failed.",
	NM_DEVICE_STATE_REASON_GSM_APN_FAILED:                 "Failed to select the specified GSM APN.",
	NM_DEVICE_STATE_REASON_GSM_REGISTRATION_NOT_SEARCHING: "Not searching for networks.",
	NM_DEVICE_STATE_REASON_GSM_REGISTRATION_DENIED:        "Network registration was denied.",
	NM_DEVICE_STATE_REASON_GSM_REGISTRATION_TIMEOUT:       "Network registration timed out.",
	NM_DEVICE_STATE_REASON_GSM_REGISTRATION_FAILED:        "Failed to register with the requested GSM network.",
	NM_DEVICE_STATE_REASON_GSM_PIN_CHECK_FAILED:           "PIN check failed.",
	NM_DEVICE_STATE_REASON_FIRMWARE_MISSING:               "Necessary firmware for the device may be missing.",
	NM_DEVICE_STATE_REASON_REMOVED:                        "The device was removed.",
	NM_DEVICE_STATE_REASON_SLEEPING:                       "NetworkManager went to sleep.",
	NM_DEVICE_STATE_REASON_CONNECTION_REMOVED:             "The device's active connection was removed or disappeared.",
	NM_DEVICE_STATE_REASON_USER_REQUESTED:                 "A user or client requested the disconnection.",
	NM_DEVICE_STATE_REASON_CARRIER:                        "The device's carrier/link changed.",
	NM_DEVICE_STATE_REASON_CONNECTION_ASSUMED:             "The device's existing connection was assumed.",
	NM_DEVICE_STATE_REASON_SUPPLICANT_AVAILABLE:           "The 802.1x supplicant is now available.",
	NM_DEVICE_STATE_REASON_MODEM_NOT_FOUND:                "The modem could not be found.",
	NM_DEVICE_STATE_REASON_BT_FAILED:                      "The Bluetooth connection timed out or failed.",
	NM_DEVICE_STATE_REASON_GSM_SIM_NOT_INSERTED:           "GSM Modem's SIM Card not inserted.",
	NM_DEVICE_STATE_REASON_GSM_SIM_PIN_REQUIRED:           "GSM Modem's SIM Pin required.",
	NM_DEVICE_STATE_REASON_GSM_SIM_PUK_REQUIRED:           "GSM Modem's SIM Puk required.",
	NM_DEVICE_STATE_REASON_GSM_SIM_WRONG:                  "GSM Modem's SIM wrong",
	NM_DEVICE_STATE_REASON_INFINIBAND_MODE:                "InfiniBand device does not support connected mode.",
	NM_DEVICE_STATE_REASON_DEPENDENCY_FAILED:              "A dependency of the connection failed.",
	NM_DEVICE_STATE_REASON_BR2684_FAILED:                  "Problem with the RFC 2684 Ethernet over ADSL bridge.",
	NM_DEVICE_STATE_REASON_MODEM_MANAGER_UNAVAILABLE:      "ModemManager was not running or quit unexpectedly.",
	NM_DEVICE_STATE_REASON_SSID_NOT_FOUND:                 "The 802.11 Wi-Fi network could not be found.",
	NM_DEVICE_STATE_REASON_SECONDARY_CONNECTION_FAILED:    "A secondary connection of the base connection failed.",
}

const (
	//don't use iota, the value is defined by networkmanager
	NM_VPN_CONNECTION_STATE_REASON_UNKNOWN               = 0
	NM_VPN_CONNECTION_STATE_REASON_NONE                  = 1
	NM_VPN_CONNECTION_STATE_REASON_USER_DISCONNECTED     = 2
	NM_VPN_CONNECTION_STATE_REASON_DEVICE_DISCONNECTED   = 3
	NM_VPN_CONNECTION_STATE_REASON_SERVICE_STOPPED       = 4
	NM_VPN_CONNECTION_STATE_REASON_IP_CONFIG_INVALID     = 5
	NM_VPN_CONNECTION_STATE_REASON_CONNECT_TIMEOUT       = 6
	NM_VPN_CONNECTION_STATE_REASON_SERVICE_START_TIMEOUT = 7
	NM_VPN_CONNECTION_STATE_REASON_SERVICE_START_FAILED  = 8
	NM_VPN_CONNECTION_STATE_REASON_NO_SECRETS            = 9
	NM_VPN_CONNECTION_STATE_REASON_LOGIN_FAILED          = 10
	NM_VPN_CONNECTION_STATE_REASON_CONNECTION_REMOVED    = 11
)

var VPNErrorTable = map[uint32]string{
	NM_VPN_CONNECTION_STATE_REASON_UNKNOWN:               "The reason for the VPN connection state change is unknown.",
	NM_VPN_CONNECTION_STATE_REASON_NONE:                  "No reason was given for the VPN connection state change.",
	NM_VPN_CONNECTION_STATE_REASON_USER_DISCONNECTED:     "The VPN connection changed state because the user disconnected it.",
	NM_VPN_CONNECTION_STATE_REASON_DEVICE_DISCONNECTED:   "The VPN connection %s changed state because the device it was using was disconnected.",
	NM_VPN_CONNECTION_STATE_REASON_SERVICE_STOPPED:       "The service providing the VPN connection was stopped.",
	NM_VPN_CONNECTION_STATE_REASON_IP_CONFIG_INVALID:     "The IP config of the VPN connection was invalid.",
	NM_VPN_CONNECTION_STATE_REASON_CONNECT_TIMEOUT:       "The connection attempt to the VPN service timed out.",
	NM_VPN_CONNECTION_STATE_REASON_SERVICE_START_TIMEOUT: "A timeout occurred while starting the service providing the VPN connection.",
	NM_VPN_CONNECTION_STATE_REASON_SERVICE_START_FAILED:  "Starting the service starting the service providing the VPN connection failed.",
	NM_VPN_CONNECTION_STATE_REASON_NO_SECRETS:            "Necessary secrets for the VPN connection were not provided.",
	NM_VPN_CONNECTION_STATE_REASON_LOGIN_FAILED:          "Authentication to the VPN server failed.",
	NM_VPN_CONNECTION_STATE_REASON_CONNECTION_REMOVED:    "The connection was deleted from settings.",
}

func handleStateChanged(id string, reason uint32) {
	notify, _ := notifications.NewNotifier("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	notify.Notify("Network", 0, "network.error", "Network failed", fmt.Sprintf(VPNErrorTable[reason], id), nil, nil, 0)
}
