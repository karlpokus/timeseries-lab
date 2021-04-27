package timeseries

import (
	"log"
	"time"

	"timeseries/lib/db"
	"timeseries/lib/telemetry"
)

type Collector interface {
	Run(c chan<- telemetry.Record)
	Info() string
}

func Start(collectors ...Collector) error {
	url := "postgres://postgres:secret@localhost:5432/test"
	pool, err := db.Connect(url)
	if err != nil {
		return err
	}
	log.Println("db connected")
	defer pool.Close()
	c := make(chan telemetry.Record)
	for _, collector := range collectors {
		log.Printf(collector.Info())
		go collector.Run(c)
		jitter()
	}
	for rcd := range c {
		err = db.Insert(pool, rcd)
		if err != nil {
			// just logging here is a temporary solution.
			// broken db connections will print an ugly error
			// and then retry with another proper connection from the pool
			log.Printf("db insert err: %s", err)
		}
	}
	return nil
}

func jitter() {
	time.Sleep(3 * time.Second)
}
