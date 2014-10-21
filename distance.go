package maps

import (
	"fmt"
	"net/url"
	"time"

	"code.google.com/p/go.net/context"
)

// DistanceMatrix requests travel distance and time for a matrix of origins and destinations.
//
// See https://developers.google.com/maps/documentation/distancematrix/
func DistanceMatrix(ctx context.Context, orig, dest []Location, opts *DistanceMatrixOpts) (*DistanceMatrixResult, error) {
	var d distanceResponse
	if err := doDecode(ctx, baseURL+distancematrix(orig, dest, opts), &d); err != nil {
		return nil, err
	}
	if d.Status != StatusOK {
		return nil, APIError{d.Status, ""}
	}
	return &d.DistanceMatrixResult, nil
}

func distancematrix(orig, dest []Location, opts *DistanceMatrixOpts) string {
	p := url.Values{}
	p.Set("origins", encodeLocations(orig))
	p.Set("destinations", encodeLocations(dest))
	opts.update(p)
	return "distancematrix/json?" + p.Encode()
}

// DistanceMatrixOpts defines options for DistanceMatrix requsts.
type DistanceMatrixOpts struct {
	// The language in which to return results.
	//
	// See https://developers.google.com/maps/faq#languagesupport
	Language string

	// Specifies the unit system to use when displaying results.
	//
	// Accepted values are UnitMetric and UnitImperial.
	//
	// See https://developers.google.com/maps/documentation/distancematrix/#unit_systems
	Units string

	// Specifies the mode of transport to use when calculating directions.
	//
	// Accepted values are ModeDriving (the default), ModeWalking and ModeBicycling.
	Mode string

	// Indicates that the calculated route(s) should avoid the indicated features.
	//
	// Accepted values are AvoidTolls, AvoidHighways or AvoidFerries
	//
	// See https://developers.google.com/maps/documentation/distancematrix/#Restrictions
	Avoid string

	// Specifies the desired time of departure.
	//
	// The departure time can be specified for Google Maps API for Work clients to receive trip duration considering current traffic conditions.
	DepartureTime time.Time
}

func (o *DistanceMatrixOpts) update(p url.Values) {
	if o == nil {
		return
	}
	if o.Mode != "" {
		p.Set("mode", o.Mode)
	}
	if o.Language != "" {
		p.Set("language", o.Language)
	}
	if o.Avoid != "" {
		p.Set("avoid", o.Avoid)
	}
	if o.Units != "" {
		p.Set("units", o.Units)
	}
	if !o.DepartureTime.IsZero() {
		p.Set("departure_time", fmt.Sprintf("%d", o.DepartureTime.Unix()))
	}
}

type distanceResponse struct {
	Status string `json:"status"`
	DistanceMatrixResult
}

// DistanceMatrixResult describes the matrix of distances between the requested origins and destinations.
type DistanceMatrixResult struct {
	// OriginAddresses contains the addresses returned by the API from your original request.
	//
	// These are formatted by the geocoder and localized according to the requested Language.
	OriginAddresses []string `json:"origin_addresses"`

	// DestinationAddresses contains the addresses returned by the API from your original request.
	//
	// These are formatted by the geocoder and localized according to the requested Language.
	DestinationAddresses []string `json:"destination_addresses"`

	// Rows describes rows in the matrix.
	Rows []struct {
		// Elements describes columns in the matrix.
		Elements []struct {
			// Status indicates the status of the request, and will be one of StatusOK, StatusNotFound or StatusZeroResults.
			Status string `json:"status"`

			// Duration indicates the total duration of this journey.
			Duration Duration `json:"duration"`

			// Distance indicates the total distance of this journey.
			Distance Distance `json:"distance"`
		} `json:"elements"`
	} `json:"rows"`
}
