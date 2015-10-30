package item

import (
	"math/rand"
	"path/filepath"
	"pkg.deepin.io/dde/daemon/launcher/item/dstore"
	"pkg.deepin.io/lib/dbus"
	"time"
)

type MockSoftcenter struct {
	count           int
	disconnectCount int
	handlers        map[int]func([][]interface{})
	softs           map[string]string
}

func (m *MockSoftcenter) GetPkgNameFromPath(path string) (string, error) {
	relPath, _ := filepath.Rel(".", path)
	for pkgName, path := range m.softs {
		if path == relPath {
			return pkgName, nil
		}
	}
	return "", nil
}

func (m *MockSoftcenter) sendMessage(msg [][]interface{}) {
	for _, fn := range m.handlers {
		fn(msg)
	}
}

func (m *MockSoftcenter) sendStartMessage(pkgName string) {
	action := []interface{}{
		dstore.ActionStart,
		dbus.MakeVariant([]interface{}{pkgName, int32(0)}),
	}
	m.sendMessage([][]interface{}{action})
}

func (m *MockSoftcenter) sendUpdateMessage(pkgName string) {
	updateTime := rand.Intn(5) + 1
	for i := 0; i < updateTime; i++ {
		action := []interface{}{
			dstore.ActionUpdate,
			dbus.MakeVariant([]interface{}{
				pkgName,
				int32(1),
				int32(int(i+1) / updateTime),
				"update",
			}),
		}
		m.sendMessage([][]interface{}{action})
		time.Sleep(time.Duration(rand.Int31n(100)+100) * time.Millisecond)
	}
}
func (m *MockSoftcenter) sendFinishedMessage(pkgName string) {
	action := []interface{}{
		dstore.ActionFinish,
		dbus.MakeVariant([]interface{}{
			pkgName,
			int32(2),
			[][]interface{}{
				[]interface{}{
					pkgName,
					true,
					false,
					false,
				},
			},
		}),
	}
	m.sendMessage([][]interface{}{action})
}

func (m *MockSoftcenter) sendFailedMessage(pkgName string) {
	action := []interface{}{
		dstore.ActionFailed,
		dbus.MakeVariant([]interface{}{
			pkgName,
			int32(3),
			[][]interface{}{
				[]interface{}{
					pkgName,
					false,
					false,
					false,
				},
			},
			"uninstall failed",
		}),
	}
	m.sendMessage([][]interface{}{action})
}

func (m *MockSoftcenter) UninstallPkg(pkgName string, purge bool) error {
	if _, ok := m.softs[pkgName]; !ok {
		m.sendFailedMessage(pkgName)
		return nil
	}
	m.sendStartMessage(pkgName)
	m.sendUpdateMessage(pkgName)
	m.sendFinishedMessage(pkgName)
	return nil
}

func (m *MockSoftcenter) Connectupdate_signal(fn func([][]interface{})) func() {
	id := m.count
	m.handlers[id] = fn
	m.count++
	return func() {
		delete(m.handlers, id)
		m.disconnectCount++
	}
}

func NewMockSoftcenter() *MockSoftcenter {
	return &MockSoftcenter{
		handlers: map[int]func([][]interface{}){},
		count:    0,
		softs: map[string]string{
			"firefox":             "../testdata/firefox.desktop",
			"deepin-music-player": "../testdata/deepin-music-player.desktop",
		},
	}
}
