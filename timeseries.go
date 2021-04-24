package timeseries

import (
	"log"
	"context"

	"timeseries/lib/db"
	"timeseries/lib/telemetry"
)

type Collector interface {
	Run(c chan<- telemetry.Record)
	Info() string
}

func Start(collectors ...Collector) error {
	url := "postgres://postgres:secret@localhost:5432/test"
	conn, err := db.Connect(url)
	if err != nil {
		return err
	}
	log.Println("db connected")
	defer conn.Close(context.Background())
	c := make(chan telemetry.Record)
	for _, collector := range collectors {
		log.Printf(collector.Info())
		go collector.Run(c)
	}
	for rcd := range c {
		err = db.Insert(conn, rcd)
		if err != nil {
			log.Printf("db insert err: %s", err)
		}
	}
	return nil
}
