package hog

import (
  "log"
  "os/exec"
  "context"
  "time"
  "strings"

  "timeseries/lib/metric"
)

type Hog struct {
	interval time.Duration
}

// Run runs the data collector every specified interval to
// collect the hog
func (hog *Hog) Run(c chan<- metric.Metric) {
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
    	data := strings.Fields(string(out))
    	if len(data) == 0 {
        log.Println("hog err: ps command output is empty")
				return
    	}
			c <- metric.Metric{
				Key:   "hog",
				Value: data[0], // data[1] is %cpu might be useful later
			}
		}(ctx)
	}
}

func New(interval string) (*Hog, error) {
	d, err := time.ParseDuration(interval)
	if err != nil {
		return nil, err
	}
	return &Hog{
		interval: d,
	}, nil
}
