package maps

import (
	"context"
	"fmt"
	"image/color"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	// MapFormatPNG requests the image in PNG format.
	MapFormatPNG = "png"
	// MapFormatPNG32 requests the image in PNG32 format.
	MapFormatPNG32 = "png32"
	// MapFormatGIF requests the image in GIF format.
	MapFormatGIF = "gif"
	// MapFormatJPG requests the image in JPG format.
	MapFormatJPG = "jpg"
	// MapFormatJPGBaseline requests the image in a non-progressive JPG format.
	MapFormatJPGBaseline = "jpg-baseline"

	// MapTypeRoadmap requests a standard roadmap image.
	MapTypeRoadmap = "roadmap"
	// MapTypeSatellite requests a satellite image.
	MapTypeSatellite = "satellite"
	// MapTypeTerrain requests a terrain image.
	MapTypeTerrain = "terrain"
	// MapTypeHybrid requests a hybrid of the satellite and roadmap image.
	MapTypeHybrid = "hybrid"

	// SizeTiny requests a small-sized marker.
	SizeTiny = "tiny"
	// SizeMid requests a mid-sized marker.
	SizeMid = "mid"
	// SizeLarge requests a large-sized marker.
	SizeLarge = "large"

	// VisibilityOn requests that the feature be shown.
	VisibilityOn = "on"
	// VisibilityOff requests that the feature not be shown.
	VisibilityOff = "off"
	// VisibilitySimplified requests that the feature be shown in a simplified manner.
	VisibilitySimplified = "simplified"
)

// StaticMap requests a static map image of a requested size.
func StaticMap(ctx context.Context, s Size, opts *StaticMapOpts) (io.ReadCloser, error) {
	resp, err := do(ctx, baseURL+staticmap(s, opts))
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

// Size specifies the height and width of the image to request in pixels.
//
// See https://developers.google.com/maps/documentation/staticmaps/#Imagesizes
type Size struct {
	H, W int
}

func (s Size) String() string {
	return fmt.Sprintf("%dx%d", s.W, s.H)
}

// StaticMapOpts defines options for StaticMap requests.
type StaticMapOpts struct {
	// Center specifies the location to place at the center of the image.
	Center Location

	// Zoom is the zoom level to request.
	//
	// The zoom level is between 0 (the lowest level, in which the entire world can be seen on one map)
	// and 21+ (down to streets and individual buildings).
	//
	// Each zoom level doubles the precision in both horizontal and vertical dimensions.
	//
	// See https://developers.google.com/maps/documentation/staticmaps/#Zoomlevels
	Zoom int

	// Scale requests a higher-resolution map image for use on high-density displays.
	//
	// A scale value of 2 will double the resulting image size. A scale value of 4 (only available to Google Maps API for Work clients) will quadruple it.
	Scale int

	// Format specifies the image format to request.
	//
	// Accepted values are MapFormatPNG (the default), MapFormatPNG32, MapFormatGIF, MapFormatJPG and MapFormatJPGBaseline.
	Format string

	// MapType specifies the map type to request.
	//
	// Accepted values are MapTypeRoadmap (the default), MapTypeSatellite, MapTypeTerrain and MapTypeHybrid.
	MapType string

	// The language in which to localize labels on the map.
	//
	// See https://developers.google.com/maps/faq#languagesupport
	Language string

	// Region defines the appropriate borders to display, based on geo-political sensitivities.
	//
	// Accepts a two-character ccTLD ("top-level domain") value.
	Region string

	// Markers defines one or more markers to attach to the image at specified locations.
	//
	// See https://developers.google.com/maps/documentation/staticmaps/#Markers
	Markers []Markers

	// Path defines one or more paths to attach to the image at specified locations.
	//
	// See https://developers.google.com/maps/documentation/staticmaps/#Paths
	Paths []Path

	// Visible specifies one or more locations that should remain visible on the map, though no markers or other indicators will be displayed.
	Visible []Location

	// Styles defines custom styles to alter the presentation of specific features on the map.
	//
	// See https://developers.google.com/maps/documentation/staticmaps/#StyledMaps
	Styles []Style
}

// Markers defines marker(s) to attach to the image.
type Markers struct {
	// Size defines the size of the marker(s).
	//
	// Accepted values are SizeTiny, SizeMid (the default) and SizeLarge.
	Size string // tiny, mid, small

	// Color defines the color of the marker(s).
	//
	// The alpha value of this color is ignored.
	Color color.Color

	// Label specifies a single uppercase alphanumeric character to place inside the marker image.
	Label string

	// IconURL specifies the URL of a custom icon image to use.
	IconURL string

	// HideShadow, if true, will not include a shadow for the marker(s).
	HideShadow bool

	// Locations specifies the locations of markers to be placed in this group.
	Locations []Location
}

func rgb(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("0x%02X%02X%02X", r>>8, g>>8, b>>8)
}

func rgba(c color.Color) string {
	r, g, b, a := c.RGBA()
	return fmt.Sprintf("0x%02X%02X%02X%02X", r>>8, g>>8, b>>8, a>>8)
}

func (m Markers) encode() string {
	s := []string{}
	if m.Size != "" {
		s = append(s, "size:"+m.Size)
	}
	if m.Color != nil {
		s = append(s, "color:"+rgb(m.Color))
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

// Path defines a path to attach to the image at specified locations.
type Path struct {
	// Weight specifies the thickness of the path in pixels. If no weight is specified, the default is 5 pixels.
	Weight int

	// Color defines the color of the path.
	Color color.Color

	// FillColor defines the color to fill the area of the path.
	FillColor color.Color

	// Geodesic, if true, indicates that the requested path should be interpreted as a geodesic line that follows the curvature of the earth.
	//
	// If false (the default), the path is rendered as a straight line in screen space.
	Geodesic bool

	// Polyline specifies an encoded polyline of points defining the path, if Locations is not provided.
	Polyline string

	// Locations specifies the points of the path.
	Locations []Location
}

func (p Path) encode() string {
	s := []string{}
	if p.Weight != 0 {
		s = append(s, fmt.Sprintf("weight:%d", p.Weight))
	}
	if p.Color != nil {
		s = append(s, "color:"+rgba(p.Color))
	}
	if p.FillColor != nil {
		s = append(s, "fillcolor:"+rgba(p.FillColor))
	}
	if p.Geodesic {
		s = append(s, "geodesic:true")
	}
	style := strings.Join(s, "|")
	if style != "" {
		style += "|"
	}
	if p.Polyline != "" {
		return style + "enc:" + p.Polyline
	}
	return style + encodeLocations(p.Locations)
}

// Style defines a set of rules to use to style the requested map image.
type Style struct {
	// Feature specifies the feature type for this style modification.
	//
	// See https://developers.google.com/maps/documentation/staticmaps/#StyledMapFeatures
	Feature string // TODO enum

	// Element indicates the subset of selected features to style.
	//
	// See https://developers.google.com/maps/documentation/staticmaps/#StyledMapElements
	Element string // TODO enum

	// Rules specifies the style rules to apply to the map.
	Rules []StyleRule
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

// StyleRule defines a style rule to apply to the map.
type StyleRule struct {
	// Hue is the hue value to apply to the selection.
	//
	// Note that while this takes a color value, it only uses this value to determine the basic color (its orientation around the color wheel),
	// not its saturation or lightness, which are indicated separately as percentage changes.
	Hue color.Color

	// Lightness (a value between -100 and 100) indicates the percentage change in brightness of the element. -100 is black, 100 is white.
	Lightness float64

	// Saturation (a value between -100 and 100) indicates the percentage change in intensity of the basic color to apply to the element.
	Saturation float64

	// Gamma (a value between 0.01 and 10.0, where 1.0 applies no correction) indicates the amount of gamma correction to apply to the element.
	Gamma *float64 // .01 to 10, default 1

	// InverseLightness, if true, inverts the Lightness value.
	InverseLightness bool

	// Visibility indicates whether and how the element appears on the map.
	//
	// Accepted values are VisibilityOn (the default), VisibilityOff and VisibilitySimplified.
	Visibility string
}

func (r StyleRule) encode() string {
	s := []string{}
	if r.Hue != nil {
		s = append(s, "hue:"+rgb(r.Hue))
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
	if r.InverseLightness {
		s = append(s, "inverse_lightness:true")
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
		p.Set("center", so.Center.Location())
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
