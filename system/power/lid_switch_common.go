/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

type InputEvent struct {
	Time  syscall.Timeval // time in seconds since epoch at which event occurred
	Type  uint16          // event type - one of ecodes.EV_*
	Code  uint16          // event code related to the event type
	Value int32           // event value related to the event type
}

// Get a useful description for an input event. Example:
//   event at 1347905437.435795, code 01, type 02, val 02
func (ev *InputEvent) String() string {
	return fmt.Sprintf("event at %d.%d, code %02d, type %02d, val %02d",
		ev.Time.Sec, ev.Time.Usec, ev.Code, ev.Type, ev.Value)
}

var eventSize = int(unsafe.Sizeof(InputEvent{}))

func (m *Manager) findLidSwitch() string {
	devices := m.gudevClient.QueryBySubsystem("input")
	var deviceFile string
	for _, device := range devices {
		name := device.GetName()
		if !strings.HasPrefix(name, "event") {
			continue
		}
		logger.Debug("name:", name)
		logger.Debug("sysfsPath:", device.GetSysfsPath())
		devName := device.GetSysfsAttr("device/name")
		logger.Debugf("dev name: %q", devName)
		if strings.Contains(strings.ToLower(devName), "lid switch") {
			deviceFile = device.GetDeviceFile()
			if deviceFile != "" {
				return deviceFile
			}
		}
	}
	return ""
}

func (m *Manager) initLidSwitchCommon() {
	devFile := m.findLidSwitch()
	if devFile == "" {
		logger.Info("Not found lid switch")
		return
	}
	logger.Debugf("find dev file %q", devFile)
	// open
	f, err := os.Open(devFile)
	if err != nil {
		logger.Debug("err:", err)
		return
	}
	m.HasLidSwitch = true

	go func() {
		for {
			events, err := readLidSwitchEvent(f)
			if err != nil {
				logger.Warning(err)
				continue
			}
			for _, ev := range events {
				logger.Debugf("%v", &ev)
				if ev.Type == EV_SW && ev.Code == SW_LID {
					logger.Debugf("lid switch event value: %v", ev.Value)
					var closed bool
					switch ev.Value {
					case 1:
						closed = true
					case 0:
						closed = false
					default:
						logger.Warningf("unknown lid switch event value %v", ev.Value)
						continue
					}
					m.handleLidSwitchEvent(closed)
				}
			}
		}
	}()
}

const (
	EV_SW  = 0x05
	SW_LID = 0x00
)

func readLidSwitchEvent(f *os.File) ([]InputEvent, error) {
	// read
	count := 16

	events := make([]InputEvent, count)
	buffer := make([]byte, eventSize*count)

	_, err := f.Read(buffer)
	if err != nil {
		logger.Debug("f read err:", err)
		return nil, err
	}

	b := bytes.NewBuffer(buffer)
	err = binary.Read(b, binary.LittleEndian, &events)
	if err != nil {
		logger.Debug("binary read err:", err)
		return nil, err
	}

	// remove trailing structures
	for i := range events {
		logger.Debug("i", i)
		if events[i].Time.Sec == 0 {
			events = append(events[:i])
			break
		}
	}
	return events, nil
}
