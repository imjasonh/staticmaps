package maps

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const baseURL = "https://maps.googleapis.com/maps/api/"

const (
	// StatusOK indicates the response contains a valid result.
	StatusOK = "OK"
	// StatusNotFound indicates at least one of the locations specified could not be geocoded.
	StatusNotFound = "NOT_FOUND"
	// StatusZeroResults indicates no route could be found between the origin and destination.
	StatusZeroResults = "ZERO_RESULTS"
	// StatusMaxWaypointsExceeded indicates that too many Waypoints were provided in the request.
	//
	// The maximum allowed waypoints is 8, plus the origin and destination. Google Maps API for Work clients may contain requests with up to 23 waypoints.
	StatusMaxWaypointsExceeded = "MAX_WAYPOINTS_EXCEEDED"
	// StatusInvalidRequest indicates that the provided request was invalid.
	StatusInvalidRequest = "INVALID_REQUEST"
	// StatusRequestDenied indicates that the service denied use of the service by your application.
	StatusRequestDenied = "REQUEST_DENIED"
	// StatusUnknownError indicates that the request could not be processed due to a server error. The request may succeed if you try again.
	StatusUnknownError = "UNKNOWN_ERROR"
	// StatusOverQueryLimit indicates that the service has received too many requests from your application within the allowed time period.
	StatusOverQueryLimit = "OVER_QUERY_LIMIT"
)

type Client struct {
	Transport                http.RoundTripper
	Key, ClientID, Signature string
}

func NewClient(key string) Client {
	return Client{Key: key, Transport: http.DefaultTransport}
}

func NewWorkClient(clientID, signature string) Client {
	return Client{ClientID: clientID, Signature: signature, Transport: http.DefaultTransport}
}

func (c Client) do(url string) (*http.Response, error) {
	cl := &http.Client{Transport: &backoff{
		Transport: &c.Transport,
	}}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	if c.Key != "" {
		q.Set("key", c.Key)
	}
	if c.ClientID != "" {
		q.Set("clientId", c.ClientID)
	}
	if c.Signature != "" {
		q.Set("signature", c.Signature)
	}
	req.URL.RawQuery = q.Encode()
	return cl.Get(url)
}

func (c Client) doDecode(url string, r interface{}) error {
	resp, err := c.do(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return HTTPError{resp}
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	return nil
}

// HTTPError indicates an error communicating with the API server, and includes the HTTP response returned from the server.
type HTTPError struct {
	Response *http.Response
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("http error %d", e.Response.StatusCode)
}

// APIError indicates a failure response from the API server, even though a successful HTTP response was returned.
// Its Status and Message fields can be consulted for more information about the specific error conditions.
type APIError struct {
	Status, Message string
}

func (e APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("API error %q: %s", e.Status, e.Message)
	}
	return fmt.Sprintf("API Error %q", e.Status)
}

type Location interface {
	Location() string
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (ll LatLng) Location() string {
	return fmt.Sprintf("%f,%f", ll.Lat, ll.Lng)
}

func (ll LatLng) String() string {
	return ll.Location()
}

type Address string

func (a Address) Location() string {
	return string(a)
}

func encodeLocations(ls []Location) string {
	s := make([]string, len(ls))
	for i, l := range ls {
		s[i] = l.Location()
	}
	return strings.Join(s, "|")
}

func encodeLatLngs(ll []LatLng) string {
	s := make([]string, len(ll))
	for i, l := range ll {
		s[i] = l.Location()
	}
	return strings.Join(s, "|")
}

func Float64(f float64) *float64 {
	return &f
}

func String(s string) *string {
	return &s
}
