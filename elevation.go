package maps

import (
	"fmt"
	"net/url"

	"golang.org/x/net/context"
)

// Elevation requests elevation data for a series of locations.
//
// See https://developers.google.com/maps/documentation/elevation/
func Elevation(ctx context.Context, ll []LatLng) ([]ElevationResult, error) {
	var r elevationResponse
	if err := doDecode(ctx, baseURL+elevation(ll), &r); err != nil {
		return nil, err
	}
	if r.Status != StatusOK {
		return nil, APIError{r.Status, ""}
	}
	return r.Results, nil
}

// ElevationPolyline requests elevation data for a series of locations as specified as an encoded polyline.
//
// See https://developers.google.com/maps/documentation/elevation/#Locations
func ElevationPolyline(ctx context.Context, p string) ([]ElevationResult, error) {
	var r elevationResponse
	if err := doDecode(ctx, baseURL+elevationpoly(p), &r); err != nil {
		return nil, err
	}
	if r.Status != StatusOK {
		return nil, APIError{r.Status, ""}
	}
	return r.Results, nil
}

// ElevationPath requests elevation data for a number of samples along a path described as a series of locations.
func ElevationPath(ctx context.Context, ll []LatLng, samples int) ([]ElevationResult, error) {
	var r elevationResponse
	if err := doDecode(ctx, baseURL+elevationpath(ll, samples), &r); err != nil {
		return nil, err
	}
	if r.Status != StatusOK {
		return nil, APIError{r.Status, ""}
	}
	return r.Results, nil
}

// ElevationPathPoly requests elevation data for a number of samples along a path described as a series of locations specified as an encoded polyline.
func ElevationPathPoly(ctx context.Context, p string, samples int) ([]ElevationResult, error) {
	var r elevationResponse
	if err := doDecode(ctx, baseURL+elevationpathpoly(p, samples), &r); err != nil {
		return nil, err
	}
	if r.Status != StatusOK {
		return nil, APIError{r.Status, ""}
	}
	return r.Results, nil
}

func elevation(ll []LatLng) string {
	p := url.Values{}
	p.Set("locations", encodeLatLngs(ll))
	return "elevation/json?" + p.Encode()
}

func elevationpoly(poly string) string {
	p := url.Values{}
	p.Set("locations", "enc:"+poly)
	return "elevation/json?" + p.Encode()
}

func elevationpath(ll []LatLng, s int) string {
	p := url.Values{}
	p.Set("path", encodeLatLngs(ll))
	p.Set("samples", fmt.Sprintf("%d", s))
	return "elevation/json?" + p.Encode()
}

func elevationpathpoly(poly string, s int) string {
	p := url.Values{}
	p.Set("path", "enc:"+poly)
	p.Set("samples", fmt.Sprintf("%d", s))
	return "elevation/json?" + p.Encode()
}

type elevationResponse struct {
	Results []ElevationResult `json:"results"`
	Status  string            `json:"status"`
}

// ElevationResult describes elevation data.
type ElevationResult struct {
	// Elevation is elevation of the location in meters.
	Elevation float64 `json:"elevation"`

	// Location is the location of the location being described.
	Location LatLng `json:"location"`

	// Resolution is the maximum distance between data points from which the elevation was interpolated, in meters.
	//
	// This property will be zero if the resolution is not known. Note that elevation data becomes
	// more coarse (larger Resolution values) when multiple points are passed. To obtain the most
	// accurate elevation value for a point, it should be queried independently.
	Resolution float64 `json:"resolution"`
}
