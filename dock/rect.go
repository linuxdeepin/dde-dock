package dock

import (
	"github.com/BurntSushi/xgbutil/xrect"
)

type Rect struct {
	X, Y          int32
	Width, Height uint32
}

func NewRect() *Rect {
	return &Rect{}
}

func (r *Rect) ToXRect() xrect.Rect {
	return xrect.New(int(r.X), int(r.Y), int(r.Width), int(r.Height))
}
