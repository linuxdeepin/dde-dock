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
	DeviceStateChanged           func(devPath dbus.ObjectPath, new_state uint32)
	AccessPointAdded             func(devPath dbus.ObjectPath, ap AccessPoint)
	AccessPointRemoved           func(devPath dbus.ObjectPath, apPath dbus.ObjectPath)
	AccessPointPropertiesChanged func(devPath dbus.ObjectPath, ap AccessPoint) // TODO

	agent *Agent
}

func (this *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{DBusDest, DBusPath, DBusIFC}
}

func _NewManager() (m *Manager) {
	this := &Manager{}
	return this
}

func (this *Manager) initManager() {
	this.WiredEnabled = true
	this.WirelessEnabled = property.NewWrapProperty(this, "WirelessEnabled", _NMManager.WirelessEnabled)
	this.NetworkingEnabled = property.NewWrapProperty(this, "NetworkingEnabled", _NMManager.NetworkingEnabled)

	this.initDeviceManage()
	this.initConnectionManage()

	// update property "State"
	this.updatePropState()
	_NMManager.State.ConnectChanged(func() {
		this.updatePropState()
	})

	// update property "ActiveConnections" after initilizing device
	this.updatePropActiveConnections()
	_NMManager.ActiveConnections.ConnectChanged(func() {
		this.updatePropActiveConnections()
	})

	this.agent = newAgent("org.snyh.agent")
}

func (this *Manager) updatePropActiveConnections() {
	this.ActiveConnections = make([]string, 0)
	for _, cpath := range _NMManager.ActiveConnections.Get() {
		if conn, err := nm.NewActiveConnection(NMDest, cpath); err == nil {
			this.ActiveConnections = append(this.ActiveConnections, conn.Uuid.Get())
		}
	}
	dbus.NotifyChange(this, "ActiveConnections")
}

func (this *Manager) updatePropState() {
	this.State = _NMManager.State.Get()
	dbus.NotifyChange(this, "State")
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
