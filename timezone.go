package maps

import (
	"fmt"
	"net/url"
	"time"
)

func (c Client) TimeZone(ll LatLng, t time.Time, opts *TimeZoneOpts) (*TimeZoneResponse, error) {
	var r TimeZoneResponse
	if err := c.doDecode(baseURL+timezone(ll, t, opts), &r); err != nil {
		return nil, err
	}
	return &r, nil
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

type TimeZoneResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
	DSTOffset    int64  `json:"dstOffset"` // secs offset for DST
	RawOffset    int64  `json:"rawOffset"`
	TimeZoneID   string `json:"timeZoneId"`
	TimeZoneName string `json:"timeZoneName"`
}

func (r TimeZoneResponse) DSTOffsetDuration() time.Duration {
	return time.Duration(r.DSTOffset) * time.Second
}

func (r TimeZoneResponse) RawOffsetDuration() time.Duration {
	return time.Duration(r.RawOffset) * time.Second
}
