package web

type ConfigCreateRequest struct {
	Data map[string]interface{} `json:"data" validate:"required"` // The configuration JSON object
}
