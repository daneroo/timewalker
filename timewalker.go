package timewalker

import (
	"fmt"
	"time"
)

type Interval struct {
	Start time.Time
	End   time.Time
}

func TryAdd() {
	now := time.Now()
	fmt.Printf("now:  %v\n", now)
	soon := now.Add(time.Second)
	fmt.Printf("now:  %v\n", now)
	fmt.Printf("soon: %v\n", soon)

}

func (i Interval) String() string {
	// layout := time.RFC3339
	layout := "2006-01-02T15:04:05.000Z07:00"
	return fmt.Sprintf("[%s, %s)", i.Start.Format(layout), i.End.Format(layout))
}

// Truncate Start, Round up End
func (i Interval) Round(d time.Duration) Interval {
	fmt.Printf("-%v\n", i)
	start := i.Start.Truncate(d)
	end := i.End.Truncate(d)
	fmt.Printf("+%v\n", i)
	fmt.Printf("=%v\n", Interval{Start: start, End: end})
	return Interval{Start: start, End: end}
}

func (i Interval) Iter(d time.Duration) (<-chan Interval, error) {
	ch := make(chan Interval)

	// truncate start to duration
	ri := i.Round(d)

	fmt.Printf("\n")
	fmt.Printf("%s\n", time.RFC3339)
	fmt.Printf(" i: %v\n", i)
	fmt.Printf("ri: %v\n", ri)
	go func() {
		start := ri.Start
		for start.Before(ri.End) {
			end := start.Add(d)
			ch <- Interval{Start: start, End: end}
			start = end
		}
		close(ch)
	}()
	return ch, nil
}
