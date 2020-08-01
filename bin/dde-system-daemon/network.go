package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"sort"

	. "pkg.deepin.io/dde/daemon/common/dsync"
	"pkg.deepin.io/gir/glib-2.0"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	nmConnDir = "/etc/NetworkManager/system-connections"

	kfSectionConnection = "connection"
	kfSectionWIFI       = "wifi"
	kfKeyType           = "type"
	kfKeyMac            = "mac-address"
	kfKeyMacBlacklist   = "mac-address-blacklist"
	kfKeySeenBSSID      = "seen-bssids"
	kfKeyInterfaceName  = "interface-name"

	nmSyncVersion = "1.0"
)

func (*Daemon) NetworkGetConnections() ([]byte, *dbus.Error) {
	list, err := getConnectionList(nmConnDir)
	if err != nil {
		return nil, dbusutil.ToError(err)
	}
	var info = NetworkData{
		Version:     nmSyncVersion,
		Connections: list,
	}
	data, err := json.Marshal(&info)
	if err != nil {
		return nil, dbusutil.ToError(err)
	}
	return data, nil
}

func (*Daemon) NetworkSetConnections(data []byte) *dbus.Error {
	var info NetworkData
	err := json.Unmarshal(data, &info)
	if err != nil {
		return dbusutil.ToError(err)
	}
	err = info.Connections.Check()
	if err != nil {
		return dbusutil.ToError(err)
	}
	list, _ := getConnectionList(nmConnDir)
	for _, conn := range info.Connections {
		tmp := list.Get(conn.Type, conn.Filename)
		if tmp != nil && tmp.Equal(conn) {
			continue
		}
		// add or modify
		err = conn.WriteFile(nmConnDir)
		if err != nil {
			// TODO(jouyouyun): handle error
			logger.Error("[Network] Failed to write connection file:",
				err)
			// return dbusutil.ToError(err)
			continue
		}
	}
	removeList := info.Connections.Diff(list)
	for _, conn := range removeList {
		err = conn.RemoveFile(nmConnDir)
		if err != nil {
			// TODO(jouyouyun): handle error
			logger.Error("[Network] Failed to remove connection file:",
				err)
			continue
		}
	}
	err = reloadConnections()
	if err != nil {
		logger.Warning("[Network] Failed to reload connections:", err)
	}
	return nil
}

func getConnectionList(dir string) (ConnectionList, error) {
	files, err := getConnectionFiles(dir)
	if err != nil {
		return nil, err
	}

	var datas ConnectionList
	for _, f := range files {
		data, err := loadConnectionFile(f)
		if err != nil {
			continue
		}
		if datas.Exists(data) {
			continue
		}
		datas = append(datas, data)
	}
	sort.Sort(datas)
	return datas, nil
}

func loadConnectionFile(filename string) (*Connection, error) {
	var kf = glib.NewKeyFile()
	// ignore comments and translations
	_, err := kf.LoadFromFile(filename, glib.KeyFileFlagsNone)
	if err != nil {
		return nil, err
	}
	defer kf.Free()

	ty, err := kf.GetString(kfSectionConnection, kfKeyType)
	if err != nil {
		return nil, err
	}
	if ty != ConnTypeWIFI {
		return nil, ErrConnUnsupportedType
	}

	// erase some keys
	_, _ = kf.RemoveKey(kfSectionConnection, kfKeyInterfaceName)
	_, _ = kf.RemoveKey(kfSectionWIFI, kfKeyMac)
	_, _ = kf.RemoveKey(kfSectionWIFI, kfKeyMacBlacklist)
	_, _ = kf.RemoveKey(kfSectionWIFI, kfKeySeenBSSID)

	_, contents, err := kf.ToData()
	if err != nil {
		return nil, err
	}

	return &Connection{
		Type:     ty,
		Filename: filepath.Base(filename),
		Contents: []byte(contents),
	}, nil
}

func getConnectionFiles(dir string) ([]string, error) {
	finfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, finfo := range finfos {
		if finfo.IsDir() {
			continue
		}
		files = append(files, filepath.Join(dir, finfo.Name()))
	}
	return files, nil
}
