package trayicon

import (
	"fmt"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/composite"
	"github.com/linuxdeepin/go-x11-client/ext/damage"
)

type TrayIcon struct {
	win    x.Window
	notify bool
	data   []byte // window pixmap data
	damage damage.Damage
}

func NewTrayIcon(win x.Window) *TrayIcon {
	return &TrayIcon{
		win:    win,
		notify: true,
	}
}

func (icon *TrayIcon) getName() string {
	wmName, _ := ewmhConn.GetWMName(icon.win).Reply(ewmhConn)
	if wmName != "" {
		return wmName
	}

	wmNameTextProp, err := icccmConn.GetWMName(icon.win).Reply(icccmConn)
	if err == nil {
		wmName, _ := wmNameTextProp.GetStr()
		if wmName != "" {
			return wmName
		}
	}

	wmClass, err := icccmConn.GetWMClass(icon.win).Reply(icccmConn)
	if err == nil {
		return fmt.Sprintf("[%s|%s]", wmClass.Class, wmClass.Instance)
	}

	return ""
}

func (icon *TrayIcon) getPixmapData() ([]byte, error) {
	pixmapId, err := XConn.GenerateID()
	if err != nil {
		return nil, err
	}
	pixmap := x.Pixmap(pixmapId)
	err = composite.NameWindowPixmapChecked(XConn, icon.win, pixmap).Check(XConn)
	if err != nil {
		return nil, err
	}
	defer x.FreePixmap(XConn, pixmap)

	geo, err := x.GetGeometry(XConn, x.Drawable(icon.win)).Reply(XConn)
	if err != nil {
		return nil, err
	}

	img, err := x.GetImage(XConn, x.ImageFormatZPixmap, x.Drawable(pixmap),
		0, 0, geo.Width, geo.Height, (1<<32)-1).Reply(XConn)
	if err != nil {
		return nil, err
	}
	return img.Data, nil
}
