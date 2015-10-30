package search

import (
	"fmt"
	"regexp"
	"sync"

	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
)

// SearchInstalledItemTransaction is a command object for searching installed items.
type SearchInstalledItemTransaction struct {
	maxGoroutineNum int

	keyMatcher   *regexp.Regexp
	nameMatchers map[*regexp.Regexp]uint32

	resChan    chan<- Result
	cancelChan chan struct{}
	cancelled  bool
}

// NewSearchInstalledItemTransaction creates a new SearchInstalledItemTransaction object.
func NewSearchInstalledItemTransaction(res chan<- Result, cancelChan chan struct{}, maxGoroutineNum int) (*SearchInstalledItemTransaction, error) {
	if res == nil {
		return nil, ErrorSearchNullChannel
	}

	if maxGoroutineNum <= 0 {
		maxGoroutineNum = DefaultGoroutineNum
	}

	return &SearchInstalledItemTransaction{
		maxGoroutineNum: maxGoroutineNum,
		keyMatcher:      nil,
		nameMatchers:    nil,
		resChan:         res,
		cancelChan:      cancelChan,
		cancelled:       false,
	}, nil
}

// Cancel cancels this transaction.
func (s *SearchInstalledItemTransaction) Cancel() {
	if !s.cancelled {
		close(s.cancelChan)
		s.cancelled = true
	}
}

func (s *SearchInstalledItemTransaction) initKeyMatchers(key string) {
	s.keyMatcher, _ = regexp.Compile(fmt.Sprintf("(?i)(%s)", key))
	s.nameMatchers = getMatchers(key)
}

func (s *SearchInstalledItemTransaction) calcScore(data ItemInfo) (score uint32) {
	for matcher, s := range s.nameMatchers {
		if matcher.MatchString(data.Name()) {
			score += s
		}
	}
	if data.LocaleName() != data.Name() {
		for matcher, s := range s.nameMatchers {
			if matcher.MatchString(data.Name()) {
				score += s
			}
		}
	}

	if s.keyMatcher == nil {
		return
	}

	for _, keyword := range data.Keywords() {
		if s.keyMatcher.MatchString(keyword) {
			score += VeryGood
		}
	}

	if s.keyMatcher.MatchString(data.Path()) {
		score += Average
	}

	if s.keyMatcher.MatchString(data.ExecCmd()) {
		score += Good
	}

	if s.keyMatcher.MatchString(data.GenericName()) {
		score += BelowAverage
	}

	if s.keyMatcher.MatchString(data.Description()) {
		score += Poor
	}

	return
}

func (s *SearchInstalledItemTransaction) isCancelled() bool {
	if s.cancelled {
		return true
	}

	select {
	case <-s.cancelChan:
		s.cancelled = true
		return true
	default:
		return false
	}
}

// ScoreItem scores item.
func (s *SearchInstalledItemTransaction) ScoreItem(dataSetChan <-chan ItemInfo) {
	if s.isCancelled() {
		return
	}

	for data := range dataSetChan {
		score := s.calcScore(data)
		if score == 0 {
			continue
		}

		if s.isCancelled() {
			return
		}

		select {
		case s.resChan <- Result{
			ID:    data.ID(),
			Name:  data.Name(),
			Score: score,
		}:
		}
	}
}

// Search executes transaction and returns searching results.
func (s *SearchInstalledItemTransaction) Search(key string, dataSet []ItemInfo) {
	s.initKeyMatchers(key)
	dataSetChan := make(chan ItemInfo)
	go func() {
		defer close(dataSetChan)
		if s.isCancelled() {
			return
		}

		for _, data := range dataSet {
			if s.isCancelled() {
				return
			}

			select {
			case dataSetChan <- data:
			}
		}
	}()
	var wg sync.WaitGroup
	wg.Add(s.maxGoroutineNum)
	for i := 0; i < s.maxGoroutineNum; i++ {
		go func(i int) {
			s.ScoreItem(dataSetChan)
			wg.Done()
		}(i)
	}

	wg.Wait()
}
