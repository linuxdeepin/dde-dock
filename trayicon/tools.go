package trayicon

import (
	"bytes"
	"crypto/md5"
	"github.com/BurntSushi/xgb/composite"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/xgraphics"
)

func isValidWindow(win xproto.Window) bool {
	r, err := xproto.GetWindowAttributes(XU.Conn(), win).Reply()
	return r != nil && err == nil
}

func findRGBAVisualID() xproto.Visualid {
	for _, dinfo := range XU.Screen().AllowedDepths {
		if dinfo.Depth == 32 {
			for _, vinfo := range dinfo.Visuals {
				return vinfo.VisualId
			}
		}
	}
	return XU.Screen().RootVisual
}

func icon2md5(win xproto.Window) []byte {
	pixmap, _ := xproto.NewPixmapId(XU.Conn())
	defer xproto.FreePixmap(XU.Conn(), pixmap)
	if err := composite.NameWindowPixmapChecked(XU.Conn(), win, pixmap).Check(); err != nil {
		logger.Warning("NameWindowPixmap failed:", err, win)
		return nil
	}
	im, err := xgraphics.NewDrawable(XU, xproto.Drawable(pixmap))
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
