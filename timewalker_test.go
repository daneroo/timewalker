package timewalker

import (
	"fmt"
	"testing"
	"time"
)

// package level example for finding DST Boundaries
func Example_daylightSavingsBoundaries() {
	loc, _ := time.LoadLocation("America/Montreal")
	i, _ := Interval{
		Start: time.Date(2000, time.January, 1, 0, 0, 0, 0, loc),
		End:   time.Date(2009, time.January, 1, 0, 0, 0, 0, loc),
	}.Round(Day)
	fmt.Printf("DST boundaries in %v\n", i)

	days, _ := i.Walk(Day)

	for day := range days {
		hours := day.End.Sub(day.Start).Hours()
		if hours != 24 {
			zs, _ := day.Start.Zone()
			ze, _ := day.End.Zone()
			fmt.Printf("%v (%s->%s) has %.0f hours\n", day.Start.Format("2006-01-02"), zs, ze, hours)
		}
	}

	// Output:
	// DST boundaries in [2000-01-01T00:00:00-05:00, 2009-01-01T00:00:00-05:00)
	// 2000-04-02 (EST->EDT) has 23 hours
	// 2000-10-29 (EDT->EST) has 25 hours
	// 2001-04-01 (EST->EDT) has 23 hours
	// 2001-10-28 (EDT->EST) has 25 hours
	// 2002-04-07 (EST->EDT) has 23 hours
	// 2002-10-27 (EDT->EST) has 25 hours
	// 2003-04-06 (EST->EDT) has 23 hours
	// 2003-10-26 (EDT->EST) has 25 hours
	// 2004-04-04 (EST->EDT) has 23 hours
	// 2004-10-31 (EDT->EST) has 25 hours
	// 2005-04-03 (EST->EDT) has 23 hours
	// 2005-10-30 (EDT->EST) has 25 hours
	// 2006-04-02 (EST->EDT) has 23 hours
	// 2006-10-29 (EDT->EST) has 25 hours
	// 2007-03-11 (EST->EDT) has 23 hours
	// 2007-11-04 (EDT->EST) has 25 hours
	// 2008-03-09 (EST->EDT) has 23 hours
	// 2008-11-02 (EDT->EST) has 25 hours
}

// package level example parsing time in Location
func Example_parseTimeInLocation() {
	// This is how we parse time literals with time.RFC3339

	// default is in UTC when literal terminates in ..Z
	fmt.Println(parseTime("2001-02-03T12:45:56Z"))

	// Now show it's equivalent in EST
	loc, _ := time.LoadLocation("America/Montreal")
	fmt.Println(parseTime("2001-02-03T12:45:56Z").In(loc))

	// Parse EST time directly
	// TODO(daneroo) this seems to depend on Local timezone being set to America/Montreal or equiv...
	fmt.Println(parseTime("2001-02-03T07:45:56-05:00"))

	// Parse EDT time directly
	fmt.Println(parseTime("2001-07-03T07:45:56-04:00"))

	// Output:
	// 2001-02-03 12:45:56 +0000 UTC
	// 2001-02-03 07:45:56 -0500 EST
	// 2001-02-03 07:45:56 -0500 EST
	// 2001-07-03 07:45:56 -0400 EDT
}

var durationTests = []struct {
	dur Duration
	exp string // expected result
}{
	{Day, "Day"},
	{Month, "Month"},
	{Year, "Year"},
}

func TestDuration(t *testing.T) {
	for _, tt := range durationTests {
		actual := tt.dur.String()
		if actual != tt.exp {
			t.Errorf("(%s): exp: %v act: %v", tt.dur, tt.exp, actual)
		}
	}
}

func ExampleDuration_String() {
	fmt.Printf("Day: %v\n", Day)
	fmt.Printf("Month: %v\n", Month)
	fmt.Printf("Year: %v\n", Year)
	// Output:
	// Day: Day
	// Month: Month
	// Year: Year
}

func TestDurationFloor(t *testing.T) {
	var testData = []struct {
		inp time.Time // input
		dur Duration  // Rounding duration
		exp time.Time // expected result
	}{
		{ // Day, already on boundary
			inp: parseTime("2001-02-03T00:00:00Z"),
			dur: Day,
			exp: parseTime("2001-02-03T00:00:00Z"),
		}, { //Day
			inp: parseTime("2001-02-03T12:45:56Z"),
			dur: Day,
			exp: parseTime("2001-02-03T00:00:00Z"),
		}, { //Month
			inp: parseTime("2001-02-03T12:45:56Z"),
			dur: Month,
			exp: parseTime("2001-02-01T00:00:00Z"),
		}, { //Year
			inp: parseTime("2001-02-03T12:45:56Z"),
			dur: Year,
			exp: parseTime("2001-01-01T00:00:00Z"),
		},
	}
	for _, tt := range testData {
		actual := tt.dur.Floor(tt.inp)
		if actual != tt.exp {
			t.Errorf("%s.Floor(%s): \nexp: %v, \nact: %v", tt.dur, tt.inp, tt.exp, actual)
		}
	}
}

func TestDurationCeil(t *testing.T) {
	var testData = []struct {
		inp time.Time // input
		dur Duration  // Rounding duration
		exp time.Time // expected result
	}{
		{ // Day, already on Boundary
			inp: parseTime("2001-02-03T00:00:00Z"),
			dur: Day,
			exp: parseTime("2001-02-03T00:00:00Z"),
		}, { // Day
			inp: parseTime("2001-02-03T12:45:56Z"),
			dur: Day,
			exp: parseTime("2001-02-04T00:00:00Z"),
		}, { // Month
			inp: parseTime("2001-02-03T12:45:56Z"),
			dur: Month,
			exp: parseTime("2001-03-01T00:00:00Z"),
		}, { // Year
			inp: parseTime("2001-02-03T12:45:56Z"),
			dur: Year,
			exp: parseTime("2002-01-01T00:00:00Z"),
		},
	}
	for _, tt := range testData {
		actual := tt.dur.Ceil(tt.inp)
		if actual != tt.exp {
			t.Errorf("%s.Ceil(%s): \nexp: %v, \nact: %v", tt.dur, tt.inp, tt.exp, actual)
		}
	}
}

func ExampleDuration_Ceil_day() {
	t := parseTime("2001-02-03T12:45:56Z")
	rt := Day.Ceil(t)
	fmt.Printf("%v -> %v\n", t, rt)

	// Now in location:
	l, _ := time.LoadLocation("America/Montreal")
	t = t.In(l)
	rt = Day.Ceil(t)
	fmt.Printf("%v -> %v (%v)\n", t, rt, rt.UTC())

	// Output:
	// 2001-02-03 12:45:56 +0000 UTC -> 2001-02-04 00:00:00 +0000 UTC
	// 2001-02-03 07:45:56 -0500 EST -> 2001-02-04 00:00:00 -0500 EST (2001-02-04 05:00:00 +0000 UTC)
}

func ExampleDuration_Floor_month() {
	// t := parseTime("2001-02-03T12:45:56Z")
	t := parseTime("2001-02-03T12:45:56Z")
	rt := Month.Floor(t)
	fmt.Printf("%v -> %v\n", t, rt)
	// Output:
	// 2001-02-03 12:45:56 +0000 UTC -> 2001-02-01 00:00:00 +0000 UTC
}

func ExampleDuration_Floor_day() {
	t := parseTime("2001-02-03T12:45:56Z")
	rt := Day.Floor(t)
	fmt.Printf("%v -> %v\n", t, rt)

	// Now in location:
	l, _ := time.LoadLocation("America/Montreal")
	t = t.In(l)
	rt = Day.Floor(t)
	fmt.Printf("%v -> %v (%v)\n", t, rt, rt.UTC())

	// Output:
	// 2001-02-03 12:45:56 +0000 UTC -> 2001-02-03 00:00:00 +0000 UTC
	// 2001-02-03 07:45:56 -0500 EST -> 2001-02-03 00:00:00 -0500 EST (2001-02-03 05:00:00 +0000 UTC)
}

func TestDurationAdding(t *testing.T) {
	var testData = []struct {
		inp time.Time // input
		dur Duration  // Rounding duration
		exp time.Time // expected result
	}{
		{ // Day
			inp: parseTime("2001-02-03T12:45:56Z"),
			dur: Day,
			exp: parseTime("2001-02-04T12:45:56Z"),
		}, { //Month
			inp: parseTime("2001-02-03T12:45:56Z"),
			dur: Month,
			exp: parseTime("2001-03-03T12:45:56Z"),
		}, { //Year
			inp: parseTime("2001-02-03T12:45:56Z"),
			dur: Year,
			exp: parseTime("2002-02-03T12:45:56Z"),
		},
	}
	for _, tt := range testData {
		actual := tt.dur.AddTo(tt.inp)
		if actual != tt.exp {
			t.Errorf("%s.AddTo(%s): exp: %v, act: %v", tt.dur, tt.inp, tt.exp, actual)
		}
	}
}

func TestWalkEmptyInterval(t *testing.T) {
	sameInstant := parseTime("2001-02-03T12:45:56Z")
	ch, _ := Walk(sameInstant, sameInstant, Day)
	count := 0
	for range ch {
		count++
	}
	if count != 1 {
		t.Error("Expect empty interval to fire once")
	}
}

func ExampleWalk_day() {
	ch, _ := Walk(parseTime("2004-02-26T12:45:56Z"), parseTime("2004-03-03T12:45:56Z"), Day)
	for t := range ch {
		fmt.Printf("%s\n", t)
	}
	// Output:
	// 2004-02-26 00:00:00 +0000 UTC
	// 2004-02-27 00:00:00 +0000 UTC
	// 2004-02-28 00:00:00 +0000 UTC
	// 2004-02-29 00:00:00 +0000 UTC
	// 2004-03-01 00:00:00 +0000 UTC
	// 2004-03-02 00:00:00 +0000 UTC
}

func ExampleWalk_monthLocalAndUTC() {
	l, _ := time.LoadLocation("America/Montreal")
	ch, _ := Walk(parseTime("2001-02-03T12:45:56Z").In(l), parseTime("2002-02-03T12:45:56Z"), Month)
	for t := range ch {
		fmt.Printf("%v (%v)\n", t, t.UTC())
	}
	// Output:
	// 2001-02-01 00:00:00 -0500 EST (2001-02-01 05:00:00 +0000 UTC)
	// 2001-03-01 00:00:00 -0500 EST (2001-03-01 05:00:00 +0000 UTC)
	// 2001-04-01 00:00:00 -0500 EST (2001-04-01 05:00:00 +0000 UTC)
	// 2001-05-01 00:00:00 -0400 EDT (2001-05-01 04:00:00 +0000 UTC)
	// 2001-06-01 00:00:00 -0400 EDT (2001-06-01 04:00:00 +0000 UTC)
	// 2001-07-01 00:00:00 -0400 EDT (2001-07-01 04:00:00 +0000 UTC)
	// 2001-08-01 00:00:00 -0400 EDT (2001-08-01 04:00:00 +0000 UTC)
	// 2001-09-01 00:00:00 -0400 EDT (2001-09-01 04:00:00 +0000 UTC)
	// 2001-10-01 00:00:00 -0400 EDT (2001-10-01 04:00:00 +0000 UTC)
	// 2001-11-01 00:00:00 -0500 EST (2001-11-01 05:00:00 +0000 UTC)
	// 2001-12-01 00:00:00 -0500 EST (2001-12-01 05:00:00 +0000 UTC)
	// 2002-01-01 00:00:00 -0500 EST (2002-01-01 05:00:00 +0000 UTC)
}

func ExampleWalk_year() {
	ch, _ := Walk(parseTime("2001-06-03T12:45:56Z"), parseTime("2005-07-03T12:45:56Z"), Year)
	for t := range ch {
		fmt.Printf("%s\n", t)
	}
	// Output:
	// 2001-01-01 00:00:00 +0000 UTC
	// 2002-01-01 00:00:00 +0000 UTC
	// 2003-01-01 00:00:00 +0000 UTC
	// 2004-01-01 00:00:00 +0000 UTC
}

func ExampleInterval_String() {
	intvl := parseIntvl("2000-01-01T00:00:00Z", "2001-01-01T00:00:00Z")
	fmt.Printf("%v\n", intvl)
	// Output:
	// [2000-01-01T00:00:00Z, 2001-01-01T00:00:00Z)
}

func TestRound(t *testing.T) {
	var roundTests = []struct {
		inp Interval // input
		dur Duration // Rounding duration
		exp Interval // expected result
	}{
		{ // already ok
			inp: parseIntvl("2000-01-01T00:00:00Z", "2001-01-01T00:00:00Z"),
			dur: Day,
			exp: parseIntvl("2000-01-01T00:00:00Z", "2001-01-01T00:00:00Z"),
		}, { //swap start, end
			inp: parseIntvl("2001-01-01T00:00:00Z", "2000-01-01T00:00:00Z"),
			dur: Day,
			exp: parseIntvl("2000-01-01T00:00:00Z", "2001-01-01T00:00:00Z"),
		}, { //round start
			inp: parseIntvl("2000-01-01T00:00:06Z", "2001-01-01T00:00:00Z"),
			dur: Day,
			exp: parseIntvl("2000-01-01T00:00:00Z", "2001-01-01T00:00:00Z"),
		}, { // round end - up
			inp: parseIntvl("2000-01-01T00:00:00Z", "2001-01-01T00:00:06Z"),
			dur: Day,
			exp: parseIntvl("2000-01-01T00:00:00Z", "2001-01-02T00:00:00Z"),
		}, { // round end - up because before start+d
			inp: parseIntvl("2000-01-01T00:00:00Z", "2000-01-01T00:00:06Z"),
			dur: Day,
			exp: parseIntvl("2000-01-01T00:00:00Z", "2000-01-02T00:00:00Z"),
		},
	}

	for _, tt := range roundTests {
		actual, _ := tt.inp.Round(tt.dur)
		if actual != tt.exp {
			t.Errorf("%v.Round(%s): \nexp: %v, \nact: %v", tt.inp, tt.dur, tt.exp, actual)
		}
	}
}

func ExampleInterval_Round() {
	i := parseIntvl("2000-01-01T12:00:00Z", "2001-01-01T12:00:00Z")
	ri, _ := i.Round(Day)
	fmt.Printf("%v.Round(Day)==%v\n", i, ri)

	// Now in location:
	loc, _ := time.LoadLocation("America/Montreal")

	i.Start = i.Start.In(loc)
	i.End = i.End.In(loc)
	ri, _ = i.Round(Day)
	fmt.Printf("%v.Round(Day)==%v\n", i, ri)

	//  Now with mismatched locations
	i.Start = i.Start.UTC()
	i.End = i.End.In(loc)
	_, err := i.Round(Day)
	fmt.Printf("%s\n", err)

	// Output:
	// [2000-01-01T12:00:00Z, 2001-01-01T12:00:00Z).Round(Day)==[2000-01-01T00:00:00Z, 2001-01-02T00:00:00Z)
	// [2000-01-01T07:00:00-05:00, 2001-01-01T07:00:00-05:00).Round(Day)==[2000-01-01T00:00:00-05:00, 2001-01-02T00:00:00-05:00)
	// Interval boundaries have in different time.Location: UTC!=America/Montreal, [2000-01-01T00:00:00Z, 2001-01-02T00:00:00-05:00)
}

func TestTimeWalker(t *testing.T) {
	i := parseIntvl("2000-01-01T00:00:00Z", "2000-01-11T00:00:00Z")

	ch, err := i.Walk(Day)
	if err != nil {
		t.Errorf("TimeWalker generated unexpected error")
	}

	var result []Interval
	for interval := range ch {
		result = append(result, interval)
	}

	if len(result) != 10 {
		t.Errorf("Generated interval has wrong length: %d", len(result))
	}
}

// Utility functions for time literals in our tests
func parseTime(ts string) time.Time {
	lyt := time.RFC3339
	t, err := time.Parse(lyt, ts)
	if err != nil {
		panic(err)
	}
	return t
}

//  Below is Interval stuff
func parseIntvl(a, b string) Interval {
	return Interval{Start: parseTime(a), End: parseTime(b)}
}
