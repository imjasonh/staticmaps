package maps

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// TimeZone requests time zone information about a location.
//
// See https://developers.google.com/maps/documentation/timezone/
func TimeZone(ctx context.Context, ll LatLng, t time.Time, opts *TimeZoneOpts) (*TimeZoneResult, error) {
	var r timeZoneResponse
	if err := doDecode(ctx, baseURL+timezone(ll, t, opts), &r); err != nil {
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

// TimeZoneOpts defines options for TimeZone requests.
type TimeZoneOpts struct {
	// The language in which to return results.
	//
	// See https://developers.google.com/maps/faq#languagesupport
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

// TimeZoneResult describes information about the time zone at the requested location.
type TimeZoneResult struct {
	// DSTOffset is the offset for daylight-savings time in seconds.
	//
	// This will be zero if the time zone is not in Daylight Savings Time during the specified time.
	DSTOffset int64 `json:"dstOffset"`

	// RawOffset is the offset from UTC (in seconds) for the given location.
	//
	// This does not take into effect daylight savings.
	RawOffset int64 `json:"rawOffset"`

	// TimeZoneID is a string containing the "tz" ID of the time zone, such as "America/Los_Angeles"
	TimeZoneID string `json:"timeZoneId"`

	// TimeZoneName is a string containing the long form name of the time zone, e.g., "Pacific Daylight Time".
	//
	// This field will be localized if the Language was specified.
	TimeZoneName string `json:"timeZoneName"`
}

// DSTOffsetDuration translates the DSTOffset into a time.Duration
func (r TimeZoneResult) DSTOffsetDuration() time.Duration {
	return time.Duration(r.DSTOffset) * time.Second
}

// RawOffsetDuration translates the RawOffset into a time.Duration
func (r TimeZoneResult) RawOffsetDuration() time.Duration {
	return time.Duration(r.RawOffset) * time.Second
}
