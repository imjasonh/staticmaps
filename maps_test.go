package maps

import (
	"testing"
	"time"
)

func TestDirections(t *testing.T) {
	c := NewClient("")
	orig, dest := Address("111 8th Ave, NYC"), Address("170 E 92nd St, NYC")
	opts := &DirectionsOpts{
		Mode:          ModeWalking,
		DepartureTime: time.Now(),
	}
	t.Logf("%s", baseURL+directions(orig, dest, opts))
	r, err := c.Directions(orig, dest, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)
}
