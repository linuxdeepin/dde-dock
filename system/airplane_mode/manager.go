package airplane_mode

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"sync"

	dbus "github.com/godbus/dbus"
	nmdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.networkmanager"
	polkit "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.policykit1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	dbusServiceName = "com.deepin.daemon.AirplaneMode"
	dbusPath        = "/com/deepin/daemon/AirplaneMode"
	dbusInterface   = dbusServiceName

	actionId = "com.deepin.daemon.airplane-mode.enable-disable-any"
)

//go:generate dbusutil-gen -type Manager manager.go

type Manager struct {
	service   *dbusutil.Service
	nmManager *nmdbus.Manager
	sigLoop   *dbusutil.SignalLoop

	PropsMu          sync.RWMutex
	Enabled          bool
	WifiEnabled      bool
	BluetoothEnabled bool

	state AirplaneModeState

	// nolint
	methods *struct {
		Enable          func() `in:"enabled"`
		EnableWifi      func() `in:"enabled"`
		EnableBluetooth func() `in:"enabled"`
	}
}

func (m *Manager) GetInterfaceName() string {
	return dbusInterface
}

func newManager(service *dbusutil.Service) *Manager {
	m := &Manager{
		service: service,
	}
	m.state.enableWifiFn = func(enabled bool) {
		logger.Debug("call enableWifiFn", enabled)
		err := m.enableWifi(enabled)
		if err != nil {
			logger.Warning(err)
		}
	}
	m.state.enableBtFn = func(enabled bool) {
		logger.Debug("call enableBtFn", enabled)
		err := m.enableBluetooth(enabled)
		if err != nil {
			logger.Warning(err)
		}
	}

	sysBus := service.Conn()
	m.nmManager = nmdbus.NewManager(sysBus)
	m.sigLoop = dbusutil.NewSignalLoop(sysBus, 10)
	m.sigLoop.Start()

	var err error
	m.BluetoothEnabled, err = getBtEnabled()
	if err != nil {
		logger.Warning(err)
	}

	m.WifiEnabled, err = m.nmManager.WirelessEnabled().Get(0)
	if err != nil {
		logger.Warning("failed to get nmManager WirelessEnabled:", err)
	}

	var cfg config
	err = loadConfig(configFile, &cfg)
	if err != nil {
		if !os.IsNotExist(err) {
			logger.Warning(err)
		}
	}
	m.Enabled = cfg.Enabled

	m.state.WifiEnabled = m.WifiEnabled
	m.state.BtEnabled = m.BluetoothEnabled
	m.state.Enabled = m.Enabled

	m.listenDBusSignals()
	m.listenRfkillEvents()
	return m
}

func (m *Manager) listenDBusSignals() {
	m.nmManager.InitSignalExt(m.sigLoop, true)
	err := m.nmManager.WirelessEnabled().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}

		logger.Debug("nmManager.WirelessEnabled changed to", value)

		m.PropsMu.Lock()
		m.setPropWifiEnabled(value)
		m.PropsMu.Unlock()
		m.state.enableWifi(value)
	})
	if err != nil {
		logger.Warning(err)
	}
}

func (m *Manager) listenRfkillEvents() {
	cmd := exec.Command("rfkill", "event")
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		logger.Warning(err)
		return
	}
	err = cmd.Start()
	if err != nil {
		logger.Warning(err)
		return
	}

	go func() {
		rd := bufio.NewReader(outPipe)
		for {
			line, err := rd.ReadBytes('\n')
			if err == io.EOF {
				break
			}
			logger.Debugf("rfkill event: %s", bytes.TrimSpace(line))
			m.handleRfkillEvent()
		}
		err = cmd.Wait()
		if err != nil {
			logger.Warning(err)
		}
	}()
}

func (m *Manager) handleRfkillEvent() {
	enabled, err := getBtEnabled()
	if err != nil {
		logger.Warning(err)
		return
	}

	m.PropsMu.Lock()
	if m.BluetoothEnabled == enabled {
		m.PropsMu.Unlock()
		return
	}
	m.setPropBluetoothEnabled(enabled)
	m.PropsMu.Unlock()

	m.state.enableBt(enabled)
}

func (m *Manager) DumpState() *dbus.Error {
	m.state.dump()
	return nil
}

func (m *Manager) Enable(sender dbus.Sender, enabled bool) *dbus.Error {
	err := checkAuthorization(actionId, string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}

	m.PropsMu.RLock()
	if m.Enabled == enabled {
		m.PropsMu.RUnlock()
		return nil
	}
	m.PropsMu.RUnlock()

	m.state.enable(enabled)

	m.PropsMu.Lock()
	m.setPropEnabled(enabled)
	m.PropsMu.Unlock()

	err = m.saveConfig(enabled)
	if err != nil {
		logger.Warning(err)
		return dbusutil.ToError(err)
	}

	return nil
}

func (m *Manager) saveConfig(enabled bool) error {
	var cfg config
	cfg.Enabled = enabled
	err := saveConfig(configFile, &cfg)
	return err
}

func (m *Manager) enableWifi(enabled bool) error {
	wEnabled, err := m.nmManager.WirelessEnabled().Get(0)
	if err != nil {
		return err
	}

	if wEnabled == enabled {
		return nil
	}

	err = m.nmManager.WirelessEnabled().Set(0, enabled)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) EnableWifi(sender dbus.Sender, enabled bool) *dbus.Error {
	err := checkAuthorization(actionId, string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}

	err = m.enableWifi(enabled)
	return dbusutil.ToError(err)
}

func (m *Manager) EnableBluetooth(sender dbus.Sender, enabled bool) *dbus.Error {
	err := checkAuthorization(actionId, string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}

	err = m.enableBluetooth(enabled)
	return dbusutil.ToError(err)
}

func (m *Manager) enableBluetooth(enabled bool) error {
	return enableBt(enabled)
}

func checkAuthorization(actionId string, sysBusName string) error {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	authority := polkit.NewAuthority(systemBus)
	subject := polkit.MakeSubject(polkit.SubjectKindSystemBusName)
	subject.SetDetail("name", sysBusName)

	ret, err := authority.CheckAuthorization(0, subject, actionId,
		nil, polkit.CheckAuthorizationFlagsAllowUserInteraction, "")
	if err != nil {
		return err
	}
	if !ret.IsAuthorized {
		return errors.New("not authorized")
	}

	return nil
}
