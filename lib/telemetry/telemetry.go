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
