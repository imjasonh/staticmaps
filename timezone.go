package maps

import (
	"fmt"
	"net/url"
	"time"
)

func (c Client) TimeZone(ll LatLng, t time.Time, opts *TimeZoneOpts) (*TimeZoneResult, error) {
	var r timeZoneResponse
	if err := c.doDecode(baseURL+timezone(ll, t, opts), &r); err != nil {
		return nil, err
	}
	if r.Status != StatusOK {
		return nil, APIError{r.Status, r.ErrorMessage}
	}
	return &r.TimeZoneResult, nil
}

func timezone(ll LatLng, t time.Time, opts *TimeZoneOpts) string {
	p := url.Values{}
	p.Set("location", ll.String())
	p.Set("timestamp", fmt.Sprintf("%d", t.Unix()))
	opts.update(p)
	return "timezone/json?" + p.Encode()
}

type TimeZoneOpts struct {
	Language string
}

func (t *TimeZoneOpts) update(p url.Values) {
	if t == nil {
		return
	}
	if t.Language != "" {
		p.Set("language", t.Language)
	}
}

type timeZoneResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
	TimeZoneResult
}

type TimeZoneResult struct {
	DSTOffset    int64  `json:"dstOffset"` // secs offset for DST
	RawOffset    int64  `json:"rawOffset"`
	TimeZoneID   string `json:"timeZoneId"`
	TimeZoneName string `json:"timeZoneName"`
}

func (r TimeZoneResult) DSTOffsetDuration() time.Duration {
	return time.Duration(r.DSTOffset) * time.Second
}

func (r TimeZoneResult) RawOffsetDuration() time.Duration {
	return time.Duration(r.RawOffset) * time.Second
}
