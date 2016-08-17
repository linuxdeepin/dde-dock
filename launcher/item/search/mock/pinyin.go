/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package mock

type PinYin struct {
	data  map[string][]string
	valid bool
}

func (self *PinYin) Search(key string) ([]string, error) {
	return self.data[key], nil
}

func (self *PinYin) Update(data []string) error {
	return nil
}

func (self *PinYin) IsValid() bool {
	return self.valid
}

func NewPinYin(data map[string][]string, valid bool) *PinYin {
	return &PinYin{
		data:  data,
		valid: valid,
	}
}
