package main

import (
	"log"
	"time"

	"timeseries/lib/telemetry"
	"timeseries/lib/store"
	"timeseries/lib/battery"
	"timeseries/lib/cpu/heat"
	"timeseries/lib/cpu/hog"
)

var version string

func main() {
	log.SetFlags(0)
	st, err := store.New("postgres://postgres:secret@localhost:5432/test")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db connected")
	defer st.Close()
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
	log.Println("collectors initialized")
	log.Printf("starting agent version %s", version)
	Start(st, bat, cpuHeat, cpuHog)
	log.Println("telemetry agent exit")
}

func Start(st store.Store, collectors ...telemetry.Collector) {
	c := make(chan telemetry.Record)
	for _, collector := range collectors {
		log.Printf(collector.Info())
		go collector.Run(c)
		jitter()
	}
	for rcd := range c {
		err := st.Insert(rcd)
		if err != nil {
			// just logging here is a temporary solution.
			// broken db connections will print an ugly error
			// and then retry with another proper connection from the pool
			log.Printf("db insert err: %s", err)
		}
	}
}

func jitter() {
	time.Sleep(3 * time.Second)
}
