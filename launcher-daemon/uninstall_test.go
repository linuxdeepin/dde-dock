package main

import (
	"strings"
	"testing"
)

func TestGetPackageNamesFromDataBase(t *testing.T) {
	names, err := getPackageNamesFromDatabase("/usr/share/applications/firefox.desktop")
	if err != nil {
		t.Error(err)
	}

	for _, name := range names {
		if strings.Contains(name, "firefox") {
			return
		}
	}
	t.Error("Get Wrong package names", names)
}
func TestGetPackageNamesFromCommandline(t *testing.T) {
	names, err :=
		getPackageNamesFromCommandline("/usr/share/applications/firefox.desktop")
	if err != nil {
		t.Error(err)
	}

	for _, name := range names {
		if strings.Contains(name, "firefox") {
			return
		}
	}
	t.Error("Get Wrong package names", names)
}
