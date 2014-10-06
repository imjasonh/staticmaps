package maps

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	MapFormatPNG         = "png"
	MapFormatPNG32       = "png32"
	MapFormatGIF         = "gif"
	MapFormatJPG         = "jpg"
	MapFormatJPGBaseline = "jpg-baseline"

	MapTypeRoadmap   = "roadmap"
	MapTypeSatellite = "satellite"
	MapTypeTerrain   = "terrain"
	MapTypeHybrid    = "hybrid"
)

func (c Client) StaticMap(s Size, opts *StaticMapOpts) (io.ReadCloser, error) {
	resp, err := c.do(baseURL + staticmap(s, opts))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, HTTPError{resp}
	}
	return resp.Body, nil
}

func staticmap(s Size, opts *StaticMapOpts) string {
	p := url.Values{}
	p.Set("size", s.String())
	opts.update(p)
	return "staticmap?" + p.Encode()
}

type Size struct {
	H, W int
}

func (s Size) String() string {
	return fmt.Sprintf("%dx%d", s.H, s.W)
}

type StaticMapOpts struct {
	Center                            *Location
	Zoom, Scale                       int
	Format, MapType, Language, Region string
	Markers                           []Markers
	Paths                             []Path
	Visible                           []Location
	Styles                            []Style
}

type Markers struct {
	Size       string // tiny, mid, small
	Color      string // hex color, or certain color names
	Label      string
	IconURL    string
	HideShadow bool
	Locations  []Location
}

func (m Markers) encode() string {
	s := []string{}
	if m.Size != "" {
		s = append(s, "size:"+m.Size)
	}
	if m.Color != "" {
		s = append(s, "color:"+m.Color)
	}
	if m.Label != "" {
		s = append(s, "label:"+m.Label)
	}
	if m.IconURL != "" {
		s = append(s, "icon:"+m.IconURL)
	}
	if m.HideShadow {
		s = append(s, "shadow:false")
	}
	style := strings.Join(s, "|")
	if style != "" {
		style += "|"
	}
	return style + encodeLocations(m.Locations)
}

type Path struct {
	Weight    int
	Color     string // with opacity, or certain color names
	FillColor string // with opacity, or certain color names
	Geodesic  bool
	Polyline  string
	Locations []Location
}

func (p Path) encode() string {
	s := []string{}
	if p.Weight != 0 {
		s = append(s, fmt.Sprintf("weight:%d", p.Weight))
	}
	if p.Color != "" {
		s = append(s, "color:"+p.Color)
	}
	if p.FillColor != "" {
		s = append(s, "fillcolor:"+p.FillColor)
	}
	if p.Geodesic {
		s = append(s, "geodesic:true")
	}
	style := strings.Join(s, "|")
	if style != "" {
		style += "|"
	}
	if p.Polyline != "" {
		return style + p.Polyline
	}
	return style + encodeLocations(p.Locations)
}

type Style struct {
	Feature string // TODO enum with sub-options
	Element string // TODO enum
	Rules   []StyleRule
}

func (t Style) encode() string {
	s := []string{}
	if t.Feature != "" {
		s = append(s, "feature:"+t.Feature)
	}
	if t.Element != "" {
		s = append(s, "element:"+t.Element)
	}
	for _, r := range t.Rules {
		s = append(s, r.encode())
	}
	return strings.Join(s, "|")
}

type StyleRule struct {
	Hue              string   // rgb color
	Lightness        float64  // -100 to 100
	Saturation       float64  // -100 to 100
	Gamma            *float64 // .01 to 10, default 1
	InverseLightness *bool
	Visibility       string // on, off, simplified
}

func (r StyleRule) encode() string {
	s := []string{}
	if r.Hue != "" {
		s = append(s, "hue:"+r.Hue)
	}
	if r.Lightness != 0 {
		s = append(s, fmt.Sprintf("lightness:%f", r.Lightness))
	}
	if r.Saturation != 0 {
		s = append(s, fmt.Sprintf("saturation:%f", r.Saturation))
	}
	if r.Gamma != nil {
		s = append(s, fmt.Sprintf("gamma:%f", *r.Gamma))
	}
	if r.InverseLightness != nil && *r.InverseLightness == false {
		s = append(s, "inverse_lightness:false")
	}
	if r.Visibility != "" {
		s = append(s, "visibility:"+r.Visibility)
	}
	return strings.Join(s, "|")
}

func (so *StaticMapOpts) update(p url.Values) {
	if so == nil {
		return
	}
	if so.Center != nil {
		p.Set("center", (*so.Center).Location())
	}
	if so.Zoom != 0 {
		p.Set("zoom", fmt.Sprintf("%d", so.Zoom))
	}
	if so.Scale != 0 {
		p.Set("scale", fmt.Sprintf("%d", so.Scale))
	}
	if so.Format != "" {
		p.Set("format", so.Format)
	}
	if so.MapType != "" {
		p.Set("maptype", so.MapType)
	}
	if so.Language != "" {
		p.Set("language", so.Language)
	}
	if so.Region != "" {
		p.Set("region", so.Region)
	}
	for _, m := range so.Markers {
		p.Add("markers", m.encode())
	}
	for _, path := range so.Paths {
		p.Add("path", path.encode())
	}
	if so.Visible != nil {
		p.Set("visible", encodeLocations(so.Visible))
	}
	for _, s := range so.Styles {
		p.Add("style", s.encode())
	}
}
