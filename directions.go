package maps

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	AvoidTolls    = "tolls"
	AvoidHighways = "highways"
	AvoidFerries  = "ferries"

	ModeDriving   = "driving"
	ModeWalking   = "walking"
	ModeTransit   = "transit"
	ModeBicycling = "bicycling"

	UnitMetric   = "metric"
	UnitImperial = "imperial"

	StatusOK                   = "OK"
	StatusNotFound             = "NOT_FOUND"
	StatusZeroResults          = "ZERO_RESULTS"
	StatusMaxWaypointsExceeded = "MAX_WAYPOINTS_EXCEEDED"
	StatusInvalidRequest       = "INVALID_REQUEST"
	StatusRequestDenied        = "REQUEST_DENIED"
	StatusUnknownError         = "UNKNOWN_ERROR"
)

// TODO: via_waypoint not documented

func (c Client) Directions(orig, dest Location, opts *DirectionsOpts) (*DirectionsResponse, error) {
	var d DirectionsResponse
	if err := c.do(baseURL+directions(orig, dest, opts), &d); err != nil {
		return nil, err
	}
	return &d, nil
}

func directions(orig, dest Location, opts *DirectionsOpts) string {
	p := url.Values{}
	p.Set("origin", orig.Location())
	p.Set("destination", dest.Location())
	opts.update(p)
	return "directions/json?" + p.Encode()
}

type DirectionsOpts struct {
	Mode          string
	Waypoints     Route
	Alternatives  bool
	Avoid         []string
	Language      string
	Units         string
	Region        string
	DepartureTime time.Time
	ArrivalTime   time.Time
}

func (do *DirectionsOpts) update(p url.Values) {
	if do == nil {
		return
	}
	if do.Mode != "" {
		p.Set("mode", do.Mode)
	}
	if do.Waypoints != nil {
		p.Set("waypoints", do.Waypoints.String())
	}
	if do.Alternatives {
		p.Set("alternatives", "true")
	}
	if do.Avoid != nil {
		p.Set("avoid", strings.Join(do.Avoid, "|"))
	}
	if do.Language != "" {
		p.Set("language", do.Language)
	}
	if do.Units != "" {
		p.Set("units", do.Units)
	}
	if do.Region != "" {
		p.Set("region", do.Region)
	}
	if !do.DepartureTime.IsZero() {
		p.Set("departure_time", fmt.Sprintf("%d", do.DepartureTime.Unix()))
	}
	if !do.ArrivalTime.IsZero() {
		p.Set("arrival_time", fmt.Sprintf("%d", do.ArrivalTime.Unix()))
	}
}

type DirectionsResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
	Routes       []struct {
		Summary string `json:"summary"`
		Legs    []struct {
			Duration          *Duration `json:"duration"`
			DurationInTraffic *Duration `json:"duration_in_traffic"`
			Distance          *Distance `json:"distance"`
			ArrivalTime       Time      `json:"arrival_time"`
			DepartureTime     Time      `json:"departure_time"`
			StartLocation     *LatLng   `json:"start_location"`
			EndLocation       *LatLng   `json:"end_location"`
			StartAddress      string    `json:"start_address"`
			EndAddress        string    `json:"end_address"`
			Steps             []Step    `json:"steps"`
		} `json:"legs"`
		Bounds           Bounds   `json:"bounds"`
		Copyrights       string   `json:"copyrights"`
		OverviewPolyline Polyline `json:"overview_polyline"`
		Warnings         []string `json:"warnings"`
		WaypointOrder    []int    `json:"waypoint_order"`
	} `json:"routes"`
}

type Time struct {
	Value    int64  `json:"value"`
	Text     string `json:"text"`
	TimeZone string `json:"time_zone"`
}

func (t Time) Time() (*time.Time, error) {
	l, err := time.LoadLocation(t.TimeZone)
	if err != nil {
		return nil, err
	}
	tm := time.Unix(t.Value, 0).In(l)
	return &tm, nil
}

type Step struct {
	TravelMode       string    `json:"travel_mode"` // TODO: enum
	StartLocation    *LatLng   `json:"start_location"`
	EndLocation      *LatLng   `json:"end_location"`
	Maneuver         string    `json:"maneuver"` // TODO: not documented?
	Polyline         *Polyline `json:"polyline"`
	Duration         *Duration `json:"duration"`
	Distance         *Distance `json:"distance"`
	HTMLInstructions string    `json:"html_instructions"`
	Steps            []Step    `json:"steps"` // sub-steps
	TransitDetails   *struct {
		ArrivalStop   string `json:"arrival_stop"`
		DepartureStop string `json:"departure_stop"`
		ArrivalTime   Time   `json:"arrival_time"`
		DepartureTime Time   `json:"departure_time"`
		Headsign      string `json:"headsign"`
		Headway       int64  `json:"headway"` // Seconds until next departure, TODO: to time.Duration?
		NumStops      int    `json:"num_stops"`
		Line          struct {
			Name      string `json:"name"`
			ShortName string `json:"short_name"`
			Color     string `json:"color"` // hex color, TODO: to image.Color?
			Agencies  []struct {
				Name  string `json:"name"`
				URL   string `json:"url"`
				Phone string `json:"phone"`
			} `json:"agencies"`
			URL       string `json:"url"`
			IconURL   string `json:"icon"`
			TextColor string `json:"text_color"` // hex color, TODO: to image.Color?
			Vehicle   *struct {
				Name    string `json:"name"`
				Type    string `json:"type"` // TODO: enum
				IconURL string `json:"icon"`
			} `json:"vehicle"`
		} `json:"line"`
	} `json:"transit_details"`
}

type Bounds struct {
	Northeast LatLng `json:"northeast"`
	Southwest LatLng `json:"southwest"`
}

type Duration struct {
	Value int64  `json:"value"`
	Text  string `json:"text"`
}

func (d Duration) Duration() time.Duration {
	return time.Duration(d.Value) * time.Second
}

type Distance struct {
	Value int64  `json:"value"` // meters
	Text  string `json:"text"`
}

type Polyline struct {
	Points string `json:"points"`
}

// TODO: methods to decode polyline points?
