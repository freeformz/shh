package outputters

import (
  "fmt"
  "github.com/freeformz/shh/mm"
)

type Outputter interface {
	Poll(measurements <-chan *mm.Measurement)
}

type L2MetStdOut struct {}

func (out L2MetStdOut) Output (measurements <-chan *mm.Measurement) {
	for measurement := range measurements {
		fmt.Println(measurement)
	}
}
