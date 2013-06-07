package utils

import (
  "time"
  "runtime"
)

type runtimeTimer struct {
	i      int32
	when   int64
	period int64
	f      func(int64, interface{}) // NOTE: must not be closure
	arg    interface{}
}

func startTimer(*runtimeTimer)

func sendTime(now int64, c interface{}) {
	// Non-blocking send of time on c.
	// Used in NewTimer, it cannot block anyway (buffer).
	// Used in NewTicker, dropping sends on the floor is
	// the desired behavior when the reader gets behind,
	// because the sends are periodic.
	select {
	case c.(chan time.Time) <- time.Unix(0, now):
	default:
	}
}

func NewPreTicker(d time.Duration) *time.Ticker {
	if d <= 0 {
		panic(errors.New("non-positive interval for NewTicker"))
	}
	// Give the channel a 1-element time buffer.
	// If the client falls behind while reading, we drop ticks
	// on the floor until the client catches up.
	c := make(chan time.Time, 1)
	t := &time.Ticker{
		C: c,
		r: time.runtimeTimer{
			when:   time.UnixNano(),
			period: int64(d),
			f:      sendTime,
			arg:    c,
		},
	}
	startTimer(&t.r)
	return t
}

func PreTick(d time.Duration) <-chan time.Time {
  if d <= 0 {
    return nil
  }
  return NewPreTicker(d).C
}


