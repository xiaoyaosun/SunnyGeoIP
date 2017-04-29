package main

import (
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
)

type GeoIPArr struct {
	Country string  `json:"country"`
	City    string  `json:"city"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	CityEn  string  `json:"cityen"`
}

// 上海IP：{"Ip":{"add":"114.80.166.240"}}
// 北京IP：{"Ip":{"add":"124.207.28.124"}}
// 郑州IP：{"Ip":{"add":"61.53.96.255"}}
// 南京IP：{"Ip":{"add":"221.229.251.45"}}
func (server *HttpServer) FindGeoInfoByIP(ip net.IP) (recordjson GeoIPArr, err error) {

	db, err := geoip2.Open("asset/geo-ip/GeoLite2-City.mmdb")
	if err != nil {
		return recordjson, err
	}
	defer db.Close()
	record, err := db.City(ip)
	// fmt.Printf("record: %+v\n", record)
	if err != nil {
		return recordjson, err
	}

	if len(record.City.Names["en"]) == 0 {
		return recordjson, fmt.Errorf("Invalid city")
	}

	recordjson.City = record.City.Names["zh-CN"]
	recordjson.Country = record.Country.Names["zh-CN"]
	recordjson.CityEn = record.City.Names["en"]
	recordjson.Lat = record.Location.Latitude
	recordjson.Lng = record.Location.Longitude
	//fmt.Println("recordjson: ", recordjson)
	return recordjson, nil
}
