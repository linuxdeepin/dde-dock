package airplane_mode

import (
	"github.com/davecgh/go-spew/spew"
	"sync"
)

type savedEnabled struct {
	WifiEnabled bool
	WifiChanged bool

	BtEnabled bool
	BtChanged bool
}

type AirplaneModeState struct {
	mu      sync.Mutex
	Enabled bool

	WifiEnabled    bool
	WifiWillChange bool

	BtEnabled    bool
	BtWillChange bool

	Saved *savedEnabled

	enableWifiFn func(enabled bool)
	enableBtFn   func(enabled bool)
}

func (s *AirplaneModeState) dump() {
	s.mu.Lock()
	spew.Dump(s)
	s.mu.Unlock()
}

func (s *AirplaneModeState) doEnableWifi(enabled bool) {
	s.WifiWillChange = true
	s.enableWifiFn(enabled)
}

func (s *AirplaneModeState) doEnableBt(enabled bool) {
	s.BtWillChange = true
	s.enableBtFn(enabled)
}

func (s *AirplaneModeState) hasWifiChanged() bool {
	return s.Saved != nil && s.Saved.WifiChanged
}

func (s *AirplaneModeState) hasBtChanged() bool {
	return s.Saved != nil && s.Saved.BtChanged
}

func toEnableDisable(enable bool) string {
	if enable {
		return "enable"
	}
	return "disable"
}

// 此方法调用后，要改变wifi或bt的状态， 就必须只能调用 doEnableWifi 或 doEnableBt 方法, 将会设置
// WifiWillChange 或 BtWillChange 为 true, 然后期待外部事件去调用 enableWifi 或 enableBt。
func (s *AirplaneModeState) enable(enabled bool) {
	logger.Debug("State enable", enabled)
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Enabled == enabled {
		return
	}
	s.Enabled = enabled
	if enabled {
		logger.Debug("enable flight mode")
		s.Saved = &savedEnabled{
			WifiEnabled: s.WifiEnabled,
			BtEnabled:   s.BtEnabled,
		}

		if s.WifiEnabled {
			logger.Debug("auto disable WIFI")
			s.doEnableWifi(false)
		}
		if s.BtEnabled {
			logger.Debug("auto disable BT")
			s.doEnableBt(false)
		}
	} else {
		// disable
		logger.Debug("disable flight mode")
		if !s.hasWifiChanged() {
			if s.Saved != nil {
				if s.Saved.WifiEnabled != s.WifiEnabled {
					sEnabled := s.Saved.WifiEnabled
					s.doEnableWifi(sEnabled)
					logger.Debugf("auto %s WIFI (recovery)", toEnableDisable(sEnabled))
				}
			} else {
				logger.Debug("auto enable WIFI")
				s.doEnableWifi(true)
			}
		} else {
			logger.Debug("do not change wifi")
		}

		if !s.hasBtChanged() {
			if s.Saved != nil {
				if s.Saved.BtEnabled != s.BtEnabled {
					sEnabled := s.Saved.BtEnabled
					s.doEnableBt(sEnabled)
					logger.Debugf("auto %s BT (recovery)", toEnableDisable(sEnabled))
				}
			} else {
				logger.Debug("auto enable BT")
				s.doEnableBt(true)
			}
		} else {
			logger.Debug("do not change bt")
		}
		s.Saved = nil
	}
}

// 此函数只能由Manager收到wifi状态改变事件后调用
func (s *AirplaneModeState) enableWifi(enabled bool) {
	logger.Debug("State enableWifi", enabled)
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.WifiEnabled == enabled {
		return
	}
	s.WifiEnabled = enabled
	if s.Saved != nil && !s.WifiWillChange {
		logger.Debug("hint wifi changed")
		s.Saved.WifiChanged = true
	}
	s.WifiWillChange = false
}

// 此函数只能由Manager收到Bt状态改变事件后调用
func (s *AirplaneModeState) enableBt(enabled bool) {
	logger.Debug("State enableBt", enabled)
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.BtEnabled == enabled {
		return
	}
	s.BtEnabled = enabled
	if s.Saved != nil && !s.BtWillChange {
		logger.Debug("hint bt changed")
		s.Saved.BtChanged = true
	}
	s.BtWillChange = false
}
