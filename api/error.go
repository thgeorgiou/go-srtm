package api

// Error represents an API error. This struct is created and then converted
// to JSON before being sent to the client
type Error struct {
	Name        string      `json:"error_name"`
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
}
