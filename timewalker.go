package timewalker

import (
	"time"
)

type Interval struct {
	Start time.Time
	End   time.Time
}

func TimeWalker() (<-chan Interval, error) {
	intervals := make(chan Interval)

	go func() {
		start := time.Now()
		for i := 0; i < 10; i++ {
			end := start.Add(time.Second)
			intervals <- Interval{Start: start, End: end}
			start = end
		}
		close(intervals)
	}()

	return intervals, nil
}
