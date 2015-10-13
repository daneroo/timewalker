package timewalker

import (
	"testing"
	"time"
)

func TestAdder(t *testing.T) {
	TryAdd()
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
