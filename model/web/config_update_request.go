package web

type ConfigUpdateRequest struct {
	Data map[string]interface{} `json:"data" validate:"required"`
}
