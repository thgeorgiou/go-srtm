package SRTM

import (
	"fmt"
	"strconv"
	"github.com/sakisds/gigahertzor/coordinates"
)

// FilenameLengthError Error thrown when a .hgt file has an invalid name due to length
type FilenameLengthError struct {
	Filename string
	length   int
}

func (e *FilenameLengthError) Error() string {
	return fmt.Sprintf("Filename %s has invalid length! (%d instead of 7)", e.Filename, e.length)
}

// FilenameToCoordinates Returns the coordinates from the starting point of a .hgt filename
func FilenameToCoordinates(filename string) (latitude float64, longitude float64, err error) {
	// Check length
	if len(filename) != 7 {
		return 0.0, 0.0, &FilenameLengthError{Filename: filename, length: len(filename)}
	}

	// Grab latitude
	latitude, err = strconv.ParseFloat(filename[1:3], 64)
	if err != nil {
		return 0.0, 0.0, err
	}
	if filename[0] == 'S' { // Make negative if south
		latitude *= -1
	}

	// Grab longitude
	longitude, err = strconv.ParseFloat(filename[4:7], 64)
	if err != nil {
		return 0.0, 0.0, err
	}
	if filename[3] == 'W' { // Make negative if west
		longitude *= -1
	}

	return
}

// CoordinatesToFilename returns the filename that covers the given point
func CoordinatesToFilename(point coordinates.Point) string {
	var NorthOrSouth, EastOrWest string

	if point.Latitude >= 0 {
		NorthOrSouth = "N"
	} else {
		NorthOrSouth = "S"
	}

	if point.Longitude >= 0 {
		EastOrWest = "E"
	} else {
		EastOrWest = "W"
	}

	return fmt.Sprintf("%s%d%s%03d", NorthOrSouth, int(point.Latitude), EastOrWest, int(point.Longitude))
}
