// Package timewalker manages sequences of time instances. It is meant to be a source channel of time.Time, or timewalker.Interval

package timewalker

import (
	"fmt"
	"time"
)

// Common Practical Durations, which do not have the same semantics as time.Duration (non-arithmetic), in the sense that not all years, month or days have the same length. For Examples are leap years, months with different number of days, and days on daylight savings boundaries which may not have 24 hours.
type Duration int

// Different pakage constants defining an enum type for Duration
const (
	Day Duration = iota
	Month
	Year
)

// Produces Human readble represations of the Duration enum values
func (d Duration) String() string {
	str := "Invalid"
	switch d {
	case Day:
		str = "Day"
	case Month:
		str = "Month"
	case Year:
		str = "Year"
	}
	return str
}

// Returns the greatest time.Time that is on recievers Duration boundary; akin to math.Floor for ints
func (d Duration) Floor(t time.Time) time.Time {
	year, month, day := t.Date()
	switch d {
	case Day:
		t = time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	case Month:
		t = time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
	case Year:
		t = time.Date(year, time.January, 1, 0, 0, 0, 0, t.Location())
	}
	return t
}

// Returns the least time.Time that is on recievers Duration boundary; akin to math.Ceil for ints
func (d Duration) Ceil(t time.Time) time.Time {
	least := d.Floor(t)
	if least.Before(t) {
		least = d.AddTo(least)
	}
	return least
}

func (dur Duration) AddTo(t time.Time) time.Time {
	var y, m, d int
	switch dur {
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
func Walk(a, b time.Time, d Duration) (<-chan time.Time, error) {
	ch := make(chan time.Time)
	ra := d.Floor(a)
	rb := d.Floor(b)
	if ra == rb {
		rb = d.AddTo(rb)
	}

	// fmt.Printf("\n")
	// fmt.Printf("%s\n", time.RFC3339)
	// fmt.Printf(" a: %v\n", a)
	// fmt.Printf("ra: %v\n", ra)
	go func() {
		start := ra
		for start.Before(rb) {
			ch <- start
			start = d.AddTo(start)
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
// BUG(daneroo): Interval Rounding behavior is not well defined yet. This is also an example of a BUG comment showing up in the godocs

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
