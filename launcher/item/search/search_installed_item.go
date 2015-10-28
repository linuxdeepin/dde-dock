package search

import (
	"fmt"
	"regexp"
	"sync"

	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
)

type SearchInstalledItemTransaction struct {
	maxGoroutineNum int

	keyMatcher   *regexp.Regexp
	nameMatchers map[*regexp.Regexp]uint32

	resChan    chan<- SearchResult
	cancelChan chan struct{}
	cancelled  bool
}

func NewSearchInstalledItemTransaction(res chan<- SearchResult, cancelChan chan struct{}, maxGoroutineNum int) (*SearchInstalledItemTransaction, error) {
	if res == nil {
		return nil, SearchErrorNullChannel
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

func (s *SearchInstalledItemTransaction) calcScore(data ItemInfoInterface) (score uint32) {
	for matcher, s := range s.nameMatchers {
		if matcher.MatchString(data.Name()) {
			score += s
		}
	}
	if data.EnName() != data.Name() {
		for matcher, s := range s.nameMatchers {
			if matcher.MatchString(data.EnName()) {
				score += s
			}
		}
	}

	if s.keyMatcher == nil {
		return
	}

	for _, keyword := range data.Keywords() {
		if s.keyMatcher.MatchString(keyword) {
			score += VERY_GOOD
		}
	}

	if s.keyMatcher.MatchString(data.Path()) {
		score += AVERAGE
	}

	if s.keyMatcher.MatchString(data.ExecCmd()) {
		score += GOOD
	}

	if s.keyMatcher.MatchString(data.GenericName()) {
		score += BELOW_AVERAGE
	}

	if s.keyMatcher.MatchString(data.Description()) {
		score += POOR
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

func (s *SearchInstalledItemTransaction) ScoreItem(dataSetChan <-chan ItemInfoInterface) {
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
		case s.resChan <- SearchResult{
			Id:    data.Id(),
			Name:  data.Name(),
			Score: score,
		}:
		}
	}
}

func (s *SearchInstalledItemTransaction) Search(key string, dataSet []ItemInfoInterface) {
	s.initKeyMatchers(key)
	dataSetChan := make(chan ItemInfoInterface)
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
