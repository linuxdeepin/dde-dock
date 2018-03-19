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
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"pkg.deepin.io/dde/daemon/network/nm"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	. "pkg.deepin.io/lib/gettext"
)

const agentTimeout = 120 // 120s

var invalidSecretsData = make(map[string]map[string]dbus.Variant)

type mapKey struct {
	connPath    dbus.ObjectPath
	settingName string
}
type agent struct {
	service          *dbusutil.Service // is system service
	pendingKeys      map[mapKey]chan interface{}
	savedKeys        map[mapKey]map[string]map[string]dbus.Variant // TODO: remove
	vpnProcesses     map[dbus.ObjectPath]*os.Process
	vpnProcessesLock sync.Mutex

	secretReceivers *secretProxyType
	receiversLocker sync.Mutex

	methods *struct {
		GetSecrets       func() `in:"connection,connectionPath,settingName,hints,flags" out:"secrets"`
		CancelGetSecrets func() `in:"connectionPath,settingName"`
		SaveSecrets      func() `in:"connection,connectionPath"`
		DeleteSecrets    func() `in:"connection,connectionPath"`
	}
}

// secretsInfo provide more detailed information for front-end to
// pop-up the password authentication dialog.
type secretsInfo struct {
	ConnectionPath dbus.ObjectPath
	SettingName    string

	// ConnectionId is just the connection name which ask user to fill
	// the secrets.
	ConnectionId string

	// AutoConnect tells whether the connection is auto-connect
	// enabled.
	AutoConnect bool

	// KeyType will be used by IsPasswordValid() so the front-end
	// could check the input value on the time.
	KeyType string

	// DevicePath tells which device are sending the signal to ask for
	// secrets. Note that DevicePath will be empty in some cases.
	DevicePath dbus.ObjectPath

	// Receiver tells which client process to ask for secrets.
	// If 0, no clent process was selected.
	Receiver uint32
}

func newAgent(service *dbusutil.Service) (a *agent) {
	a = &agent{}
	a.pendingKeys = make(map[mapKey]chan interface{})
	a.vpnProcesses = make(map[dbus.ObjectPath]*os.Process)
	a.savedKeys = make(map[mapKey]map[string]map[string]dbus.Variant)
	a.secretReceivers = new(secretProxyType)
	a.service = service

	err := a.service.Export("/org/freedesktop/NetworkManager/SecretAgent",
		a)

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
	a.service.StopExport(a)
}

func (*agent) GetInterfaceName() string {
	return "org.freedesktop.NetworkManager.SecretAgent"
}

// TODO: refactor code
// isSecretKey check if target setting key is a secret key which should be stored in keyring
func isSecretKey(connectionData map[string]map[string]dbus.Variant, settingName, keyName string) (isSecret bool) {
	switch settingName {
	case nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME:
		switch keyName {
		case nm.NM_SETTING_WIRELESS_SECURITY_WEP_KEY1, nm.NM_SETTING_WIRELESS_SECURITY_PSK:
			isSecret = true
		}
	case nm.NM_SETTING_802_1X_SETTING_NAME:
		switch keyName {
		case nm.NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD, nm.NM_SETTING_802_1X_PASSWORD:
			isSecret = true
		}
	case nm.NM_SETTING_PPPOE_SETTING_NAME:
		switch keyName {
		case nm.NM_SETTING_PPPOE_PASSWORD:
			isSecret = true
		}
	case nm.NM_SETTING_GSM_SETTING_NAME:
		switch keyName {
		case nm.NM_SETTING_GSM_PASSWORD, nm.NM_SETTING_GSM_PIN:
			isSecret = true
		}
	case nm.NM_SETTING_CDMA_SETTING_NAME:
		switch keyName {
		case nm.NM_SETTING_CDMA_PASSWORD:
			isSecret = true
		}
	case nm.NM_SETTING_VPN_SETTING_NAME:
		if keyName == nm.NM_SETTING_VPN_SECRETS {
			isSecret = true
		}
	}
	return
}

func buildSecretData(connectionData map[string]map[string]dbus.Variant, settingName string, keyValue interface{}) (secretsData map[string]map[string]dbus.Variant) {
	secretsData = make(map[string]map[string]dbus.Variant)
	secretsData[settingName] = make(map[string]dbus.Variant)
	fillSecretData(connectionData, secretsData, settingName, keyValue)
	return secretsData
}
func fillSecretData(connectionData, secretsData map[string]map[string]dbus.Variant, settingName string, keyValueIfc interface{}) {
	// FIXME: some sections support multiple secret keys such as 8021x
	switch settingName {
	case nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME:
		keyValue, _ := keyValueIfc.(string)
		switch getSettingVkWirelessSecurityKeyMgmt(connectionData) {
		case "none": // ignore
		case "wep":
			setSettingWirelessSecurityWepKey0(secretsData, keyValue)
		case "wpa-psk":
			setSettingWirelessSecurityPsk(secretsData, keyValue)
		case "wpa-eap":
			// If the user chose an 802.1x-based auth method, return
			// 802.1x secrets together.
			secretsData[nm.NM_SETTING_802_1X_SETTING_NAME] = make(map[string]dbus.Variant)
			doFillSecret8021x(connectionData, secretsData, keyValue)
		}
	case nm.NM_SETTING_802_1X_SETTING_NAME:
		// wired 8021x
		keyValue, _ := keyValueIfc.(string)
		doFillSecret8021x(connectionData, secretsData, keyValue)
	case nm.NM_SETTING_PPPOE_SETTING_NAME:
		keyValue, _ := keyValueIfc.(string)
		setSettingPppoePassword(secretsData, keyValue)
	case nm.NM_SETTING_VPN_SETTING_NAME:
		keyValue, _ := keyValueIfc.(map[string]string)
		setSettingVpnSecrets(secretsData, keyValue)
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

func buildKeyringSecret(connectionData map[string]map[string]dbus.Variant, settingName string, keyValues map[string]string) (secretsData map[string]map[string]dbus.Variant) {
	secretsData = make(map[string]map[string]dbus.Variant)
	fillKeyringSecret(secretsData, settingName, keyValues)
	return secretsData
}
func fillKeyringSecret(secretsData map[string]map[string]dbus.Variant, settingName string, keyValues map[string]string) {
	if !isSettingExists(secretsData, settingName) {
		addSetting(secretsData, settingName)
	}
	if settingName == nm.NM_SETTING_VPN_SETTING_NAME {
		// FIXME: looks vpn secrets should be ignored here
		vpnSecretData := make(map[string]string)
		// if vpnSecretData, ok := doGetSettingVpnPluginData(secretsData, true); ok {
		for k, v := range keyValues {
			// secret keyValues for vpn should always are string type
			valueStr := marshalVpnPluginKey(v, ktypeString)
			vpnSecretData[k] = valueStr
		}
		// }
		setSettingVpnSecrets(secretsData, vpnSecretData)
	} else {
		for k, v := range keyValues {
			doSetSettingKey(secretsData, settingName, k, v)
		}
	}
}

func (a *agent) GetSecrets(connectionData map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath, settingName string, hints []string, flags uint32) (secretsData map[string]map[string]dbus.Variant, busErr *dbus.Error) {
	logger.Info("GetSecrets:", connectionPath, settingName, hints, flags)

	var ask = false

	// try to get secrets from keyring firstly
	values, ok := secretGetAll(getSettingConnectionUuid(connectionData), settingName)

	// if queried keyring failed will ask for user if we're allowed to
	if !ok && flags&nm.NM_SECRET_AGENT_GET_SECRETS_FLAG_ALLOW_INTERACTION != 0 {
		ask = true
	}

	secretsData = buildKeyringSecret(connectionData, settingName, values)

	// besides, the following cases will ask for user, too
	if flags != nm.NM_SECRET_AGENT_GET_SECRETS_FLAG_NONE {
		if flags&nm.NM_SECRET_AGENT_GET_SECRETS_FLAG_REQUEST_NEW != 0 {
			// the previous secrets are wrong, so ask for user is necessary
			ask = true
		} else if flags&nm.NM_SECRET_AGENT_GET_SECRETS_FLAG_ALLOW_INTERACTION != 0 && isConnectionAlwaysAsk(connectionData, settingName) {
			ask = true
		}
	}

	if !ask {
		return
	}

	logger.Info("askForSecrets:", connectionPath, settingName)

	keyId := mapKey{connPath: connectionPath, settingName: settingName}
	if _, ok := a.pendingKeys[keyId]; ok {
		logger.Info("GetSecrets repeatly, cancel last one", keyId)
		a.cancelGetSecrets(connectionPath, settingName, false)
	}
	select {
	case value, ok := <-a.createPendingKey(connectionData, keyId, hints, flags):
		if ok {
			secretsData = buildSecretData(connectionData, settingName, value)
			a.SaveSecrets(secretsData, connectionPath)
		} else {
			logger.Info("failed to get secretes", keyId)
		}
		if !isVpnConnection(connectionData) {
			manager.service.Emit(manager, "NeedSecretsFinished",
				string(connectionPath), settingName)
		}
	case <-time.After(agentTimeout * time.Second):
		a.cancelGetSecrets(connectionPath, settingName, true)
		logger.Info("get secrets timeout", keyId)
	}
	return
}
func (a *agent) createPendingKey(connectionData map[string]map[string]dbus.Variant, keyId mapKey, hints []string, flags uint32) chan interface{} {
	autoConnect := nmGeneralGetConnectionAutoconnect(keyId.connPath)
	connectionId := getSettingConnectionId(connectionData)
	logger.Debug("createPendingKey:", keyId, connectionId, autoConnect)

	a.pendingKeys[keyId] = make(chan interface{})
	if isVpnConnection(connectionData) {
		// for vpn connections, ask password for vpn auth dialogs
		vpnAuthDilogBin := getVpnAuthDialogBin(connectionData)
		go func() {
			args := []string{
				"-u", getSettingConnectionUuid(connectionData),
				"-n", connectionId,
				"-s", getSettingVpnServiceType(connectionData),
			}
			if flags&nm.NM_SECRET_AGENT_GET_SECRETS_FLAG_ALLOW_INTERACTION != 0 {
				args = append(args, "-i")
			}
			if flags&nm.NM_SECRET_AGENT_GET_SECRETS_FLAG_REQUEST_NEW != 0 {
				args = append(args, "-r")
			}
			// add hints
			for _, h := range hints {
				args = append(args, "-t", h)
			}

			// run vpn auth dialog
			logger.Info("run vpn auth dialog:", vpnAuthDilogBin, args)
			process, stdin, stdout, _, err := execWithIO(vpnAuthDilogBin, args...)
			stdinWriter := bufio.NewWriter(stdin)
			stdoutReader := bufio.NewReader(stdout)

			a.vpnProcessesLock.Lock()
			a.vpnProcesses[keyId.connPath] = process
			a.vpnProcessesLock.Unlock()

			// try to get vpn secrets data from keyring or network manager dbus interface
			vpnData := getSettingVpnData(connectionData)
			var vpnSecretData map[string]string
			vpnSecretData, ok := secretGetAll(getSettingConnectionUuid(connectionData), nm.NM_SETTING_VPN_SETTING_NAME)
			if !ok {
				if secretsData, err := nmGetConnectionSecrets(keyId.connPath, nm.NM_SETTING_VPN_SETTING_NAME); err == nil {
					vpnSecretData = getSettingVpnSecrets(secretsData)
				}
			}

			// send vpn connection data to the authentication dialog binary
			for key, value := range vpnData {
				stdinWriter.WriteString("DATA_KEY=" + key + "\n")
				stdinWriter.WriteString("DATA_VAL=" + value + "\n\n")
			}
			for key, value := range vpnSecretData {
				stdinWriter.WriteString("SECRET_KEY=" + key + "\n")
				stdinWriter.WriteString("SECRET_VAL=" + value + "\n\n")
			}
			stdinWriter.WriteString("DONE\n\n")
			stdinWriter.Flush()

			stdoutData := make(map[string]string)
			lastKey := ""
			// read output until there are two empty lines printed
			empty_lines := 0
			for {
				lineBytes, _, err := stdoutReader.ReadLine()
				if err != nil {
					break
				}
				line := string(lineBytes)

				if len(line) == 0 {
					empty_lines++
				} else {
					// the secrets key and value are split as lines
					if len(lastKey) == 0 {
						lastKey = line
					} else {
						stdoutData[lastKey] = line
						lastKey = ""
					}
				}
				if empty_lines == 2 {
					break
				}
			}

			// notify auth dialog to quit
			stdinWriter.WriteString("QUIT\n\n")
			err = stdinWriter.Flush()
			if err == nil {
				a.feedSecret(keyId.connPath, keyId.settingName, stdoutData, autoConnect)
			} else {
				// mostly, if error occurred for input/output
				// operation, the vpn auth dialog should be killed by
				// cancelVpnAuthDialog() which is triggered for user
				// disconnected the vpn connection
				a.cancelGetSecrets(keyId.connPath, keyId.settingName, false)
			}
		}()
	} else {
		// for none vpn connections, ask password for front-end
		secretsInfo := secretsInfo{
			ConnectionPath: keyId.connPath,
			SettingName:    keyId.settingName,
			ConnectionId:   connectionId,
			AutoConnect:    autoConnect,
			KeyType:        getSettingPassKeyType(connectionData, keyId.settingName),
			DevicePath:     a.guessDevice(connectionData),
		}
		a.receiversLocker.Lock()
		secretsInfo.Receiver = a.secretReceivers.Last()
		a.receiversLocker.Unlock()
		secretsInfoJSON, _ := marshalJSON(secretsInfo)
		notify(notifyIconWirelessDisconnected, "", fmt.Sprintf(Tr("Password required to connect %q"), connectionId))
		manager.service.Emit(manager, "NeedSecrets", secretsInfoJSON)
	}
	return a.pendingKeys[keyId]
}
func (a *agent) guessDevice(connectionData map[string]map[string]dbus.Variant) (devicePath dbus.ObjectPath) {
	switch getSettingConnectionType(connectionData) {
	case nm.NM_SETTING_WIRED_SETTING_NAME:
		return a.doGuessDevice(connectionData, deviceEthernet)
	case nm.NM_SETTING_WIRELESS_SETTING_NAME:
		return a.doGuessDevice(connectionData, deviceWifi)
	}
	return
}
func (a *agent) doGuessDevice(connectionData map[string]map[string]dbus.Variant, deviceType string) (devicePath dbus.ObjectPath) {
	manager.devicesLock.Lock()
	defer manager.devicesLock.Unlock()

	devices := manager.devices[deviceType]

	// check for the hardware address
	var hwAddressBytes []byte
	switch deviceType {
	case deviceEthernet:
		hwAddressBytes = getSettingWiredMacAddress(connectionData)
	case deviceWifi:
		hwAddressBytes = getSettingWirelessMacAddress(connectionData)
	}

	if len(hwAddressBytes) != 0 {
		hwAddress := convertMacAddressToString(hwAddressBytes)
		for _, device := range devices {
			if hwAddress == device.HwAddress {
				return device.Path
			}
		}
	}

	// check for the device state
	for _, device := range devices {
		if isDeviceStateInActivating(device.State) {
			return device.Path
		}
	}

	// if all failed, and there is only one device now, just return it
	if len(devices) == 1 {
		return devices[0].Path
	}

	return
}

func (a *agent) cancelVpnAuthDialog(connPath dbus.ObjectPath) {
	a.vpnProcessesLock.Lock()
	defer a.vpnProcessesLock.Unlock()

	for p, process := range a.vpnProcesses {
		if p == connPath {
			process.Kill()
			break
		}
	}
	delete(a.vpnProcesses, connPath)
}

func (a *agent) CancelGetSecrets(connectionPath dbus.ObjectPath, settingName string) *dbus.Error {
	logger.Info("CancelGetSecrets:", connectionPath, settingName)
	a.cancelGetSecrets(connectionPath, settingName, true)
	return nil
}

func (a *agent) cancelGetSecrets(connectionPath dbus.ObjectPath, settingName string, notifyFinished bool) {
	logger.Debug("cancelGetSecrets", connectionPath, settingName, notifyFinished)
	keyId := mapKey{connPath: connectionPath, settingName: settingName}

	if notifyFinished {
		manager.service.Emit(manager, "NeedSecretsFinished",
			string(connectionPath), settingName)
	}

	if pendingChan, ok := a.pendingKeys[keyId]; ok {
		close(pendingChan)
		delete(a.pendingKeys, keyId)
	} else {
		logger.Warning("CancelGetSecrets unknown PendingKey", keyId)
	}
}

func (a *agent) SaveSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath) *dbus.Error {
	logger.Info("SaveSecretes:", connectionPath)
	return nil
}

func (a *agent) DeleteSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath) *dbus.Error {
	// TODO delete secrets from keyring
	logger.Info("DeleteSecrets:", connectionPath)
	if _, ok := connection["802-11-wireless-security"]; ok {
		keyId := mapKey{connPath: connectionPath, settingName: "802-11-wireless-security"}
		delete(a.savedKeys, keyId)
	}
	return nil
}

func (a *agent) feedSecret(path dbus.ObjectPath, settingName string, keyValue interface{}, autoConnect bool) {
	keyId := mapKey{connPath: path, settingName: settingName}
	if ch, ok := a.pendingKeys[keyId]; ok {
		ch <- keyValue
		delete(a.pendingKeys, keyId)
	} else {
		logger.Warning("feedSecret, unknown PendingKey", keyId)
	}

	// update secret data in connection settings manually to fix
	// password popup issue when editing such connections
	data, err := nmGetConnectionData(path)
	if err != nil {
		return
	}
	generalSetSettingAutoconnect(data, autoConnect)
	fillSecretData(data, data, settingName, keyValue)
	nmUpdateConnectionData(path, data)
}

func (m *Manager) FeedSecret(path string, settingName, keyValue string, autoConnect bool) *dbus.Error {
	logger.Info("FeedSecret:", path, settingName, "xxxx")
	m.agent.feedSecret(dbus.ObjectPath(path), settingName, keyValue, autoConnect)
	return nil
}

func (m *Manager) CancelSecret(path string, settingName string) *dbus.Error {
	logger.Info("CancelSecret:", path, settingName)
	m.agent.cancelGetSecrets(dbus.ObjectPath(path), settingName, true)
	return nil
}

func (m *Manager) RegisterSecretReceiver(sender dbus.Sender) *dbus.Error {
	if m.agent == nil {
		logger.Info("Agent object no created")
		return nil
	}
	pid, err := m.service.GetConnPID(string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}

	m.agent.receiversLocker.Lock()
	m.agent.secretReceivers.Add(pid)
	m.agent.receiversLocker.Unlock()
	return nil
}
