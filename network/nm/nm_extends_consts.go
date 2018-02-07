/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package nm

// VPN L2TP
const (
	NM_DBUS_SERVICE_L2TP   = "org.freedesktop.NetworkManager.l2tp"
	NM_DBUS_INTERFACE_L2TP = "org.freedesktop.NetworkManager.l2tp"
	NM_DBUS_PATH_L2TP      = "/org/freedesktop/NetworkManager/l2tp"
)

const (
	NM_DBUS_SERVICE_L2TP_PPP   = "org.freedesktop.NetworkManager.l2tp-ppp"
	NM_DBUS_PATH_L2TP_PPP      = "/org/freedesktop/NetworkManager/l2tp/ppp"
	NM_DBUS_INTERFACE_L2TP_PPP = "org.freedesktop.NetworkManager.l2tp.ppp"
)

const (
	NM_L2TP_SECRET_FLAG_NONE         = 0 // system saved
	NM_L2TP_SECRET_FLAG_AGENT_OWNED  = 1
	NM_L2TP_SECRET_FLAG_NOT_SAVED    = 3
	NM_L2TP_SECRET_FLAG_NOT_REQUIRED = 5
)

// VPN OpenConnect
const (
	NM_DBUS_SERVICE_OPENCONNECT   = "org.freedesktop.NetworkManager.openconnect"
	NM_DBUS_INTERFACE_OPENCONNECT = "org.freedesktop.NetworkManager.openconnect"
	NM_DBUS_PATH_OPENCONNECT      = "/org/freedesktop/NetworkManager/openconnect"
)

// VPN OpenVPN
const (
	NM_DBUS_SERVICE_OPENVPN   = "org.freedesktop.NetworkManager.openvpn"
	NM_DBUS_INTERFACE_OPENVPN = "org.freedesktop.NetworkManager.openvpn"
	NM_DBUS_PATH_OPENVPN      = "/org/freedesktop/NetworkManager/openvpn"
)

const (
	/* Internal auth-dialog -> service token indicating that no secrets are
	 * required for the connection.
	 */
	NM_SETTING_VPN_OPENVPN_KEY_NOSECRET = "no-secret"
)

const (
	NM_OPENVPN_CONTYPE_TLS          = "tls"
	NM_OPENVPN_CONTYPE_STATIC_KEY   = "static-key"
	NM_OPENVPN_CONTYPE_PASSWORD     = "password"
	NM_OPENVPN_CONTYPE_PASSWORD_TLS = "password-tls"
)
const (
	NM_OPENVPN_AUTH_NONE      = "none"
	NM_OPENVPN_AUTH_RSA_MD4   = "RSA-MD4"
	NM_OPENVPN_AUTH_MD5       = "MD5"
	NM_OPENVPN_AUTH_SHA1      = "SHA1"
	NM_OPENVPN_AUTH_SHA224    = "SHA224"
	NM_OPENVPN_AUTH_SHA256    = "SHA256"
	NM_OPENVPN_AUTH_SHA384    = "SHA384"
	NM_OPENVPN_AUTH_SHA512    = "SHA512"
	NM_OPENVPN_AUTH_RIPEMD160 = "RIPEMD160"
)
const (
	NM_OPENVPN_REM_CERT_TLS_CLIENT = "client"
	NM_OPENVPN_REM_CERT_TLS_SERVER = "server"
)

// Define secret flags
const (
	NM_OPENVPN_SECRET_FLAG_SAVE   = 0
	NM_OPENVPN_SECRET_FLAG_ASK    = 2
	NM_OPENVPN_SECRET_FLAG_UNUSED = 4
)

// VPN PP2P
const (
	NM_DBUS_SERVICE_PPTP   = "org.freedesktop.NetworkManager.pptp"
	NM_DBUS_INTERFACE_PPTP = "org.freedesktop.NetworkManager.pptp"
	NM_DBUS_PATH_PPTP      = "/org/freedesktop/NetworkManager/pptp"
)

// Define secret flags
const (
	NM_PPTP_SECRET_FLAG_NONE         = 0
	NM_PPTP_SECRET_FLAG_AGENT_OWNED  = 1
	NM_PPTP_SECRET_FLAG_NOT_SAVED    = 3
	NM_PPTP_SECRET_FLAG_NOT_REQUIRED = 5
)

// VPN StrongSwan
const (
	NM_DBUS_SERVICE_STRONGSWAN = "org.freedesktop.NetworkManager.strongswan"
)

const (
	NM_STRONGSWAN_METHOD_KEY       = "key"
	NM_STRONGSWAN_METHOD_AGENT     = "agent"
	NM_STRONGSWAN_METHOD_SMARTCARD = "smartcard"
	NM_STRONGSWAN_METHOD_EAP       = "eap"
	NM_STRONGSWAN_METHOD_PSK       = "psk"
)

// VPN VPNC
const (
	NM_DBUS_SERVICE_VPNC   = "org.freedesktop.NetworkManager.vpnc"
	NM_DBUS_INTERFACE_VPNC = "org.freedesktop.NetworkManager.vpnc"
	NM_DBUS_PATH_VPNC      = "/org/freedesktop/NetworkManager/vpnc"
)

const (
	NM_VPNC_NATT_MODE_NATT        = "natt"
	NM_VPNC_NATT_MODE_NONE        = "none"
	NM_VPNC_NATT_MODE_NATT_ALWAYS = "force-natt"
	NM_VPNC_NATT_MODE_CISCO       = "cisco-udp"
)
const (
	NM_VPNC_PW_TYPE_SAVE   = "save"   // -> flags 1
	NM_VPNC_PW_TYPE_ASK    = "ask"    // -> flags 3
	NM_VPNC_PW_TYPE_UNUSED = "unused" // -> flags 5
)
const (
	NM_VPNC_DHGROUP_DH1 = "dh1"
	NM_VPNC_DHGROUP_DH2 = "dh2"
	NM_VPNC_DHGROUP_DH5 = "dh5"
)
const (
	NM_VPNC_PFS_SERVER = "server"
	NM_VPNC_PFS_NOPFS  = "nopfs"
	NM_VPNC_PFS_DH1    = "dh1"
	NM_VPNC_PFS_DH2    = "dh2"
	NM_VPNC_PFS_DH5    = "dh5"
)
const (
	NM_VPNC_VENDOR_CISCO     = "cisco"
	NM_VPNC_VENDOR_NETSCREEN = "netscreen"
)

// Define secret flags
const (
	NM_VPNC_SECRET_FLAG_NONE   = 0
	NM_VPNC_SECRET_FLAG_SAVE   = 1
	NM_VPNC_SECRET_FLAG_ASK    = 3
	NM_VPNC_SECRET_FLAG_UNUSED = 5
)
