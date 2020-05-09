package power

import "math"

/* 记录数量 */
func (bat *Battery) getHistoryLength() int {
	return len(bat.batteryHistory)
}

/* 插入记录 */
func (bat *Battery) appendToHistory(percentage float64) {
	bat.batteryHistory = append(bat.batteryHistory, percentage)
	if len(bat.batteryHistory) > 10 {
		bat.batteryHistory = bat.batteryHistory[1:]
	}
}

/* 计算方差 */
func (bat *Battery) calcHistoryVariance() float64 {
	var average float64 = 0.0
	for i := range bat.batteryHistory {
		average += bat.batteryHistory[i]
	}
	average /= float64(len(bat.batteryHistory))

	var variance = 0.0
	for i := range bat.batteryHistory {
		variance += math.Pow(bat.batteryHistory[i]-average, 2)
	}
	variance /= float64(len(bat.batteryHistory))

	return variance
}
