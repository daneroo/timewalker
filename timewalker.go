// go package to generate sequential time intervals
// for days,months, and years, accounting for Timezones.
package timewalker

import (
	"fmt"
	"time"
)

// Human Readable Durations (non-arithmetic)
type HDuration int

// Different pakage constant defineing an enu typr fot HDuration
const (
	Day HDuration = iota
	Month
	Year
)

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

func Add(t time.Time, hd HDuration) time.Time {
	var y, m, d int
	switch hd {
	case Day:
		y, m, d = 0, 0, 1
	case Month:
		y, m, d = 0, 1, 0
	case Year:
		y, m, d = 1, 0, 0
	}
	return t.AddDate(y, m, d)
}

// produce times from a (incl) to b (excl)
func Walk(a, b time.Time, hd HDuration) (<-chan time.Time, error) {
	ch := make(chan time.Time)
	ra := Round(a, hd)
	rb := Round(b, hd)
	if ra == rb {
		rb = Add(rb, hd)
	}

	// fmt.Printf("\n")
	// fmt.Printf("%s\n", time.RFC3339)
	// fmt.Printf(" a: %v\n", a)
	// fmt.Printf("ra: %v\n", ra)
	go func() {
		start := ra
		for start.Before(rb) {
			ch <- start
			start = Add(start, hd)
		}
		close(ch)
	}()
	return ch, nil
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
