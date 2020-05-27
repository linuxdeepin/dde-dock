/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package network

import (
	"os"
	"path/filepath"

	"pkg.deepin.io/lib/keyfile"
)

const (
	nmVpnL2tpNameFile        = "nm-l2tp-service.name"
	nmVpnOpenconnectNameFile = "nm-openconnect-service.name"
	nmVpnOpenvpnNameFile     = "nm-openvpn-service.name"
	nmVpnPptpNameFile        = "nm-pptp-service.name"
	nmVpnStrongswanNameFile  = "nm-strongswan-service.name"
	nmVpnVpncNameFile        = "nm-vpnc-service.name"
)

const (
	nmOpenConnectServiceType = "org.freedesktop.NetworkManager.openconnect"
)

func getVpnAuthDialogBin(data connectionData) (authdialog string) {
	vpnType := getCustomConnectionType(data)
	return doGetVpnAuthDialogBin(vpnType)
}

func doGetVpnAuthDialogBin(vpnType string) (authdialog string) {
	k := keyfile.NewKeyFile()
	k.LoadFromFile(getVpnNameFile(vpnType))
	authdialog, _ = k.GetString("GNOME", "auth-dialog")

	return
}

func getVpnNameFile(vpnType string) (nameFile string) {
	var baseName string
	switch vpnType {
	case connectionVpnL2tp:
		baseName = nmVpnL2tpNameFile
	case connectionVpnOpenconnect:
		baseName = nmVpnOpenconnectNameFile
	case connectionVpnOpenvpn:
		baseName = nmVpnOpenvpnNameFile
	case connectionVpnPptp:
		baseName = nmVpnPptpNameFile
	case connectionVpnStrongswan:
		baseName = nmVpnStrongswanNameFile
	case connectionVpnVpnc:
		baseName = nmVpnVpncNameFile
	default:
		return ""
	}

	for _, dir := range []string{"/etc/NetworkManager/VPN", "/usr/lib/NetworkManager/VPN"} {
		nameFile = filepath.Join(dir, baseName)
		_, err := os.Stat(nameFile)
		if err == nil {
			return nameFile
		}
	}

	return ""
}
