/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
