package heat

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"timeseries/lib/telemetry"
)

type Heat struct {
	interval time.Duration
	name     string
}

// Run runs the data collector every specified interval to
// collect current cpu heat
func (heat *Heat) Run(c chan<- telemetry.Record) {
	bg := context.Background()
	for {
		time.Sleep(heat.interval) // kiss
		ctx, cancel := context.WithTimeout(bg, 3*time.Second)
		go func(ctx context.Context) {
			defer cancel()
			cmd := "/usr/bin/sensors | grep 'Package id 0' | cut -d ' ' -f 5"
			out, err := exec.CommandContext(ctx, "bash", "-c", cmd).Output()
			if err != nil {
				log.Printf("heat err: %s", err)
				return
			}
			c <- telemetry.Record{
				Time:  time.Now(),
				Key:   heat.name,
				Value: telemetry.StringToFloat(trim(string(out))),
			}
		}(ctx)
	}
}

func (heat *Heat) Info() string {
	return fmt.Sprintf("%s collector interval set to %s", heat.name, heat.interval)
}

func trim(s string) string {
	return strings.Trim(strings.TrimSpace(s), "+°C")
}

func New(interval string) (*Heat, error) {
	d, err := time.ParseDuration(interval)
	if err != nil {
		return nil, err
	}
	return &Heat{
		interval: d,
		name:     "heat",
	}, nil
}
