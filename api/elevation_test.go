package api

import (
	"github.com/sakisds/go-srtm/coordinates"
	"testing"
)

type testLocation struct {
	locationString string
	expectedOutput []coordinates.Point
}

var testLocations = []testLocation{
	{locationString: "18.145,-42.174", expectedOutput: []coordinates.Point{
		{Latitude: 18.145, Longitude: -42.174},
	}},
	{locationString: "18.145,-42.174|19.124,12.111", expectedOutput: []coordinates.Point{
		{Latitude: 18.145, Longitude: -42.174},
		{Latitude: 19.124, Longitude: 12.111},
	}},
	{locationString: "18.145,-42.174|19.124,12.111|8.21,-1.22", expectedOutput: []coordinates.Point{
		{Latitude: 18.145, Longitude: -42.174},
		{Latitude: 19.124, Longitude: 12.111},
		{Latitude: 8.21, Longitude: -1.22},
	}},
}

func TestParseLocationString(t *testing.T) {
	for _, testLocation := range testLocations {
		output := parseLocationsString(testLocation.locationString)
		for index := range output {
			if testLocation.expectedOutput[index] != output[index] {
				t.Errorf("Expected: %v, Result: %v", testLocation.expectedOutput[index], output[index])
			}
		}
	}
}
