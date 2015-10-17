package timewalker

import (
	"testing"
	"time"
)

var durationTests = []struct {
	hdur HDuration
	exp  string // expected result
}{
	{Day, "Day"},
	{Month, "Month"},
	{Year, "Year"},
}

func TestHDuration(t *testing.T) {
	for _, tt := range durationTests {
		actual := tt.hdur.String()
		if actual != tt.exp {
			t.Errorf("(%s): exp: %v act: %v", tt.hdur, tt.exp, actual)
		}
	}
}

func parseTime(ts string) time.Time {
	lyt := time.RFC3339
	t, err := time.Parse(lyt, ts)
	if err != nil {
		panic(err)
	}
	return t
}

var roundingTests = []struct {
	inp  time.Time // input
	hdur HDuration // Rounding duration
	exp  time.Time // expected result
}{
	{ // Day
		inp:  parseTime("2001-02-03T12:45:56Z"),
		hdur: Day,
		exp:  parseTime("2001-02-03T00:00:00Z"),
	}, { //Month
		inp:  parseTime("2001-02-03T12:45:56Z"),
		hdur: Month,
		exp:  parseTime("2001-02-01T00:00:00Z"),
	}, { //Year
		inp:  parseTime("2001-02-03T12:45:56Z"),
		hdur: Year,
		exp:  parseTime("2001-01-01T00:00:00Z"),
	},
}

func TestRounding(t *testing.T) {
	for _, tt := range roundingTests {
		actual := Round(tt.inp, tt.hdur)
		if actual != tt.exp {
			t.Errorf("Round(%s,%s): \nexp: %v, \nact: %v", tt.inp, tt.hdur, tt.exp, actual)
		}
	}
}

var addingTests = []struct {
	inp  time.Time // input
	hdur HDuration // Rounding duration
	exp  time.Time // expected result
}{
	{ // Day
		inp:  parseTime("2001-02-03T12:45:56Z"),
		hdur: Day,
		exp:  parseTime("2001-02-04T12:45:56Z"),
	}, { //Month
		inp:  parseTime("2001-02-03T12:45:56Z"),
		hdur: Month,
		exp:  parseTime("2001-03-03T12:45:56Z"),
	}, { //Year
		inp:  parseTime("2001-02-03T12:45:56Z"),
		hdur: Year,
		exp:  parseTime("2002-02-03T12:45:56Z"),
	},
}

func TestAdding(t *testing.T) {
	for _, tt := range addingTests {
		actual := Add(tt.inp, tt.hdur)
		if actual != tt.exp {
			t.Errorf("Add(%s,%s): exp: %v, act: %v", tt.inp, tt.hdur, tt.exp, actual)
		}
	}
}

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
