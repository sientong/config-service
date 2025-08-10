package web

type ConfigUpdateRequest struct {
	Schema string                 `validate:"required"`
	Name   string                 `validate:"required"`
	Data   map[string]interface{} `json:"data" validate:"required"`
}
