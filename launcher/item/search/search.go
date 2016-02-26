/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package search

import (
	"regexp"
	"strings"
	"sync"

	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
)

// default values.
const (
	DefaultGoroutineNum = 20
)

type FreqGetter interface {
	GetFrequency(string) uint64
}

// Result stores items information for searching.
type Result struct {
	ID    ItemID
	Name  string
	Score uint32
	Freq  uint64
}

// Transaction is a command object used for search.
type Transaction struct {
	maxGoroutineNum int
	pinyinObj       PinYin
	freqGetter      FreqGetter
	result          chan<- Result
	cancelChan      chan struct{}
	cancelled       bool
}

// NewTransaction creates a new Transaction object.
func NewTransaction(pinyinObj PinYin, result chan<- Result, cancelChan chan struct{}, maxGoroutineNum int) (*Transaction, error) {
	if result == nil {
		return nil, ErrorSearchNullChannel
	}
	if maxGoroutineNum <= 0 {
		maxGoroutineNum = DefaultGoroutineNum
	}
	return &Transaction{
		maxGoroutineNum: maxGoroutineNum,
		pinyinObj:       pinyinObj,
		result:          result,
		cancelChan:      cancelChan,
		cancelled:       false,
	}, nil
}

// Cancel cancels this transaction.
func (s *Transaction) Cancel() {
	if !s.cancelled {
		close(s.cancelChan)
		s.cancelled = true
	}
}

func (s *Transaction) SetFreqGetter(freqGetter FreqGetter) *Transaction {
	s.freqGetter = freqGetter
	return s
}

func (s *Transaction) getFreq(id ItemID) uint64 {
	if s.freqGetter == nil {
		return 0
	}
	return s.freqGetter.GetFrequency(string(id))
}

// Search executes this transaction and returns the searching result.
func (s *Transaction) Search(key string, dataSet []ItemInfo) {
	trimedKey := strings.TrimSpace(key)
	escapedKey := regexp.QuoteMeta(trimedKey)

	enablePinYinSearch := s.pinyinObj != nil && s.pinyinObj.IsValid()
	keys := make(chan string)
	go func() {
		defer close(keys)
		if enablePinYinSearch {
			pinyins, _ := s.pinyinObj.Search(escapedKey)
			for _, pinyin := range pinyins {
				keys <- pinyin
			}
		}
	}()

	transaction, _ := NewSearchInstalledItemTransaction(s.result, s.cancelChan, s.maxGoroutineNum)
	transaction.SetFreqGetter(s.freqGetter)
	transaction.Search(escapedKey, dataSet)

	if !enablePinYinSearch {
		return
	}

	const MaxKeyGoroutineNum = 5
	var wg sync.WaitGroup
	wg.Add(MaxKeyGoroutineNum)
	for i := 0; i < MaxKeyGoroutineNum; i++ {
		go func(i int) {
			for key := range keys {
				// pinyin is the LocaleDisplayName.
				for _, info := range dataSet {
					if info.LocaleName() == key {
						s.result <- Result{
							ID:    info.ID(),
							Name:  info.LocaleName(),
							Score: Excellent,
							Freq:  s.getFreq(info.ID()),
						}
					}
				}
				select {
				case <-s.cancelChan:
					break
				default:
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
