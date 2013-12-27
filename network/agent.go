package main

import "dlib/dbus"
import nm "dbus/org/freedesktop/networkmanager"
import "log"

type Agent struct {
	keys map[string]chan string
}

func (c *Agent) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		".",
		"/org/freedesktop/NetworkManager/SecretAgent",
		"org.freedesktop.NetworkManager.SecretAgent",
	}
}

func (a *Agent) GetSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath, settingName string, hints []string, flags uint32) map[string]map[string]dbus.Variant {
	log.Println("GetSecrtes:", connectionPath, settingName, hints, flags)
	r := make(map[string]map[string]dbus.Variant)
	if uuid, ok := connection[fieldConnection]["uuid"].Value().(string); ok {
		r["802-11-wireless-security"] = make(map[string]dbus.Variant)
		_Manager.NeedMoreConfigure(uuid, "PASSWORD")
		log.Println("Begin getsecrtes...", connection[fieldConnection]["id"], uuid)
		if _, ok := a.keys[uuid]; !ok {
			//only wait the key when there is no other GetSecrtes runing with this uuid
			a.keys[uuid] = make(chan string)
			key, ok := <-a.keys[uuid]
			if ok {
				r["802-11-wireless-security"]["psk"] = dbus.MakeVariant(key)
			}
			log.Println("End getsecrtes...", connection[fieldConnection]["id"])
		}
	}
	return r
}

func (a *Agent) CancelGetSecrtes(connectionPath dbus.ObjectPath, settingName string) {
	if setting, err := nm.NewSettingsConnection(connectionPath); err == nil {
		s, err := setting.GetSettings()
		if err == nil {
			uuid, ok := s[fieldConnection]["uuid"].Value().(string)
			if ok {
				close(a.keys[uuid])
				delete(a.keys, uuid)
			}
		} else {
			log.Println("CancelGetSecrtes Error:", err)
		}
	} else {
		log.Println("CancelGetSecrtes Failed:", err)
	}
}

func (a *Agent) SaveSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath) {
	log.Println("SaveSecretes")
}

func (a *Agent) DeleteSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath) {
	log.Println("DeleteSecretes")
}

func (a *Agent) handleNewKeys(uuid string, key string) {
	ch, ok := a.keys[uuid]
	if ok {
		//ignore handle key request when there is no GetSecrets waiting.
		ch <- key
		delete(a.keys, uuid)
	}
}

func newAgent(identify string) *Agent {
	c := &Agent{}
	c.keys = make(map[string]chan string)
	dbus.InstallOnSystem(c)
	manager, err := nm.NewAgentManager("/org/freedesktop/NetworkManager/AgentManager")
	if err != nil {
		panic(err)
	}
	manager.Register("org.snyh.test")
	return c
}

func (m *Manager) SetKey(id string, key string) {
	m.agent.handleNewKeys(id, key)
}
