package trayicon

import (
	"bytes"
	"crypto/md5"
	"github.com/BurntSushi/xgb/composite"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/xgraphics"
)

func isValidWindow(xid xproto.Window) bool {
	r, err := xproto.GetWindowAttributes(TrayXU.Conn(), xid).Reply()
	return r != nil && err == nil
}

func findRGBAVisualID() xproto.Visualid {
	for _, dinfo := range TrayXU.Screen().AllowedDepths {
		if dinfo.Depth == 32 {
			for _, vinfo := range dinfo.Visuals {
				return vinfo.VisualId
			}
		}
	}
	return TrayXU.Screen().RootVisual
}

func icon2md5(xid xproto.Window) []byte {
	pixmap, _ := xproto.NewPixmapId(TrayXU.Conn())
	defer xproto.FreePixmap(TrayXU.Conn(), pixmap)
	if err := composite.NameWindowPixmapChecked(TrayXU.Conn(), xid, pixmap).Check(); err != nil {
		logger.Warning("NameWindowPixmap failed:", err, xid)
		return nil
	}
	im, err := xgraphics.NewDrawable(TrayXU, xproto.Drawable(pixmap))
	if err != nil {
		logger.Warning("Create xgraphics.Image failed:", err, pixmap)
		return nil
	}
	buf := bytes.NewBuffer(nil)
	im.WritePng(buf)
	hasher := md5.New()
	hasher.Write(buf.Bytes())
	return hasher.Sum(nil)
}

func md5Equal(a []byte, b []byte) bool {
	if len(a) != 16 || len(b) != 16 {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
