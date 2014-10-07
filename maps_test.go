package maps

import (
	"testing"
	"time"
)

func TestDirections(t *testing.T) {
	c := NewClient("")
	orig, dest := Address("111 8th Ave, NYC"), Address("170 E 92nd St, NYC")
	opts := &DirectionsOpts{
		Mode:          ModeTransit,
		DepartureTime: time.Now(),
		Alternatives:  true,
	}
	t.Logf("%s", baseURL+directions(orig, dest, opts))
	r, err := c.Directions(orig, dest, opts)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)
}

func TestStaticMap(t *testing.T) {
	c := NewClient("")
	s := Size{512, 512}
	opts := &StaticMapOpts{
		Center: LatLng{-5, -5},
		Markers: []Markers{
			{
				Size:  "small",
				Color: "blue",
				Locations: []Location{
					LatLng{1, 1},
					LatLng{2, 2},
				},
			}, {
				Size:  "mid",
				Color: "red",
				Locations: []Location{
					LatLng{3, 3},
				},
			},
		},
		Paths: []Path{
			{
				Weight: 10,
				Color:  "green",
				Locations: []Location{
					LatLng{4, 4},
					LatLng{5, 5},
				},
			}, {
				Color:     "0x00000000",
				FillColor: "0x00000033",
				Locations: []Location{
					LatLng{6, 6}, LatLng{7, 7}, LatLng{7, 3},
				},
			},
		},
		Visible: []Location{
			LatLng{15, 15},
		},
		Styles: []Style{
			{
				Feature: "water",
				Element: "geometry.fill",
				Rules: []StyleRule{
					{Hue: "0x0000FF"},
				},
			},
		},
	}
	t.Logf("%s", baseURL+staticmap(s, opts))
	if _, err := c.StaticMap(s, opts); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestStreetView(t *testing.T) {
	c := NewClient("")
	s := Size{600, 300}
	opts := &StreetViewOpts{
		Location: &LatLng{46.414382, 10.013988},
		Heading:  Float64(151.78),
		Pitch:    -0.76,
	}
	t.Logf("%s", baseURL+streetview(s, opts))
	if _, err := c.StreetView(s, opts); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTimeZone(t *testing.T) {
	c := NewClient("")
	ll := LatLng{40.7142700, -74.0059700}
	tm := time.Now()
	t.Logf("%s", baseURL+timezone(ll, tm, nil))
	r, err := c.TimeZone(ll, tm, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)
}

func TestElevation(t *testing.T) {
	c := NewClient("")
	ll := []LatLng{{39.7391536, -104.9847034}}
	t.Logf("%s", baseURL+elevation(ll))
	r, err := c.Elevation(ll)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)

	p := "gfo}EtohhU"
	t.Logf("%s", baseURL+elevationpoly(p))
	r, err = c.ElevationPolyline(p)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)

	samples := 3
	ll = []LatLng{{36.578581, -118.291994}, {36.23998, -116.83171}}
	t.Logf("%s", baseURL+elevationpath(ll, samples))
	r, err = c.ElevationPath(ll, samples)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)

	p = "gfo}EtohhUxD@bAxJmGF"
	t.Logf("%s", baseURL+elevationpathpoly(p, samples))
	r, err = c.ElevationPathPoly(p, samples)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)
}
