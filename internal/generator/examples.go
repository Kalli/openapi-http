package generator

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// generateExample recursively creates example values for parameters from an OpenAPI schema
// Using, in order of preference: explicit examples, defaults, or type-based generation.
func (g *Generator) generateExample(schema *openapi3.Schema) interface{} {
	// use the provided example if present
	if schema.Example != nil {
		return schema.Example
	}

	// use default value if present
	if schema.Default != nil {
		return schema.Default
	}

	// for enums, return first value
	if len(schema.Enum) > 0 {
		return schema.Enum[0]
	}

	// check type via helper
	schemaType := getSchemaType(schema)

	switch schemaType {
	case "string":
		if schema.Format == "date" {
			return "2024-01-01"
		}
		if schema.Format == "date-time" {
			return "2024-01-01T00:00:00Z"
		}
		if schema.Format == "email" {
			return "user@example.com"
		}
		return "string"

	case "integer":
		if schema.Min != nil {
			return int(*schema.Min)
		}
		return 0

	case "number":
		if schema.Min != nil {
			return *schema.Min
		}
		return 0.0

	case "boolean":
		return false

	case "array":
		if schema.Items != nil && schema.Items.Value != nil {
			item := g.generateExample(schema.Items.Value)
			return []interface{}{item}
		}
		return []interface{}{}

	case "object":
		obj := make(map[string]interface{})

		// generate for all properties
		for propName, propSchema := range schema.Properties {
			if propSchema.Value != nil {
				obj[propName] = g.generateExample(propSchema.Value)
			}
		}

		// if no properties but additionalProperties, show example
		if len(obj) == 0 && schema.AdditionalProperties.Has != nil && *schema.AdditionalProperties.Has {
			obj["key"] = "value"
		}

		return obj

	default:
		// no type specified, try object
		if len(schema.Properties) > 0 {
			obj := make(map[string]interface{})
			for propName, propSchema := range schema.Properties {
				if propSchema.Value != nil {
					obj[propName] = g.generateExample(propSchema.Value)
				}
			}
			return obj
		}
		return nil
	}
}

// helper to extract primary type from *openapi3.Types
func getSchemaType(schema *openapi3.Schema) string {
	if schema.Type == nil {
		return ""
	}

	// Type is *Types which is a slice wrapper
	if schema.Type.Includes("string") {
		return "string"
	}
	if schema.Type.Includes("integer") {
		return "integer"
	}
	if schema.Type.Includes("number") {
		return "number"
	}
	if schema.Type.Includes("boolean") {
		return "boolean"
	}
	if schema.Type.Includes("array") {
		return "array"
	}
	if schema.Type.Includes("object") {
		return "object"
	}
	if schema.Type.Includes("null") {
		return "null"
	}

	return ""
}
