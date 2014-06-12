/**
 * Copyright (c) 2013 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package grub2

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xwindow"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func quoteString(str string) string {
	return strconv.Quote(str)
}

func unquoteString(str string) string {
	if strings.HasPrefix(str, `"`) && strings.HasSuffix(str, `"`) {
		s, _ := strconv.Unquote(str)
		return s
	} else if strings.HasPrefix(str, `'`) && strings.HasSuffix(str, `'`) {
		return str[1 : len(str)-1]
	}
	return str
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func isFileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
}

// Get all screen's best resolution and choose a smaller one for there
// is no screen is primary.
func getPrimaryScreenBestResolution() (w uint16, h uint16) {
	w, h = 1024, 768 // default value

	XU, err := xgbutil.NewConn()
	if err != nil {
		return
	}
	err = randr.Init(XU.Conn())
	if err != nil {
		return
	}
	_, err = randr.QueryVersion(XU.Conn(), 1, 4).Reply()
	if err != nil {
		return
	}
	Root := xproto.Setup(XU.Conn()).DefaultScreen(XU.Conn()).Root
	resources, err := randr.GetScreenResources(XU.Conn(), Root).Reply()
	if err != nil {
		return
	}

	bestModes := make([]uint32, 0)
	for _, output := range resources.Outputs {
		reply, err := randr.GetOutputInfo(XU.Conn(), output, 0).Reply()
		if err == nil && reply.NumModes > 1 {
			bestModes = append(bestModes, uint32(reply.Modes[0]))
		}
	}

	w, h = 0, 0
	for _, m := range resources.Modes {
		for _, id := range bestModes {
			if id == m.Id {
				bw, bh := m.Width, m.Height
				if w == 0 || h == 0 {
					w, h = bw, bh
				} else if uint32(bw)*uint32(bh) < uint32(w)*uint32(h) {
					w, h = bw, bh
				}
			}
		}
	}

	if w == 0 || h == 0 {
		// get resource failed, use root window's geometry
		rootRect := xwindow.RootGeometry(XU)
		w, h = uint16(rootRect.Width()), uint16(rootRect.Height())
	}

	if w == 0 || h == 0 {
		w, h = 1024, 768 // default value
	}

	logger.Debugf("primary screen's best resolution is %dx%d", w, h)
	return
}

func delta(v1, v2 float64) float64 {
	if v1 > v2 {
		return v1 - v2
	}
	return v2 - v1
}

func isSymlink(file string) bool {
	f, err := os.Lstat(file)
	if err != nil {
		return false
	}
	if f.Mode()&os.ModeSymlink == os.ModeSymlink {
		// This is a symlink
		return true
	}

	// Not a symlink
	return false
}

func copyFile(src, dest string) (written int64, err error) {
	if dest == src {
		return -1, fmt.Errorf("source and destination are same file")
	}

	sf, err := os.Open(src)
	if err != nil {
		return
	}
	defer sf.Close()
	df, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return
	}
	defer df.Close()
	return io.Copy(df, sf)
}

func execAndWait(timeout int, name string, arg ...string) (stdout, stderr string, err error) {
	cmd := exec.Command(name, arg...)
	var bufStdout, bufStderr bytes.Buffer
	cmd.Stdout = &bufStdout
	cmd.Stderr = &bufStderr
	err = cmd.Start()
	if err != nil {
		return
	}

	// wait for process finished
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		if err = cmd.Process.Kill(); err != nil {
			return
		}
		<-done
		err = fmt.Errorf("time out and process was killed")
	case err = <-done:
		stdout = bufStdout.String()
		stderr = bufStderr.String()
		if err != nil {
			return
		}
	}
	return
}
