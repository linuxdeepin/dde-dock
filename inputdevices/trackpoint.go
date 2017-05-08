/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package inputdevices

import (
	"gir/gio-2.0"
	"pkg.deepin.io/lib/dbus/property"
)

const (
	trackPointSchema              = "com.deepin.dde.trackpoint"
	trackPointKeyMidButton        = "middle-button-enabled"
	trackPointKeyMidButtonTimeout = "middle-button-timeout"
	trackPointKeyWheel            = "wheel-emulation"
	trackPointKeyWheelButton      = "wheel-emulation-button"
	trackPointKeyWheelTimeout     = "wheel-emulation-timeout"
	trackPointKeyWheelHorizScroll = "wheel-horiz-scroll"
	trackPointKeyAcceleration     = "motion-acceleration"
	trackPointKeyThreshold        = "motion-threshold"
	trackPointKeyScaling          = "motion-scaling"
	trackPointKeyLeftHanded       = "left-handed"
)

type TrackPoint struct {
	MiddleButtonEnabled *property.GSettingsBoolProperty `access:"readwrite"`
	WheelEmulation      *property.GSettingsBoolProperty `access:"readwrite"`
	WheelHorizScroll    *property.GSettingsBoolProperty `access:"readwrite"`

	MiddleButtonTimeout   *property.GSettingsIntProperty `access:"readwrite"`
	WheelEmulationButton  *property.GSettingsIntProperty `access:"readwrite"`
	WheelEmulationTimeout *property.GSettingsIntProperty `access:"readwrite"`

	MotionAcceleration *property.GSettingsFloatProperty `access:"readwrite"`
	MotionThreshold    *property.GSettingsFloatProperty `access:"readwrite"`
	MotionScaling      *property.GSettingsFloatProperty `access:"readwrite"`

	LeftHanded *property.GSettingsBoolProperty `access:"readwrite"`

	DeviceList string
	Exist      bool

	devInfos dxMouses
	setting  *gio.Settings
}

var _trackpoint *TrackPoint

func getTrackPoint() *TrackPoint {
	if _trackpoint == nil {
		_trackpoint = NewTrackPoint()
	}

	return _trackpoint
}

func NewTrackPoint() *TrackPoint {
	var tp = new(TrackPoint)

	tp.setting = gio.NewSettings(trackPointSchema)
	tp.MiddleButtonEnabled = property.NewGSettingsBoolProperty(
		tp, "MiddleButtonEnabled",
		tp.setting, trackPointKeyMidButton)
	tp.WheelEmulation = property.NewGSettingsBoolProperty(
		tp, "WheelEmulation",
		tp.setting, trackPointKeyWheel)
	tp.WheelHorizScroll = property.NewGSettingsBoolProperty(
		tp, "WheelHorizScroll",
		tp.setting, trackPointKeyWheelHorizScroll)

	tp.MotionAcceleration = property.NewGSettingsFloatProperty(
		tp, "MotionAcceleration",
		tp.setting, trackPointKeyAcceleration)
	tp.MotionThreshold = property.NewGSettingsFloatProperty(
		tp, "MotionThreshold",
		tp.setting, trackPointKeyThreshold)
	tp.MotionScaling = property.NewGSettingsFloatProperty(
		tp, "MotionScaling",
		tp.setting, trackPointKeyScaling)

	tp.MiddleButtonTimeout = property.NewGSettingsIntProperty(
		tp, "MiddleButtonTimeout",
		tp.setting, trackPointKeyMidButtonTimeout)
	tp.WheelEmulationButton = property.NewGSettingsIntProperty(
		tp, "WheelEmulationButton",
		tp.setting, trackPointKeyWheelButton)
	tp.WheelEmulationTimeout = property.NewGSettingsIntProperty(
		tp, "WheelEmulationTimeout",
		tp.setting, trackPointKeyWheelTimeout)

	tp.LeftHanded = property.NewGSettingsBoolProperty(
		tp, "LeftHanded",
		tp.setting, trackPointKeyLeftHanded)

	tp.updateDXMouses()

	return tp
}

func (tp *TrackPoint) init() {
	if !tp.Exist {
		return
	}

	tp.enableMiddleButton()
	tp.enableWheelEmulation()
	tp.enableWheelHorizScroll()
	tp.enableLeftHanded()
	tp.middleButtonTimeout()
	tp.wheelEmulationButton()
	tp.wheelEmulationTimeout()
	tp.motionAcceleration()
	tp.motionThreshold()
	tp.motionScaling()
}

func (tp *TrackPoint) handleDeviceChanged() {
	tp.updateDXMouses()
	tp.init()
}

func (tp *TrackPoint) updateDXMouses() {
	tp.devInfos = dxMouses{}
	for _, info := range getMouseInfos(false) {
		if !info.TrackPoint {
			continue
		}

		tmp := tp.devInfos.get(info.Id)
		if tmp != nil {
			continue
		}
		tp.devInfos = append(tp.devInfos, info)
	}

	var v string
	if len(tp.devInfos) == 0 {
		tp.setPropExist(false)
	} else {
		tp.setPropExist(true)
		v = tp.devInfos.string()
	}
	setPropString(tp, &tp.DeviceList, "DeviceList", v)
}

func (tp *TrackPoint) enableMiddleButton() {
	enabled := tp.MiddleButtonEnabled.Get()
	for _, info := range tp.devInfos {
		err := info.EnableMiddleButtonEmulation(enabled)
		if err != nil {
			logger.Warningf("Enable middle button for '%v %s' failed: %v",
				info.Id, info.Name, err)
		}
	}
}

func (tp *TrackPoint) enableWheelEmulation() {
	enabled := tp.WheelEmulation.Get()
	for _, info := range tp.devInfos {
		err := info.EnableWheelEmulation(enabled)
		if err != nil {
			logger.Warningf("Enable wheel emulation for '%v %s' failed: %v",
				info.Id, info.Name, err)
		}
	}
}

func (tp *TrackPoint) enableWheelHorizScroll() {
	enabled := tp.WheelHorizScroll.Get()
	for _, info := range tp.devInfos {
		err := info.EnableWheelHorizScroll(enabled)
		if err != nil {
			logger.Warningf("Enable wheel horiz scroll for '%v %s' failed: %v",
				info.Id, info.Name, err)
		}
	}
}

func (tp *TrackPoint) enableLeftHanded() {
	enabled := tp.LeftHanded.Get()
	for _, info := range tp.devInfos {
		err := info.EnableLeftHanded(enabled)
		if err != nil {
			logger.Warningf("Enable left-handed for '%v %s' failed: %v",
				info.Id, info.Name, err)
		}
	}
}

func (tp *TrackPoint) middleButtonTimeout() {
	timeout := tp.MiddleButtonTimeout.Get()
	for _, info := range tp.devInfos {
		err := info.SetMiddleButtonEmulationTimeout(int16(timeout))
		if err != nil {
			logger.Warningf("Set middle button timeout for '%v %s' failed: %v",
				info.Id, info.Name, err)
		}
	}
}

func (tp *TrackPoint) wheelEmulationButton() {
	button := tp.WheelEmulationButton.Get()
	for _, info := range tp.devInfos {
		err := info.SetWheelEmulationButton(int8(button))
		if err != nil {
			logger.Warningf("Set wheel button for '%v %s' failed: %v",
				info.Id, info.Name, err)
		}
	}
}

func (tp *TrackPoint) wheelEmulationTimeout() {
	timeout := tp.WheelEmulationTimeout.Get()
	for _, info := range tp.devInfos {
		err := info.SetWheelEmulationTimeout(int16(timeout))
		if err != nil {
			logger.Warningf("Enable wheel timeout for '%v %s' failed: %v",
				info.Id, info.Name, err)
		}
	}
}

func (tp *TrackPoint) motionAcceleration() {
	accel := float32(tp.MotionAcceleration.Get())
	for _, v := range tp.devInfos {
		err := v.SetMotionAcceleration(accel)
		if err != nil {
			logger.Debugf("Set acceleration for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (tp *TrackPoint) motionThreshold() {
	thres := float32(tp.MotionThreshold.Get())
	for _, v := range tp.devInfos {
		err := v.SetMotionThreshold(thres)
		if err != nil {
			logger.Debugf("Set threshold for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (tp *TrackPoint) motionScaling() {
	scaling := float32(tp.MotionScaling.Get())
	for _, v := range tp.devInfos {
		err := v.SetMotionScaling(scaling)
		if err != nil {
			logger.Debugf("Set scaling for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}
