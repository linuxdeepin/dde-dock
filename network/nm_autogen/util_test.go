package main

import (
	"testing"
)

func TestToClassName(t *testing.T) {
	tests := []struct {
		str, result string
	}{
		{"vk-key-mgmt", "VkKeyMgmt"},
		{"vk key mgmt", "VkKeyMgmt"},
		{"vk_key_mgmt", "VkKeyMgmt"},
	}
	for _, test := range tests {
		if test.result != ToClassName(test.str) {
			t.Fail()
		}
	}
}
