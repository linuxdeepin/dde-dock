package main

import "dlib/dbus"
import "time"

type mapKey struct {
	path dbus.ObjectPath
	name string
}
type Agent struct {
	pendingKeys map[mapKey]chan string

	savedKeys map[mapKey]map[string]map[string]dbus.Variant
}

var (
	invalidKey = make(map[string]map[string]dbus.Variant)
)

func (a *Agent) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		".",
		"/org/freedesktop/NetworkManager/SecretAgent",
		"org.freedesktop.NetworkManager.SecretAgent",
	}
}

func fillSecret(settingName string, key string) map[string]map[string]dbus.Variant {
	r := make(map[string]map[string]dbus.Variant)
	r[settingName] = make(map[string]dbus.Variant)
	switch settingName {
	case "802-11-wireless-security":
		r[settingName]["psk"] = dbus.MakeVariant(key) // TODO
	default:
		logger.Warning("Unknow secrety setting name", settingName, ",please report it to linuxdeepin")
	}
	return r
}

func (a *Agent) GetSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath, settingName string, hints []string, flags uint32) map[string]map[string]dbus.Variant {
	logger.Debug("GetSecrets:", connectionPath, settingName, hints, flags)
	keyId := mapKey{connectionPath, settingName}

	// TODO fixme
	// if keyValue, ok := a.savedKeys[keyId]; ok && (flags&NM_SECRET_AGENT_GET_SECRETS_FLAG_REQUEST_NEW == 0) {
	// logger.Debug("GetSecrets return ", keyValue) // TODO test
	// return keyValue
	// }
	// if flags&NM_SECRET_AGENT_GET_SECRETS_FLAG_USER_REQUESTED == 0 {
	// logger.Debug("GetSecrets return") // TODO test
	// return invalidKey
	// }

	if flags&NM_SECRET_AGENT_GET_SECRETS_FLAG_ALLOW_INTERACTION == 0 &&
		flags&NM_SECRET_AGENT_GET_SECRETS_FLAG_USER_REQUESTED == 0 {
		logger.Debug("GetSecrets: invalid key", flags)
		return invalidKey
	}

	// if _, ok := a.pendingKeys[keyId]; ok {
	// 	//only wait the key when there is no other GetSecrtes runing with this uuid
	// 	logger.Info("Repeat GetSecrets", keyId)
	// } else {
	select {
	case keyValue, ok := <-a.createPendingKey(keyId, getSettingConnectionId(connection)):
		if ok {
			keyValue := fillSecret(settingName, keyValue)
			a.SaveSecrets(keyValue, connectionPath)
			return keyValue
		}
		logger.Info("failed getsecrtes...", keyId)
	case <-time.After(120 * time.Second):
		a.CancelGetSecrtes(connectionPath, settingName)
		logger.Info("get secrets timeout:", keyId)
	}
	// }
	return invalidKey
}
func (a *Agent) createPendingKey(keyId mapKey, connectionId string) chan string {
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

func (a *Agent) CancelGetSecrtes(connectionPath dbus.ObjectPath, settingName string) {
	logger.Debug("CancelGetSecrtes:", connectionPath, settingName)
	keyId := mapKey{connectionPath, settingName}

	if pendingChan, ok := a.pendingKeys[keyId]; ok {
		close(pendingChan)
		delete(a.pendingKeys, keyId)
	} else {
		logger.Warning("CancelGetSecrtes an unknow PendingKey:", keyId)
	}
}

func (a *Agent) SaveSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath) {
	logger.Debug("SaveSecretes:", connectionPath)
	// TODO fixme
	// if _, ok := connection["802-11-wireless-security"]; ok {
	// keyId := mapKey{connectionPath, "802-11-wireless-security"}
	// a.savedKeys[keyId] = connection
	// logger.Debug("SaveSecrets:", connection, connectionPath) // TODO test
	// }
}

func (a *Agent) DeleteSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath) {
	logger.Debug("DeleteSecrets:", connectionPath) // TODO test
	if _, ok := connection["802-11-wireless-security"]; ok {
		keyId := mapKey{connectionPath, "802-11-wireless-security"}
		delete(a.savedKeys, keyId)
	}
}

func (a *Agent) feedSecret(path string, name string, key string) {
	keyId := mapKey{dbus.ObjectPath(path), name}
	if ch, ok := a.pendingKeys[keyId]; ok {
		ch <- key
		delete(a.pendingKeys, keyId)
	}
}

func newAgent() (a *Agent) {
	a = &Agent{}
	a.pendingKeys = make(map[mapKey]chan string)
	a.savedKeys = make(map[mapKey]map[string]map[string]dbus.Variant)

	err := dbus.InstallOnSystem(a)
	if err != nil {
		logger.Error("install network agent failed:", err)
		return
	}

	nmAgentRegister("com.deepin.daemon.Network.Agent")
	return
}

func destroyAgent(a *Agent) {
	nmAgentUnregister()
	dbus.UnInstallObject(a)
}

func (m *Manager) FeedSecret(path string, name, key string) {
	logger.Debug("FeedSecret:", path, name, key)
	m.agent.feedSecret(path, name, key)
}
func (m *Manager) CancelSecret(path string, name string) {
	logger.Debug("CancelSecret:", path, name)
	m.agent.CancelGetSecrtes(dbus.ObjectPath(path), name)
}
