package v1

import "github.com/gromnsk/money.io/pkg/errors"

// HTTP headers
const (
	HeaderContentType = "Content-Type"
)

// ContentTypes
const (
	// API response mime type
	CTApplicationJSON = "application/json"

	// API Blueprint mime type
	CTApiB = "text/vnd.apiblueprint"
)

// Response is JSON API response structure
type Response struct {
	Data  interface{}           `json:"data"`
	Error *errors.ResponseError `json:"error"`
}