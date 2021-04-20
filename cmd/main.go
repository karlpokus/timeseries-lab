package main

import (
	"log"

	"timeseries"
	"timeseries/lib/battery"
)

func main() {
	bat, err := battery.New("5s")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("main start")
	timeseries.Start(bat)
}
