package main

// Virtual fields, used for vpn connectionns.

// vpn-l2tp
const (
	NM_SETTING_VF_VPN_L2TP_SETTING_NAME       = "vf-vpn-l2tp"
	NM_SETTING_VF_VPN_L2TP_PPP_SETTING_NAME   = "vf-vpn-l2tp-ppp"
	NM_SETTING_VF_VPN_L2TP_IPSEC_SETTING_NAME = "vf-vpn-l2tp-ipsec"
)

// vpn-openconnect
const (
	NM_SETTING_VF_VPN_OPENCONNECT_SETTING_NAME = "vf-vpn-openconnect"
)

// vpn-openvpn TODO
const (
	NM_SETTING_VF_VPN_OPENVPN_SETTING_NAME = "vf-vpn-openvpn"
)

// vpn-pptp
const (
	NM_SETTING_VF_VPN_PPTP_SETTING_NAME     = "vf-vpn-pptp"
	NM_SETTING_VF_VPN_PPTP_PPP_SETTING_NAME = "vf-vpn-pptp-ppp"
)

// vpn-vpnc TODO
const (
	NM_SETTING_VF_VPN_VPNC_SETTING_NAME = "vf-vpn-vpnc"
)

func getRealFieldName(name string) (realName string) {
	realName = name
	switch name {
	case NM_SETTING_VF_VPN_L2TP_SETTING_NAME:
		realName = fieldVpn
	case NM_SETTING_VF_VPN_L2TP_PPP_SETTING_NAME:
		realName = fieldVpn
	case NM_SETTING_VF_VPN_L2TP_IPSEC_SETTING_NAME:
		realName = fieldVpn
	case NM_SETTING_VF_VPN_OPENCONNECT_SETTING_NAME:
		realName = fieldVpn
	case NM_SETTING_VF_VPN_OPENVPN_SETTING_NAME:
		realName = fieldVpn
	case NM_SETTING_VF_VPN_PPTP_SETTING_NAME:
		realName = fieldVpn
	case NM_SETTING_VF_VPN_PPTP_PPP_SETTING_NAME:
		realName = fieldVpn
	case NM_SETTING_VF_VPN_VPNC_SETTING_NAME:
		realName = fieldVpn
	}
	return
}
