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
	"strings"
)

type Trie struct {
	Key      byte
	Values   []string
	NextNode [_TRIE_CHILD_LEN]*Trie
}

const (
	_TRIE_CHILD_LEN = 26
)

var (
	trieMD5Map map[string]*Trie
	strsMD5Map map[string][]*TrieInfo
)

func getNode(ch byte) *Trie {
	node := new(Trie)
	node.Key = ch
	return node
}

func newTrie() *Trie {
	root := getNode(' ')
	return root
}

func (root *Trie) insertTrieInfo(values []*TrieInfo) {
	for _, v := range values {
		root.insertStringArray(v.Pinyins, v.Key)
	}
}

func (root *Trie) insertStringArray(strs []string, id string) {
	if strs == nil {
		return
	}

	for _, v := range strs {
		root.insertString(v, id)
	}
}

func (root *Trie) insertString(str, id string) {
	if l := len(str); l == 0 {
		return
	}
	low := strings.ToLower(str)
	curNode := root

	for i, _ := range str {
		index := low[i] - 'a'
		if curNode.NextNode[index] == nil {
			curNode.NextNode[index] = getNode(low[i])
		}
		if !isIdExist(id, curNode.NextNode[index].Values) {
			curNode.NextNode[index].Values = append(curNode.NextNode[index].Values, id)
		}
		curNode = curNode.NextNode[index]
	}
}

func (node *Trie) traversalTrie() {
	if node == nil {
		Logger.Info("trie is nil")
		return
	}

	for i := 0; i < _TRIE_CHILD_LEN; i++ {
		if node.NextNode[i] != nil {
			node.NextNode[i].traversalTrie()
			/*Logger.Info(node.NextNode[i].Key)*/
			Logger.Infof("%v", node.NextNode[i].Values)
		}
	}
}

func (root *Trie) searchTrie(keys string) []string {
	if root == nil {
		return nil
	}
	if len(keys) <= 0 {
		return nil
	}

	curNode := root
	low := strings.ToLower(keys)
	for _, i := range low {
		if i >= 'a' && i <= 'z' {
			index := i - 'a'
			if curNode.NextNode[index] == nil {
				return nil
			}
			curNode = curNode.NextNode[index]
		} else {
			return nil
		}
	}

	retArray := curNode.Values
	Logger.Info("ret array:", retArray)
	return retArray
}
