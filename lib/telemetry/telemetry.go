package telemetry

import (
	"strconv"
	"strings"
	"time"
)

type Record struct {
	Time  time.Time
	Key   string
	Value float64
}

func ByteToFloat(b []byte) float64 {
	return StringToFloat(string(b))
}

func StringToFloat(s string) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f
}

func Datapoints(rcds []Record, key string) (out [][]interface{}) {
	for _, rcd := range rcds {
		if rcd.Key == key {
			x := make([]interface{}, 2)
			x[0] = rcd.Value
			x[1] = unixMs(rcd.Time)
			out = append(out, x)
		}
	}
	return
}

func unixMs(t time.Time) int64 {
	return t.UnixNano() / 1e6
}
