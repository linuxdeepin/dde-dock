package main

import "dlib/dbus"
import "dlib/dbus/property"
import "dlib/logger"
import "os"
import nm "dbus/org/freedesktop/networkmanager"
import "flag"

const (
	DBusDest = "com.deepin.daemon.Network"
	DBusPath = "/com/deepin/daemon/Network"
	DBusIFC  = "com.deepin.daemon.Network"

	NMDest = "org.freedesktop.NetworkManager"
)

const (
	OpAdded = iota
	OpRemoved
)

var (
	_NMManager, _  = nm.NewManager(NMDest, "/org/freedesktop/NetworkManager")
	_NMSettings, _ = nm.NewSettings(NMDest, "/org/freedesktop/NetworkManager/Settings")
	_Manager       *Manager
	LOGGER         = logger.NewLogger("com.deepin.daemon.Network")
	argDebug       bool
)

type Manager struct {
	//update by manager.go
	WiredEnabled      bool          `access:"readwrite"`
	VPNEnabled        bool          `access:"readwrite"`
	WirelessEnabled   dbus.Property `access:"readwrite"`
	NetworkingEnabled dbus.Property `access:"readwrite"`
	ActiveConnections []string      // uuid collection of connections that activated
	State             uint32        // networking state

	//update by devices.go
	WirelessDevices []*Device // TODO is "Device" struct still needed?
	WiredDevices    []*Device
	OtherDevices    []*Device

	//update by connections.go
	WiredConnections    []string
	WirelessConnections []string
	VPNConnections      []string
	uuid2connectionType map[string]string // TODO

	//signals
	NeedSecrets                  func(string, string, string)
	DeviceStateChanged           func(devPath string, newState uint32)
	AccessPointAdded             func(devPath string, apPath string)
	AccessPointRemoved           func(devPath string, apPath string)
	AccessPointPropertiesChanged func(devPath string, apPath string)

	agent *Agent
}

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{DBusDest, DBusPath, DBusIFC}
}

func _NewManager() (m *Manager) {
	m = &Manager{}
	return
}

func (m *Manager) initManager() {
	m.WiredEnabled = true
	m.WirelessEnabled = property.NewWrapProperty(m, "WirelessEnabled", _NMManager.WirelessEnabled)
	m.NetworkingEnabled = property.NewWrapProperty(m, "NetworkingEnabled", _NMManager.NetworkingEnabled)

	m.initDeviceManage()
	m.initConnectionManage()

	// update property "ActiveConnections" after initilizing device
	m.updatePropActiveConnections()
	_NMManager.ActiveConnections.ConnectChanged(func() {
		m.updatePropActiveConnections()
	})

	// update property "State"
	m.updatePropState()
	_NMManager.State.ConnectChanged(func() {
		m.updatePropState()
	})

	m.agent = newAgent("org.snyh.agent")
}

func (m *Manager) updatePropActiveConnections() {
	m.ActiveConnections = make([]string, 0)
	for _, cpath := range _NMManager.ActiveConnections.Get() {
		if conn, err := nm.NewActiveConnection(NMDest, cpath); err == nil {
			m.ActiveConnections = append(m.ActiveConnections, conn.Uuid.Get())
			LOGGER.Debugf("ActiveConnections, uuid=%s, state=%d", conn.Uuid.Get(), conn.State.Get()) // TODO test
		}
	}
	dbus.NotifyChange(m, "ActiveConnections")
}

func (m *Manager) updatePropState() {
	m.State = _NMManager.State.Get()
	dbus.NotifyChange(m, "State")
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			LOGGER.Fatal(err)
		}
	}()

	// configure logger
	flag.BoolVar(&argDebug, "d", false, "debug mode")
	flag.BoolVar(&argDebug, "debug", false, "debug mode")
	flag.Parse()
	if argDebug {
		LOGGER.SetLogLevel(logger.LEVEL_DEBUG)
	}

	_Manager = _NewManager()
	err := dbus.InstallOnSession(_Manager)
	if err != nil {
		LOGGER.Error("register dbus interface failed: ", err)
		os.Exit(1)
	}

	// initialize manager after configuring dbus
	_Manager.initManager()

	dbus.DealWithUnhandledMessage()
	if err := dbus.Wait(); err != nil {
		LOGGER.Error("lost dbus session: ", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
