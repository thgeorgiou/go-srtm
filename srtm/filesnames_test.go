package SRTM

import (
	"testing"
)

type testFilename struct {
	filename    string
	expectedLat float64
	expectedLon float64
}

var testFilenames = []testFilename{
	{filename: "N18E033", expectedLat: 18.0, expectedLon: 33.0},
	{filename: "S16E037", expectedLat: -16.0, expectedLon: 37.0},
	{filename: "N22W044", expectedLat: 22.0, expectedLon: -44.0},
	{filename: "S18W033", expectedLat: -18.0, expectedLon: -33.0},
	{filename: "N02W110", expectedLat: 2.0, expectedLon: -110.0},
}

// TestFilenameToCoordinates to return the correct lat/lon pair every time
func TestFilenameToCoordinates(t *testing.T) {
	for _, testData := range testFilenames {
		lat, lon, err := FilenameToCoordinates(testData.filename)
		if err != nil {
			t.Errorf("Expected no error with '%s', got '%s'", testData.filename, err)
		}
		if lat != testData.expectedLat {
			t.Errorf("Expected lat = '%.1f' with '%s', got '%.1f", testData.expectedLat, testData.filename, lat)
		}
		if lon != testData.expectedLon {
			t.Errorf("Expected lon = '%.1f' with '%s', got '%.1f", testData.expectedLon, testData.filename, lon)
		}
	}
}
