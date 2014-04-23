/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
        "dlib"
        "dlib/dbus"
        "dlib/gio-2.0"
        "dlib/gobject-2.0"
        "dlib/logger"
        "fmt"
        "os"
        "sync"
)

type DiskInfo struct {
        Id         int32
        Name       string
        Type       string
        CanUnmount bool
        CanEject   bool
        UsableCap  int64
        TotalCap   int64
        Path       string
        UUID       string
}

type ObjectInfo struct {
        Object  interface{}
        Type    string
}

type Manager struct {
        DiskList []DiskInfo
        Error    func(string)
}

const (
        DEVICE_KIND = "unix-device"
)

var (
        monitor   = gio.VolumeMonitorGet()
        objectMap = make(map[int32]*ObjectInfo)
        logObject = logger.NewLogger("daemon/mounts")
        mutex     = new(sync.Mutex)

        genID, destroyID = func() (func() int32, func()) {
                count := int32(0)
                return func() int32 {
                                mutex.Lock()
                                tmp := count
                                count += 1
                                mutex.Unlock()
                                return tmp
                        }, func() {
                                mutex.Lock()
                                count = 0
                                mutex.Unlock()
                        }
        }()
)

func (m *Manager) DeviceEject(id int32) (bool, string) {
        mutex.Lock()
        defer mutex.Unlock()
        info, ok := objectMap[id]
        if !ok {
                logObject.Infof("Eject id - %d not in objectMap.", id)
                return false, fmt.Sprintf("Invalid Id: %d\n", id)
        }

        logObject.Infof("Eject type: %s", info.Type)
        switch info.Type {
        case "drive":
                op := info.Object.(*gio.Drive)
                op.Eject(gio.MountUnmountFlagsNone, nil, gio.AsyncReadyCallback(func(o *gobject.Object, res *gio.AsyncResult) {
                        _, err := op.EjectFinish(res)
                        if err != nil {
                                m.Error(err.Error())
                                logObject.Infof("drive eject failed: %d, %s", id, err)
                        }
                }))
        case "volume":
                op := info.Object.(*gio.Volume)
                op.Eject(gio.MountUnmountFlagsNone, nil, gio.AsyncReadyCallback(func(o *gobject.Object, res *gio.AsyncResult) {
                        _, err := op.EjectFinish(res)
                        if err != nil {
                                m.Error(err.Error())
                                logObject.Infof("volume eject failed: %d, %s", id, err)
                        }
                }))
        case "mount":
                op := info.Object.(*gio.Mount)
                op.Eject(gio.MountUnmountFlagsNone, nil, gio.AsyncReadyCallback(func(o *gobject.Object, res *gio.AsyncResult) {
                        _, err := op.EjectFinish(res)
                        if err != nil {
                                m.Error(err.Error())
                                logObject.Infof("mount eject failed: %d, %s", id, err)
                        }
                }))
        default:
                logObject.Infof("'%s' invalid type", info.Type)
                return false, fmt.Sprintf("Invalid type: '%s'\n", info.Type)
        }

        return true, ""
}

func (m *Manager) DeviceMount(id int32) (bool, string) {
        mutex.Lock()
        defer mutex.Unlock()
        info, ok := objectMap[id]
        if !ok {
                logObject.Infof("Mount id - %d not in objectMap.", id)
                return false, fmt.Sprintf("Invalid Id: %d\n", id)
        }

        logObject.Infof("Mount type: %s", info.Type)
        switch info.Type {
        case "volume":
                op := info.Object.(*gio.Volume)
                op.Mount(gio.MountMountFlagsNone, nil, nil, gio.AsyncReadyCallback(func(o *gobject.Object, res *gio.AsyncResult) {
                        _, err := op.MountFinish(res)
                        if err != nil {
                                m.Error(err.Error())
                                logObject.Infof("volume mount failed: %d, %s", id, err)
                        }
                }))
        case "mount":
                op := info.Object.(*gio.Mount)
                op.Remount(gio.MountMountFlagsNone, nil, nil, gio.AsyncReadyCallback(func(o *gobject.Object, res *gio.AsyncResult) {
                        _, err := op.RemountFinish(res)
                        if err != nil {
                                m.Error(err.Error())
                                logObject.Infof("mount remount failed: %d, %s", id, err)
                        }
                }))
        default:
                logObject.Infof("'%s' invalid type", info.Type)
                return false, fmt.Sprintf("Invalid type: '%s'\n", info.Type)
        }

        return true, ""
}

func (m *Manager) DeviceUnmount(id int32) (bool, string) {
        mutex.Lock()
        defer mutex.Unlock()
        info, ok := objectMap[id]
        if !ok {
                logObject.Infof("Unmount id - %d not in objectMap.", id)
                return false, fmt.Sprintf("Invalid Id: %d\n", id)
        }

        logObject.Infof("Unmount type: %s", info.Type)
        switch info.Type {
        case "mount":
                op := info.Object.(*gio.Mount)
                op.Unmount(gio.MountUnmountFlagsNone, nil, gio.AsyncReadyCallback(func(o *gobject.Object, res *gio.AsyncResult) {
                        _, err := op.UnmountFinish(res)
                        if err != nil {
                                m.Error(err.Error())
                                logObject.Infof("mount unmount failed: %d, %s", id, err)
                        }
                }))
        default:
                logObject.Infof("'%s' invalid type", info.Type)
                return false, fmt.Sprintf("Invalid type: '%s'\n", info.Type)
        }

        return true, ""
}

func newDiskInfo(value interface{}, t string, id int32) DiskInfo {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Warning("Received Error: ", err)
                }
        }()

        info := DiskInfo{}
        info.Id = id

        switch t {
        case "volume":
                v := value.(*gio.Volume)
                info.Name = v.GetName()
                info.CanEject = v.CanEject()
                id := v.GetIdentifier(DEVICE_KIND)
                info.TotalCap, info.UsableCap = getDiskCap(id)
                info.Path = v.GetIdentifier(gio.VolumeIdentifierKindUnixDevice)
                info.UUID = v.GetIdentifier(gio.VolumeIdentifierKindUuid)
                if mount := v.GetMount(); mount != nil {
                        info.CanUnmount = mount.CanUnmount()
                }
                if containStart("network", id) {
                        info.Type = "network"
                } else if info.CanEject {
                        info.Type = "removable"
                } else {
                        info.Type = "native"
                }
        case "drive":
                v := value.(*gio.Drive)
                info.Name = v.GetName()
                info.CanEject = v.CanEject()
                id := v.GetIdentifier(DEVICE_KIND)
                info.TotalCap, info.UsableCap = getDiskCap(id)
                info.Path = v.GetIdentifier(gio.VolumeIdentifierKindUnixDevice)
                info.UUID = v.GetIdentifier(gio.VolumeIdentifierKindUuid)
                if containStart("network", id) {
                        info.Type = "network"
                } else if info.CanEject {
                        info.Type = "removable"
                } else {
                        info.Type = "native"
                }
        case "mount":
                v := value.(*gio.Mount)
                info.Name = v.GetName()
                info.CanEject = v.CanEject()
                info.CanUnmount = v.CanUnmount()
                root := v.GetRoot()
                info.TotalCap, info.UsableCap = getDiskCap(root.GetPath())
                if info.CanEject {
                        info.Type = "removable"
                } else if root.IsNative() {
                        info.Type = "native"
                } else {
                        info.Type = "network"
                }
                if volume := v.GetVolume(); volume != nil {
                        info.Path = volume.GetIdentifier(gio.VolumeIdentifierKindUnixDevice)
                        info.UUID = volume.GetIdentifier(gio.VolumeIdentifierKindUuid)
                }
        default:
                logObject.Infof("'%s' invalid type", t)
        }

        return info
}

func newObjectInfo(v interface{}, t string) *ObjectInfo {
        return &ObjectInfo{Object: v, Type: t}
}

func driverList() []DiskInfo {
        list := []DiskInfo{}
        drivers := monitor.GetConnectedDrives()
        for _, driver := range drivers {
                volumes := driver.GetVolumes()
                if volumes == nil {
                        if driver.IsMediaRemovable() &&
                                !driver.IsMediaCheckAutomatic() {
                                info := newDiskInfo(driver, "drive", genID())
                                objectMap[info.Id] = newObjectInfo(driver, "drive")
                                list = append(list, info)
                        }
                        continue
                }
                for _, volume := range volumes {
                        mount := volume.GetMount()
                        if mount != nil {
                                info := newDiskInfo(mount, "mount", genID())
                                objectMap[info.Id] = newObjectInfo(mount, "mount")
                                list = append(list, info)
                        } else {
                                info := newDiskInfo(volume, "volume", genID())
                                objectMap[info.Id] = newObjectInfo(volume, "volume")
                                list = append(list, info)
                        }
                }
        }

        return list
}

func volumeList() []DiskInfo {
        list := []DiskInfo{}
        volumes := monitor.GetVolumes()
        for _, volume := range volumes {
                driver := volume.GetDrive()
                if driver != nil {
                        continue
                }
                //id := volume.GetIdentifier("unix-device")
                mount := volume.GetMount()
                if mount != nil {
                        info := newDiskInfo(mount, "mount", genID())
                        objectMap[info.Id] = newObjectInfo(mount, "mount")
                        list = append(list, info)
                } else {
                        info := newDiskInfo(volume, "volume", genID())
                        objectMap[info.Id] = newObjectInfo(volume, "volume")
                        list = append(list, info)
                }
        }
        return list
}

func mountList() []DiskInfo {
        list := []DiskInfo{}
        mounts := monitor.GetMounts()
        for _, mount := range mounts {
                if mount.IsShadowed() {
                        continue
                }

                volume := mount.GetVolume()
                if volume != nil {
                        continue
                }
                info := newDiskInfo(mount, "mount", genID())
                objectMap[info.Id] = newObjectInfo(mount, "mount")
                list = append(list, info)
        }
        return list
}

func containStart(str1, str2 string) bool {
        for i, _ := range str1 {
                if str1[i] != str2[i] {
                        return false
                }
        }

        return true
}

func getDiskInfoList() []DiskInfo {
        list := []DiskInfo{}

        destroyObjectMap()
        l1 := driverList()
        l2 := volumeList()
        l3 := mountList()
        list = append(list, l1...)
        list = append(list, l2...)
        list = append(list, l3...)

        return list
}

func destroyObjectMap() {
        for k, info := range objectMap {
                switch info.Type {
                case "drive":
                        op := info.Object.(*gio.Drive)
                        op.Unref()
                case "volume":
                        op := info.Object.(*gio.Volume)
                        op.Unref()
                case "mount":
                        op := info.Object.(*gio.Mount)
                        op.Unref()
                }
                delete(objectMap, k)
        }
        destroyID()
}

func NewManager() *Manager {
        m := &Manager{}
        m.setPropName("DiskList")
        m.listenSignalChanged()

        //printDiskInfo(m.DiskList)
        return m
}

func main() {
        defer logObject.EndTracing()

        if !dlib.UniqueOnSession(DISK_INFO_DEST) {
                logObject.Warning("There already has an Mount daemon running.")
                return
        }

        logObject.SetRestartCommand("/usr/lib/deepin-daemon/mounts")

        m := NewManager()
        err := dbus.InstallOnSession(m)
        if err != nil {
                logObject.Info("Install DBus Session Failed:", err)
                panic(err)
        }
        dbus.DealWithUnhandledMessage()

        go dlib.StartLoop()
        if err = dbus.Wait(); err != nil {
                logObject.Info("lost dbus session:", err)
                os.Exit(1)
        } else {
                os.Exit(0)
        }
}

func printDiskInfo(infos []DiskInfo) {
        for _, v := range infos {
                logObject.Infof("Id: %d", v.Id)
                logObject.Infof("Name: %s", v.Name)
                logObject.Infof("Type: %s", v.Type)
                logObject.Info("CanEject:", v.CanEject)
                logObject.Info("CanUnmount:", v.CanUnmount)
        }
}
