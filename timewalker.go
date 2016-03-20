// Package timewalker manages sequences of time instances. It is meant to be a source channel of time.Time, or timewalker.Interval
package timewalker

import (
	"fmt"
	"time"
)

// Duration represents common practical Durations, which do not have the same semantics as time.Duration (non-arithmetic), in the sense that not all years, month or days have the same length. For Examples are leap years, months with different number of days, and days on daylight savings boundaries which may not have 24 hours.
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

// Floor returns the greatest time.Time that is on receivers Duration boundary; akin to math.Floor for ints
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

// Ceil returns the least time.Time that is on receivers Duration boundary; akin to math.Ceil for ints
func (d Duration) Ceil(t time.Time) time.Time {
	least := d.Floor(t)
	if least.Before(t) {
		least = d.AddTo(least)
	}
	return least
}

// AddTo returns a new Time by adding the receiver's duration to the passed time parameter
func (d Duration) AddTo(t time.Time) time.Time {
	var yr, mo, dy int
	switch d {
	case Day:
		yr, mo, dy = 0, 0, 1
	case Month:
		yr, mo, dy = 0, 1, 0
	case Year:
		yr, mo, dy = 1, 0, 0
	}
	return t.AddDate(yr, mo, dy)
}

// Walk produces times from a (incl) to b (excl)
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

// Interval represents a time interval from [Start,End)
type Interval struct {
	Start time.Time
	End   time.Time
}

func (i Interval) String() string {
	layout := time.RFC3339
	return fmt.Sprintf("[%s, %s)", i.Start.Format(layout), i.End.Format(layout))
}

/*
Round normalizes an Interval's representations
	-if Start,End are not in same location, throw
	-Swap Start,End if appropriate (if End.Before(Start))
	-Round down (Floor) Start, Round up (Ceil) End, both on Duration boundary
	-Make sure we have at least one interval.
*/
func (i Interval) Round(d Duration) (Interval, error) {
	// BUG(daneroo): Interval Rounding behavior is not well defined yet. This is also an example of a BUG comment showing up in the godocs
	if i.End.Before(i.Start) {
		i.End, i.Start = i.Start, i.End
	}
	i.Start = d.Floor(i.Start)
	i.End = d.Ceil(i.End)

	// minimum End is d.AddTo(i.Start)
	// not sure this is actually possible...
	// because we know i.Start <= i.End
	// d.Floor(i.Start) <= d.Ceil(i.End)
	minEnd := d.AddTo(i.Start)
	if i.End.Before(minEnd) {
		i.End = minEnd
	}
	if i.Start.Location() != i.End.Location() {
		return i, fmt.Errorf("Interval boundaries have in different time.Location: %s!=%s, %v", i.Start.Location(), i.End.Location(), i)
	}
	return i, nil
}

// Walk traverses the receiver's interval in steps of  the given duration
func (i Interval) Walk(d Duration) (<-chan Interval, error) {
	// Round interval
	ri, err := i.Round(d)
	// TODO(daneroo) What is the idomatic way of returning the channel on error condition
	if err != nil {
		return nil, err
	}

	ch := make(chan Interval)

	go func() {
		start := ri.Start
		for start.Before(ri.End) {
			end := d.AddTo(start)
			ch <- Interval{Start: start, End: end}
			start = end
		}
		close(ch)
	}()
	return ch, nil
}
