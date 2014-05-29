package network

import "dbus/org/freedesktop/notifications"
import "dlib/dbus"
import "sync"
import nm "dbus/org/freedesktop/networkmanager"
import "dlib"

const (
	NM_STATE_UNKNOWN          = uint32(0)
	NM_STATE_ASLEEP           = uint32(10)
	NM_STATE_DISCONNECTED     = uint32(20)
	NM_STATE_DISCONNECTING    = uint32(30)
	NM_STATE_CONNECTING       = uint32(40)
	NM_STATE_CONNECTED_LOCAL  = uint32(50)
	NM_STATE_CONNECTED_SITE   = uint32(60)
	NM_STATE_CONNECTED_GLOBAL = uint32(70)
)

const (
	NM_DEVICE_STATE_UNKNOWN      = 0   // The device is in an unknown state.
	NM_DEVICE_STATE_UNMANAGED    = 10  // The device is recognized but not managed by NetworkManager.
	NM_DEVICE_STATE_UNAVAILABLE  = 20  // The device cannot be used (carrier off, rfkill, etc).
	NM_DEVICE_STATE_DISCONNECTED = 30  // The device is not connected.
	NM_DEVICE_STATE_PREPARE      = 40  // The device is preparing to connect.
	NM_DEVICE_STATE_CONFIG       = 50  // The device is being configured.
	NM_DEVICE_STATE_NEED_AUTH    = 60  // The device is awaiting secrets necessary to continue connection.
	NM_DEVICE_STATE_IP_CONFIG    = 70  // The IP settings of the device are being requested and configured.
	NM_DEVICE_STATE_IP_CHECK     = 80  // The device's IP connectivity ability is being determined.
	NM_DEVICE_STATE_SECONDARIES  = 90  // The device is waiting for secondary connections to be activated.
	NM_DEVICE_STATE_ACTIVATED    = 100 // The device is active.
	NM_DEVICE_STATE_DEACTIVATING = 110 // The device's network connection is being torn down.
	NM_DEVICE_STATE_FAILED       = 120 // The device is in a failure state following an attempt to activate it.
)

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

var VPNErrorTable = make(map[uint32]string)
var DEVICEErrorTable = make(map[uint32]string)

func initReasons() {
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_NOW_MANAGED] = dlib.Tr("The device is now managed.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_NOW_UNMANAGED] = dlib.Tr("The device is no longer managed.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_CONFIG_FAILED] = dlib.Tr("The device could not be readied for configuration.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_CONFIG_UNAVAILABLE] = dlib.Tr("IP configuration could not be reserved (no available address timeout etc).")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_CONFIG_EXPIRED] = dlib.Tr("The IP configuration is no longer valid.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_NO_SECRETS] = dlib.Tr("Secrets were required but not provided.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_SUPPLICANT_DISCONNECT] = dlib.Tr("The 802.1X supplicant disconnected from the access point or authentication server.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_SUPPLICANT_CONFIG_FAILED] = dlib.Tr("Configuration of the 802.1X supplicant failed.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_SUPPLICANT_FAILED] = dlib.Tr("The 802.1X supplicant quit or failed unexpectedly.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_SUPPLICANT_TIMEOUT] = dlib.Tr("The 802.1X supplicant took too long to authenticate.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_PPP_START_FAILED] = dlib.Tr("The PPP service failed to start within the allowed time.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_PPP_DISCONNECT] = dlib.Tr("The PPP service disconnected unexpectedly.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_PPP_FAILED] = dlib.Tr("The PPP service quit or failed unexpectedly.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_DHCP_START_FAILED] = dlib.Tr("The DHCP service failed to start within the allowed time.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_DHCP_ERROR] = dlib.Tr("The DHCP service reported an unexpected error.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_DHCP_FAILED] = dlib.Tr("The DHCP service quit or failed unexpectedly.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_SHARED_START_FAILED] = dlib.Tr("The shared connection service failed to start.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_SHARED_FAILED] = dlib.Tr("The shared connection service quit or failed unexpectedly.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_AUTOIP_START_FAILED] = dlib.Tr("The AutoIP service failed to start.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_AUTOIP_ERROR] = dlib.Tr("The AutoIP service reported an unexpected error.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_AUTOIP_FAILED] = dlib.Tr("The AutoIP service quit or failed unexpectedly.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_MODEM_BUSY] = dlib.Tr("Dialing failed because the line was busy.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_MODEM_NO_DIAL_TONE] = dlib.Tr("Dialing failed because there was no dial tone.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_MODEM_NO_CARRIER] = dlib.Tr("Dialing failed because there was carrier.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_MODEM_DIAL_TIMEOUT] = dlib.Tr("Dialing timed out.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_MODEM_DIAL_FAILED] = dlib.Tr("Dialing failed.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_MODEM_INIT_FAILED] = dlib.Tr("Modem initialization failed.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_GSM_APN_FAILED] = dlib.Tr("Failed to select the specified GSM APN.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_GSM_REGISTRATION_NOT_SEARCHING] = dlib.Tr("Not searching for networks.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_GSM_REGISTRATION_DENIED] = dlib.Tr("Network registration was denied.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_GSM_REGISTRATION_TIMEOUT] = dlib.Tr("Network registration timed out.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_GSM_REGISTRATION_FAILED] = dlib.Tr("Failed to register with the requested GSM network.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_GSM_PIN_CHECK_FAILED] = dlib.Tr("PIN check failed.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_FIRMWARE_MISSING] = dlib.Tr("Necessary firmware for the device may be missing.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_REMOVED] = dlib.Tr("The device was removed.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_SLEEPING] = dlib.Tr("NetworkManager went to sleep.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_CONNECTION_REMOVED] = dlib.Tr("The device's active connection was removed or disappeared.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_USER_REQUESTED] = dlib.Tr("A user or client requested the disconnection.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_CARRIER] = dlib.Tr("The device's carrier/link changed.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_CONNECTION_ASSUMED] = dlib.Tr("The device's existing connection was assumed.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_SUPPLICANT_AVAILABLE] = dlib.Tr("The 802.1x supplicant is now available.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_MODEM_NOT_FOUND] = dlib.Tr("The modem could not be found.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_BT_FAILED] = dlib.Tr("The Bluetooth connection timed out or failed.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_GSM_SIM_NOT_INSERTED] = dlib.Tr("GSM Modem's SIM Card not inserted.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_GSM_SIM_PIN_REQUIRED] = dlib.Tr("GSM Modem's SIM Pin required.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_GSM_SIM_PUK_REQUIRED] = dlib.Tr("GSM Modem's SIM Puk required.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_GSM_SIM_WRONG] = dlib.Tr("GSM Modem's SIM wrong")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_INFINIBAND_MODE] = dlib.Tr("InfiniBand device does not support connected mode.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_DEPENDENCY_FAILED] = dlib.Tr("A dependency of the connection failed.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_BR2684_FAILED] = dlib.Tr("Problem with the RFC 2684 Ethernet over ADSL bridge.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_MODEM_MANAGER_UNAVAILABLE] = dlib.Tr("ModemManager was not running or quit unexpectedly.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_SSID_NOT_FOUND] = dlib.Tr("The 802.11 Wi-Fi network could not be found.")
	DEVICEErrorTable[NM_DEVICE_STATE_REASON_SECONDARY_CONNECTION_FAILED] = dlib.Tr("A secondary connection of the base connection failed.")

	VPNErrorTable[NM_VPN_CONNECTION_STATE_REASON_USER_DISCONNECTED] = dlib.Tr("The VPN connection changed state because the user disconnected it.")
	VPNErrorTable[NM_VPN_CONNECTION_STATE_REASON_DEVICE_DISCONNECTED] = dlib.Tr("The VPN connection %s changed state because the device it was using was disconnected.")
	VPNErrorTable[NM_VPN_CONNECTION_STATE_REASON_SERVICE_STOPPED] = dlib.Tr("The service providing the VPN connection was stopped.")
	VPNErrorTable[NM_VPN_CONNECTION_STATE_REASON_IP_CONFIG_INVALID] = dlib.Tr("The IP config of the VPN connection was invalid.")
	VPNErrorTable[NM_VPN_CONNECTION_STATE_REASON_CONNECT_TIMEOUT] = dlib.Tr("The connection attempt to the VPN service timed out.")
	VPNErrorTable[NM_VPN_CONNECTION_STATE_REASON_SERVICE_START_TIMEOUT] = dlib.Tr("A timeout occurred while starting the service providing the VPN connection.")
	VPNErrorTable[NM_VPN_CONNECTION_STATE_REASON_SERVICE_START_FAILED] = dlib.Tr("Starting the service starting the service providing the VPN connection failed.")
	VPNErrorTable[NM_VPN_CONNECTION_STATE_REASON_NO_SECRETS] = dlib.Tr("Necessary secrets for the VPN connection were not provided.")
	VPNErrorTable[NM_VPN_CONNECTION_STATE_REASON_LOGIN_FAILED] = dlib.Tr("Authentication to the VPN server failed.")
	VPNErrorTable[NM_VPN_CONNECTION_STATE_REASON_CONNECTION_REMOVED] = dlib.Tr("The connection was deleted from settings.")
}

func init() {
	var notify *notifications.Notifier
	notify, _ = notifications.NewNotifier("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	currentDevices := make(map[dbus.ObjectPath]*nm.Device)
	var locker sync.Mutex

	var watch = func(path dbus.ObjectPath) {
		defer func() {
			if err := recover(); err != nil {
				locker.Lock()
				defer locker.Unlock()
				delete(currentDevices, path)
				logger.Error(err)
			}
		}()
		if dev, err := nm.NewDevice("org.freedesktop.NetworkManager", path); err == nil {
			locker.Lock()
			defer locker.Unlock()
			currentDevices[path] = dev
			dev.ConnectStateChanged(func(newState, oldState, reason uint32) {
				switch newState {
				case NM_DEVICE_STATE_ACTIVATED:
					ac, _ := nm.NewActiveConnection("org.freedesktop.NetworkManager", dev.ActiveConnection.Get())
					cc, _ := nm.NewSettingsConnection("org.freedesktop.NetworkManager", ac.Connection.Get())
					data, _ := cc.GetSettings()
					var icon string
					switch getCustomConnectionType(data) {
					case connectionWired:
						icon = "notification-network-ethernet-connected"
					case connectionWireless:
						// TODO
						// icon = "notification-network-wireless-connected"
						icon = "notification-network-wireless-full"
					default:
						icon = "network-transmit-receive"
					}
					notify.Notify("Network", 0, icon, dlib.Tr("Connected"), getSettingConnectionId(data), nil, nil, 0)
				case NM_DEVICE_STATE_FAILED, NM_DEVICE_STATE_DISCONNECTED,
					NM_DEVICE_STATE_UNMANAGED, NM_DEVICE_STATE_UNAVAILABLE:
					switch oldState {
					case NM_DEVICE_STATE_FAILED, NM_DEVICE_STATE_DISCONNECTED,
						NM_DEVICE_STATE_UNMANAGED, NM_DEVICE_STATE_UNAVAILABLE:
						// this was not a disconnect
					default:
						//TODO: icon name can be different by device type
						if reason != NM_DEVICE_STATE_REASON_NONE && reason != NM_DEVICE_STATE_REASON_UNKNOWN {
							var icon string
							switch dev.DeviceType.Get() {
							case NM_DEVICE_TYPE_ETHERNET:
								icon = "notification-network-wired-disconnected"
							case NM_DEVICE_TYPE_WIFI:
								icon = "notification-network-wireless-disconnected"
							default:
								icon = "network-error"
							}
							// TODO: show connection name when disconnected
							notify.Notify("Network", 0, icon, dlib.Tr("Disconnect"), DEVICEErrorTable[reason], nil, nil, 0)
						}
					}
				}
			})
		}
	}
	var remove = func(path dbus.ObjectPath) {
		locker.Lock()
		defer locker.Unlock()
		if dev, ok := currentDevices[path]; ok {
			nmDestroyDevice(dev)
		}
	}

	nmManager.ConnectDeviceRemoved(func(path dbus.ObjectPath) {
		remove(path)
	})
	nmManager.ConnectDeviceAdded(func(path dbus.ObjectPath) {
		watch(path)
	})

	paths, _ := nmManager.GetDevices()
	for _, path := range paths {
		watch(path)
	}

	nmManager.ConnectStateChanged(func(state uint32) {
		switch state {
		case NM_STATE_DISCONNECTED, NM_STATE_ASLEEP:
			notify.Notify("Network", 0, "network-offline", dlib.Tr("Offline"), dlib.Tr("Disconnected - You are now offline."), nil, nil, 0)
		}
	})

	//TODO: VPN state
}
