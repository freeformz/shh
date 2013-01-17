package pollers

import (
	"bufio"
	"github.com/freeformz/shh/mm"
	"log"
	"os"
	"strings"
)

type Load struct{}

func (poller Load) Poll(measurements chan<- *mm.Measurement) {
	file, err := os.Open("/proc/loadavg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fields := strings.Fields(line)
	measurements <- &mm.Measurement{poller.Name(), []string{"1m"}, fields[0]}
	measurements <- &mm.Measurement{poller.Name(), []string{"5m"}, fields[1]}
	measurements <- &mm.Measurement{poller.Name(), []string{"15m"}, fields[2]}
	entities := strings.Split(fields[3], "/")
	measurements <- &mm.Measurement{poller.Name(), []string{"scheduling", "entities", "executing"}, entities[0]}
	measurements <- &mm.Measurement{poller.Name(), []string{"scheduling", "entities", "total"}, entities[1]}
	measurements <- &mm.Measurement{poller.Name(), []string{"pid", "last"}, fields[4]}
}

func (poller Load) Name() string {
	return "load"
}
