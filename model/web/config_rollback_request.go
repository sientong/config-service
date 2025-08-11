package web

type ConfigRollbackRequest struct {
	Version int `validate:"required" json:"version"`
}
