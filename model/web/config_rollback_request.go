package web

type ConfigRollbackRequest struct {
	Schema  string `validate:"required"`
	Name    string `validate:"required"`
	Version int    `validate:"required"`
}
