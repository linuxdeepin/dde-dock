package main

import (
	"bufio"
	"io/ioutil"
	"strings"
)

func (grub *Grub2) readEntries() {
	fileContent, err := ioutil.ReadFile(grub.grubMenuFile)
	if err != nil {
		logError(err.Error()) // TODO
		return
	}
	grub.entries = make([]string, 10)
	grub.parseEntries(string(fileContent))
}

func (grub *Grub2) readSettings() {
	fileContent, err := ioutil.ReadFile(grub.grubConfigFile)
	if err != nil {
		logError(err.Error()) // TODO
		return
	}
	grub.settings = make(map[string]string)
	grub.parseSettings(string(fileContent))
}

func (grub *Grub2) parseEntries(fileContent string) {
	// TODO
}

func (grub *Grub2) parseSettings(fileContent string) {
	// fileContent := []rune(fileContent) // TODO
	// br = bufio.NewReader(strings.NewReader(fileContent))
}
