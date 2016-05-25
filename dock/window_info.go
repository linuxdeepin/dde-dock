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
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

type WindowInfo struct {
	window xproto.Window
	Title  string
	Icon   string

	x                        int16
	y                        int16
	width                    uint16
	height                   uint16
	lastConfigureNotifyEvent *xevent.ConfigureNotifyEvent
	updateConfigureTimer     *time.Timer

	propertyNotifyTimer     *time.Timer
	propertyNotifyAtomTable map[xproto.Atom]bool
	propertyNotifyEnabled   bool

	wmState          []string
	wmWindowType     []string
	wmAllowedActions []string
	hasXEmbedInfo    bool
	mapState         byte
	wmClass          *icccm.WmClass
	wmName           string

	process *ProcessInfo
	app     *RuntimeApp

	firstUpdate bool
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
	return contains(winInfo.wmAllowedActions, "_NET_WM_ACTION_MINIMIZE")
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
	return contains(winInfo.wmState, "_NET_WM_STATE_SKIP_TASKBAR")
}

func (winInfo *WindowInfo) hasWmStateModal() bool {
	return contains(winInfo.wmState, "_NET_WM_STATE_MODAL")
}

// map state
func (winInfo *WindowInfo) updateMapState() {
	windowAttributes, err := xproto.GetWindowAttributes(XU.Conn(), winInfo.window).Reply()
	if err != nil {
		logger.Debug(err)
		return
	}
	winInfo.mapState = windowAttributes.MapState
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

// wm name
func (winInfo *WindowInfo) updateWmName() {
	winInfo.wmName = getWmName(XU, winInfo.window)
	winInfo.Title = winInfo.getTitle()
	if winInfo.app != nil {
		winInfo.app.notifyChanged()
	}
}

func (winInfo *WindowInfo) getTitle() string {
	wmName := winInfo.wmName
	if wmName == "" || !utf8.ValidString(wmName) {
		if winInfo.app != nil {
			appInfo := winInfo.app.appInfo
			if appInfo != nil {
				return appInfo.GetDisplayName()
			}
			return winInfo.app.Id
		}
		// winInfo.app is nil
		return fmt.Sprintf("window: %v", winInfo.window)
	} else {
		return wmName
	}
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

	if !winInfo.isMapStateViewable() || !winInfo.isWmClassOk() ||
		winInfo.hasWmStateSkipTaskbar() || winInfo.hasWmStateModal() ||
		winInfo.hasXEmbedInfo {
		return false
	}

	for _, winType := range winInfo.wmWindowType {
		if winType == "_NET_WM_WINDOW_TYPE_DIALOG" &&
			!winInfo.isActionMinimizeAllowed() {
			return false
		} else if contains(skipTaskbarWindowTypes, winType) {
			return false
		}
	}
	return true
}

func (winInfo *WindowInfo) initProcessInfo() {
	win := winInfo.window
	pid := getWmPid(XU, win)
	var err error
	winInfo.process, err = NewProcessInfo(pid)
	if err != nil {
		logger.Warning(err)
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
	winInfo.updateWmName()
	winInfo.initProcessInfo()
}

func (winInfo *WindowInfo) guessAppId(filterGroup *AppIdFilterGroup) string {
	var appId string
	if winInfo.process == nil {
		return ""
	}
	process := winInfo.process
	execName := filepath.Base(process.exe)
	// wm name
	appId = findAppIdByFilter(execName, winInfo.wmName, filterGroup.KeyFileWmName)
	if appId != "" {
		return appId
	}

	if winInfo.wmClass != nil {
		// wmclass instance
		appId = findAppIdByFilter(execName, winInfo.wmClass.Instance, filterGroup.KeyFileWmInstance)
		if appId != "" {
			return appId
		}

		// wmclass class
		appId = findAppIdByFilter(execName, winInfo.wmClass.Class, filterGroup.KeyFileWmClass)
		if appId != "" {
			return appId
		}
	}

	// args
	argsJoined := strings.Join(process.args, " ")
	appId = findAppIdByFilter(execName, argsJoined, filterGroup.KeyFileArgs)
	if appId != "" {
		return appId
	}

	// icon name
	iconName, _ := ewmh.WmIconNameGet(XU, winInfo.window)
	appId = findAppIdByFilter(execName, iconName, filterGroup.KeyFileIconName)
	if appId != "" {
		return appId
	}

	// exec name
	appId = findAppIdByFilter(execName, execName, filterGroup.KeyFileExecName)
	return appId
}

func (winInfo *WindowInfo) createAppId() string {
	var appId string
	if winInfo.wmClass != nil {
		appId = winInfo.wmClass.Class
		// it is possible that getting invalid string which might be xgb implementation's bug.
		// for instance: xdemineur's WMClass
		if appId != "" && utf8.ValidString(appId) {
			return normalizeAppID(appId)
		}

		appId = winInfo.wmClass.Instance
		if appId != "" && utf8.ValidString(appId) {
			// wmclass instance 可能是文件路径，比如 xdemineur 的
			appId = filepath.Base(appId)
			return normalizeAppID(appId)
		}
	}

	if winInfo.process != nil {
		appId = filepath.Base(winInfo.process.exe)
		if appId != "" {
			return normalizeAppID(appId)
		}
	}
	// TODO: try WM_COMMAND
	// TODO: try ICON NAME
	return ""
}

func (winInfo *WindowInfo) initPropertyNotifyEventHandler(entryManager *EntryManager) {
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
			entryManager.attachOrDetachRuntimeAppWindow(winInfo)
		}

		// end
		winInfo.propertyNotifyAtomTable = make(map[xproto.Atom]bool)
		winInfo.propertyNotifyEnabled = true
	})
}

func (winInfo *WindowInfo) handlePropertyNotifyEvent(ev xevent.PropertyNotifyEvent) {
	winInfo.propertyNotifyAtomTable[ev.Atom] = true
	if winInfo.propertyNotifyEnabled {
		winInfo.propertyNotifyTimer.Reset(300 * time.Millisecond)
		winInfo.propertyNotifyEnabled = false
	}
}

func (winInfo *WindowInfo) handlePropertyNotifyAtom(atom xproto.Atom) bool {
	app := winInfo.app

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
		if app != nil {
			app.updateIcon(winInfo)
			app.notifyChanged()
		}
		return false
		// case ATOM_DOCK_APP_ID:
		// 	// TODO DOCK_APP_ID?
		// 	logger.Debug("PropertyNotifyEvent ATOM_DOCK_APP_ID")
		// 	if app != nil {
		// 		app.updateAppid(xid)
		// 		app.notifyChanged()
		// 	}
	default:
		return false
	}
}
