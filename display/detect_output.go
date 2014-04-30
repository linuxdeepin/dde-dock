package main

import "github.com/BurntSushi/xgb/randr"
import "github.com/BurntSushi/xgb"
import "strings"
import "sort"

import "github.com/BurntSushi/xgb/xproto"

func guestBuiltIn(ops []*Monitor) *Monitor {
	// It's a bug if there hasn't any Output.
	var mirrorOP *Monitor = ops[0]
	currentType := unknownAtom
	for _, op := range ops {
		t := getContentorType(GetDisplayInfo().outputNames[op.Outputs[0]])
		if !greaterConnectorType(t, currentType) {
			currentType = t
			mirrorOP = op
		}
	}
	return mirrorOP
}

func guestPrimaryName() string {
	//TODO:
	return "TODO"
}

func listModes(op randr.Output) []Mode {
	oinfo, err := randr.GetOutputInfo(xcon, op, LastConfigTimeStamp).Reply()
	if err != nil {
		return nil
	}
	var r []Mode
	for _, id := range oinfo.Modes {
		r = append(r, GetDisplayInfo().modes[id])
	}
	return r
}

func getMatchedSize(ops []randr.Output) (uint16, uint16) {
	if len(ops) == 0 {
		Logger.Error("getMatchedSize with zero length []randr.Output")
		return 0, 0
	}
	defer func() {
		if err := recover(); err != nil {
			Logger.Error(err)
		}
	}()
	var allMatched []Mode
	for _, minfo := range listModes(ops[0]) {
		allHave := true
		for _, op := range ops[1:] {
			has := false
			for _, minfo2 := range listModes(op) {
				if minfo.Width == minfo2.Width && minfo.Height == minfo2.Height {
					has = true
					break
				}
			}
			if !has {
				allHave = false
				break
			}

		}
		if allHave {
			allMatched = append(allMatched, minfo)
		}
	}
	sort.Sort(Modes(allMatched))
	return allMatched[0].Width, allMatched[0].Height
}

func getMatchedSize1(ops []*Monitor) (uint16, uint16) {
	switch len(ops) {
	case 0:
		panic("getMatchedSize received an ops with zero length")
	case 1:
		bestMode := ops[0].BestMode
		return parseRotationSize(ops[0].Rotation, bestMode.Width, bestMode.Height)
	}
	//TODO: calc rotation
	sameModes := make([]Mode, 0)
	first := ops[0]
	for _, modeA := range first.ListModes() {
		allHave := true
		for _, op := range ops[1:] {
			found := false
			for _, modeB := range op.ListModes() {
				if modeA.Width == modeB.Width && modeA.Height == modeA.Height {
					found = true
					break
				}
			}
			if found == false {
				allHave = false
				break
			}
		}
		if allHave {
			sameModes = append(sameModes, modeA)
		}
	}

	bestMode := Mode{}
	for _, mode := range sameModes {
		if bestMode.Width+bestMode.Height <= mode.Width+mode.Height {
			bestMode = mode
		}
	}
	return bestMode.Width, bestMode.Height
}

func getMirrorSize(ops []*Monitor) (uint16, uint16) {
	switch len(ops) {
	case 0:
		return 0, 0
	case 1:
		return parseRotationSize(ops[0].Rotation, ops[0].CurrentMode.Width, ops[0].CurrentMode.Height)
	default:
		builtin := guestBuiltIn(ops)
		oth := make([]*Monitor, 0)
		for _, op := range ops {
			if op != builtin {
				oth = append(oth, op)
			}
		}
		if len(oth) == 0 {
			return parseRotationSize(ops[0].Rotation, builtin.CurrentMode.Width, builtin.CurrentMode.Height)
		}
		return getMatchedSize1(oth)
	}
}

var (
	_VGAAtom          = getAtom(xcon, "VGA")
	_DVIAtom          = getAtom(xcon, "DVI")
	_DVIIAtom         = getAtom(xcon, "DVI-I")
	_DVIAAtom         = getAtom(xcon, "DVI-A")
	_DVIDAtom         = getAtom(xcon, "DVI-D")
	_HDMIAtom         = getAtom(xcon, "HDMI")
	_PanelAtom        = getAtom(xcon, "Panel")
	_TVAtom           = getAtom(xcon, "TV")
	_TVCompositeAtom  = getAtom(xcon, "TV-Composite")
	_TVSVidoeAtom     = getAtom(xcon, "TV-SVideo")
	_TVSComponentAtom = getAtom(xcon, "TV-Component")
	_TVSCARTAtom      = getAtom(xcon, "TV-SCART")
	_TVC4Atom         = getAtom(xcon, "TV-C4")
	_DisplayPort      = getAtom(xcon, "DisplayPort")

	connectorTypeAtom = getAtom(xcon, "ConnectorType")
)

var connectorTypeMap = map[xproto.Atom]int{
	_PanelAtom:        0,
	_VGAAtom:          1,
	_DVIAtom:          2,
	_DVIIAtom:         2,
	_DVIAAtom:         2,
	_DVIDAtom:         2,
	_HDMIAtom:         3,
	_DisplayPort:      3,
	_TVAtom:           4,
	_TVCompositeAtom:  4,
	_TVSVidoeAtom:     4,
	_TVSComponentAtom: 4,
	_TVSCARTAtom:      4,
	_TVC4Atom:         4,
}

func greaterConnectorType(a xproto.Atom, b xproto.Atom) bool {
	if connectorTypeMap[a] > connectorTypeMap[b] {
		return true
	} else {
		return false
	}
}

func getContentorType(op randr.Output) xproto.Atom {
	prop, err := randr.GetOutputProperty(xcon, op, connectorTypeAtom, xproto.AtomAtom, 0, 1, false, false).Reply()
	if err != nil {
		return unknownAtom
	}
	if prop.NumItems == 1 {
		return xproto.Atom(xgb.Get32(prop.Data))
	}

	//many drivers don't implement the ConnectorType property *and* Xserver don't thorw an error when that happend!
	//fallback method: resort the op name
	oinfo, err := randr.GetOutputInfo(xcon, op, xproto.TimeCurrentTime).Reply()
	if err != nil {
		return unknownAtom
	}
	switch {
	case strings.Contains(string(oinfo.Name), "HDMI"):
		return _HDMIAtom
	case strings.Contains(string(oinfo.Name), "VGA"):
		return _VGAAtom
	case strings.Contains(string(oinfo.Name), "LVDS"), strings.Contains(string(oinfo.Name), "LCD"), strings.Contains(string(oinfo.Name), "Lvds"):
		return _PanelAtom
	case strings.Contains(string(oinfo.Name), "eDP"):
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
