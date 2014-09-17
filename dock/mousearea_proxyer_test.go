package dock

import (
	"errors"
	"testing"
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

var mockXMouseArea, err = NewMockXMouseArea()
var xmouseArea, _ = NewXMouseArea(mockXMouseArea, err)

func TestNewXMouseArea(t *testing.T) {
	_, err := NewXMouseArea(NewMockXMouseArea())
	if err != nil {
		t.FailNow()
	}
	_, err = NewXMouseArea(&MockXMouseArea{}, errors.New("create MockXMouseArea failed"))
	if err == nil {
		t.FailNow()
	}
}

func Test_unregister(t *testing.T) {
	xmouseArea.unregister()
	if xmouseArea.idValid != false {
		t.FailNow()
	}

	xmouseArea.RegisterFullScreen()
	xmouseArea.unregister()
	if xmouseArea.idValid != false {
		t.FailNow()
	}
}

func Test_connectMotionInto(t *testing.T) {
	c := make(chan struct{})
	xmouseArea.connectMotionInto(func(_, _ int32, id string) {
		close(c)
	})
	xmouseArea.RegisterFullScreen()
	mockXMouseArea.emitMotionInto(0, 0, "0")
	select {
	case <-c:
	case <-time.After(time.Second):
		t.FailNow()
	}
}

func Test_connectMotionOut(t *testing.T) {
	c := make(chan struct{})
	xmouseArea.connectMotionOut(func(_, _ int32, id string) {
		close(c)
	})
	xmouseArea.RegisterFullScreen()
	mockXMouseArea.emitMotionOut(0, 0, "0")
	select {
	case <-c:
	case <-time.After(time.Second):
		t.FailNow()
	}
}
