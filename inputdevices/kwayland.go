package inputdevices

import (
	"os"
	"strconv"
	"strings"

	kwin "github.com/linuxdeepin/go-dbus-factory/org.kde.kwin"
	"pkg.deepin.io/dde/api/dxinput"
	"pkg.deepin.io/dde/api/dxinput/common"
	"pkg.deepin.io/dde/api/dxinput/kwayland"
	"pkg.deepin.io/lib/dbusutil"
)

var (
	globalWayland bool
	kwinManager   *kwin.InputDeviceManager
	kwinIdList    []dbusutil.SignalHandlerId
)

func init() {
	if len(os.Getenv("WAYLAND_DISPLAY")) != 0 {
		globalWayland = true
	}
}

func handleInputDeviceChanged(srv *dbusutil.Service, stop bool) {
	if kwinManager == nil && srv != nil {
		kwinManager = kwin.NewInputDeviceManager(srv.Conn())
	}
	if stop && kwinManager != nil {
		for _, id := range kwinIdList {
			kwinManager.RemoveHandler(id)
		}
		return
	}

	kwinManager.InitSignalExt(_manager.sessionSigLoop, true)
	addedID, err := kwinManager.ConnectDeviceAdded(func(sysName string) {
		doHandleKWinDeviceAdded(sysName)
	})
	if err == nil {
		kwinIdList = append(kwinIdList, addedID)
	}
	removedID, err := kwinManager.ConnectDeviceRemoved(func(sysName string) {
		logger.Info("[Device Removed]:", sysName)
		doHandleKWinDeviceRemoved(sysName)
	})
	if err == nil {
		kwinIdList = append(kwinIdList, removedID)
	}
}

func doHandleKWinDeviceAdded(sysName string) {
	info, err := kwayland.NewDeviceInfo(sysName)
	if err != nil {
		logger.Debug("[Device Added] Failed to new device:", err)
		return
	}
	if info == nil {
		return
	}
	logger.Info("[Device Changed] added:", info.Id, info.Type, info.Name, info.Enabled)
	switch info.Type {
	case common.DevTypeMouse:
		v, _ := dxinput.NewMouseFromDeviceInfo(info)
		mouseInfos = append(mouseInfos, v)
		_manager.mouse.handleDeviceChanged()
	case common.DevTypeTouchpad:
		v, _ := dxinput.NewTouchpadFromDevInfo(info)
		tpadInfos = append(tpadInfos, v)
		_manager.tpad.handleDeviceChanged()
	}
}

func doHandleKWinDeviceRemoved(sysName string) {
	str := strings.TrimLeft(sysName, kwayland.SysNamePrefix)
	id, _ := strconv.Atoi(str)
	logger.Info("----------------items:", sysName, str, id)

	var minfos dxMouses
	var changed bool
	for _, tmp := range mouseInfos {
		if tmp.Id == int32(id) {
			changed = true
			continue
		}
		minfos = append(minfos, tmp)
	}
	if changed {
		logger.Info("[Device Removed] mouse:", sysName)
		mouseInfos = minfos
		_manager.mouse.handleDeviceChanged()
		return
	}

	var tinfos dxTouchpads
	for _, tmp := range tpadInfos {
		if tmp.Id == int32(id) {
			changed = true
			continue
		}
		tinfos = append(tinfos, tmp)
	}
	if changed {
		logger.Info("[Device Removed] touchpad:", sysName)
		tpadInfos = tinfos
		_manager.tpad.handleDeviceChanged()
	}
}
