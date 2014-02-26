package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	pinyin "dbus/com/deepin/api/search"
)

var tree *pinyin.Search = nil
var treeId string

type SearchFunc func(key string, res chan<- SearchResult, end chan<- bool)

// func registerSearchFunc(searchFunc SearchFunc) {
// 	searchFuncs = append(searchFuncs, searchFunc)
// }

type SearchResult struct {
	Id    ItemId
	Score uint32
}

type Result struct {
	sync.RWMutex
	Res map[ItemId]SearchResult
}

type ResultList []SearchResult

func (res ResultList) Len() int {
	return len(res)
}

func (res ResultList) Swap(i, j int) {
	res[i], res[j] = res[j], res[i]
}

func (res ResultList) Less(i, j int) bool {
	if res[i].Score > res[j].Score {
		return true
	} else if res[i].Score == res[j].Score {
		return itemTable[res[i].Id].Name < itemTable[res[j].Id].Name
	} else {
		return false
	}
}

// TODO:
// 1. cancellable
func search(key string) []ItemId {
	key = strings.TrimSpace(key)
	res := Result{Res: make(map[ItemId]SearchResult, 0)}
	resChan := make(chan SearchResult)
	go func(r *Result, c <-chan SearchResult) {
		for {
			select {
			case d := <-c:
				r.Lock()
				if _, ok := r.Res[d.Id]; !ok {
					r.Res[d.Id] = d
				} else {
					d.Score = r.Res[d.Id].Score + d.Score
					r.Res[d.Id] = d
				}
				r.Unlock()
			case <-time.After(2 * time.Second):
				return
			}
		}
	}(&res, resChan)

	keys := []string{}
	if tree != nil {
		keys, _ = tree.SearchKeys(key, treeId)
	}

	for _, v := range keys {
		if v != key {
			keys = append(keys, key)
		}
	}

	done := make(chan bool, 1)
	for _, k := range keys {
		for _, fn := range searchFuncs {
			go fn(k, resChan, done)
		}
	}

	for _ = range keys {
		for _ = range searchFuncs {
			select {
			case <-done:
				// fmt.Println("done")
			case <-time.After(1 * time.Second):
				fmt.Println("wait search result time out")
			}
		}
	}

	resList := make(ResultList, 0)
	for _, v := range res.Res {
		resList = append(resList, v)
	}
	sort.Sort(resList)

	ids := make([]ItemId, 0)
	for _, v := range resList {
		// fmt.Println(itemTable[v.Id].Name, v.Score)
		ids = append(ids, v.Id)
	}
	return ids
}

// 2. add a weight for frequency.
func searchInstalled(key string, res chan<- SearchResult, end chan<- bool) {
	matchers := getMatchers(key)
	for id, v := range itemTable {
		var score uint32 = 0
		var weight uint32 = 1

		for matcher, s := range matchers {
			if matcher.MatchString(v.Name) {
				score += s * weight
			}
			for _, keyword := range v.xinfo.keywords {
				if matcher.MatchString(keyword) {
					score += s * weight
				}
			}
			// if matcher.MatchString(v.xinfo.exec) {
			// 	score += s * weight
			// }
			// if matcher.MatchString(v.xinfo.genericName) {
			// 	score += s * weight
			// }
			// if matcher.MatchString(v.xinfo.description) {
			// 	score += s * weight
			// }
		}

		if score > 0 {
			res <- SearchResult{id, score}
		}
	}
	end <- true
}

var searchFuncs = []SearchFunc{
	searchInstalled,
}
