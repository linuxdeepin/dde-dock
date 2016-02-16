/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package search

// ResultList is type alias for []Result used for sort.
type ResultList []Result

func (self ResultList) Len() int {
	return len(self)
}

func (self ResultList) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self ResultList) Less(i, j int) bool {
	if self[i].Score > self[j].Score {
		return true
	}

	if self[j].Score > self[i].Score {
		return false
	}

	if self[i].Freq > self[j].Freq {
		return true
	}

	if self[j].Freq > self[i].Freq {
		return false
	}

	return self[i].Name < self[j].Name
}
