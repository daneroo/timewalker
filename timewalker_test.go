package timewalker

import (
	"testing"
)

func TestTimeWalker(t *testing.T) {
	tw, err := TimeWalker()
	if err != nil {
		t.Errorf("TimeWalker generated unexpected error")
	}
	result := make([]Interval, 0)
	for interval := range tw {
		result = append(result, interval)
	}

	if len(result) != 10 {
		t.Errorf("Generated interval has wrong length")
	}
}
