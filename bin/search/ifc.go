/*
 * Copyright (C) 2013 ~ 2018 Deepin Technology Co., Ltd.
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

package main

import (
	"path"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/pinyin"
	dutils "pkg.deepin.io/lib/utils"
)

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}

func (m *Manager) NewSearchWithStrList(list []string) (string, bool, *dbus.Error) {
	m.service.DelayAutoQuit()
	var datas []dataInfo
	strs := ""

	for _, v := range list {
		strs += v + "+"
		pyList := pinyin.HansToPinyin(v)
		if len(pyList) == 1 && pyList[0] == v {
			info := dataInfo{v, v}
			datas = append(datas, info)
			continue
		}
		for _, k := range pyList {
			info := dataInfo{k, v}
			datas = append(datas, info)
		}
		info := dataInfo{v, v}
		datas = append(datas, info)
	}

	md5Str, ok := dutils.SumStrMd5(strs)
	if !ok {
		logger.Warning("Sum MD5 Failed")
		return "", false, nil
	}

	cachePath, ok1 := getCachePath()
	if !ok1 {
		logger.Warning("Get Cache Path Failed")
		return "", false, nil
	}

	filename := path.Join(cachePath, md5Str)
	m.writeStart = true
	m.writeEnd = make(chan bool)
	go func() {
		writeDatasToFile(&datas, filename)
		m.writeEnd <- true
		m.writeStart = false
	}()

	return md5Str, true, nil
}

func (m *Manager) NewSearchWithStrDict(dict map[string]string) (string, bool, *dbus.Error) {
	m.service.DelayAutoQuit()

	var datas []dataInfo
	strs := ""

	for k, v := range dict {
		strs += k + "+"
		pyList := pinyin.HansToPinyin(v)
		if len(pyList) == 1 && pyList[0] == v {
			info := dataInfo{v, k}
			datas = append(datas, info)
			continue
		}

		for _, l := range pyList {
			info := dataInfo{l, k}
			datas = append(datas, info)
		}
		info := dataInfo{v, k}
		datas = append(datas, info)
	}

	md5Str, ok := dutils.SumStrMd5(strs)
	if !ok {
		logger.Warning("Sum MD5 Failed")
		return "", false, nil
	}

	cachePath, ok1 := getCachePath()
	if !ok1 {
		logger.Warning("Get Cache Path Failed")
		return "", false, nil
	}

	filename := path.Join(cachePath, md5Str)
	m.writeStart = true
	m.writeEnd = make(chan bool)
	go func() {
		writeDatasToFile(&datas, filename)
		m.writeEnd <- true
		m.writeStart = false
	}()

	return md5Str, true, nil
}

func (m *Manager) SearchString(str, md5 string) ([]string, *dbus.Error) {
	m.service.DelayAutoQuit()

	var list []string
	if len(str) < 1 || len(md5) < 1 {
		return list, nil
	}

	list = searchString(str, md5)
	var tmp []string
	for _, v := range list {
		if !strIsInList(v, tmp) {
			tmp = append(tmp, v)
		}
	}

	return tmp, nil
}

func (m *Manager) SearchStartWithString(str, md5 string) ([]string, *dbus.Error) {
	m.service.DelayAutoQuit()

	var list []string
	if len(str) < 1 || len(md5) < 1 {
		return list, nil
	}

	list = searchStartWithString(str, md5)
	var tmp []string
	for _, v := range list {
		if !strIsInList(v, tmp) {
			tmp = append(tmp, v)
		}
	}

	return tmp, nil
}

func strIsInList(str string, list []string) bool {
	for _, l := range list {
		if str == l {
			return true
		}
	}

	return false
}
