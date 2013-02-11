package pollers

import (
	"github.com/freeformz/shh/mm"
	"github.com/freeformz/shh/utils"
	"sync"
	"time"
)

const (
	DEFAULT_POLLERS = "load,cpu,df,disk,nif,mem,ntpdate" // Default pollers
)

var (
	pollers = utils.GetEnvWithDefaultStrings("SHH_POLLERS", DEFAULT_POLLERS)
)

type Poller interface {
	Name() string
	Poll(tick time.Time)
}

func NewMultiPoller(measurements chan<- *mm.Measurement) Multi {
	mp := Multi{pollers: make(map[string]Poller), measurements: measurements, counts: make(map[string]uint64)}

	for _, poller := range pollers {
		switch poller {
		case "load":
			mp.RegisterPoller(NewLoadPoller(measurements))
		case "cpu":
			mp.RegisterPoller(NewCpuPoller(measurements))
		case "df":
			mp.RegisterPoller(NewDfPoller(measurements))
		case "disk":
			mp.RegisterPoller(NewDiskPoller(measurements))
		case "nif":
			mp.RegisterPoller(NewNetworkInterfacePoller(measurements))
		case "mem":
			mp.RegisterPoller(NewMemoryPoller(measurements))
		case "ntpdate":
			mp.RegisterPoller(NewNtpdatePoller(measurements))
		}
	}

	return mp
}

type Multi struct {
	sync.WaitGroup
	measurements chan<- *mm.Measurement
	pollers      map[string]Poller
	counts       map[string]uint64
}

func (p Multi) RegisterPoller(poller Poller) {
	p.pollers[poller.Name()] = poller
	p.counts[poller.Name()] = 0
}

func (p Multi) Poll(tick time.Time) {
	for name, poller := range p.pollers {
		count := p.counts[name]
		count += 1
		p.counts[name] = count
		p.measurements <- &mm.Measurement{tick, p.Name(), []string{"ticks", name, "count"}, count}
		p.Add(1)
		go func(poller Poller) {
			poller.Poll(tick)
			p.Done()
		}(poller)
	}
	p.Wait()
}

func (p Multi) Name() string {
	return "multi_poller"
}
