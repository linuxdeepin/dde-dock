/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/
package screenedge

type areaRange struct {
	X1 int32
	Y1 int32
	X2 int32
	Y2 int32
}

func (area *areaRange) Contains(x, y int32) bool {
	if x >= (area.X1) && x < (area.X2) &&
		y >= (area.Y1) && y < (area.Y2) {
		return true
	}
	return false
}

func (area *areaRange) GetCornerSquares(sideLength int32) []*areaRange {
	var (
		startX = area.X1
		startY = area.Y1
		endX   = area.X2
		endY   = area.Y2
	)

	return []*areaRange{
		// TopLeft
		&areaRange{
			startX, startY, startX + sideLength, startY + sideLength,
		},
		// TopRight
		&areaRange{
			endX - sideLength, startY, endX, startY + sideLength,
		},
		// BottomRight
		&areaRange{
			endX - sideLength, endY - sideLength, endX, endY,
		},
		// BottomLeft
		&areaRange{
			startX, endY - sideLength, startX + sideLength, endY,
		},
	}
}
