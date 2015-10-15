package timewalker

import (
	"fmt"
	"time"
)

// Human Readable Durations (non-arithmetic)
type HDuration int

const (
	Day HDuration = iota
	Month
	Year
)

func Round(t time.Time, hd HDuration) time.Time {
	year, month, day := t.Date()
	switch hd {
	case Day:
		return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	case Month:
		return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
	case Year:
		return time.Date(year, time.January, 1, 0, 0, 0, 0, t.Location())
	}
	return t

}
func (hd HDuration) String() string {
	switch hd {
	case Day:
		return "Day"
	case Month:
		return "Month"
	case Year:
		return "Year"
	}
	return ""
}

type Interval struct {
	Start time.Time
	End   time.Time
}

func (i Interval) String() string {
	// layout := time.RFC3339
	layout := "2006-01-02T15:04:05.000Z07:00"
	return fmt.Sprintf("[%s, %s)", i.Start.Format(layout), i.End.Format(layout))
}

// This normalizes an Interval's representations
// -swap Start,End if appropriate
// Truncate Start, Round up End
// Make sure we have at least one interval.
func (i Interval) Round(d time.Duration) Interval {
	if i.End.Before(i.Start) {
		i.End, i.Start = i.Start, i.End
	}
	i.Start = i.Start.Truncate(d)
	// truncate End, if truncEnd<End, or truncEnd<Start, add 1 duration
	truncEnd := i.End.Truncate(d)
	if truncEnd.Before(i.End) {
		i.End = truncEnd.Add(d)
	} else if truncEnd.Before(i.Start.Add(d)) {
		i.End = i.Start.Add(d)
	} else {
		i.End = truncEnd
	}
	return i
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
