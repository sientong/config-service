package web

import "time"

type ConfigResponse struct {
	Schema    string                 `json:"schema"`
	Name      string                 `json:"name"`
	Version   int                    `json:"version"`
	Data      map[string]interface{} `json:"data"` // raw JSON
	CreatedAt time.Time              `json:"created_at"`
}

type ConfigResponses struct {
	Schema         string           `json:"schema"`
	Name           string           `json:"name"`
	ConfigVersions []ConfigResponse `json:"configVersions"`
}
