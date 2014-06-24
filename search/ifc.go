/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
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

package search

import (
	"pkg.linuxdeepin.com/lib/dbus"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"strings"
)

const (
	SEARCH_DEST = "com.deepin.daemon.Search"
	SEARCH_PATH = "/com/deepin/daemon/Search"
	SEARCH_IFC  = "com.deepin.daemon.Search"
)

func (s *Search) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		SEARCH_DEST,
		SEARCH_PATH,
		SEARCH_IFC,
	}
}

func (s *Search) NewTrieWithString(values map[string]string, name string) string {
	md5Str, ok := dutils.SumStrMd5(getStringFromArray(values))
	if !ok {
		return ""
	}
	if isMd5Exist(md5Str) {
		return md5Str
	}

	if isNameExist(name) {
		str, _ := s.nameMD5Map[name]
		s.DestroyTrie(str)
	}
	s.nameMD5Map[name] = md5Str

	root := newTrie()
	if values == nil {
		return ""
	}
	go func() {
		infos := getPinyinArray(values)
		s.strsMD5Map[md5Str] = infos
		root.insertTrieInfo(infos)
		s.trieMD5Map[md5Str] = root
	}()
	return md5Str
}

/*
func (s *Search) TraversalTrie(str string) {
	root := s.trieMD5Map[str]
	root.traversalTrie()
}
*/

func (s *Search) SearchKeys(keys string, str string) []string {
	root, ok := s.trieMD5Map[str]
	if !ok {
		return nil
	}
	keys = strings.ToLower(keys)
	rets := root.searchTrie(keys)
	Logger.Debug("trie rets:", rets)
	tmp := searchKeyFromString(keys, str)
	Logger.Debug("array rets:", tmp)
	for _, v := range tmp {
		if !isIdExist(v, rets) {
			rets = append(rets, v)
		}
	}

	return rets
}

func (s *Search) SearchKeysByFirstLetter(keys string, md5Str string) []string {
	root, ok := s.trieMD5Map[md5Str]
	if !ok {
		return nil
	}

	keys = strings.ToLower(keys)
	rets := root.searchTrie(keys)

	return rets
}

func (s *Search) DestroyTrie(md5Str string) {
	/*
		root, ok := s.trieMD5Map[md5Str]
		if !ok {
			return
		}
	*/
	delete(s.trieMD5Map, md5Str)
	delete(s.strsMD5Map, md5Str)
}
