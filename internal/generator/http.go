package generator

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kalli/openapi-http/internal/parser"

	"github.com/getkin/kin-openapi/openapi3"
)

type Generator struct {
	spec *openapi3.T
}

func NewGenerator(spec *openapi3.T) *Generator {
	return &Generator{spec: spec}
}

// generates an HTTP request in rfc9110 compliant .http file format from an OpenAPI operation,
// including method, URL, headers, and body.
// Adds @name attributes based on the operationID
func (g *Generator) BuildHTTPRequest(op parser.Operation) (string, error) {
	var sb strings.Builder

	sb.WriteString("###\n")

	// add @name if operationId exists
	if op.Operation.OperationID != "" {
		sb.WriteString(fmt.Sprintf("# @name %s\n", op.Operation.OperationID))
	}

	// add summary as comment if present
	if op.Operation.Summary != "" {
		sb.WriteString(fmt.Sprintf("# %s\n", op.Operation.Summary))
	}

	sb.WriteString("\n")

	// request line
	baseURL := g.getBaseURL()
	path := g.buildPath(op)
	query := g.buildQueryString(op)

	sb.WriteString(fmt.Sprintf("%s %s%s", op.Method, baseURL, path))
	if query != "" {
		sb.WriteString("?" + query)
	}
	sb.WriteString("\n")

	// headers
	headers := g.buildHeaders(op)
	for k, v := range headers {
		sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}

	// request body
	body, err := g.buildRequestBody(op)
	if err != nil {
		return "", err
	}
	if body != "" {
		sb.WriteString("\n")
		sb.WriteString(body)
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// gets a baseurl based on the servers in the spec, uses generic hostname
// if no servers exist
func (g *Generator) getBaseURL() string {
	if len(g.spec.Servers) > 0 {
		return g.spec.Servers[0].URL // use placeholder for first server
	}
	return "{{hostname}}"
}

// builds path with example values for parameters 
func (g *Generator) buildPath(op parser.Operation) string {
	path := op.Path

	// collect path params
	params := g.collectParameters(op, "path")

	for _, param := range params {
		placeholder := fmt.Sprintf("{{%s}}", param.Name)

		// try to get example value
		if param.Example != nil {
			placeholder = fmt.Sprintf("%v", param.Example)
		} else if param.Schema != nil && param.Schema.Value != nil {
			if example := g.generateExample(param.Schema.Value); example != nil {
				placeholder = fmt.Sprintf("%v", example)
			}
		}

		path = strings.ReplaceAll(path, fmt.Sprintf("{%s}", param.Name), placeholder)
	}

	return path
}

// builds a querystring with example values
func (g *Generator) buildQueryString(op parser.Operation) string {
	params := g.collectParameters(op, "query")
	if len(params) == 0 {
		return ""
	}

	var parts []string
	for _, param := range params {
		value := "{{" + param.Name + "}}"

		if param.Example != nil {
			value = fmt.Sprintf("%v", param.Example)
		} else if param.Schema != nil && param.Schema.Value != nil {
			if example := g.generateExample(param.Schema.Value); example != nil {
				value = fmt.Sprintf("%v", example)
			}
		}

		parts = append(parts, fmt.Sprintf("%s=%s", param.Name, value))
	}

	return strings.Join(parts, "&")
}

// builds headers defined for the request
func (g *Generator) buildHeaders(op parser.Operation) map[string]string {
	headers := make(map[string]string)

	// content-type from request body
	if op.Operation.RequestBody != nil && op.Operation.RequestBody.Value != nil {

		// if request body has json, prefer that. Seems likely to be most common use case.
		if _, ok := op.Operation.RequestBody.Value.Content["application/json"]; ok {
			headers["Content-Type"] = "application/json"
		} else {
			for contentType := range op.Operation.RequestBody.Value.Content {
				headers["Content-Type"] = contentType
				break
			}
		}
	}

	// header params
	params := g.collectParameters(op, "header")
	for _, param := range params {
		value := "{{" + param.Name + "}}"
		if param.Example != nil {
			value = fmt.Sprintf("%v", param.Example)
		}
		headers[param.Name] = value
	}

	return headers
}

// builds a request body based on the operation schema
// todo: Currently limited to json requests, add support for other types
func (g *Generator) buildRequestBody(op parser.Operation) (string, error) {
	if op.Operation.RequestBody == nil || op.Operation.RequestBody.Value == nil {
		return "", nil
	}

	rb := op.Operation.RequestBody.Value

	// prefer application/json
	var mediaType *openapi3.MediaType
	if mt, ok := rb.Content["application/json"]; ok {
		mediaType = mt
	} else {
		// fallback to first available
		for _, mt := range rb.Content {
			mediaType = mt
			break
		}
	}

	if mediaType == nil {
		return "", nil
	}

	// try to get example
	var data interface{}
	if mediaType.Example != nil {
		data = mediaType.Example
	} else if mediaType.Examples != nil && len(mediaType.Examples) > 0 {
		// use first example
		for _, ex := range mediaType.Examples {
			if ex.Value != nil {
				data = ex.Value.Value
				break
			}
		}
	} else if mediaType.Schema != nil && mediaType.Schema.Value != nil {
		// generate from schema
		data = g.generateExample(mediaType.Schema.Value)
	}

	if data == nil {
		return "{}", nil
	}

	// format as json
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// collects all parameters for an operation, either path or op level
func (g *Generator) collectParameters(op parser.Operation, in string) []*openapi3.Parameter {
	var params []*openapi3.Parameter

	// path-level params
	if op.PathItem.Parameters != nil {
		for _, paramRef := range op.PathItem.Parameters {
			if paramRef.Value != nil && paramRef.Value.In == in {
				params = append(params, paramRef.Value)
			}
		}
	}

	// operation-level params (override path-level)
	if op.Operation.Parameters != nil {
		for _, paramRef := range op.Operation.Parameters {
			if paramRef.Value != nil && paramRef.Value.In == in {
				params = append(params, paramRef.Value)
			}
		}
	}

	return params
}
