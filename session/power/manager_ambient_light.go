package power

import (
	"sort"
)

func (m *Manager) claimOrReleaseAmbientLight() {
	logger.Debug("call claimOrReleaseAmbientLight")
	var shouldClaim bool

	autoAdjustEnabled := m.AmbientLightAdjustBrightness.Get()

	m.PropsMu.RLock()
	if m.HasAmbientLightSensor &&
		m.lightLevelUnit == "lux" &&
		autoAdjustEnabled &&
		m.sessionActive {
		shouldClaim = true
	}

	if m.lidSwitchState == lidSwitchStateClose {
		shouldClaim = false
	}

	logger.Debugf("hasAmbientLightSensor: %v, autoAdjustEnabled: %v,"+
		" session active: %v, lidSwitchState: %v -> shouldClaim: %v",
		m.HasAmbientLightSensor, autoAdjustEnabled, m.sessionActive,
		m.lidSwitchState, shouldClaim)

	m.PropsMu.RUnlock()

	if shouldClaim {
		m.claimAmbientLight()
	} else {
		m.releaseAmbientLight()
	}
}

func (m *Manager) claimAmbientLight() {
	m.PropsMu.Lock()
	defer m.PropsMu.Unlock()

	if m.ambientLightClaimed {
		return
	}

	logger.Debug("claim ambient light")
	err := m.helper.SensorProxy.ClaimLight(0)
	if err != nil {
		logger.Warning(err)
		return
	}

	m.ambientLightClaimed = true

	lightLevel, _ := m.helper.SensorProxy.LightLevel().Get(0)
	m.handleLightLevelChanged(lightLevel)
}

func (m *Manager) releaseAmbientLight() {
	m.PropsMu.Lock()
	defer m.PropsMu.Unlock()

	if !m.ambientLightClaimed {
		return
	}

	logger.Debug("release ambient light")
	err := m.helper.SensorProxy.ReleaseLight(0)
	if err != nil {
		logger.Warning(err)
		return
	}
	m.ambientLightClaimed = false
}

func (m *Manager) handleLightLevelChanged(lightLevel float64) {
	if !m.AmbientLightAdjustBrightness.Get() {
		return
	}

	if lightLevel <= 0 {
		logger.Warning("invalid light level:", lightLevel)
		return
	}
	logger.Debug("light level changed to", lightLevel)

	display := m.helper.Display
	outputNames, err := display.ListOutputNames(0)
	if err != nil {
		logger.Warning(err)
		return
	}

	//outputNames
	var builtinOutputName string
	for _, name := range outputNames {
		if isBuiltinOutput(name) {
			builtinOutputName = name
			break
		}
	}

	if builtinOutputName == "" {
		// not found builtin output
		return
	}

	br := float64(calcBrWithLightLevel(lightLevel)) / 255
	logger.Debugf("auto set brightness to %v\n", br)
	err = display.SetBrightness(0, builtinOutputName, br)
	if err != nil {
		logger.Warning("failed to set brightness:", err)
	}
}

type lightLevelBr struct {
	lightLevel int // unit lux
	brightness byte
}

var lightLevelBrTable = []lightLevelBr{
	{0, 0},
	{1, 2},
	{2, 3},
	{3, 5},
	{4, 7},
	{5, 8},
	{6, 10},
	{7, 12},
	{8, 14},
	{9, 15},
	{10, 17},
	{11, 19},
	{12, 20},
	{13, 22},
	{14, 24},
	{15, 25},
	{16, 27},
	{17, 29},
	{18, 31},
	{19, 32},
	{20, 34},
	{21, 36},
	{22, 37},
	{23, 39},
	{24, 41},
	{25, 42},
	{30, 47},
	{50, 48},
	{100, 50},
	{200, 55},
	{500, 69},
	{600, 74},
	{700, 79},
	{800, 84},
	{900, 88},
	{1000, 93},
	{1100, 98},
	{1200, 103},
	{1300, 107},
	{1400, 112},
	{1500, 117},
	{1600, 122},
	{1700, 127},
	{1800, 131},
	{1900, 136},
	{1995, 141},
	{1996, 141},
	{2000, 141},
	{2100, 153},
	{2200, 164},
	{2300, 175},
	{2400, 187},
	{2500, 198},
	{2600, 210},
	{2700, 221},
	{2800, 232},
	{2900, 244},
	{3000, 255},
	{3100, 255},
	{3200, 255},
	{3300, 255},
	{3400, 255},
	{3500, 255},
	{3600, 255},
	{3700, 255},
	{3800, 255},
	{3900, 255},
	{4000, 255},
	{4100, 255},
	{4200, 255},
	{4300, 255},
	{4400, 255},
	{4500, 255},
	{4600, 255},
	{4700, 255},
	{4800, 255},
	{4900, 255},
	{5000, 255},
	{6000, 255},
	{7000, 255},
	{8000, 255},
	{9000, 255},
	{10000, 255},
}

func calcBrWithLightLevel(lightLevel float64) byte {
	if lightLevel > 10000 {
		return 255
	}
	if lightLevel < 0 {
		return 0
	}
	lightLv := int(lightLevel)

	i := sort.Search(len(lightLevelBrTable), func(i int) bool {
		return lightLevelBrTable[i].lightLevel >= lightLv
	})
	if i < len(lightLevelBrTable) && lightLevelBrTable[i].lightLevel == lightLv {
		return lightLevelBrTable[i].brightness
	} else {
		x1 := float64(lightLevelBrTable[i-1].lightLevel)
		y1 := float64(lightLevelBrTable[i-1].brightness)

		x2 := float64(lightLevelBrTable[i].lightLevel)
		y2 := float64(lightLevelBrTable[i].brightness)

		x := float64(lightLv)
		y := (x-x1)/(x2-x1)*(y2-y1) + y1
		return byte(y)
	}
}
