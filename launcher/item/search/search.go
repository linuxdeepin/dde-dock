package search

import (
	// "fmt"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"regexp"
	"strings"
	"sync"
)

const (
	DefaultGoroutineNum = 20
)

type SearchResult struct {
	Id    ItemId
	Name  string
	Score uint32
}

type SearchTransaction struct {
	maxGoroutineNum int
	pinyinObj       PinYinInterface
	result          chan<- SearchResult
	cancelChan      chan struct{}
	cancelled       bool
}

func NewSearchTransaction(pinyinObj PinYinInterface, result chan<- SearchResult, cancelChan chan struct{}, maxGoroutineNum int) (*SearchTransaction, error) {
	if result == nil {
		return nil, SearchErrorNullChannel
	}
	if maxGoroutineNum <= 0 {
		maxGoroutineNum = DefaultGoroutineNum
	}
	return &SearchTransaction{
		maxGoroutineNum: maxGoroutineNum,
		pinyinObj:       pinyinObj,
		result:          result,
		cancelChan:      cancelChan,
		cancelled:       false,
	}, nil
}

func (s *SearchTransaction) Cancel() {
	if !s.cancelled {
		close(s.cancelChan)
		s.cancelled = true
	}
}

func (s *SearchTransaction) Search(key string, dataSet []ItemInfoInterface) {
	trimedKey := strings.TrimSpace(key)
	escapedKey := regexp.QuoteMeta(trimedKey)

	keys := make(chan string)
	go func() {
		defer close(keys)
		if s.pinyinObj != nil && s.pinyinObj.IsValid() {
			pinyins, _ := s.pinyinObj.Search(escapedKey)
			for _, pinyin := range pinyins {
				keys <- pinyin
			}
		}
		keys <- escapedKey
	}()

	const MaxKeyGoroutineNum = 5
	var wg sync.WaitGroup
	wg.Add(MaxKeyGoroutineNum)
	for i := 0; i < MaxKeyGoroutineNum; i++ {
		go func() {
			for key := range keys {
				transaction, _ := NewSearchInstalledItemTransaction(s.result, s.cancelChan, s.maxGoroutineNum)
				transaction.Search(key, dataSet)
				select {
				case <-s.cancelChan:
					break
				default:
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
