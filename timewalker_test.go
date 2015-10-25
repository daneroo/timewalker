package timewalker

import (
	"fmt"
	"testing"
	"time"
)

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

var durationRoundingTests = []struct {
	inp time.Time // input
	dur Duration  // Rounding duration
	exp time.Time // expected result
}{
	{ // Day
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

func TestDurationRounding(t *testing.T) {
	for _, tt := range durationRoundingTests {
		actual := tt.dur.Round(tt.inp)
		if actual != tt.exp {
			t.Errorf("%s.Round(%s): \nexp: %v, \nact: %v", tt.dur, tt.inp, tt.exp, actual)
		}
	}
}

func ExampleDuration_Round_month() {
	// t := parseTime("2001-02-03T12:45:56Z")
	t := parseTime("2001-02-03T12:45:56Z")
	rt := Month.Round(t)
	fmt.Printf("%v -> %v", t, rt)
	// Output:
	// 2001-02-03 12:45:56 +0000 UTC -> 2001-02-01 00:00:00 +0000 UTC
}

func ExampleDuration_Round_day() {
	t := parseTime("2001-02-03T12:45:56Z")
	rt := Day.Round(t)
	fmt.Printf("%v -> %v", t, rt)
	// Output:
	// 2001-02-03 12:45:56 +0000 UTC -> 2001-02-03 00:00:00 +0000 UTC
}

func ExampleDuration_Round_dayInLocation() {
	l, _ := time.LoadLocation("America/Montreal")
	t := parseTime("2001-02-03T12:45:56Z").In(l)
	rt := Day.Round(t)
	fmt.Printf("%v -> %v", t, rt)
	// Output:
	// 2001-02-03 07:45:56 -0500 EST -> 2001-02-03 00:00:00 -0500 EST
}

var durationAddingTests = []struct {
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

func TestDurationAdding(t *testing.T) {
	for _, tt := range durationAddingTests {
		actual := tt.dur.AddTo(tt.inp)
		if actual != tt.exp {
			t.Errorf("%s.AddTo(%s): exp: %v, act: %v", tt.dur, tt.inp, tt.exp, actual)
		}
	}
}

func TestWalkEmptyInterval(t *testing.T) {
	ch, _ := Walk(parseTime("2001-02-03T12:45:56Z"), parseTime("2001-02-03T12:45:56Z"), Day)
	count := 0
	for _ = range ch {
		count++
	}
	if count != 1 {
		t.Error("Expect empty interval to fire once")
	}
}

func ExampleWalk_month() {
	ch, _ := Walk(parseTime("2004-02-03T12:45:56Z"), parseTime("2004-03-03T12:45:56Z"), Day)
	for t := range ch {
		fmt.Printf("%s\n", t)
	}
	// Output:
	// 2004-02-03 00:00:00 +0000 UTC
	// 2004-02-04 00:00:00 +0000 UTC
	// 2004-02-05 00:00:00 +0000 UTC
	// 2004-02-06 00:00:00 +0000 UTC
	// 2004-02-07 00:00:00 +0000 UTC
	// 2004-02-08 00:00:00 +0000 UTC
	// 2004-02-09 00:00:00 +0000 UTC
	// 2004-02-10 00:00:00 +0000 UTC
	// 2004-02-11 00:00:00 +0000 UTC
	// 2004-02-12 00:00:00 +0000 UTC
	// 2004-02-13 00:00:00 +0000 UTC
	// 2004-02-14 00:00:00 +0000 UTC
	// 2004-02-15 00:00:00 +0000 UTC
	// 2004-02-16 00:00:00 +0000 UTC
	// 2004-02-17 00:00:00 +0000 UTC
	// 2004-02-18 00:00:00 +0000 UTC
	// 2004-02-19 00:00:00 +0000 UTC
	// 2004-02-20 00:00:00 +0000 UTC
	// 2004-02-21 00:00:00 +0000 UTC
	// 2004-02-22 00:00:00 +0000 UTC
	// 2004-02-23 00:00:00 +0000 UTC
	// 2004-02-24 00:00:00 +0000 UTC
	// 2004-02-25 00:00:00 +0000 UTC
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
		fmt.Printf("%v %v\n", t, t.UTC())
	}
	// Output:
	// 2001-02-01 00:00:00 -0500 EST 2001-02-01 05:00:00 +0000 UTC
	// 2001-03-01 00:00:00 -0500 EST 2001-03-01 05:00:00 +0000 UTC
	// 2001-04-01 00:00:00 -0500 EST 2001-04-01 05:00:00 +0000 UTC
	// 2001-05-01 00:00:00 -0400 EDT 2001-05-01 04:00:00 +0000 UTC
	// 2001-06-01 00:00:00 -0400 EDT 2001-06-01 04:00:00 +0000 UTC
	// 2001-07-01 00:00:00 -0400 EDT 2001-07-01 04:00:00 +0000 UTC
	// 2001-08-01 00:00:00 -0400 EDT 2001-08-01 04:00:00 +0000 UTC
	// 2001-09-01 00:00:00 -0400 EDT 2001-09-01 04:00:00 +0000 UTC
	// 2001-10-01 00:00:00 -0400 EDT 2001-10-01 04:00:00 +0000 UTC
	// 2001-11-01 00:00:00 -0500 EST 2001-11-01 05:00:00 +0000 UTC
	// 2001-12-01 00:00:00 -0500 EST 2001-12-01 05:00:00 +0000 UTC
	// 2002-01-01 00:00:00 -0500 EST 2002-01-01 05:00:00 +0000 UTC
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

//  Below is Interval stuff
func parseIntvl(a, b string) Interval {
	lyt := time.RFC3339
	ta, err := time.Parse(lyt, a)
	if err != nil {
		panic(err)
	}
	tb, err := time.Parse(lyt, b)
	if err != nil {
		panic(err)
	}
	return Interval{Start: ta, End: tb}
}

var roundTests = []struct {
	inp Interval      // input
	dur time.Duration // Rounding duration
	exp Interval      // expected result
}{
	{ // already ok
		inp: parseIntvl("2000-01-01T00:00:00Z", "2001-01-01T00:00:00Z"),
		dur: time.Second * 10,
		exp: parseIntvl("2000-01-01T00:00:00Z", "2001-01-01T00:00:00Z"),
	}, { //swap start, end
		inp: parseIntvl("2001-01-01T00:00:00Z", "2000-01-01T00:00:00Z"),
		dur: time.Second * 10,
		exp: parseIntvl("2000-01-01T00:00:00Z", "2001-01-01T00:00:00Z"),
	}, { //round start
		inp: parseIntvl("2000-01-01T00:00:06Z", "2001-01-01T00:00:00Z"),
		dur: time.Second * 10,
		exp: parseIntvl("2000-01-01T00:00:00Z", "2001-01-01T00:00:00Z"),
	}, { // round end - up
		inp: parseIntvl("2000-01-01T00:00:00Z", "2001-01-01T00:00:06Z"),
		dur: time.Second * 10,
		exp: parseIntvl("2000-01-01T00:00:00Z", "2001-01-01T00:00:10Z"),
	}, { // round end - up because before start+d
		inp: parseIntvl("2000-01-01T00:00:00Z", "2000-01-01T00:00:06Z"),
		dur: time.Second * 10,
		exp: parseIntvl("2000-01-01T00:00:00Z", "2000-01-01T00:00:10Z"),
	}, { // Try with Hour - round start, end ok
		inp: parseIntvl("2000-01-01T01:23:45Z", "2001-01-01T00:00:00Z"),
		dur: time.Hour,
		exp: parseIntvl("2000-01-01T01:00:00Z", "2001-01-01T00:00:00Z"),
	}, { // Hour - round end up
		inp: parseIntvl("2000-01-01T01:23:45Z", "2001-01-01T01:23:45Z"),
		dur: time.Hour,
		exp: parseIntvl("2000-01-01T01:00:00Z", "2001-01-01T02:00:00Z"),
	},
}

func TestRound(t *testing.T) {
	t.Skip()

	for _, tt := range roundTests {
		actual := tt.inp.Round(tt.dur)
		t.Logf("-inp %v\n", tt.inp)
		t.Logf("-dur %v\n", tt.dur)
		t.Logf("-act %v\n", actual)
		t.Logf("-exp %v\n", tt.exp)
		if actual != tt.exp {
			t.Errorf("Round(%s):\ninp :%v \nexp: %v, \nact: %v", tt.dur, tt.inp, tt.exp, actual)
		}
	}
	// i := Interval{Start: time.Unix(1e8, 0), End: time.Unix(10, 5e8)}
	// t.Logf("-i %v\n", i)
	// r := i.Round(time.Minute)
	// t.Logf("+i %v\n", i)
	// t.Logf("+r %v\n", r)
}

func TestTimeWalker(t *testing.T) {
	t.Skip()
	return
	// i := Interval{Start: time.Now(), End: time.Now().Add(10 * time.Second)}
	i := Interval{Start: time.Unix(0, 0), End: time.Unix(10, 5e8)}

	ch, err := i.Iter(time.Second)
	if err != nil {
		t.Errorf("TimeWalker generated unexpected error")
	}

	result := make([]Interval, 0)
	for interval := range ch {
		result = append(result, interval)
	}

	if len(result) != 10 {
		t.Errorf("Generated interval has wrong length: %d", len(result))
	}
}

// Utility function for time literals in our tests
func parseTime(ts string) time.Time {
	lyt := time.RFC3339
	t, err := time.Parse(lyt, ts)
	if err != nil {
		panic(err)
	}
	return t
}
