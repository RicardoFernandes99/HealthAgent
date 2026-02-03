package collectors

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

type CPUCollector struct {
	interval time.Duration
}

func NewCPUCollector(interval time.Duration) *CPUCollector {
	return &CPUCollector{interval: interval}
}

func (c *CPUCollector) Name() string {
	return "CPU"
}

func (c *CPUCollector) Interval() time.Duration {
	return c.interval
}

func (c *CPUCollector) Collect() (Metric, error) {
	percent, err := cpu.Percent(0, false)
	if err != nil {
		return Metric{}, err
	}

	metric := Metric{
		Name:      c.Name(),
		Value:     percent[0],
		Timestamp: time.Now(),
	}
	log.Printf("CpuCollector : %.2f%%", metric.Value)

	return metric, nil
}
