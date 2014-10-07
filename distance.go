package maps

import (
	"fmt"
	"net/url"
	"time"
)

func (c Client) DistanceMatrix(orig, dest []Location, opts *DistanceMatrixOpts) (*DistanceMatrixResponse, error) {
	var d DistanceMatrixResponse
	if err := c.doDecode(baseURL+distancematrix(orig, dest, opts), &d); err != nil {
		return nil, err
	}
	return &d, nil
}

func distancematrix(orig, dest []Location, opts *DistanceMatrixOpts) string {
	p := url.Values{}
	p.Set("origins", encodeLocations(orig))
	p.Set("destinations", encodeLocations(dest))
	opts.update(p)
	return "distancematrix/json?" + p.Encode()
}

type DistanceMatrixOpts struct {
	Language, Units string
	Mode, Avoid     *string
	DepartureTime   time.Time
}

func (o *DistanceMatrixOpts) update(p url.Values) {
	if o == nil {
		return
	}
	if o.Mode != nil {
		p.Set("mode", *o.Mode)
	}
	if o.Language != "" {
		p.Set("language", o.Language)
	}
	if o.Avoid != nil {
		p.Set("avoid", *o.Avoid)
	}
	if o.Units != "" {
		p.Set("units", o.Units)
	}
	if !o.DepartureTime.IsZero() {
		p.Set("departure_time", fmt.Sprintf("%d", o.DepartureTime.Unix()))
	}
}

type DistanceMatrixResponse struct {
	Status               string   `json:"status"`
	OriginAddresses      []string `json:"origin_addresses"`
	DestinationAddresses []string `json:"destination_addresses"`
	Rows                 []struct {
		Elements []struct {
			Status   string   `json:"status"`
			Duration Duration `json:"duration"`
			Distance Distance `json:"distance"`
		} `json:"elements"`
	} `json:"rows"`
}
