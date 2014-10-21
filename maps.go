package maps

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	baseURL = "https://maps.googleapis.com/maps/api/"

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

// Client provides methods to make requests of various Maps APIs.
type Client struct {
	// The underlying http.RoundTripper to use. If this is not specified, http.DefaultTransport will be used.
	Transport http.RoundTripper

	// Key is the API key to use to authorize requests.
	Key string

	// ClientID is the Google Maps API for Work client ID to use.
	//
	// See https://developers.google.com/maps/documentation/business/webservices/auth
	ClientID string

	// PrivateKey is the base64-encoded private key to use for Google Maps API for Work requests.
	//
	// See https://developers.google.com/maps/documentation/business/webservices/auth
	PrivateKey string
}

// New returns a Client using the specified API key and RoundTripper.
func New(key string, rt http.RoundTripper) Client {
	return Client{Key: key, Transport: rt}
}

// NewForWork returns a Client using the specified Google Maps API for Work client ID and signature, and RoundTripper.
//
// See https://developers.google.com/maps/documentation/business/
func NewForWork(clientID, privateKey string, rt http.RoundTripper) Client {
	return Client{ClientID: clientID, PrivateKey: privateKey, Transport: rt}
}

func (c Client) do(url string) (*http.Response, error) {
	t := c.Transport
	if t == nil {
		t = http.DefaultTransport
	}
	cl := &http.Client{Transport: &backoff{
		Transport: t,
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
		q.Set("client", c.ClientID)
	}
	enc := q.Encode()
	if c.PrivateKey != "" {
		sig, err := c.genSig(req.URL.Path, enc)
		if err != nil {
			return nil, err
		}
		enc += "&signature=" + sig
	}
	req.URL.RawQuery = enc
	return cl.Get(url)
}

// See https://developers.google.com/maps/documentation/business/webservices/auth
func (c Client) genSig(path, query string) (string, error) {
	toSign := path + "?" + query
	decodedKey, err := base64.URLEncoding.DecodeString(c.PrivateKey)
	if err != nil {
		return "", err
	}
	d := hmac.New(sha1.New, decodedKey)
	if _, err := d.Write([]byte(toSign)); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(d.Sum(nil)), nil
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
	// Response is the http.Response returned from the API request.
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

// Location represents the general concept of a location in various methods.
type Location interface {
	Location() string
}

// LatLng represents a Location that is identified by its longitude and latitude.
type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Location returns the latitude/longitude pair as a comma-separated string.
func (ll LatLng) Location() string {
	return fmt.Sprintf("%f,%f", ll.Lat, ll.Lng)
}

func (ll LatLng) String() string {
	return ll.Location()
}

// Address represents a Location that is identified by its name or address, e.g., "New York, NY" or "111 8th Ave, NYC"
type Address string

// Location returns the Address as a string.
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

// Float64 returns a pointer to the given float64 value.
//
// This is a convenience method since some options structs take a *float64.
func Float64(f float64) *float64 {
	return &f
}
