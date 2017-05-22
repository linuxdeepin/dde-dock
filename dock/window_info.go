/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const windowHashPrefix = "w:"

type WindowInfo struct {
	innerId string
	window  xproto.Window
	Title   string
	Icon    string

	x                        int16
	y                        int16
	width                    uint16
	height                   uint16
	lastConfigureNotifyEvent *xevent.ConfigureNotifyEvent
	updateConfigureTimer     *time.Timer

	propertyNotifyTimer          *time.Timer
	propertyNotifyAtomTable      map[xproto.Atom]bool
	propertyNotifyAtomTableMutex sync.Mutex
	propertyNotifyEnabled        bool

	wmState           []string
	wmWindowType      []string
	wmAllowedActions  []string
	hasXEmbedInfo     bool
	hasWmTransientFor bool
	mapState          byte
	wmClass           *icccm.WmClass
	wmName            string

	gtkAppId string
	wmRole   string
	pid      uint
	process  *ProcessInfo
	entry    *AppEntry

	firstUpdate bool

	entryInnerId string
	appInfo      *AppInfo
}

func NewWindowInfo(win xproto.Window) *WindowInfo {
	winInfo := &WindowInfo{
		window:   win,
		mapState: xproto.MapStateUnmapped,
	}
	return winInfo
}

// window type
func (winInfo *WindowInfo) updateWmWindowType() {
	var err error
	winInfo.wmWindowType, err = ewmh.WmWindowTypeGet(XU, winInfo.window)
	if err != nil {
		logger.Debug(err)
	}
}

// wm allowed actions
func (winInfo *WindowInfo) updateWmAllowedActions() {
	var err error
	winInfo.wmAllowedActions, err = ewmh.WmAllowedActionsGet(XU, winInfo.window)
	if err != nil {
		logger.Debug(err)
	}
}
func (winInfo *WindowInfo) isActionMinimizeAllowed() bool {
	logger.Debugf("wmAllowedActions: %#v", winInfo.wmAllowedActions)
	return strSliceContains(winInfo.wmAllowedActions, "_NET_WM_ACTION_MINIMIZE")
}

// wm state
func (winInfo *WindowInfo) updateWmState() {
	var err error
	winInfo.wmState, err = ewmh.WmStateGet(XU, winInfo.window)
	if err != nil {
		logger.Debug(err)
	}
}

func (winInfo *WindowInfo) hasWmStateSkipTaskbar() bool {
	return strSliceContains(winInfo.wmState, "_NET_WM_STATE_SKIP_TASKBAR")
}

func (winInfo *WindowInfo) hasWmStateModal() bool {
	return strSliceContains(winInfo.wmState, "_NET_WM_STATE_MODAL")
}

func (winInfo *WindowInfo) isValidModal() bool {
	return winInfo.hasWmTransientFor && winInfo.hasWmStateModal()
}

// map state
func (winInfo *WindowInfo) updateMapState() {
	windowAttributes, err := xproto.GetWindowAttributes(XU.Conn(), winInfo.window).Reply()
	if err != nil {
		logger.Debug(err)
		return
	}
	winInfo.mapState = windowAttributes.MapState
	logger.Debug("update map state:", winInfo.mapState)
}

func (winInfo *WindowInfo) isMapStateViewable() bool {
	logger.Debugf("mapState: %v", winInfo.mapState)
	return winInfo.mapState == xproto.MapStateViewable
}

// wm class
func (winInfo *WindowInfo) updateWmClass() {
	var err error
	winInfo.wmClass, err = icccm.WmClassGet(XU, winInfo.window)
	if err != nil {
		logger.Debug(err)
	}
}

// 通过 wmClass 判断是否需要隐藏此窗口
func (winInfo *WindowInfo) isWmClassOk() bool {
	wmClass := winInfo.wmClass
	logger.Debugf("wmClass: %#v", wmClass)
	if wmClass == nil {
		return true
	}
	if wmClass.Instance == "explorer.exe" && wmClass.Class == "Wine" {
		return false
	}
	return true
}

// xembed info
// 一般 trayicon 会带有 _XEMBED_INFO 属性
func (winInfo *WindowInfo) updateHasXEmbedInfo() {
	_, err := xprop.GetProperty(XU, winInfo.window, "_XEMBED_INFO")
	winInfo.hasXEmbedInfo = (err == nil)
}

// WM_TRANSIENT_FOR
func (winInfo *WindowInfo) updateHasWmTransientFor() {
	_, err := xprop.GetProperty(XU, winInfo.window, "WM_TRANSIENT_FOR")
	winInfo.hasWmTransientFor = (err == nil)
}

// wm name
func (winInfo *WindowInfo) updateWmName() {
	winInfo.wmName = getWmName(XU, winInfo.window)
	winInfo.Title = winInfo.getTitle()
	entry := winInfo.entry
	if entry != nil {
		entry.updateWindowTitles()
	}
}

func (winInfo *WindowInfo) getDisplayName() string {
	return strings.Title(winInfo._getDisplayName())
}

func (winInfo *WindowInfo) _getDisplayName() string {
	win := winInfo.window
	role := winInfo.wmRole
	if !utf8.ValidString(role) {
		role = ""
	}

	var class, instance string
	if winInfo.wmClass != nil {
		class = winInfo.wmClass.Class
		if !utf8.ValidString(class) {
			class = ""
		}

		instance = filepath.Base(winInfo.wmClass.Instance)
		if !utf8.ValidString(instance) {
			instance = ""
		}
	}
	logger.Debugf("getDisplayName class: %q, instance: %q", class, instance)

	if role != "" && class != "" {
		return class + " " + role
	}

	if class != "" {
		return class
	}

	if instance != "" {
		return instance
	}

	wmName := winInfo.wmName
	if wmName != "" {
		var shortWmName string
		rindex := strings.LastIndex(wmName, "-")
		if rindex > 0 {
			shortWmName = wmName[rindex:]
			if shortWmName != "" && utf8.ValidString(shortWmName) {
				return shortWmName
			}
		}
	}

	if winInfo.process != nil {
		exeBasename := filepath.Base(winInfo.process.exe)
		if utf8.ValidString(exeBasename) {
			return exeBasename
		}
	}

	return fmt.Sprintf("window: %v", win)
}

func (winInfo *WindowInfo) getTitle() string {
	wmName := winInfo.wmName
	if wmName == "" || !utf8.ValidString(wmName) {
		if winInfo.entry != nil {
			appInfo := winInfo.entry.appInfo
			if appInfo != nil {
				return appInfo.GetDisplayName()
			}
		}
		// winInfo.entry is nil
		return winInfo.getDisplayName()
	}
	return wmName
}

func (winInfo *WindowInfo) getIcon() string {
	if winInfo.Icon == "" {
		logger.Debug("get icon from window", winInfo.window)
		winInfo.Icon = getIconFromWindow(XU, winInfo.window)
	}
	return winInfo.Icon
}

var skipTaskbarWindowTypes []string = []string{
	"_NET_WM_WINDOW_TYPE_UTILITY",
	"_NET_WM_WINDOW_TYPE_COMBO",
	"_NET_WM_WINDOW_TYPE_DESKTOP",
	"_NET_WM_WINDOW_TYPE_DND",
	"_NET_WM_WINDOW_TYPE_DOCK",
	"_NET_WM_WINDOW_TYPE_DROPDOWN_MENU",
	"_NET_WM_WINDOW_TYPE_MENU",
	"_NET_WM_WINDOW_TYPE_NOTIFICATION",
	"_NET_WM_WINDOW_TYPE_POPUP_MENU",
	"_NET_WM_WINDOW_TYPE_SPLASH",
	"_NET_WM_WINDOW_TYPE_TOOLBAR",
	"_NET_WM_WINDOW_TYPE_TOOLTIP",
}

func (winInfo *WindowInfo) canShowOnDock() bool {
	logger.Debug("canShowOnDock win", winInfo.window)
	if !winInfo.firstUpdate {
		winInfo.update()
		winInfo.firstUpdate = true
	}

	logger.Debugf("hasXEmbedInfo: %v", winInfo.hasXEmbedInfo)
	logger.Debugf("wmWindowType: %#v", winInfo.wmWindowType)
	logger.Debugf("wmState: %#v", winInfo.wmState)

	if winInfo.hasWmStateSkipTaskbar() || winInfo.isValidModal() ||
		winInfo.hasXEmbedInfo || !winInfo.isWmClassOk() {
		return false
	}

	for _, winType := range winInfo.wmWindowType {
		if winType == "_NET_WM_WINDOW_TYPE_DIALOG" &&
			!winInfo.isActionMinimizeAllowed() {
			return false
		} else if strSliceContains(skipTaskbarWindowTypes, winType) {
			return false
		}
	}
	return true
}

func (winInfo *WindowInfo) initProcessInfo() {
	win := winInfo.window
	winInfo.pid = getWmPid(XU, win)
	var err error
	winInfo.process, err = NewProcessInfo(winInfo.pid)
	if err != nil {
		logger.Debug(err)
		// Try WM_COMMAND
		wmCommand, err := getWmCommand(XU, win)
		if err == nil {
			winInfo.process = NewProcessInfoWithCmdline(wmCommand)
		}
	}
	logger.Debugf("process: %#v", winInfo.process)
}

func (winInfo *WindowInfo) update() {
	win := winInfo.window
	logger.Debugf("update window %v info", win)
	winInfo.updateMapState()
	winInfo.updateWmClass()
	winInfo.updateWmState()
	winInfo.updateWmWindowType()
	winInfo.updateWmAllowedActions()
	if len(winInfo.wmWindowType) == 0 {
		winInfo.updateHasXEmbedInfo()
	}
	winInfo.updateHasWmTransientFor()
	winInfo.initProcessInfo()
	winInfo.wmRole = getWmWindowRole(XU, win)
	winInfo.gtkAppId = getWindowGtkApplicationId(XU, win)
	winInfo.updateWmName()
	winInfo.genInnerId()
}

func _filterFilePath(args []string) string {
	var filtered []string
	for _, arg := range args {
		if strings.Contains(arg, "/") || arg == "." || arg == ".." {
			filtered = append(filtered, "%F")
		} else {
			filtered = append(filtered, arg)
		}
	}
	return strings.Join(filtered, " ")
}

func (winInfo *WindowInfo) genInnerId() {
	win := winInfo.window
	var wmClass string
	var wmInstance string
	if winInfo.wmClass != nil {
		wmClass = winInfo.wmClass.Class
		wmInstance = filepath.Base(winInfo.wmClass.Instance)
	}
	var exe string
	var args string
	if winInfo.process != nil {
		exe = winInfo.process.exe
		args = _filterFilePath(winInfo.process.args)
	}
	hasPid := (winInfo.pid != 0)

	var str string
	// NOTE: 不要使用 wmRole，有些程序总会改变这个值比如 GVim
	if wmInstance == "" && wmClass == "" && exe == "" && winInfo.gtkAppId == "" {
		if winInfo.wmName != "" {
			str = fmt.Sprintf("wmName:%q", winInfo.wmName)
		} else {
			str = fmt.Sprintf("windowId:%v", winInfo.window)
		}
	} else {
		str = fmt.Sprintf("wmInstance:%q,wmClass:%q,exe:%q,args:%q,hasPid:%v,gtkAppId:%q",
			wmInstance, wmClass, exe, args, hasPid, winInfo.gtkAppId)
	}

	hasher := md5.New()
	hasher.Write([]byte(str))
	winInfo.innerId = windowHashPrefix + hex.EncodeToString(hasher.Sum(nil))
	logger.Debugf("genInnerId win: %v str: %s, md5sum: %s", win, str, winInfo.innerId)
}

func (winInfo *WindowInfo) initPropertyNotifyEventHandler(dockManager *DockManager) {
	if winInfo.propertyNotifyTimer != nil {
		return
	}

	winInfo.propertyNotifyAtomTable = make(map[xproto.Atom]bool)
	winInfo.propertyNotifyEnabled = false
	// simulate first property notify event
	winInfo.propertyNotifyAtomTable[ATOM_WINDOW_STATE] = true

	winInfo.propertyNotifyTimer = time.AfterFunc(300*time.Millisecond, func() {
		var atomNames []string
		var needUpdate bool
		winInfo.propertyNotifyAtomTableMutex.Lock()
		defer winInfo.propertyNotifyAtomTableMutex.Unlock()

		for atom, _ := range winInfo.propertyNotifyAtomTable {
			atomName, _ := xprop.AtomName(XU, atom)
			atomNames = append(atomNames, atomName)
			if winInfo.handlePropertyNotifyAtom(atom) {
				// may changed
				needUpdate = true
			}
		}
		logger.Debugf("propertyNotifyAtomTable win %v atom: %v", winInfo.window, atomNames)

		if needUpdate {
			dockManager.attachOrDetachWindow(winInfo)
		}

		// end
		winInfo.propertyNotifyAtomTable = make(map[xproto.Atom]bool)
		winInfo.propertyNotifyEnabled = true
	})
}

func (winInfo *WindowInfo) handlePropertyNotifyEvent(ev xevent.PropertyNotifyEvent) {
	winInfo.propertyNotifyAtomTableMutex.Lock()
	winInfo.propertyNotifyAtomTable[ev.Atom] = true
	winInfo.propertyNotifyAtomTableMutex.Unlock()

	if winInfo.propertyNotifyEnabled {
		winInfo.propertyNotifyTimer.Reset(300 * time.Millisecond)
		winInfo.propertyNotifyEnabled = false
	}
}

func (winInfo *WindowInfo) handlePropertyNotifyAtom(atom xproto.Atom) bool {
	switch atom {
	case ATOM_WINDOW_STATE:
		winInfo.updateWmState()
		return true

	case ATOM_WINDOW_TYPE:
		winInfo.updateWmWindowType()
		return true

	case ATOM_XEMBED_INFO:
		winInfo.updateHasXEmbedInfo()
		return true

	case ATOM_WINDOW_NAME:
		winInfo.updateWmName()
		return false

	case ATOM_WINDOW_ICON:
		//  update icon cache
		winInfo.Icon = getIconFromWindow(XU, winInfo.window)
		entry := winInfo.entry
		if entry != nil && entry.current == winInfo {
			entry.updateIcon()
		}
		return false
	default:
		return false
	}
}
