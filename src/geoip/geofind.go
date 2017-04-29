package main

import (
	geo "github.com/paulmach/go.geo"
	// geojson "github.com/paulmach/go.geojson"
	"errors"
	"math"
)

type Address struct {
	Geo        [2]float64             `json:"geo"`
	Properties map[string]interface{} `json:"name"`
}

func (geoc *Geocoder) Locate(lat float64, lng float64) (Address, error) {

	found := false
	var properties map[string]interface{}
	point := geo.NewPoint(lng, lat)
	// fmt.Println("locate",lat,lng)
	for _, feature := range geoc.features {

		if feature.contains(point) {

			if found {
				return Address{}, errors.New("Multiple results")
			}

			properties = feature.properties
			found = true

		}

	}

	if !found {
		return Address{}, errors.New("No result")
	}

	addr := Address{
		Geo:        [2]float64{lat, lng},
		Properties: properties,
	}
	return addr, nil
}

func (feature *GeoFeature) contains(point *geo.Point) bool {

	for _, polygon := range feature.polygons {

		if !polygon[0].contains(point) {
			continue
		}

		/* The first polygon is the enclosing box, the
		 * next ones are the holes */

		j := 1
		for ; j < len(polygon); j++ {
			if polygon[j].contains(point) {
				break
			}
		}

		if j == len(polygon) {
			return true
		}

	}
	return false
}

func (polygon *GeoPolygon) contains(point *geo.Point) bool {

	if !polygon.bounds.Contains(point) {
		return false
	}

	points := polygon.path.Points()
	start := len(points) - 1
	end := 0

	contains := intersectsWithRaycast(point, points[start], points[end])

	for i := 1; i < len(points); i++ {
		if intersectsWithRaycast(point, points[i-1], points[i]) {
			contains = !contains
		}
	}

	return contains
}

// From: https://github.com/kellydunn/golang-geo/blob/master/polygon.go
// Using the raycast algorithm, this returns whether or not the passed in point
// Intersects with the edge drawn by the passed in start and end points.
// Original implementation: http://rosettacode.org/wiki/Ray-casting_algorithm#Go
func intersectsWithRaycast(point *geo.Point, start geo.Point, end geo.Point) bool {
	// Always ensure that the the first point
	// has a y coordinate that is less than the second point
	if start.Lng() > end.Lng() {

		// Switch the points if otherwise.
		start, end = end, start

	}

	// Move the point's y coordinate
	// outside of the bounds of the testing region
	// so we can start drawing a ray
	for point.Lng() == start.Lng() || point.Lng() == end.Lng() {
		newLng := math.Nextafter(point.Lng(), math.Inf(1))
		point = geo.NewPoint(point.Lat(), newLng)
	}

	// If we are outside of the polygon, indicate so.
	if point.Lng() < start.Lng() || point.Lng() > end.Lng() {
		return false
	}

	if start.Lat() > end.Lat() {
		if point.Lat() > start.Lat() {
			return false
		}
		if point.Lat() < end.Lat() {
			return true
		}

	} else {
		if point.Lat() > end.Lat() {
			return false
		}
		if point.Lat() < start.Lat() {
			return true
		}
	}

	raySlope := (point.Lng() - start.Lng()) / (point.Lat() - start.Lat())
	diagSlope := (end.Lng() - start.Lng()) / (end.Lat() - start.Lat())

	return raySlope >= diagSlope
}
