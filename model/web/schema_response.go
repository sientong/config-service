package web

type SchemaResponse struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Directory string `json:"directory"`
}

type SchemaListResponse struct {
	Schemas []SchemaResponse `json:"schemas"`
}
