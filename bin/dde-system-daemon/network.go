package main

import (
	"encoding/json"
	"gir/glib-2.0"
	"io/ioutil"
	"path/filepath"
	. "pkg.deepin.io/dde/daemon/common/dsync"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"sort"
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
)

func (*Daemon) NetworkGetConnections() ([]byte, *dbus.Error) {
	list, err := getConnectionList(nmConnDir)
	if err != nil {
		return nil, dbusutil.ToError(err)
	}
	var info = NetworkData{
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
	kf.RemoveKey(kfSectionConnection, kfKeyInterfaceName)
	kf.RemoveKey(kfSectionWIFI, kfKeyMac)
	kf.RemoveKey(kfSectionWIFI, kfKeyMacBlacklist)
	kf.RemoveKey(kfSectionWIFI, kfKeySeenBSSID)

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
