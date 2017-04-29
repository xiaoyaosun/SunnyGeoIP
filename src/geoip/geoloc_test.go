package main

import (
	"testing"
)

func TestLocateAPI(t *testing.T) {

	srv := NewGeolocationServer()
	if srv == nil {
		t.Error("Unable to open the location server")
		return
	}

	location, err1 := srv.Locate(50.58, 123.7)

	if err1 != nil {
		t.Fatal("Expected to find the position (error: ", err1, ")")
	}

	if location.Province != "内蒙古自治区" {
		t.Error("The location should be 内蒙古自治区. It is ", location.Province)
	}

	if location.Prefecture != "呼伦贝尔市" {
		t.Error("The location should be 呼伦贝尔市. It is ", location.Prefecture)
	}

}
