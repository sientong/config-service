package domain

import (
	"fmt"
	"os"
	"path/filepath"
)

// Directory to store all schemas
var SchemaDir = "schemas"

// Schemas holds schema name -> JSON schema string
var Schemas = map[string]string{}

// LoadSchemas loads all JSON schema files from a given directory into Schemas map
func LoadSchemas(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read schema directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		content, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to read schema file %s: %w", file.Name(), err)
		}

		// Use the file name (without .json) as the schema key
		key := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
		Schemas[key] = string(content)
	}

	return nil
}
