package maps

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	// ComponentRoute matches long or short name of a route.
	ComponentRoute = "route"
	// ComponentLocality matches against most locality and sublocality types.
	ComponentLocality = "locality"
	// ComponentAdministrativeArea matches all the administrative area levels.
	ComponentAdministrativeArea = "administrative_area"
	// ComponentPostalCode matches postcal code and postal code prefix.
	ComponentPostalCode = "postal_code"
	// ComponentCountry matches a country name or a two-letter ISO 3166-1 country code.
	ComponentCountry = "country"

	// LocationTypeRooftop restricts the results to addresses for which we have location information accurate down to street address precision.
	LocationTypeRooftop = "ROOFTOP"
	// LocationTypeRangeInterpolated restricts the results to those that reflect an approximation (usually a road) interpolated between two precise points (such as intersections).
	LocationTypeRangeInterpolated = "RANGE_INTERPOLATED"
	// LocationTypeGeometricCenter restricts the results to geometric centers of a location such as a polyline (for example, a street) or polygon (region).
	LocationTypeGeometricCenter = "GEOMETRIC_CENTER"
	// LocationTypeApproximate restricts the results to those that are characterized as approximate.
	LocationTypeApproximate = "APPROXIMATE"

	// TODO: enums for https://developers.google.com/maps/documentation/geocoding/#Types
)

// Geocode requests conversion of an address to a latitude/longitude pair.
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

// GeocodeOpts defines options for Geocode requests.
type GeocodeOpts struct {
	// The address to geocode.
	Address Address

	// Components specifies a component filter for which you wish to obtain a geocode.
	//
	// Only results that match all filters will be returned. Filter values support the same methods of spelling correction and partial matching as other geocoding requests.
	//
	// See https://developers.google.com/maps/documentation/geocoding/#ComponentFiltering
	Components []Component

	// The language in which to return results.
	//
	// See https://developers.google.com/maps/faq#languagesupport
	Language string

	// Region specifies a region to bias results toward, specified as a ccTLD ("top-level domain" two-character value.
	//
	// This will influence, not fully restrict, results from the geocoder.
	//
	// See https://developers.google.com/maps/documentation/geocoding/#RegionCodes
	Region string

	// Bounds specifies the bounding box of the viewport within which to bias geocode results more prominently.
	//
	// This will influence, not fully restrict, results from the geocoder.
	Bounds *Bounds
}

// Component describes a single component filter.
//
// See https://developers.google.com/maps/documentation/geocoding/#ComponentFiltering
type Component struct {
	// Key is the type of component to filter on.
	//
	// Accepted values are ComponentRoute, ComponentLocality, ComponentAdministrativeArea, ComponentPostalCode and ComponentCountry.
	Key string

	// Value is the value of the filter to enforce.
	Value string
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

// GeocodeResult represents the geocoded result for the given input.
type GeocodeResult struct {
	// AddressComponents contains the separate address components.
	AddressComponents []struct {
		// LongName is the full text description or name of the address component returned by the geocoder.
		LongName string `json:"long_name"`
		// ShortName is the abbreviated textual name for the address component, if available.
		//
		// For example, an address component for the state of Alaska may have a LongName of "Alaska" and a ShortName of "AK" using the 2-letter postal abbreviation.
		ShortName string `json:"short_name"`

		// Types contains the types of address components.
		Types []string `json:"types"`
	} `json:"address_components"`

	// PostcodeLocalities denotes all the localities contained in a postal code.
	//
	// This is only present when the result is a postal code that contains multiple localities.
	PostcodeLocalities []string `json:"postcode_localities"`

	// FormattedAddress contains the human-readable address of this location.
	//
	// Often this is equivalent to the "postal address" which sometimes differs from country to country.
	FormattedAddress string `json:"formatted_address"`

	// Geometry contains information about the location.
	Geometry struct {
		// Location is the geocoded latitude and longitude of the location.
		Location LatLng `json:"location"`

		// LocationType specifies additional data about the specified location.
		//
		// Its value will be one of:
		// - "ROOFTOP" indicating that the returned result is a precise geocode for which we have location information accurate down to street address precision.
		// - "RANGE_INTERPOLATED" indicating that the returned result reflects and approximation (usually on a road) interpolated between two precise points (such as intersections). Interpolated results are generally returned when rooftop geocodes are unavailable for a street address.
		// - "GEOMETRIC_CENTER" indicating that the returned result is the geometric center of a result such as a polyline (for example, a street) or polygon (region).
		LocationType string `json:"location_type"`

		// Viewport contains the recommended viewport bounding box for displaying the returned result.
		Viewport Bounds `json:"viewport"`

		// Bounds, if included, contains the bounding box which can fully contain the returned result.
		//
		// Note that this may not match Viewport.
		Bounds Bounds `json:"bounds"`
	} `json:"geometry"`

	// Types indicates the types of the returned result.
	//
	// This contains a set of zero or more tags identifying the type of feature returned in a result.
	// For example, a geocode of "Chicago" returns "locality" which indicates that "Chicago" is a city,
	// and also returns "political" which indicates it is a political entity.
	Types []string `json:"types"`

	// PartialMatch indicates that the geocoder did not return an exact match for the original request,
	// though it was able to match part of the requested address. You may wish to examine the original
	// request for misspellings and/or an incomplete address.
	PartialMatch string `json:"partial_match"`
}

// ReverseGeocode requests conversion of a location to its nearest address.
//
// See https://developers.google.com/maps/documentation/geocoding/#ReverseGeocoding
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

// ReverseGeocodeOpts defines options for ReverseGeocode requests.
type ReverseGeocodeOpts struct {
	// The language in which to return results.
	//
	// See https://developers.google.com/maps/faq#languagesupport
	Language string

	// ResultTypes specifies filters on result types, e.g., "country", "street address"
	ResultTypes []string

	// LocationTypes specifies filters on location types.
	//
	// Accepted values are LocationTypeRooftop, LocationTypeRangeInterpolated, LocationTypeGeometricCenter and LocationTypeApproximate.
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
