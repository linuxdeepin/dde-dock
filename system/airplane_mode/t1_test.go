package airplane_mode

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"pkg.deepin.io/lib/log"
)

func assertIt(t *testing.T, opId int, s *AirplaneModeState, enabled, wifiEnabled, btEnable bool) {
	assert.Equal(t, enabled, s.Enabled, "op %d enabled", opId)
	assert.Equal(t, wifiEnabled, s.WifiEnabled, "op %d wifi", opId)
	assert.Equal(t, btEnable, s.BtEnabled, "op %d bt", opId)
}

func init() {
	logger.SetLogLevel(log.LevelDebug)
}

func (s *AirplaneModeState) setupForTest() {
	s.enableWifiFn = func(enabled bool) {
		s.WifiEnabled = enabled
	}
	s.enableBtFn = func(enabled bool) {
		s.BtEnabled = enabled
	}
}

func TestScene1(t *testing.T) {
	// scene 1
	s := &AirplaneModeState{
		WifiEnabled: true,
		BtEnabled:   true,
	}
	s.setupForTest()
	// op 1 飞行模式开启
	s.enable(true)
	assertIt(t, 1, s, true, false, false)
	// op 2 飞行模式关闭
	s.enable(false)
	assertIt(t, 2, s, false, true, true)
}

func TestScene2(t *testing.T) {
	// scene 2
	s := &AirplaneModeState{
		WifiEnabled: true,
		BtEnabled:   true,
	}
	s.setupForTest()
	// op 1 飞行模式开启
	s.enable(true)
	assertIt(t, 1, s, true, false, false)
	// op 2 wifi 开启
	s.enableWifi(true)
	assertIt(t, 2, s, true, true, false)
	// op 3 飞行模式关闭
	s.enable(false)
	assertIt(t, 3, s, false, true, true)
}

func TestScene3(t *testing.T) {
	// scene 3
	s := &AirplaneModeState{
		WifiEnabled: true,
		BtEnabled:   true,
	}
	s.setupForTest()
	// op 1 飞行模式开启
	s.enable(true)
	assertIt(t, 1, s, true, false, false)
	// op 2 BT 开启
	s.enableBt(true)
	assertIt(t, 2, s, true, false, true)
	// op 3 飞行模式关闭
	s.enable(false)
	assertIt(t, 3, s, false, true, true)
}

func TestScene4(t *testing.T) {
	// scene 4
	s := &AirplaneModeState{
		WifiEnabled: true,
		BtEnabled:   true,
	}
	s.setupForTest()
	// op 1 飞行模式开启
	s.enable(true)
	assertIt(t, 1, s, true, false, false)
	// op 2 wifi, BT 开启
	s.enableWifi(true)
	s.enableBt(true)
	assertIt(t, 2, s, true, true, true)
	// op 3 飞行模式关闭
	s.enable(false)
	assertIt(t, 3, s, false, true, true)
}

func TestScene5(t *testing.T) {
	// scene 5
	s := &AirplaneModeState{
		WifiEnabled: true,
		BtEnabled:   true,
	}
	s.setupForTest()
	// op 1 wifi, BT 关闭
	s.enableWifi(false)
	s.enableBt(false)
	assertIt(t, 1, s, false, false, false)
	// op 2 飞行模式开启
	s.enable(true)
	assertIt(t, 2, s, true, false, false)
}

func TestScene6(t *testing.T) {
	// scene 6
	s := &AirplaneModeState{
		WifiEnabled: false,
		BtEnabled:   true,
	}
	s.setupForTest()
	// op 1 飞行模式开启
	s.enable(true)
	assertIt(t, 1, s, true, false, false)
	// op 2 飞行模式关闭
	s.enable(false)
	assertIt(t, 2, s, false, false, true)
}

func TestScene7(t *testing.T) {
	// scene 7
	s := &AirplaneModeState{
		WifiEnabled: true,
		BtEnabled:   false,
	}
	s.setupForTest()
	// op 1 飞行模式开启
	s.enable(true)
	assertIt(t, 1, s, true, false, false)
	// op 2 飞行模式关闭
	s.enable(false)
	assertIt(t, 2, s, false, true, false)
}

func TestScene8_1(t *testing.T) {
	// scene 8.1
	s := &AirplaneModeState{
		WifiEnabled: false,
		BtEnabled:   false,
	}
	s.setupForTest()
	// op 1 飞行模式开启
	s.enable(true)
	assertIt(t, 1, s, true, false, false)
	// op 2 飞行模式关闭
	s.enable(false)
	assertIt(t, 2, s, false, false, false)
}

func TestScene8_2(t *testing.T) {
	// scene 8.2
	s := &AirplaneModeState{
		WifiEnabled: false,
		BtEnabled:   false,
	}
	s.setupForTest()
	// op 1 飞行模式开启
	s.enable(true)
	assertIt(t, 1, s, true, false, false)
	// op 2 wifi 开启
	s.enableWifi(true)
	assertIt(t, 2, s, true, true, false)
	// op 3 飞行模式关闭
	s.enable(false)
	assertIt(t, 3, s, false, true, false)
}

func TestScene8_3(t *testing.T) {
	// scene 8.3
	s := &AirplaneModeState{
		WifiEnabled: false,
		BtEnabled:   false,
	}
	s.setupForTest()
	// op 1 飞行模式开启
	s.enable(true)
	assertIt(t, 1, s, true, false, false)
	// op 2 BT 开启
	s.enableBt(true)
	assertIt(t, 2, s, true, false, true)
	// op 3 飞行模式关闭
	s.enable(false)
	assertIt(t, 3, s, false, false, true)
}

func TestScene8_4(t *testing.T) {
	// scene 8.4
	s := &AirplaneModeState{
		WifiEnabled: false,
		BtEnabled:   false,
	}
	s.setupForTest()
	// op 1 飞行模式开启
	s.enable(true)
	assertIt(t, 1, s, true, false, false)
	// op 2 wifi、BT 开启
	s.enableWifi(true)
	s.enableBt(true)
	assertIt(t, 2, s, true, true, true)
	// op 3 飞行模式关闭
	s.enable(false)
	assertIt(t, 3, s, false, true, true)
}
