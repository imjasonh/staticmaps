package maps

import (
	"math/rand"
	"net/http"
	"time"
)

type backoff struct {
	Transport http.RoundTripper
	MaxTries  int
	tries     int
	sleep     func(time.Duration)
}

func (bt *backoff) RoundTrip(req *http.Request) (*http.Response, error) {
	if bt.MaxTries == 0 {
		bt.MaxTries = 5
	}
	if bt.sleep == nil {
		bt.sleep = time.Sleep
	}
	wait := time.Second
	var err error
	for ; bt.tries < bt.MaxTries; bt.tries++ {
		var resp *http.Response
		resp, err = bt.Transport.RoundTrip(req)
		if err != nil {
			// fallthrough, retry
		} else if resp.StatusCode >= http.StatusInternalServerError {
			err = HTTPError{resp}
		} else {
			return resp, nil
		}
		if bt.tries == bt.MaxTries-1 {
			// last try failed, just give up
			continue
		}
		bt.sleep(wait)
		jitter := time.Duration(rand.Intn(1000))*time.Millisecond - 500*time.Millisecond // jitter is +/- .5s
		wait = wait*2 + jitter
	}
	return nil, err
}
