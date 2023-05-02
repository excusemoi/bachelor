package kafka

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type Metrics struct {
	latency              *Metric
	lag                  *Metric
	inputMessagesPerSec  *Metric
	outputMessagesPerSec *Metric
}

type Metric struct {
	count    bool
	value    int64
	duration time.Duration
	m        prometheus.Gauge
}

func (mps *Metric) Observe(handler func()) {
	if mps.duration != 0 {
		for {
			select {
			case <-time.After(mps.duration):
				{
					handler()
				}
			}
		}
	}
}

func (mps *Metric) messagePerSecondHandler() {
	mps.m.Set(float64(mps.value) / mps.duration.Seconds())
	mps.value = 0
}
