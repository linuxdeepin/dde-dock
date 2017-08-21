package trayicon

import (
	x "github.com/linuxdeepin/go-x11-client"
)

func isValidWindow(win x.Window) bool {
	reply, err := x.GetWindowAttributes(XConn, win).Reply(XConn)
	return reply != nil && err == nil
}

func findRGBAVisualID() x.VisualID {
	screen := XConn.GetDefaultScreen()
	for _, dinfo := range screen.AllowedDepths {
		if dinfo.Depth == 32 {
			for _, vinfo := range dinfo.Visuals {
				return vinfo.VisualId
			}
		}
	}
	return screen.RootVisual
}
