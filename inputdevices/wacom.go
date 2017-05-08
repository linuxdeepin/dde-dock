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
	"fmt"
	"gir/gio-2.0"
	"pkg.deepin.io/dde/api/dxinput"
	"pkg.deepin.io/lib/dbus/property"
	"strings"
)

const (
	wacomSchema       = "com.deepin.dde.wacom"
	wacomStylusSchema = wacomSchema + ".stylus"
	wacomEraserSchema = wacomSchema + ".eraser"

	wacomKeyLeftHanded        = "left-handed"
	wacomKeyCursorMode        = "cursor-mode"
	wacomKeyUpAction          = "keyup-action"
	wacomKeyDownAction        = "keydown-action"
	wacomKeySuppress          = "suppress"
	wacomKeyPressureSensitive = "pressure-sensitive"
	wacomKeyMapOutput         = "map-output"
	wacomKeyRawSample         = "raw-sample"
	wacomKeyThreshold         = "threshold"
)

const (
	btnNumUpKey   = 3
	btnNumDownKey = 2
)

var actionMap = map[string]string{
	"LeftClick":   "button 1",
	"MiddleClick": "button 2",
	"RightClick":  "button 3",
	"PageUp":      "key KP_Page_Up",
	"PageDown":    "key KP_Page_Down",
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

	Suppress                *property.GSettingsUintProperty `access:"readwrite"`
	StylusPressureSensitive *property.GSettingsUintProperty `access:"readwrite"`
	EraserPressureSensitive *property.GSettingsUintProperty `access:"readwrite"`
	StylusRawSample         *property.GSettingsUintProperty `access:"readwrite"`
	EraserRawSample         *property.GSettingsUintProperty `access:"readwrite"`
	StylusThreshold         *property.GSettingsUintProperty `access:"readwrite"`
	EraserThreshold         *property.GSettingsUintProperty `access:"readwrite"`

	DeviceList  string
	ActionInfos ActionInfos
	Exist       bool

	devInfos      dxWacoms
	setting       *gio.Settings
	stylusSetting *gio.Settings
	eraserSetting *gio.Settings
}

var _wacom *Wacom

func getWacom() *Wacom {
	if _wacom == nil {
		_wacom = NewWacom()
	}

	return _wacom
}

func NewWacom() *Wacom {
	var w = new(Wacom)

	w.setting = gio.NewSettings(wacomSchema)
	w.stylusSetting = gio.NewSettings(wacomStylusSchema)
	w.eraserSetting = gio.NewSettings(wacomEraserSchema)

	w.LeftHanded = property.NewGSettingsBoolProperty(
		w, "LeftHanded",
		w.setting, wacomKeyLeftHanded)
	w.CursorMode = property.NewGSettingsBoolProperty(
		w, "CursorMode",
		w.setting, wacomKeyCursorMode)

	w.KeyUpAction = property.NewGSettingsStringProperty(
		w, "KeyUpAction",
		w.stylusSetting, wacomKeyUpAction)
	w.KeyDownAction = property.NewGSettingsStringProperty(
		w, "KeyDownAction",
		w.stylusSetting, wacomKeyDownAction)

	w.MapOutput = property.NewGSettingsStringProperty(
		w, "MapOutput",
		w.setting, wacomKeyMapOutput)

	w.Suppress = property.NewGSettingsUintProperty(
		w, "Suppress",
		w.setting, wacomKeySuppress)

	w.StylusPressureSensitive = property.NewGSettingsUintProperty(
		w, "StylusPressureSensitive",
		w.stylusSetting, wacomKeyPressureSensitive)

	w.EraserPressureSensitive = property.NewGSettingsUintProperty(
		w, "EraserPressureSensitive",
		w.eraserSetting, wacomKeyPressureSensitive)

	w.StylusRawSample = property.NewGSettingsUintProperty(
		w, "StylusRawSample",
		w.stylusSetting, wacomKeyRawSample)

	w.EraserRawSample = property.NewGSettingsUintProperty(
		w, "EraserRawSample",
		w.eraserSetting, wacomKeyRawSample)

	w.StylusThreshold = property.NewGSettingsUintProperty(
		w, "StylusThreshold",
		w.stylusSetting, wacomKeyThreshold)

	w.EraserThreshold = property.NewGSettingsUintProperty(
		w, "EraserThreshold",
		w.eraserSetting, wacomKeyThreshold)

	w.updateDXWacoms()

	return w
}

func (w *Wacom) init() {
	if !w.Exist {
		return
	}

	w.enableCursorMode()
	w.enableLeftHanded()
	w.setStylusButtonAction(btnNumUpKey, w.KeyUpAction.Get())
	w.setStylusButtonAction(btnNumDownKey, w.KeyDownAction.Get())
	w.setPressureSensitive()
	w.setSuppress()
	w.mapToOutput()
	w.setRawSample()
	w.setThreshold()
}

func (w *Wacom) handleDeviceChanged() {
	w.updateDXWacoms()
	w.init()
}

func (w *Wacom) updateDXWacoms() {
	w.devInfos = dxWacoms{}
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

func (w *Wacom) setStylusButtonAction(btnNum int, action string) {
	value, ok := actionMap[action]
	if !ok {
		logger.Warningf("Invalid button action %q, actionMap: %v", action, actionMap)
		return
	}
	// set button action for stylus
	for _, v := range w.devInfos {
		if v.QueryType() == dxinput.WacomTypeStylus {
			err := v.SetButton(btnNum, value)
			if err != nil {
				logger.Warningf("Set button mapping for '%v - %v' failed: %v",
					v.Id, v.Name, err)
			}
		}
	}
}

func (w *Wacom) enableLeftHanded() {
	var rotate string = "none"
	if w.LeftHanded.Get() {
		rotate = "half"
	}
	// set rotate for stylus and eraser
	// Rotation is a tablet-wide option:
	// rotation of one tool affects all other tools associated with the same tablet.
	for _, v := range w.devInfos {
		devType := v.QueryType()
		if devType == dxinput.WacomTypeStylus || devType == dxinput.WacomTypeEraser {
			err := v.SetRotate(rotate)
			if err != nil {
				logger.Warningf("Set rotate for '%v - %v' failed: %v",
					v.Id, v.Name, err)
			}
		}
	}
}

func (w *Wacom) enableCursorMode() {
	var mode string = "Absolute"
	if w.CursorMode.Get() {
		mode = "Relative"
	}

	// set mode for stylus and eraser
	// NOTE: set mode Relative for pad  will cause X error
	for _, v := range w.devInfos {
		devType := v.QueryType()
		if devType == dxinput.WacomTypeStylus || devType == dxinput.WacomTypeEraser {
			err := v.SetMode(mode)
			if err != nil {
				logger.Warningf("Set mode for '%v - %v' failed: %v",
					v.Id, v.Name, err)
			}
		}
	}
}

func getPressureCurveControlPoints(level int) []int {
	// level 1 ~ 7 Soften ~ Firmer
	const seg = 6.0
	const d = 30.0
	x := (100-2*d)/seg*(float64(level)-1) + d
	y := 100 - x
	x1 := x - d
	y1 := y - d
	x2 := x + d
	y2 := y + d
	return []int{int(x1), int(y1), int(x2), int(y2)}
}

func (w *Wacom) getPressureCurveArray(devType int) ([]int, error) {
	// level is float value
	var level uint32
	switch devType {
	case dxinput.WacomTypeStylus:
		level = w.StylusPressureSensitive.Get()
	case dxinput.WacomTypeEraser:
		level = w.EraserPressureSensitive.Get()
	default:
		return nil, fmt.Errorf("Invalid wacom device type")
	}
	logger.Debug("pressure level:", level)
	if 1 <= level && level <= 7 {
		points := getPressureCurveControlPoints(int(level))
		return points, nil
	} else {
		return nil, fmt.Errorf("Invalid pressure sensitive level %v, range: [1, 7]", level)
	}
}

func (w *Wacom) setPressureSensitiveForType(devType int) {
	for _, v := range w.devInfos {
		if v.QueryType() == devType {
			array, err := w.getPressureCurveArray(devType)
			if err != nil {
				logger.Warning(err)
				continue
			}

			logger.Debug("set curve array:", array)
			err = v.SetPressureCurve(array[0], array[1], array[2], array[3])
			if err != nil {
				logger.Warningf("Set pressure curve for '%v - %v' failed: %v",
					v.Id, v.Name, err)
			}
		}
	}
}

func (w *Wacom) setPressureSensitive() {
	w.setPressureSensitiveForType(dxinput.WacomTypeStylus)
	w.setPressureSensitiveForType(dxinput.WacomTypeEraser)
}

func (w *Wacom) setSuppress() {
	delta := int(w.Suppress.Get())
	for _, dev := range w.devInfos {
		err := dev.SetSuppress(delta)
		if err != nil {
			logger.Debugf("Set suppress for '%v - %v' to %v failed: %v",
				dev.Id, dev.Name, delta, err)
		}
	}
}

func (w *Wacom) mapToOutput() {
	output := strings.Trim(w.MapOutput.Get(), " ")
	if output == "" {
		output = "desktop"
	}

	for _, v := range w.devInfos {
		err := v.MapToOutput(output)
		if err != nil {
			logger.Warningf("Map output for '%v - %v' to %v failed: %v",
				v.Id, v.Name, output, err)
		}
	}
}

func (w *Wacom) getRawSample(devType int) (uint32, error) {
	var rawSample uint32
	switch devType {
	case dxinput.WacomTypeStylus:
		rawSample = w.StylusRawSample.Get()
	case dxinput.WacomTypeEraser:
		rawSample = w.EraserRawSample.Get()
	default:
		return 0, fmt.Errorf("Invalid wacom device type")
	}
	return rawSample, nil
}

func (w *Wacom) setRawSampleForType(devType int) {
	for _, v := range w.devInfos {
		if v.QueryType() == devType {
			rawSample, err := w.getRawSample(devType)
			if err != nil {
				logger.Warning(err)
				continue
			}
			err = v.SetRawSample(rawSample)
			if err != nil {
				logger.Warningf("Set raw sample for '%v - %v' to %v failed: %v",
					v.Id, v.Name, rawSample, err)
			}
		}
	}
}

func (w *Wacom) setRawSample() {
	w.setRawSampleForType(dxinput.WacomTypeStylus)
	w.setRawSampleForType(dxinput.WacomTypeEraser)
}

func (w *Wacom) getThreshold(devType int) (int, error) {
	var threshold int
	switch devType {
	case dxinput.WacomTypeStylus:
		threshold = int(w.StylusThreshold.Get())
	case dxinput.WacomTypeEraser:
		threshold = int(w.EraserThreshold.Get())
	default:
		return 0, fmt.Errorf("Invalid wacom device type")
	}
	return threshold, nil
}

func (w *Wacom) setThresholdForType(devType int) {
	for _, v := range w.devInfos {
		if v.QueryType() == devType {
			threshold, err := w.getThreshold(devType)
			if err != nil {
				logger.Warning(err)
				continue
			}
			err = v.SetThreshold(threshold)
			if err != nil {
				logger.Warningf("Set threshold for '%v - %v' to %v failed: %v",
					v.Id, v.Name, threshold, err)
			}
		}
	}
}

func (w *Wacom) setThreshold() {
	w.setThresholdForType(dxinput.WacomTypeStylus)
	w.setThresholdForType(dxinput.WacomTypeEraser)
}
