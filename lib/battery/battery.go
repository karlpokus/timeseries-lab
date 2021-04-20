package battery

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"timeseries/lib/metric"
)

type Battery struct {
	interval time.Duration
}

// Run runs the data collector every specified interval to
// collect current battery charge percentage
func (bat *Battery) Run(c chan<- metric.Metric) {
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
			c <- metric.Metric{
				Key:   "battery",
				Value: fmt.Sprintf("%.0f", toFloat(now)/toFloat(full)*100),
			}
			log.Println("battery sent metric") // debug
		}(ctx)
	}
}

func toFloat(b []byte) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSpace(string(b)), 64)
	return f
}

func New(interval string) (*Battery, error) {
	d, err := time.ParseDuration(interval)
	if err != nil {
		return nil, err
	}
	return &Battery{
		interval: d,
	}, nil
}
