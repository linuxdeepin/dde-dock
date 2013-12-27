package main 

import (
	"io/ioutil"
	"regexp"
	"strconv"
)

const (
	_SIMPLE_CONFIG_PATH = "/etc/default/grub"
	_TIMEOUT_REGEXP = "^GRUB_TIMEOUT=([0-9]+)"
)

type SimpleParser struct {
	timeout int32
	resolution string
}

func NewSimpleParser() *SimpleParser {
	return &SimpleParser {}
}

func (sp *SimpleParser) Parse() (error) {
	file_content, err := ioutil.ReadFile(_SIMPLE_CONFIG_PATH)
	if (err != nil) {
		return err
	}
	timeout_pattern := regexp.MustCompile(_TIMEOUT_REGEXP)
	matches := timeout_pattern.FindStringSubmatch(string(file_content))
	if len(matches) == 0 {
		timeout, err := strconv.Atoi(matches[0])
		sp.timeout = int32(timeout)
		if (err != nil) {
			return err
		}
	}
	
	return nil
}