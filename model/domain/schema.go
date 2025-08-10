package domain

// Hardcoded JSON Schemas for different config types. Add more as needed.

var Schemas = map[string]string{
	"payment_config": `{
        "$schema": "http://json-schema.org/draft-07/schema#",
        "type": "object",
        "properties": {
            "max_limit": { "type": "integer" },
            "enabled": { "type": "boolean" }
        },
        "required": ["max_limit", "enabled"],
        "additionalProperties": false
    }`,

	"feature_flags": `{
        "$schema": "http://json-schema.org/draft-07/schema#",
        "type": "object",
        "patternProperties": {
            "^[a-zA-Z0-9_\\-]+$": { "type": "boolean" }
        },
        "additionalProperties": false
    }`,
}
