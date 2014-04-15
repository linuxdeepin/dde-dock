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

func (c *Agent) GetDBusInfo() dbus.DBusInfo {
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
		Logger.Warning("Unknow secrety setting name", settingName, ",please report it to linuxdeepin")
	}
	return r
}

func (a *Agent) GetSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath, settingName string, hints []string, flags uint32) map[string]map[string]dbus.Variant {
	Logger.Info("GetSecrets:", connectionPath, settingName, hints, flags)
	keyId := mapKey{connectionPath, settingName}

	// TODO fixme
	// if keyValue, ok := a.savedKeys[keyId]; ok && (flags&NM_SECRET_AGENT_GET_SECRETS_FLAG_REQUEST_NEW == 0) {
	// Logger.Debug("GetSecrets return ", keyValue) // TODO test
	// return keyValue
	// }
	// if flags&NM_SECRET_AGENT_GET_SECRETS_FLAG_USER_REQUESTED == 0 {
	// Logger.Debug("GetSecrets return") // TODO test
	// return invalidKey
	// }

	if flags&NM_SECRET_AGENT_GET_SECRETS_FLAG_ALLOW_INTERACTION == 0 {
		return invalidKey
	}

	// if _, ok := a.pendingKeys[keyId]; ok {
	// 	//only wait the key when there is no other GetSecrtes runing with this uuid
	// 	Logger.Info("Repeat GetSecrets", keyId)
	// } else {
	select {
	// TODO
	case keyValue, ok := <-a.createPendingKey(keyId, pageGeneralGetId(connection)):
		if ok {
			keyValue := fillSecret(settingName, keyValue)
			a.SaveSecrets(keyValue, connectionPath)
			return keyValue
		}
		Logger.Info("failed getsecrtes...", keyId)
	case <-time.After(120 * time.Second):
		a.CancelGetSecrtes(connectionPath, settingName)
		Logger.Info("get secrets timeout:", keyId)
	}
	// }
	return invalidKey
}
func (a *Agent) createPendingKey(keyId mapKey, connectionId string) chan string {
	Logger.Debug("createPendingKey:", keyId, connectionId) // TODO test
	if manager.NeedSecrets != nil {
		defer manager.NeedSecrets(string(keyId.path), keyId.name, connectionId)
	} else {
		Logger.Warning("createPendingKey when DNetworkManager hasn't init")
	}
	a.pendingKeys[keyId] = make(chan string)
	return a.pendingKeys[keyId]
}

func (a *Agent) CancelGetSecrtes(connectionPath dbus.ObjectPath, settingName string) {
	keyId := mapKey{connectionPath, settingName}

	if pendingChan, ok := a.pendingKeys[keyId]; ok {
		close(pendingChan)
		delete(a.pendingKeys, keyId)
	} else {
		Logger.Warning("CancelGetSecrtes an unknow PendingKey:", keyId)
	}
	Logger.Info("CancelGetSecrtes")
}

func (a *Agent) SaveSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath) {
	// TODO fixme
	// if _, ok := connection["802-11-wireless-security"]; ok {
	// keyId := mapKey{connectionPath, "802-11-wireless-security"}
	// a.savedKeys[keyId] = connection
	// Logger.Debug("SaveSecrets:", connection, connectionPath) // TODO test
	// }
	Logger.Info("SaveSecretes")
}

func (a *Agent) DeleteSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath) {
	if _, ok := connection["802-11-wireless-security"]; ok {
		keyId := mapKey{connectionPath, "802-11-wireless-security"}
		delete(a.savedKeys, keyId)
	}
	Logger.Info("DeleteSecretes")
}

func (a *Agent) feedSecret(path string, name string, key string) {
	keyId := mapKey{dbus.ObjectPath(path), name}
	if ch, ok := a.pendingKeys[keyId]; ok {
		ch <- key
		delete(a.pendingKeys, keyId)
	}
}

func newAgent(identify string) *Agent {
	c := &Agent{}
	c.pendingKeys = make(map[mapKey]chan string)
	c.savedKeys = make(map[mapKey]map[string]map[string]dbus.Variant)

	dbus.InstallOnSystem(c)

	if manager, err := nmNewAgentManager(); err != nil {
		panic(err)
	} else {
		manager.Register("com.deepin.daemon.Network.Agent")
	}
	return c
}

func (m *Manager) FeedSecret(path string, name, key string) {
	m.agent.feedSecret(path, name, key)
}
func (m *Manager) CancelSecret(path string, name string) {
	m.agent.CancelGetSecrtes(dbus.ObjectPath(path), name)
}
