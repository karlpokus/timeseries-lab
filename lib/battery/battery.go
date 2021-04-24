package battery

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	"timeseries/lib/telemetry"
)

type Battery struct {
	interval time.Duration
	name     string
}

// Run runs the data collector every specified interval to
// collect current battery charge percentage
func (bat *Battery) Run(c chan<- telemetry.Record) {
	bg := context.Background()
	for {
		time.Sleep(bat.interval) // kiss
		ctx, cancel := context.WithTimeout(bg, 3*time.Second)
		go func(ctx context.Context) {
			defer cancel()
			now, err := exec.CommandContext(ctx, "cat", "/sys/class/power_supply/BAT0/charge_now").Output()
			if err != nil {
				log.Printf("battery err: %s", err)
				return
			}
			full, err := exec.CommandContext(ctx, "cat", "/sys/class/power_supply/BAT0/charge_full").Output()
			if err != nil {
				log.Printf("battery err: %s", err)
				return
			}
			c <- telemetry.Record{
				Time:  time.Now(),
				Key:   bat.name,
				Value: telemetry.ByteToFloat(now) / telemetry.ByteToFloat(full) * 100,
			}
		}(ctx)
	}
}

func (bat *Battery) Info() string {
	return fmt.Sprintf("%s collector interval set to %s", bat.name, bat.interval)
}

func New(interval string) (*Battery, error) {
	d, err := time.ParseDuration(interval)
	if err != nil {
		return nil, err
	}
	return &Battery{
		interval: d,
		name:     "battery",
	}, nil
}
