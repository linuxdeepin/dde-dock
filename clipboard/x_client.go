package clipboard

import (
	"errors"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/xfixes"
)

type XClient interface {
	Conn() *x.Conn
	GetAtom(name string) (x.Atom, error)
	GetAtomName(atom x.Atom) (string, error)
	GetSelectionOwner(selection x.Atom) (x.Window, error)
	CreateWindow() (x.Window, error)
	SetSelectionOwner(owner x.Window, selection x.Atom,
		time x.Timestamp)
	SetSelectionOwnerE(owner x.Window, selection x.Atom,
		time x.Timestamp) error
	ChangePropertyE(mode uint8, window x.Window, property, type0 x.Atom, format uint8, data []byte) error
	GetProperty(delete bool, window x.Window, property, type0 x.Atom, longOffset, longLength uint32) (*x.GetPropertyReply, error)
	DeletePropertyE(window x.Window, property x.Atom) error
	SendEventE(propagate bool, destination x.Window, eventMask uint32, event interface{}) error
	ConvertSelection(requestor x.Window, selection, target, property x.Atom, time x.Timestamp)
	ConvertSelectionE(requestor x.Window, selection, target, property x.Atom, time x.Timestamp) error
	Flush() error
	SelectSelectionInputE(window x.Window, selection x.Atom, eventMask uint32) error
	ChangeWindowEventMask(win x.Window, evMask uint32) error
}

//go:generate mockery -name XClient
type xClient struct {
	conn *x.Conn
}

func (xc *xClient) GetAtom(name string) (x.Atom, error) {
	return xc.conn.GetAtom(name)
}

func (xc *xClient) GetAtomName(atom x.Atom) (string, error) {
	return xc.conn.GetAtomName(atom)
}

func (xc *xClient) GetSelectionOwner(selection x.Atom) (x.Window, error) {
	return getSelectionOwner(xc.conn, selection)
}

func (xc *xClient) CreateWindow() (x.Window, error) {
	return createWindow(xc.conn)
}

func (xc *xClient) SetSelectionOwner(owner x.Window, selection x.Atom,
	time x.Timestamp) {
	x.SetSelectionOwner(xc.conn, owner, selection, time)
}

func (xc *xClient) SetSelectionOwnerE(owner x.Window, selection x.Atom,
	time x.Timestamp) error {
	return x.SetSelectionOwnerChecked(xc.conn, owner, selection, time).Check(xc.conn)
}

func (xc *xClient) GetTimestamp() {

}

func (xc *xClient) ChangePropertyE(mode uint8, window x.Window, property, type0 x.Atom, format uint8, data []byte) error {
	return x.ChangePropertyChecked(xc.conn, mode, window, property, type0, format,
		data).Check(xc.conn)
}

func (xc *xClient) GetProperty(delete bool, window x.Window, property, type0 x.Atom, longOffset, longLength uint32) (*x.GetPropertyReply, error) {
	return x.GetProperty(xc.conn, delete, window, property, type0, longOffset, longLength).Reply(xc.conn)
}

func (xc *xClient) DeletePropertyE(window x.Window, property x.Atom) error {
	return x.DeletePropertyChecked(xc.conn, window, property).Check(xc.conn)
}

func (xc *xClient) SendEventE(propagate bool, destination x.Window, eventMask uint32, event interface{}) error {
	w := x.NewWriter()
	switch e := event.(type) {
	case *x.SelectionNotifyEvent:
		x.WriteSelectionNotifyEvent(w, e)
	default:
		return errors.New("unsupported event type")
	}

	return x.SendEventChecked(xc.conn, propagate, destination, eventMask,
		w.Bytes()).Check(xc.conn)
}

func (xc *xClient) ConvertSelectionE(requestor x.Window, selection, target, property x.Atom, time x.Timestamp) error {
	return x.ConvertSelectionChecked(xc.conn, requestor, selection, target, property, time).Check(xc.conn)
}

func (xc *xClient) ConvertSelection(requestor x.Window, selection, target, property x.Atom, time x.Timestamp) {
	x.ConvertSelection(xc.conn, requestor, selection, target, property, time)
}

func (xc *xClient) Flush() error {
	return xc.conn.Flush()
}

func (xc *xClient) SelectSelectionInputE(window x.Window, selection x.Atom, eventMask uint32) error {
	return xfixes.SelectSelectionInputChecked(xc.conn, window, selection, eventMask).Check(xc.conn)
}

func (xc *xClient) ChangeWindowEventMask(win x.Window, evMask uint32) error {
	const valueMask = x.CWEventMask
	return x.ChangeWindowAttributesChecked(xc.conn, win, valueMask, []uint32{evMask}).Check(xc.conn)
}

func (xc *xClient) Conn() *x.Conn {
	return xc.conn
}
