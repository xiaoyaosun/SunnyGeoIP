[![Build Status](https://travis-ci.org/xiaoyaosun/SunnyGeoIP.svg?branch=master)](https://travis-ci.org/xiaoyaosun/SunnyGeoIP)

# SunnyGeoIP

## Description

	You can get the GPS/Province/City info from this API by GeoLite2 (GeoLite2-City.mmdb is a free one)

## Code

    github: https://github.com/xiaoyaosun/SunnyGeoIP.git

## Golang Version

    1.8

## Vendor directory management

    https://glide.sh/    

## Usage

To build and run the server:

    make run

Then 


**Case (1): Your client IP is unknown**

	curl --data "user=ABC" http://localhost:8087/geoip/location

Will give:

	{
	    "message": "Can not get the location information",
	    "status": "error"
	}
	
**Case (2): Your client IP is BeiJing**

	curl --data "user=ABC" http://localhost:8087/geoip/location

Will give:

	{
	    "data": {
	        "province": "北京市",
	        "gps": {
	            "lat": 39.9289,
	            "lng": 116.3883
	        }
	    },
	    "refresh": 60,
	    "status": "ok"
	}

**Case (3): You set ipaddr is ShangHai**

	curl --data "user=ABC&ipaddr=114.80.166.240" http://localhost:8087/geoip/location

Will give:

	{
	    "data": {
	        "province": "上海市",
	        "gps": {
	            "lat": 31.0456,
	            "lng": 121.3997
	        }
	    },
	    "refresh": 60,
	    "status": "ok"
	}

**Case (4): You set the GPS location info**

	curl --data "lat=39.379436&lng=116.091230&user=ABC" http://localhost:8087/geoip/location

Will give:

	{
	    "data": {
	        "province": "河北省",
	        "prefecture": "保定市",
	        "gps": {
	            "lat": 39.379436,
	            "lng": 116.09123
	        }
	    },
	    "refresh": 60,
	    "status": "ok"
	}

## Test

To run the test, just enter

    make test

You will see 

	[geos] decoding  7 MB in  1.042079445s  #  36 features
	[geos] decoding  24 MB in  3.526207237s  #  335 features
	[geos] decoding  7 MB in  984.411265ms  #  36 features
	Found in  376.918µs
	[geos] decoding  7 MB in  984.019947ms  #  36 features
	[geos] decoding  24 MB in  3.37118092s  #  335 features
	PASS

## Production server

In order to push the production server, you just need to push to the `prod` branch.
    
    git push origin prod

## Data source

The geojson files are downloaded from 

    https://mapzen.com/data/borders/ (https://s3.amazonaws.com/osm-polygons.mapzen.com/china_geojson.tgz)

Official prefecture/province list - 最新县及县以上行政区划代码

    http://www.stats.gov.cn/tjsj/tjbz/xzqhdm/201401/t20140116_501070.html

## Contact Me

Email

    xiaoyaosun AT qq DOT com

