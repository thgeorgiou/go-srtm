package coordinates

import "math"

// EarthRadius is the radium of the earth in meters
const EarthRadius = 6371e3

// Point represents a point on earth
type Point struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

// ToRadians converts a number in degrees to radians
func ToRadians(f float64) float64 {
	return f * math.Pi / 180
}

// DistanceTo target from point (meters)
func (coords Point) DistanceTo(target Point) float64 {
	// Convert coordinates to radians
	var φ1 = ToRadians(coords.Latitude)
	var φ2 = ToRadians(target.Latitude)

	// Get distances in radians
	var Δφ = ToRadians(target.Latitude - coords.Latitude)
	var Δλ = ToRadians(target.Longitude - coords.Longitude)

	// Calculate components
	var a = math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*math.Sin(Δλ/2)*math.Sin(Δλ/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Return distance
	return EarthRadius * c
}

// Unwrap returns the latitude and longitude of a Coordinates struct
func (coords Point) Unwrap() (latitude, longitude float64) {
	return coords.Latitude, coords.Longitude
}

// Lerp Linear interpolation between two points for a parameter t
func Lerp(from, to Point, t float64) (output Point) {
	output = Point{
		Latitude:  from.Latitude + t*(to.Latitude-from.Latitude),
		Longitude: from.Longitude + t*(to.Longitude-from.Longitude),
	}
	return
}
