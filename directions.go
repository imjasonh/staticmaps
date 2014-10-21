package maps

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"code.google.com/p/go.net/context"
)

const (
	// AvoidTolls indicates that the route should avoid toll roads.
	AvoidTolls = "tolls"
	// AvoidHighways indicates that the route should avoid highways.
	AvoidHighways = "highways"
	// AvoidFerries indicates that the route should avoid ferries.
	AvoidFerries = "ferries"

	// ModeDriving indicates that driving directions are requested.
	ModeDriving = "driving"
	// ModeWalking indicates that walking directions are requested.
	ModeWalking = "walking"
	// ModeTransit indicates that transit directions are requested.
	ModeTransit = "transit"
	// ModeBicycling indicates that bicycling directions are requested.
	ModeBicycling = "bicycling"

	// UnitMetric indicates that the results should be stated in metric units.
	UnitMetric = "metric"
	// UnitImperial indicates that the results should be stated in imperial units.
	UnitImperial = "imperial"
)

// TODO: via_waypoint not documented

// Directions requests routes between orig and dest Locations.
//
// See https://developers.google.com/maps/documentation/directions/
func (c Client) Directions(ctx context.Context, orig, dest Location, opts *DirectionsOpts) ([]Route, error) {
	var d directionsResponse
	if err := c.doDecode(ctx, baseURL+directions(orig, dest, opts), &d); err != nil {
		return nil, err
	}
	if d.Status != StatusOK {
		return nil, APIError{d.Status, d.ErrorMessage}
	}
	return d.Routes, nil
}

func directions(orig, dest Location, opts *DirectionsOpts) string {
	p := url.Values{}
	p.Set("origin", orig.Location())
	p.Set("destination", dest.Location())
	opts.update(p)
	return "directions/json?" + p.Encode()
}

// DirectionsOpts defines options for Directions requests.
type DirectionsOpts struct {
	// Waypoints alter a route by routing it through the specified location(s).
	//
	// Waypoints are only supported for driving, walking and bicycling directions.
	//
	// See https://developers.google.com/maps/documentation/directions/#Waypoints
	Waypoints []Location

	// If true and if waypoints are specified, the Directions API may optimize the provided route by rearranging the waypoints in a more efficient order.
	OptimizeWaypoints bool

	// If true, multiple Routes may be returned. This may increase the response time from the server.
	Alternatives bool

	// Avoid indicates that the calculated route(s) should avoid the indicated features.
	//
	// Accepted values are AvoidTolls, AvoidHighways and AvoidFerries
	//
	// See https://developers.google.com/maps/documentation/directions/#Restrictions
	Avoid []string

	// Specifies the mode of transport to use when calculating directions.
	//
	// Accepted values are ModeDriving (the default), ModeWalking, ModeTransit and ModeBicycling.
	// If ModeTransit is specified, either DepartureTime or ArrivalTime must also be specified.
	//
	// See https://developers.google.com/maps/documentation/directions/#TravelModes
	Mode string

	// The language in which to return results.
	//
	// See https://developers.google.com/maps/faq#languagesupport
	Language string

	// Specifies the unit system to use when displaying results.
	//
	// Accepted values are UnitMetric and UnitImperial.
	//
	// See https://developers.google.com/maps/documentation/directions/#UnitSystems
	Units string

	// The region code, specified as a ccTLD ("top-level domain") two-character value.
	//
	// See https://developers.google.com/maps/documentation/directions/#RegionBiasing
	Region string

	// Specifies the desired time of departure.
	//
	// The departure time can be specified for transit directions, or for Google Maps API for Work
	// clients to receive trip duration considering current traffic conditions.
	DepartureTime time.Time

	// Specifies the desired time of arrival for transit directions.
	ArrivalTime time.Time
}

func (do *DirectionsOpts) update(p url.Values) {
	if do == nil {
		return
	}
	if do.Mode != "" {
		p.Set("mode", do.Mode)
	}
	if do.Waypoints != nil {
		v := ""
		if do.OptimizeWaypoints {
			v = "optimize:true|"
		}
		p.Set("waypoints", v+encodeLocations(do.Waypoints))
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

type directionsResponse struct {
	Status       string  `json:"status"`
	ErrorMessage string  `json:"error_message"`
	Routes       []Route `json:"routes"`
}

// Route describes a possible route between the requested origin and destination.
type Route struct {
	// Summary contains a short textual description of the route, suitable for naming and disambiguating the route from alternatives.
	Summary string `json:"summary"`

	// Legs contains information about the legs of a route, between two locations within the route.
	// A separate leg will be present for each waypoint or destination specified. If no waypoints were
	// requested, this will contain one element.
	//
	// See https://developers.google.com/maps/documentation/directions/#Legs
	Legs []struct {
		// Duration indicates the total duration of this leg.
		Duration *Duration `json:"duration"`

		// DurationInTraffic indicates the total duration of this leg, taking into account current traffic conditions.
		//
		// It will only be included if all of the following are true:
		// - The directions request includes a DepartureTime parameter set to a value within a few minutes of the current time.
		// - The request is made using a Google Maps API for Work client.
		// - Traffic conditions are available for the requested route.
		// - The directions request doesnot include stopover waypoints.
		DurationInTraffic *Duration `json:"duration_in_traffic"`

		// Distance is the total distance covered by this leg.
		Distance *Distance `json:"distance"`

		// ArrivalTime indicates the estimated time of arrival for this leg. This is only included for transit directions.
		ArrivalTime Time `json:"arrival_time"`

		// DepartureTime indicates the estimated time of arrival for this leg. This is only included for transit directions.
		DepartureTime Time `json:"departure_time"`

		// StartLocation indicates the origin of this leg.
		//
		// Because the Directions API calculates directions between locations by using the nearest transportation option (usually a road)
		// at the start and end points, StartLocation may be different than the provided origin of this leg if, for example, a road is not near the origin.
		StartLocation *LatLng `json:"start_location"`

		// EndLocation indicates the destination of this leg.
		//
		// Because the Directions API calculates directions between locations by using the nearest transportation option (usually a road)
		// at the start and end points, EndLocation may be different than the provided destination of this leg if, for example, a road is not near the destination.
		EndLocation *LatLng `json:"end_location"`

		// StartAddress contains the human-readable address (typically a street address) reflecting the StartLocation of this leg.
		StartAddress string `json:"start_address"`

		// EndAddress contains the human-readable address (typically a street address) reflecting the EndLocation of this leg.
		EndAddress string `json:"end_address"`

		// Steps describes each step of the leg of the journey.
		//
		// See https://developers.google.com/maps/documentation/directions/#Steps
		Steps []Step `json:"steps"`
	} `json:"legs"`

	// Bounds describes the viewport bounding box of the OverviewPolyline
	Bounds Bounds `json:"bounds"`

	// Copyrights contains the copyright text to be displayed for this route. You must display this information yourself.
	Copyrights string `json:"copyrights"`

	// OverviewPolyline contains the encoded polyline representation of the route. This polyline is an approximate (smoothed) path of the resulting directions.
	OverviewPolyline Polyline `json:"overview_polyline"`

	// Warnings contains any warnings to display to the user when showing these directions. You must display these warnings yourself.
	Warnings []string `json:"warnings"`

	// WaypointOrder indicates the order of any waypoints in the calculated route. The waypoints may be reordered if OptimizeWaypoints was specified.
	WaypointOrder []int `json:"waypoint_order"`
}

// Time represents a time.
type Time struct {
	// Value indicates the number of seconds since epoch of this time.
	Value int64 `json:"value"`

	// Text contains a human-readable time, displayed in the time zone of the transit stop.
	Text string `json:"text"`

	// TimeZone indicates the time zone of this station, e.g., "America/New_York"
	TimeZone string `json:"time_zone"`
}

// Time returns the Time as a time.Time, in the specified TimeZone.
func (t Time) Time() (*time.Time, error) {
	l, err := time.LoadLocation(t.TimeZone)
	if err != nil {
		return nil, err
	}
	tm := time.Unix(t.Value, 0).In(l)
	return &tm, nil
}

// Step describes a step of a leg of a journey.
//
// See https://developers.google.com/maps/documentation/directions/#Steps
type Step struct {
	TravelMode string `json:"travel_mode"` // TODO: not documented? (all caps)

	// StartLocation contains the location of the starting point of this step.
	StartLocation *LatLng `json:"start_location"`

	// EndLocation contains the location of the ending point of this step.
	EndLocation *LatLng `json:"end_location"`

	Maneuver string `json:"maneuver"` // TODO: not documented?

	// Contains an encoded polyline representation of the step. This is an approximate (smoothed) path of the step.
	Polyline *Polyline `json:"polyline"`

	// Duration indicates the total duration of this step.
	Duration *Duration `json:"duration"`

	// Distance is the total distance covered by this step.
	Distance *Distance `json:"distance"`

	// HTMLInstructions contains formatted instructions for this step, presented as an HTML text string.
	HTMLInstructions string `json:"html_instructions"`

	// Steps contains detailed directions for walking or driving steps in transit directions. These are only available when the requested Mode is ModeTransit
	Steps []Step `json:"steps"`

	// TransitDetails contains transit-specific information. This is only included when the requested Mode is ModeTransit
	//
	// See https://developers.google.com/maps/documentation/directions/#TransitDetails
	TransitDetails *struct {
		// ArrivalStop contains information about the arrival stop for this part of the trip.
		ArrivalStop Stop `json:"arrival_stop"`

		// DepartureStop contains information about the departure stop for this part of the trip.
		DepartureStop Stop `json:"departure_stop"`

		// ArrivalTime indicates the arrival time for this leg of the journey.
		ArrivalTime Time `json:"arrival_time"`

		// DepartureTime indicates the departure time for this leg of the journey.
		DepartureTime Time `json:"departure_time"`

		// Headsign specifies the direction in which to travel on this line, as it is marked on the vehicle or at the departure stop. This will often be the terminus station.
		Headsign string `json:"headsign"`

		// Headway specifies the expected number of seconds between departure from the same stop at this time.
		//
		// For example, with a Headway of 600, you would expect a ten minute way if you should miss your bus.
		Headway int64 `json:"headway"` // TODO: to time.Duration?

		// NumStops indicates the number of stops in this step, counting the arrival stop, but not the departure stop.
		//
		// For example, if your directions involve leaving Stop A, passing through Stops B and C, and arriving at Stop D, NumStops will be 3.
		NumStops int `json:"num_stops"`

		// Line contains information about the transit line used in this step.
		Line struct {
			// Name is the full name of this transit line, e.g., "7 Avenue Express"
			Name string `json:"name"`

			// ShortName is the short name of this transit line. This will normally be a line number, such as "M7" or "355"
			ShortName string `json:"short_name"`

			// Color is the color commonly used in signage for this transit line, specified as a hex string such as "#FF0033"
			Color string `json:"color"` // TODO: to image.Color?

			// Agencies contains information about the operator of the line.
			//
			// You must display the names and URLs of the transit agencies servicing the trip results.
			Agencies []struct {
				// Name is the name of the transit agency.
				Name string `json:"name"`

				// URL is the URL for the transit agency.
				URL string `json:"url"`

				// Phone is the phone number of the transit agency.
				Phone string `json:"phone"`
			} `json:"agencies"`

			// URL is the URL for this transit line as provided by the transit agency.
			URL string `json:"url"`

			// IconURL is the URL for the icon associated with this line.
			IconURL string `json:"icon"`

			// TextColor is the color of text commonly used for signage of this line, specified as a hex string such as "#FF0033"
			TextColor string `json:"text_color"` // TODO: to image.Color?

			// Vehicle contains the type of vehicle used on this line.
			Vehicle *struct {
				// Name is the name of the vehicle used on this line, e.g., "Subway"
				Name string `json:"name"`

				// Type is the type of the vehicle that runs on this line.
				//
				// See https://developers.google.com/maps/documentation/directions/#VehicleType
				Type string `json:"type"`

				// IconURL is the URL for an icon associated with this vehicle type.
				IconURL string `json:"icon"`
			} `json:"vehicle"`
		} `json:"line"`
	} `json:"transit_details"`
}

// Stop describes a transit station or stop.
type Stop struct {
	// Location is the location of the station/stop.
	Location LatLng `json:"location"`

	// Name is the name of the transit station/stop, e.g., "Union Square".
	Name string `json:"name"`
}

// Bounds describes a viewport bounding box.
type Bounds struct {
	// The northeasternmost point of the bounding box.
	Northeast LatLng `json:"northeast"`

	// The southwesternmost point of the bounding box.
	Southwest LatLng `json:"southwest"`
}

func (b Bounds) String() string {
	return fmt.Sprintf("%s|%s", b.Northeast, b.Southwest)
}

// Duration describes an amount of time for a leg or step.
type Duration struct {
	// Value indicates the duration in seconds.
	Value int64 `json:"value"`

	// Text contains a human-readable representation of the duration.
	Text string `json:"text"`
}

// Duration translates the Duration into a time.Duration.
func (d Duration) Duration() time.Duration {
	return time.Duration(d.Value) * time.Second
}

// Distance describes a distance for a leg or step.
type Distance struct {
	// Value is the distance in meters
	Value int64 `json:"value"`

	// Text contains a human-readable representation of the distance, displayed in the requested unit and language, or in the unit and language used at the origin.
	Text string `json:"text"`
}

// Polyline contains data describing an encoded polyline.
type Polyline struct {
	// Points is an encoded polyline describing some path.
	//
	// See https://developers.google.com/maps/documentation/utilities/polylinealgorithm
	Points string `json:"points"`
}

// TODO: methods to decode polyline points?
