/*
For more information about GPX, see:

    https://www.topografix.com/gpx.asp
*/
package model

import (
	"time"
)

// GPX is the top-level object of a GPX document.
type GPX struct {
	Meta      Metadata   `xml:"metadata"`
	Waypoints []Waypoint `xml:"wpt"`
	Routes    []Route    `xml:"rte"`
	Tracks    []Track    `xml:"trk"`
}

// Link represents a URL, its textual description and the MIME
// type of the content it returns.
type Link struct {
	URL  string `xml:"href,attr"`
	Text string `xml:"text"`
	Type string `xml:"type"`
}

// Metadata contains general information about the GPX file.
type Metadata struct {
	Link      Link      `xml:"link"`
	Bounds    Bounds    `xml:"bounds,omitempty"`
	CreatedAt time.Time `xml:"time"`
}

// Bounds encodes a bounding box around all waypoints, tracks
// and routes in the GPX file.
type Bounds struct {
	MinLat  float64 `xml:"minlat"`
	MinLong float64 `xml:"minlon"`
	MaxLat  float64 `xml:"maxlat"`
	MaxLong float64 `xml:"maxlon"`
}

// Waypoint represents a single point.
type Waypoint struct {
	Name        string  `xml:"name"`
	Description string  `xml:"desc"`
	Type        string  `xml:"type"`
	Link        Link    `xml:"link"`
	Latitude    float64 `xml:"lat,attr"`
	Longitude   float64 `xml:"lon,attr"`
	Elevation   float64 `xml:"ele,omitempty"`
	// Other fields, forthcoming.
}

// Route represents a pre-determined route, plotted by a user.
type Route struct {
	Name        string     `xml:"name"`
	Description string     `xml:"desc"`
	Type        string     `xml:"type"`
	Link        Link       `xml:"link"`
	Points      []Waypoint `xml:"rtept"`
}

// Track represents an actual route recorded by a GPS unit.
type Track struct {
	Name        string     `xml:"name"`
	Description string     `xml:"desc"`
	Type        string     `xml:"type"`
	Link        Link       `xml:"link"`
	Segments    []Waypoint `xml:"trkseg>trkpt"`
}

// ToGeoJSON converts the GPX object into a GeoJSON one.
//
// Routes  and Track Segments are converted into GeoJSON
// MultiLineString features and Waypoints as GeoJSON
// MultiPoint features.
//
// All generated geometries are returned as individual
// Features within a FeatureCollection.
func (g *GPX) ToGeoJSON() (*FeatureCollection, error) {
	features := []Feature{}

	// Waypoints
	if len(g.Waypoints) > 0 {
		arr := make([]Position, len(g.Waypoints))
		for i, v := range g.Waypoints {
			pos, err := NewPosition([]float64{v.Longitude, v.Latitude, v.Elevation})
			if err != nil {
				return nil, err
			}
			arr[i] = pos
		}
		pt := NewMultiPoint(arr)
		f := NewFeature(pt)
		features = append(features, f)
	}

	// Tracks
	if len(g.Tracks) > 0 {
		arr := make([][]Position, len(g.Tracks))
		for i, t := range g.Tracks {
			positions, err := convertToPositions(t.Segments)
			if err != nil {
				return nil, err
			}
			arr[i] = positions
		}
		ml := NewMultiLineString(arr)
		f := NewFeature(ml)
		features = append(features, f)
	}

	// Routes
	if len(g.Routes) > 0 {
		arr := make([][]Position, len(g.Routes))
		for i, r := range g.Routes {
			positions, err := convertToPositions(r.Points)
			if err != nil {
				return nil, err
			}
			arr[i] = positions
		}
		ml := NewMultiLineString(arr)
		f := NewFeature(ml)
		features = append(features, f)
	}

	coll := NewFeatureCollection(features)

	return &coll, nil
}

// convertToPositions turns a slice of GPX Wayppints into a slice
// of GeoJSON Positions.
func convertToPositions(wpts []Waypoint) ([]Position, error) {
	positions := make([]Position, len(wpts))
	for i, wp := range wpts {
		pos, err := NewPosition([]float64{wp.Longitude, wp.Latitude, wp.Elevation})
		if err != nil {
			return nil, err
		}
		positions[i] = pos
	}

	return positions, nil
}
