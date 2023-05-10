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
	filtrationParams     *Metric
}

type Metric struct {
	count bool
	value int64
	start time.Time
	m     prometheus.Gauge
}

func (mps *Metric) messagePerSecondHandler() {
	//mps.m.Set(float64(mps.value) / mps.start.Seconds())
}
