package main

import (
	"fmt"
	"testing"
	"time"
	// "runtime"
)

func TestLocate(t *testing.T) {

	geoc, err := NewGeocoder(4)
	if err != nil {
		t.Error("Expected to decode the geojson file. (error: ", err, ")")
		return
	}

	startTime := time.Now()
	location, err1 := geoc.Locate(39.913818, 116.363625)
	dt := time.Now().Sub(startTime)

	if err1 != nil {
		t.Fatal("Expected to find the position (error: ", err1, ")")
	}

	if location.Properties["name"] != "北京市" {
		t.Error("The location should be 北京市. It is ", location.Properties["name"])
	}

	fmt.Println("Found in ", dt)

}
