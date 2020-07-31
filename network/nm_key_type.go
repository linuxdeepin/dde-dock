/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package network

type ipv4AddressesWrapper []ipv4AddressWrapper
type ipv4AddressWrapper struct {
	Address string
	Mask    string
	Gateway string
}

// Ipv6AddressesWrapper
type ipv6AddressesWrapper []ipv6AddressWrapper
type ipv6AddressWrapper struct {
	Address string
	Prefix  uint32
	Gateway string
}

// Ipv6Addresses is an array of (byte array, uint32, byte array)
type ipv6Addresses []ipv6Address
type ipv6Address struct {
	Address []byte
	Prefix  uint32
	Gateway []byte
}

// ipv6Routes is an array of (byte array, uint32, byte array, uint32)
type ipv6Route struct {
	Address []byte
	Prefix  uint32
	NextHop []byte
	Metric  uint32
}
type ipv6Routes []ipv6Route
