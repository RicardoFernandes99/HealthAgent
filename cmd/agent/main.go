package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/RicardoFernandes99/health-agent/internal/collectors"
	"github.com/RicardoFernandes99/health-agent/internal/state"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	log.Println("Agent Running")
	store := state.NewStore()
	cpuColector := collectors.NewCPUCollector(5 * time.Second)

	go func() {
		ticker := time.NewTicker(cpuColector.Interval())
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				metric, err := cpuColector.Collect()
				if err != nil {
					log.Printf("Error collecting CPU metrics: %v", err)
					continue
				}
				store.Set(metric)

			}
		}
	}()
	go func() {
		ticker := time.NewTicker(7 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				metric, exists := store.Get("CPU_Usage")
				cputemp, tempExists := store.Get("CPU_Temperature")
				if tempExists {
					log.Printf("Stored CPU Temperature Metric: %.2f°C at %s", cputemp.Value, cputemp.Timestamp.Format("02/01/2006 15:04:05"))
					if !exists {
						log.Println("CPU metric not found")
						continue
					}
					log.Printf("Stored CPU Metric: %.2f%% at %s", metric.Value, metric.Timestamp.Format("02/01/2006 15:04:05"))
					log.Printf("CPU Temp : %.2f%% at %s\n", cputemp.Value, metric.Timestamp.Format("02/01/2006 15:04:05"))
				}
			}
		}
	}()
	<-ctx.Done()
	log.Println("Agent Stopped")
}
