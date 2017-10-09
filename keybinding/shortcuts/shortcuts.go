/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package shortcuts

import (
	"strings"
	"sync"

	"gir/gio-2.0"
	"pkg.deepin.io/lib/log"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/record"
	"github.com/linuxdeepin/go-x11-client/util/keybind"
	"github.com/linuxdeepin/go-x11-client/util/keysyms"
	"github.com/linuxdeepin/go-x11-client/util/wm/ewmh"
)

var logger *log.Logger

const (
	SKLCtrlShift uint32 = 1 << iota
	SKLAltShift
	SKLSuperSpace
)

func SetLogger(l *log.Logger) {
	logger = l
}

type KeyEventFunc func(ev *KeyEvent)

type ShortcutManager struct {
	conn     *x.Conn
	dataConn *x.Conn // conn for receive record event

	idShortcutMap     map[string]Shortcut
	grabedKeyAccelMap map[Key]*Accel
	keySymbols        *keysyms.KeySymbols
	ewmhConn          *ewmh.Conn

	recordEnable        bool
	recordContext       record.Context
	xRecordEventHandler *XRecordEventHandler
	eventCb             KeyEventFunc
	eventCbMu           sync.Mutex
}

type KeyEvent struct {
	Mods     Modifiers
	Code     Keycode
	Shortcut Shortcut
}

func NewShortcutManager(conn *x.Conn, keySymbols *keysyms.KeySymbols, eventCb KeyEventFunc) *ShortcutManager {
	ss := &ShortcutManager{
		idShortcutMap:     make(map[string]Shortcut),
		grabedKeyAccelMap: make(map[Key]*Accel),
		eventCb:           eventCb,
		conn:              conn,
		keySymbols:        keySymbols,
		recordEnable:      true,
	}

	ss.ewmhConn, _ = ewmh.NewConn(conn)

	ss.xRecordEventHandler = NewXRecordEventHandler(keySymbols)
	ss.xRecordEventHandler.modKeyReleasedCb = func(code uint8, mods uint16) {
		if isKbdAlreadyGrabbed(ss.conn, ss.ewmhConn) {
			return
		}
		switch mods {
		case keysyms.ModMaskCapsLock, keysyms.ModMaskNumLock, keysyms.ModMaskSuper:
			// caps_lock, num_lock, super
			ss.emitKeyEvent(0, Key{Code: Keycode(code)})

		case keysyms.ModMaskControl | keysyms.ModMaskShift:
			// ctrl-shift
			ss.emitFakeKeyEvent(&Action{Type: ActionTypeSwitchKbdLayout, Arg: SKLCtrlShift})

		case keysyms.ModMaskAlt | keysyms.ModMaskShift:
			// alt-shift
			ss.emitFakeKeyEvent(&Action{Type: ActionTypeSwitchKbdLayout, Arg: SKLAltShift})
		}
	}
	// init record
	err := ss.initRecord()
	if err == nil {
		logger.Debug("start record event loop")
		go ss.recordEventLoop()
	} else {
		logger.Warning("init record failed: ", err)
	}

	return ss
}

func (ss *ShortcutManager) recordEventLoop() {
	// enable context
	cookie := record.EnableContext(ss.dataConn, ss.recordContext)

	for {
		reply, err := cookie.Reply(ss.dataConn)
		if err != nil {
			logger.Warning(err)
			return
		}
		if !ss.recordEnable {
			logger.Warning("record disabled!")
			continue
		}

		if reply.ClientSwapped != 0 {
			logger.Warning("reply.ClientSwapped is true")
			continue
		}
		if len(reply.Data) == 0 {
			continue
		}

		ge := x.GenericEvent(reply.Data)

		switch ge.GetEventCode() {
		case x.KeyPressEventCode:
			event, _ := x.NewKeyPressEvent(ge)
			//logger.Debug(event)
			ss.handleXRecordKeyEvent(true, uint8(event.Detail), event.State)

		case x.KeyReleaseEventCode:
			event, _ := x.NewKeyReleaseEvent(ge)
			//logger.Debug(event)
			ss.handleXRecordKeyEvent(false, uint8(event.Detail), event.State)

		case x.ButtonPressEventCode:
			//event, _ := x.NewButtonPressEvent(ge)
			//logger.Debug(event)
			ss.handleXRecordButtonEvent(true)
		case x.ButtonReleaseEventCode:
			//event, _ := x.NewButtonReleaseEvent(ge)
			//logger.Debug(event)
			ss.handleXRecordButtonEvent(false)
		default:
			logger.Debug(ge)
		}

	}
}

func (ss *ShortcutManager) initRecord() error {
	ctrlConn := ss.conn
	dataConn, err := x.NewConn()
	if err != nil {
		return err
	}

	_, err = record.QueryVersion(ctrlConn, record.MajorVersion, record.MinorVersion).Reply(ctrlConn)
	if err != nil {
		return err
	}
	_, err = record.QueryVersion(dataConn, record.MajorVersion, record.MinorVersion).Reply(dataConn)
	if err != nil {
		return err
	}

	xid, err := ctrlConn.GenerateID()
	if err != nil {
		return err
	}
	ctx := record.Context(xid)
	logger.Debug("record context id:", ctx)

	// create context
	clientSpec := []record.ClientSpec{record.ClientSpec(record.CSAllClients)}
	ranges := []record.Range{
		{
			DeviceEvents: record.Range8{
				First: x.KeyPressEventCode,
				Last:  x.ButtonReleaseEventCode,
			},
		},
	}

	err = record.CreateContextChecked(ctrlConn, ctx, record.ElementHeader(0), 1, 1, clientSpec, ranges).Check(ctrlConn)
	if err != nil {
		return err
	}

	ss.recordContext = ctx
	ss.dataConn = dataConn
	return nil
}

func (ss *ShortcutManager) EnableRecord(val bool) {
	ss.recordEnable = val
}

func (ss *ShortcutManager) Destroy() {
	// TODO
}

func (ss *ShortcutManager) List() (list []Shortcut) {
	// RLock ss.idShortcutMap
	for _, shortcut := range ss.idShortcutMap {
		list = append(list, shortcut)
	}
	return
}

func (ss *ShortcutManager) grabAccel(shortcut Shortcut, pa ParsedAccel, dummy bool) {
	key, err := pa.QueryKey(ss.keySymbols)
	if err != nil {
		logger.Debugf("getAccel failed shortcut: %v, pa: %v, err: %v", shortcut.GetId(), pa, err)
		return
	}
	//logger.Debugf("grabAccel shortcut: %s, pa: %s, key: %s, dummy: %v", shortcut.GetId(), pa, key, dummy)

	// RLock ss.grabedKeyAccelMap
	if confAccel, ok := ss.grabedKeyAccelMap[key]; ok {
		// conflict
		logger.Debugf("key %v is grabed by %v", key, confAccel.Shortcut.GetId())
		return
	}

	// no conflict
	if !dummy {
		err = key.Grab(ss.conn)
		if err != nil {
			logger.Debug(err)
			return
		}
	}
	accel := &Accel{
		Parsed:    pa,
		Shortcut:  shortcut,
		GrabedKey: key,
	}
	// Lock ss.grabedKeyAccelMap
	// attach key <-> accel
	ss.grabedKeyAccelMap[key] = accel
}

func (ss *ShortcutManager) ungrabAccel(pa ParsedAccel, dummy bool) {
	key, err := pa.QueryKey(ss.keySymbols)
	if err != nil {
		logger.Debug(err)
		return
	}

	// Lock ss.grabedKeyAccelMap
	delete(ss.grabedKeyAccelMap, key)
	key.Ungrab(ss.conn)
}

func (ss *ShortcutManager) grabShortcut(shortcut Shortcut) {
	//logger.Debug("grabShortcut shortcut id:", shortcut.GetId())
	for _, pa := range shortcut.GetAccels() {
		dummy := dummyGrab(shortcut, pa)
		ss.grabAccel(shortcut, pa, dummy)
	}
}

func (ss *ShortcutManager) ungrabShortcut(shortcut Shortcut) {

	for _, pa := range shortcut.GetAccels() {
		dummy := dummyGrab(shortcut, pa)
		ss.ungrabAccel(pa, dummy)
	}
}

func (ss *ShortcutManager) ModifyShortcutAccels(shortcut Shortcut, newAccels []ParsedAccel) {
	logger.Debug("ShortcutManager.ModifyShortcutAccels", shortcut, newAccels)
	ss.ungrabShortcut(shortcut)
	shortcut.setAccels(newAccels)
	ss.grabShortcut(shortcut)
}

func (ss *ShortcutManager) AddShortcutAccel(shortcut Shortcut, pa ParsedAccel) {
	logger.Debug("ShortcutManager.AddShortcutAccel", shortcut, pa)
	newAccels := shortcut.GetAccels()
	newAccels = append(newAccels, pa)
	shortcut.setAccels(newAccels)

	// grab accel
	dummy := dummyGrab(shortcut, pa)
	ss.grabAccel(shortcut, pa, dummy)
}

func (ss *ShortcutManager) RemoveShortcutAccel(shortcut Shortcut, pa ParsedAccel) {
	logger.Debug("ShortcutManager.RemoveShortcutAccel", shortcut, pa)
	logger.Debug("shortcut.GetAccel", shortcut.GetAccels())
	var newAccels []ParsedAccel
	for _, accel := range shortcut.GetAccels() {
		if !accel.Equal(ss.keySymbols, pa) {
			newAccels = append(newAccels, accel)
		}
	}
	shortcut.setAccels(newAccels)
	logger.Debug("shortcut.setAccels", newAccels)

	// ungrab accel
	dummy := dummyGrab(shortcut, pa)
	ss.ungrabAccel(pa, dummy)
}

func dummyGrab(shortcut Shortcut, pa ParsedAccel) bool {
	if shortcut.GetType() == ShortcutTypeWM {
		return true
	}

	switch strings.ToLower(pa.Key) {
	case "super_l", "super_r", "caps_lock", "num_lock":
		return true
	}
	return false
}

func (ss *ShortcutManager) UngrabAll() {
	// ungrab all grabed keys
	// Lock ss.grabedKeyAccelMap
	for grabedKey, accel := range ss.grabedKeyAccelMap {
		dummy := dummyGrab(accel.Shortcut, accel.Parsed)
		if !dummy {
			grabedKey.Ungrab(ss.conn)
		}
	}
	// new map
	count := len(ss.grabedKeyAccelMap)
	ss.grabedKeyAccelMap = make(map[Key]*Accel, count)
}

func (ss *ShortcutManager) GrabAll() {
	// re-grab all shortcuts
	for _, shortcut := range ss.idShortcutMap {
		ss.grabShortcut(shortcut)
	}
}

func (ss *ShortcutManager) regrabAll() {
	logger.Debug("regrabAll")
	ss.UngrabAll()
	ss.GrabAll()
}

func (ss *ShortcutManager) ReloadAllShortcutAccels() []Shortcut {
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
func getConcernedMods(state uint16) uint16 {
	var mods uint16
	if state&keysyms.ModMaskShift > 0 {
		mods |= keysyms.ModMaskShift
	}
	if state&keysyms.ModMaskControl > 0 {
		mods |= keysyms.ModMaskControl
	}
	if state&keysyms.ModMaskAlt > 0 {
		mods |= keysyms.ModMaskAlt
	}
	if state&keysyms.ModMaskSuper > 0 {
		mods |= keysyms.ModMaskSuper
	}
	return mods
}

func GetConcernedModifiers(state uint16) Modifiers {
	return Modifiers(getConcernedMods(state))
}

func combineStateCode2Key(state uint16, code uint8) Key {
	mods := GetConcernedModifiers(state)
	_code := Keycode(code)
	key := Key{
		Mods: mods,
		Code: _code,
	}
	return key
}

func (ss *ShortcutManager) callEventCallback(ev *KeyEvent) {
	ss.eventCbMu.Lock()
	ss.eventCb(ev)
	ss.eventCbMu.Unlock()
}

func (ss *ShortcutManager) handleKeyEvent(pressed bool, detail x.Keycode, state uint16) {
	key := combineStateCode2Key(state, uint8(detail))
	logger.Debug("event key:", key)

	if pressed {
		// key press
		ss.emitKeyEvent(Modifiers(state), key)
	}
}

func (ss *ShortcutManager) emitFakeKeyEvent(action *Action) {
	keyEvent := &KeyEvent{
		Shortcut: NewFakeShortcut(action),
	}
	ss.callEventCallback(keyEvent)
}

func (ss *ShortcutManager) emitKeyEvent(mods Modifiers, key Key) {
	// RLock ss.grabedKeyAccelMap
	accel, ok := ss.grabedKeyAccelMap[key]
	if ok {
		logger.Debugf("accel: %#v", accel)
		keyEvent := &KeyEvent{
			Mods:     mods,
			Code:     key.Code,
			Shortcut: accel.Shortcut,
		}

		ss.callEventCallback(keyEvent)
	} else {
		logger.Debug("accel not found")
	}
}

func isKbdAlreadyGrabbed(conn *x.Conn, ewmhConn *ewmh.Conn) bool {
	var grabWin x.Window

	rootWin := conn.GetDefaultScreen().Root
	if activeWin, _ := ewmhConn.GetActiveWindow().Reply(ewmhConn); activeWin == 0 {
		grabWin = rootWin
	} else {
		// check viewable
		attrs, err := x.GetWindowAttributes(conn, activeWin).Reply(conn)
		if err != nil {
			grabWin = rootWin
		} else if attrs.MapState != x.MapStateViewable {
			// err is nil and activeWin is not viewable
			grabWin = rootWin
		} else {
			// err is nil, activeWin is viewable
			grabWin = activeWin
		}
	}

	err := keybind.GrabKeyboard(conn, grabWin)
	if err == nil {
		// grab keyboard successful
		keybind.UngrabKeyboard(conn)
		return false
	}

	logger.Warningf("GrabKeyboard win %d failed: %v", grabWin, err)

	gkErr, ok := err.(keybind.GrabKeyboardError)
	if ok && gkErr.Status == x.GrabStatusAlreadyGrabbed {
		return true
	}
	return false
}

func (ss *ShortcutManager) SetAllModKeysReleasedCallback(cb func()) {
	ss.xRecordEventHandler.allModKeysReleasedCb = cb
}

func (ss *ShortcutManager) handleXRecordKeyEvent(pressed bool, code uint8, state uint16) {
	ss.xRecordEventHandler.handleKeyEvent(pressed, code, state)
	if pressed {
		// Special handling screenshot* shortcuts
		key := combineStateCode2Key(state, code)
		accel, ok := ss.grabedKeyAccelMap[key]
		if ok {
			shortcut := accel.Shortcut
			if shortcut != nil && shortcut.GetType() == ShortcutTypeSystem &&
				strings.HasPrefix(shortcut.GetId(), "screenshot") {
				keyEvent := &KeyEvent{
					Mods:     key.Mods,
					Code:     key.Code,
					Shortcut: shortcut,
				}
				logger.Debug("handleXRecordKeyEvent: emit key event for screenshot* shortcuts")
				ss.callEventCallback(keyEvent)
			}
		}
	}
}

func (ss *ShortcutManager) handleXRecordButtonEvent(pressed bool) {
	ss.xRecordEventHandler.handleButtonEvent(pressed)
}

func (ss *ShortcutManager) EventLoop() {
	for {
		ev := ss.conn.WaitForEvent()
		switch ev.GetEventCode() {
		case x.KeyPressEventCode:
			event, _ := x.NewKeyPressEvent(ev)
			logger.Debug(event)
			ss.handleKeyEvent(true, event.Detail, event.State)
		case x.KeyReleaseEventCode:
			event, _ := x.NewKeyReleaseEvent(ev)
			logger.Debug(event)
			ss.handleKeyEvent(false, event.Detail, event.State)
		case x.MappingNotifyEventCode:
			event, _ := x.NewMappingNotifyEvent(ev)
			logger.Debug(event)
			if ss.keySymbols.RefreshKeyboardMapping(event) {
				ss.regrabAll()
			}
		}
	}
}

func (ss *ShortcutManager) Add(shortcut Shortcut) {
	ss.AddWithoutLock(shortcut)
}

// TODO private
func (ss *ShortcutManager) AddWithoutLock(shortcut Shortcut) {
	logger.Debug("AddWithoutLock", shortcut)
	uid := shortcut.GetUid()
	// Lock ss.idShortcutMap
	ss.idShortcutMap[uid] = shortcut
	ss.grabShortcut(shortcut)
}

func (ss *ShortcutManager) Delete(shortcut Shortcut) {
	uid := shortcut.GetUid()
	// Lock ss.idShortcutMap
	delete(ss.idShortcutMap, uid)
	ss.ungrabShortcut(shortcut)
}

func (ss *ShortcutManager) GetByIdType(id string, _type int32) Shortcut {
	uid := idType2Uid(id, _type)
	// Lock ss.idShortcutMap
	shortcut := ss.idShortcutMap[uid]
	return shortcut
}

// ret0: Conflicting Accel
// ret1: pa parse key error
func (ss *ShortcutManager) FindConflictingAccel(pa ParsedAccel) (*Accel, error) {
	key, err := pa.QueryKey(ss.keySymbols)
	if err != nil {
		return nil, err
	}

	logger.Debug("ShortcutManager.FindConflictingAccel pa:", pa)
	logger.Debug("key:", key)

	// RLock ss.grabedKeyAccelMap
	accel, ok := ss.grabedKeyAccelMap[key]
	if ok {
		return accel, nil
	}
	return nil, nil
}

func (ss *ShortcutManager) AddSystem(gsettings *gio.Settings) {
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

func (ss *ShortcutManager) AddWM(gsettings *gio.Settings) {
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

func (ss *ShortcutManager) AddMedia(gsettings *gio.Settings) {
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

func (ss *ShortcutManager) AddCustom(csm *CustomShortcutManager) {
	logger.Debug("AddCustom")
	for _, shortcut := range csm.List() {
		ss.AddWithoutLock(shortcut)
	}
}

func (ss *ShortcutManager) AddSpecial() {
	idNameMap := getSpecialIdNameMap()

	// add SwitchKbdLayout <Super>space
	s0 := NewFakeShortcut(&Action{Type: ActionTypeSwitchKbdLayout, Arg: SKLSuperSpace})
	pa, err := ParseStandardAccel("<Super>space")
	if err != nil {
		panic(err)
	}
	s0.Id = "switch-kbd-layout"
	s0.Name = idNameMap[s0.Id]
	s0.Accels = []ParsedAccel{pa}
	ss.AddWithoutLock(s0)
}
