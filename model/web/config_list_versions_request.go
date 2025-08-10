package web

type ConfigListVersionsRequest struct {
	Schema string `validate:"required"`
	Name   string `validate:"required"`
}
