package main

// Virtual fields, used for vpn connectionns.

// vpn-l2tp
const (
	NM_SETTING_VF_VPN_L2TP_SETTING_NAME       = "vf-vpn-l2tp"
	NM_SETTING_VF_VPN_L2TP_PPP_SETTING_NAME   = "vf-vpn-l2tp-ppp"
	NM_SETTING_VF_VPN_L2TP_IPSEC_SETTING_NAME = "vf-vpn-l2tp-ipsec"
)

func getRealFiledName(name string) (realName string) {
	realName = name
	switch name {
	case NM_SETTING_VF_VPN_L2TP_SETTING_NAME:
		realName = NM_SETTING_VPN_SETTING_NAME
	case NM_SETTING_VF_VPN_L2TP_PPP_SETTING_NAME:
		realName = NM_SETTING_VPN_SETTING_NAME
	case NM_SETTING_VF_VPN_L2TP_IPSEC_SETTING_NAME:
		realName = NM_SETTING_VPN_SETTING_NAME
	}
	return
}
