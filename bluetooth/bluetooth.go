/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package bluetooth

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	dbus "github.com/godbus/dbus"
	apidevice "github.com/linuxdeepin/go-dbus-factory/com.deepin.api.device"
	bluez "github.com/linuxdeepin/go-dbus-factory/org.bluez"
	obex "github.com/linuxdeepin/go-dbus-factory/org.bluez.obex"
	ofdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.dbus"
	"pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/gsprop"
	"pkg.deepin.io/lib/dbusutil/proxy"
)

const (
	bluezDBusServiceName      = "org.bluez"
	bluezAdapterDBusInterface = "org.bluez.Adapter1"
	bluezDeviceDBusInterface  = "org.bluez.Device1"

	dbusServiceName = "com.deepin.daemon.Bluetooth"
	dbusPath        = "/com/deepin/daemon/Bluetooth"
	dbusInterface   = dbusServiceName

	daemonSysService = "com.deepin.daemon.Daemon"
	daemonSysPath    = "/com/deepin/daemon/Daemon"
	daemonSysIFC     = daemonSysService

	methodSysBlueGetDeviceTech = daemonSysIFC + ".BluetoothGetDeviceTechnologies"
)
const (
	bluetoothSchema = "com.deepin.dde.bluetooth"
	displaySwitch   = "display-switch"
)

const (
	StateUnavailable = 0
	StateAvailable   = 1
	StateConnected   = 2
)

// device type index
const (
	Computer = iota
	Phone
	Modem
	NetworkWireless
	AudioCard
	CameraVideo
	InputGaming
	InputKeyboard
	InputTablet
	InputMouse
	Printer
	CameraPhone
)

// nolint
const (
	transferStatusQueued    = "queued"
	transferStatusActive    = "active"
	transferStatusSuspended = "suspended"
	transferStatusComplete  = "complete"
	transferStatusError     = "error"
)

var DeviceTypes []string = []string{
	"computer",
	"phone",
	"modem",
	"network-wireless",
	"audio-card",
	"camera-video",
	"input-gaming",
	"input-keyboard",
	"input-tablet",
	"input-mouse",
	"printer",
	"camera-photo",
}

//go:generate dbusutil-gen -type Bluetooth bluetooth.go

type Bluetooth struct {
	service       *dbusutil.Service
	sigLoop       *dbusutil.SignalLoop
	systemSigLoop *dbusutil.SignalLoop
	config        *config
	objectManager *bluez.ObjectManager
	sysDBusDaemon *ofdbus.DBus
	apiDevice     *apidevice.Device
	agent         *agent
	obexAgent     *obexAgent
	obexManager   *obex.Manager

	// adapter
	adaptersLock sync.Mutex
	adapters     map[dbus.ObjectPath]*adapter

	// device
	devicesLock sync.Mutex
	devices     map[dbus.ObjectPath][]*device

	PropsMu sync.RWMutex
	State   uint32 // StateUnavailable/StateAvailable/StateConnected

	// 当发起设备连接成功后，应该把连接的设备添加进设备列表
	connectedDevices map[dbus.ObjectPath][]*device
	connectedLock    sync.RWMutex

	sessionCancelChMap   map[dbus.ObjectPath]chan struct{}
	sessionCancelChMapMu sync.Mutex

	settings      *gio.Settings
	DisplaySwitch gsprop.Bool `prop:"access:rw"`

	// nolint
	methods *struct {
		DebugInfo                     func() `out:"info"`
		GetDevices                    func() `in:"adapter" out:"devicesJSON"`
		ConnectDevice                 func() `in:"device"`
		DisconnectDevice              func() `in:"device"`
		RemoveDevice                  func() `in:"adapter,device"`
		SetDeviceAlias                func() `in:"device,alias"`
		SetDeviceTrusted              func() `in:"device,trusted"`
		Confirm                       func() `in:"device,accept"`
		FeedPinCode                   func() `in:"device,accept,pinCode"`
		FeedPasskey                   func() `in:"device,accept,passkey"`
		GetAdapters                   func() `out:"adaptersJSON"`
		RequestDiscovery              func() `in:"adapter"`
		SendFiles                     func() `in:"device,files" out:"sessionPath"`
		CancelTransferSession         func() `in:"sessionPath"`
		SetAdapterPowered             func() `in:"adapter,powered"`
		SetAdapterAlias               func() `in:"adapter,alias"`
		SetAdapterDiscoverable        func() `in:"adapter,discoverable"`
		SetAdapterDiscovering         func() `in:"adapter,discovering"`
		SetAdapterDiscoverableTimeout func() `in:"adapter,timeout"`
	}

	// nolint
	signals *struct {
		// adapter/device properties changed signals
		AdapterAdded, AdapterRemoved, AdapterPropertiesChanged struct {
			adapterJSON string
		}

		DeviceAdded, DeviceRemoved, DevicePropertiesChanged struct {
			devJSON string
		}

		// pair request signals
		DisplayPinCode struct {
			device  dbus.ObjectPath
			pinCode string
		}
		DisplayPasskey struct {
			device  dbus.ObjectPath
			passkey uint32
			entered uint32
		}

		// RequestConfirmation you should call Confirm with accept
		RequestConfirmation struct {
			device  dbus.ObjectPath
			passkey string
		}

		// RequestAuthorization you should call Confirm with accept
		RequestAuthorization struct {
			device dbus.ObjectPath
		}

		// RequestPinCode you should call FeedPinCode with accept and key
		RequestPinCode struct {
			device dbus.ObjectPath
		}

		// RequestPasskey you should call FeedPasskey with accept and key
		RequestPasskey struct {
			device dbus.ObjectPath
		}

		Cancelled struct {
			device dbus.ObjectPath
		}

		ObexSessionCreated struct {
			sessionPath dbus.ObjectPath
		}

		ObexSessionRemoved struct {
			sessionPath dbus.ObjectPath
		}

		ObexSessionProgress struct {
			sessionPath dbus.ObjectPath
			totalSize   uint64
			transferred uint64
			currentIdx  int
		}

		TransferCreated struct {
			file         string
			transferPath dbus.ObjectPath
			sessionPath  dbus.ObjectPath
		}

		TransferRemoved struct {
			file         string
			transferPath dbus.ObjectPath
			sessionPath  dbus.ObjectPath
			done         bool
		}
		TransferFailed struct {
			file        string
			sessionPath dbus.ObjectPath
			errInfo     string
		}
	}
}

func newBluetooth(service *dbusutil.Service) (b *Bluetooth) {
	systemConn, err := dbus.SystemBus()
	if err != nil {
		logger.Warning(err)
		return nil
	}

	b = &Bluetooth{
		service:       service,
		sigLoop:       dbusutil.NewSignalLoop(service.Conn(), 0),
		systemSigLoop: dbusutil.NewSignalLoop(systemConn, 10),
		obexManager:   obex.NewManager(service.Conn()),
	}

	b.adapters = make(map[dbus.ObjectPath]*adapter)
	return
}

func (b *Bluetooth) destroy() {
	b.agent.destroy()

	b.objectManager.RemoveHandler(proxy.RemoveAllHandlers)
	b.sysDBusDaemon.RemoveHandler(proxy.RemoveAllHandlers)

	b.devicesLock.Lock()
	for _, devices := range b.devices {
		for _, device := range devices {
			device.destroy()
		}
	}
	b.devicesLock.Unlock()

	b.adaptersLock.Lock()
	for _, adapter := range b.adapters {
		adapter.destroy()
	}
	b.adaptersLock.Unlock()

	err := b.service.StopExport(b)
	if err != nil {
		logger.Warning(err)
	}
	b.systemSigLoop.Stop()
}

func (*Bluetooth) GetInterfaceName() string {
	return dbusInterface
}

func (b *Bluetooth) init() {
	b.sigLoop.Start()

	b.systemSigLoop.Start()
	b.config = newConfig()
	b.devices = make(map[dbus.ObjectPath][]*device)

	systemBus := b.systemSigLoop.Conn()

	b.connectedLock.Lock()
	b.connectedDevices = make(map[dbus.ObjectPath][]*device, len(DeviceTypes))
	b.connectedLock.Unlock()

	b.sessionCancelChMap = make(map[dbus.ObjectPath]chan struct{})

	// start bluetooth goroutine
	// monitor click signal or time out signal to close notification window
	go beginTimerNotify(globalTimerNotifier)

	b.apiDevice = apidevice.NewDevice(systemBus)
	b.sysDBusDaemon = ofdbus.NewDBus(systemBus)
	b.sysDBusDaemon.InitSignalExt(b.systemSigLoop, true)
	_, err := b.sysDBusDaemon.ConnectNameOwnerChanged(b.handleDBusNameOwnerChanged)
	if err != nil {
		logger.Warning(err)
	}

	b.settings = gio.NewSettings(bluetoothSchema)
	b.DisplaySwitch.Bind(b.settings, displaySwitch)

	// initialize dbus object manager
	b.objectManager = bluez.NewObjectManager(systemBus)

	// connect signals
	b.objectManager.InitSignalExt(b.systemSigLoop, true)
	_, err = b.objectManager.ConnectInterfacesAdded(b.handleInterfacesAdded)
	if err != nil {
		logger.Warning(err)
	}

	_, err = b.objectManager.ConnectInterfacesRemoved(b.handleInterfacesRemoved)
	if err != nil {
		logger.Warning(err)
	}

	b.agent.init()
	b.loadObjects()

	b.obexAgent.init()

	b.config.clearSpareConfig(b)
	b.config.save()
	go b.tryConnectPairedDevices()
	// move to power module
	// b.wakeupWorkaround()
}

func (b *Bluetooth) unblockBluetoothDevice() {
	has, err := b.apiDevice.HasBluetoothDeviceBlocked(0)
	if err != nil {
		logger.Warning(err)
		return
	}

	if has {
		err = b.apiDevice.UnblockBluetoothDevices(0)
		if err != nil {
			logger.Warning(err)
		}
	}
}

func (b *Bluetooth) loadObjects() {
	// add exists adapters and devices
	objects, err := b.objectManager.GetManagedObjects(0)
	if err != nil {
		logger.Error(err)
		return
	}

	b.unblockBluetoothDevice()
	// add adapters
	for path, obj := range objects {
		if _, ok := obj[bluezAdapterDBusInterface]; ok {
			b.addAdapter(path)
		}
	}

	// then add devices
	for path, obj := range objects {
		if _, ok := obj[bluezDeviceDBusInterface]; ok {
			b.addDevice(path)
		}
	}
}

func (b *Bluetooth) removeAllObjects() {
	b.devicesLock.Lock()
	for _, devices := range b.devices {
		for _, device := range devices {
			device.notifyDeviceRemoved()
			device.destroy()
		}
	}
	b.devices = make(map[dbus.ObjectPath][]*device)
	b.devicesLock.Unlock()

	b.adaptersLock.Lock()
	for _, adapter := range b.adapters {
		adapter.notifyAdapterRemoved()
		adapter.destroy()
	}
	b.adapters = make(map[dbus.ObjectPath]*adapter)
	b.adaptersLock.Unlock()
}

func (b *Bluetooth) handleInterfacesAdded(path dbus.ObjectPath, data map[string]map[string]dbus.Variant) {
	if _, ok := data[bluezAdapterDBusInterface]; ok {
		b.unblockBluetoothDevice()
		b.addAdapter(path)
	}
	if _, ok := data[bluezDeviceDBusInterface]; ok {
		b.addDevice(path)
	}
}

func (b *Bluetooth) handleInterfacesRemoved(path dbus.ObjectPath, interfaces []string) {
	if isStringInArray(bluezAdapterDBusInterface, interfaces) {
		b.removeAdapter(path)
	}
	if isStringInArray(bluezDeviceDBusInterface, interfaces) {
		b.removeDevice(path)
	}
}

func (b *Bluetooth) handleDBusNameOwnerChanged(name, oldOwner, newOwner string) {
	// if a new dbus session was installed, the name and newOwner
	// will be not empty, if a dbus session was uninstalled, the
	// name and oldOwner will be not empty
	if name != bluezDBusServiceName {
		return
	}
	if newOwner != "" {
		logger.Info("bluetooth is starting")
		time.AfterFunc(1*time.Second, func() {
			b.loadObjects()
			b.agent.registerDefaultAgent()
		})
	} else {
		logger.Info("bluetooth stopped")
		b.removeAllObjects()
	}
}

func (b *Bluetooth) addDevice(dpath dbus.ObjectPath) {
	if b.isDeviceExists(dpath) {
		logger.Warning("repeat add device", dpath)
		return
	}

	d := newDevice(b.systemSigLoop, dpath)
	b.adaptersLock.Lock()
	d.adapter = b.adapters[d.AdapterPath]
	b.adaptersLock.Unlock()

	if d.adapter == nil {
		logger.Warningf("failed to add device %s, not found adapter", dpath)
		return
	}

	// device detail info is needed to write into config file
	b.config.addDeviceConfig(d)

	b.devicesLock.Lock()
	b.devices[d.AdapterPath] = append(b.devices[d.AdapterPath], d)
	b.devicesLock.Unlock()

	connected := b.config.getDeviceConfigConnected(d.getAddress())
	if connected {
		d, _ := b.getDevice(dpath)
		if d == nil {
			return
		}

		adapterPowered, err := d.adapter.core.Powered().Get(0)
		if err != nil {
			logger.Warning(err)
			return
		}

		if !adapterPowered {
			return
		}

		paired, err := d.core.Paired().Get(0)
		if err != nil {
			logger.Warning(err)
			return
		}

		connected, err := d.core.Connected().Get(0)
		if err != nil {
			logger.Warning(err)
			return
		}
		if paired && connected {
			b.addConnectedDevice(d)
		}
	}

	d.notifyDeviceAdded()
}

func (b *Bluetooth) removeDevice(dpath dbus.ObjectPath) {
	apath, i := b.getDeviceIndex(dpath)
	if i < 0 {
		logger.Error("repeat remove device", dpath)
		return
	}

	b.devicesLock.Lock()
	defer b.devicesLock.Unlock()
	b.devices[apath] = b.doRemoveDevice(b.devices[apath], i)
}

func (b *Bluetooth) doRemoveDevice(devices []*device, i int) []*device {
	// NOTE: do not remove device from config
	d := devices[i]
	d.notifyDeviceRemoved()
	d.destroy()
	copy(devices[i:], devices[i+1:])
	devices[len(devices)-1] = nil
	devices = devices[:len(devices)-1]
	return devices
}

func (b *Bluetooth) isDeviceExists(dpath dbus.ObjectPath) bool {
	_, i := b.getDeviceIndex(dpath)
	return i >= 0
}

func (b *Bluetooth) findDevice(dpath dbus.ObjectPath) (apath dbus.ObjectPath, index int) {
	for p, devices := range b.devices {
		for i, d := range devices {
			if d.Path == dpath {
				return p, i
			}
		}
	}
	return "", -1
}

func (b *Bluetooth) getDeviceIndex(dpath dbus.ObjectPath) (apath dbus.ObjectPath, index int) {
	b.devicesLock.Lock()
	defer b.devicesLock.Unlock()
	return b.findDevice(dpath)
}

func (b *Bluetooth) getDevice(dpath dbus.ObjectPath) (*device, error) {
	b.devicesLock.Lock()
	defer b.devicesLock.Unlock()
	apath, index := b.findDevice(dpath)
	if index < 0 {
		return nil, errInvalidDevicePath
	}
	return b.devices[apath][index], nil
}

func (b *Bluetooth) getAdapterDevices(adapterAddress string) []*device {
	var aPath dbus.ObjectPath
	b.adaptersLock.Lock()
	for adapterPath, adapter := range b.adapters {
		if adapter.address == adapterAddress {
			aPath = adapterPath
			break
		}
	}
	b.adaptersLock.Unlock()

	if aPath == "" {
		return nil
	}

	b.devicesLock.Lock()
	defer b.devicesLock.Unlock()

	devices := b.devices[aPath]
	if devices == nil {
		return nil
	}

	result := make([]*device, 0, len(devices))
	result = append(result, devices...)
	return result
}

func (b *Bluetooth) addAdapter(apath dbus.ObjectPath) {
	if b.isAdapterExists(apath) {
		logger.Warning("repeat add adapter", apath)
		return
	}

	a := newAdapter(b.systemSigLoop, apath)
	// initialize adapter power state
	b.config.addAdapterConfig(a.address)
	cfgPowered := b.config.getAdapterConfigPowered(a.address)
	err := a.core.Powered().Set(0, cfgPowered)
	if err != nil {
		logger.Warning(err)
	}

	if cfgPowered {
		err = a.core.DiscoverableTimeout().Set(0, 0)
		if err != nil {
			logger.Warning(err)
		}

		err = a.core.Discoverable().Set(0, b.config.Discoverable)
		if err != nil {
			logger.Warning(err)
		}
	}

	b.adaptersLock.Lock()
	b.adapters[apath] = a
	b.adaptersLock.Unlock()

	a.notifyAdapterAdded()
}

func (b *Bluetooth) removeAdapter(apath dbus.ObjectPath) {
	b.adaptersLock.Lock()
	defer b.adaptersLock.Unlock()

	if b.adapters[apath] == nil {
		logger.Warning("repeat remove adapter", apath)
		return
	}

	b.doRemoveAdapter(apath)
}

func (b *Bluetooth) doRemoveAdapter(apath dbus.ObjectPath) {
	// NOTE: do not remove adapter from config file
	removeAdapter := b.adapters[apath]
	delete(b.adapters, apath)
	removeAdapter.notifyAdapterRemoved()
	removeAdapter.destroy()
}

func (b *Bluetooth) getAdapter(apath dbus.ObjectPath) (a *adapter, err error) {
	b.adaptersLock.Lock()
	defer b.adaptersLock.Unlock()

	a = b.adapters[apath]
	if a == nil {
		err = fmt.Errorf("adapter not exists %s", apath)
		logger.Error(err)
		return
	}
	return
}
func (b *Bluetooth) isAdapterExists(apath dbus.ObjectPath) bool {
	b.adaptersLock.Lock()
	defer b.adaptersLock.Unlock()
	return b.adapters[apath] != nil
}

func (b *Bluetooth) feed(devPath dbus.ObjectPath, accept bool, key string) (err error) {
	_, err = b.getDevice(devPath)
	if nil != err {
		logger.Warningf("FeedRequest can not find device: %v, %v", devPath, err)
		return err
	}

	b.agent.mu.Lock()
	if b.agent.requestDevice != devPath {
		b.agent.mu.Unlock()
		logger.Warningf("FeedRequest can not find match device: %q, %q", b.agent.requestDevice, devPath)
		return errBluezCanceled
	}
	b.agent.mu.Unlock()

	select {
	case b.agent.rspChan <- authorize{path: devPath, accept: accept, key: key}:
		return nil
	default:
		return errors.New("rspChan no reader")
	}
}

func (b *Bluetooth) updateState() {
	newState := StateUnavailable
	if len(b.adapters) > 0 {
		newState = StateAvailable
	}

	for _, devices := range b.devices {
		for _, d := range devices {
			if d.connected {
				newState = StateConnected
				break
			}
		}
	}

	b.PropsMu.Lock()
	b.setPropState(uint32(newState))
	b.PropsMu.Unlock()
}

func (b *Bluetooth) tryConnectPairedDevices() {
	//input and audio devices counter
	typeMap := make(map[string]uint8)
	typeMap["audio-card"] = 0
	typeMap["input-keyboard"] = 0
	typeMap["input-mouse"] = 0
	typeMap["input-tablet"] = 0

	var devList = b.getPairedDeviceList()
	for _, dev := range devList {
		// make sure dev always exist
		if dev == nil {
			continue
		}
		//connect back to a device
		switch dev.Icon {
		case "audio-card", "input-keyboard", "input-mouse", "input-tablet":
			if typeMap[dev.Icon] == 0 {
				if b.tryConnectPairedDevice(dev) {
					typeMap[dev.Icon]++
				}
			}
		default:
			b.tryConnectPairedDevice(dev)
		}
	}
}

func (b *Bluetooth) tryConnectPairedDevice(dev *device) bool {
	logger.Info("[DEBUG] Auto connect device:", dev.Path)

	// if device using LE mode, will suspend, try connect should be failed, filter it.
	if !b.isBREDRDevice(dev) {
		return false
	}
	logger.Debug("Will auto connect device:", dev.String(), dev.adapter.address, dev.Address)
	err := dev.doConnect(false)
	if err != nil {
		logger.Debug("failed to connect:", dev.String(), err)
		return false
	} else {
		// if auto connect success, add device into map connectedDevices
		if dev.ConnectState {
			b.addConnectedDevice(dev)
		}
	}
	return true
}

// get paired device list
func (b *Bluetooth) getPairedDeviceList() []*device {
	// memory lock
	b.adaptersLock.Lock()
	defer b.adaptersLock.Unlock()
	b.devicesLock.Lock()
	defer b.devicesLock.Unlock()

	// get all paired devices list from adapters
	var devAddressMap = make(map[string]*device)
	for _, aobj := range b.adapters {
		logger.Info("[DEBUG] Auto connect adapter:", aobj.Path)

		// check if devices list in current adapter is legal
		list := b.devices[aobj.Path]
		if len(list) == 0 || !b.config.getAdapterConfigPowered(aobj.address) {
			continue
		}

		// add devices info to list
		for _, value := range list {
			// select already paired but not connected device from list
			if value == nil || !value.Paired || value.connected {
				continue
			}
			devAddressMap[value.getAddress()] = value
			logger.Debug("devAddressMap", value)
		}
	}
	// select the latest devices of each deviceType and add them into list
	devList := b.config.filterDemandedTypeDevices(devAddressMap)

	return devList
}

func (b *Bluetooth) getTechnologies(dev *device) ([]string, error) {
	var technologies []string
	err := b.systemSigLoop.Conn().Object(daemonSysService,
		daemonSysPath).Call(methodSysBlueGetDeviceTech, 0,
		dev.adapter.address, dev.Address).Store(&technologies)
	if err != nil {
		return nil, err
	}
	return technologies, nil
}

func (b *Bluetooth) isBREDRDevice(dev *device) bool {
	technologies, err := b.getTechnologies(dev)
	if err != nil {
		logger.Warningf("failed to get device(%s -- %s) technologies: %v",
			dev.adapter.address, dev.Address, err)
		return false
	}
	for _, tech := range technologies {
		if tech == "BR/EDR" {
			return true
		}
	}
	return false
}

func (b *Bluetooth) addConnectedDevice(connectedDev *device) {
	b.connectedLock.Lock()
	b.connectedDevices[connectedDev.AdapterPath] = append(b.connectedDevices[connectedDev.AdapterPath], connectedDev)
	b.connectedLock.Unlock()
}

func (b *Bluetooth) removeConnectedDevice(disconnectedDev *device) {
	b.connectedLock.Lock()
	// check if dev exist in connectedDevices map, if exist, remove device
	if connectedDevices, ok := globalBluetooth.connectedDevices[disconnectedDev.AdapterPath]; ok {
		var tempDevices []*device
		for _, dev := range connectedDevices {
			// check if disconnected device exist in connected devices, if exist, abandon this
			if dev.Address != disconnectedDev.Address {
				tempDevices = append(tempDevices, dev)
			}
		}
		globalBluetooth.connectedDevices[disconnectedDev.AdapterPath] = tempDevices
	}
	b.connectedLock.Unlock()
}

func (b *Bluetooth) getConnectedDeviceByAddress(address string) *device {
	b.connectedLock.Lock()
	defer b.connectedLock.Unlock()

	for _, v := range b.connectedDevices {
		for _, dev := range v {
			if dev.Address == address {
				return dev
			}
		}
	}

	return nil
}

func (b *Bluetooth) sendFiles(dev *device, files []string) (dbus.ObjectPath, error) {
	var totalSize uint64

	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			return "", err
		}

		totalSize += uint64(info.Size())
	}
	// 创建 OBEX session
	args := make(map[string]dbus.Variant)
	args["Source"] = dbus.MakeVariant(dev.adapter.address) // 蓝牙适配器地址
	args["Target"] = dbus.MakeVariant("opp")               // 连接方式「OPP」
	sessionPath, err := b.obexManager.CreateSession(0, dev.Address, args)
	if err != nil {
		logger.Warning("failed to create obex session:", err)
		return "", err
	}
	b.emitObexSessionCreated(sessionPath)

	session, err := obex.NewSession(b.service.Conn(), sessionPath)
	if err != nil {
		logger.Warning("failed to get session bus:", err)
		return "", err
	}

	go b.doSendFiles(session, files, totalSize)

	return sessionPath, nil
}

func (b *Bluetooth) doSendFiles(session *obex.Session, files []string, totalSize uint64) {
	sessionPath := session.Path_()
	cancelCh := make(chan struct{})

	b.sessionCancelChMapMu.Lock()
	b.sessionCancelChMap[sessionPath] = cancelCh
	b.sessionCancelChMapMu.Unlock()

	var transferredBase uint64

	for i, f := range files {
		_, err := os.Stat(f)
		if err != nil {
			b.emitTransferFailed(f, sessionPath, err.Error())
			break
		}
		transferPath, properties, err := session.ObjectPush().SendFile(0, f)
		if err != nil {
			logger.Warningf("failed to send file: %s: %s", f, err)
			continue
		}
		logger.Infof("properties: %v", properties)

		transfer, err := obex.NewTransfer(b.service.Conn(), transferPath)
		if err != nil {
			logger.Warningf("failed to send file: %s: %s", f, err)
			continue
		}

		transfer.InitSignalExt(b.sigLoop, true)

		b.emitTransferCreated(f, transferPath, sessionPath)

		ch := make(chan bool)
		err = transfer.Status().ConnectChanged(func(hasValue bool, value string) {
			if !hasValue {
				return
			}

			// 成功或者失败，说明这个传输结束
			if value == transferStatusComplete || value == transferStatusError {
				ch <- value == transferStatusComplete
			}
		})
		if err != nil {
			logger.Warning("connect to status changed failed:", err)
		}

		err = transfer.Transferred().ConnectChanged(func(hasValue bool, value uint64) {
			if !hasValue {
				return
			}

			transferred := transferredBase + value
			b.emitObexSessionProgress(sessionPath, totalSize, transferred, i+1)
		})
		if err != nil {
			logger.Warning("connect to transferred changed failed:", err)
		}

		var res bool
		var cancel bool
		select {
		case res = <-ch:
		case <-cancelCh:
			b.sessionCancelChMapMu.Lock()
			delete(b.sessionCancelChMap, sessionPath)
			b.sessionCancelChMapMu.Unlock()

			cancel = true
			err = transfer.Cancel(0)
			if err != nil {
				logger.Warning("failed to cancel transfer:", err)
			}
		}

		b.emitTransferRemoved(f, transferPath, sessionPath, res)

		if cancel {
			break
		}

		info, err := os.Stat(f)
		if err != nil {
			logger.Warning("failed to stat file:", err)
			continue
		} else {
			transferredBase += uint64(info.Size())
		}

		b.emitObexSessionProgress(sessionPath, totalSize, transferredBase, i+1)
	}

	b.sessionCancelChMapMu.Lock()
	delete(b.sessionCancelChMap, sessionPath)
	b.sessionCancelChMapMu.Unlock()

	b.emitObexSessionRemoved(sessionPath)

	objs, err := obex.NewObjectManager(b.service.Conn()).GetManagedObjects(0)
	if err != nil {
		logger.Warning("failed to get managed objects:", err)
	} else {
		_, pathExists := objs[sessionPath]
		if !pathExists {
			logger.Debugf("session %s not exists", sessionPath)
			return
		}
	}

	err = b.obexManager.RemoveSession(0, sessionPath)
	if err != nil {
		logger.Warning("failed to remove session:", err)
	}
}

func (b *Bluetooth) emitObexSessionCreated(sessionPath dbus.ObjectPath) {
	err := b.service.Emit(b, "ObexSessionCreated", sessionPath)
	if err != nil {
		logger.Warning("failed to emit ObexSessionCreated:", err)
	}
}

func (b *Bluetooth) emitObexSessionRemoved(sessionPath dbus.ObjectPath) {
	err := b.service.Emit(b, "ObexSessionRemoved", sessionPath)
	if err != nil {
		logger.Warning("failed to emit ObexSessionRemoved:", err)
	}
}

func (b *Bluetooth) emitObexSessionProgress(sessionPath dbus.ObjectPath, totalSize uint64, transferred uint64, currentIdx int) {
	err := b.service.Emit(b, "ObexSessionProgress", sessionPath, totalSize, transferred, currentIdx)
	if err != nil {
		logger.Warning("failed to emit ObexSessionProgress:", err)
	}
}

func (b *Bluetooth) emitTransferCreated(file string, transferPath dbus.ObjectPath, sessionPath dbus.ObjectPath) {
	err := b.service.Emit(b, "TransferCreated", file, transferPath, sessionPath)
	if err != nil {
		logger.Warning("failed to emit TransferCreated:", err)
	}
}

func (b *Bluetooth) emitTransferRemoved(file string, transferPath dbus.ObjectPath, sessionPath dbus.ObjectPath, done bool) {
	err := b.service.Emit(b, "TransferRemoved", file, transferPath, sessionPath, done)
	if err != nil {
		logger.Warning("failed to emit TransferRemoved:", err)
	}
}

func (b *Bluetooth) emitTransferFailed(file string, sessionPath dbus.ObjectPath, errInfo string) {
	err := b.service.Emit(b, "TransferFailed", file, sessionPath, errInfo)
	if err != nil {
		logger.Warning("failed to emit TransferFailed:", err)
	}
}
