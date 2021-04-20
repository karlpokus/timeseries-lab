package heat

import (
  "log"
  "os/exec"
  "context"
  "time"
  "strings"

  "timeseries/lib/metric"
)

type Heat struct {
	interval time.Duration
}

// Run runs the data collector every specified interval to
// collect current cpu heat
func (heat *Heat) Run(c chan<- metric.Metric) {
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
			c <- metric.Metric{
				Key:   "heat",
				Value: strings.TrimSpace(string(out)), // TODO: remove +°C
			}
		}(ctx)
	}
}

func New(interval string) (*Heat, error) {
	d, err := time.ParseDuration(interval)
	if err != nil {
		return nil, err
	}
	return &Heat{
		interval: d,
	}, nil
}
