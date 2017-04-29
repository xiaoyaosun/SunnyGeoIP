package main

import (
	_ "fmt"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"github.com/qiniu/mockhttp.v2"
	"github.com/qiniu/rpc"
)

func NewMockHttpServer() *HttpServer {

	server := new(HttpServer)

	server.startTime = time.Now()

	server.waiting = true
	server.ready = sync.NewCond(&sync.Mutex{})

	server.geoloc = nil

	return server

}

// For the mockhttp using
func (server *HttpServer) testListenGeoLocation() {

	router := mux.NewRouter()
	server.router = router
	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.HandleFunc("/geoip/location", server.geoipHandler)

	// Listening on the Mock Server
	mockhttp.ListenAndServe(testDomain, router)
}

type TestGeoIP2Result struct {
	Status  string `json:"status"`
	Refresh int    `json:"refresh"`
	Data    struct {
		Province   string `json:"province"`
		Prefecture string `json:"prefecture"`
		Gps        struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"gps"`
	} `json:"data"`
}

const testUserID = "MockTestUserID"

// 山西大同的坐标
const testLat = "39.379436"
const testLng = "114.091230"

// 上海的IP
const testSHIpaddr = "114.80.166.240"

// 北京的IP
const testBJIpaddr = "124.42.72.120"

// 回环IP
const testLocalipaddr = "127.0.0.1"

// 定义当前服务的IP
const testDomain = "test.com"
const testServer = "http://test.com/geoip/location"

func init() {

	server := NewMockHttpServer()
	server.testListenGeoLocation()
	server.bind(NewGeolocationServer())
}

func TestGeoIP2_basic(t *testing.T) {

	t.Log("TestGeoIP2_basic")
	recordjson, err1 := server.FindGeoInfoByIP(net.ParseIP(testBJIpaddr))

	if err1 != nil {
		t.Fatal("Expected to find the position by IP (error: ", err1, ")")
	}

	if recordjson.City != "北京" {
		t.Error("The IP location City should be 北京. It is ", recordjson.City)
	}

	if recordjson.CityEn != "Beijing" {
		t.Error("The IP location english name should be Beijing. It is ", recordjson.CityEn)
	}

	if recordjson.Country != "中国" {
		t.Error("The IP location Country should be 中国. It is ", recordjson.Country)
	}

	if recordjson.Lat != 39.9289 {
		t.Error("The IP location Lat should be 39.9289. It is ", recordjson.Lat)
	}

	if recordjson.Lng != 116.3883 {
		t.Error("The IP location Lng should be 116.3883. It is ", recordjson.Lng)
	}
}

func TestMock_basic(t *testing.T) {

	t.Log("TestMock_basic")
	//var l rpc.Logger
	c := rpc.Client{mockhttp.DefaultClient}
	{
		var ret TestGeoIP2Result
		req := make(map[string][]string)
		req["user"] = []string{testUserID}
		req["ipaddr"] = []string{testSHIpaddr}
		err := c.CallWithForm(nil, &ret, testServer, req)
		if err != nil {
			t.Fatal("call ret failed:", err)
		}

		t.Log(ret)
		//t.Log(l)
		if ret.Data.Province != "上海市" {
			t.Error("The Province should be 上海市. It is ", ret.Data.Province)
		}

		if ret.Data.Gps.Lat != 31.0456 {
			t.Error("The IP location Lat should be 31.0456. It is ", ret.Data.Gps.Lat)
		}

		if ret.Data.Gps.Lng != 121.3997 {
			t.Error("The IP location Lng should be 121.3997. It is ", ret.Data.Gps.Lng)
		}

	}
}

// Case1: With lat/lng and User
func TestMock_case1(t *testing.T) {

	t.Log("Case1: With lat/lng and User")
	c := rpc.Client{mockhttp.DefaultClient}
	{
		var ret TestGeoIP2Result
		req := make(map[string][]string)
		req["user"] = []string{testUserID}
		req["lat"] = []string{testLat}
		req["lng"] = []string{testLng}
		err := c.CallWithForm(nil, &ret, testServer, req)
		if err != nil {
			t.Fatal("call ret failed:", err)
		}

		t.Log(ret)
		if ret.Data.Province != "山西省" {
			t.Error("The Province should be 山西省. It is ", ret.Data.Province)
		}

		if ret.Data.Prefecture != "大同市" {
			t.Error("The Province should be 大同市. It is ", ret.Data.Province)
		}
	}
}

// Case2: With ipaddr and User
func TestMock_case2(t *testing.T) {

	t.Log("Case2: With ipaddr and User")
	c := rpc.Client{mockhttp.DefaultClient}
	{
		var ret TestGeoIP2Result
		req := make(map[string][]string)
		req["user"] = []string{testUserID}
		req["ipaddr"] = []string{testSHIpaddr}

		err := c.CallWithForm(nil, &ret, testServer, req)
		if err != nil {
			t.Fatal("call ret failed:", err)
		}

		t.Log(ret)
		if ret.Data.Province != "上海市" {
			t.Error("The Province should be 上海市. It is ", ret.Data.Province)
		}

		if ret.Data.Prefecture != "" {
			t.Error("The Province should be . It is ", ret.Data.Province)
		}
	}
}

// Case3: No param
func TestMock_case3(t *testing.T) {

	t.Log("Case3: No param")
	c := rpc.Client{mockhttp.DefaultClient}
	{
		var ret TestGeoIP2Result
		req := make(map[string][]string)

		err := c.CallWithForm(nil, &ret, testServer, req)
		if err != nil {
			t.Fatal("call ret failed:", err)
		}

		t.Log(ret)
		if ret.Status != "error" {
			t.Error("The Province should be error. It is ", ret.Status)
		}
	}
}

// Case4: ipaddr=127.0.0.1 with User
func TestMock_case4(t *testing.T) {

	t.Log("Case4: ipaddr=127.0.0.1 with User")
	c := rpc.Client{mockhttp.DefaultClient}
	{
		var ret TestGeoIP2Result
		req := make(map[string][]string)
		req["user"] = []string{testUserID}
		req["ipaddr"] = []string{testLocalipaddr}

		err := c.CallWithForm(nil, &ret, testServer, req)
		if err != nil {
			t.Fatal("call ret failed:", err)
		}

		t.Log(ret)
		if ret.Status != "error" {
			t.Error("The Province should be error. It is ", ret.Status)
		}
	}
}
