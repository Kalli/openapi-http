package generator

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestGenerateExample_String(t *testing.T) {
	gen := &Generator{}

	schema := &openapi3.Schema{
		Type: &openapi3.Types{"string"},
	}

	result := gen.generateExample(schema)

	if result != "string" {
		t.Errorf("expected 'string', got: %v", result)
	}
}

func TestGenerateExample_StringWithFormat(t *testing.T) {
	tests := []struct {
		format   string
		expected string
	}{
		{"date", "2024-01-01"},
		{"date-time", "2024-01-01T00:00:00Z"},
		{"email", "user@example.com"},
	}

	gen := &Generator{}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			schema := &openapi3.Schema{
				Type:   &openapi3.Types{"string"},
				Format: tt.format,
			}

			result := gen.generateExample(schema)

			if result != tt.expected {
				t.Errorf("expected %s, got: %v", tt.expected, result)
			}
		})
	}
}

func TestGenerateExample_Integer(t *testing.T) {
	gen := &Generator{}

	schema := &openapi3.Schema{
		Type: &openapi3.Types{"integer"},
	}

	result := gen.generateExample(schema)

	if result != 0 {
		t.Errorf("expected 0, got: %v", result)
	}
}

func TestGenerateExample_IntegerWithMin(t *testing.T) {
	gen := &Generator{}

	minVal := 10.0
	schema := &openapi3.Schema{
		Type: &openapi3.Types{"integer"},
		Min:  &minVal,
	}

	result := gen.generateExample(schema)

	if result != 10 {
		t.Errorf("expected 10, got: %v", result)
	}
}

func TestGenerateExample_Number(t *testing.T) {
	gen := &Generator{}

	schema := &openapi3.Schema{
		Type: &openapi3.Types{"number"},
	}

	result := gen.generateExample(schema)

	if result != 0.0 {
		t.Errorf("expected 0.0, got: %v", result)
	}
}

func TestGenerateExample_Boolean(t *testing.T) {
	gen := &Generator{}

	schema := &openapi3.Schema{
		Type: &openapi3.Types{"boolean"},
	}

	result := gen.generateExample(schema)

	if result != false {
		t.Errorf("expected false, got: %v", result)
	}
}

func TestGenerateExample_Array(t *testing.T) {
	gen := &Generator{}

	schema := &openapi3.Schema{
		Type: &openapi3.Types{"array"},
		Items: &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:    &openapi3.Types{"string"},
				Example: "item",
			},
		},
	}

	result := gen.generateExample(schema)

	arr, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected array, got: %T", result)
	}

	if len(arr) != 1 {
		t.Errorf("expected array with 1 item, got: %d", len(arr))
	}

	if arr[0] != "item" {
		t.Errorf("expected 'item', got: %v", arr[0])
	}
}

func TestGenerateExample_Object(t *testing.T) {
	gen := &Generator{}

	schema := &openapi3.Schema{
		Type: &openapi3.Types{"object"},
		Properties: openapi3.Schemas{
			"name": &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:    &openapi3.Types{"string"},
					Example: "John",
				},
			},
			"age": &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:    &openapi3.Types{"integer"},
					Example: 30,
				},
			},
		},
	}

	result := gen.generateExample(schema)

	obj, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected object, got: %T", result)
	}

	if obj["name"] != "John" {
		t.Errorf("expected name 'John', got: %v", obj["name"])
	}

	if obj["age"] != 30 {
		t.Errorf("expected age 30, got: %v", obj["age"])
	}
}

func TestGenerateExample_NestedObject(t *testing.T) {
	gen := &Generator{}

	schema := &openapi3.Schema{
		Type: &openapi3.Types{"object"},
		Properties: openapi3.Schemas{
			"user": &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{"object"},
					Properties: openapi3.Schemas{
						"name": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{"string"},
								Example: "Jane",
							},
						},
					},
				},
			},
		},
	}

	result := gen.generateExample(schema)

	obj, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected object, got: %T", result)
	}

	user, ok := obj["user"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected nested user object, got: %T", obj["user"])
	}

	if user["name"] != "Jane" {
		t.Errorf("expected name 'Jane', got: %v", user["name"])
	}
}

func TestGenerateExample_WithExplicitExample(t *testing.T) {
	gen := &Generator{}

	schema := &openapi3.Schema{
		Type:    &openapi3.Types{"string"},
		Example: "custom-example",
	}

	result := gen.generateExample(schema)

	if result != "custom-example" {
		t.Errorf("expected 'custom-example', got: %v", result)
	}
}

func TestGenerateExample_WithDefault(t *testing.T) {
	gen := &Generator{}

	schema := &openapi3.Schema{
		Type:    &openapi3.Types{"string"},
		Default: "default-value",
	}

	result := gen.generateExample(schema)

	if result != "default-value" {
		t.Errorf("expected 'default-value', got: %v", result)
	}
}

func TestGenerateExample_WithEnum(t *testing.T) {
	gen := &Generator{}

	schema := &openapi3.Schema{
		Type: &openapi3.Types{"string"},
		Enum: []interface{}{"active", "inactive", "pending"},
	}

	result := gen.generateExample(schema)

	if result != "active" {
		t.Errorf("expected first enum value 'active', got: %v", result)
	}
}

func TestGenerateExample_PriorityOrder(t *testing.T) {
	gen := &Generator{}

	// Example should take priority over default and enum
	schema := &openapi3.Schema{
		Type:    &openapi3.Types{"string"},
		Example: "explicit-example",
		Default: "default-value",
		Enum:    []interface{}{"enum1", "enum2"},
	}

	result := gen.generateExample(schema)

	if result != "explicit-example" {
		t.Errorf("expected example to take priority, got: %v", result)
	}
}

func TestGenerateExample_EmptyArray(t *testing.T) {
	gen := &Generator{}

	schema := &openapi3.Schema{
		Type: &openapi3.Types{"array"},
		// No items defined
	}

	result := gen.generateExample(schema)

	arr, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected array, got: %T", result)
	}

	if len(arr) != 0 {
		t.Errorf("expected empty array, got length: %d", len(arr))
	}
}

func TestGenerateExample_EmptyObject(t *testing.T) {
	gen := &Generator{}

	schema := &openapi3.Schema{
		Type: &openapi3.Types{"object"},
		// No properties defined
	}

	result := gen.generateExample(schema)

	obj, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected object, got: %T", result)
	}

	if len(obj) != 0 {
		t.Errorf("expected empty object, got: %v", obj)
	}
}

func TestGenerateExample_AdditionalProperties(t *testing.T) {
	gen := &Generator{}

	hasAdditional := true
	schema := &openapi3.Schema{
		Type: &openapi3.Types{"object"},
		AdditionalProperties: openapi3.AdditionalProperties{
			Has: &hasAdditional,
		},
	}

	result := gen.generateExample(schema)

	obj, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected object, got: %T", result)
	}

	if obj["key"] != "value" {
		t.Errorf("expected example key-value pair, got: %v", obj)
	}
}

func TestGetSchemaType(t *testing.T) {
	tests := []struct {
		name     string
		types    *openapi3.Types
		expected string
	}{
		{"string", &openapi3.Types{"string"}, "string"},
		{"integer", &openapi3.Types{"integer"}, "integer"},
		{"number", &openapi3.Types{"number"}, "number"},
		{"boolean", &openapi3.Types{"boolean"}, "boolean"},
		{"array", &openapi3.Types{"array"}, "array"},
		{"object", &openapi3.Types{"object"}, "object"},
		{"null", &openapi3.Types{"null"}, "null"},
		{"nil type", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := &openapi3.Schema{
				Type: tt.types,
			}

			result := getSchemaType(schema)

			if result != tt.expected {
				t.Errorf("expected %s, got: %s", tt.expected, result)
			}
		})
	}
}
