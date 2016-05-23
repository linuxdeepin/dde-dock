package power

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func newBatteryDeviceWithFile(file string) (batteryDevice, error) {
	typeFileContent, err := ioutil.ReadFile(filepath.Join(file, "type"))
	if err != nil {
		return nil, err
	}

	typeStr := strings.TrimSpace(strings.ToLower(string(typeFileContent)))
	logger.Debugf("type is %q", typeStr)
	if typeStr != "battery" {
		return nil, errors.New("device type is not battery")
	}
	batDevice, err := newPollBatteryDevice(filepath.Join(file, "uevent"))
	if err != nil {
		return nil, err
	}
	return batDevice, nil
}

func (batGroup *batteryGroup) initPollBatteryDevices() {
	files, err := ioutil.ReadDir(sysPowerSupplyDir)
	if err != nil {
		logger.Warning(err)
		return
	}

	for _, fileInfo := range files {
		logger.Debugf("fileInfo: %#v", fileInfo)
		batDevice, err := newBatteryDeviceWithFile(
			filepath.Join(sysPowerSupplyDir, fileInfo.Name()))
		if err == nil {
			batGroup.Add(batDevice)
		} else {
			logger.Debug(err)
		}
	}
}

func (batGroup *batteryGroup) initPollBackend() {
	batGroup.initPollBatteryDevices()
}
