// +build sw

/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func init() {
	submoduleList = append(submoduleList, newBatteryStateListener)
}

type batteryStateListener struct {
	manager *Manager
	ticker  *time.Ticker
}

func readBatteryInfoFile(file string) (map[string]string, error) {
	content, err := ioutil.ReadFile(file)
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

	return data, nil
}

func getStringValue(data map[string]string, key string) (string, error) {
	valStr, ok := data[key]
	if !ok {
		return "", fmt.Errorf("no this key %q", key)
	}
	return valStr, nil
}

func getFloatValue(data map[string]string, key string) (float64, error) {
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

func getInfoFromPowerSupplyUEvent(battery *batteryDevice) (*batteryInfo, error) {
	nativePath := battery.NativePath.Get()
	file := "/sys/class/power_supply/" + nativePath + "/uevent"
	logger.Debug("getInfoFromPowerSupplyUEvent", nativePath)
	data, err := readBatteryInfoFile(file)
	if err != nil {
		return nil, err
	}

	isPresentStr, err := getStringValue(data, "POWER_SUPPLY_PRESENT")
	if err != nil {
		return nil, err
	}
	isPresent := (isPresentStr == "1")

	status, err := getStringValue(data, "POWER_SUPPLY_STATUS")
	if err != nil {
		return nil, err
	}
	state, ok := batteryStateMap[status]
	if !ok {
		return nil, fmt.Errorf("Invalid battery state  %q", status)
	}
	logger.Debugf("status: %v, state %v", status, state)

	percentage, err := getFloatValue(data, "POWER_SUPPLY_CAPACITY")
	if err != nil {
		return nil, err
	}

	energyFullDesign, err := getFloatValue(data, "POWER_SUPPLY_ENERGY_FULL_DESIGN")
	if err != nil {
		return nil, err
	}

	energyFull, err := getFloatValue(data, "POWER_SUPPLY_ENERGY_FULL")
	if err != nil {
		return nil, err
	}

	energy, err := getFloatValue(data, "POWER_SUPPLY_ENERGY_NOW")
	if err != nil {
		return nil, err
	}

	const unit = 1000000.0
	return &batteryInfo{
		IsPresent:        isPresent,
		State:            state,
		Percentage:       percentage,
		EnergyFullDesign: energyFullDesign / unit,
		EnergyFull:       energyFull / unit,
		Energy:           energy / unit,
	}, nil
}

func newBatteryStateListener(m *Manager) (string, submodule, error) {
	m.batteryGroup.getInfoMethod = getInfoFromPowerSupplyUEvent
	// force use percentage for policy
	m.usePercentageForPolicy = true
	bsl := &batteryStateListener{
		manager: m,
		ticker:  time.NewTicker(time.Second * 5),
	}
	return "BatteryStateListener", bsl, nil
}

func (bsl *batteryStateListener) Start() error {
	go func() {
		for range bsl.ticker.C {
			logger.Debug("battery state tick")
			bsl.manager.updateBatteryGroupInfo()
		}
	}()
	return nil
}

func (bsl *batteryStateListener) Destroy() {
	bsl.ticker.Stop()
}
