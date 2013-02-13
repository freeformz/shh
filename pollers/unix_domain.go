package pollers

import (
	"github.com/freeformz/shh/mm"
	"github.com/freeformz/shh/utils"
	"log"
  "fmt"
	"net"
	"sync"
	"time"
  "strings"
)

type ListenStats struct {
	sync.RWMutex
	connectionCount float64
}

func (ls *ListenStats) IncrementConnectionCount() {
    ls.Lock()
    fmt.Println(ls.connectionCount)
    ls.connectionCount++
    ls.Unlock()
}

func (ls *ListenStats) DecrementConnectionCount() {
    ls.Lock()
    fmt.Println(ls.connectionCount)
    ls.connectionCount--
    ls.Unlock()
}

func (ls *ListenStats) ConnectionCount() float64 {
	defer ls.RUnlock()
	ls.RLock()
	return ls.connectionCount
}

type Listen struct {
	measurements    chan<- *mm.Measurement
	listener        net.Listener
	stats           *ListenStats
}

var (
	listen = utils.GetEnvWithDefault("SHH_LISTEN", "unix,/tmp/shh")
  listenNet string
  listenLaddr string
)

func init() {
  tmp := strings.Split(listen,",")

	if len(tmp) != 2 {
		log.Fatal("SHH_LISTEN is not in the format: 'unix,/tmp/shh'")
	}

  listenNet = tmp[0]
  listenLaddr = tmp[1]

	switch listenNet{
	case "tcp", "tcp4", "tcp6", "unix", "unixpacket":
		break
	default:
		log.Fatalf("SHH_LISTEN format (%s,%s) is not correct", listenNet, listenLaddr)
	}

}

func NewListenPoller(measurements chan<- *mm.Measurement) Listen {
	listener, err := net.Listen(listenNet, listenLaddr)

	if err != nil {
		log.Fatal(err)
	}

	poller := Listen{measurements: measurements, listener: listener, stats: &ListenStats{}}

	go func(poller *Listen) {
		for {
			conn, err := poller.listener.Accept()
			if err != nil {
				log.Println(err)
				continue
			}

			go poller.handleConnection(conn)
		}
	}(&poller)

	return poller
}

func (poller Listen) Poll(tick time.Time) {
	poller.measurements <- &mm.Measurement{tick, poller.Name(), []string{"connection", "count"}, poller.stats.ConnectionCount()}
}

func (poller Listen) handleConnection(conn net.Conn) {
	defer func() {
    conn.Close()
		poller.stats.DecrementConnectionCount()
  }()

	poller.stats.IncrementConnectionCount()

	time.Sleep(time.Second * 20)

}

func (poller Listen) Name() string {
	return "listen"
}

func (poller Listen) Exit() {
	poller.listener.Close()
}
