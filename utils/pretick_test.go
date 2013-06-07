package utils_test

import (
  "testing"
  . "github.com/freeformz/shh/utils"
)

func TestTicker(t *testing.T) {
	const Count = 10
	Delta := 100 * Millisecond
	ticker := NewPreTicker(Delta)
	t0 := Now()
	for i := 0; i < Count; i++ {
		<-ticker.C
	}
	ticker.Stop()
	t1 := Now()
	dt := t1.Sub(t0)
	target := Delta * Count
	slop := target * 2 / 10
	if dt < target-slop || (!testing.Short() && dt > target+slop) {
		t.Fatalf("%d %s ticks took %s, expected [%s,%s]", Count, Delta, dt, target-slop, target+slop)
	}
	// Now test that the ticker stopped
	Sleep(2 * Delta)
	select {
	case <-ticker.C:
		t.Fatal("Ticker did not shut down")
	default:
		// ok
	}
}
