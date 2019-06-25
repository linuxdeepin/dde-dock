package clipboard

import (
	"errors"
	"testing"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/stretchr/testify/assert"
	"pkg.deepin.io/dde/daemon/clipboard/mocks"
)

func initAtomsForTest() {
	const base = 100
	atomClipboardManager = base + 1
	atomClipboard = base + 2
	atomSaveTargets = base + 3
	atomTargets = base + 4
	atomMultiple = base + 5
	atomDelete = base + 6
	atomInsertProperty = base + 7
	atomInsertSelection = base + 8
	atomAtomPair = base + 9
	atomIncr = base + 10
	atomTimestamp = base + 11
	atomTimestampProp = base + 12
	atomNull = base + 13
}

func TestManager_finishSelectionRequest(t *testing.T) {
	ts := x.Timestamp(11)
	owner := x.Window(1)
	reqWin := x.Window(2)
	selection := x.Atom(3)
	target := x.Atom(4)
	prop := x.Atom(5)

	ev := &x.SelectionRequestEvent{
		Time:      ts,
		Owner:     owner,
		Requestor: reqWin,
		Selection: selection,
		Target:    target,
		Property:  prop,
	}

	m := &Manager{}
	m.window = owner
	xc := &mocks.XClient{}
	m.xc = xc
	xc.On("SendEventE", false, reqWin, uint32(x.EventMaskNoEvent), &x.SelectionNotifyEvent{
		Time:      ts,
		Requestor: reqWin,
		Selection: selection,
		Target:    target,
		Property:  prop,
	}).Return(nil).Once()

	m.finishSelectionRequest(ev, true)

	xc.On("SendEventE", false, reqWin, uint32(x.EventMaskNoEvent), &x.SelectionNotifyEvent{
		Time:      ts,
		Requestor: reqWin,
		Selection: selection,
		Target:    target,
		Property:  x.AtomNone,
	}).Return(nil).Once()
	m.finishSelectionRequest(ev, false)

	xc.AssertExpectations(t)
}

func TestManager_getProperty(t *testing.T) {
	win := x.Window(1)
	prop := x.Atom(2)
	m := &Manager{}
	xc := &mocks.XClient{}
	m.xc = xc

	w := x.NewWriter()
	w.Write4b(1)
	w.Write4b(2)
	w.Write4b(3)
	w.Write4b(4)

	var err1 = errors.New("E1")

	tests := []struct {
		desc            string
		propReply0      *x.GetPropertyReply
		propReply1      *x.GetPropertyReply
		err0            error
		err1            error
		resultPropReply *x.GetPropertyReply
		resultErr       error
		longLength      uint32
	}{
		{
			desc: "normal: data type: []uint32, len: 4",
			propReply0: &x.GetPropertyReply{
				Format:     32,
				Type:       x.AtomCardinal,
				BytesAfter: 16,
			},
			propReply1: &x.GetPropertyReply{
				Format:   32,
				Type:     x.AtomCardinal,
				Value:    w.Bytes(),
				ValueLen: 4,
			},
			resultPropReply: &x.GetPropertyReply{
				Format:   32,
				Type:     x.AtomCardinal,
				Value:    w.Bytes(),
				ValueLen: 4,
			},
			longLength: 4,
		},
		{
			desc: "normal: data type: string, len: 5",
			propReply0: &x.GetPropertyReply{
				Format:     8,
				Type:       x.AtomString,
				BytesAfter: 5,
			},
			propReply1: &x.GetPropertyReply{
				Format:   8,
				Type:     x.AtomString,
				Value:    []byte("abcde"),
				ValueLen: 5,
			},
			resultPropReply: &x.GetPropertyReply{
				Format:   8,
				Type:     x.AtomString,
				Value:    []byte("abcde"),
				ValueLen: 5,
			},
			longLength: 2,
		},
		{
			desc:      "the first call to get property failed",
			err0:      err1,
			resultErr: err1,
		},
		{
			desc: "the second call to get property failed",
			propReply0: &x.GetPropertyReply{
				Format:     8,
				Type:       x.AtomString,
				BytesAfter: 5,
			},
			propReply1: &x.GetPropertyReply{
				Format:   8,
				Type:     x.AtomString,
				Value:    []byte("hello"),
				ValueLen: 5,
			},
			longLength: 2,
			err1:       err1,
			resultErr:  err1,
		},
	}

	for _, test := range tests {
		xc.On("GetProperty", false, win, prop, x.Atom(x.GetPropertyTypeAny),
			uint32(0), uint32(0)).Return(test.propReply0, test.err0).Once()

		if test.err0 == nil {
			xc.On("GetProperty", false, win, prop, x.Atom(x.GetPropertyTypeAny),
				uint32(0), test.longLength).Return(test.propReply1, test.err1).Once()
		}

		propReply, err := m.getProperty(win, prop, false)
		assert.Equal(t, test.resultPropReply, propReply, test.desc)
		assert.Equal(t, test.resultErr, err, test.desc)
		_ = err
	}

	xc.AssertExpectations(t)
}

func Test_shouldIgnoreSaveTarget(t *testing.T) {
	initAtomsForTest()
	fn := shouldIgnoreSaveTarget
	assert.True(t, fn(atomTimestamp, "TIMESTAMP"))
	assert.True(t, fn(atomTargets, "TARGETS"))

	assert.False(t, fn(200, "image/jpeg"))
	assert.False(t, fn(200, "image/png"))
	assert.False(t, fn(200, "image/bmp"))

	assert.True(t, fn(200, "image/xpm"))
	assert.True(t, fn(200, "image/webp"))

	assert.False(t, fn(200, "text/plain"))
	assert.False(t, fn(200, "application/x-qt-image"))
}

func TestManagerAddGetTargetData(t *testing.T) {
	m := &Manager{}
	td0 := &TargetData{
		Target: 1,
		Type:   x.AtomAtom,
		Data:   []byte{1, 2, 3, 4},
	}
	m.addTargetData(td0)
	assert.Equal(t, td0, m.getTargetData(1))
	assert.Len(t, m.content, 1)

	td1 := &TargetData{
		Target: 1,
		Type:   x.AtomAtom,
		Data:   []byte{1, 2, 3, 4, 5, 6, 7, 8},
	}
	m.addTargetData(td1)
	assert.Equal(t, td1, m.getTargetData(1))
	assert.Len(t, m.content, 1)
}
