/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package dock

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/util/wm/ewmh"
	"github.com/linuxdeepin/go-x11-client/util/wm/icccm"
)

const windowHashPrefix = "w:"

type WindowInfo struct {
	innerId string
	window  x.Window
	Title   string
	Icon    string

	x                        int16
	y                        int16
	width                    uint16
	height                   uint16
	lastConfigureNotifyEvent *x.ConfigureNotifyEvent
	mu                       sync.Mutex
	updateConfigureTimer     *time.Timer

	wmState           []x.Atom
	wmWindowType      []x.Atom
	wmAllowedActions  []x.Atom
	hasXEmbedInfo     bool
	hasWmTransientFor bool
	wmClass           *icccm.WMClass
	wmName            string

	gtkAppId     string
	flatpakAppID string
	wmRole       string
	pid          uint
	process      *ProcessInfo
	entry        *AppEntry

	firstUpdate bool

	entryInnerId string
	appInfo      *AppInfo
}

func NewWindowInfo(win x.Window) *WindowInfo {
	winInfo := &WindowInfo{
		window: win,
	}
	return winInfo
}

// window type
func (winInfo *WindowInfo) updateWmWindowType() {
	var err error
	winInfo.wmWindowType, err = ewmh.GetWMWindowType(globalXConn, winInfo.window).Reply(globalXConn)
	if err != nil {
		logger.Debugf("failed to get WMWindowType for window %d: %v", winInfo.window, err)
	}
}

// wm allowed actions
func (winInfo *WindowInfo) updateWmAllowedActions() {
	var err error
	winInfo.wmAllowedActions, err = ewmh.GetWMAllowedActions(globalXConn,
		winInfo.window).Reply(globalXConn)
	if err != nil {
		logger.Debugf("failed to get WMAllowedActions for window %d: %v", winInfo.window, err)
	}
}

// wm state
func (winInfo *WindowInfo) updateWmState() {
	var err error
	winInfo.wmState, err = ewmh.GetWMState(globalXConn, winInfo.window).Reply(globalXConn)
	if err != nil {
		logger.Debugf("failed to get WMState for window %d: %v", winInfo.window, err)
	}
}

// wm class
func (winInfo *WindowInfo) updateWmClass() {
	var err error
	winInfo.wmClass, err = getWmClass(winInfo.window)
	if err != nil {
		logger.Debugf("failed to get wmClass for window %d: %v", winInfo.window, err)
	}
}

// wm name
func (winInfo *WindowInfo) updateWmName() {
	winInfo.wmName = getWmName(winInfo.window)
	winInfo.Title = winInfo.getTitle()
}

func (winInfo *WindowInfo) updateIcon() {
	winInfo.Icon = getIconFromWindow(winInfo.window)
}

// XEmbed info
// 一般 tray icon 会带有 _XEMBED_INFO 属性
func (winInfo *WindowInfo) updateHasXEmbedInfo() {
	reply, err := x.GetProperty(globalXConn, false, winInfo.window, atomXEmbedInfo, x.AtomAny, 0, 2).Reply(globalXConn)
	if err != nil {
		logger.Debug(err)
		return
	}
	if reply.Format != 0 {
		// has property
		winInfo.hasXEmbedInfo = true
	}
}

// WM_TRANSIENT_FOR
func (winInfo *WindowInfo) updateHasWmTransientFor() {
	_, err := icccm.GetWMTransientFor(globalXConn, winInfo.window).Reply(globalXConn)
	winInfo.hasWmTransientFor = err == nil
}

func (winInfo *WindowInfo) isActionMinimizeAllowed() bool {
	logger.Debugf("wmAllowedActions: %#v", winInfo.wmAllowedActions)
	return atomsContains(winInfo.wmAllowedActions, atomNetWmActionMinimize)
}

func (winInfo *WindowInfo) hasWmStateDemandsAttention() bool {
	return atomsContains(winInfo.wmState, atomWmStateDemandsAttention)
}

func (winInfo *WindowInfo) hasWmStateSkipTaskBar() bool {
	return atomsContains(winInfo.wmState, atomNetWmStateSkipTaskbar)
}

func (winInfo *WindowInfo) hasWmStateModal() bool {
	return atomsContains(winInfo.wmState, atomNetWmStateModal)
}

func (winInfo *WindowInfo) isValidModal() bool {
	return winInfo.hasWmTransientFor && winInfo.hasWmStateModal()
}

// 通过 wmClass 判断是否需要隐藏此窗口
func (winInfo *WindowInfo) shouldSkipWithWMClass() bool {
	wmClass := winInfo.wmClass
	if wmClass == nil {
		return false
	}
	if wmClass.Instance == "explorer.exe" && wmClass.Class == "Wine" {
		return true
	} else if wmClass.Class == "dde-launcher" {
		return true
	}

	return false
}

func (winInfo *WindowInfo) getDisplayName() string {
	return strings.Title(winInfo.getDisplayName0())
}

func (winInfo *WindowInfo) getDisplayName0() string {
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
		lastIndex := strings.LastIndex(wmName, "-")
		if lastIndex > 0 {
			shortWmName = wmName[lastIndex:]
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
		return winInfo.getDisplayName()
	}
	return wmName
}

func (winInfo *WindowInfo) getIcon() string {
	if winInfo.Icon == "" {
		logger.Debug("get icon from window", winInfo.window)
		winInfo.Icon = getIconFromWindow(winInfo.window)
	}
	return winInfo.Icon
}

var skipTaskBarWindowTypes = []string{
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

func (winInfo *WindowInfo) shouldSkip() bool {
	logger.Debugf("win %d shouldSkip?", winInfo.window)
	if !winInfo.firstUpdate {
		winInfo.update()
		winInfo.firstUpdate = true
	}

	logger.Debugf("hasXEmbedInfo: %v", winInfo.hasXEmbedInfo)
	logger.Debugf("wmWindowType: %#v", winInfo.wmWindowType)
	logger.Debugf("wmState: %#v", winInfo.wmState)
	logger.Debugf("wmClass: %#v", winInfo.wmClass)

	if winInfo.hasWmStateSkipTaskBar() || winInfo.isValidModal() ||
		winInfo.hasXEmbedInfo || winInfo.shouldSkipWithWMClass() {
		return true
	}

	for _, winType := range winInfo.wmWindowType {
		winTypeStr, _ := getAtomName(winType)
		if winType == atomNetWmWindowTypeDialog &&
			!winInfo.isActionMinimizeAllowed() {
			return true
		} else if strSliceContains(skipTaskBarWindowTypes, winTypeStr) {
			return true
		}
	}
	return false
}

func (winInfo *WindowInfo) initProcessInfo() {
	win := winInfo.window
	winInfo.pid = getWmPid(win)
	var err error
	winInfo.process, err = NewProcessInfo(winInfo.pid)
	if err != nil {
		logger.Debug(err)
		// Try WM_COMMAND
		wmCommand, err := getWmCommand(win)
		if err == nil {
			winInfo.process = NewProcessInfoWithCmdline(wmCommand)
		}
	}
	logger.Debugf("process: %#v", winInfo.process)
}

func (winInfo *WindowInfo) update() {
	win := winInfo.window
	logger.Debugf("update window %v info", win)
	winInfo.updateWmClass()
	winInfo.updateWmState()
	winInfo.updateWmWindowType()
	winInfo.updateWmAllowedActions()
	if len(winInfo.wmWindowType) == 0 {
		winInfo.updateHasXEmbedInfo()
	}
	winInfo.updateHasWmTransientFor()
	winInfo.initProcessInfo()
	winInfo.wmRole = getWmWindowRole(win)
	winInfo.gtkAppId = getWindowGtkApplicationId(win)
	winInfo.flatpakAppID = getWindowFlatpakAppID(win)
	winInfo.updateWmName()
	winInfo.genInnerId()
}

func filterFilePath(args []string) string {
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
		args = filterFilePath(winInfo.process.args)
	}
	hasPid := winInfo.pid != 0

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

	md5hash := md5.New()
	md5hash.Write([]byte(str))
	winInfo.innerId = windowHashPrefix + hex.EncodeToString(md5hash.Sum(nil))
	logger.Debugf("genInnerId win: %v str: %s, md5sum: %s", win, str, winInfo.innerId)
}
