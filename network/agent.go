package main

import "dlib/dbus"
import nm "dbus/org/freedesktop/networkmanager"

import "fmt"

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
	fmt.Println("GetSecrtes:", connectionPath, settingName, hints, flags)
	r := make(map[string]map[string]dbus.Variant)
	if uuid, ok := connection[fieldConnection]["uuid"].Value().(string); ok {
		r["802-11-wireless-security"] = make(map[string]dbus.Variant)
		_Manager.NeedMoreConfigure(uuid, "PASSWORD")
		fmt.Println("Begin getsecrtes...", connection[fieldConnection]["id"], uuid)
		if _, ok := a.keys[uuid]; !ok {
			//only wait the key when there is no other GetSecrtes runing with this uuid
			a.keys[uuid] = make(chan string)
			key, ok := <-a.keys[uuid]
			if ok {
				r["802-11-wireless-security"]["psk"] = dbus.MakeVariant(key)
			}
			fmt.Println("End getsecrtes...", connection[fieldConnection]["id"])
		}
	}
	return r
}

func (a *Agent) CancelGetSecrtes(connectionPath dbus.ObjectPath, settingName string) {
	uuid, ok := nm.GetSettingsConnection(string(connectionPath)).GetSettings()[fieldConnection]["uuid"].Value().(string)
	if ok {
		close(a.keys[uuid])
		delete(a.keys, uuid)
	}
	fmt.Println("CancelGetSecrtes")
}

func (a *Agent) SaveSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath) {
	fmt.Println("SaveSecretes")
}

func (a *Agent) DeleteSecrets(connection map[string]map[string]dbus.Variant, connectionPath dbus.ObjectPath) {
	fmt.Println("DeleteSecretes")
}

func (a *Agent) handleNewKeys(uuid string, key string) {
	ch, ok := a.keys[uuid]
	if ok {
		//ignore handle key request when there is no GetSecrets waiting.
		ch <- key
		delete(a.keys, uuid)
	}
}

func NewAgent(identify string) *Agent {
	c := &Agent{}
	c.keys = make(map[string]chan string)
	dbus.InstallOnSystem(c)
	nm.GetAgentManager("/org/freedesktop/NetworkManager/AgentManager").Register("org.snyh.test")
	return c
}

func (m *Manager) SetKey(id string, key string) {
	m.agent.handleNewKeys(id, key)
}
