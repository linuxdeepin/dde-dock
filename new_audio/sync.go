package main

import "unsafe"
import "sync"

/*
#include "dde-pulse.h"
*/
import "C"
import "fmt"

type paInfo struct {
	core unsafe.Pointer
	Type int
}

type cookie struct {
	id   int64
	data chan *paInfo
}

func NewCookie(id int64) *cookie {
	return &cookie{int64(id), make(chan *paInfo, 1024)}
}
func (c *cookie) Reply() *paInfo {
	defer deleteCookie(c.id)
	return <-c.data
}
func (c *cookie) ReplyList() []*paInfo {
	defer deleteCookie(c.id)
	var infos []*paInfo
	for info := range c.data {
		infos = append(infos, info)
	}
	return infos
}
func (c *cookie) Free() {
}
func (c *cookie) Feed(infoType int, info unsafe.Pointer) {
	c.data <- &paInfo{info, infoType}
}
func (c *cookie) EndOfList() {
	close(c.data)
}

//export receive_some_info
func receive_some_info(cookie int64, infoType int, info unsafe.Pointer, end bool) {
	c := fetchCookie(cookie)
	if end {
		c.EndOfList()
	} else {
		c.Feed(infoType, info)
	}
}

var newCookie, fetchCookie, deleteCookie = func() (func() *cookie,
	func(int64) *cookie,
	func(int64)) {

	cookies := make(map[int64]*cookie)
	id := int64(0)
	var locker sync.Mutex
	return func() *cookie {
			locker.Lock()
			id++
			locker.Unlock()

			c := NewCookie(id)
			cookies[c.id] = c
			return c
		}, func(i int64) *cookie {
			locker.Lock()
			c := cookies[i]
			locker.Unlock()

			return c
		}, func(i int64) {
			delete(cookies, i)
		}
}()

var ctx *C.pa_context

func init() {
	ml := C.pa_mainloop_new()
	ctx = C.pa_init(ml)
	go C.pa_mainloop_run(ml, nil)
}

//export go_handle_changed
func go_handle_changed(facility int, event_type int, idx uint32) {
	switch facility {
	case C.PA_SUBSCRIPTION_EVENT_CARD:
		if event_type == C.PA_SUBSCRIPTION_EVENT_NEW {
			fmt.Printf("DEBUG card %d new\n", idx)
		} else if event_type == C.PA_SUBSCRIPTION_EVENT_CHANGE {
			fmt.Printf("DEBUG card %d state changed\n", idx)
		} else if event_type == C.PA_SUBSCRIPTION_EVENT_REMOVE {
			fmt.Printf("DEBUG card %d removed\n", idx)
		}
	case C.PA_SUBSCRIPTION_EVENT_SINK:
		if event_type == C.PA_SUBSCRIPTION_EVENT_NEW {
			fmt.Printf("DEBUG sink %d new\n", idx)
		} else if event_type == C.PA_SUBSCRIPTION_EVENT_CHANGE {
			fmt.Printf("DEBUG sink %d state changed\n", idx)
		} else if event_type == C.PA_SUBSCRIPTION_EVENT_REMOVE {
			fmt.Printf("DEBUG sink %d removed\n", idx)
		}
	case C.PA_SUBSCRIPTION_EVENT_SOURCE:
		if event_type == C.PA_SUBSCRIPTION_EVENT_NEW {
			fmt.Printf("DEBUG source %d new\n", idx)
		} else if event_type == C.PA_SUBSCRIPTION_EVENT_CHANGE {
			fmt.Printf("DEBUG source %d changed\n", idx)
		} else if event_type == C.PA_SUBSCRIPTION_EVENT_REMOVE {
			fmt.Printf("DEBUG source %d removed\n", idx)
		}
	case C.PA_SUBSCRIPTION_EVENT_SINK_INPUT:
		if event_type == C.PA_SUBSCRIPTION_EVENT_NEW {
			fmt.Printf("DEBUG sink input %d new\n", idx)
		} else if event_type == C.PA_SUBSCRIPTION_EVENT_CHANGE {
			fmt.Printf("DEBUG sink input %d changed\n", idx)
		} else if event_type == C.PA_SUBSCRIPTION_EVENT_REMOVE {
			fmt.Printf("DEBUG sink input %d removed\n", idx)
		}
	case C.PA_SUBSCRIPTION_EVENT_SOURCE_OUTPUT:
		if event_type == C.PA_SUBSCRIPTION_EVENT_NEW {
			fmt.Printf("DEBUG source output %d new\n", idx)
		} else if event_type == C.PA_SUBSCRIPTION_EVENT_CHANGE {
			fmt.Printf("DEBUG source output %d changed\n", idx)
		} else if event_type == C.PA_SUBSCRIPTION_EVENT_REMOVE {
			fmt.Printf("DEBUG source output %d removed\n", idx)
		}
	case C.PA_SUBSCRIPTION_EVENT_SAMPLE_CACHE:
	}
}
