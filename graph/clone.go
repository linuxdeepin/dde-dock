/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package graph

// Clone clones(deep copy) the graph Data. (changing the cloned Data would not affect the original Data.)
// It traverses every single node with depth-first-search.
func (d *Data) Clone() *Data {
	clonedData := New()
	src := &Node{}
	for nd := range d.NodeMap {
		src = nd
		break
	}
	d.cloneDfs(src, clonedData)
	return clonedData
}

func (d *Data) cloneDfs(src *Node, clonedData *Data) {
	if src.Color == "black" {
		return
	}

	src.Color = "black"
	srcClone := NewNode(src.ID)

	for ov, weight := range src.WeightTo {
		ovClone := NewNode(ov.ID)
		clonedData.Connect(srcClone, ovClone, weight)
		if ov.Color == "white" {
			d.cloneDfs(ov, clonedData)
		}
	}

	for iv, weight := range src.WeightFrom {
		ivClone := NewNode(iv.ID)
		clonedData.Connect(ivClone, srcClone, weight)
		if iv.Color == "white" {
			d.cloneDfs(iv, clonedData)
		}
	}
}
