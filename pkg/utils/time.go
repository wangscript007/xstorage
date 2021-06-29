package utils

import "time"

func Unix2Str(timestamp int64, layout string, offset int) string {
	var cstzone = time.FixedZone("CST", offset)
	timeStr := time.Unix(timestamp, 0).In(cstzone).Format(layout)
	return timeStr
}
