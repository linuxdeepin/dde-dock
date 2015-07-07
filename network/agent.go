/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package network

import "pkg.deepin.io/lib/dbus"
import "time"

var invalidKeyData = make(map[string]map[string]dbus.Variant)

type mapKey struct {
	path dbus.ObjectPath
	name string
}
type agent struct {
	pendingKeys map[mapKey]chan string
	savedKeys   map[mapKey]map[string]map[string]dbus.Variant
}

func newAgent() (a *agent) {
	a = &agent{}
	a.pendingKeys = make(map[mapKey]chan string)
	a.savedKeys = make(map[mapKey]map[string]map[string]dbus.Variant)

	err := dbus.InstallOnSystem(a)
	if err != nil {
		logger.Error("install network agent failed:", err)
		return
	}

	nmAgentRegister("com.deepin.daemon.Network.agent")
	return
}

func destroyAgent(a *agent) {
	for key, ch := range a.pendingKeys {
		close(ch)
		delete(a.pendingKeys, key)
	}
	nmAgentUnregister()
	dbus.UnInstallObject(a)
}

func (a *agent) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       ".",
		ObjectPath: "/org/freedesktop/NetworkManager/SecretAgent",
		Interface:  "org.freedesktop.NetworkManager.SecretAgent",
	}
}

// FIXME: some section support multiple secret keys like 8021x and vpn
func fillSecret(connectionData map[string]map[string]dbus.Variant, settingName, key string) (keyData map[string]map[string]dbus.Variant) {
	keyData = make(map[string]map[string]dbus.Variant)
	keyData[settingName] = make(map[string]dbus.Variant)
	doFillSecret(connectionData, keyData, settingName, key)
	return keyData
}
func doFillSecret(refData, keyData map[string]map[string]dbus.Variant, settingName, key string) {
	switch settingName {
	case sectionWirelessSecurity:
		switch getSettingVkWirelessSecurityKeyMgmt(refData) {
		case "none": // ignore
		case "wep":
			setSettingWirelessSecurityWepKey0(keyData, key)
		case "wpa-psk":
			setSettingWirelessSecurityPsk(keyData, key)
		case "wpa-eap":
			// If the user chose an 802.1x-based auth method, return
			// 802.1x secrets together.
			keyData[section8021x] = make(map[string]dbus.Variant)
			doFillSecret8021x(refData, keyData, key)
		}
	case section8021x:
		// wired 8021x
		doFillSecret8021x(refData, keyData, key)
	case sectionPppoe:
		setSettingPppoePassword(keyData, key)
	case sectionVpn:
		setSettingVpnSecrets(keyData, make(map[string]string))
		switch getCustomConnectionType(refData) {
		case connectionVpnL2tp:
			setSettingVpnL2tpKeyPassword(keyData, key)
		case connectionVpnOpenconnect: // ignore
		case connectionVpnOpenvpn:
			setSettingVpnOpenvpnKeyPassword(keyData, key)
			// setSettingVpnOpenvpnKeyCertpass
			// setSettingVpnOpenvpnKeyHttpProxyPassword
		case connectionVpnPptp:
			setSettingVpnPptpKeyPassword(keyData, key)
		case connectionVpnVpnc:
			setSettingVpnVpncKeySecret(keyData, key)
			// setSettingVpnVpncKeyXauthPassword(keyData, key)
		}
	default:
		logger.Error("Unknown secretly setting name", settingName, ", please report it to linuxdeepin")
	}
}
func doFillSecret8021x(refData, keyData map[string]map[string]dbus.Variant, key string) {
	switch getSettingVk8021xEap(refData) {
	case "tls":
		setSetting8021xPrivateKeyPassword(keyData, key)
	case "md5":
		setSetting8021xPassword(keyData, key)
	case "leap":
		// LEAP secrets aren't in the 802.1x setting, just ignore
	case "fast":
		setSetting8021xPassword(keyData, key)
	case "ttls":
		setSetting8021xPassword(keyData, key)
	case "peap":
		setSetting8021xPassword(keyData, key)
	}
}

func (a *agent) GetSecrets(connectionData map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath, settingName string, hints []string, flags uint32) (keyData map[string]map[string]dbus.Variant) {
	logger.Info("GetSecrets:", connectionPath, settingName, hints, flags)
	keyId := mapKey{connectionPath, settingName}

	if flags&NM_SECRET_AGENT_GET_SECRETS_FLAG_ALLOW_INTERACTION == 0 &&
		flags&NM_SECRET_AGENT_GET_SECRETS_FLAG_USER_REQUESTED == 0 {
		logger.Info("GetSecrets, invalid key flag", flags)
		keyData = invalidKeyData
		return
	}

	if _, ok := a.pendingKeys[keyId]; ok {
		logger.Info("GetSecrets repeatly, cancel last one", keyId)
		a.CancelGetSecrets(connectionPath, settingName)
	}
	select {
	case key, ok := <-a.createPendingKey(connectionData, keyId):
		if ok {
			keyData = fillSecret(connectionData, settingName, key)
			a.SaveSecrets(keyData, connectionPath)
		} else {
			logger.Info("failed to get secretes", keyId)
		}
	case <-time.After(120 * time.Second):
		a.CancelGetSecrets(connectionPath, settingName)
		logger.Info("get secrets timeout", keyId)
	}
	return
}
func (a *agent) createPendingKey(connectionData map[string]map[string]dbus.Variant, keyId mapKey) chan string {
	autoConnect := nmGeneralGetConnectionAutoconnect(keyId.path)
	connectionId := getSettingConnectionId(connectionData)
	logger.Debug("createPendingKey:", keyId, connectionId, autoConnect)

	a.pendingKeys[keyId] = make(chan string)
	dbus.Emit(manager, "NeedSecrets", string(keyId.path), keyId.name, connectionId, autoConnect)
	return a.pendingKeys[keyId]
}

func (a *agent) CancelGetSecrets(connectionPath dbus.ObjectPath, settingName string) {
	logger.Debug("CancelGetSecrets:", connectionPath, settingName)
	keyId := mapKey{connectionPath, settingName}

	if pendingChan, ok := a.pendingKeys[keyId]; ok {
		close(pendingChan)
		delete(a.pendingKeys, keyId)
	} else {
		logger.Warning("CancelGetSecrets unknown PendingKey", keyId)
	}
}

func (a *agent) SaveSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath) {
	logger.Debug("SaveSecretes:", connectionPath)
}

func (a *agent) DeleteSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath) {
	if _, ok := connection["802-11-wireless-security"]; ok {
		keyId := mapKey{connectionPath, "802-11-wireless-security"}
		delete(a.savedKeys, keyId)
	}
}

func (a *agent) feedSecret(path dbus.ObjectPath, name string, key string) {
	keyId := mapKey{path, name}
	if ch, ok := a.pendingKeys[keyId]; ok {
		ch <- key
		delete(a.pendingKeys, keyId)
	} else {
		logger.Warning("feedSecret, unknown PendingKey", keyId)
	}
}

func (m *Manager) FeedSecret(path string, name, key string, autoConnect bool) {
	logger.Debug("FeedSecret:", path, name, "xxxx")

	opath := dbus.ObjectPath(path)
	m.agent.feedSecret(opath, name, key)

	// FIXME: update secret data in connection settings manually to fix
	// password popup issue when editing such connections
	data, err := nmGetConnectionData(opath)
	if err != nil {
		return
	}
	generalSetSettingAutoconnect(data, autoConnect)
	doFillSecret(data, data, name, key)
	nmUpdateConnectionData(opath, data)
}
func (m *Manager) CancelSecret(path string, name string) {
	logger.Debug("CancelSecret:", path, name)
	m.agent.CancelGetSecrets(dbus.ObjectPath(path), name)
}
