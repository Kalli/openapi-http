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

// ANSI color codes
const (
	colorReset      = "\033[0m"
	colorBlue       = "\033[34m"  // GET
	colorRed        = "\033[31m"  // DELETE
	colorLightGreen = "\033[92m"  // PATCH
	colorDarkGreen  = "\033[32m"  // POST
	colorYellow     = "\033[93m"  // PUT (bright yellow, universally supported)
	colorBold       = "\033[1m"   // Bold for headers
)

// getMethodColor returns the ANSI color code for a given HTTP method
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return colorBlue
	case "DELETE":
		return colorRed
	case "PATCH":
		return colorLightGreen
	case "POST":
		return colorDarkGreen
	case "PUT":
		return colorYellow
	default:
		return colorReset
	}
}

// List the available operations in a spec, format:
// $HTTPMethod $Path $OperationId - $summary
func ListOperations(spec *openapi3.T) {
	fmt.Print("available operations:\n\n")

	// Print table headers
	fmt.Printf("  %s%-8s %-30s %-25s %s%s\n", colorBold, "METHOD", "PATH", "OPERATIONID", "SUMMARY", colorReset)
	fmt.Printf("  %s%s%s\n", colorBold, "────────────────────────────────────────────────────────────────────────────────────", colorReset)

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

			color := getMethodColor(method)
			fmt.Printf("  %s%-8s%s %-30s %s%-25s%s %s\n", color, method, colorReset, path, colorBold, opID, colorReset, summary)
		}
	}
}

// hasTag checks if an operation has a specific tag
func hasTag(op *openapi3.Operation, tag string) bool {
	for _, t := range op.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// find all the operations for a given path, operationId, or tag
func FindOperations(spec *openapi3.T, operationID, path, tag string) []Operation {
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

			// filter by tag if specified
			if tag != "" && !hasTag(op, tag) {
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
