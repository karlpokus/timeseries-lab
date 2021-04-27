package main

import (
	"log"

	"timeseries"
	"timeseries/lib/battery"
	"timeseries/lib/cpu/heat"
	"timeseries/lib/cpu/hog"
)

var version string

func main() {
	log.SetFlags(0)
	bat, err := battery.New("60s")
	if err != nil {
		log.Fatal(err)
	}
	cpuHeat, err := heat.New("60s")
	if err != nil {
		log.Fatal(err)
	}
	cpuHog, err := hog.New("60s")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("telemetry version %s start", version)
	err = timeseries.Start(bat, cpuHeat, cpuHog)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("telemetry exit")
}
