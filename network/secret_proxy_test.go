/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

import (
	C "gopkg.in/check.v1"
)

func (*testWrapper) TestSecretProxyType(c *C.C) {
	var proxy = new(secretProxyType)
	proxy.Add(1)
	c.Check(len(*proxy), C.Equals, 1)
	c.Check(proxy.Last(), C.Equals, uint32(1))
	proxy.Add(1)
	c.Check(len(*proxy), C.Equals, 1)
	proxy.delete(1)
	c.Check(len(*proxy), C.Equals, 0)
}
