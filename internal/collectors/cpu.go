package collectors

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
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

func (c *CPUCollector) Collect() ([]Metric, error) {
	metrics := []Metric{}
	percent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	metrics = append(metrics, Metric{
		Name:      "CPU_Usage",
		Value:     percent[0],
		Timestamp: time.Now(),
	})

	temps, err := host.SensorsTemperatures()
	log.Println("temps:", temps)

	if err != nil {
		log.Println("Temp not found")
		return nil, err
	}

	for _, temperature := range temps {
		metrics = append(metrics, Metric{
			Name:      "CPU_Temperature",
			Value:     temperature.Temperature,
			Timestamp: time.Now(),
		})

	}

	log.Println(metrics)
	log.Printf("CpuCollector : %.2f%%", metrics[0].Value)

	return metrics, nil
}
