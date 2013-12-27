package main 

import (
	"io/ioutil"
	"regexp"
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

func (sp *SimpleParser) Parse(path string) (*SimpleParser, error) {
	file_content, err := ioutil.ReadFile(path)
	if (err != nil) {
		return sp, err
	}
	timeout_pattern := regexp.MustCompile(_TIMEOUT_REGEX)
	matches := timeout_pattern.FindStringSubmatch(file_content)
	if matches.len = 0 {
		sp.timeout = matches[0]
	}
	
	return sp
}