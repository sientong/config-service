package web

type ConfigFetchRequest struct {
	Schema  string `validate:"required"`
	Name    string `validate:"required"`
	Version int
}
