package maps

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	ComponentRoute              = "route"
	ComponentLocality           = "locality"
	ComponentAdministrativeArea = "administrative_area"
	ComponentPostalCode         = "postal_code"
	ComponentCountry            = "country"

	LocationTypeRooftop           = "ROOFTOP"
	LocationTypeRangeInterpolated = "RANGE_INTERPOLATED"
	LocationTypeGeometricCenter   = "GEOMETRIC_CENTER"
	LocationTypeApproximate       = "APPROXIMATE"

	// TODO: enums for https://developers.google.com/maps/documentation/geocoding/#Types
)

func (c Client) Geocode(opts *GeocodeOpts) ([]GeocodeResult, error) {
	var r geocodeResponse
	if err := c.doDecode(baseURL+geocode(opts), &r); err != nil {
		return nil, err
	}
	if r.Status != StatusOK {
		return nil, APIError{r.Status, ""}
	}
	return r.Results, nil
}

func geocode(opts *GeocodeOpts) string {
	p := url.Values{}
	opts.update(p)
	return "geocode/json?" + p.Encode()
}

type GeocodeOpts struct {
	Address          Address
	Components       []Component
	Language, Region string
	Bounds           *Bounds
}

type Component struct {
	Key, Value string
}

func (c Component) encode() string {
	return fmt.Sprintf("%s:%s", c.Key, c.Value)
}

func encodeComponents(cs []Component) string {
	s := make([]string, len(cs))
	for i, c := range cs {
		s[i] = c.encode()
	}
	return strings.Join(s, "|")
}

func (g *GeocodeOpts) update(p url.Values) {
	if g == nil {
		return
	}
	if string(g.Address) != "" {
		p.Set("address", g.Address.Location())
	}
	if g.Components != nil {
		p.Set("components", encodeComponents(g.Components))
	}
	if g.Language != "" {
		p.Set("language", g.Language)
	}
	if g.Region != "" {
		p.Set("region", g.Region)
	}
	if g.Bounds != nil {
		p.Set("bounds", g.Bounds.String())
	}
}

type geocodeResponse struct {
	Results []GeocodeResult `json:"results"`
	Status  string          `json:"status"`
}
type GeocodeResult struct {
	AddressComponents []struct {
		LongName  string   `json:"long_name"`
		ShortName string   `json:"short_name"`
		Types     []string `json:"types"`
	} `json:"address_components"`
	PostcodeLocalities []string `json:"postcode_localities"`
	FormattedAddress   string   `json:"formatted_address"`
	Geometry           struct {
		Location     LatLng `json:"location"`
		LocationType string `json:"location_type"`
		Viewport     Bounds `json:"viewport"`
		Bounds       Bounds `json:"bounds"`
	} `json:"geometry"`
	Types        []string `json:"types"`
	PartialMatch string   `json:"partial_match"`
}

func (c Client) ReverseGeocode(ll LatLng, opts *ReverseGeocodeOpts) ([]GeocodeResult, error) {
	var r geocodeResponse
	if err := c.doDecode(baseURL+reversegeocode(ll, opts), &r); err != nil {
		return nil, err
	}
	if r.Status != StatusOK {
		return nil, APIError{r.Status, ""}
	}
	return r.Results, nil
}

type ReverseGeocodeOpts struct {
	Language      string
	ResultTypes   []string
	LocationTypes []string
}

func (r *ReverseGeocodeOpts) update(p url.Values) {
	if r == nil {
		return
	}
	if r.ResultTypes != nil {
		p.Set("result_type", strings.Join(r.ResultTypes, "|"))
	}
	if r.LocationTypes != nil {
		p.Set("location_type", strings.Join(r.LocationTypes, "|"))
	}
}

func reversegeocode(ll LatLng, opts *ReverseGeocodeOpts) string {
	p := url.Values{}
	p.Set("latlng", ll.String())
	opts.update(p)
	return "geocode/json?" + p.Encode()
}
