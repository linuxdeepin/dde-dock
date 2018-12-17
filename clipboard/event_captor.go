package clipboard

import (
	"container/list"
	"errors"
	"sync"
	"time"

	"github.com/linuxdeepin/go-x11-client"
)

type eventCaptorHook struct {
	time time.Time
	cond func(ev interface{}) bool
	ch   chan interface{}
}

type eventCaptor struct {
	mu       sync.Mutex
	hooks    *list.List
	interval time.Duration
	timer    *time.Timer
	quit     chan struct{}
}

func newEventCaptor() *eventCaptor {
	interval := 5 * time.Second
	timer := time.NewTimer(interval)
	timer.Stop()
	l := &eventCaptor{
		hooks:    list.New(),
		timer:    timer,
		interval: interval,
	}
	go l.loopCheck()
	return l
}

func (ec *eventCaptor) destroy() {
	close(ec.quit)
}

func (ec *eventCaptor) loopCheck() {
	const hookExpires = 11 * time.Second
	for {
		select {
		case <-ec.timer.C:
			ec.mu.Lock()

			var expiredElems []*list.Element
			for e := ec.hooks.Front(); e != nil; e = e.Next() {
				hook := e.Value.(*eventCaptorHook)
				if time.Since(hook.time) > hookExpires {
					expiredElems = append(expiredElems, e)
				}
			}
			for _, e := range expiredElems {
				hook := e.Value.(*eventCaptorHook)
				close(hook.ch)
				ec.hooks.Remove(e)
			}

			if ec.hooks.Len() > 0 {
				ec.timer.Reset(ec.interval)
			}
			ec.mu.Unlock()
		case <-ec.quit:
			return
		}
	}
}

func (ec *eventCaptor) captureEvent(begin func() error, cond func(interface{}) bool) (interface{},
	error) {
	ch := make(chan interface{})
	hook := &eventCaptorHook{
		time: time.Now(),
		cond: cond,
		ch:   ch,
	}

	ec.mu.Lock()
	ec.hooks.PushBack(hook)
	ec.timer.Reset(ec.interval)
	ec.mu.Unlock()
	err := begin()
	if err != nil {
		return nil, err
	}
	event, ok := <-ch
	if !ok {
		return nil, errTimeout
	}
	return event, nil
}

var errTimeout = errors.New("wait event timed out")

func (ec *eventCaptor) capturePropertyNotifyEvent(begin func() error,
	cond func(*x.PropertyNotifyEvent) bool) (*x.PropertyNotifyEvent, error) {
	event, err := ec.captureEvent(begin, func(ev interface{}) bool {
		event, ok := ev.(*x.PropertyNotifyEvent)
		if !ok {
			return false
		}
		return cond(event)
	})
	if err != nil {
		return nil, err
	}
	return event.(*x.PropertyNotifyEvent), nil
}

func (ec *eventCaptor) captureSelectionNotifyEvent(begin func() error,
	cond func(event *x.SelectionNotifyEvent) bool) (*x.SelectionNotifyEvent, error) {
	event, err := ec.captureEvent(begin, func(ev interface{}) bool {
		event, ok := ev.(*x.SelectionNotifyEvent)
		if !ok {
			return false
		}
		return cond(event)
	})
	if err != nil {
		return nil, err
	}
	return event.(*x.SelectionNotifyEvent), nil
}

func (ec *eventCaptor) handleEvent(event interface{}) bool {
	var matchedElem *list.Element
	ec.mu.Lock()
	for e := ec.hooks.Front(); e != nil; e = e.Next() {
		hook := e.Value.(*eventCaptorHook)
		if hook.cond != nil {
			if hook.cond(event) {
				matchedElem = e
				break
			}
		}
	}
	ec.mu.Unlock()

	if matchedElem == nil {
		return false
	}

	hook := matchedElem.Value.(*eventCaptorHook)
	hook.ch <- event
	close(hook.ch)
	ec.mu.Lock()
	ec.hooks.Remove(matchedElem)
	ec.mu.Unlock()
	return true
}
