package maps

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c Client) StreetView(s Size, opts *StreetViewOpts) (io.ReadCloser, error) {
	resp, err := c.do(baseURL + streetview(s, opts))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, HTTPError{resp}
	}
	return resp.Body, nil
}

func streetview(s Size, opts *StreetViewOpts) string {
	p := url.Values{}
	p.Set("size", s.String())
	opts.update(p)
	return "streetview?" + p.Encode()
}

type StreetViewOpts struct {
	Location Location
	Pano     string   // panoramio ID
	Heading  *float64 // 0-360
	FOV      *float64 // 0-120, 90 default
	Pitch    float64  // -90 to +90
}

func (s *StreetViewOpts) update(p url.Values) {
	if s == nil {
		return
	}
	if s.Location != nil {
		p.Set("location", s.Location.Location())
	}
	if s.Pano != "" {
		p.Set("pano", s.Pano)
	}
	if s.Heading != nil {
		p.Set("heading", fmt.Sprintf("%f", *s.Heading))
	}
	if s.FOV != nil {
		p.Set("fov", fmt.Sprintf("%f", *s.FOV))
	}
	if s.Pitch != 0 {
		p.Set("pitch", fmt.Sprintf("%f", s.Pitch))
	}
}
