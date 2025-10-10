package parser

import (
	"fmt"
	"sort"

	"github.com/getkin/kin-openapi/openapi3"
)

type Operation struct {
	Path      string
	Method    string
	Operation *openapi3.Operation
	PathItem  *openapi3.PathItem
}

// List the available operations in a spec, format:
// $HTTPMethod $Path $OperationId - $summary
func ListOperations(spec *openapi3.T) {
	fmt.Print("available operations:\n\n")

	var paths []string
	for path := range spec.Paths.Map() {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		pathItem := spec.Paths.Map()[path]

		for _, method := range []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"} {
			op := pathItem.GetOperation(method)
			if op == nil {
				continue
			}

			opID := op.OperationID
			if opID == "" {
				opID = "(no id)"
			}
			summary := op.Summary
			if summary == "" {
				summary = "(no summary)"
			}

			fmt.Printf("  %-8s %-30s %s - %s\n", method, path, opID, summary)
		}
	}
}

// find all the operations for a given path or identified by an operationId
func FindOperations(spec *openapi3.T, operationID, path string) []Operation {
	var results []Operation

	for p, pathItem := range spec.Paths.Map() {
		// filter by path if specified
		if path != "" && p != path {
			continue
		}

		for _, method := range []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"} {
			op := pathItem.GetOperation(method)
			if op == nil {
				continue
			}

			// filter by operation id if specified
			if operationID != "" && op.OperationID != operationID {
				continue
			}

			results = append(results, Operation{
				Path:      p,
				Method:    method,
				Operation: op,
				PathItem:  pathItem,
			})
		}
	}

	return results
}