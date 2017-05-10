package SRTM

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sakisds/gigahertzor/coordinates"
)

// SquareSize The size of each row and column in the SRTM files
const SquareSize = 3601

// File Represents a SRTM file in memory
type File struct {
	Name                string
	Contents            []byte
	latitude, longitude float64
}

// LoadFile Load a SRTM (.hgt) file into the memory
func LoadFile(path string) (file *File, err error) {
	// Parse path
	_, filename := filepath.Split(path)
	latitude, longitude, err := FilenameToCoordinates(filename[0:7])
	if err != nil {
		return
	}

	// Create File struct
	file = &File{}
	file.Name = filename[:7]
	file.latitude = latitude
	file.longitude = longitude

	// Check if file exists/is not a directory
	stats, err := os.Stat(path)
	if err != nil {
		return
	}
	if stats.IsDir() {
		err = errors.New("Path is a directory instead of a file: " + path)
		return
	}

	// Load content
	file.Contents, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}

	return
}

// IsCovered checks if the given coordinates are covered by this file
func (file File) IsCovered(latitude, longitude float64) bool {
	var latCovered, lonCovered bool

	if file.latitude > 0 {
		latCovered = file.latitude <= latitude && file.latitude+1 > latitude
	} else {
		latCovered = file.latitude >= latitude && file.latitude+1 < latitude
	}

	if file.longitude > 0 {
		lonCovered = file.longitude <= longitude && file.longitude+1 > longitude
	} else {
		lonCovered = file.longitude >= longitude && file.longitude+1 < longitude
	}

	return latCovered && lonCovered
}

// GetElevation returns elevation at a specific point
func (file File) GetElevation(latitude, longitude float64) (int, error) {
	// Check if coordinates are out of bounds
	if !file.IsCovered(latitude, longitude) {
		return 0, fmt.Errorf("Coordinates (%f, %f) are out of bounds for file %s", latitude, longitude, file.Name)
	}

	// Determinate row and column of file
	row := int((file.latitude + 1.0 - latitude) * (float64(SquareSize - 1.0)))
	column := int((longitude - file.longitude) * (float64(SquareSize - 1.0)))

	// Get the two bytes and return the elevation value
	index := row*SquareSize + column
	return int(file.Contents[index*2])*256 + int(file.Contents[index*2+1]), nil
}

// GetElevationPath returns an array of elevation information for the path that connects the two pairs of coordinates
func (file File) GetElevationPath(from, to coordinates.Point) (path []int, err error) {
	path = make([]int, 10)
	for i := 0; i < 10; i++ {
		currentPoint := coordinates.Lerp(from, to, float64(i)/10.0)
		elevation, err := file.GetElevation(currentPoint.Latitude, currentPoint.Longitude)
		if err != nil {
			return path, err
		}

		path[i] = elevation
	}

	return
}
