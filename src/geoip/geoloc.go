package main

import (
	"fmt"
)

type GeolocationServer struct {
	provinces   *Geocoder
	prefectures *Geocoder
}

func NewGeolocationServer() *GeolocationServer {

	server := new(GeolocationServer)

	provinces, err1 := NewGeocoder(4)
	if err1 != nil {
		fmt.Println("Ooops..", err1)
		return nil
	}

	prefectures, err2 := NewGeocoder(5)
	if err2 != nil {
		fmt.Println("Ooops..", err2)
		return nil
	}

	/* For each prefecture, find which province it belongs to */
	for i, prefecture := range prefectures.features {
		center := prefecture.polygons[0][0].bounds.Center()
		found := false
		for _, province := range provinces.features {
			if province.contains(center) {
				if found {
					fmt.Println("Duplicated province for ", prefecture.properties["name"].(string))
				}

				name := getName(province.properties)
				prefectures.features[i].parent = name
				found = true
			}
		}
		if !found {
			// fmt.Println("No province for ", prefecture.properties["name"].(string))
		}
	}

	server.provinces = provinces
	server.prefectures = prefectures
	return server

}

type Geolocation struct {
	// Geo [2]float64 `json:"geo"`
	Province   string `json:"province,omitempty"`
	Prefecture string `json:"prefecture,omitempty"`

	Gps struct {
		Lat  float64 `json:"lat"`
		Lng  float64 `json:"lng"`
		Time int64   `json:"time,omitempty"`
	} `json:"gps"`
}

func (server *GeolocationServer) Locate(lat float64, lng float64) (Geolocation, error) {

	var geoloc Geolocation
	// geoloc.Geo = [2]float64{ lat,lng }
	if province, err1 := server.provinces.Locate(lat, lng); err1 == nil {
		geoloc.Province = getName(province.Properties)
	} else {
		return geoloc, err1
	}

	if prefecture, err1 := server.prefectures.Locate(lat, lng); err1 == nil {
		geoloc.Prefecture = getName(prefecture.Properties)
	}

	return geoloc, nil

}

func (server *GeolocationServer) Provinces() map[string][]string {

	l := make(map[string][]string, 0)
	for _, feature := range server.provinces.features {
		name := getName(feature.properties)
		if name == "Border Henan - Hubei" {
			continue
		}
		if _, has := l[name]; !has {
			l[name] = make([]string, 0)
		}
	}

	for _, feature := range server.prefectures.features {
		province := feature.parent
		if len(province) == 0 {
			continue
		}
		if _, has := l[province]; !has {
			l[province] = make([]string, 0)
		}
		if feature.parent == province {

			name := getName(feature.properties)
			if len(name) == 0 {
				fmt.Println(feature.properties)
				continue
			}
			l[province] = append(l[province], name)
		}

	}
	return l

}

func (server *GeolocationServer) Prefectures(province string) []string {

	l := make([]string, 0)
	for _, feature := range server.prefectures.features {

		if feature.parent == province {

			name := getName(feature.properties)
			if len(name) == 0 {
				fmt.Println(feature.properties)
				continue
			}
			l = append(l, name)
		}

	}
	return l

}

func getName(properties map[string]interface{}) string {
	cn := properties["name:zh"]
	if cn == nil {
		cn = properties["name"]
	}
	name := cn.(string)
	return name

}
