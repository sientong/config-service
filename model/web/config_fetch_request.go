package web

type ConfigFetchRequest struct {
	// Version is optional
	Version *int `json:"version,omitempty"`
}
