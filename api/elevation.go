package api

import (
	"encoding/json"
	"errors"
	"github.com/go-ini/ini"
	"github.com/sakisds/go-srtm/coordinates"
	"github.com/sakisds/go-srtm/srtm"
	"net/http"
	"strconv"
	"strings"
)

var srtmManager SRTM.Manager

func InitElevation(cfg ini.File) (err error) {
	if !cfg.Section("").HasKey("dataset") {
		return errors.New("No dataset path set in config")
	}

	srtmManager, err = SRTM.CreateManager(cfg.Section("").Key("dataset").String())
	if err != nil {
		return err
	}
	return nil
}

func parseLocationsString(locationsString string) (locations []coordinates.Point) {
	// Split string at each separator
	stringFragments := strings.Split(locationsString, "|")
	locations = make([]coordinates.Point, 0, len(stringFragments))

	// Parse each location
	for _, stringFragment := range stringFragments {
		// Split each coordinate into the two numbers and parse them
		coordinatesStrings := strings.Split(stringFragment, ",")

		latitude, err := strconv.ParseFloat(coordinatesStrings[0], 64)
		if err != nil {
			continue
		}

		longitude, err := strconv.ParseFloat(coordinatesStrings[1], 64)
		if err != nil {
			continue
		}

		// Add to list of locations
		locations = append(locations, coordinates.Point{Latitude: latitude, Longitude: longitude})
	}

	return
}

func HandleElevation(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	// Check if the user passed a location string
	locations := r.URL.Query()["locations"]
	if len(locations) == 0 {
		apiError := Error{Name: "parameter_missing", Description: "Parameter 'locations' not optional", Data: nil}
		encoder.Encode(apiError)
		return
	}

	// Parse locations
	points := parseLocationsString(locations[0])
	if len(points) == 0 {
		apiError := Error{Name: "parameter_invalid", Description: "No valid points in parameter 'locations'", Data: locations}
		encoder.Encode(apiError)
		return
	}

	// Get elevation for each point and send to client as a JSON array
	elevationArray := make([]int, len(points))
	defer srtmManager.Clear()
	for i, point := range points {
		elevation, err := srtmManager.GetElevation(point)
		if err != nil {
			elevation = -1
		}

		elevationArray[i] = elevation
	}

	encoder.Encode(elevationArray)
}

func HandleElevationPath(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	// Check if the user passed the two parameters
	urlQuery := r.URL.Query()

	fromStr := urlQuery["from"]
	if len(fromStr) == 0 {
		apiError := Error{Name: "parameter_missing", Description: "Parameter 'from' not optional", Data: nil}
		encoder.Encode(apiError)
		return
	}

	toStr := urlQuery["to"]
	if len(toStr) == 0 {
		apiError := Error{Name: "parameter_missing", Description: "Parameter 'to' not optional", Data: nil}
		encoder.Encode(apiError)
		return
	}

	// Parse the two locations
	locations := parseLocationsString(fromStr[0] + "|" + toStr[0])
	if len(locations) != 2 {
		apiError := Error{Name: "invalid_location", Description: "Could not parse locations", Data: nil}
		encoder.Encode(apiError)
		return
	}

	// Get the elevation data for a path connecting the two points
	count := int(locations[0].DistanceTo(locations[1]) / 30) // One elevation point every 30m
	step := 1 / float64(count)
	elevationData := make([]int, 0, count)

	for i := 0; i < count; i++ {
		currentPoint := coordinates.Lerp(locations[0], locations[1], float64(i)*step)
		elevation, err := srtmManager.GetElevation(currentPoint)
		if err == nil {
			elevationData = append(elevationData, elevation)
		} else {
			elevationData = append(elevationData, -1)
		}
	}

	// Return to client as a JSON array
	encoder.Encode(elevationData)
}
