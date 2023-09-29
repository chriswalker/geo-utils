/*
For more information about GeoJSON, see:

	https://datatracker.ietf.org/doc/html/rfc7946
*/
package model

import (
	"errors"
)

// Point represents a single position on a map.
type Point struct {
	Type        string   `json:"type"`
	Coordinates Position `json:"coordinates"`
}

// NewPoint creates and returns a new Point geometry.
func NewPoint(p Position) *Point {
	return &Point{
		Type:        "Point",
		Coordinates: p,
	}
}

// MultiPoint represents two or more points.
type MultiPoint struct {
	Type        string     `json:"type"`
	Coordinates []Position `json:"coordinates"`
}

// NewLMultiPoint creates and returns a new MultiPoint geometry.
func NewMultiPoint(p []Position) *MultiPoint {
	return &MultiPoint{
		Type:        "MultiPoint",
		Coordinates: p,
	}
}

// LineString is an ordered collection of Positions; visually
// each position is connected by a straight line.
type LineString struct {
	Type        string     `json:"type"`
	Coordinates []Position `json:"coordinates"`
}

// NewLineString creates and returns a new LineString geometry.
func NewLineString(p []Position) *LineString {
	return &LineString{
		Type:        "LineString",
		Coordinates: p,
	}
}

// MultiLineString is a collection of LineStrings.
type MultiLineString struct {
	Type        string       `json:"type"`
	Coordinates [][]Position `json:"coordinates"`
}

// NewMultiLineString creates and returns a new MultiLineString geometry.
func NewMultiLineString(p [][]Position) *MultiLineString {
	return &MultiLineString{
		Type:        "MultiLineString",
		Coordinates: p,
	}
}

// A Position is a slice of float64 - it must contain two elements
// minimum, and three at most (where the third indicates elevation,
// if present).
type Position []float64

// NewPosition converts an existing []float64 into a Position
// provided it is valid (contains at least two and no more than
// three elements).
func NewPosition(vals []float64) (Position, error) {
	if len(vals) < 2 || len(vals) > 3 {
		return nil, errors.New("invalid number of positional elements")
	}

	p := make(Position, len(vals))
	copy(p, vals)

	return p, nil
}

// Feature represents a single GeoJSON feature.
type Feature struct {
	Type        string         `json:"type"`
	ID          string         `json:"id,omitempty"`
	BoundingBox []float64      `json:"bbox,omitempty"`
	Geometry    any            `json:"geometry"` // Workable, but 'any' sucks a bit
	Properties  map[string]any `json:"properties"`
}

// NewFeature creates a new GeoJSON Feature.
func NewFeature(geometry any) Feature {
	feature := Feature{
		Type:     "Feature",
		Geometry: geometry,
	}
	return feature
}

// FeatureCollection is a collection of Features.
type FeatureCollection struct {
	Type        string    `json:"type"`
	BoundingBox []float64 `json:"bbox,omitempty"`
	Features    []Feature `json:"features"`
}

// NewFeatureCollection creates a FeatureCollection using the
// supplied slice of Features.
func NewFeatureCollection(f []Feature) FeatureCollection {
	return FeatureCollection{
		Type:     "FeatureCollection",
		Features: f,
	}
}
