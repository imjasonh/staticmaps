package maps

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBackoff(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "fail", http.StatusInternalServerError)
	}))
	defer s.Close()
	tries := 10
	sleeps := 0
	bt := &backoff{
		MaxTries:  tries,
		Transport: http.DefaultTransport,
		sleep: func(d time.Duration) {
			sleeps += 1
			t.Logf("sleeping %s", d)
		},
	}
	c := &http.Client{Transport: bt}
	r, _ := http.NewRequest("GET", s.URL, nil)
	if _, err := c.Do(r); err == nil {
		t.Errorf("expected error")
	}
	if sleeps != tries-1 {
		t.Errorf("unexpected # of calls to sleep, got %d, want %d", tries-1, sleeps)
	}
	if bt.tries != tries {
		t.Errorf("got %d tries, want %d", bt.tries, tries)
	}
}
