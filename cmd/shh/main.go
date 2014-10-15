package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/heroku/shh"
)

var (
	signalChannel = make(chan os.Signal, 1)
	versionFlag = flag.Bool("version", false, "Display version info and exit")
)

func main() {
	flag.Parse()

	if *versionFlag {
		fmt.Println(shh.VERSION)
		os.Exit(0)
	}

	measurements := make(chan shh.Measurement, 100)
	config := shh.GetConfig()

	mp := shh.NewMultiPoller(measurements, config)

	signal.Notify(signalChannel, syscall.SIGINT)
	signal.Notify(signalChannel, syscall.SIGTERM)

	go func() {
		for sig := range signalChannel {
			mp.Exit()
			log.Fatal(shh.Slog{"signal": sig, "finished": time.Now(), "duration": time.Since(config.Start)})
		}
	}()

	if config.ProfilePort != shh.DEFAULT_PROFILE_PORT {
		go func() {
			log.Println(http.ListenAndServe("localhost:"+config.ProfilePort, nil))
		}()
	}

	ctx := shh.Slog{"shh_start": true, "at": config.Start.Format(time.RFC3339Nano), "interval": config.Interval}
	fmt.Println(ctx)

	outputter, err := shh.NewOutputter(config.Outputter, measurements, config)
	if err != nil {
		ctx.FatalError(err, "creating outputter")
	}
	outputter.Start()

	start := make(chan time.Time, 1)
	start <- time.Now()
	ticks := time.Tick(config.Interval)

	for {
		select {
		case tick := <-start:
			mp.Poll(tick)
			start = nil
		case tick := <-ticks:
			mp.Poll(tick)
		}
	}
}
