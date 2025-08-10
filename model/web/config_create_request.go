package web

type ConfigCreateRequest struct {
	Schema string                 `validate:"required"`
	Name   string                 `validate:"required"`
	Data   map[string]interface{} `json:"data" validate:"required"` // The configuration JSON object
}
