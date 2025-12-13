package parser

import (
	"slices"
	"testing"
)

func TestListOperations(t *testing.T) {
	spec, err := LoadSpec("../../test/petstore.yml")
	if err != nil {
		t.Fatalf("failed to load spec: %v", err)
	}

	// This function prints to stdout, so we just ensure it doesn't panic
	ListOperations(spec)
}

func TestFindOperations_All(t *testing.T) {
	spec, err := LoadSpec("../../test/petstore.yml")
	if err != nil {
		t.Fatalf("failed to load spec: %v", err)
	}

	// Test finding all operations (empty filters)
	ops := FindOperations(spec, "", "", "")

	if len(ops) == 0 {
		t.Fatal("expected to find operations in spec")
	}

	if len(ops) != 20 {
		t.Errorf("expected 20 operations in petstore, got %d", len(ops))
	}

	// Verify all returned items are valid operations
	for _, op := range ops {
		if op.Path == "" {
			t.Error("found operation with empty path")
		}
		if op.Method == "" {
			t.Error("found operation with empty method")
		}
		if op.Operation == nil {
			t.Error("found operation with nil Operation")
		}
	}

	// Verify we have diverse methods (not just one method)
	methods := make(map[string]bool)
	for _, op := range ops {
		methods[op.Method] = true
	}
	if len(methods) != 4 {
		t.Errorf("expected 4 HTTP methods, got only %v", methods)
	}

	// Verify we have diverse paths (not just one path)
	paths := make(map[string]bool)
	for _, op := range ops {
		paths[op.Path] = true
	}
	if len(paths) < 2 {
		t.Errorf("expected multiple paths, got only %v", paths)
	}
}

func TestFindOperations_ByOperationID(t *testing.T) {
	spec, err := LoadSpec("../../test/petstore.yml")
	if err != nil {
		t.Fatalf("failed to load spec: %v", err)
	}

	// Test finding a specific operation by ID
	ops := FindOperations(spec, "addPet", "", "")

	if len(ops) == 0 {
		t.Fatal("expected to find addPet operation")
	}

	if len(ops) != 1 {
		t.Errorf("expected 1 operation, got %d", len(ops))
	}

	op := ops[0]
	if op.Operation.OperationID != "addPet" {
		t.Errorf("expected operationId 'addPet', got '%s'", op.Operation.OperationID)
	}

	if op.Method != "POST" {
		t.Errorf("expected method POST, got %s", op.Method)
	}

	if op.Path != "/pet" {
		t.Errorf("expected path '/pet', got '%s'", op.Path)
	}
}

func TestFindOperations_ByPath(t *testing.T) {
	spec, err := LoadSpec("../../test/petstore.yml")
	if err != nil {
		t.Fatalf("failed to load spec: %v", err)
	}

	// Test finding all operations for a path
	ops := FindOperations(spec, "", "/pet", "")

	if len(ops) == 0 {
		t.Fatal("expected to find operations for /pet path")
	}

	// /pet should have POST and PUT operations
	if len(ops) < 2 {
		t.Errorf("expected at least 2 operations for /pet, got %d", len(ops))
	}

	// Verify all operations have the correct path
	for _, op := range ops {
		if op.Path != "/pet" {
			t.Errorf("expected path '/pet', got '%s'", op.Path)
		}
	}
}

func TestFindOperations_ByBoth(t *testing.T) {
	spec, err := LoadSpec("../../test/petstore.yml")
	if err != nil {
		t.Fatalf("failed to load spec: %v", err)
	}

	// Test finding by both operationId and path
	ops := FindOperations(spec, "updatePet", "/pet", "")

	if len(ops) != 1 {
		t.Fatalf("expected 1 operation, got %d", len(ops))
	}

	op := ops[0]
	if op.Operation.OperationID != "updatePet" {
		t.Errorf("expected operationId 'updatePet', got '%s'", op.Operation.OperationID)
	}

	if op.Path != "/pet" {
		t.Errorf("expected path '/pet', got '%s'", op.Path)
	}

	if op.Method != "PUT" {
		t.Errorf("expected method PUT, got %s", op.Method)
	}
}

func TestFindOperations_ByTag(t *testing.T) {
	spec, err := LoadSpec("../../test/petstore.yml")
	if err != nil {
		t.Fatalf("failed to load spec: %v", err)
	}

	// Test finding operations by tag
	ops := FindOperations(spec, "", "", "pet")

	if len(ops) == 0 {
		t.Fatal("expected to find operations with 'pet' tag")
	}

	// Verify all operations have the correct tag
	for _, op := range ops {
		if !slices.Contains(op.Operation.Tags, "pet") {
			t.Errorf("operation %s does not have 'pet' tag", op.Operation.OperationID)
		}
	}

}

func TestFindOperations_ByTagAndPath(t *testing.T) {
	spec, err := LoadSpec("../../test/petstore.yml")
	if err != nil {
		t.Fatalf("failed to load spec: %v", err)
	}

	// Test finding operations by both tag and path
	ops := FindOperations(spec, "", "/pet", "pet")

	if len(ops) == 0 {
		t.Fatal("expected to find operations with 'pet' tag at /pet path")
	}

	// Verify all operations have the correct path and tag
	for _, op := range ops {
		if op.Path != "/pet" {
			t.Errorf("expected path '/pet', got '%s'", op.Path)
		}
		if !slices.Contains(op.Operation.Tags, "pet") {
			t.Errorf("operation %s does not have 'pet' tag", op.Operation.OperationID)
		}
	}
}

func TestFindOperations_NotFound(t *testing.T) {
	spec, err := LoadSpec("../../test/petstore.yml")
	if err != nil {
		t.Fatalf("failed to load spec: %v", err)
	}

	// Test finding non-existent operation
	ops := FindOperations(spec, "nonExistentOperation", "", "")

	if len(ops) != 0 {
		t.Errorf("expected 0 operations, got %d", len(ops))
	}

	// Test finding non-existent tag
	ops = FindOperations(spec, "", "", "nonExistentTag")

	if len(ops) != 0 {
		t.Errorf("expected 0 operations for non-existent tag, got %d", len(ops))
	}
}

func TestFindOperations_WithPathParameters(t *testing.T) {
	spec, err := LoadSpec("../../test/petstore.yml")
	if err != nil {
		t.Fatalf("failed to load spec: %v", err)
	}

	// Test finding operation with path parameters
	ops := FindOperations(spec, "getPetById", "", "")

	if len(ops) != 1 {
		t.Fatalf("expected 1 operation, got %d", len(ops))
	}

	op := ops[0]
	if op.Path != "/pet/{petId}" {
		t.Errorf("expected path '/pet/{petId}', got '%s'", op.Path)
	}

	// Verify it has parameters
	if len(op.Operation.Parameters) == 0 {
		t.Error("expected operation to have parameters")
	}
}
