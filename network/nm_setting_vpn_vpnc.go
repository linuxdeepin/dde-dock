package main

// #define NM_DBUS_SERVICE_VPNC    "org.freedesktop.NetworkManager.vpnc"
// #define NM_DBUS_INTERFACE_VPNC  "org.freedesktop.NetworkManager.vpnc"
// #define NM_DBUS_PATH_VPNC       "/org/freedesktop/NetworkManager/vpnc"

// #define NM_VPNC_KEY_GATEWAY "IPSec gateway"
// #define NM_VPNC_KEY_ID "IPSec ID"
// #define NM_VPNC_KEY_SECRET "IPSec secret"
// #define NM_VPNC_KEY_SECRET_TYPE "ipsec-secret-type"
// #define NM_VPNC_KEY_XAUTH_USER "Xauth username"
// #define NM_VPNC_KEY_XAUTH_PASSWORD "Xauth password"
// #define NM_VPNC_KEY_XAUTH_PASSWORD_TYPE "xauth-password-type"
// #define NM_VPNC_KEY_DOMAIN "Domain"
// #define NM_VPNC_KEY_DHGROUP "IKE DH Group"
// #define NM_VPNC_KEY_PERFECT_FORWARD "Perfect Forward Secrecy"
// #define NM_VPNC_KEY_VENDOR "Vendor"
// #define NM_VPNC_KEY_APP_VERSION "Application Version"
// #define NM_VPNC_KEY_SINGLE_DES "Enable Single DES"
// #define NM_VPNC_KEY_NO_ENCRYPTION "Enable no encryption"
// #define NM_VPNC_KEY_NAT_TRAVERSAL_MODE "NAT Traversal Mode"
// #define NM_VPNC_KEY_DPD_IDLE_TIMEOUT "DPD idle timeout (our side)"
// #define NM_VPNC_KEY_CISCO_UDP_ENCAPS_PORT "Cisco UDP Encapsulation Port"
// #define NM_VPNC_KEY_LOCAL_PORT "Local Port"
// #define NM_VPNC_KEY_AUTHMODE "IKE Authmode"
// #define NM_VPNC_KEY_CA_FILE "CA-File"

// #define NM_VPNC_NATT_MODE_NATT        "natt"
// #define NM_VPNC_NATT_MODE_NONE        "none"
// #define NM_VPNC_NATT_MODE_NATT_ALWAYS "force-natt"
// #define NM_VPNC_NATT_MODE_CISCO       "cisco-udp"

// #define NM_VPNC_PW_TYPE_SAVE   "save"
// #define NM_VPNC_PW_TYPE_ASK    "ask"
// #define NM_VPNC_PW_TYPE_UNUSED "unused"

// #define NM_VPNC_DHGROUP_DH1 "dh1"
// #define NM_VPNC_DHGROUP_DH2 "dh2"
// #define NM_VPNC_DHGROUP_DH5 "dh5"

// #define NM_VPNC_PFS_SERVER "server"
// #define NM_VPNC_PFS_NOPFS  "nopfs"
// #define NM_VPNC_PFS_DH1    "dh1"
// #define NM_VPNC_PFS_DH2    "dh2"
// #define NM_VPNC_PFS_DH5    "dh5"

// #define NM_VPNC_VENDOR_CISCO     "cisco"
// #define NM_VPNC_VENDOR_NETSCREEN "netscreen"

// static ValidProperty valid_properties[] = {
// 	{ NM_VPNC_KEY_GATEWAY,               ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_ID,                    ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_XAUTH_USER,            ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_DOMAIN,                ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_DHGROUP,               ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_PERFECT_FORWARD,       ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_VENDOR,                ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_APP_VERSION,           ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_SINGLE_DES,            ITEM_TYPE_BOOLEAN, 0, 0 },
// 	{ NM_VPNC_KEY_NO_ENCRYPTION,         ITEM_TYPE_BOOLEAN, 0, 0 },
// 	{ NM_VPNC_KEY_DPD_IDLE_TIMEOUT,      ITEM_TYPE_INT, 0, 86400 },
// 	{ NM_VPNC_KEY_NAT_TRAVERSAL_MODE,    ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_CISCO_UDP_ENCAPS_PORT, ITEM_TYPE_INT, 0, 65535 },
// 	{ NM_VPNC_KEY_LOCAL_PORT,            ITEM_TYPE_INT, 0, 65535 },
// 	/* Hybrid Auth */
// 	{ NM_VPNC_KEY_AUTHMODE,              ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_CA_FILE,               ITEM_TYPE_PATH, 0, 0 },
// 	/* Ignored option for internal use */
// 	{ NM_VPNC_KEY_SECRET_TYPE,           ITEM_TYPE_IGNORED, 0, 0 },
// 	{ NM_VPNC_KEY_XAUTH_PASSWORD_TYPE,   ITEM_TYPE_IGNORED, 0, 0 },
// 	{ NM_VPNC_KEY_SECRET"-flags",        ITEM_TYPE_IGNORED, 0, 0 },
// 	{ NM_VPNC_KEY_XAUTH_PASSWORD"-flags",ITEM_TYPE_IGNORED, 0, 0 },
// 	/* Legacy options that are ignored */
// 	{ LEGACY_NAT_KEEPALIVE,              ITEM_TYPE_STRING, 0, 0 },
// 	{ NULL,                              ITEM_TYPE_UNKNOWN, 0, 0 }
// }
// static ValidProperty valid_secrets[] = {
// 	{ NM_OPENVPN_KEY_PASSWORD,             G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_CERTPASS,             G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_NOSECRET,             G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_HTTP_PROXY_PASSWORD,  G_TYPE_STRING, 0, 0, FALSE },
// 	{ NULL,                                G_TYPE_NONE, FALSE }
// };
// static ValidProperty valid_secrets[] = {
// 	{ NM_VPNC_KEY_SECRET,                ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_XAUTH_PASSWORD,        ITEM_TYPE_STRING, 0, 0 },
// 	{ NULL,                              ITEM_TYPE_UNKNOWN, 0, 0 }
// };
