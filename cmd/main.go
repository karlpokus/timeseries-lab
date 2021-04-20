package main

import (
	"log"

	"timeseries"
	"timeseries/lib/battery"
	"timeseries/lib/cpu/heat"
)

func main() {
	bat, err := battery.New("5s")
	if err != nil {
		log.Fatal(err)
	}
	cpuHeat, err := heat.New("7s")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("main start")
	timeseries.Start(bat, cpuHeat)
}
