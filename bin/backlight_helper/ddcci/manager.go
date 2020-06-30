package ddcci

import (
	"fmt"
	"sync"

	x "github.com/linuxdeepin/go-x11-client"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/log"
)

const (
	DbusPath      = "/com/deepin/daemon/helper/Backlight/DDCCI"
	dbusInterface = "com.deepin.daemon.helper.Backlight.DDCCI"
)

var logger = log.NewLogger("backlight_helper/ddcci")

type Manager struct {
	service *dbusutil.Service
	ddcci   *DDCCI

	PropsMu         sync.RWMutex
	configTimestamp x.Timestamp

	methods *struct {
		CheckSupport    func() `in:"edidChecksum" out:"support"`
		GetBrightness   func() `in:"edidChecksum" out:"value"`
		SetBrightness   func() `in:"edidChecksum,value"`
		RefreshDisplays func()
	}
}

func NewManager() (*Manager, error) {
	m := &Manager{}

	var err error
	m.ddcci, err = newDDCCI()
	if err != nil {
		return nil, fmt.Errorf("failed to init ddc/ci: %s", err)
	}

	return m, nil
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}

func (m *Manager) CheckSupport(edidChecksum string) (bool, *dbus.Error) {
	return m.ddcci.SupportBrightness(edidChecksum), nil
}

func (m *Manager) GetBrightness(edidChecksum string) (int32, *dbus.Error) {
	if !m.ddcci.SupportBrightness(edidChecksum) {
		err := fmt.Errorf("not support ddc/ci: %s", edidChecksum)
		return 0, dbusutil.ToError(err)
	}

	brightness, err := m.ddcci.GetBrightness(edidChecksum)
	return int32(brightness), dbusutil.ToError(err)
}

func (m *Manager) SetBrightness(edidChecksum string, value int32) *dbus.Error {
	if !m.ddcci.SupportBrightness(edidChecksum) {
		err := fmt.Errorf("not support ddc/ci: %s", edidChecksum)
		return dbusutil.ToError(err)
	}

	err := m.ddcci.SetBrightness(edidChecksum, int(value))
	return dbusutil.ToError(err)
}

func (m *Manager) RefreshDisplays() *dbus.Error {
	m.ddcci.RefreshDisplays()
	return dbusutil.ToError(nil)
}
