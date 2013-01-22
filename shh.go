package main

import (
	"fmt"
	"github.com/freeformz/shh/utils"
	"github.com/freeformz/shh/mm"
	"github.com/freeformz/shh/pollers"
  "github.com/freeformz/shh/outputters"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	DEFAULT_INTERVAL = "10s" // Default tick interval for pollers
)

var (
  Interval = utils.GetEnvWithDefault("SHH_INTERVAL", DEFAULT_INTERVAL) // Polling Interval
	Start = time.Now()                                                   // Start time
)

// Get a time.Duration for the Interval
func getDuration() time.Duration {
	duration, err := time.ParseDuration(Interval)

	if err != nil {
		log.Fatal("unable to parse $SHH_INTERVAL: " + Interval)
	}

	return duration
}

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		for sig := range c {
			fmt.Printf("signal=%s finished=%s duration=%s\n", sig, time.Now().Format(time.RFC3339Nano), time.Since(Start))
			os.Exit(1)
		}
	}()
}

func main() {
	duration := getDuration()
	fmt.Printf("shh_start=true at=%s interval=%s\n", Start.Format(time.RFC3339Nano), duration)

	mp := pollers.NewMultiPoller()
	mp.RegisterPoller(pollers.Load{})
	mp.RegisterPoller(pollers.Cpu{})
	mp.RegisterPoller(pollers.Df{})
	mp.RegisterPoller(pollers.Disk{})

	measurements := make(chan *mm.Measurement, 100)
	//go outputters.L2MetStdOut{}.Output(measurements)
	go outputters.Librato{}.Output(measurements)

	// do a tick at start
	go mp.Poll(measurements)

	ticks := time.Tick(duration)
	for {
		select {
		case <-ticks:
			go mp.Poll(measurements)
		}
	}
}
