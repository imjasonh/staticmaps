package maps

import (
	"fmt"
	"net/url"
)

func (c Client) Elevation(ll []LatLng) (*ElevationResponse, error) {
	var r ElevationResponse
	if err := c.doDecode(baseURL+elevation(ll), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (c Client) ElevationPolyline(p string) (*ElevationResponse, error) {
	var r ElevationResponse
	if err := c.doDecode(baseURL+elevationpoly(p), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (c Client) ElevationPath(ll []LatLng, samples int) (*ElevationResponse, error) {
	var r ElevationResponse
	if err := c.doDecode(baseURL+elevationpath(ll, samples), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (c Client) ElevationPathPoly(p string, samples int) (*ElevationResponse, error) {
	var r ElevationResponse
	if err := c.doDecode(baseURL+elevationpathpoly(p, samples), &r); err != nil {
		return nil, err
	}
	return &r, nil
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

type ElevationResponse struct {
	Results []struct {
		Elevation  float64 `json:"elevation"`
		Location   LatLng  `json:"location"`
		Resolution float64 `json:"resolution"`
	}
	Status string `json:"status"`
}
