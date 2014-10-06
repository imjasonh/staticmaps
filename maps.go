package maps

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const baseURL = "https://maps.googleapis.com/maps/api/"

type Client struct {
	Transport http.RoundTripper
	Key       string
}

func NewClient(key string) Client {
	return Client{Key: key, Transport: http.DefaultTransport}
}

func (c Client) do(url string, r interface{}) error {
	cl := &http.Client{Transport: c.Transport}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if c.Key != "" {
		q := req.URL.Query()
		q.Set("key", c.Key)
		req.URL.RawQuery = q.Encode()
	}

	resp, err := cl.Get(url)
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

type HTTPError struct {
	Response *http.Response
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("http error %d", e.Response.StatusCode)
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

type Route []Location

func (r Route) String() string {
	s := make([]string, len(r))
	for i, l := range r {
		s[i] = l.Location()
	}
	return strings.Join(s, "|")
}
