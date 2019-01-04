package network

import (
	"encoding/json"
	"pkg.deepin.io/dde/daemon/common/dsync"
	"pkg.deepin.io/lib/dbus1"
)

type syncConfig struct {
	m *Manager
}

const (
	daemonSysService = "com.deepin.daemon.Daemon"
	daemonSysPath    = "/com/deepin/daemon/Daemon"
	daemonSysIFC     = daemonSysService

	methodSysNetGetConnections = daemonSysIFC + ".NetworkGetConnections"
	methodSysNetSetConnections = daemonSysIFC + ".NetworkSetConnections"
)

func (sc *syncConfig) Get() (interface{}, error) {
	obj, err := getDaemonSysBus()
	if err != nil {
		return nil, err
	}
	var data []byte
	err = obj.Call(methodSysNetGetConnections, 0).Store(&data)
	if err != nil {
		return nil, err
	}
	var info dsync.NetworkData
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (sc *syncConfig) Set(data []byte) error {
	obj, err := getDaemonSysBus()
	if err != nil {
		return err
	}
	return obj.Call(methodSysNetSetConnections, 0, data).Store()
}

func getDaemonSysBus() (dbus.BusObject, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	return conn.Object(daemonSysService, daemonSysPath), nil
}
