/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package shortcuts

import (
	"gir/gio-2.0"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	//"sync"
	"pkg.deepin.io/lib/log"
	"strings"
)

var logger *log.Logger

func SetLogger(l *log.Logger) {
	logger = l
}

type KeyEventFunc func(ev *KeyEvent)

type Shortcuts struct {
	idShortcutMap     map[string]Shortcut
	grabedKeyAccelMap map[Key]*Accel
	eventCb           KeyEventFunc
	xu                *xgbutil.XUtil
}

type KeyEvent struct {
	Mods     Modifiers
	Code     Keycode
	Shortcut Shortcut
}

func NewShortcuts(xu *xgbutil.XUtil, eventCb KeyEventFunc) *Shortcuts {
	ss := &Shortcuts{
		idShortcutMap:     make(map[string]Shortcut),
		grabedKeyAccelMap: make(map[Key]*Accel),
		eventCb:           eventCb,
		xu:                xu,
	}

	return ss
}

func (ss *Shortcuts) List() (list []Shortcut) {
	// RLock ss.idShortcutMap
	for _, shortcut := range ss.idShortcutMap {
		list = append(list, shortcut)
	}
	return
}

func (ss *Shortcuts) grabAccel(shortcut Shortcut, pa ParsedAccel, dummy bool) {
	keys, err := pa.QueryKeys(ss.xu)
	if err != nil {
		return
	}

	// RLock ss.grabedKeyAccelMap
	for _, key := range keys {
		accel, ok := ss.grabedKeyAccelMap[key]
		if ok {
			// conflict
			logger.Debugf("key %v is grabed by %v", key, accel.Shortcut.GetId())
			return
		}
	}

	// dummy grab Super_L and Super_R
	// We use X RECORD extension to receive "Super" key press and release events.
	keyLower := strings.ToLower(pa.Key)
	if keyLower == "super_l" || keyLower == "super_r" {
		dummy = true
	}

	// no conflict
	if !dummy {
		err = keys.Grab(ss.xu)
		if err != nil {
			logger.Debug(err)
			return
		}
	}
	accel := &Accel{
		Parsed:     pa,
		Shortcut:   shortcut,
		GrabedKeys: keys,
	}
	// Lock ss.grabedKeyAccelMap
	// attach key <-> accel
	for _, key := range keys {
		ss.grabedKeyAccelMap[key] = accel
	}
}

func (ss *Shortcuts) ungrabAccel(pa ParsedAccel, dummy bool) {
	keys, err := pa.QueryKeys(ss.xu)
	if err != nil {
		logger.Debug(err)
		return
	}

	// Lock ss.grabedKeyAccelMap
	for _, key := range keys {
		delete(ss.grabedKeyAccelMap, key)
	}

	keys.Ungrab(ss.xu)
}

func (ss *Shortcuts) grabShortcut(shortcut Shortcut) {
	dummy := dummyGrab(shortcut.GetType())

	logger.Debug("grabShortcut shortcut id:", shortcut.GetId())
	for _, pa := range shortcut.GetAccels() {
		logger.Debug("grabAccel accel:", pa)
		ss.grabAccel(shortcut, pa, dummy)
	}
}

func (ss *Shortcuts) ungrabShortcut(shortcut Shortcut) {
	dummy := dummyGrab(shortcut.GetType())

	for _, pa := range shortcut.GetAccels() {
		ss.ungrabAccel(pa, dummy)
	}
}

func (ss *Shortcuts) ModifyShortcutAccels(shortcut Shortcut, newAccels []ParsedAccel) {
	logger.Debug("Shortcuts.ModifyShortcutAccels", shortcut, newAccels)
	ss.ungrabShortcut(shortcut)
	shortcut.setAccels(newAccels)
	ss.grabShortcut(shortcut)
}

func (ss *Shortcuts) AddShortcutAccel(shortcut Shortcut, pa ParsedAccel) {
	logger.Debug("Shortcuts.AddShortcutAccel", shortcut, pa)
	newAccels := shortcut.GetAccels()
	newAccels = append(newAccels, pa)
	shortcut.setAccels(newAccels)

	// grab accel
	dummy := dummyGrab(shortcut.GetType())
	ss.grabAccel(shortcut, pa, dummy)
}

func (ss *Shortcuts) RemoveShortcutAccel(shortcut Shortcut, pa ParsedAccel) {
	logger.Debug("Shortcuts.RemoveShortcutAccel", shortcut, pa)
	logger.Debug("shortcut.GetAccel", shortcut.GetAccels())
	var newAccels []ParsedAccel
	for _, accel := range shortcut.GetAccels() {
		if !accel.Equal(ss.xu, pa) {
			newAccels = append(newAccels, accel)
		}
	}
	shortcut.setAccels(newAccels)
	logger.Debug("shortcut.setAccels", newAccels)

	// ungrab accel
	dummy := dummyGrab(shortcut.GetType())
	ss.ungrabAccel(pa, dummy)
}

func dummyGrab(shortcutType int32) bool {
	if shortcutType == ShortcutTypeWM ||
		shortcutType == ShortcutTypeMetacity {
		return true
	}
	return false
}

func (ss *Shortcuts) UngrabAll() {
	// ungrab all grabed keys
	// Lock ss.grabedKeyAccelMap
	for grabedKey, accel := range ss.grabedKeyAccelMap {
		dummy := dummyGrab(accel.Shortcut.GetType())
		if !dummy {
			grabedKey.Ungrab(ss.xu)
		}
	}
	// new map
	count := len(ss.grabedKeyAccelMap)
	ss.grabedKeyAccelMap = make(map[Key]*Accel, count)
}

func (ss *Shortcuts) GrabAll() {
	// re-grab all shortcuts
	for _, shortcut := range ss.idShortcutMap {
		ss.grabShortcut(shortcut)
	}
}

func (ss *Shortcuts) updateKeymap() {
	// update map before re-bind
	keyMap, modMap := keybind.MapsGet(ss.xu)
	keybind.KeyMapSet(ss.xu, keyMap)
	keybind.ModMapSet(ss.xu, modMap)

	ss.UngrabAll()
	ss.GrabAll()
}

func (ss *Shortcuts) ReloadAllShortcutAccels() []Shortcut {
	var changes []Shortcut
	for _, shortcut := range ss.idShortcutMap {
		changed := shortcut.ReloadAccels()
		if changed {
			changes = append(changes, shortcut)
		}
	}
	return changes
}

// shift, control, alt(mod1), super(mod4)
func GetConcernedModifiers(state uint16) Modifiers {
	var mods Modifiers
	if state&xproto.ModMaskShift > 0 {
		mods |= xproto.ModMaskShift
	}
	if state&xproto.ModMaskControl > 0 {
		mods |= xproto.ModMaskControl
	}
	if state&xproto.ModMask1 > 0 {
		mods |= xproto.ModMask1
	}
	if state&xproto.ModMask4 > 0 {
		mods |= xproto.ModMask4
	}
	return mods
}

func (ss *Shortcuts) handleKeyEvent(xu *xgbutil.XUtil, state uint16, detail xproto.Keycode, pressed bool) {
	mods := GetConcernedModifiers(state)
	code := Keycode(detail)
	logger.Debug("event mods:", Modifiers(state))
	key := Key{
		Mods: mods,
		Code: code,
	}
	logger.Debug("event key:", key)

	if pressed {
		// key press
		ss.emitKeyEvent(Modifiers(state), key)
	}
}

func (ss *Shortcuts) emitKeyEvent(mods Modifiers, key Key) {
	// RLock ss.grabedKeyAccelMap
	accel, ok := ss.grabedKeyAccelMap[key]
	if ok {
		logger.Debugf("accel: %#v", accel)
		keyEvent := &KeyEvent{
			Mods:     mods,
			Code:     key.Code,
			Shortcut: accel.Shortcut,
		}

		ss.eventCb(keyEvent)
	} else {
		logger.Debug("accel not found")
	}
}

// Returns whether the operation was successful
func tryGrabKeyboard(xu *xgbutil.XUtil) bool {
	if err := keybind.GrabKeyboard(xu, xu.RootWin()); err != nil {
		return false
	}
	keybind.UngrabKeyboard(xu)
	return true
}

func (ss *Shortcuts) HandleXRecordKeyRelease(code int) {
	str := keybind.LookupString(ss.xu, 0, xproto.Keycode(code))
	switch strings.ToLower(str) {
	case "super_r", "super_l":
		if ok := tryGrabKeyboard(ss.xu); !ok {
			return
		}

		ss.emitKeyEvent(0, Key{Code: Keycode(code)})
	}
}

func (ss *Shortcuts) ListenXEvents() {
	xevent.KeyPressFun(func(xu *xgbutil.XUtil, ev xevent.KeyPressEvent) {
		logger.Debug(ev)
		ss.handleKeyEvent(xu, ev.State, ev.Detail, true)
	}).Connect(ss.xu, ss.xu.RootWin())

	xevent.KeyReleaseFun(func(xu *xgbutil.XUtil, ev xevent.KeyReleaseEvent) {
		logger.Debug(ev)
		ss.handleKeyEvent(xu, ev.State, ev.Detail, false)
	}).Connect(ss.xu, ss.xu.RootWin())

	xevent.MappingNotifyFun(func(xu *xgbutil.XUtil, ev xevent.MappingNotifyEvent) {
		logger.Debug("MappingNotifyEvent")

		if ev.Request == xproto.MappingKeyboard {
			logger.Debug("Shortcuts.updateKeymap")
			ss.updateKeymap()
		}
	}).Connect(ss.xu, xevent.NoWindow)
}

func (ss *Shortcuts) Add(shortcut Shortcut) {
	ss.AddWithoutLock(shortcut)
}

// TODO private
func (ss *Shortcuts) AddWithoutLock(shortcut Shortcut) {
	logger.Debug("AddWithoutLock", shortcut)
	uid := shortcut.GetUid()
	// Lock ss.idShortcutMap
	ss.idShortcutMap[uid] = shortcut
	ss.grabShortcut(shortcut)
}

func (ss *Shortcuts) Delete(shortcut Shortcut) {
	uid := shortcut.GetUid()
	// Lock ss.idShortcutMap
	delete(ss.idShortcutMap, uid)
	ss.ungrabShortcut(shortcut)
}

func (ss *Shortcuts) GetByIdType(id string, _type int32) Shortcut {
	uid := idType2Uid(id, _type)
	// Lock ss.idShortcutMap
	shortcut := ss.idShortcutMap[uid]
	return shortcut
}

// ret0: Conflicting Accel
// ret1: pa parse key error
func (ss *Shortcuts) FindConflictingAccel(pa ParsedAccel) (*Accel, error) {
	keys, err := pa.QueryKeys(ss.xu)
	if err != nil {
		return nil, err
	}

	logger.Debug("Shortcuts.FindConflictingAccel pa:", pa)
	logger.Debug("keys:", keys)

	// RLock ss.grabedKeyAccelMap
	for _, key := range keys {
		accel, ok := ss.grabedKeyAccelMap[key]
		if ok {
			return accel, nil
		}
	}
	return nil, nil
}

func (ss *Shortcuts) AddSystem(gsettings *gio.Settings) {
	logger.Debug("AddSystem")
	idNameMap := getSystemIdNameMap()
	for _, id := range gsettings.ListKeys() {
		name := idNameMap[id]
		if name == "" {
			name = id
		}
		accels := gsettings.GetStrv(id)
		gs := NewGSettingsShortcut(gsettings, id, ShortcutTypeSystem, accels, name)
		sysShortcut := &SystemShortcut{
			GSettingsShortcut: gs,
			arg: &ActionExecCmdArg{
				Cmd: getSystemAction(id),
			},
		}
		ss.AddWithoutLock(sysShortcut)
	}
}

func (ss *Shortcuts) AddWM(gsettings *gio.Settings) {
	logger.Debug("AddWM")
	idNameMap := getWMIdNameMap()
	for _, id := range gsettings.ListKeys() {
		name := idNameMap[id]
		if name == "" {
			name = id
		}
		accels := gsettings.GetStrv(id)
		gs := NewGSettingsShortcut(gsettings, id, ShortcutTypeWM, accels, name)
		ss.AddWithoutLock(gs)
	}
}

func (ss *Shortcuts) AddMetacity(gsettings *gio.Settings) {
	logger.Debug("AddMetacity")
	idNameMap := getMetacityIdNameMap()
	for _, id := range gsettings.ListKeys() {
		name := idNameMap[id]
		if name == "" {
			name = id
		}
		accels := gsettings.GetStrv(id)
		gs := NewGSettingsShortcut(gsettings, id, ShortcutTypeMetacity, accels, name)
		ss.AddWithoutLock(gs)
	}
}

func (ss *Shortcuts) AddMedia(gsettings *gio.Settings) {
	logger.Debug("AddMedia")
	idNameMap := getMediaIdNameMap()
	for _, id := range gsettings.ListKeys() {
		name := idNameMap[id]
		if name == "" {
			name = id
		}
		accels := gsettings.GetStrv(id)
		gs := NewGSettingsShortcut(gsettings, id, ShortcutTypeMedia, accels, name)
		mediaShortcut := &MediaShortcut{
			GSettingsShortcut: gs,
		}
		ss.AddWithoutLock(mediaShortcut)
	}
}

func (ss *Shortcuts) AddCustom(csm *CustomShortcutManager) {
	logger.Debug("AddCustom")
	for _, shortcut := range csm.List() {
		ss.AddWithoutLock(shortcut)
	}
}
