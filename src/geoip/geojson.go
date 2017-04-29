package main

import (
	"errors"
	"fmt"
	geo "github.com/paulmach/go.geo"
	geojson "github.com/paulmach/go.geojson"
	"io/ioutil"
	"time"
)

type GeoPolygon struct {
	path   *geo.Path
	bounds *geo.Bound
}

type GeoFeature struct {
	polygons   [][]GeoPolygon
	properties map[string]interface{}
	parent     string
}

type Geocoder struct {
	geojson *geojson.FeatureCollection

	features []GeoFeature
}

func NewGeocoder(level int) (*Geocoder, error) {

	geoc := new(Geocoder)
	if err := geoc.decode(level); err != nil {
		return nil, err
	}
	return geoc, nil

}

func (geoc *Geocoder) Dump() {

	for _, feature := range geoc.geojson.Features {
		bbox := feature.BoundingBox
		fmt.Println("\t", feature.Properties["name:zh"], bbox, "\n")
		// geom := feature.Geometry
		// if (!geom.IsMultiPolygon()) {
		// 	fmt.Println("Unknown geom type",geom.Type)
		// 	continue
		// }
		// if false {
		// 	for _, polygon := range(geom.MultiPolygon) {
		// 		j := 0
		// 		for j<len(polygon) {
		// 	 		path := geo.NewPathFromXYSlice(polygon[j])
		// 	 		bounds := path.Bound()
		// 	 		fmt.Print("\t\t#/",j,":",bounds,"\n")
		// 	 		j = j+1
		// 		}

		// 	}
		// }

	}
}

func (geoc *Geocoder) List() []string {

	l := make([]string, 0)
	for _, feature := range geoc.features {

		cn := feature.properties["name:zh"]
		if cn == nil {
			cn = feature.properties["name"]
		}
		name := cn.(string)
		if name == "Border Henan - Hubei" {
			continue
		}
		l = append(l, name)

	}
	return l

}

func (geoc *Geocoder) decode(level int) error {

	stime := time.Now()
	filename := fmt.Sprintf("asset/china-geojson/admin_level_%d.geojson", level)
	rawjson, e1 := ioutil.ReadFile(filename)
	if e1 != nil {
		return errors.New("Can not find the geojson file (" + filename + ")")
	}
	g, e2 := geojson.UnmarshalFeatureCollection([]byte(rawjson))
	if e2 != nil {
		return errors.New("Unable to decode the geojson file")
	}

	features := make([]GeoFeature, len(g.Features))
	geoc.geojson = g
	geoc.features = features

	for i, feature := range g.Features {
		geom := feature.Geometry
		if !geom.IsMultiPolygon() {
			fmt.Println("Unknown geom type", geom.Type)
			continue
		}

		polygons := make([][]GeoPolygon, len(geom.MultiPolygon))
		features[i].properties = feature.Properties
		features[i].polygons = polygons

		for j, polygon := range geom.MultiPolygon {

			polygons[j] = make([]GeoPolygon, len(polygon))
			for k := 0; k < len(polygon); k++ {
				path := geo.NewPathFromXYSlice(polygon[k])
				bounds := path.Bound()
				polygons[j][k] = GeoPolygon{
					bounds: bounds,
					path:   path,
				}
			}

		}

	}

	dt := time.Now().Sub(stime)
	fmt.Println("[geos] decoding ", (len(rawjson) / (1024 * 1024)), "MB in ", dt, " # ", len(g.Features), "features")
	return nil

}
