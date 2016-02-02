package mpris

import (
	"fmt"
	"os/exec"
	"gir/gio-2.0"
)

const (
	mimeTypeBrowser = "x-scheme-handler/http"
	mimeTypeEmail   = "x-scheme-handler/mailto"
	mimeTypeCalc    = "x-scheme-handler/calculator"
	mimeTypeAudioMedia = "audio/mpeg"
)

func execByMime(mime string, pressed bool) error {
	if !pressed {
		return nil
	}

	cmd := queryCommand(mime)
	if len(cmd) == 0 {
		return fmt.Errorf("Not found executable for: %s", mime)
	}
	return doAction(cmd)
}

func queryCommand(mime string) string {
	if mime == mimeTypeCalc {
		return "gnome-calculator"
	}

	app := gio.AppInfoGetDefaultForType(mime, false)
	if app == nil {
		return ""
	}
	defer app.Unref()

	return app.GetExecutable()
}

func doAction(cmd string) error {
	logger.Debug("execute command: ", cmd)
	return exec.Command("/bin/sh", "-c", cmd).Run()
}
