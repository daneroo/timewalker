package timewalker

import (
	"testing"
	"time"
)

func BenchmarkDateConstructor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC)
	}
}
func BenchmarkDateConstructorInLocation(b *testing.B) {
	l, _ := time.LoadLocation("America/Montreal")
	for i := 0; i < b.N; i++ {
		_ = time.Date(2001, time.January, 1, 0, 0, 0, 0, l)
	}
}
