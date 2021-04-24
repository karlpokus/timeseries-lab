package hog

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"timeseries/lib/telemetry"
)

type Hog struct {
	interval time.Duration
	name     string
}

// Run runs the data collector every specified interval to
// collect the hog
func (hog *Hog) Run(c chan<- telemetry.Record) {
	bg := context.Background()
	for {
		time.Sleep(hog.interval) // kiss
		ctx, cancel := context.WithTimeout(bg, 3*time.Second)
		go func(ctx context.Context) {
			defer cancel()
			cmd := "ps -Ao comm,pcpu --sort=-pcpu --no-headers | head -n 1"
			out, err := exec.CommandContext(ctx, "bash", "-c", cmd).Output()
			if err != nil {
				log.Printf("hog err: %s", err)
				return
			}
			data := strings.Fields(string(out)) // command pcpu
			if len(data) == 0 {
				log.Println("hog err: ps command output is empty")
				return
			}
			c <- telemetry.Record{
				Time:  time.Now(),
				Key:   hog.name,
				Value: telemetry.StringToFloat(data[1]),
			}
		}(ctx)
	}
}

func (hog *Hog) Info() string {
	return fmt.Sprintf("%s collector interval set to %s", hog.name, hog.interval)
}

func New(interval string) (*Hog, error) {
	d, err := time.ParseDuration(interval)
	if err != nil {
		return nil, err
	}
	return &Hog{
		interval: d,
		name:     "hog",
	}, nil
}
