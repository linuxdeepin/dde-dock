package power

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type pollBatteryDevice struct {
	ueventFile string
	info       *batteryInfo
	ticker     *time.Ticker
}

func newPollBatteryDevice(file string) (*pollBatteryDevice, error) {
	// check file exist
	_, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	batDevice := &pollBatteryDevice{
		ueventFile: file,
	}
	logger.Debugf("newPollBatteryDevice file: %q", file)
	return batDevice, nil
}

func (dev *pollBatteryDevice) Destroy() {
	if dev.ticker != nil {
		dev.ticker.Stop()
	}
}

func (dev *pollBatteryDevice) GetPath() string {
	return "file://" + dev.ueventFile
}

func (dev *pollBatteryDevice) GetInfo() *batteryInfo {
	return dev.info
}

type batteryInfoMap map[string]string

func (dev *pollBatteryDevice) readBatteryInfoFile() (batteryInfoMap, error) {
	content, err := ioutil.ReadFile(dev.ueventFile)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var data = make(map[string]string)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		array := strings.Split(line, "=")
		if len(array) != 2 {
			continue
		}
		data[array[0]] = strings.TrimSpace(array[1])
	}

	return batteryInfoMap(data), nil
}

func (data batteryInfoMap) getStringValue(key string) (string, error) {
	valStr, ok := data[key]
	if !ok {
		return "", fmt.Errorf("no this key %q", key)
	}
	return valStr, nil
}

func (data batteryInfoMap) getFloatValue(key string) (float64, error) {
	valStr, ok := data[key]
	if !ok {
		return 0, fmt.Errorf("no this key %q", key)
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (data batteryInfoMap) getIsPresent() bool {
	isPresentStr, err := data.getStringValue("POWER_SUPPLY_PRESENT")
	if err != nil {
		return false
	}
	return isPresentStr == "1"
}

func (data batteryInfoMap) getState() batteryStateType {
	stateStr, err := data.getStringValue("POWER_SUPPLY_STATUS")
	if err != nil {
		return BatteryStateUnknown
	}
	state, ok := batteryStateMap[stateStr]
	if !ok {
		return BatteryStateUnknown
	}
	return state
}

func (data batteryInfoMap) getPercentage() float64 {
	percentage, _ := data.getFloatValue("POWER_SUPPLY_CAPACITY")
	return percentage
}

const energyUnit = 1000000.0

func (data batteryInfoMap) getEneryFullDesign() float64 {
	energy, _ := data.getFloatValue("POWER_SUPPLY_ENERGY_FULL_DESIGN")
	return energy / energyUnit
}

func (data batteryInfoMap) getEneryFull() float64 {
	energy, _ := data.getFloatValue("POWER_SUPPLY_ENERGY_FULL")
	return energy / energyUnit
}

func (data batteryInfoMap) getEnery() float64 {
	energy, _ := data.getFloatValue("POWER_SUPPLY_ENERGY_NOW")
	return energy / energyUnit
}

func (dev *pollBatteryDevice) fillInfo() {
	bi := dev.info
	batInfoMap, err := dev.readBatteryInfoFile()
	if err != nil {
		logger.Warning(err)
		return
	}
	bi.setIsPresent(batInfoMap.getIsPresent())
	bi.setState(batInfoMap.getState())
	bi.setEnergyFullDesign(batInfoMap.getEneryFullDesign())
	bi.setEnergyFull(batInfoMap.getEneryFull())
	bi.setEnergy(batInfoMap.getEnery())
	bi.setPercentage(batInfoMap.getPercentage())
}

func (dev *pollBatteryDevice) SetInfo(bi *batteryInfo) {
	dev.info = bi
	dev.ticker = time.NewTicker(time.Second * 10)
	go func() {
		for range dev.ticker.C {
			dev.fillInfo()
		}
	}()
	dev.fillInfo()
}
