package util

import "time"

func GetTimeByConverFloat64(t float64)time.Time{
	i := int64(t)
	return time.Unix(i,0)
}