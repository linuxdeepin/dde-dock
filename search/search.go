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

package search

import (
	dpinyin "pkg.linuxdeepin.com/lib/pinyin"
	"strings"
)

type Search struct {
	trieMD5Map map[string]*Trie
	strsMD5Map map[string][]*TrieInfo
	nameMD5Map map[string]string
}

type TrieInfo struct {
	Pinyins []string
	Key     string
	Value   string
}

var _search *Search

func GetManager() *Search {
	if _search == nil {
		_search = newSearch()
	}

	return _search
}

func newSearch() *Search {
	s := &Search{}

	s.trieMD5Map = make(map[string]*Trie)
	s.strsMD5Map = make(map[string][]*TrieInfo)
	s.nameMD5Map = make(map[string]string)

	return s
}

func getStringFromArray(strs map[string]string) string {
	str := ""

	for _, v := range strs {
		str += v
	}

	return str
}

func getPinyinArray(strs map[string]string) []*TrieInfo {
	rets := []*TrieInfo{}
	for k, v := range strs {
		array := dpinyin.HansToPinyin(v)
		v = strings.ToLower(v)
		tmp := &TrieInfo{Pinyins: array, Key: k, Value: v}
		rets = append(rets, tmp)
	}

	return rets
}

func searchKeyFromString(key, md5Str string) []string {
	rets := []string{}

	infos := GetManager().strsMD5Map[md5Str]
	for _, v := range infos {
		if strings.Contains(v.Value, key) {
			rets = append(rets, v.Key)
		}
	}

	return rets
}

func isIdExist(id string, list []string) bool {
	for _, v := range list {
		if v == id {
			return true
		}
	}

	return false
}

func isMd5Exist(md5Str string) bool {
	_, ok := GetManager().strsMD5Map[md5Str]
	if ok {
		return true
	}

	return false
}

func isNameExist(name string) bool {
	_, ok := GetManager().nameMD5Map[name]
	if ok {
		return true
	}

	return false
}
