package parser

import (
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

func TestFindOperations_ByOperationID(t *testing.T) {
	spec, err := LoadSpec("../../test/petstore.yml")
	if err != nil {
		t.Fatalf("failed to load spec: %v", err)
	}

	// Test finding a specific operation by ID
	ops := FindOperations(spec, "addPet", "")
	
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
	ops := FindOperations(spec, "", "/pet")
	
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
	ops := FindOperations(spec, "updatePet", "/pet")
	
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

func TestFindOperations_NotFound(t *testing.T) {
	spec, err := LoadSpec("../../test/petstore.yml")
	if err != nil {
		t.Fatalf("failed to load spec: %v", err)
	}

	// Test finding non-existent operation
	ops := FindOperations(spec, "nonExistentOperation", "")
	
	if len(ops) != 0 {
		t.Errorf("expected 0 operations, got %d", len(ops))
	}
}

func TestFindOperations_WithPathParameters(t *testing.T) {
	spec, err := LoadSpec("../../test/petstore.yml")
	if err != nil {
		t.Fatalf("failed to load spec: %v", err)
	}

	// Test finding operation with path parameters
	ops := FindOperations(spec, "getPetById", "")
	
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
