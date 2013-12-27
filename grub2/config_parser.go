package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
)

const (
	_SIMPLE_CONFIG_PATH = "/etc/default/grub"
	_TIMEOUT_REGEXP     = "GRUB_TIMEOUT=([0-9]+)"
	_RESOLUTION_REGEXP  = "GRUB_GFXMODE=([0-9]+)x([0-9]+)"
	_THEME_REGEXP       = "GRUB_THEME=/boot/grub/themes/(.*?)/theme.txt"
)

type SimpleParser struct {
	timeout    int32
	resolution string
	theme string
}

func NewSimpleParser() *SimpleParser {
	return &SimpleParser{-1, "", ""}
}

func (sp *SimpleParser) Parse() error {
	file_content, err := ioutil.ReadFile(_SIMPLE_CONFIG_PATH)
	if err != nil {
		return err
	}

	timeout_pattern := regexp.MustCompile(_TIMEOUT_REGEXP)
	timeout_matches := timeout_pattern.FindStringSubmatch(string(file_content))
	if len(timeout_matches) != 0 {
		timeout, err := strconv.Atoi(timeout_matches[1])
		if err != nil {
			return err
		}
		sp.timeout = int32(timeout)
	}

	resolution_pattern := regexp.MustCompile(_RESOLUTION_REGEXP)
	resolution_matches := resolution_pattern.FindStringSubmatch(string(file_content))
	if len(resolution_matches) != 0 {
		res_width := resolution_matches[1]
		res_height := resolution_matches[2]
		sp.resolution = fmt.Sprintf("%sx%s", res_width, res_height)
	}
	
	theme_pattern := regexp.MustCompile(_THEME_REGEXP)
	theme_matches := theme_pattern.FindStringSubmatch(string(file_content))
	if len(theme_matches) != 0 {
		sp.theme = theme_matches[1]
	}

	return nil
}
