package timeseries

import (
	"log"

	"timeseries/lib/metric"
)

type Collector interface {
	Run(c chan<- metric.Metric)
}

func Start(collectors ...Collector) {
	c := make(chan metric.Metric)
	for _, collector := range collectors {
		go collector.Run(c)
	}
	for x := range c {
		log.Printf("%s reported %s", x.Key, x.Value)
	}
}
