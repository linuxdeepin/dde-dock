package launcher

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	pinyin "dbus/com/deepin/daemon/search"
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

	keys := []string{key}
	var tkeys []string
	if tree != nil && treeId != "" {
		var err error
		tkeys, err = tree.SearchKeys(key, treeId)
		if err != nil {
			logger.Warning("Search Keys failed:", err)
		}
		logger.Debug("get tree searchKeys:", tkeys)
	}

	for _, v := range tkeys {
		if v != key {
			keys = append(keys, v)
		}
	}

	logger.Debug("searchKeys:", keys)
	done := make(chan bool, 1)
	for _, k := range keys {
		escapedKey := regexp.QuoteMeta(k)
		for _, fn := range searchFuncs {
			go fn(escapedKey, resChan, done)
		}
	}

	for _ = range keys {
		for _ = range searchFuncs {
			select {
			case <-done:
				// logger.Info("done")
			case <-time.After(1 * time.Second):
				logger.Info("wait search result time out")
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
		// logger.Info(itemTable[v.Id].Name, v.Score)
		ids = append(ids, v.Id)
	}
	return ids
}

// 2. add a weight for frequency.
func searchInstalled(key string, res chan<- SearchResult, end chan<- bool) {
	logger.Debug("SearchKey:", key)
	keyMatcher, err := regexp.Compile(fmt.Sprintf("(?i)(%s)", key))
	if err != nil {
		logger.Warning("get key matcher failed:", err)
	}
	matchers := getMatchers(key) // just use these to name.
	for id, v := range itemTable {
		var score uint32 = 0

		logger.Debug("search", v.Name)
		for matcher, s := range matchers {
			if matcher.MatchString(v.Name) {
				logger.Debug("\tName:", v.Name, "match", matcher)
				score += s
			}
		}
		if v.enName != v.Name {
			for matcher, s := range matchers {
				if matcher.MatchString(v.enName) {
					logger.Debug("\tEnName:", v.enName, "match", matcher)
					score += s
				}
			}
		}

		for _, keyword := range v.xinfo.keywords {
			if keyMatcher.MatchString(keyword) {
				logger.Debug("\tKeyword:", keyword, "match", keyMatcher)
				score += VERY_GOOD
			}
		}
		if keyMatcher.MatchString(v.Path) {
			logger.Debug("\tPath:", v.Path, "match", keyMatcher)
			score += AVERAGE
		}
		if keyMatcher.MatchString(v.xinfo.exec) {
			logger.Debug("\tExec:", v.xinfo.exec, "match", keyMatcher)
			score += GOOD
		}
		if keyMatcher.MatchString(v.xinfo.genericName) {
			logger.Debug("\tGenericName:", v.xinfo.genericName, "match", keyMatcher)
			score += BELOW_AVERAGE
		}
		if keyMatcher.MatchString(v.xinfo.description) {
			logger.Debug("\tDescription:", v.xinfo.description, "match", keyMatcher)
			score += POOR
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
