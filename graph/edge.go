/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package graph

// Edge connects from Src to Dst with weight.
type Edge struct {
	Src    *Node
	Dst    *Node
	Weight float32
}

// GetEdges returns all edges of a graph.
func (d *Data) GetEdges() []Edge {
	rs := []Edge{}
	for nd1 := range d.NodeMap {
		for nd2, v := range nd1.WeightTo {
			one := Edge{}
			one.Src = nd1
			one.Dst = nd2
			one.Weight = v
			rs = append(rs, one)
		}
	}
	return rs
}
