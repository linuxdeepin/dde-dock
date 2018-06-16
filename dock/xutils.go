/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package dock

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/png"
	"math"
	"strings"

	"github.com/nfnt/resize"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/util/wm/ewmh"
	"github.com/linuxdeepin/go-x11-client/util/wm/icccm"
)

func atomsContains(slice []x.Atom, atom x.Atom) bool {
	for _, a := range slice {
		if a == atom {
			return true
		}
	}
	return false
}

func getWmClass(win x.Window) (*icccm.WMClass, error) {
	wmClass, err := globalIcccmConn.GetWMClass(win).Reply(globalIcccmConn)
	if err != nil {
		return nil, err
	}
	return &wmClass, nil
}

func getAtomName(atom x.Atom) (string, error) {
	reply, err := x.GetAtomName(globalXConn, atom).Reply(globalXConn)
	if err != nil {
		return "", err
	}
	return reply.Name, nil
}

func getAtom(name string) (x.Atom, error) {
	reply, err := x.InternAtom(globalXConn, false, name).Reply(globalXConn)
	if err != nil {
		return 0, err
	}
	return reply.Atom, nil
}

func maximizeWindow(win x.Window) error {
	return globalEwmhConn.RequestChangeWMState(win, ewmh.WMStateAdd, atomNetWmStateMaximizedVert, atomNetWmStateMaximizedHorz, 2).Check(globalXConn)
	return nil
}

func minimizeWindow(win x.Window) error {
	return globalIcccmConn.RequestChangeWMState(win, icccm.StateIconic).Check(globalXConn)
}

func makeWindowAbove(win x.Window) error {
	return globalEwmhConn.RequestChangeWMState(win, ewmh.WMStateAdd, atomNetWmStateAbove,
		0, 2).Check(globalXConn)
}

func moveWindow(win x.Window) error {
	// TODO:
	//return ewmh.WmMoveresize(xu, win, ewmh.MoveKeyboard)
	//globalEwmhConn.RequestWMMoveResize(win, )
	return nil
}

func closeWindow(win x.Window, ts x.Timestamp) error {
	return globalEwmhConn.RequestCloseWindow(win, ts, 2).Check(globalXConn)
}

type windowFrameExtents struct {
	Left, Right, Top, Bottom uint
}

func getWindowFrameExtents(xConn *x.Conn, win x.Window) (*windowFrameExtents, error) {
	reply, err := x.GetProperty(xConn, false, win, atomNetFrameExtents, x.AtomCardinal,
		0, 4).Reply(xConn)
	if err != nil || reply.Format == 0 {
		// try _GTK_FRAME_EXTENTS
		reply, err = x.GetProperty(xConn, false, win, atomGtkFrameExtents, x.AtomCardinal,
			0, 4).Reply(xConn)
		if err != nil {
			return nil, err
		}
	}

	return getFrameExtentsFromReply(reply)
}

func getCardinalsFromReply(r *x.GetPropertyReply) ([]uint32, error) {
	if r.Format != 32 {
		return nil, errors.New("bad reply")
	}
	count := len(r.Value) / 4
	ret := make([]uint32, count)
	rdr := x.NewReaderFromData(r.Value)
	for i := 0; i < count; i++ {
		ret[i] = uint32(rdr.Read4b())
	}
	return ret, nil
}

func getFrameExtentsFromReply(reply *x.GetPropertyReply) (*windowFrameExtents, error) {
	list, err := getCardinalsFromReply(reply)
	if err != nil {
		return nil, err
	}

	if len(list) != 4 {
		return nil, errors.New("length of list is not 4")
	}
	return &windowFrameExtents{
		Left:   uint(list[0]),
		Right:  uint(list[1]),
		Top:    uint(list[2]),
		Bottom: uint(list[3]),
	}, nil
}

func getWindowGeometry(xConn *x.Conn, win x.Window) (*Rect, error) {
	winRect, err := getDecorGeometry(xConn, win)
	if err != nil {
		return nil, err
	}
	frameExtents, _ := getWindowFrameExtents(xConn, win)
	if frameExtents != nil {
		x := winRect.X + int32(frameExtents.Left)
		y := winRect.Y + int32(frameExtents.Top)
		w := winRect.Width - uint32(frameExtents.Left+frameExtents.Right)
		h := winRect.Height - uint32(frameExtents.Top+frameExtents.Bottom)
		return &Rect{x, y, w, h}, nil
	}
	return winRect, nil
}

func getWindowParent(xConn *x.Conn, win x.Window) (x.Window, error) {
	reply, err := x.QueryTree(xConn, win).Reply(xConn)
	if err != nil {
		return 0, err
	}
	return reply.Parent, nil
}

func getDecorGeometry(xConn *x.Conn, win x.Window) (*Rect, error) {
	rootWin := xConn.GetDefaultScreen().Root
	parent := win
	for {
		tempParent, err := getWindowParent(xConn, parent)
		if err != nil || tempParent == rootWin {
			break
		}
		parent = tempParent
	}
	return getRawGeometry(xConn, parent)
}

func getRawGeometry(xConn *x.Conn, win x.Window) (*Rect, error) {
	geo, err := x.GetGeometry(xConn, x.Drawable(win)).Reply(xConn)
	if err != nil {
		return nil, err
	}
	return &Rect{
		X:      int32(geo.X),
		Y:      int32(geo.Y),
		Width:  uint32(geo.Width),
		Height: uint32(geo.Height),
	}, nil
}

func getWmName(win x.Window) string {
	// get _NET_WM_NAME
	name, err := globalEwmhConn.GetWMName(win).Reply(globalEwmhConn)
	if err != nil || name == "" {
		// get WM_NAME
		nameTp, _ := globalIcccmConn.GetWMName(win).Reply(globalIcccmConn)
		name, _ = nameTp.GetStr()
	}

	return strings.Replace(name, "\x00", "", -1)
}

func getWmPid(win x.Window) uint {
	pid, _ := globalEwmhConn.GetWMPid(win).Reply(globalEwmhConn)
	return uint(pid)
}

// WM_CLIENT_LEADER
func getWmClientLeader(win x.Window) (x.Window, error) {
	reply, err := x.GetProperty(globalXConn, false, win, atomWmClientLeader, atomUTF8String,
		0, 1).Reply(globalXConn)
	if err != nil {
		return 0, err
	}
	leader, err := getWindowFromReply(reply)
	if err != nil {
		return 0, err
	}
	return leader, nil
}

// WM_TRANSIENT_FOR
func getWmTransientFor(win x.Window) (x.Window, error) {
	return globalIcccmConn.GetWMTransientFor(win).Reply(globalIcccmConn)
}

// _NET_WM_WINDOW_OPACITY
func getWmWindowOpacity(win x.Window) (uint, error) {
	reply, err := x.GetProperty(globalXConn, false, win, atomNetWmWindowOpacity, atomUTF8String,
		0, 1).Reply(globalXConn)
	if err != nil {
		return 0, err
	}

	opacity, err := getCardinalFromReply(reply)
	if err != nil {
		return 0, err
	}

	return uint(opacity), nil
}

func getWmCommand(win x.Window) ([]string, error) {
	reply, err := x.GetProperty(globalXConn, false, win, atomWmCommand, atomUTF8String,
		0, lengthMax).Reply(globalXConn)
	if err != nil {
		return nil, err
	}
	return getUTF8StrsFromReply(reply)
}

func getWindowGtkApplicationId(win x.Window) string {
	gtkAppId, _ := getWindowPropertyString(win, atomGtkApplicationId)
	return gtkAppId
}

func getWindowFlatpakAppID(win x.Window) string {
	id, _ := getWindowPropertyString(win, atomFlatpakAppId)
	return id
}

const lengthMax = 0xffff

func getWmWindowRole(win x.Window) string {
	role, _ := getWindowPropertyString(win, atomWmWindowRole)
	return role
}

func getWindowPropertyString(win x.Window, atom x.Atom) (string, error) {
	reply, err := x.GetProperty(globalXConn, false, win, atom, atomUTF8String,
		0, lengthMax).Reply(globalXConn)
	if err != nil {
		return "", err
	}
	return getUTF8StrFromReply(reply)
}

func getCardinalFromReply(r *x.GetPropertyReply) (uint32, error) {
	if r.Format != 32 || len(r.Value) != 4 {
		return 0, errors.New("bad reply")
	}
	return uint32(x.Get32(r.Value)), nil
}

func getWindowFromReply(r *x.GetPropertyReply) (x.Window, error) {
	if r.Format != 32 || len(r.Value) != 4 {
		return 0, errors.New("bad reply")
	}
	return x.Window(x.Get32(r.Value)), nil
}

func getUTF8StrFromReply(reply *x.GetPropertyReply) (string, error) {
	if reply.Format != 8 {
		return "", errors.New("bad reply")
	}

	return string(reply.Value), nil
}

func getUTF8StrsFromReply(reply *x.GetPropertyReply) ([]string, error) {
	if reply.Format != 8 {
		return nil, errors.New("bad reply")
	}

	data := reply.Value
	var strs []string
	sstart := 0
	for i, c := range data {
		if c == 0 {
			strs = append(strs, string(data[sstart:i]))
			sstart = i + 1
		}
	}
	if sstart < len(data) {
		strs = append(strs, string(data[sstart:]))
	}
	return strs, nil
}

const bestIconSize = 48

func getIconFromWindow(win x.Window) string {
	img, err := findIconEwmh(win)
	if err != nil {
		logger.Warning(err)
		// try icccm
		img, err = findIconIcccm(win)
		if err != nil {
			logger.Warning(err)
			// get icon failed
			return ""
		}
	}

	img = resize.Thumbnail(bestIconSize, bestIconSize, img, resize.NearestNeighbor)

	// encode image to png, then to base64 string
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		logger.Warning(err)
		return ""
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

func findIconEwmh(win x.Window) (image.Image, error) {
	icon, err := getBestEwmhIcon(win)
	if err != nil {
		return nil, err
	}
	return NewNRGBAImageFromEwmhIcon(icon), nil
}

// findIconIcccm helps FindIcon by trying to return an icccm-style icon.
func findIconIcccm(wid x.Window) (image.Image, error) {
	// TODO
	return nil, errors.New("todo")
	//hints, err := icccm.WmHintsGet(xu, wid)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// Only continue if the WM_HINTS flags say an icon is specified and
	//// if at least one of icon pixmap or icon mask is non-zero.
	//if hints.Flags&icccm.HintIconPixmap == 0 ||
	//	(hints.IconPixmap == 0 && hints.IconMask == 0) {
	//
	//	return nil, errors.New("No icon found in WM_HINTS.")
	//}
	//
	//return xgraphics.NewIcccmIcon(xu, hints.IconPixmap, hints.IconMask)
}

func getBestEwmhIcon(win x.Window) (*ewmh.WMIcon, error) {
	icons, err := globalEwmhConn.GetWMIcon(win).Reply(globalEwmhConn)
	if err != nil {
		return nil, err
	}

	best := findBestEwmhIcon(bestIconSize, bestIconSize, icons)
	if best == nil {
		return nil, errors.New("ewmh icon not found")
	}
	return best, nil
}

// findBestEwmhIcon takes width/height dimensions and a slice of *ewmh.WmIcon
// and finds the best matching icon of the bunch. We always prefer bigger.
// If no icons are bigger than the preferred dimensions, use the biggest
// available. Otherwise, use the smallest icon that is greater than or equal
// to the preferred dimensions. The preferred dimensions is essentially
// what you'll likely scale the resulting icon to.
// If width and height are 0, then the largest icon found will be returned.
func findBestEwmhIcon(width, height int, icons []ewmh.WMIcon) *ewmh.WMIcon {
	// nada nada limonada
	if len(icons) == 0 {
		return nil
	}

	parea := width * height // preferred size
	best := -1

	// If zero area, set it to the largest possible.
	if parea == 0 {
		parea = math.MaxInt32
	}

	var bestArea, iconArea int

	for i, icon := range icons {
		// the first valid icon we've seen; use it!
		if best == -1 {
			best = i
			continue
		}

		// load areas for comparison
		bestArea = int(icons[best].Width * icons[best].Height)
		iconArea = int(icon.Width * icon.Height)

		// We don't always want to accept bigger icons if our best is
		// already bigger. But we always want something bigger if our best
		// is insufficient.
		if (iconArea >= parea && iconArea <= bestArea) ||
			(bestArea < parea && iconArea > bestArea) {
			best = i
		}
	}

	if best > -1 {
		return &icons[best]
	}
	return nil
}

func NewNRGBAImageFromEwmhIcon(icon *ewmh.WMIcon) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, int(icon.Width), int(icon.Height)))
	// icon.Data []uint8 BGRA
	// img.Pix []uint8 RGBA

	for i := 0; i < len(icon.Data); i += 4 {
		b := icon.Data[i]
		g := icon.Data[i+1]
		r := icon.Data[i+2]
		a := icon.Data[i+3]

		img.Pix[i] = r
		img.Pix[i+1] = g
		img.Pix[i+2] = b
		img.Pix[i+3] = a
	}
	return img
}

func getWindowUserTime(win x.Window) (uint, error) {
	timestamp, err := globalEwmhConn.GetWMUserTime(win).Reply(globalEwmhConn)
	if err != nil {
		userTimeWindow, err := globalEwmhConn.GetWMUserTimeWindow(win).Reply(globalEwmhConn)
		if err != nil {
			return 0, err
		}

		timestamp, err = globalEwmhConn.GetWMUserTime(userTimeWindow).Reply(globalEwmhConn)
		if err != nil {
			return 0, err
		}
	}
	return uint(timestamp), nil
}

func changeCurrentWorkspaceToWindowWorkspace(win x.Window) error {
	winWorkspace, err := globalEwmhConn.GetWMDesktop(win).Reply(globalEwmhConn)
	if err != nil {
		return err
	}

	currentWorkspace, err := globalEwmhConn.GetCurrentDesktop().Reply(globalEwmhConn)
	if err != nil {
		return err
	}

	if currentWorkspace == winWorkspace {
		logger.Debugf("No need to change workspace, the current desktop is already %v", currentWorkspace)
		return nil
	}
	logger.Debug("Change workspace")

	winUserTime, err := getWindowUserTime(win)
	logger.Debug("window user time:", winUserTime)
	if err != nil {
		// only warning not return
		logger.Warning("getWindowUserTime failed:", err)
	}
	err = globalEwmhConn.RequestChangeCurrentDestkop(winWorkspace, x.Timestamp(winUserTime)).Check(globalXConn)

	return nil
}

func activateWindow(win x.Window) error {
	logger.Debug("activateWindow", win)
	err := changeCurrentWorkspaceToWindowWorkspace(win)
	if err != nil {
		return err
	}
	return globalEwmhConn.RequestChangeActiveWindow(win,
		2, 0, 0).Check(globalXConn)
}

func isHiddenPre(win x.Window) bool {
	state, _ := globalEwmhConn.GetWMState(win).Reply(globalEwmhConn)
	return atomsContains(state, atomNetWmStateHidden)
}

// works for new deepin wm.
func isWindowOnCurrentWorkspace(win x.Window) (bool, error) {
	winWorkspace, err := globalEwmhConn.GetWMDesktop(win).Reply(globalEwmhConn)
	if err != nil {
		return false, err
	}

	currentWorkspace, err := globalEwmhConn.GetCurrentDesktop().Reply(globalEwmhConn)
	if err != nil {
		return false, err
	}

	return winWorkspace == currentWorkspace, nil
}

func onCurrentWorkspacePre(win x.Window) bool {
	isOnCurrentWorkspace, err := isWindowOnCurrentWorkspace(win)
	if err != nil {
		logger.Warning(err)
		// 也许是窗口跳过窗口管理器了，如 dde-control-center
		return true
	}
	return isOnCurrentWorkspace
}

func isGoodWindow(win x.Window) bool {
	_, err := x.GetGeometry(globalXConn, x.Drawable(win)).Reply(globalXConn)
	return err == nil
}

func killClient(win x.Window) {
	err := x.KillClientChecked(globalXConn, uint32(win)).Check(globalXConn)
	if err != nil {
		logger.Warning(err)
	}
}
