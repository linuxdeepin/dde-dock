package main

import (
	"dlib/glib-2.0"
	"os"
	p "path"
	"testing"
)

func TestIsOnDesktop(t *testing.T) {
	isOnDesktop("firefox.desktop")
}

func TestSendToDesktop(t *testing.T) {
	target := "/usr/share/applications/firefox.desktop"
	sendToDesktop(target)
	path :=
		p.Join(glib.GetUserSpecialDir(glib.UserDirectoryDirectoryDesktop),
			p.Base(target))
	state, err := os.Stat(path)
	if err != nil {
		t.Log(err)
		return
	}

	if state.Mode().Perm()&0100 == 0 {
		t.Error("Permision failed")
	}
}
