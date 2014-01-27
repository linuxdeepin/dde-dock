package main

import "github.com/BurntSushi/xgb/randr"
import "github.com/BurntSushi/xgb"
import "strings"

import "github.com/BurntSushi/xgb/xproto"

func (dpy *Display) SetMirrorMode(v bool) {
	dpy.setPropMirrorMode(v)
	if v && dpy.MirrorOutput == nil {
		dpy.SetMirrorOutput(deduceMirrorOutput(dpy.Outputs))
	}
}
func (dpy *Display) SetMirrorOutput(op *Output) {
	if op.Opened {
		op.pendingConfig = NewPendingConfig(op).SetPos(0, 0).SetScale(1, 1).SetRotation(randr.RotationRotate0)
		dpy.setPropMirrorOutput(op)
		DPY.ApplyChanged()
	}
}

func deduceMirrorOutput(ops []*Output) *Output {
	// It's a bug if there isn't any Output.
	var mirrorOP *Output = ops[0]
	currentType := unknownAtom
	for _, op := range ops {
		if op.Opened {
			t := getContentorType(op.Identify)
			if greterConnectorType(t, currentType) {
				currentType = t
				mirrorOP = op
			}
		}
	}
	return mirrorOP
}

var (
	_VGAAtom          = getAtom(X, "VGA")
	_DVIAtom          = getAtom(X, "DVI")
	_DVIIAtom         = getAtom(X, "DVI-I")
	_DVIAAtom         = getAtom(X, "DVI-A")
	_DVIDAtom         = getAtom(X, "DVI-D")
	_HDMIAtom         = getAtom(X, "HDMI")
	_PanelAtom        = getAtom(X, "Panel")
	_TVAtom           = getAtom(X, "TV")
	_TVCompositeAtom  = getAtom(X, "TV-Composite")
	_TVSVidoeAtom     = getAtom(X, "TV-SVideo")
	_TVSComponentAtom = getAtom(X, "TV-Component")
	_TVSCARTAtom      = getAtom(X, "TV-SCART")
	_TVC4Atom         = getAtom(X, "TV-C4")
	_DisplayPort      = getAtom(X, "DisplayPort")
)

var connectorTypeMap = map[xproto.Atom]int{
	_PanelAtom:        0,
	_VGAAtom:          1,
	_DVIAtom:          2,
	_DVIIAtom:         2,
	_DVIAAtom:         2,
	_DVIDAtom:         2,
	_HDMIAtom:         3,
	_TVAtom:           4,
	_TVCompositeAtom:  4,
	_TVSVidoeAtom:     4,
	_TVSComponentAtom: 4,
	_TVSCARTAtom:      4,
	_TVC4Atom:         4,
	_DisplayPort:      5,
}

func greterConnectorType(a xproto.Atom, b xproto.Atom) bool {
	if connectorTypeMap[a] > connectorTypeMap[b] {
		return true
	} else {
		return false
	}
}

func getContentorType(op randr.Output) xproto.Atom {
	prop, err := randr.GetOutputProperty(X, op, connectorTypeAtom, xproto.AtomAtom, 0, 1, false, false).Reply()
	if err != nil {
		return unknownAtom
	}
	if prop.NumItems == 1 {
		return xproto.Atom(xgb.Get32(prop.Data))
	}

	//many drivers don't implement the ConnectorType property *and* Xserver don't thorw an error when that happend!
	//fallback method: resort the op name
	oinfo, err := randr.GetOutputInfo(X, op, xproto.TimeCurrentTime).Reply()
	if err != nil {
		return unknownAtom
	}
	switch {
	case strings.Contains(string(oinfo.Name), "VGA"):
		return _VGAAtom
	case strings.Contains(string(oinfo.Name), "LVDS"), strings.Contains(string(oinfo.Name), "LCD"), strings.Contains(string(oinfo.Name), "Lvds"):
		return _PanelAtom
	case strings.Contains(string(oinfo.Name), "DP"):
		return _DisplayPort
	case strings.Contains(string(oinfo.Name), "TV"):
		return _TVAtom
	case strings.Contains(string(oinfo.Name), "TMDS"), strings.Contains(string(oinfo.Name), "DVI"):
		return _DVIAtom
	case strings.Contains(string(oinfo.Name), "S-video"):
		return _TVAtom
	default:
		return unknownAtom
	}
}
