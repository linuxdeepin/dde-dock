package main

import "dlib/dbus"
import "time"
import nm "dbus/org/freedesktop/networkmanager"

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
		r[settingName]["psk"] = dbus.MakeVariant(key)
	default:
		LOGGER.Warning("Unknow secrety setting name", settingName, ",please report it to linuxdeepin")
	}
	return r
}

func (a *Agent) GetSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath, settingName string, hints []string, flags uint32) map[string]map[string]dbus.Variant {
	LOGGER.Info("GetSecrtes:", connectionPath, settingName, hints, flags)
	keyId := mapKey{connectionPath, settingName}
	if keyValue, ok := a.savedKeys[keyId]; ok && (flags&NM_SECRET_AGENT_GET_SECRETS_FLAG_REQUEST_NEW == 0) {
		return keyValue
	}

	if flags&NM_SECRET_AGENT_GET_SECRETS_FLAG_ALLOW_INTERACTION == 0 {
		return invalidKey
	}

	if _, ok := a.pendingKeys[keyId]; ok {
		//only wait the key when there is no other GetSecrtes runing with this uuid
		LOGGER.Info("Repeat GetSecrets", keyId)
	} else {
		select {
		case keyValue, ok := <-a.createPendingKey(keyId, pageGeneralGetId(connection)):
			if ok {
				keyValue := fillSecret(settingName, keyValue)
				a.SaveSecrets(keyValue, connectionPath)
				return keyValue
			}
			LOGGER.Info("failed getsecrtes...", keyId)
		case <-time.After(120 * time.Second):
			a.CancelGetSecrtes(connectionPath, settingName)
			LOGGER.Info("get secrets timeout:", keyId)
		}
	}
	return invalidKey
}
func (a *Agent) createPendingKey(keyId mapKey, connectionId string) chan string {
	if _Manager.NeedSecrets != nil {
		defer _Manager.NeedSecrets(keyId.path, keyId.name, connectionId)
	} else {
		LOGGER.Warning("createPendingKey when DNetworkManager hasn't init")
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
		LOGGER.Warning("CancelGetSecrtes an unknow PendingKey:", keyId)
	}
}

func (a *Agent) SaveSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath) {
	if _, ok := connection["802-11-wireless-security"]; ok {
		keyId := mapKey{connectionPath, "802-11-wireless-security"}
		a.savedKeys[keyId] = connection
	}
	LOGGER.Info("SaveSecretes")
}

func (a *Agent) DeleteSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath) {
	if _, ok := connection["802-11-wireless-security"]; ok {
		keyId := mapKey{connectionPath, "802-11-wireless-security"}
		delete(a.savedKeys, keyId)
	}
	LOGGER.Info("DeleteSecretes")
}

func (a *Agent) feedSecret(path dbus.ObjectPath, name string, key string) {
	keyId := mapKey{path, name}
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

	if manager, err := nm.NewAgentManager(NMDest, "/org/freedesktop/NetworkManager/AgentManager"); err != nil {
		panic(err)
	} else {
		manager.Register("com.deepin.daemon.Network.Agent")
	}
	return c
}

func (m *Manager) FeedSecret(path dbus.ObjectPath, name, key string) {
	m.agent.feedSecret(path, name, key)
}
func (m *Manager) CancelSecret(path dbus.ObjectPath, name string) {
	m.agent.CancelGetSecrtes(path, name)
}
