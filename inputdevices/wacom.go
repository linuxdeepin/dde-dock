/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
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
	wacomSchema = "com.deepin.dde.wacom"

	wacomKeyLeftHanded        = "left-handed"
	wacomKeyCursorMode        = "cursor-mode"
	wacomKeyUpAction          = "keyup-action"
	wacomKeyDownAction        = "keydown-action"
	wacomKeyDoubleDelta       = "double-delta"
	wacomKeyPressureSensitive = "pressure-sensitive"
	wacomKeyMapOutput         = "map-output"
	wacomKeyRawSample         = "raw-sample"
	wacomKeyThreshold         = "threshold"
)

const (
	btnNumUpKey   int32 = 3
	btnNumDownKey       = 2
)

var actionMap = map[string]string{
	"LeftClick":   "button 1",
	"MiddleClick": "button 2",
	"RightClick":  "button 3",
	"PageUp":      "key KP_Page_Up",
	"PageDown":    "key KP_Page_Down",
}

// Soften(x1<y1 x2<y2) --> Firmer(x1>y1 x2>y2)
var pressureLevel = map[uint32][]int{
	1:  []int{0, 100, 0, 100},
	2:  []int{20, 80, 20, 80},
	3:  []int{30, 70, 30, 70},
	4:  []int{0, 0, 100, 100},
	5:  []int{60, 40, 60, 40},
	6:  []int{70, 30, 70, 30}, // default
	7:  []int{75, 25, 75, 25},
	8:  []int{80, 20, 80, 20},
	9:  []int{90, 10, 90, 10},
	10: []int{100, 0, 100, 0},
}

type ActionInfo struct {
	Action string
	Desc   string
}
type ActionInfos []*ActionInfo

type Wacom struct {
	LeftHanded *property.GSettingsBoolProperty `access:"readwrite"`
	CursorMode *property.GSettingsBoolProperty `access:"readwrite"`

	KeyUpAction   *property.GSettingsStringProperty `access:"readwrite"`
	KeyDownAction *property.GSettingsStringProperty `access:"readwrite"`
	MapOutput     *property.GSettingsStringProperty `access:"readwrite"`

	DoubleDelta       *property.GSettingsUintProperty `access:"readwrite"`
	PressureSensitive *property.GSettingsUintProperty `access:"readwrite"`
	RawSample         *property.GSettingsUintProperty `access:"readwrite"`
	Threshold         *property.GSettingsUintProperty `access:"readwrite"`

	DeviceList  string
	ActionInfos ActionInfos
	Exist       bool

	devInfos dxWacoms
	setting  *gio.Settings
}

var _wacom *Wacom

func getWacom() *Wacom {
	if _wacom == nil {
		_wacom = NewWacom()

		_wacom.init()
		_wacom.handleGSettings()
	}

	return _wacom
}

func NewWacom() *Wacom {
	var w = new(Wacom)

	w.setting = gio.NewSettings(wacomSchema)
	w.LeftHanded = property.NewGSettingsBoolProperty(
		w, "LeftHanded",
		w.setting, wacomKeyLeftHanded)
	w.CursorMode = property.NewGSettingsBoolProperty(
		w, "CursorMode",
		w.setting, wacomKeyCursorMode)

	w.KeyUpAction = property.NewGSettingsStringProperty(
		w, "KeyUpAction",
		w.setting, wacomKeyUpAction)
	w.KeyDownAction = property.NewGSettingsStringProperty(
		w, "KeyDownAction",
		w.setting, wacomKeyDownAction)
	w.MapOutput = property.NewGSettingsStringProperty(
		w, "MapOutput",
		w.setting, wacomKeyMapOutput)

	w.DoubleDelta = property.NewGSettingsUintProperty(
		w, "DoubleDelta",
		w.setting, wacomKeyDoubleDelta)
	w.PressureSensitive = property.NewGSettingsUintProperty(
		w, "PressureSensitive",
		w.setting, wacomKeyPressureSensitive)
	w.RawSample = property.NewGSettingsUintProperty(
		w, "RawSample",
		w.setting, wacomKeyRawSample)
	w.Threshold = property.NewGSettingsUintProperty(
		w, "Threshold",
		w.setting, wacomKeyThreshold)

	w.updateDXWacoms()

	return w
}

func (w *Wacom) init() {
	if !w.Exist {
		return
	}

	w.enableCursorMode()
	w.enableLeftHanded()
	w.setKeyAction(btnNumUpKey, w.KeyUpAction.Get())
	w.setKeyAction(btnNumDownKey, w.KeyDownAction.Get())
	w.setPressure()
	w.setClickDelta()
	w.mapToOutput()
	w.setRawSample()
	w.setThreshold()
}

func (w *Wacom) handleDeviceChanged() {
	w.updateDXWacoms()
	w.init()
}

func (w *Wacom) updateDXWacoms() {
	for _, info := range getWacomInfos(false) {
		tmp := w.devInfos.get(info.Id)
		if tmp != nil {
			continue
		}
		w.devInfos = append(w.devInfos, info)
	}

	var v string
	if len(w.devInfos) == 0 {
		w.setPropExist(false)
	} else {
		w.setPropExist(true)
		v = w.devInfos.string()
	}
	setPropString(w, &w.DeviceList, "DeviceList", v)
}

func (w *Wacom) setKeyAction(btnNum int32, action string) {
	value, ok := actionMap[action]
	if !ok {
		return
	}

	for _, v := range w.devInfos {
		err := v.SetButton(int(btnNum), value)
		if err != nil {
			logger.Debugf("Set btn mapping for '%v - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (w *Wacom) enableLeftHanded() {
	var rotate string = "none"
	if w.LeftHanded.Get() {
		rotate = "half"
	}

	for _, v := range w.devInfos {
		err := v.SetRotate(rotate)
		if err != nil {
			logger.Debugf("Enable left handed for '%v - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (w *Wacom) enableCursorMode() {
	var mode string = "Absolute"
	if w.CursorMode.Get() {
		mode = "Relative"
	}

	for _, v := range w.devInfos {
		err := v.SetMode(mode)
		if err != nil {
			logger.Debugf("Enable cursor mode for '%v - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (w *Wacom) setPressure() {
	array, ok := pressureLevel[w.PressureSensitive.Get()]
	if !ok {
		return
	}

	for _, v := range w.devInfos {
		err := v.SetPressureCurve(array[0], array[1], array[2], array[3])
		if err != nil {
			logger.Debugf("Set pressure curve for '%v - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (w *Wacom) setClickDelta() {
	delta := int(w.DoubleDelta.Get())
	for _, v := range w.devInfos {
		err := v.SetSuppress(delta)
		if err != nil {
			logger.Debugf("Set double click delta for '%v - %v' to %v failed: %v",
				v.Id, v.Name, delta, err)
		}
	}
}

func (w *Wacom) mapToOutput() {
	output := w.MapOutput.Get()
	for _, v := range w.devInfos {
		err := v.MapToOutput(output)
		if err != nil {
			logger.Debugf("Map output for '%v - %v' to %v failed: %v",
				v.Id, v.Name, output, err)
		}
	}
}

func (w *Wacom) setRawSample() {
	sample := w.RawSample.Get()
	for _, v := range w.devInfos {
		err := v.SetRawSample(sample)
		if err != nil {
			logger.Debugf("Set raw sample for '%v - %v' to %v failed: %v",
				v.Id, v.Name, sample, err)
		}
	}
}

func (w *Wacom) setThreshold() {
	thres := int(w.Threshold.Get())
	for _, v := range w.devInfos {
		err := v.SetThreshold(thres)
		if err != nil {
			logger.Debugf("Set threshold for '%v - %v' to %v failed: %v",
				v.Id, v.Name, thres, err)
		}
	}
}
