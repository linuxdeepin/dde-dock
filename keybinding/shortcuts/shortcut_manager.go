/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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
	"time"

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
	keyKeystrokeMap   map[Key]*Keystroke
	keyKeystrokeMapMu sync.Mutex
	keySymbols        *keysyms.KeySymbols
	ewmhConn          *ewmh.Conn

	recordEnable        bool
	recordContext       record.Context
	xRecordEventHandler *XRecordEventHandler
	eventCb             KeyEventFunc
	eventCbMu           sync.Mutex
	layoutChanged       chan struct{}

	ConflictingKeystrokes []*Keystroke
	EliminateConflictDone bool
}

type KeyEvent struct {
	Mods     Modifiers
	Code     Keycode
	Shortcut Shortcut
}

func NewShortcutManager(conn *x.Conn, keySymbols *keysyms.KeySymbols, eventCb KeyEventFunc) *ShortcutManager {
	ss := &ShortcutManager{
		idShortcutMap:   make(map[string]Shortcut),
		eventCb:         eventCb,
		conn:            conn,
		keySymbols:      keySymbols,
		recordEnable:    true,
		keyKeystrokeMap: make(map[Key]*Keystroke),
		layoutChanged:   make(chan struct{}),
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

func (sm *ShortcutManager) recordEventLoop() {
	// enable context
	cookie := record.EnableContext(sm.dataConn, sm.recordContext)

	for {
		reply, err := cookie.Reply(sm.dataConn)
		if err != nil {
			logger.Warning(err)
			return
		}
		if !sm.recordEnable {
			logger.Debug("record disabled!")
			continue
		}

		if reply.ClientSwapped {
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
			sm.handleXRecordKeyEvent(true, uint8(event.Detail), event.State)

		case x.KeyReleaseEventCode:
			event, _ := x.NewKeyReleaseEvent(ge)
			//logger.Debug(event)
			sm.handleXRecordKeyEvent(false, uint8(event.Detail), event.State)

		case x.ButtonPressEventCode:
			//event, _ := x.NewButtonPressEvent(ge)
			//logger.Debug(event)
			sm.handleXRecordButtonEvent(true)
		case x.ButtonReleaseEventCode:
			//event, _ := x.NewButtonReleaseEvent(ge)
			//logger.Debug(event)
			sm.handleXRecordButtonEvent(false)
		default:
			logger.Debug(ge)
		}

	}
}

func (sm *ShortcutManager) initRecord() error {
	ctrlConn := sm.conn
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

	err = record.CreateContextChecked(ctrlConn, ctx, record.ElementHeader(0),
		clientSpec, ranges).Check(ctrlConn)
	if err != nil {
		return err
	}

	sm.recordContext = ctx
	sm.dataConn = dataConn
	return nil
}

func (sm *ShortcutManager) EnableRecord(val bool) {
	sm.recordEnable = val
}

func (sm *ShortcutManager) NotifyLayoutChanged() {
	sm.layoutChanged <- struct{}{}
}

func (sm *ShortcutManager) Destroy() {
	// TODO
}

func (sm *ShortcutManager) List() (list []Shortcut) {
	// RLock ss.idShortcutMap
	for _, shortcut := range sm.idShortcutMap {
		list = append(list, shortcut)
	}
	return
}

func (sm *ShortcutManager) ListByType(type0 int32) (list []Shortcut) {
	for _, shortcut := range sm.idShortcutMap {
		if type0 == shortcut.GetType() {
			list = append(list, shortcut)
		}
	}
	return
}

func (sm *ShortcutManager) storeConflictingKeystroke(ks *Keystroke) {
	sm.ConflictingKeystrokes = append(sm.ConflictingKeystrokes, ks)
}

func (sm *ShortcutManager) grabKeystroke(shortcut Shortcut, ks *Keystroke, dummy bool) {
	key, err := ks.ToKey(sm.keySymbols)
	if err != nil {
		logger.Debugf("grabKeystroke failed, shortcut: %v, ks: %v, err: %v", shortcut.GetId(), ks, err)
		return
	}
	//logger.Debugf("grabKeystroke shortcut: %s, ks: %s, key: %s, dummy: %v", shortcut.GetId(), ks, key, dummy)

	sm.keyKeystrokeMapMu.Lock()
	conflictKeystroke, ok := sm.keyKeystrokeMap[key]
	sm.keyKeystrokeMapMu.Unlock()

	if ok {
		// conflict
		logger.Debugf("key %v is grabed by %v", key, conflictKeystroke.Shortcut.GetId())
		if !sm.EliminateConflictDone {
			sm.storeConflictingKeystroke(ks)
		}
		return
	}

	// no conflict
	if !dummy {
		err = key.Grab(sm.conn)
		if err != nil {
			logger.Debug(err)
			return
		}
	}
	sm.keyKeystrokeMapMu.Lock()
	sm.keyKeystrokeMap[key] = ks
	sm.keyKeystrokeMapMu.Unlock()
}

func (sm *ShortcutManager) ungrabKeystroke(ks *Keystroke, dummy bool) {
	key, err := ks.ToKey(sm.keySymbols)
	if err != nil {
		logger.Debug(err)
		return
	}

	sm.keyKeystrokeMapMu.Lock()
	delete(sm.keyKeystrokeMap, key)
	sm.keyKeystrokeMapMu.Unlock()
	if !dummy {
		key.Ungrab(sm.conn)
	}
}

func (sm *ShortcutManager) grabShortcut(shortcut Shortcut) {
	//logger.Debug("grabShortcut shortcut id:", shortcut.GetId())
	for _, ks := range shortcut.GetKeystrokes() {
		dummy := dummyGrab(shortcut, ks)
		sm.grabKeystroke(shortcut, ks, dummy)
		ks.Shortcut = shortcut
	}
}

func (sm *ShortcutManager) ungrabShortcut(shortcut Shortcut) {

	for _, ks := range shortcut.GetKeystrokes() {
		dummy := dummyGrab(shortcut, ks)
		sm.ungrabKeystroke(ks, dummy)
		ks.Shortcut = nil
	}
}

func (sm *ShortcutManager) ModifyShortcutKeystrokes(shortcut Shortcut, newVal []*Keystroke) {
	logger.Debug("ShortcutManager.ModifyShortcutKeystrokes", shortcut, newVal)
	sm.ungrabShortcut(shortcut)
	shortcut.setKeystrokes(newVal)
	sm.grabShortcut(shortcut)
}

func (sm *ShortcutManager) AddShortcutKeystroke(shortcut Shortcut, ks *Keystroke) {
	logger.Debug("ShortcutManager.AddShortcutKeystroke", shortcut, ks.DebugString())
	oldVal := shortcut.GetKeystrokes()
	notExist := true
	for _, ks0 := range oldVal {
		if ks.Equal(sm.keySymbols, ks0) {
			notExist = false
			break
		}
	}
	if notExist {
		shortcut.setKeystrokes(append(oldVal, ks))
		logger.Debug("shortcut.Keystrokes append", ks.DebugString())

		// grab keystroke
		dummy := dummyGrab(shortcut, ks)
		sm.grabKeystroke(shortcut, ks, dummy)
	}
	ks.Shortcut = shortcut
}

func (sm *ShortcutManager) DeleteShortcutKeystroke(shortcut Shortcut, ks *Keystroke) {
	logger.Debug("ShortcutManager.DeleteShortcutKeystroke", shortcut, ks.DebugString())
	oldVal := shortcut.GetKeystrokes()
	var newVal []*Keystroke
	for _, ks0 := range oldVal {
		// Leaving unequal values
		if !ks.Equal(sm.keySymbols, ks0) {
			newVal = append(newVal, ks0)
		}
	}
	shortcut.setKeystrokes(newVal)
	logger.Debugf("shortcut.Keystrokes  %v -> %v", oldVal, newVal)

	// ungrab keystroke
	dummy := dummyGrab(shortcut, ks)
	sm.ungrabKeystroke(ks, dummy)
	ks.Shortcut = nil
}

func dummyGrab(shortcut Shortcut, ks *Keystroke) bool {
	if shortcut.GetType() == ShortcutTypeWM {
		return true
	}

	switch strings.ToLower(ks.Keystr) {
	case "super_l", "super_r", "caps_lock", "num_lock":
		return true
	}
	return false
}

func (sm *ShortcutManager) UngrabAll() {
	sm.keyKeystrokeMapMu.Lock()
	// ungrab all grabed keys
	for key, keystroke := range sm.keyKeystrokeMap {
		dummy := dummyGrab(keystroke.Shortcut, keystroke)
		if !dummy {
			key.Ungrab(sm.conn)
		}
	}
	// new map
	count := len(sm.keyKeystrokeMap)
	sm.keyKeystrokeMap = make(map[Key]*Keystroke, count)
	sm.keyKeystrokeMapMu.Unlock()
}

func (sm *ShortcutManager) GrabAll() {
	// re-grab all shortcuts
	for _, shortcut := range sm.idShortcutMap {
		sm.grabShortcut(shortcut)
	}
}

func (sm *ShortcutManager) regrabAll() {
	logger.Debug("regrabAll")
	sm.UngrabAll()
	sm.GrabAll()
}

func (sm *ShortcutManager) ReloadAllShortcutsKeystrokes() []Shortcut {
	var changes []Shortcut
	for _, shortcut := range sm.idShortcutMap {
		changed := shortcut.ReloadKeystrokes()
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
	key := Key{
		Mods: mods,
		Code: Keycode(code),
	}
	return key
}

func (sm *ShortcutManager) callEventCallback(ev *KeyEvent) {
	sm.eventCbMu.Lock()
	sm.eventCb(ev)
	sm.eventCbMu.Unlock()
}

func (sm *ShortcutManager) handleKeyEvent(pressed bool, detail x.Keycode, state uint16) {
	key := combineStateCode2Key(state, uint8(detail))
	logger.Debug("event key:", key)

	if pressed {
		// key press
		sm.emitKeyEvent(Modifiers(state), key)
	}
}

func (sm *ShortcutManager) emitFakeKeyEvent(action *Action) {
	keyEvent := &KeyEvent{
		Shortcut: NewFakeShortcut(action),
	}
	sm.callEventCallback(keyEvent)
}

func (sm *ShortcutManager) emitKeyEvent(mods Modifiers, key Key) {
	sm.keyKeystrokeMapMu.Lock()
	keystroke, ok := sm.keyKeystrokeMap[key]
	sm.keyKeystrokeMapMu.Unlock()
	if ok {
		logger.Debugf("emitKeyEvent keystroke: %#v", keystroke)
		keyEvent := &KeyEvent{
			Mods:     mods,
			Code:     key.Code,
			Shortcut: keystroke.Shortcut,
		}

		sm.callEventCallback(keyEvent)
	} else {
		logger.Debug("keystroke not found")
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

func (sm *ShortcutManager) SetAllModKeysReleasedCallback(cb func()) {
	sm.xRecordEventHandler.allModKeysReleasedCb = cb
}

func (sm *ShortcutManager) handleXRecordKeyEvent(pressed bool, code uint8, state uint16) {
	sm.xRecordEventHandler.handleKeyEvent(pressed, code, state)
	if pressed {
		// Special handling screenshot* shortcuts
		key := combineStateCode2Key(state, code)
		sm.keyKeystrokeMapMu.Lock()
		keystroke, ok := sm.keyKeystrokeMap[key]
		sm.keyKeystrokeMapMu.Unlock()
		if ok {
			shortcut := keystroke.Shortcut
			if shortcut != nil && shortcut.GetType() == ShortcutTypeSystem &&
				strings.HasPrefix(shortcut.GetId(), "screenshot") {
				keyEvent := &KeyEvent{
					Mods:     key.Mods,
					Code:     key.Code,
					Shortcut: shortcut,
				}
				logger.Debug("handleXRecordKeyEvent: emit key event for screenshot* shortcuts")
				sm.callEventCallback(keyEvent)
			}
		}
	}
}

func (sm *ShortcutManager) handleXRecordButtonEvent(pressed bool) {
	sm.xRecordEventHandler.handleButtonEvent(pressed)
}

func (sm *ShortcutManager) EventLoop() {
	eventChan := make(chan x.GenericEvent, 500)
	sm.conn.AddEventChan(eventChan)
	for ev := range eventChan {
		switch ev.GetEventCode() {
		case x.KeyPressEventCode:
			event, _ := x.NewKeyPressEvent(ev)
			logger.Debug(event)
			sm.handleKeyEvent(true, event.Detail, event.State)
		case x.KeyReleaseEventCode:
			event, _ := x.NewKeyReleaseEvent(ev)
			logger.Debug(event)
			sm.handleKeyEvent(false, event.Detail, event.State)
		case x.MappingNotifyEventCode:
			event, _ := x.NewMappingNotifyEvent(ev)
			logger.Debug(event)
			if sm.keySymbols.RefreshKeyboardMapping(event) {
				go func() {
					select {
					case _, ok := <-sm.layoutChanged:
						if !ok {
							logger.Error("Invalid layout changed event")
							return
						}

						sm.regrabAll()
					case _, ok := <-time.After(3 * time.Second):
						if !ok {
							logger.Error("Invalid time event")
							return
						}

						logger.Debug("layout not changed")
					}
				}()
			}
		}
	}
}

func (sm *ShortcutManager) Add(shortcut Shortcut) {
	sm.AddWithoutLock(shortcut)
}

// TODO private
func (sm *ShortcutManager) AddWithoutLock(shortcut Shortcut) {
	logger.Debug("AddWithoutLock", shortcut)
	uid := shortcut.GetUid()
	// Lock ss.idShortcutMap
	sm.idShortcutMap[uid] = shortcut
	sm.grabShortcut(shortcut)
}

func (sm *ShortcutManager) Delete(shortcut Shortcut) {
	uid := shortcut.GetUid()
	// Lock ss.idShortcutMap
	delete(sm.idShortcutMap, uid)
	sm.ungrabShortcut(shortcut)
}

func (sm *ShortcutManager) GetByIdType(id string, type0 int32) Shortcut {
	uid := idType2Uid(id, type0)
	// Lock ss.idShortcutMap
	shortcut := sm.idShortcutMap[uid]
	return shortcut
}

// ret0: Conflicting keystroke
// ret1: error
func (sm *ShortcutManager) FindConflictingKeystroke(ks *Keystroke) (*Keystroke, error) {
	key, err := ks.ToKey(sm.keySymbols)
	if err != nil {
		return nil, err
	}

	logger.Debug("ShortcutManager.FindConflictingKeystroke", ks.DebugString())
	logger.Debug("key:", key)

	sm.keyKeystrokeMapMu.Lock()
	ks1, ok := sm.keyKeystrokeMap[key]
	sm.keyKeystrokeMapMu.Unlock()
	if ok {
		return ks1, nil
	}
	return nil, nil
}

func (sm *ShortcutManager) AddSystem(gsettings *gio.Settings) {
	logger.Debug("AddSystem")
	idNameMap := getSystemIdNameMap()
	for _, id := range gsettings.ListKeys() {
		name := idNameMap[id]
		if name == "" {
			name = id
		}
		keystrokes := gsettings.GetStrv(id)
		gs := NewGSettingsShortcut(gsettings, id, ShortcutTypeSystem, keystrokes, name)
		sysShortcut := &SystemShortcut{
			GSettingsShortcut: gs,
			arg: &ActionExecCmdArg{
				Cmd: getSystemActionCmd(id),
			},
		}
		sm.AddWithoutLock(sysShortcut)
	}
}

func (sm *ShortcutManager) AddWM(gsettings *gio.Settings) {
	logger.Debug("AddWM")
	idNameMap := getWMIdNameMap()
	for _, id := range gsettings.ListKeys() {
		name := idNameMap[id]
		if name == "" {
			name = id
		}
		keystrokes := gsettings.GetStrv(id)
		gs := NewGSettingsShortcut(gsettings, id, ShortcutTypeWM, keystrokes, name)
		sm.AddWithoutLock(gs)
	}
}

func (sm *ShortcutManager) AddMedia(gsettings *gio.Settings) {
	logger.Debug("AddMedia")
	idNameMap := getMediaIdNameMap()
	for _, id := range gsettings.ListKeys() {
		name := idNameMap[id]
		if name == "" {
			name = id
		}
		keystrokes := gsettings.GetStrv(id)
		gs := NewGSettingsShortcut(gsettings, id, ShortcutTypeMedia, keystrokes, name)
		mediaShortcut := &MediaShortcut{
			GSettingsShortcut: gs,
		}
		sm.AddWithoutLock(mediaShortcut)
	}
}

func (sm *ShortcutManager) AddCustom(csm *CustomShortcutManager) {
	logger.Debug("AddCustom")
	for _, shortcut := range csm.List() {
		sm.AddWithoutLock(shortcut)
	}
}

func (sm *ShortcutManager) AddSpecial() {
	idNameMap := getSpecialIdNameMap()

	// add SwitchKbdLayout <Super>space
	s0 := NewFakeShortcut(&Action{Type: ActionTypeSwitchKbdLayout, Arg: SKLSuperSpace})
	ks, err := ParseKeystroke("<Super>space")
	if err != nil {
		panic(err)
	}
	s0.Id = "switch-kbd-layout"
	s0.Name = idNameMap[s0.Id]
	s0.Keystrokes = []*Keystroke{ks}
	sm.AddWithoutLock(s0)
}
