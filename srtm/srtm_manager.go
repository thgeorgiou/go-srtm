package SRTM

import (
	"github.com/sakisds/gigahertzor/coordinates"
	"io/ioutil"
	"fmt"
	"path/filepath"
)

type Manager struct {
	files     map[string]File
	coverage  map[string]bool
	dataStore string
}

// DataUnavailableError Error thrown when data is not available for a requested point
type DataUnavailableError struct {
	Point    coordinates.Point
	Filename string
}

func (e *DataUnavailableError) Error() string {
	return fmt.Sprintf("Elevation data is not available for (%f, %f) (filename %s not found)",
		e.Point.Latitude, e.Point.Longitude, e.Filename)
}

// CreateManager returns a new SRTM file manager for the given directory
func CreateManager(dataStore string) (manager Manager, err error) {
	manager.dataStore = dataStore
	manager.coverage = make(map[string]bool)
	manager.files = make(map[string]File)

	// Get a list of files
	files, err := ioutil.ReadDir(dataStore)
	if err != nil {
		return
	}

	// Add each file to the coverage map
	for _, file := range files {
		filename := file.Name()
		if len(filename) < 7 {
			fmt.Printf("[SRTM-Manager] File '%s' is not named correctly, skipping...\n", filename)
			continue
		}

		// Check if the filename is correct
		_, _, err = FilenameToCoordinates(filename[:7])
		if err != nil {
			fmt.Printf("[SRTM-Manager] File '%s' is not named correctly, skipping...\n", filename)
			continue
		}

		manager.coverage[filename[:7]] = true
	}

	return manager, nil
}

// IsDataAvailable returns true if the given coordinates are covered by the data of this SRTM manager. It also returns
// the filename that is supposed to contain the requested data
func (manager Manager) IsDataAvailable(point coordinates.Point) (availability bool, filename string) {
	// Convert the given coordinates to a filename so we know which file to look for
	filename = CoordinatesToFilename(point)

	// Check coverage
	i, ok := manager.coverage[filename]
	return ok && i, filename // For coverage the key must exist and it must be true
}

// GetElevation returns the elevation at a given point. This function will check if the required files are in memory
// and will load anything that's missing
func (manager Manager) GetElevation(point coordinates.Point) (elevation int, err error) {
	// First, let's check availability at the point given
	availability, filename := manager.IsDataAvailable(point)
	if !availability {
		return -1, &DataUnavailableError{Point: point, Filename: filename}
	}

	// The file might already be loaded so let's check
	_, loaded := manager.files[filename]
	if !loaded {
		file, err := LoadFile(filepath.Join(manager.dataStore, filename + ".hgt"))
		if err != nil {
			return -1, err
		}

		manager.files[filename] = *file
	}

	// Grab the file and get elevation
	return manager.files[filename].GetElevation(point.Unwrap())
}

// Clear removes all loaded files from memory
func (manager Manager) Clear() {
	for key := range manager.files {
		delete(manager.files, key)
	}
}
