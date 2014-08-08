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

import "pkg.linuxdeepin.com/lib/dbus"
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
	nmAgentUnregister()
	dbus.UnInstallObject(a)
}

func (a *agent) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		".",
		"/org/freedesktop/NetworkManager/SecretAgent",
		"org.freedesktop.NetworkManager.SecretAgent",
	}
}

func fillSecret(connectionData map[string]map[string]dbus.Variant, settingName, key string) (keyData map[string]map[string]dbus.Variant) {
	keyData = make(map[string]map[string]dbus.Variant)
	keyData[settingName] = make(map[string]dbus.Variant)
	switch settingName {
	case sectionWired: // TODO 8021x
	case sectionWirelessSecurity:
		switch getSettingVkWirelessSecurityKeyMgmt(connectionData) {
		case "none": // ignore
		case "wep":
			setSettingWirelessSecurityWepKey0(keyData, key)
		case "wpa-psk":
			setSettingWirelessSecurityPsk(keyData, key)
		case "wpa-eap":
			// If the user chose an 802.1x-based auth method, return
			// 802.1x secrets, not wireless secrets.
			delete(keyData, settingName)
			keyData[section8021x] = make(map[string]dbus.Variant)
			switch getSettingVk8021xEap(connectionData) {
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
	case sectionVpn: // TODO
	default:
		logger.Error("Unknown secretly setting name", settingName, ", please report it to linuxdeepin")
	}
	return keyData
}

func (a *agent) GetSecrets(connectionData map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath, settingName string, hints []string, flags uint32) (keyData map[string]map[string]dbus.Variant) {
	logger.Info("GetSecrets:", connectionPath, settingName, hints, flags)
	keyId := mapKey{connectionPath, settingName}

	// TODO fixme
	// if key, ok := a.savedKeys[keyId]; ok && (flags&NM_SECRET_AGENT_GET_SECRETS_FLAG_REQUEST_NEW == 0) {
	// logger.Debug("GetSecrets return ", key) // TODO test
	// return key
	// }
	// if flags&NM_SECRET_AGENT_GET_SECRETS_FLAG_USER_REQUESTED == 0 {
	// logger.Debug("GetSecrets return") // TODO test
	// return invalidKeyData
	// }

	if flags&NM_SECRET_AGENT_GET_SECRETS_FLAG_ALLOW_INTERACTION == 0 &&
		flags&NM_SECRET_AGENT_GET_SECRETS_FLAG_USER_REQUESTED == 0 {
		logger.Info("GetSecrets, invalid key flag", flags)
		keyData = invalidKeyData
		return
	}

	if _, ok := a.pendingKeys[keyId]; ok {
		// only wait the key when there is no other GetSecrtes runing with this uuid
		logger.Info("GetSecrets repeatly, cancel last one", keyId)
		a.CancelGetSecrtes(connectionPath, settingName)
	}
	select {
	case key, ok := <-a.createPendingKey(keyId, getSettingConnectionId(connectionData)):
		if ok {
			keyData = fillSecret(connectionData, settingName, key)
			a.SaveSecrets(keyData, connectionPath)
		} else {
			logger.Info("failed to get secretes,", keyId)
		}
	case <-time.After(120 * time.Second):
		a.CancelGetSecrtes(connectionPath, settingName)
		logger.Info("get secrets timeout,", keyId)
	}
	return
}
func (a *agent) createPendingKey(keyId mapKey, connectionId string) chan string {
	// TODO vpn
	// /usr/lib/NetworkManager/nm-pptp-auth-dialog -u fec2a72f-db65-4e76-be37-995932b64bb7 -n pptp -s org.freedesktop.NetworkManager.pptp -i

	logger.Debug("createPendingKey:", keyId, connectionId) // TODO test
	if manager.NeedSecrets != nil {
		logger.Debug("OnNeedSecrets:", string(keyId.path), keyId.name, connectionId)
		defer manager.NeedSecrets(string(keyId.path), keyId.name, connectionId)
	} else {
		logger.Warning("createPendingKey when DNetworkManager hasn't init")
	}
	a.pendingKeys[keyId] = make(chan string)
	return a.pendingKeys[keyId]
}

func (a *agent) CancelGetSecrtes(connectionPath dbus.ObjectPath, settingName string) {
	logger.Debug("CancelGetSecrtes:", connectionPath, settingName)
	keyId := mapKey{connectionPath, settingName}

	if pendingChan, ok := a.pendingKeys[keyId]; ok {
		close(pendingChan)
		delete(a.pendingKeys, keyId)
	} else {
		logger.Warning("CancelGetSecrtes unknown PendingKey", keyId)
	}
}

func (a *agent) SaveSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath) {
	logger.Debug("SaveSecretes:", connectionPath)
	// TODO
	// if _, ok := connection["802-11-wireless-security"]; ok {
	// keyId := mapKey{connectionPath, "802-11-wireless-security"}
	// a.savedKeys[keyId] = connection
	// logger.Debug("SaveSecrets:", connection, connectionPath) // TODO test
	// }
}

func (a *agent) DeleteSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath) {
	logger.Debug("DeleteSecrets:", connectionPath) // TODO test
	if _, ok := connection["802-11-wireless-security"]; ok {
		keyId := mapKey{connectionPath, "802-11-wireless-security"}
		delete(a.savedKeys, keyId)
	}
}

func (a *agent) feedSecret(path string, name string, key string) {
	keyId := mapKey{dbus.ObjectPath(path), name}
	if ch, ok := a.pendingKeys[keyId]; ok {
		ch <- key
		delete(a.pendingKeys, keyId)
	} else {
		logger.Warning("feedSecret, unknown PendingKey", keyId)
	}
}

func (m *Manager) FeedSecret(path string, name, key string) {
	logger.Debug("FeedSecret:", path, name, key)
	m.agent.feedSecret(path, name, key)
}
func (m *Manager) CancelSecret(path string, name string) {
	logger.Debug("CancelSecret:", path, name)
	m.agent.CancelGetSecrtes(dbus.ObjectPath(path), name)
}
