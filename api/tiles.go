package api

import (
	"encoding/json"
	"net/http"
)

// HandleTiles handles the '/tiles' endpoint which returns the list of available
// SRTM tiles to the caller
func HandleTiles(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	encoder.Encode(srtmManager.GetCoverage())
}
