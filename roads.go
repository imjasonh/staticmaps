package maps

import (
	"net/url"

	"code.google.com/p/go.net/context"
)

const roadsAPIBaseURL = "https://roads.googleapis.com/v1/"

func SnapToRoads(ctx context.Context, path []LatLng, opts *SnapToRoadsOpts) ([]SnappedPoint, error) {
	var d snapToRoadsResponse
	if err := doDecode(ctx, roadsAPIBaseURL+snapToRoads(path, opts), &d); err != nil {
		return nil, err
	}
	return d.SnappedPoints, nil
}

func snapToRoads(path []LatLng, opts *SnapToRoadsOpts) string {
	p := url.Values{}
	p.Set("path", encodeLatLngs(path))
	opts.update(p)
	return "snapToRoads?" + p.Encode()
}

// SnapToRoadsOpts defines options for SnapToRoads requests.
type SnapToRoadsOpts struct {
	// Whether to interpolate a path to include all points forming the full road-geometry.
	Interpolate bool
}

func (o SnapToRoadsOpts) update(p url.Values) {
	if o.Interpolate {
		p.Set("interpolate", "true")
	}
}

type snapToRoadsResponse struct {
	SnappedPoints []SnappedPoint `json:"snappedPoints"`
}

// SnappedPoint represents a location derived from a SnapToRoads request.
type SnappedPoint struct {
	// The latitude and longitude of the snapped location.
	Location SnappedLocation `json:"location"`

	// The index of the corresponding value in the original request.
	//
	// Each value in the request should map to a snapped value in the response. However, if you've set Interpolate to be true, then it's possible that the response will contain more coordinates than the request. Interpolated values will have this value set to nil.
	OriginalIndex *int `json:"originalIndex,omitempty"`

	// A unique identifier for a Place corresponding to a road segment.
	PlaceID string `json:"placeId"`
}

// SnappedLocation is equivalent to LatLng
//
// TODO(jasonhall): Reconcile this with LatLng which expects JSON fields lat/lng
type SnappedLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
