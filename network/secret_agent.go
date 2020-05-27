package network

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	secrets "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.secrets"
	"pkg.deepin.io/dde/daemon/network/nm"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/strv"
)

const (
	nmSecretDialogBin              = "/usr/lib/deepin-daemon/dnetwork-secret-dialog"
	getSecretsFlagAllowInteraction = 0x1
	getSecretsFlagRequestNew       = 0x2
	getSecretsFlagUserRequested    = 0x4

	secretFlagNone          = 0 // save for all user
	secretFlagNoneStr       = "0"
	secretFlagAgentOwned    = 1 // save for me
	secretFlagAgentOwnedStr = "1"
	secretFlagAsk           = 2 // always ask
	secretFlagAskStr        = "2"
	secretFlagNotRequired   = 4 // no need password

	// keep keyring tags same with nm-applet
	keyringTagConnUUID    = "connection-uuid"
	keyringTagSettingName = "setting-name"
	keyringTagSettingKey  = "setting-key"
)

type saveSecretsTaskKey struct {
	connPath    dbus.ObjectPath
	settingName string
}

type saveSecretsTask struct {
	process *os.Process
}

type SecretAgent struct {
	secretService       *secrets.Service
	secretSessionPath   dbus.ObjectPath
	defaultCollection   *secrets.Collection
	defaultCollectionMu sync.Mutex

	saveSecretsTasks   map[saveSecretsTaskKey]saveSecretsTask
	saveSecretsTasksMu sync.Mutex

	m *Manager

	methods *struct {
		GetSecrets        func() `in:"connection,connectionPath,settingName,hints,flags" out:"secrets"`
		CancelGetSecrets  func() `in:"connectionPath,settingName"`
		SaveSecrets       func() `in:"connection,connectionPath"`
		SaveSecretsDeepin func() `in:"connection,connectionPath"`
		DeleteSecrets     func() `in:"connection,connectionPath"`
	}
}

func (sa *SecretAgent) addSaveSecretsTask(connPath dbus.ObjectPath,
	settingName string, process *os.Process) {
	sa.saveSecretsTasksMu.Lock()

	sa.saveSecretsTasks[saveSecretsTaskKey{
		connPath:    connPath,
		settingName: settingName,
	}] = saveSecretsTask{process: process}

	sa.saveSecretsTasksMu.Unlock()
}

func (sa *SecretAgent) removeSaveSecretsTask(connPath dbus.ObjectPath,
	settingName string) {
	sa.saveSecretsTasksMu.Lock()
	delete(sa.saveSecretsTasks, saveSecretsTaskKey{
		connPath:    connPath,
		settingName: settingName,
	})
	sa.saveSecretsTasksMu.Unlock()
}

func (sa *SecretAgent) getSaveSecretsTaskProcess(connPath dbus.ObjectPath,
	settingName string) *os.Process {
	sa.saveSecretsTasksMu.Lock()

	task := sa.saveSecretsTasks[saveSecretsTaskKey{
		connPath:    connPath,
		settingName: settingName,
	}]

	sa.saveSecretsTasksMu.Unlock()
	return task.process
}

func (sa *SecretAgent) getDefaultCollection() (*secrets.Collection, error) {
	sa.defaultCollectionMu.Lock()
	defer sa.defaultCollectionMu.Unlock()

	if sa.defaultCollection != nil {
		return sa.defaultCollection, nil
	}

	collectionPath, err := sa.secretService.ReadAlias(0, "default")
	if err != nil {
		return nil, err
	}

	if collectionPath == "/" {
		return nil, errors.New("failed to get default collection path")
	}

	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}

	collectionObj, err := secrets.NewCollection(sessionBus, collectionPath)
	if err == nil {
		sa.defaultCollection = collectionObj
	}
	return collectionObj, err
}

func newSecretAgent(secServiceObj *secrets.Service, manager *Manager) (*SecretAgent, error) {
	_, sessionPath, err := secServiceObj.OpenSession(0, "plain", dbus.MakeVariant(""))
	if err != nil {
		return nil, err
	}

	sa := &SecretAgent{}
	sa.secretSessionPath = sessionPath
	sa.secretService = secServiceObj
	sa.saveSecretsTasks = make(map[saveSecretsTaskKey]saveSecretsTask)
	sa.m = manager
	logger.Debug("session path:", sessionPath)
	return sa, nil
}

func (sa *SecretAgent) deleteAll(uuid string) error {
	attributes := map[string]string{
		keyringTagConnUUID: uuid,
	}

	defaultCollection, err := sa.getDefaultCollection()
	if err != nil {
		return err
	}

	items, err := defaultCollection.SearchItems(0, attributes)
	if err != nil {
		return err
	}
	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return err
	}

	for _, itemPath := range items {
		itemObj, err := secrets.NewItem(sessionBus, itemPath)
		if err != nil {
			continue
		}
		_, err = itemObj.Delete(0)
		if err != nil {
			logger.Warningf("delete item %q failed: %v\n", itemPath, err)
		}
	}
	return nil
}

func (sa *SecretAgent) getAll(uuid, settingName string) (map[string]string, error) {
	attributes := map[string]string{
		keyringTagConnUUID:    uuid,
		keyringTagSettingName: settingName,
	}
	defaultCollection, err := sa.getDefaultCollection()
	if err != nil {
		return nil, err
	}

	items, err := defaultCollection.SearchItems(0, attributes)
	if err != nil {
		return nil, err
	}

	secretsData, err := sa.secretService.GetSecrets(0, items, sa.secretSessionPath)
	if err != nil {
		return nil, err
	}

	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	var result = make(map[string]string)
	for itemPath, itemSecret := range secretsData {
		itemObj, err := secrets.NewItem(sessionBus, itemPath)
		if err != nil {
			return nil, err
		}

		attributes, _ := itemObj.Attributes().Get(0)
		settingKey := attributes[keyringTagSettingKey]
		if settingKey != "" {
			result[settingKey] = string(itemSecret.Value)
		}
	}
	return result, nil
}

func (sa *SecretAgent) delete(uuid, settingName, settingKey string) error {
	attributes := map[string]string{
		keyringTagConnUUID:    uuid,
		keyringTagSettingName: settingName,
		keyringTagSettingKey:  settingKey,
	}
	defaultCollection, err := sa.getDefaultCollection()
	if err != nil {
		return err
	}

	items, err := defaultCollection.SearchItems(0, attributes)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return nil
	}

	logger.Debugf("delete uuid: %q, setting name: %q, setting key: %q\n",
		uuid, settingName, settingKey)

	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	item := items[0]
	itemObj, err := secrets.NewItem(sessionBus, item)
	if err != nil {
		return err
	}
	_, err = itemObj.Delete(0)
	return err
}

func (sa *SecretAgent) set(label, uuid, settingName, settingKey, value string) error {
	logger.Debugf("set label: %q, uuid: %q, setting name: %q, setting key: %q, value: %q",
		label, uuid, settingName, settingKey, value)
	itemSecret := secrets.Secret{
		Session:     sa.secretSessionPath,
		Value:       []byte(value),
		ContentType: "text/plain",
	}

	properties := map[string]dbus.Variant{
		"org.freedesktop.Secret.Item.Label": dbus.MakeVariant(label),
		"org.freedesktop.Secret.Item.Type":  dbus.MakeVariant("org.freedesktop.Secret.Generic"),
		"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
			keyringTagConnUUID:    uuid,
			keyringTagSettingName: settingName,
			keyringTagSettingKey:  settingKey,
		}),
	}

	defaultCollection, err := sa.getDefaultCollection()
	if err != nil {
		return err
	}
	_, _, err = defaultCollection.CreateItem(0, properties, itemSecret, true)
	return err
}

func (*SecretAgent) GetInterfaceName() string {
	return "org.freedesktop.NetworkManager.SecretAgent"
}

type getSecretsRequest struct {
	ConnId      string   `json:"connId"`
	ConnType    string   `json:"connType"`
	ConnUUID    string   `json:"connUUID"`
	VpnService  string   `json:"vpnService"`
	SettingName string   `json:"settingName"`
	Secrets     []string `json:"secrets"`
	RequestNew  bool     `json:"requestNew"`
}

type getSecretsReply struct {
	Secrets []string `json:"secrets"`
}

func (sa *SecretAgent) askPasswords(connPath dbus.ObjectPath,
	connectionData map[string]map[string]dbus.Variant,
	connUUID, settingName string, settingKeys []string, requestNew bool) (map[string]string, error) {

	logger.Debugf("askPasswords settingName: %v, settingKeys: %v",
		settingName, settingKeys)
	connId, _ := getConnectionDataString(connectionData, "connection", "id")

	connType, _ := getConnectionDataString(connectionData, "connection", "type")

	vpnService, _ := getConnectionDataString(connectionData, "vpn", "service")

	var req getSecretsRequest
	req.ConnId = connId
	req.ConnType = connType
	req.ConnUUID = connUUID
	req.VpnService = vpnService
	req.SettingName = settingName
	req.Secrets = settingKeys
	req.RequestNew = requestNew
	reqJSON, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}
	logger.Debugf("reqJSON: %s", reqJSON)

	cmd := exec.Command(nmSecretDialogBin)
	cmd.Stdin = bytes.NewReader(reqJSON)
	var cmdOutBuf bytes.Buffer
	cmd.Stdout = &cmdOutBuf
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	sa.addSaveSecretsTask(connPath, settingName, cmd.Process)
	err = cmd.Wait()
	sa.removeSaveSecretsTask(connPath, settingName)
	if err != nil {
		return nil, err
	}
	var reply getSecretsReply
	err = json.Unmarshal(cmdOutBuf.Bytes(), &reply)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)

	if len(settingKeys) != len(reply.Secrets) {
		return nil, errors.New("secretAgent.askPasswords: length not equal")
	}

	for i := 0; i < len(settingKeys); i++ {
		result[settingKeys[i]] = reply.Secrets[i]
	}
	return result, nil
}

func (sa *SecretAgent) GetSecrets(connectionData map[string]map[string]dbus.Variant,
	connectionPath dbus.ObjectPath, settingName string, hints []string, flags uint32) (
	secretsData map[string]map[string]dbus.Variant, busErr *dbus.Error) {
	var err error
	secretsData, err = sa.getSecrets(connectionData, connectionPath, settingName, hints, flags)
	if err != nil {
		return nil, dbusutil.ToError(err)
	}

	logger.Debugf("secretsData: %#v", secretsData)
	return secretsData, nil
}

func getSecretFlagsKeyName(key string) string {
	if strings.HasPrefix(key, "wep-key") {
		num, err := strconv.Atoi(string(key[len(key)-1]))
		if err == nil && 0 <= num && num <= 3 {
			// num in range [0,3]
			return "wep-key-flags"
		}
	}
	return key + "-flags"
}

// 根据当前连接设置，找出必要的密码key。
func isMustAsk(data connectionData, settingName, secretKey string) bool {
	mgmt := getSettingWirelessSecurityKeyMgmt(data)
	switch settingName {
	case nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME:
		wepTxKeyIdx := getSettingWirelessSecurityWepTxKeyidx(data)
		switch mgmt {
		case "wpa-psk":
			if secretKey == "psk" {
				return true
			}

		case "none":
			if secretKey == "wep-key0" && wepTxKeyIdx == 0 {
				return true
			} else if secretKey == "wep-key1" && wepTxKeyIdx == 1 {
				return true
			} else if secretKey == "wep-key2" && wepTxKeyIdx == 2 {
				return true
			} else if secretKey == "wep-key3" && wepTxKeyIdx == 3 {
				return true
			}
		}

	case nm.NM_SETTING_802_1X_SETTING_NAME:
		eap := getSetting8021xEap(data)
		var eap0 string
		if len(eap) >= 1 {
			eap0 = eap[0]
		}
		switch eap0 {
		case "md5", "fast", "ttls", "peap", "leap":
			if secretKey == "password" {
				return true
			}
		case "tls":
			if secretKey == "private-key-password" {
				return true
			}
		}

	}

	return false
}

func (sa *SecretAgent) getSecrets(connectionData map[string]map[string]dbus.Variant,
	connectionPath dbus.ObjectPath, settingName string, hints []string, flags uint32) (
	secretsData map[string]map[string]dbus.Variant, err error) {

	logger.Debug("call getSecrets")

	var allowInteraction bool
	var requestNew bool
	if flags&getSecretsFlagAllowInteraction != 0 {
		logger.Debug("allow interaction")
		allowInteraction = true
	}
	if flags&getSecretsFlagRequestNew != 0 {
		logger.Debug("request new")
		requestNew = true
	}

	logger.Debug("connection path:", connectionPath)
	logger.Debug("setting Name:", settingName)
	logger.Debug("hints:", hints)
	logger.Debug("flags:", flags)

	printConnectionData(connectionData)

	connUUID, ok := getConnectionDataString(connectionData, "connection", "uuid")
	if !ok {
		return nil, errors.New("not found connection uuid")
	}

	logger.Debug("uuid:", connUUID)
	secretsData = make(map[string]map[string]dbus.Variant)
	setting := make(map[string]dbus.Variant)
	secretsData[settingName] = setting

	if settingName == "vpn" {
		vpnSecretsData := make(map[string]string)
		if getSettingVpnServiceType(connectionData) == nmOpenConnectServiceType {
			vpnSecretsData, ok = <- sa.createPendingKey(connectionData, hints, flags)
			if !ok {
				return nil, errors.New("failed to createPendingKey")
			}
		} else {
			vpnSecretsData, _ = getConnectionDataMapStrStr(connectionData, "vpn",
				"secrets")
			vpnDataMap, _ := getConnectionDataMapStrStr(connectionData, "vpn",
				"data")

			var askItems []string
			for _, secretKey := range vpnSecretKeys {
				secretFlag := vpnDataMap[getSecretFlagsKeyName(secretKey)]
				if secretFlag == secretFlagAskStr {
					logger.Debug("ask for password", settingName, secretKey)
					askItems = append(askItems, secretKey)
				}
			}

			if allowInteraction && len(askItems) > 0 {
				resultAsk, err := sa.askPasswords(connectionPath, connectionData, connUUID,
					settingName, askItems, requestNew)
				if err != nil {
					logger.Debug("waring askPasswords error:", err)
					return nil, err
				}
				for key, value := range resultAsk {
					vpnSecretsData[key] = value
				}
			}
		}

		resultSaved, err := sa.getAll(connUUID, settingName)
		if err != nil {
			return nil, err
		}
		logger.Debug("getAll resultSaved:", resultSaved)

		for key, value := range resultSaved {
			if _, ok := vpnSecretsData[key]; !ok {
				vpnSecretsData[key] = value
			} else {
				logger.Debug("not override key", key)
			}
		}

		setting["secrets"] = dbus.MakeVariant(vpnSecretsData)
	} else if secretKeys, ok := secretSettingKeys[settingName]; ok {

		var askItems []string
		for _, secretKey := range secretKeys {
			secretFlags, _ := getConnectionDataUint32(connectionData, settingName,
				getSecretFlagsKeyName(secretKey))

			if secretFlags == secretFlagAsk {
				if allowInteraction && isMustAsk(connectionData, settingName, secretKey) {
					askItems = append(askItems, secretKey)
				}
			} else if secretFlags == secretFlagNone {
				secretStr, _ := getConnectionDataString(connectionData, settingName,
					secretKey)

				if requestNew {
					secretStr = ""
				}

				if secretStr != "" {
					setting[secretKey] = dbus.MakeVariant(secretStr)
				} else if allowInteraction &&
					isMustAsk(connectionData, settingName, secretKey) {
					askItems = append(askItems, secretKey)
				}
			}
		}

		if allowInteraction && len(askItems) > 0 {
			resultAsk, err := sa.askPasswords(connectionPath, connectionData, connUUID,
				settingName, askItems, requestNew)
			if err != nil {
				logger.Warning("askPasswords error:", err)
				if sa.m.ActiveConnectSettingPath == connectionPath {
					sa.m.DisconnectDevice(sa.m.ActiveConnectDevpath)
				}
			} else {
				for key, value := range resultAsk {
					setting[key] = dbus.MakeVariant(value)
				}
			}
		}

		resultSaved, err := sa.getAll(connUUID, settingName)
		if err != nil {
			return nil, err
		}
		logger.Debugf("getAll resultSaved: %#v", resultSaved)

		for key, value := range resultSaved {
			secretFlags, _ := getConnectionDataUint32(connectionData, settingName,
				getSecretFlagsKeyName(key))
			if secretFlags == secretFlagAgentOwned {
				setting[key] = dbus.MakeVariant(value)
			}
		}
	}
	return
}

func printConnectionData(data map[string]map[string]dbus.Variant) {
	for settingName, setting := range data {
		for key, value := range setting {
			logger.Debugf("> %s.%s: %v", settingName, key, value)
		}
	}
}

func (sa *SecretAgent) CancelGetSecrets(connectionPath dbus.ObjectPath, settingName string) *dbus.Error {
	logger.Debug("call CancelGetSecrets")

	logger.Debug("connection path:", connectionPath)
	logger.Debug("setting name:", settingName)

	process := sa.getSaveSecretsTaskProcess(connectionPath, settingName)
	if process != nil {
		logger.Debug("kill process", process.Pid)
		process.Kill()
	}

	return nil
}

func (a *SecretAgent) createPendingKey(connectionData map[string]map[string]dbus.Variant, hints []string, flags uint32) chan map[string]string {
	ch := make(chan map[string]string)

	// for vpn connections, ask password for vpn auth dialogs
	vpnAuthDilogBin := getVpnAuthDialogBin(connectionData)
	go func() {
		args := []string{
			"-u", getSettingConnectionUuid(connectionData),
			"-n", getSettingConnectionId(connectionData),
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
		// process, stdin, stdout, _, err := execWithIO(vpnAuthDilogBin, args...)
		_, stdin, stdout, _, err := execWithIO(vpnAuthDilogBin, args...)
		if err != nil {
			logger.Warning("failed to run vpn auth dialog", err)
			close(ch)
			return
		}

		stdinWriter := bufio.NewWriter(stdin)
		stdoutReader := bufio.NewReader(stdout)

		vpnData := getSettingVpnData(connectionData)
		vpnSecretData := getSettingVpnSecrets(connectionData)

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

		newVpnSecretData := make(map[string]string)
		lastKey := ""
		// read output until there are two empty lines printed
		emptyLines := 0
		for {
			lineBytes, _, err := stdoutReader.ReadLine()
			if err != nil {
				break
			}
			line := string(lineBytes)

			if len(line) == 0 {
				emptyLines++
			} else {
				// the secrets key and value are split as line
				if len(lastKey) == 0 {
					lastKey = line
				} else {
					newVpnSecretData[lastKey] = line
					lastKey = ""
				}
			}
			if emptyLines == 2 {
				break
			}
		}

		// notify auth dialog to quit
		stdinWriter.WriteString("QUIT\n\n")
		err = stdinWriter.Flush()
		if err == nil {
			ch <- newVpnSecretData
		} else {
			logger.Warning("failed to flush auth dialog data", err)
			close(ch)
		}
	}()

	return ch
}

type settingItem struct {
	settingName string
	settingKey  string
	value       string
	label       string
}

func getConnectionDataVariant(connectionData map[string]map[string]dbus.Variant,
	settingName, settingKey string) (dbus.Variant, bool) {

	setting, ok := connectionData[settingName]
	if !ok {
		return dbus.Variant{}, false
	}
	value, ok := setting[settingKey]
	if !ok {
		return dbus.Variant{}, false
	}
	return value, true
}

func getConnectionData(connectionData map[string]map[string]dbus.Variant,
	settingName, settingKey string) (interface{}, bool) {

	variant, ok := getConnectionDataVariant(connectionData, settingName, settingKey)
	if !ok {
		return nil, false
	}
	return variant.Value(), true
}

func getConnectionDataString(connectionData map[string]map[string]dbus.Variant,
	settingName, settingKey string) (string, bool) {
	val, ok := getConnectionData(connectionData, settingName, settingKey)
	if ok {
		valStr, ok := val.(string)
		if ok {
			return valStr, true
		}
	}
	return "", false
}

func getConnectionDataMapStrStr(connectionData map[string]map[string]dbus.Variant,
	settingName, settingKey string) (map[string]string, bool) {

	val, ok := getConnectionData(connectionData, settingName, settingKey)
	if ok {
		valMap, ok := val.(map[string]string)
		if ok {
			return valMap, true
		}
	}
	return nil, false
}

func getConnectionDataUint32(connectionData map[string]map[string]dbus.Variant,
	settingName, settingKey string) (uint32, bool) {

	val, ok := getConnectionData(connectionData, settingName, settingKey)
	if ok {
		valUint, ok := val.(uint32)
		if ok {
			return valUint, true
		}
	}
	return 0, false
}

var secretSettingKeys = map[string][]string{
	"802-11-wireless-security": {"psk", "wep-key0", "wep-key1", "wep-key2", "wep-key3",
		"leap-password"},
	"802-1x": {"password", "password-raw", "ca-cert-password",
		"client-cert-password", "phase2-ca-cert-password", "phase2-client-cert-password",
		"private-key-password", "phase2-private-key-password", "pin"},
	// temporarily not supported password-raw
	"pppoe": {"password"},
	"gsm":   {"password", "pin"},
	"cdma":  {"password"},
}

var vpnSecretKeys = []string{
	"password", "proxy-password", "IPSec secret", "Xauth password",
}

func (sa *SecretAgent) SaveSecretsDeepin(connectionData map[string]map[string]dbus.Variant,
	connectionPath dbus.ObjectPath) *dbus.Error {
	err := sa.saveSecrets(connectionData, connectionPath)
	return dbusutil.ToError(err)
}

func (sa *SecretAgent) SaveSecrets(connectionData map[string]map[string]dbus.Variant,
	connectionPath dbus.ObjectPath) *dbus.Error {
	err := sa.saveSecrets(connectionData, connectionPath)
	return dbusutil.ToError(err)
}

func (sa *SecretAgent) saveSecrets(connectionData map[string]map[string]dbus.Variant,
	connectionPath dbus.ObjectPath) error {
	logger.Debug("call saveSecrets")

	printConnectionData(connectionData)
	logger.Debug("connection path:", connectionPath)

	connUUID, ok := getConnectionDataString(connectionData, "connection",
		"uuid")
	if !ok {
		return dbusutil.ToError(errors.New("not found connection uuid"))
	}
	logger.Debug("uuid:", connUUID)
	connId, _ := getConnectionDataString(connectionData, "connection", "id")
	logger.Debug("conn id:", connId)

	vpnServiceType, _ := getConnectionDataString(connectionData, "vpn", "service-type")
	dotLastIdx := strings.LastIndex(vpnServiceType, ".")
	if dotLastIdx != -1 {
		vpnServiceType = vpnServiceType[dotLastIdx+1:]
	}

	var arr []settingItem

	for settingName, setting := range connectionData {

		if settingName == "vpn" {
			var vpnDataMap map[string]string
			vpnData, ok := setting["data"]
			if ok {
				vpnDataMap, _ = vpnData.Value().(map[string]string)
				logger.Debug("vpn.data map:", vpnDataMap)
			}

			secret, ok := setting["secrets"]
			if ok {
				logger.Debug("vpn.secret value:", secret)
				secretMap, ok := secret.Value().(map[string]string)
				if ok {
					for key, value := range secretMap {
						secretFlags := vpnDataMap[getSecretFlagsKeyName(key)]

						if secretFlags == secretFlagAgentOwnedStr {
							label := fmt.Sprintf("VPN password secret for %s/%s/%s",
								connId, vpnServiceType, key)
							arr = append(arr, settingItem{
								settingName: settingName,
								settingKey:  key,
								value:       value,
								label:       label,
							})
						}
					}
				}
			}
			continue
		}

		secretKeys := secretSettingKeys[settingName]
		for key, value := range setting {
			if strv.Strv(secretKeys).Contains(key) {
				// key is secret key

				secretFlags, _ := getConnectionDataUint32(connectionData,
					settingName, getSecretFlagsKeyName(key))
				if secretFlags != secretFlagAgentOwned {
					// not agent owned
					continue
				}

				valueStr, ok := value.Value().(string)
				if ok {
					arr = append(arr, settingItem{
						settingName: settingName,
						settingKey:  key,
						value:       valueStr,
					})
				}
			}
		}
	}

	for _, item := range arr {
		label := item.label
		if label == "" {
			label = fmt.Sprintf("Network secret for %s/%s/%s", connId,
				item.settingName, item.settingKey)
		}
		sa.set(label, connUUID, item.settingName, item.settingKey, item.value)
	}

	// delete
	for settingName, secretKeys := range secretSettingKeys {
		for _, secretKey := range secretKeys {
			secretFlags, _ := getConnectionDataUint32(connectionData, settingName,
				getSecretFlagsKeyName(secretKey))

			if secretFlags != secretFlagAgentOwned {
				sa.delete(connUUID, settingName, secretKey)
			}
		}
	}

	vpnData, ok := getConnectionData(connectionData, "vpn", "data")
	if ok {
		vpnDataMap, ok := vpnData.(map[string]string)
		if ok {
			for _, secretKey := range vpnSecretKeys {
				secretFlags := vpnDataMap[getSecretFlagsKeyName(secretKey)]
				if secretFlags != secretFlagAgentOwnedStr {
					sa.delete(connUUID, "vpn", secretKey)
				}
			}
		}
	}

	return nil
}

func (sa *SecretAgent) DeleteSecrets(connectionData map[string]map[string]dbus.Variant,
	connectionPath dbus.ObjectPath) *dbus.Error {
	logger.Debug("call DeleteSecrets")
	printConnectionData(connectionData)

	connUUID, ok := getConnectionDataString(connectionData, "connection",
		"uuid")
	if !ok {
		return dbusutil.ToError(errors.New("not found connection uuid"))
	}

	err := sa.deleteAll(connUUID)
	return dbusutil.ToError(err)
}

func (*SecretAgentSession) GetSystemBusName() (string, *dbus.Error) {
	sysBus, err := dbus.SystemBus()
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	return sysBus.Names()[0], nil
}

type SecretAgentSession struct {
	methods *struct {
		GetSystemBusName func() `out:"name"`
	}
}

func (*SecretAgentSession) GetInterfaceName() string {
	return "com.deepin.daemon.Network.SecretAgent"
}
