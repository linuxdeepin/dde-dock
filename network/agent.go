/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

import "pkg.deepin.io/lib/dbus"
import "time"

var invalidSecretsData = make(map[string]map[string]dbus.Variant)

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

// TODO: refactor code
// isSecretKey check if target setting key is a secret key which should be stored in keyring
func isSecretKey(connectionData map[string]map[string]dbus.Variant, settingName, keyName string) (isSecret bool) {
	switch settingName {
	case NM_SETTING_WIRELESS_SECURITY_SETTING_NAME:
		switch keyName {
		case NM_SETTING_WIRELESS_SECURITY_WEP_KEY1, NM_SETTING_WIRELESS_SECURITY_PSK:
			isSecret = true
		}
	case NM_SETTING_802_1X_SETTING_NAME:
		switch keyName {
		case NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD, NM_SETTING_802_1X_PASSWORD:
			isSecret = true
		}
	case NM_SETTING_PPPOE_SETTING_NAME:
		switch keyName {
		case NM_SETTING_PPPOE_PASSWORD:
			isSecret = true
		}
	case NM_SETTING_GSM_SETTING_NAME:
		switch keyName {
		case NM_SETTING_GSM_PASSWORD, NM_SETTING_GSM_PIN:
			isSecret = true
		}
	case NM_SETTING_CDMA_SETTING_NAME:
		switch keyName {
		case NM_SETTING_CDMA_PASSWORD:
			isSecret = true
		}
	case NM_SETTING_VPN_SETTING_NAME:
		if keyName == NM_SETTING_VPN_SECRETS {
			isSecret = true
		}
	}
	return
}

// FIXME: some sections support multiple secret keys such as 8021x and vpn
func buildSecretData(connectionData map[string]map[string]dbus.Variant, settingName, value string) (secretsData map[string]map[string]dbus.Variant) {
	secretsData = make(map[string]map[string]dbus.Variant)
	secretsData[settingName] = make(map[string]dbus.Variant)
	fillSecretData(connectionData, secretsData, settingName, value)
	return secretsData
}
func fillSecretData(connectionData, secretsData map[string]map[string]dbus.Variant, settingName, value string) {
	switch settingName {
	case sectionWirelessSecurity:
		switch getSettingVkWirelessSecurityKeyMgmt(connectionData) {
		case "none": // ignore
		case "wep":
			setSettingWirelessSecurityWepKey0(secretsData, value)
		case "wpa-psk":
			setSettingWirelessSecurityPsk(secretsData, value)
		case "wpa-eap":
			// If the user chose an 802.1x-based auth method, return
			// 802.1x secrets together.
			secretsData[section8021x] = make(map[string]dbus.Variant)
			doFillSecret8021x(connectionData, secretsData, value)
		}
	case section8021x:
		// wired 8021x
		doFillSecret8021x(connectionData, secretsData, value)
	case sectionPppoe:
		setSettingPppoePassword(secretsData, value)
	case sectionVpn:
		setSettingVpnSecrets(secretsData, make(map[string]string))
		switch getCustomConnectionType(connectionData) {
		case connectionVpnL2tp:
			setSettingVpnL2tpKeyPassword(secretsData, value)
		case connectionVpnOpenconnect: // ignore
		case connectionVpnOpenvpn:
			setSettingVpnOpenvpnKeyPassword(secretsData, value)
			// setSettingVpnOpenvpnKeyCertpass
			// setSettingVpnOpenvpnKeyHttpProxyPassword
		case connectionVpnPptp:
			setSettingVpnPptpKeyPassword(secretsData, value)
		case connectionVpnStrongswan:
			setSettingVpnStrongswanKeyPassword(secretsData, value)
		case connectionVpnVpnc:
			setSettingVpnVpncKeySecret(secretsData, value)
			// setSettingVpnVpncKeyXauthPassword(secretsData, value)
		}
	default:
		logger.Error("Unknown secretly setting name", settingName, ", please report it to linuxdeepin")
	}
}
func doFillSecret8021x(connectionData, secretsData map[string]map[string]dbus.Variant, value string) {
	switch getSettingVk8021xEap(connectionData) {
	case "tls":
		setSetting8021xPrivateKeyPassword(secretsData, value)
	case "md5":
		setSetting8021xPassword(secretsData, value)
	case "leap":
		// LEAP secrets aren't in the 802.1x setting, just ignore
	case "fast":
		setSetting8021xPassword(secretsData, value)
	case "ttls":
		setSetting8021xPassword(secretsData, value)
	case "peap":
		setSetting8021xPassword(secretsData, value)
	}
}

func buildKeyringSecret(connectionData map[string]map[string]dbus.Variant, settingName string, values map[string]string) (secretsData map[string]map[string]dbus.Variant) {
	secretsData = make(map[string]map[string]dbus.Variant)
	fillKeyringSecret(secretsData, settingName, values)
	return secretsData
}
func fillKeyringSecret(secretsData map[string]map[string]dbus.Variant, settingName string, values map[string]string) {
	if !isSettingSectionExists(secretsData, settingName) {
		addSettingSection(secretsData, settingName)
	}
	if settingName == NM_SETTING_VPN_SETTING_NAME {
		// FIXME: looks vpn secrets should be ignored here
		vpnSecretData := make(map[string]string)
		// if vpnSecretData, ok := doGetSettingVpnPluginData(secretsData, true); ok {
		for k, v := range values {
			// secret values for vpn should always are string type
			valueStr := marshalVpnPluginKey(v, ktypeString)
			vpnSecretData[k] = valueStr
		}
		// }
		setSettingVpnSecrets(secretsData, vpnSecretData)
	} else {
		for k, v := range values {
			doSetSettingKey(secretsData, settingName, k, v)
		}
	}
}

func (a *agent) GetSecrets(connectionData map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath, settingName string, hints []string, flags uint32) (secretsData map[string]map[string]dbus.Variant) {
	logger.Info("GetSecrets:", connectionPath, settingName, hints, flags)

	// TODO: VPN passwords should be handled by the VPN plugin's auth dialog themselves

	var ask = false

	// try to get secrets from keyring firstly
	values, ok := secretGetAll(getSettingConnectionUuid(connectionData), settingName)

	// if queried keyring failed will ask for user if we're allowed to
	if !ok && flags&NM_SECRET_AGENT_GET_SECRETS_FLAG_ALLOW_INTERACTION != 0 {
		ask = true
	}

	secretsData = buildKeyringSecret(connectionData, settingName, values)

	// besides, the following cases will ask for user, too
	if flags != NM_SECRET_AGENT_GET_SECRETS_FLAG_NONE {
		if flags&NM_SECRET_AGENT_GET_SECRETS_FLAG_REQUEST_NEW != 0 {
			// the previous secrets are wrong, so ask for user is necessary
			ask = true
		} else if flags&NM_SECRET_AGENT_GET_SECRETS_FLAG_ALLOW_INTERACTION != 0 && isConnectionAlwaysAsk(connectionData, settingName) {
			ask = true
		}
	}

	if !ask {
		return
	}

	logger.Info("askForSecrets:", connectionPath, settingName)

	keyId := mapKey{connectionPath, settingName}
	if _, ok := a.pendingKeys[keyId]; ok {
		logger.Info("GetSecrets repeatly, cancel last one", keyId)
		a.CancelGetSecrets(connectionPath, settingName, false)
	}
	select {
	case value, ok := <-a.createPendingKey(connectionData, keyId):
		if ok {
			secretsData = buildSecretData(connectionData, settingName, value)
			a.SaveSecrets(secretsData, connectionPath)
		} else {
			logger.Info("failed to get secretes", keyId)
		}
		dbus.Emit(manager, "NeedSecretsFinished", string(connectionPath), settingName)
	case <-time.After(120 * time.Second):
		a.CancelGetSecrets(connectionPath, settingName, true)
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

func (a *agent) CancelGetSecrets(connectionPath dbus.ObjectPath, settingName string, notifyFinished bool) {
	logger.Debug("CancelGetSecrets:", connectionPath, settingName)
	keyId := mapKey{connectionPath, settingName}

	if notifyFinished {
		dbus.Emit(manager, "NeedSecretsFinished", string(connectionPath), settingName)
	}

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

func (a *agent) feedSecret(path dbus.ObjectPath, settingName string, keyName string) {
	keyId := mapKey{path, settingName}
	if ch, ok := a.pendingKeys[keyId]; ok {
		ch <- keyName
		delete(a.pendingKeys, keyId)
	} else {
		logger.Warning("feedSecret, unknown PendingKey", keyId)
	}
}

func (m *Manager) FeedSecret(path string, settingName, keyName string, autoConnect bool) {
	logger.Debug("FeedSecret:", path, settingName, "xxxx")

	opath := dbus.ObjectPath(path)
	m.agent.feedSecret(opath, settingName, keyName)

	// FIXME: update secret data in connection settings manually to fix
	// password popup issue when editing such connections
	data, err := nmGetConnectionData(opath)
	if err != nil {
		return
	}
	generalSetSettingAutoconnect(data, autoConnect)
	fillSecretData(data, data, settingName, keyName)
	nmUpdateConnectionData(opath, data)
}
func (m *Manager) CancelSecret(path string, settingName string) {
	logger.Debug("CancelSecret:", path, settingName)
	m.agent.CancelGetSecrets(dbus.ObjectPath(path), settingName, true)
}
