package dock

import (
	"bytes"
	"dlib/gio-2.0"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func isEntryNameValid(name string) bool {
	if !strings.HasPrefix(name, entryDestPrefix) {
		return false
	}
	return true
}

func getEntryId(name string) (string, bool) {
	a := strings.SplitN(name, entryDestPrefix, 2)
	if len(a) >= 1 {
		return a[len(a)-1], true
	}
	return "", false
}

func guess_desktop_id(oldId string) string {
	allApp := gio.AppInfoGetAll()
	for _, app := range allApp {
		baseName := filepath.Base(gio.ToDesktopAppInfo(app).GetFilename())
		if oldId == strings.ToLower(baseName) {
			return baseName
		}
	}

	return ""
}

func getAppIcon(core *gio.DesktopAppInfo) string {
	gioIcon := core.GetIcon()
	if gioIcon == nil {
		logger.Warning("get icon from appinfo failed")
		return ""
	}

	logger.Debug("GetIcon:", gioIcon.ToString())
	icon := get_theme_icon(gioIcon.ToString(), 48)
	if icon == "" {
		logger.Warning("get icon from theme failed")
		return ""
	}

	logger.Debug("get_theme_icon:", icon)
	// the filepath.Ext return ".xxx"
	ext := filepath.Ext(icon)[1:]
	logger.Debug("ext:", ext)
	if strings.EqualFold(ext, "xpm") {
		logger.Debug("change xpm to data uri")
		return xpm_to_dataurl(icon)
	}

	return icon
}

func unsetEnv(env string) {
	_, _, err := execAndWait(5, "/usr/bin/env", "-u", env)
	if err != nil {
		logger.Error(err)
	}
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
