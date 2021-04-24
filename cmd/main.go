package main

import (
	"log"

	"timeseries"
	"timeseries/lib/battery"
	"timeseries/lib/cpu/heat"
	"timeseries/lib/cpu/hog"
)

func main() {
	bat, err := battery.New("10s")
	if err != nil {
		log.Fatal(err)
	}
	cpuHeat, err := heat.New("5s")
	if err != nil {
		log.Fatal(err)
	}
	cpuHog, err := hog.New("7s")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("main start")
	err = timeseries.Start(bat, cpuHeat, cpuHog)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("main exit")
}
