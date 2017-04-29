package main

import (
	"encoding/json"
	_ "fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (server *HttpServer) geoipHandler(w http.ResponseWriter, r *http.Request) {

	server.ready.L.Lock()
	if server.waiting {
		server.ready.Wait()
	}
	server.ready.L.Unlock()

	startTime := time.Now()
	response := map[string]interface{}{}

	slat := r.FormValue("lat")
	slng := r.FormValue("lng")
	user := r.FormValue("user")

	lat, e1 := strconv.ParseFloat(slat, 64)
	lng, e2 := strconv.ParseFloat(slng, 64)

	hasLocation := false

	isUser := len(user) != 0

	if len(user) == 0 {

		response["status"] = "error"
		response["message"] = "Device or User ID must be specified"

	} else if e1 != nil || e2 != nil {

		if !isUser {

			response["status"] = "error"
			response["message"] = "Can not decode the lat/lng"
		}

	} else {

		hasLocation = true
	}

	/* Only allow to get the location from IP for the users */
	if !hasLocation && isUser {

		// Add the ipaddr parameter
		sipaddr := ""
		var ip net.IP
		sipaddr = r.FormValue("ipaddr")
		if sipaddr == "" {

			// http://nginx.org/en/docs/http/ngx_http_proxy_module.html
			// If X-Forwarded-For more than 1, the first one is the real IP
			// Use the "X-Forwarded-For" to insead the "X-Real-Ip"
			// For example
			// X-Real-Ip:[139.217.17.36]
			// X-Forwarded-For:[123.117.169.255, 139.217.17.36]
			forwardIP := r.Header.Get("X-Forwarded-For")
			splitForward := strings.Split(forwardIP, ",")
			if len(splitForward) >= 1 {
				sipaddr = splitForward[0]
			}
		}

		if len(sipaddr) == 0 {
			sipaddr, _, _ = net.SplitHostPort(r.RemoteAddr)
			ip = net.ParseIP(sipaddr)
		} else {
			ip = net.ParseIP(sipaddr)
		}

		// No GPS info and User
		// Get the lat/lng info from GeoIP2
		record, err := server.FindGeoInfoByIP(ip)

		if err != nil {
			response["status"] = "error"
			response["message"] = "Can not decode the lat/lng"
		} else {
			lat = record.Lat
			lng = record.Lng
			hasLocation = true
		}
	}

	if hasLocation {

		region, err := server.geoloc.Locate(lat, lng)
		if err != nil {

			response["status"] = "ok"

			var region Geolocation
			response["data"] = region
			region.Gps.Lat = lat
			region.Gps.Lng = lng
			response["data"] = region
			response["refresh"] = 3 * 60

		} else {

			region.Gps.Lat = lat
			region.Gps.Lng = lng
			response["status"] = "ok"
			response["refresh"] = 60
			response["data"] = region
		}
	} else {

		response["status"] = "error"
		response["message"] = "Can not get the location information"
	}

	body, _ := json.Marshal(response)
	writeResponse(w, body, startTime)

}
