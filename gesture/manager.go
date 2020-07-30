package gesture

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.daemon"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.display"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.gesture"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.dde.daemon.dock"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.wm"
	"pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
	"pkg.deepin.io/lib/gsettings"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	tsSchemaID            = "com.deepin.dde.touchscreen"
	tsSchemaKeyLongPress  = "longpress-duration"
	tsSchemaKeyShortPress = "shortpress-duration"
	tsSchemaKeyBlacklist  = "longpress-blacklist"
)

type Manager struct {
	wm            *wm.Wm
	sysDaemon     *daemon.Daemon
	systemSigLoop *dbusutil.SignalLoop
	mu            sync.RWMutex
	userFile      string
	builtinSets   map[string]func() error
	gesture       *gesture.Gesture
	dock          *dock.Dock
	display       *display.Display
	setting       *gio.Settings
	tsSetting     *gio.Settings
	enabled       bool
	Infos         gestureInfos
	lastAction    int  // 四、五指的最新一次动作

	methods *struct {
		SetLongPressDuration  func() `in:"duration"`
		GetLongPressDuration  func() `out:"duration"`
		SetShortPressDuration func() `in:"duration"`
		GetShortPressDuration func() `out:"duration"`
	}
}

const (
	lastActionNone = iota
	lastActionShowWorkspace
	lastActionHideWorkspace
	lastActionShowDesktop
	lastActionHideDesktop
)

func (m *Manager) setLastAction(v int) {
	m.mu.Lock()
	m.lastAction = v
	m.mu.Unlock()
}

func (m *Manager) getLastAction() int {
	m.mu.Lock()
	value := m.lastAction
	m.mu.Unlock()
	return value
}

func (m *Manager) getEnable() bool {
	m.mu.RLock()
	value := m.enabled
	m.mu.RUnlock()
	return value
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
		dock:      dock.NewDock(sessionConn),
		display:   display.NewDisplay(sessionConn),
		sysDaemon: daemon.NewDaemon(systemConn),
		lastAction: lastActionNone,
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
	err := m.gesture.SetShortPressDuration(0, uint32(m.tsSetting.GetInt(tsSchemaKeyShortPress)))
	if err != nil {
		logger.Warning("call SetShortPressDuration failed:", err)
	}

	m.systemSigLoop.Start()
	m.gesture.InitSignalExt(m.systemSigLoop, true)
	_, err = m.gesture.ConnectEvent(func(name string, direction string, fingers int32) {
		logger.Debug("[Event] received:", name, direction, fingers)
		err = m.Exec(name, direction, fingers)
		if err != nil {
			logger.Error("Exec failed:", err)
		}
	})
	if err != nil {
		logger.Error("connect gesture event failed:", err)
	}

	_, err = m.gesture.ConnectTouchEdgeMoveStopLeave(func(direction string, scaleX float64, scaleY float64, duration int32) {
		logger.Debug("----------[ConnectTouchEdgeMoveStopLeave]:", direction, scaleX, scaleY, duration)
		err = m.handleTouchEdgeMoveStopLeave(direction, scaleX, scaleY, duration)
		if err != nil {
			logger.Error("handleTouchEdgeMoveStopLeave failed:", err)
		}
	})
	if err != nil {
		logger.Error("connect TouchEdgeMoveStopLeave failed:", err)
	}

	m.listenGSettingsChanged()
}

func (m *Manager) getMatchGesture(direction, actionType string) bool {
	lastAc := m.getLastAction()
	isShowMultiTask, err := m.wm.GetMultiTaskingStatus(0)
	if err != nil {
		logger.Debug("GetMultiTaskingStatus err: ",err)
		return false
	}
	isShowDesktop, err := m.wm.GetIsShowDesktop(0)
	if err != nil {
		logger.Debug("GetIsShowDesktop err: ",err)
		return false
	}

	if direction == "down" {
		// 不存在桌面窗口，下滑执行显示桌面
		if !isShowDesktop && actionType == ActionTypeCommandline {
			m.setLastAction(lastActionShowDesktop)
			return true
		}
		// 处于多任务视图，下滑隐藏多任务视图
		if isShowMultiTask && actionType == ActionTypeBuiltin {
			m.setLastAction(lastActionHideWorkspace)
			return true
		}
	}

	if direction == "up" && !isShowMultiTask {
		// 不存在桌面窗口，上滑显示多任务视图
		// 最后一次没有操作，直接上滑显示多任务视图
		// 最后一次操作是对多任务视图操作（显示/隐藏），上滑显示多任务视图
		if (!isShowDesktop || lastAc < lastActionShowDesktop) && actionType == ActionTypeBuiltin {
			m.setLastAction(lastActionShowWorkspace)
			return true
		}
		// 存在桌面窗口，并且上一次操作为显示桌面，上滑隐藏桌面
		if isShowDesktop && lastAc >= lastActionShowDesktop && actionType == ActionTypeCommandline {
			m.setLastAction(lastActionHideDesktop)
			return true
		}
	}
	return false
}

func (m *Manager) getMatchGestureInfo(name, direction string, fingers int32) (*gestureInfo, error) {
	infoArray := m.Infos.Get(name, direction, fingers)
	if len(infoArray) == 0 {
		return nil, fmt.Errorf("not found gesture info for: %s, %s, %d", name, direction, fingers)
	}
	if len(infoArray) == 1 {
		return infoArray[0], nil
	}
	for _, info := range infoArray {
		// 四五指手势相同，需要依据上次操作判断当前操作
		if info.Name == "swipe" && (info.Fingers == 4 || info.Fingers == 5) {
			if m.getMatchGesture(info.Direction, info.Action.Type) {
				return info, nil
			}
		}
	}
	return nil, fmt.Errorf("not found gesture info for: %s, %s, %d", name, direction, fingers)
}

func (m *Manager) Exec(name, direction string, fingers int32) error {
	isEnable := m.getEnable()
	if !isEnable || !isSessionActive() {
		logger.Debug("Gesture had been disabled or session inactive")
		return nil
	}

	info, err := m.getMatchGestureInfo(name, direction, fingers)
	if err != nil {
		logger.Debug("getMatchGestureInfo error:",err)
		return err
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
	} else if strings.HasPrefix(info.Name, "touch") {
		return m.handleTouchScreenEvent(info)
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

//handle touchscreen event
func (*Manager) handleTouchScreenEvent(info *gestureInfo) error {
	return nil
}

//param @edge: swipe to touchscreen edge
func (m *Manager) handleTouchEdgeMoveStopLeave(edge string, scaleX float64, scaleY float64, duration int32) error {
	if edge == "bot" {
		position, err := m.dock.Position().Get(0)
		if err != nil {
			logger.Error("get dock.Position failed:", err)
			return err
		}

		if position >= 0 {
			rect, err := m.dock.FrontendWindowRect().Get(0)
			if err != nil {
				logger.Error("get dock.FrontendWindowRect failed:", err)
				return err
			}

			var dockPly uint32 = 0
			if position == positionTop || position == positionBottom {
				dockPly = rect.Height
			} else if position == positionRight || position == positionLeft {
				dockPly = rect.Width
			}

			screenHeight, err := m.display.ScreenHeight().Get(0)
			if err != nil {
				logger.Error("get display.ScreenHeight failed:", err)
				return err
			}

			if screenHeight > 0 && float64(dockPly)/float64(screenHeight) < scaleY {
				return m.handleBuiltinAction("ShowWorkspace")
			}
		}
	}
	return nil
}
