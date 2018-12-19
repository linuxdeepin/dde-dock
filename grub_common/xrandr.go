package grub_common

import (
	"github.com/linuxdeepin/go-x11-client"
	"pkg.deepin.io/dde/api/drandr"
)

func getGfxmodesWithOutput(modes drandr.ModeInfos, output drandr.OutputInfo) (result Gfxmodes) {
	for _, modeId := range output.Modes {
		modeInfo := modes.Query(modeId)
		if modeInfo.Width < 1024 || modeInfo.Height < 720 {
			continue
		}

		mode := Gfxmode{
			Width:  int(modeInfo.Width),
			Height: int(modeInfo.Height),
		}
		result = result.Add(mode)
	}
	return
}

func GetGfxmodesFromXRandr() (Gfxmodes, error) {
	xConn, err := x.NewConn()
	if err != nil {
		return nil, err
	}
	defer xConn.Close()
	screenInfo, err := drandr.GetScreenInfo(xConn)
	if err != nil {
		return nil, err
	}
	connectedOutputs := screenInfo.Outputs.ListConnectionOutputs()

	gfxmodeMap := make(map[Gfxmode]int)

	for _, output := range connectedOutputs {
		modes := getGfxmodesWithOutput(screenInfo.Modes, output)
		for _, mode := range modes {
			gfxmodeMap[mode]++
		}
	}

	var result Gfxmodes
	for mode, count := range gfxmodeMap {
		if count == len(connectedOutputs) {
			result = result.Add(mode)
		}
	}

	if len(result) == 0 {
		result = Gfxmodes{Gfxmode{Width: 1024, Height: 768}}
	}

	return result, nil
}
