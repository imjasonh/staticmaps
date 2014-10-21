package maps

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"code.google.com/p/go.net/context"
)

// StreetView requests a static StreetView image of the requested size.
func (c Client) StreetView(ctx context.Context, s Size, opts *StreetViewOpts) (io.ReadCloser, error) {
	resp, err := c.do(ctx, baseURL+streetview(s, opts))
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

// StreetViewOpts defines options for StreetView requests.
type StreetViewOpts struct {
	// Location specifies the location where the image should be snapped from.
	Location Location

	// Pano specifies a specific Panoramio ID.
	Pano string

	// Heading specifies the compass heading (between 0 and 360) of the camera.
	//
	// Both 0 and 360 indicate North, 90 indicates East, 180 indicates South, and so on.
	//
	// If no heading is specified, a value will be calculated that directs the camera towards the specified Location, from a point at which the closest photograph was taken.
	Heading *float64

	// FOV specifies the field of view of the image.
	//
	// Accepted values are between 0 and 120, and if no value is given, a field of view of 90 degrees will be used.
	FOV *float64

	// Pitch specifies the up or down angle of the camera relative to the Street View vehicle.
	//
	// This is often, but not always, flat horizontal. Positive values angle the camera up (90 degrees indicating straight up); negative values angle the camera down (-90 indicates straight down).
	Pitch float64
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
