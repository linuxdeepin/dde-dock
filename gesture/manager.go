package gesture

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"pkg.deepin.io/gir/gio-2.0"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.daemon"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.gesture"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.wm"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
	"pkg.deepin.io/lib/gsettings"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	tsSchemaID           = "com.deepin.dde.touchscreen"
	tsSchemaKeyLongPress = "longpress-duration"
	tsSchemaKeyBlacklist = "longpress-blacklist"
)

type Manager struct {
	wm            *wm.Wm
	sysDaemon     *daemon.Daemon
	systemSigLoop *dbusutil.SignalLoop
	mu            sync.RWMutex
	userFile      string
	builtinSets   map[string]func() error
	gesture       *gesture.Gesture
	setting       *gio.Settings
	tsSetting     *gio.Settings
	enabled       bool
	Infos         gestureInfos

	methods *struct {
		SetLongPressDuration func() `in:"duration"`
		GetLongPressDuration func() `out:"duration"`
	}
}

func newManager() (*Manager, error) {
	sessionConn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}

	systemConn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	var filename = configUserPath
	if !dutils.IsFileExist(configUserPath) {
		filename = configSystemPath
	}

	infos, err := newGestureInfosFromFile(filename)
	if err != nil {
		return nil, err
	}
	// for touch long press
	infos = append(infos, &gestureInfo{
		Name:      "touch right button",
		Direction: "down",
		Fingers:   0,
		Action: ActionInfo{
			Type:   ActionTypeCommandline,
			Action: "xdotool mousedown 3",
		},
	})
	infos = append(infos, &gestureInfo{
		Name:      "touch right button",
		Direction: "up",
		Fingers:   0,
		Action: ActionInfo{
			Type:   ActionTypeCommandline,
			Action: "xdotool mouseup 3",
		},
	})

	setting, err := dutils.CheckAndNewGSettings(gestureSchemaId)
	if err != nil {
		return nil, err
	}

	tsSetting, err := dutils.CheckAndNewGSettings(tsSchemaID)
	if err != nil {
		return nil, err
	}

	m := &Manager{
		userFile:  configUserPath,
		Infos:     infos,
		setting:   setting,
		tsSetting: tsSetting,
		enabled:   setting.GetBoolean(gsKeyEnabled),
		wm:        wm.NewWm(sessionConn),
		sysDaemon: daemon.NewDaemon(systemConn),
	}

	m.gesture = gesture.NewGesture(systemConn)
	m.systemSigLoop = dbusutil.NewSignalLoop(systemConn, 10)
	return m, nil
}

func (m *Manager) destroy() {
	m.gesture.RemoveHandler(proxy.RemoveAllHandlers)
	m.systemSigLoop.Stop()
	m.setting.Unref()
}

func (m *Manager) init() {
	m.initBuiltinSets()
	m.sysDaemon.SetLongPressDuration(0, uint32(m.tsSetting.GetInt(tsSchemaKeyLongPress)))
	m.systemSigLoop.Start()
	m.gesture.InitSignalExt(m.systemSigLoop, true)
	m.gesture.ConnectEvent(func(name string, direction string, fingers int32) {
		logger.Debug("[Event] received:", name, direction, fingers)
		err := m.Exec(name, direction, fingers)
		if err != nil {
			logger.Error("Exec failed:", err)
		}
	})
	m.listenGSettingsChanged()
}

func (m *Manager) Exec(name, direction string, fingers int32) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.enabled || !isSessionActive() {
		logger.Debug("Gesture had been disabled or session inactive")
		return nil
	}

	info := m.Infos.Get(name, direction, fingers)
	if info == nil {
		return fmt.Errorf("not found gesture info for: %s, %s, %d", name, direction, fingers)
	}

	logger.Debug("[Exec] action info:", info.Name, info.Direction, info.Fingers,
		info.Action.Type, info.Action.Action)
	// allow right button up when kbd grabbed
	if (info.Name != "touch right button" || info.Direction != "up") && isKbdAlreadyGrabbed() {
		return fmt.Errorf("another process grabbed keyboard, not exec action")
	}
	// TODO(jouyouyun): improve touch right button handler
	if info.Name == "touch right button" {
		// filter google chrome
		if isInWindowBlacklist(getCurrentActionWindowCmd(), m.tsSetting.GetStrv(tsSchemaKeyBlacklist)) {
			logger.Debug("The current active window in blacklist")
			return nil
		}
	}
	var cmd = info.Action.Action
	switch info.Action.Type {
	case ActionTypeCommandline:
		break
	case ActionTypeShortcut:
		cmd = fmt.Sprintf("xdotool key %s", cmd)
		break
	case ActionTypeBuiltin:
		return m.handleBuiltinAction(cmd)
	default:
		return fmt.Errorf("invalid action type: %s", info.Action.Type)
	}

	out, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(out))
	}
	return nil
}

func (m *Manager) Write() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	err := os.MkdirAll(filepath.Dir(m.userFile), 0755)
	if err != nil {
		return err
	}
	data, err := json.Marshal(m.Infos)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(m.userFile, data, 0644)
}

func (m *Manager) listenGSettingsChanged() {
	gsettings.ConnectChanged(gestureSchemaId, gsKeyEnabled, func(key string) {
		m.mu.Lock()
		m.enabled = m.setting.GetBoolean(key)
		m.mu.Unlock()
	})
}

func (m *Manager) handleBuiltinAction(cmd string) error {
	fn := m.builtinSets[cmd]
	if fn == nil {
		return fmt.Errorf("invalid built-in action %q", cmd)
	}
	return fn()
}

func (*Manager) GetInterfaceName() string {
	return dbusServiceIFC
}
