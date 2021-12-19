package discordgo

import (
	"fmt"
	"time"
)

type sleepCT struct {
	d     time.Duration // desired duration between targets
	t     time.Time     // last time target
	wake  time.Time     // last wake time
	drift int64         // last wake drift microseconds
}

func NewSleepCT(d time.Duration) sleepCT {

	s := sleepCT{}

	s.d = d
	s.t = time.Now()

	return s
}

func (s *sleepCT) SleepNext() int64 {

	now := time.Now()

	// if target is zero safety net
	if s.t.IsZero() {
		fmt.Println("TickerCT reset")
		s.t = now.Add(-s.d)
	}

	// Move forward the sleep target by the duration
	s.t = s.t.Add(s.d)

	// Compute the desired sleep time to reach the target
	d := time.Until(s.t)

	// Sleep
	time.Sleep(d)

	// record the wake time
	s.wake = time.Now()
	s.drift = s.wake.Sub(s.t).Microseconds()

	// fmt.Println(s.t.UnixMilli(), d.Milliseconds(), wake.UnixMilli(), drift, pause, len(s.resume))

	// return the drift for monitoring purposes
	return s.drift
}
