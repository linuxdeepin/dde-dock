/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

type batteryStateType uint32

const (
	//defined at http://upower.freedesktop.org/docs/Device.html#Device:State
	BatteryStateUnknown          = batteryStateType(0)
	BatteryStateCharging         = batteryStateType(1)
	BatteryStateDischarging      = batteryStateType(2)
	BatteryStateEmpty            = batteryStateType(3)
	BatteryStateFullyCharged     = batteryStateType(4)
	BatteryStatePendingCharge    = batteryStateType(5)
	BatteryStatePendingDischarge = batteryStateType(6)
)

var batteryStateMap = map[string]batteryStateType{
	"Unknown":          BatteryStateUnknown,
	"Charging":         BatteryStateCharging,
	"Discharging":      BatteryStateDischarging,
	"Empty":            BatteryStateEmpty,
	"FullCharged":      BatteryStateFullyCharged,
	"PendingCharge":    BatteryStatePendingCharge,
	"PendingDischarge": BatteryStatePendingDischarge,
}

func (state batteryStateType) String() string {
	switch state {
	case BatteryStateCharging:
		return "Charging"
	case BatteryStateDischarging:
		return "Discharging"
	case BatteryStateEmpty:
		return "Empty"
	case BatteryStateFullyCharged:
		return "FullyCharged"
	case BatteryStatePendingCharge:
		return "PendingCharge"
	case BatteryStatePendingDischarge:
		return "PendingDischarge"
	default:
		return "Unknown"
	}
}
