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

type Collector interface {
	Run(c chan<- Record)
	Info() string
}

type Range struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type Target struct {
	Data   interface{} `json:"data"` // empty string or {key: "", operator: "=", value: x}
	Target string      `json:"target"`
	Type   string      `json:"type"`
}

type Query struct {
	Range         `json:"range"`
	Targets       []Target `json:"targets"`
	AdhocFilters  []string `json:"adhocFilters"`
	MaxDataPoints int      `json:"maxDataPoints"`
}

type Search struct {
	Type   string `json:"type"`
	Target string `json:"target"`
}

type Response struct {
	Target     string          `json:"target"`
	Datapoints [][]interface{} `json:"datapoints"`
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
