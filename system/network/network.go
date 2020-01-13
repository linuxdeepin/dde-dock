package network

import (
	"errors"
	"strings"
	"sync"

	networkmanager "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.networkmanager"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/dde/daemon/network/nm"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
	"pkg.deepin.io/lib/log"
)

const (
	dbusServiceName = "com.deepin.system.Network"
	dbusPath        = "/com/deepin/system/Network"
	dbusInterface   = dbusServiceName
)

type Module struct {
	*loader.ModuleBase
	network *Network
}

func (m *Module) GetDependencies() []string {
	return nil
}

func (m *Module) Start() error {
	if m.network != nil {
		return nil
	}
	logger.Debug("start network")
	m.network = newNetwork()

	service := loader.GetService()
	m.network.service = service

	err := m.network.init()
	if err != nil {
		return err
	}

	serverObj, err := service.NewServerObject(dbusPath, m.network)
	if err != nil {
		return err
	}

	err = serverObj.SetWriteCallback(m.network, "VpnEnabled", m.network.vpnEnabledWriteCb)
	if err != nil {
		return err
	}

	err = serverObj.Export()
	if err != nil {
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		return err
	}

	go func() {
		m.network.PropsMu.RLock()
		enabled := m.network.VpnEnabled
		m.network.PropsMu.RUnlock()

		// auto connect vpn connections
		if enabled {
			m.network.enableVpn(true)
		}
	}()

	return nil
}

func (m *Module) Stop() error {
	// TODO:
	return nil
}

var logger = log.NewLogger("daemon/system/network")

func newModule(logger *log.Logger) *Module {
	m := new(Module)
	m.ModuleBase = loader.NewModuleBase("network", m, logger)
	return m
}

func init() {
	loader.Register(newModule(logger))
}

//go:generate dbusutil-gen -type Network network.go

type Network struct {
	service        *dbusutil.Service
	PropsMu        sync.RWMutex
	VpnEnabled     bool `prop:"access:rw"`
	delayEnableVpn bool
	config         *Config
	configMu       sync.Mutex
	devices        map[dbus.ObjectPath]*device
	devicesMu      sync.Mutex
	nmManager      *networkmanager.Manager
	nmSettings     *networkmanager.Settings
	sigLoop        *dbusutil.SignalLoop
	methods        *struct {
		IsDeviceEnabled       func() `in:"pathOrIface" out:"enabled"`
		EnableDevice          func() `in:"pathOrIface,enabled"`
		Ping                  func() `in:"host"`
		ToggleWirelessEnabled func() `out:"enabled"`
	}

	signals *struct {
		DeviceEnabled struct {
			devPath dbus.ObjectPath
			enabled bool
		}
	}
}

func (n *Network) init() error {
	sysBus := n.service.Conn()
	n.sigLoop = dbusutil.NewSignalLoop(sysBus, 10)
	n.sigLoop.Start()
	n.nmManager = networkmanager.NewManager(sysBus)
	n.nmSettings = networkmanager.NewSettings(sysBus)
	devicePaths, err := n.nmManager.GetDevices(0)
	if err != nil {
		logger.Warning(err)
	} else {
		for _, devPath := range devicePaths {
			err = n.addDevice(devPath)
			if err != nil {
				logger.Warning(err)
				continue
			}
		}

		for iface := range n.config.Devices {
			if n.getDeviceByIface(iface) == nil {
				delete(n.config.Devices, iface)
			}
		}
	}
	n.connectSignal()

	return nil
}

type device struct {
	iface    string
	nmDevice *networkmanager.Device
	type0    uint32
}

func (n *Network) getSysBus() *dbus.Conn {
	return n.service.Conn()
}

func (n *Network) connectSignal() {
	err := dbusutil.NewMatchRuleBuilder().Type("signal").
		PathNamespace("/org/freedesktop/NetworkManager/Devices").
		Interface("org.freedesktop.NetworkManager.Device").
		Member("StateChanged").Build().AddTo(n.getSysBus())
	if err != nil {
		logger.Warning(err)
	}

	err = dbusutil.NewMatchRuleBuilder().Type("signal").
		PathNamespace("/org/freedesktop/NetworkManager/ActiveConnection").
		Interface("org.freedesktop.NetworkManager.VPN.Connection").
		Member("VpnStateChanged").Build().AddTo(n.getSysBus())
	if err != nil {
		logger.Warning(err)
	}

	n.nmManager.InitSignalExt(n.sigLoop, true)
	_, err = n.nmManager.ConnectDeviceAdded(func(devPath dbus.ObjectPath) {
		logger.Debug("device added", devPath)
		n.devicesMu.Lock()

		err := n.addDevice(devPath)
		if err != nil {
			logger.Warning(err)
		}

		n.devicesMu.Unlock()
	})
	if err != nil {
		logger.Warning(err)
	}

	_, err = n.nmManager.ConnectDeviceRemoved(func(devPath dbus.ObjectPath) {
		logger.Debug("device removed", devPath)
		n.devicesMu.Lock()

		n.removeDevice(devPath)

		n.devicesMu.Unlock()
	})
	if err != nil {
		logger.Warning(err)
	}

	_, err = n.nmManager.ConnectStateChanged(func(state uint32) {
		n.PropsMu.RLock()
		delay := n.delayEnableVpn
		n.PropsMu.RUnlock()
		if !delay {
			return
		}

		avail, err := n.isNetworkAvailable()
		if err != nil {
			logger.Warning(err)
			return
		}

		if avail {
			n.PropsMu.Lock()
			n.delayEnableVpn = false
			n.PropsMu.Unlock()
			go n.enableVpn1()
		}

	})
	if err != nil {
		logger.Warning(err)
	}

	n.sigLoop.AddHandler(&dbusutil.SignalRule{
		Name: "org.freedesktop.NetworkManager.VPN.Connection.VpnStateChanged",
	}, func(sig *dbus.Signal) {
		if strings.HasPrefix(string(sig.Path),
			"/org/freedesktop/NetworkManager/ActiveConnection/") &&
			len(sig.Body) >= 2 {

			state, ok := sig.Body[0].(uint32)
			if !ok {
				return
			}
			reason, ok := sig.Body[1].(uint32)
			if !ok {
				return
			}
			logger.Debug(sig.Path, "vpn state changed", state, reason)
			n.handleVpnStateChanged(state)
		}
	})
}

func (n *Network) getWirelessDevices() (devices []*device) {
	n.devicesMu.Lock()

	for _, d := range n.devices {
		if d.type0 == nm.NM_DEVICE_TYPE_WIFI {
			devices = append(devices, d)
		}
	}

	n.devicesMu.Unlock()
	return
}

func (n *Network) handleVpnStateChanged(state uint32) {
	if state >= nm.NM_VPN_CONNECTION_STATE_PREPARE &&
		state <= nm.NM_VPN_CONNECTION_STATE_ACTIVATED {

		n.PropsMu.Lock()
		changed := n.setPropVpnEnabled(true)
		n.PropsMu.Unlock()

		if changed {
			n.configMu.Lock()
			n.config.VpnEnabled = true
			err := n.saveConfig()
			n.configMu.Unlock()

			if err != nil {
				logger.Warning(err)
			}
		}
	}
}

func (n *Network) addDevice(devPath dbus.ObjectPath) error {
	_, ok := n.devices[devPath]
	if ok {
		return nil
	}

	d, err := networkmanager.NewDevice(n.getSysBus(), devPath)
	if err != nil {
		return err
	}
	iface, err := d.Interface().Get(0)
	if err != nil {
		return err
	}

	deviceType, err := d.DeviceType().Get(0)
	if err != nil {
		return err
	}

	d.InitSignalExt(n.sigLoop, false)
	_, err = d.ConnectStateChanged(func(newState uint32, oldState uint32, reason uint32) {
		//logger.Debugf("device state changed %v newState %d", d.Path_(), newState)
		enabled := n.isIfaceEnabled(iface)
		state, err := d.State().Get(0)
		if err != nil {
			logger.Warning(err)
			return
		}

		if !enabled {
			if state >= nm.NM_DEVICE_STATE_PREPARE &&
				state <= nm.NM_DEVICE_STATE_ACTIVATED {
				logger.Debug("disconnect device", d.Path_())
				err = d.Disconnect(0)
				if err != nil {
					logger.Warning(err)
				}
			}
		}
	})
	if err != nil {
		logger.Warning(err)
	}

	n.devices[devPath] = &device{
		iface:    iface,
		nmDevice: d,
		type0:    deviceType,
	}

	return nil
}

func (n *Network) removeDevice(devPath dbus.ObjectPath) {
	d, ok := n.devices[devPath]
	if !ok {
		return
	}
	d.nmDevice.RemoveHandler(proxy.RemoveAllHandlers)
	delete(n.devices, devPath)
}

func (n *Network) GetInterfaceName() string {
	return dbusInterface
}

func newNetwork() *Network {
	n := new(Network)
	cfg := loadConfigSafe(configFile)
	n.VpnEnabled = cfg.VpnEnabled
	n.config = cfg
	n.devices = make(map[dbus.ObjectPath]*device)
	return n
}

func (n *Network) EnableDevice(pathOrIface string, enabled bool) *dbus.Error {
	logger.Debug("call EnableDevice", pathOrIface, enabled)
	err := n.enableDevice(pathOrIface, enabled)
	return dbusutil.ToError(err)
}

func (n *Network) enableDevice(pathOrIface string, enabled bool) error {
	d := n.findDevice(pathOrIface)
	if d == nil {
		return errors.New("not found device")
	}

	n.enableIface(d.iface, enabled)

	err := n.service.Emit(n, "DeviceEnabled", d.nmDevice.Path_(), enabled)
	if err != nil {
		logger.Warning(err)
	}

	if enabled {
		err = n.enableDevice1(d)
		if err != nil {
			logger.Warning(err)
		}
	} else {
		err = n.disableDevice(d)
		if err != nil {
			logger.Warning(err)
		}
	}

	err = n.saveConfig()
	if err != nil {
		logger.Warning(err)
	}

	return nil
}

func (n *Network) enableDevice1(d *device) error {
	err := n.enableNetworking()
	if err != nil {
		return err
	}

	if d.type0 == nm.NM_DEVICE_TYPE_WIFI {
		err = n.enableWireless()
		if err != nil {
			return err
		}
	}

	err = setDeviceAutoConnect(d.nmDevice, true)
	if err != nil {
		return err
	}

	err = setDeviceManaged(d.nmDevice, true)
	if err != nil {
		return err
	}

	connPaths, err := d.nmDevice.AvailableConnections().Get(0)
	if err != nil {
		return err
	}
	logger.Debug("available connections:", connPaths)

	var connPath0 dbus.ObjectPath
	var maxTs uint64
	for _, connPath := range connPaths {
		connObj, err := networkmanager.NewConnectionSettings(n.getSysBus(), connPath)
		if err != nil {
			logger.Warning(err)
			continue
		}

		settings, err := connObj.GetSettings(0)
		if err != nil {
			logger.Warning(err)
			continue
		}

		auto := getSettingConnectionAutoconnect(settings)
		if !auto {
			continue
		}

		ts := getSettingConnectionTimestamp(settings)
		if maxTs < ts || connPath0 == "" {
			maxTs = ts
			connPath0 = connObj.Path_()
		}
	}

	if connPath0 != "" {
		logger.Debug("activate connection", connPath0)
		_, err = n.nmManager.ActivateConnection(0, connPath0,
			d.nmDevice.Path_(), "/")
		return err
	}
	return nil
}

func (n *Network) disableDevice(d *device) error {
	err := setDeviceAutoConnect(d.nmDevice, false)
	if err != nil {
		return err
	}

	state, err := d.nmDevice.State().Get(0)
	if err != nil {
		return err
	}

	if state >= nm.NM_DEVICE_STATE_PREPARE &&
		state <= nm.NM_DEVICE_STATE_ACTIVATED {
		return d.nmDevice.Disconnect(0)
	}
	return nil
}

func (n *Network) saveConfig() error {
	return saveConfig(configFile, n.config)
}

func (n *Network) IsDeviceEnabled(pathOrIface string) (bool, *dbus.Error) {
	b, err := n.isDeviceEnabled(pathOrIface)
	return b, dbusutil.ToError(err)
}

func (n *Network) isDeviceEnabled(pathOrIface string) (bool, error) {
	d := n.findDevice(pathOrIface)
	if d == nil {
		return false, errors.New("not found device")
	}

	return n.isIfaceEnabled(d.iface), nil
}

func (n *Network) isIfaceEnabled(iface string) bool {
	n.configMu.Lock()
	defer n.configMu.Unlock()

	devCfg, ok := n.config.Devices[iface]
	if !ok {
		// new device default enabled
		return true
	}
	return devCfg.Enabled
}

func (n *Network) enableIface(iface string, enabled bool) {
	n.configMu.Lock()
	deviceConfig := n.config.Devices[iface]
	if deviceConfig == nil {
		deviceConfig = new(DeviceConfig)
		n.config.Devices[iface] = deviceConfig
	}
	deviceConfig.Enabled = enabled
	n.configMu.Unlock()
}

func (n *Network) getDeviceByIface(iface string) *device {
	for _, value := range n.devices {
		if value.iface == iface {
			return value
		}
	}
	return nil
}

func (n *Network) findDevice(pathOrIface string) *device {
	n.devicesMu.Lock()
	defer n.devicesMu.Unlock()

	if strings.HasPrefix(pathOrIface, "/org/freedesktop/NetworkManager") {
		return n.devices[dbus.ObjectPath(pathOrIface)]
	}
	return n.getDeviceByIface(pathOrIface)
}

func (n *Network) enableNetworking() error {
	enabled, err := n.nmManager.NetworkingEnabled().Get(0)
	if err != nil {
		return err
	}

	if enabled {
		return nil
	}

	return n.nmManager.Enable(0, true)
}

func (n *Network) enableWireless() error {
	enabled, err := n.nmManager.WirelessEnabled().Get(0)
	if err != nil {
		return err
	}

	if enabled {
		return nil
	}

	return n.nmManager.WirelessEnabled().Set(0, true)
}

func (n *Network) ToggleWirelessEnabled() (bool, *dbus.Error) {
	enabled, err := n.toggleWirelessEnabled()
	return enabled, dbusutil.ToError(err)
}

func (n *Network) toggleWirelessEnabled() (bool, error) {
	enabled, err := n.nmManager.WirelessEnabled().Get(0)
	if err != nil {
		return false, err
	}
	enabled = !enabled

	err = n.nmManager.WirelessEnabled().Set(0, enabled)
	if err != nil {
		return false, err
	}

	device := n.getWirelessDevices()
	for _, d := range device {
		devPath := d.nmDevice.Path_()
		err = n.enableDevice(string(devPath), enabled)
		if err != nil {
            logger.Warningf("failed to enable %v device %s: %v", enabled, devPath, err)
		}
	}

	return enabled, nil
}

type connSettings struct {
	nmConn   *networkmanager.ConnectionSettings
	uuid     string
	settings map[string]map[string]dbus.Variant
}

func (n *Network) enableVpn1() {
	connSettingsList, err := n.getConnSettingsListByConnType("vpn")
	if err != nil {
		logger.Warning(err)
		return
	}
	for _, connSettings := range connSettingsList {
		autoConnect := getSettingConnectionAutoconnect(connSettings.settings)
		if !autoConnect {
			continue
		}

		connPath := connSettings.nmConn.Path_()
		logger.Debug("activate vpn conn", connPath)
		_, err := n.nmManager.ActivateConnection(0, connPath,
			"/", "/")
		if err != nil {
			logger.Warning(err)
		}
	}
}

func (n *Network) disableVpn() {
	connSettingsList, err := n.getConnSettingsListByConnType("vpn")
	if err != nil {
		logger.Warning(err)
		return
	}

	for _, connSettings := range connSettingsList {
		n.deactivateConnectionByUuid(connSettings.uuid)
	}
}

func (n *Network) enableVpn(enabled bool) {
	if enabled {
		avail, err := n.isNetworkAvailable()
		if err != nil {
			logger.Warning(err)
		}

		if avail {
			n.enableVpn1()
		} else {
			n.PropsMu.Lock()
			n.delayEnableVpn = true
			n.PropsMu.Unlock()
		}

	} else {
		n.disableVpn()

		n.PropsMu.Lock()
		n.delayEnableVpn = false
		n.PropsMu.Unlock()
	}
}

func (n *Network) vpnEnabledWriteCb(write *dbusutil.PropertyWrite) *dbus.Error {
	enabled := write.Value.(bool)
	logger.Debug("set VpnEnabled", enabled)

	if enabled {
		err := n.enableNetworking()
		if err != nil {
			logger.Warning(err)
			return nil
		}
	}

	n.enableVpn(enabled)

	n.configMu.Lock()
	n.config.VpnEnabled = enabled
	err := n.saveConfig()
	n.configMu.Unlock()
	if err != nil {
		logger.Warning(err)
		return dbusutil.ToError(err)
	}

	return nil
}

func (n *Network) isNetworkAvailable() (bool, error) {
	state, err := n.nmManager.State().Get(0)
	if err != nil {
		return false, err
	}

	return state >= nm.NM_STATE_CONNECTED_SITE, nil
}

func (n *Network) deactivateConnectionByUuid(uuid string) {
	activeConns, err := n.getActiveConnectionsByUuid(uuid)
	if err != nil {
		return
	}
	for _, activeConn := range activeConns {
		logger.Debug("DeactivateConnection:", uuid, activeConn.Path_())

		state, err := activeConn.State().Get(0)
		if err != nil {
			logger.Warning(err)
			return
		}

		if state == nm.NM_ACTIVE_CONNECTION_STATE_ACTIVATING ||
			state == nm.NM_ACTIVE_CONNECTION_STATE_ACTIVATED {

			err = n.nmManager.DeactivateConnection(0, activeConn.Path_())
			if err != nil {
				logger.Warning(err)
				return
			}
		}
	}
}

func (n *Network) getActiveConnectionsByUuid(uuid string) ([]*networkmanager.ActiveConnection,
	error) {
	activeConnPaths, err := n.nmManager.ActiveConnections().Get(0)
	if err != nil {
		return nil, err
	}
	var result []*networkmanager.ActiveConnection
	for _, activeConnPath := range activeConnPaths {
		activeConn, err := networkmanager.NewActiveConnection(n.getSysBus(), activeConnPath)
		if err != nil {
			logger.Warning(err)
			continue
		}

		uuid0, err := activeConn.Uuid().Get(0)
		if err != nil {
			logger.Warning(err)
			continue
		}

		if uuid0 == uuid {
			result = append(result, activeConn)
		}
	}
	return result, nil
}

func (n *Network) getConnSettingsListByConnType(connType string) ([]*connSettings, error) {
	connPaths, err := n.nmSettings.ListConnections(0)
	if err != nil {
		return nil, err
	}

	var result []*connSettings
	for _, connPath := range connPaths {
		conn, err := networkmanager.NewConnectionSettings(n.getSysBus(), connPath)
		if err != nil {
			logger.Warning(err)
			continue
		}

		settings, err := conn.GetSettings(0)
		if err != nil {
			logger.Warning(err)
			continue
		}

		if getSettingConnectionType(settings) != connType {
			continue
		}

		uuid := getSettingConnectionUuid(settings)
		if uuid != "" {
			cs := &connSettings{
				nmConn:   conn,
				uuid:     uuid,
				settings: settings,
			}
			result = append(result, cs)
		}
	}
	return result, nil
}
