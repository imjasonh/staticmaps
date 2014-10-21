package maps

import (
	"image/color"
	"testing"
	"time"

	"code.google.com/p/go.net/context"
)

var ctx = context.Background()

func TestDirections(t *testing.T) {
	c := New("", nil)
	orig, dest := Address("111 8th Ave, NYC"), Address("170 E 92nd St, NYC")
	opts := &DirectionsOpts{
		Mode:          ModeTransit,
		DepartureTime: time.Now(),
		Alternatives:  true,
	}
	t.Logf("%s", baseURL+directions(orig, dest, opts))
	r, err := c.Directions(ctx, orig, dest, opts)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)
}

func TestStaticMap(t *testing.T) {
	c := New("", nil)
	s := Size{512, 512}
	opts := &StaticMapOpts{
		Center: LatLng{-5, -5},
		Markers: []Markers{
			{
				Size:  "small",
				Color: color.Black,
				Locations: []Location{
					LatLng{1, 1},
					LatLng{2, 2},
				},
			}, {
				Size:  "mid",
				Color: color.White,
				Locations: []Location{
					LatLng{3, 3},
				},
			},
		},
		Paths: []Path{
			{
				Weight: 10,
				Color:  color.Black,
				Locations: []Location{
					LatLng{4, 4},
					LatLng{5, 5},
				},
			}, {
				Color:     color.Black,
				FillColor: color.White,
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
					{Hue: color.White, Saturation: -100.0},
				},
			},
		},
	}
	t.Logf("%s", baseURL+staticmap(s, opts))
	if _, err := c.StaticMap(ctx, s, opts); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestStreetView(t *testing.T) {
	c := New("", nil)
	s := Size{600, 300}
	opts := &StreetViewOpts{
		Location: &LatLng{46.414382, 10.013988},
		Heading:  Float64(151.78),
		Pitch:    -0.76,
	}
	t.Logf("%s", baseURL+streetview(s, opts))
	if _, err := c.StreetView(ctx, s, opts); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTimeZone(t *testing.T) {
	c := New("", nil)
	ll := LatLng{40.7142700, -74.0059700}
	tm := time.Now()
	t.Logf("%s", baseURL+timezone(ll, tm, nil))
	r, err := c.TimeZone(ctx, ll, tm, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)
}

func TestElevation(t *testing.T) {
	c := New("", nil)
	ll := []LatLng{{39.7391536, -104.9847034}}
	t.Logf("%s", baseURL+elevation(ll))
	r, err := c.Elevation(ctx, ll)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)

	p := "gfo}EtohhU"
	t.Logf("%s", baseURL+elevationpoly(p))
	r, err = c.ElevationPolyline(ctx, p)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)

	samples := 3
	ll = []LatLng{{36.578581, -118.291994}, {36.23998, -116.83171}}
	t.Logf("%s", baseURL+elevationpath(ll, samples))
	r, err = c.ElevationPath(ctx, ll, samples)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)

	p = "gfo}EtohhUxD@bAxJmGF"
	t.Logf("%s", baseURL+elevationpathpoly(p, samples))
	r, err = c.ElevationPathPoly(ctx, p, samples)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)
}

func TestGeocode(t *testing.T) {
	c := New("", nil)
	opts := &GeocodeOpts{
		Address: Address("1600 Amphitheatre Parkway, Mountain View, CA"),
	}
	t.Logf("%s", baseURL+geocode(opts))
	r, err := c.Geocode(ctx, opts)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)

	ll := LatLng{40.714224, -73.961452}
	t.Logf("%s", baseURL+reversegeocode(ll, nil))
	r, err = c.ReverseGeocode(ctx, ll, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)
}

func TestDistanceMatrix(t *testing.T) {
	c := New("", nil)
	orig := []Location{Address("Vancouver, BC"), Address("Seattle")}
	dst := []Location{Address("San Francisco"), Address("Victoria, BC")}
	opts := &DistanceMatrixOpts{
		Mode:     ModeBicycling,
		Language: "fr-FR",
	}
	t.Logf("%s", baseURL+distancematrix(orig, dst, opts))
	r, err := c.DistanceMatrix(ctx, orig, dst, opts)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)

}

// Based on https://developers.google.com/maps/documentation/business/webservices/auth#signature_examples
func TestSignature(t *testing.T) {
	c := NewForWork("", "vNIXE0xscrmjlyV-12Nj_BvUPaw=", nil)
	sig, err := c.genSig("/maps/api/geocode/json", "address=New+York&client=clientID")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	exp := "chaRF2hTJKOScPr-RQCEhZbSzIE="
	if sig != exp {
		t.Errorf("wrong signature, got %q, want %q", sig, exp)
	}
}
