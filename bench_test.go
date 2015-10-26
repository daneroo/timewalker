package timewalker

import (
	"testing"
	"time"
)

///////////////////////////
/// Duration Adding --bench Add
///////////////////////////
func BenchmarkAddYear(b *testing.B) {
	t := parseTime("2001-02-03T12:45:56Z")
	for i := 0; i < b.N; i++ {
		_ = Year.AddTo(t)
	}
}
func BenchmarkAddMonth(b *testing.B) {
	t := parseTime("2001-02-03T12:45:56Z")
	for i := 0; i < b.N; i++ {
		_ = Month.AddTo(t)
	}
}

func BenchmarkAddDay(b *testing.B) {
	t := parseTime("2001-02-03T12:45:56Z")
	for i := 0; i < b.N; i++ {
		_ = Day.AddTo(t)
	}
}

// Performs same benchmark as Day.AddTo(t) but not using our Duration.AddTo() method
func BenchmarkAddDayExplicit(b *testing.B) {
	t := parseTime("2001-02-03T12:45:56Z")
	for i := 0; i < b.N; i++ {
		_ = t.AddDate(0, 0, 1)
	}
}

func BenchmarkAddDayInLocation(b *testing.B) {
	l, _ := time.LoadLocation("America/Montreal")
	t := parseTime("2001-02-03T12:45:56Z").In(l)
	for i := 0; i < b.N; i++ {
		_ = Day.AddTo(t)
	}
}

// Performs same benchmark as Day.AddTo(t) in Location but not using our Duration.AddTo() method
func BenchmarkAddDayInLocationExplicit(b *testing.B) {
	l, _ := time.LoadLocation("America/Montreal")
	t := parseTime("2001-02-03T12:45:56Z").In(l)
	for i := 0; i < b.N; i++ {
		_ = t.AddDate(0, 0, 1)
	}
}

///////////////////////////
/// Duration Rounding --bench 'Floor|Ceil'
///////////////////////////

func BenchmarkCeilDay(b *testing.B) {
	t := parseTime("2001-02-03T12:45:56Z")
	for i := 0; i < b.N; i++ {
		_ = Day.Ceil(t)
	}
}

func BenchmarkCeilDayInLocation(b *testing.B) {
	l, _ := time.LoadLocation("America/Montreal")
	t := parseTime("2001-02-03T12:45:56Z").In(l)
	for i := 0; i < b.N; i++ {
		_ = Day.Ceil(t)
	}
}

func BenchmarkFloorDay(b *testing.B) {
	t := parseTime("2001-02-03T12:45:56Z")
	for i := 0; i < b.N; i++ {
		_ = Day.Floor(t)
	}
}

// Performs same benchmark as Day.Floor(t) but not using our Duration.Floor() method
func BenchmarkFloorDayExplicit(b *testing.B) {
	t := parseTime("2001-02-03T12:45:56Z")
	for i := 0; i < b.N; i++ {
		year, month, day := t.Date()
		_ = time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	}
}

func BenchmarkFloorDayInLocation(b *testing.B) {
	l, _ := time.LoadLocation("America/Montreal")
	t := parseTime("2001-02-03T12:45:56Z").In(l)
	for i := 0; i < b.N; i++ {
		_ = Day.Floor(t)
	}
}

// Performs same benchmark as Day.Floor(t) in Location but not using our Duration.Floor() method
func BenchmarkFloorDayInLocationExplicit(b *testing.B) {
	l, _ := time.LoadLocation("America/Montreal")
	t := parseTime("2001-02-03T12:45:56Z").In(l)
	for i := 0; i < b.N; i++ {
		year, month, day := t.Date()
		_ = time.Date(year, month, day, 0, 0, 0, 0, l)
	}
}

func BenchmarkFloorYear(b *testing.B) {
	t := parseTime("2001-02-03T12:45:56Z")
	for i := 0; i < b.N; i++ {
		_ = Year.Floor(t)
	}
}

// Performs same benchmark as Year.Floor(t) but not using our Duration.Floor() method
func BenchmarkFloorYearExplicit(b *testing.B) {
	t := parseTime("2001-02-03T12:45:56Z")
	for i := 0; i < b.N; i++ {
		year, _, _ := t.Date()
		_ = time.Date(year, time.January, 1, 0, 0, 0, 0, t.Location())
	}
}

func BenchmarkFloorYearInLocation(b *testing.B) {
	l, _ := time.LoadLocation("America/Montreal")
	t := parseTime("2001-02-03T12:45:56Z").In(l)
	for i := 0; i < b.N; i++ {
		_ = Year.Floor(t)
	}
}

// Performs same benchmark as Year.Floor(t) in Location but not using our Duration.Floor() method
func BenchmarkFloorYearInLocationExplicit(b *testing.B) {
	l, _ := time.LoadLocation("America/Montreal")
	t := parseTime("2001-02-03T12:45:56Z").In(l)
	for i := 0; i < b.N; i++ {
		year, _, _ := t.Date()
		_ = time.Date(year, time.January, 1, 0, 0, 0, 0, t.Location())
	}
}

///////////////////////////
/// time Parsing --bench Parse
///////////////////////////

// time.Time parsing with time.RFC3339 layout
func BenchmarkParseTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = parseTime("2001-02-03T12:45:56Z")
	}
}

// time.Time parsing with time.RFC3339 layout with location
func BenchmarkParseTimeInLocation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = parseTime("2001-02-03T07:45:56-05:00")
	}
}

// time.Time parsing with time.RFC3339 layout with location, not using our helper method
func BenchmarkParseTimeExplicit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = time.Parse(time.RFC3339, "2001-02-03T12:45:56Z")
	}
}

///////////////////////////
/// time Constructor --bench Constr
///////////////////////////

// time.Time constrcutor
func BenchmarkDateConstructor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC)
	}
}

// time.Time constrcutor with Location
func BenchmarkDateConstructorInLocation(b *testing.B) {
	l, _ := time.LoadLocation("America/Montreal")
	for i := 0; i < b.N; i++ {
		_ = time.Date(2001, time.January, 1, 0, 0, 0, 0, l)
	}
}
