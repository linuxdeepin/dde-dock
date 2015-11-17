/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
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
