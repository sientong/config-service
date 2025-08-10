package domain

import "time"

type ConfigRecord struct {
	Schema    string                 `json:"schema"`
	Name      string                 `json:"name"`
	Version   int                    `json:"version"`
	Data      map[string]interface{} `json:"data"` // raw JSON
	CreatedAt time.Time              `json:"created_at"`
}
