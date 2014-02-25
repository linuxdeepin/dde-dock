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

package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"io"
	"os"
	"os/exec"
	"path"
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

// TODO move to go-dlib/os, dde-api/os
func execAndWait(timeout int, name string, arg ...string) (stdout, stderr string, err error) {
	cmd := exec.Command(name, arg...)
	var buf_stdout, buf_stderr bytes.Buffer
	cmd.Stdout = &buf_stdout
	cmd.Stderr = &buf_stderr
	err = cmd.Start()
	if err != nil {
		_LOGGER.Error(err.Error())
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
			_LOGGER.Error(err.Error())
			return
		}
		<-done
		_LOGGER.Info("time out and process was killed")
	case err = <-done:
		stdout = buf_stdout.String()
		stderr = buf_stderr.String()
		if err != nil {
			_LOGGER.Error("process done with error = %v", err)
			return
		}
	}
	return
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// TODO move dde-api/file
func unTarGz(archiveFile string, destDir string, prefix string) error {
	destDir = path.Clean(destDir) + string(os.PathSeparator)

	// open the archive file
	fr, err := os.Open(archiveFile)
	if err != nil {
		return err
	}
	defer fr.Close()

	// create a gzip reader
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return err
	}
	defer gr.Close()

	// create a tar reader
	tr := tar.NewReader(gr)

	// loop files
	for hdr, err := tr.Next(); err != io.EOF; hdr, err = tr.Next() {
		if err != nil {
			return err
		}

		if !strings.HasPrefix(hdr.Name, prefix) {
			continue
		}

		fi := hdr.FileInfo()
		destFullPath := destDir + hdr.Name
		_LOGGER.Info("UnTarGzing file: " + hdr.Name)

		if hdr.Typeflag == tar.TypeDir {
			// create dir
			os.MkdirAll(destFullPath, fi.Mode().Perm())
			os.Chmod(destFullPath, fi.Mode().Perm())
		} else {
			// create the parent dir for file
			os.MkdirAll(path.Dir(destFullPath), fi.Mode().Perm())

			// write data to file
			fw, err := os.Create(destFullPath)
			if err != nil {
				return err
			}
			_, err = io.Copy(fw, tr)
			if err != nil {
				return err
			}
			fw.Close()

			os.Chmod(destFullPath, fi.Mode().Perm())
		}
	}
	return nil
}

// find if a file in archive and return its path
func findFileInTarGz(archiveFile string, targetFile string) (string, error) {
	// open the archive file
	fr, err := os.Open(archiveFile)
	if err != nil {
		return "", err
	}
	defer fr.Close()

	// create a gzip reader
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return "", err
	}
	defer gr.Close()

	// create a tar reader
	tr := tar.NewReader(gr)

	// loop files
	targetPath := ""
	for hdr, err := tr.Next(); err != io.EOF; hdr, err = tr.Next() {
		if err != nil {
			return "", err
		}

		if hdr.Typeflag != tar.TypeDir && strings.HasSuffix(hdr.Name, targetFile) {
			targetPath = hdr.Name
			break
		}
	}
	return targetPath, nil
}

func isFileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	} else {
		return false
	}
}

func copyFile(src, dest string) (written int64, err error) {
	if dest == src {
		return -1, newError("source and destination are same file")
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

func getPathLevel(p string) int {
	p = path.Clean(p)
	if len(p) == 0 {
		return 0
	}
	lv := len(strings.Split(p, string(os.PathSeparator)))
	if strings.HasPrefix(p, "/") || strings.HasPrefix(p, ".") {
		lv--
	}
	return lv
}

func newError(format string, v ...interface{}) error {
	return errors.New(fmt.Sprintf(format, v...))
}

// Get all screen's best resolution and choose a smaller one for there
// is no screen is primary.
func getPrimaryScreenBestResolution() (w uint16, h uint16) {
	w, h = 1024, 768 // default value

	X, err := xgb.NewConn()
	if err != nil {
		return
	}
	err = randr.Init(X)
	if err != nil {
		return
	}
	_, err = randr.QueryVersion(X, 1, 4).Reply()
	if err != nil {
		return
	}
	Root := xproto.Setup(X).DefaultScreen(X).Root
	resources, err := randr.GetScreenResources(X, Root).Reply()
	if err != nil {
		return
	}

	bestModes := make([]uint32, 0)
	for _, output := range resources.Outputs {
		reply, err := randr.GetOutputInfo(X, output, 0).Reply()
		if err == nil && reply.NumModes > 1 {
			bestModes = append(bestModes, uint32(reply.Modes[0]))
		}
	}

	w, h = 0, 0
	for _, m := range resources.Modes {
		for _, id := range bestModes {
			if id == m.Id {
				bw, bh := m.Width, m.Height
				if w*h == 0 {
					w, h = bw, bh
				} else if bw*bh < w*h {
					w, h = bw, bh
				}
			}
		}
	}

	_LOGGER.Info("primary screen's best resolution is %dx%d", w, h)
	return
}

func getImgClipSizeByResolution(screenWidth, screenHeight uint16, imgWidth, imgHeight int32) (w int32, h int32) {
	if imgWidth >= int32(screenWidth) && imgHeight >= int32(screenHeight) {
		w = int32(screenWidth)
		h = int32(screenHeight)
	} else {
		scale := float32(screenWidth) / float32(screenHeight)
		w = imgWidth
		h = int32(float32(w) / scale)
		if h > imgHeight {
			h = imgHeight
			w = int32(float32(h) * scale)
		}
	}
	return
}
