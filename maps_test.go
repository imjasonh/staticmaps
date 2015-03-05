package maps

import (
	"flag"
	"image/color"
	"net/http"
	"testing"
	"time"
)

var (
	apiKey = flag.String("apiKey", "", "Google Maps API key")

	ctx  = NewContext("key", &http.Client{})
	wctx = NewWorkContext("clientID", "privKey", &http.Client{})
)

func init() {
	flag.Parse()
}

func TestDirections(t *testing.T) {
	orig, dest := Address("111 8th Ave, NYC"), Address("170 E 92nd St, NYC")
	opts := &DirectionsOpts{
		Mode:          ModeTransit,
		DepartureTime: time.Now(),
		Alternatives:  true,
	}
	t.Logf("%s", baseURL+directions(orig, dest, opts))
	r, err := Directions(ctx, orig, dest, opts)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)
}

func TestStaticMap(t *testing.T) {
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
	if _, err := StaticMap(ctx, s, opts); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestStreetView(t *testing.T) {
	s := Size{600, 300}
	opts := &StreetViewOpts{
		Location: &LatLng{46.414382, 10.013988},
		Heading:  Float64(151.78),
		Pitch:    -0.76,
	}
	t.Logf("%s", baseURL+streetview(s, opts))
	if _, err := StreetView(ctx, s, opts); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTimeZone(t *testing.T) {
	ll := LatLng{40.7142700, -74.0059700}
	tm := time.Now()
	t.Logf("%s", baseURL+timezone(ll, tm, nil))
	r, err := TimeZone(ctx, ll, tm, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)
}

func TestElevation(t *testing.T) {
	ll := []LatLng{{39.7391536, -104.9847034}}
	t.Logf("%s", baseURL+elevation(ll))
	r, err := Elevation(ctx, ll)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)

	p := "gfo}EtohhU"
	t.Logf("%s", baseURL+elevationpoly(p))
	r, err = ElevationPolyline(ctx, p)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)

	samples := 3
	ll = []LatLng{{36.578581, -118.291994}, {36.23998, -116.83171}}
	t.Logf("%s", baseURL+elevationpath(ll, samples))
	r, err = ElevationPath(ctx, ll, samples)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)

	p = "gfo}EtohhUxD@bAxJmGF"
	t.Logf("%s", baseURL+elevationpathpoly(p, samples))
	r, err = ElevationPathPoly(ctx, p, samples)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)
}

func TestGeocode(t *testing.T) {
	opts := &GeocodeOpts{
		Address: Address("1600 Amphitheatre Parkway, Mountain View, CA"),
	}
	t.Logf("%s", baseURL+geocode(opts))
	r, err := Geocode(ctx, opts)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)

	ll := LatLng{40.714224, -73.961452}
	t.Logf("%s", baseURL+reversegeocode(ll, nil))
	r, err = ReverseGeocode(ctx, ll, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)
}

func TestDistanceMatrix(t *testing.T) {
	orig := []Location{Address("Vancouver, BC"), Address("Seattle")}
	dst := []Location{Address("San Francisco"), Address("Victoria, BC")}
	opts := &DistanceMatrixOpts{
		Mode:     ModeBicycling,
		Language: "fr-FR",
	}
	t.Logf("%s", baseURL+distancematrix(orig, dst, opts))
	r, err := DistanceMatrix(ctx, orig, dst, opts)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)
}

func TestSnapToRoads(t *testing.T) {
	ctx = NewContext(*apiKey, &http.Client{})

	path := []LatLng{{-35.27801, 149.12958},
		{-35.28032, 149.12907},
		{-35.28099, 149.12929},
		{-35.28144, 149.12984},
		{-35.28194, 149.13003},
		{-35.28282, 149.12956},
		{-35.28302, 149.12881},
		{-35.28473, 149.12836}}
	opts := &SnapToRoadsOpts{
		Interpolate: true,
	}
	r, err := SnapToRoads(ctx, path, opts)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v", r)
}

// Based on https://developers.google.com/maps/documentation/business/webservices/auth#signature_examples
func TestSignature(t *testing.T) {
	clientID := "clientID"
	privateKey := "vNIXE0xscrmjlyV-12Nj_BvUPaw="
	sig, err := genSig(privateKey, "/maps/api/geocode/json", "address=New+York&client="+clientID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	exp := "chaRF2hTJKOScPr-RQCEhZbSzIE="
	if sig != exp {
		t.Errorf("unexpected signature, got %q, want %q", sig, exp)
	}

	ctx := NewWorkContext(clientID, privateKey, nil)
	if gcid, gkey := workCreds(ctx); gcid != clientID || gkey != privateKey {
		t.Errorf("unepxected credentials from context, got %q and %q, want %q and %q", gcid, gkey, clientID, privateKey)
	}
}
