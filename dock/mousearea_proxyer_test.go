/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import (
	"errors"
	C "launchpad.net/gocheck"
	"time"
)

type MockXMouseArea struct {
	motionIntoCB func(int32, int32, string)
	motionOutCB  func(int32, int32, string)
}

func NewMockXMouseArea() (*MockXMouseArea, error) {
	return &MockXMouseArea{}, nil
}

func (m *MockXMouseArea) emitMotionInto(x, y int32, id string) {
	if m.motionIntoCB != nil {
		m.motionIntoCB(x, y, id)
	}
}

func (m *MockXMouseArea) emitMotionOut(x, y int32, id string) {
	if m.motionOutCB != nil {
		m.motionOutCB(x, y, id)
	}
}

func (m *MockXMouseArea) ConnectCursorInto(cb func(int32, int32, string)) func() {
	m.motionIntoCB = cb
	return func() {}
}

func (m *MockXMouseArea) ConnectCursorOut(cb func(int32, int32, string)) func() {
	m.motionOutCB = cb
	return func() {}
}

func (m *MockXMouseArea) UnregisterArea(id string) error {
	return nil
}

func (m *MockXMouseArea) RegisterAreas(areas interface{}, eventMask int32) (string, error) {
	return "0", nil
}

func (m *MockXMouseArea) RegisterFullScreen() (string, error) {
	return "0", nil
}

type XMouseAreaTestSuite struct{}

var _ = C.Suite(&XMouseAreaTestSuite{})

var mockXMouseArea, err = NewMockXMouseArea()
var xmouseArea, _ = NewXMouseAreaProxyer(mockXMouseArea, err)

func (s *XMouseAreaTestSuite) TestNewXMouseArea(c *C.C) {
	_, err := NewXMouseAreaProxyer(NewMockXMouseArea())
	c.Check(err, C.IsNil)
	_, err = NewXMouseAreaProxyer(&MockXMouseArea{}, errors.New("create MockXMouseArea failed"))
	c.Check(err, C.NotNil)
}

func (s *XMouseAreaTestSuite) Test_unregister(c *C.C) {
	xmouseArea.unregister()
	c.Check(xmouseArea.idValid, C.Equals, false)

	xmouseArea.RegisterFullScreen()
	xmouseArea.unregister()
	c.Check(xmouseArea.idValid, C.Equals, false)
}

func (s *XMouseAreaTestSuite) Test_connectMotionInto(c *C.C) {
	ch := make(chan struct{})
	xmouseArea.connectMotionInto(func(_, _ int32, id string) {
		close(ch)
	})
	xmouseArea.RegisterFullScreen()
	mockXMouseArea.emitMotionInto(0, 0, "0")
	select {
	case <-ch:
	case <-time.After(time.Second):
		c.FailNow()
	}
}

func (s *XMouseAreaTestSuite) Test_connectMotionOut(c *C.C) {
	ch := make(chan struct{})
	xmouseArea.connectMotionOut(func(_, _ int32, id string) {
		close(ch)
	})
	xmouseArea.RegisterFullScreen()
	mockXMouseArea.emitMotionOut(0, 0, "0")
	select {
	case <-ch:
	case <-time.After(time.Second):
		c.FailNow()
	}
}
