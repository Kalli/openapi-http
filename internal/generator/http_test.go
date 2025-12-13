package generator

import (
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/kalli/openapi-http/internal/parser"
)

func TestBuildHTTPRequest_Basic(t *testing.T) {
	spec := &openapi3.T{
		Servers: []*openapi3.Server{
			{URL: "https://api.example.com"},
		},
		Paths: openapi3.NewPaths(),
	}

	pathItem := &openapi3.PathItem{
		Get: &openapi3.Operation{
			OperationID: "getUser",
			Summary:     "Get a user",
		},
	}

	op := parser.Operation{
		Path:      "/users",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	result, err := gen.BuildHTTPRequest(op)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify basic structure
	if !strings.Contains(result, "###") {
		t.Error("expected request to start with ###")
	}

	if !strings.Contains(result, "# @name getUser") {
		t.Error("expected @name annotation")
	}

	if !strings.Contains(result, "# Get a user") {
		t.Error("expected summary as comment")
	}

	if !strings.Contains(result, "GET https://api.example.com/users") {
		t.Error("expected correct request line")
	}
}

func TestBuildHTTPRequest_WithPathParameters(t *testing.T) {
	spec := &openapi3.T{
		Servers: []*openapi3.Server{
			{URL: "https://api.example.com"},
		},
	}

	pathItem := &openapi3.PathItem{
		Get: &openapi3.Operation{
			OperationID: "getUserById",
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:    "userId",
						In:      "path",
						Example: "12345",
					},
				},
			},
		},
	}

	op := parser.Operation{
		Path:      "/users/{userId}",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	result, err := gen.BuildHTTPRequest(op)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "GET https://api.example.com/users/12345") {
		t.Errorf("expected path parameter to be replaced with example value, got: %s", result)
	}
}

func TestBuildHTTPRequest_WithQueryParameters(t *testing.T) {
	spec := &openapi3.T{
		Servers: []*openapi3.Server{
			{URL: "https://api.example.com"},
		},
	}

	pathItem := &openapi3.PathItem{
		Get: &openapi3.Operation{
			OperationID: "searchUsers",
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:    "status",
						In:      "query",
						Example: "active",
					},
				},
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:    "limit",
						In:      "query",
						Example: 10,
					},
				},
			},
		},
	}

	op := parser.Operation{
		Path:      "/users",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	result, err := gen.BuildHTTPRequest(op)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "?") {
		t.Error("expected query string")
	}

	if !strings.Contains(result, "status=active") {
		t.Error("expected status query parameter")
	}

	if !strings.Contains(result, "limit=10") {
		t.Error("expected limit query parameter")
	}
}

func TestBuildHTTPRequest_WithHeaders(t *testing.T) {
	spec := &openapi3.T{
		Servers: []*openapi3.Server{
			{URL: "https://api.example.com"},
		},
	}

	pathItem := &openapi3.PathItem{
		Get: &openapi3.Operation{
			OperationID: "getSecureData",
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:    "X-API-Key",
						In:      "header",
						Example: "secret-key-123",
					},
				},
			},
		},
	}

	op := parser.Operation{
		Path:      "/secure",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	result, err := gen.BuildHTTPRequest(op)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "X-API-Key: secret-key-123") {
		t.Error("expected X-API-Key header")
	}
}

func TestBuildHTTPRequest_WithJSONBody(t *testing.T) {
	spec := &openapi3.T{
		Servers: []*openapi3.Server{
			{URL: "https://api.example.com"},
		},
	}

	schema := openapi3.NewSchema()
	schema.Type = &openapi3.Types{"object"}
	schema.Properties = openapi3.Schemas{
		"name": &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:    &openapi3.Types{"string"},
				Example: "John Doe",
			},
		},
		"age": &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:    &openapi3.Types{"integer"},
				Example: 30,
			},
		},
	}

	pathItem := &openapi3.PathItem{
		Post: &openapi3.Operation{
			OperationID: "createUser",
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{
								Value: schema,
							},
						},
					},
				},
			},
		},
	}

	op := parser.Operation{
		Path:      "/users",
		Method:    "POST",
		Operation: pathItem.Post,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	result, err := gen.BuildHTTPRequest(op)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Content-Type: application/json") {
		t.Error("expected Content-Type header")
	}

	if !strings.Contains(result, `"name"`) {
		t.Error("expected name field in body")
	}

	if !strings.Contains(result, `"age"`) {
		t.Error("expected age field in body")
	}

	if !strings.Contains(result, `"John Doe"`) {
		t.Error("expected name value in body")
	}
}

func TestBuildHTTPRequest_NoOperationID(t *testing.T) {
	spec := &openapi3.T{
		Servers: []*openapi3.Server{
			{URL: "https://api.example.com"},
		},
	}

	pathItem := &openapi3.PathItem{
		Get: &openapi3.Operation{
			// No OperationID
			Summary: "Some operation",
		},
	}

	op := parser.Operation{
		Path:      "/test",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	result, err := gen.BuildHTTPRequest(op)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should not have @name annotation
	if strings.Contains(result, "# @name") {
		t.Error("should not have @name annotation when operationID is missing")
	}

	// Should still have summary
	if !strings.Contains(result, "# Some operation") {
		t.Error("expected summary comment")
	}
}

func TestGetBaseURL_WithServers(t *testing.T) {
	spec := &openapi3.T{
		Servers: []*openapi3.Server{
			{URL: "https://api.example.com/v1"},
			{URL: "https://staging.example.com/v1"},
		},
	}

	gen := NewGenerator(spec)
	baseURL := gen.getBaseURL()

	if baseURL != "https://api.example.com/v1" {
		t.Errorf("expected first server URL, got: %s", baseURL)
	}
}

func TestGetBaseURL_NoServers(t *testing.T) {
	spec := &openapi3.T{
		Servers: []*openapi3.Server{},
	}

	gen := NewGenerator(spec)
	baseURL := gen.getBaseURL()

	if baseURL != "{{hostname}}" {
		t.Errorf("expected placeholder hostname, got: %s", baseURL)
	}
}

func TestBuildPath_WithMultipleParameters(t *testing.T) {
	spec := &openapi3.T{}

	pathItem := &openapi3.PathItem{
		Get: &openapi3.Operation{
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:    "orgId",
						In:      "path",
						Example: "org-123",
					},
				},
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:    "projectId",
						In:      "path",
						Example: "proj-456",
					},
				},
			},
		},
	}

	op := parser.Operation{
		Path:      "/orgs/{orgId}/projects/{projectId}",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	path := gen.buildPath(op)

	expected := "/orgs/org-123/projects/proj-456"
	if path != expected {
		t.Errorf("expected path %s, got: %s", expected, path)
	}
}

func TestBuildQueryString_NoParameters(t *testing.T) {
	spec := &openapi3.T{}

	pathItem := &openapi3.PathItem{
		Get: &openapi3.Operation{},
	}

	op := parser.Operation{
		Path:      "/test",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	query := gen.buildQueryString(op)

	if query != "" {
		t.Errorf("expected empty query string, got: %s", query)
	}
}

func TestBuildHeaders_WithContentType(t *testing.T) {
	spec := &openapi3.T{}

	pathItem := &openapi3.PathItem{
		Post: &openapi3.Operation{
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{},
					},
				},
			},
		},
	}

	op := parser.Operation{
		Path:      "/test",
		Method:    "POST",
		Operation: pathItem.Post,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	headers := gen.buildHeaders(op)

	if headers["Content-Type"] != "application/json" {
		t.Errorf("expected Content-Type: application/json, got: %s", headers["Content-Type"])
	}
}

func TestBuildHeaders_PreferJSON(t *testing.T) {
	spec := &openapi3.T{}

	pathItem := &openapi3.PathItem{
		Post: &openapi3.Operation{
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Content: openapi3.Content{
						"application/xml":  &openapi3.MediaType{},
						"application/json": &openapi3.MediaType{},
					},
				},
			},
		},
	}

	op := parser.Operation{
		Path:      "/test",
		Method:    "POST",
		Operation: pathItem.Post,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	headers := gen.buildHeaders(op)

	if headers["Content-Type"] != "application/json" {
		t.Errorf("expected to prefer application/json, got: %s", headers["Content-Type"])
	}
}

func TestBuildRequestBody_WithExample(t *testing.T) {
	spec := &openapi3.T{}

	exampleData := map[string]interface{}{
		"id":   "123",
		"name": "Test",
	}

	pathItem := &openapi3.PathItem{
		Post: &openapi3.Operation{
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Example: exampleData,
						},
					},
				},
			},
		},
	}

	op := parser.Operation{
		Path:      "/test",
		Method:    "POST",
		Operation: pathItem.Post,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	body, err := gen.buildRequestBody(op)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(body, `"id"`) || !strings.Contains(body, `"123"`) {
		t.Error("expected example data in body")
	}
}

func TestBuildRequestBody_NoRequestBody(t *testing.T) {
	spec := &openapi3.T{}

	pathItem := &openapi3.PathItem{
		Get: &openapi3.Operation{
			// No request body
		},
	}

	op := parser.Operation{
		Path:      "/test",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	body, err := gen.buildRequestBody(op)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if body != "" {
		t.Errorf("expected empty body, got: %s", body)
	}
}

func TestCollectParameters_PathLevel(t *testing.T) {
	spec := &openapi3.T{}

	pathItem := &openapi3.PathItem{
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name: "pathLevelParam",
					In:   "query",
				},
			},
		},
		Get: &openapi3.Operation{},
	}

	op := parser.Operation{
		Path:      "/test",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	params := gen.collectParameters(op, "query")

	if len(params) != 1 {
		t.Fatalf("expected 1 parameter, got %d", len(params))
	}

	if params[0].Name != "pathLevelParam" {
		t.Errorf("expected pathLevelParam, got: %s", params[0].Name)
	}
}

func TestCollectParameters_OperationOverridesPath(t *testing.T) {
	spec := &openapi3.T{}

	pathItem := &openapi3.PathItem{
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name: "pathParam",
					In:   "query",
				},
			},
		},
		Get: &openapi3.Operation{
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name: "opParam",
						In:   "query",
					},
				},
			},
		},
	}

	op := parser.Operation{
		Path:      "/test",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	params := gen.collectParameters(op, "query")

	// Should have both path-level and operation-level params
	if len(params) != 2 {
		t.Fatalf("expected 2 parameters, got %d", len(params))
	}
}

func TestBuildSecurityHeaders_BearerAuth(t *testing.T) {
	spec := &openapi3.T{
		Components: &openapi3.Components{
			SecuritySchemes: openapi3.SecuritySchemes{
				"bearerAuth": &openapi3.SecuritySchemeRef{
					Value: &openapi3.SecurityScheme{
						Type:   "http",
						Scheme: "bearer",
					},
				},
			},
		},
	}

	pathItem := &openapi3.PathItem{
		Get: &openapi3.Operation{
			Security: &openapi3.SecurityRequirements{
				{"bearerAuth": []string{}},
			},
		},
	}

	op := parser.Operation{
		Path:      "/secure",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	headers := gen.buildSecurityHeaders(op)

	if headers["Authorization"] != "Bearer {{token}}" {
		t.Errorf("expected Bearer token auth, got: %s", headers["Authorization"])
	}
}

func TestBuildSecurityHeaders_BasicAuth(t *testing.T) {
	spec := &openapi3.T{
		Components: &openapi3.Components{
			SecuritySchemes: openapi3.SecuritySchemes{
				"basicAuth": &openapi3.SecuritySchemeRef{
					Value: &openapi3.SecurityScheme{
						Type:   "http",
						Scheme: "basic",
					},
				},
			},
		},
	}

	pathItem := &openapi3.PathItem{
		Get: &openapi3.Operation{
			Security: &openapi3.SecurityRequirements{
				{"basicAuth": []string{}},
			},
		},
	}

	op := parser.Operation{
		Path:      "/secure",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	headers := gen.buildSecurityHeaders(op)

	if headers["Authorization"] != "Basic {{credentials}}" {
		t.Errorf("expected Basic auth, got: %s", headers["Authorization"])
	}
}

func TestBuildSecurityHeaders_APIKeyInHeader(t *testing.T) {
	spec := &openapi3.T{
		Components: &openapi3.Components{
			SecuritySchemes: openapi3.SecuritySchemes{
				"apiKey": &openapi3.SecuritySchemeRef{
					Value: &openapi3.SecurityScheme{
						Type: "apiKey",
						In:   "header",
						Name: "X-API-Key",
					},
				},
			},
		},
	}

	pathItem := &openapi3.PathItem{
		Get: &openapi3.Operation{
			Security: &openapi3.SecurityRequirements{
				{"apiKey": []string{}},
			},
		},
	}

	op := parser.Operation{
		Path:      "/secure",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	headers := gen.buildSecurityHeaders(op)

	if headers["X-API-Key"] != "{{X-API-Key}}" {
		t.Errorf("expected X-API-Key header, got: %s", headers["X-API-Key"])
	}
}

func TestBuildSecurityHeaders_OAuth2(t *testing.T) {
	spec := &openapi3.T{
		Components: &openapi3.Components{
			SecuritySchemes: openapi3.SecuritySchemes{
				"oauth2": &openapi3.SecuritySchemeRef{
					Value: &openapi3.SecurityScheme{
						Type: "oauth2",
						Flows: &openapi3.OAuthFlows{
							AuthorizationCode: &openapi3.OAuthFlow{
								AuthorizationURL: "https://example.com/oauth/authorize",
								TokenURL:         "https://example.com/oauth/token",
								Scopes:           map[string]string{"read": "Read access"},
							},
						},
					},
				},
			},
		},
	}

	pathItem := &openapi3.PathItem{
		Get: &openapi3.Operation{
			Security: &openapi3.SecurityRequirements{
				{"oauth2": []string{"read"}},
			},
		},
	}

	op := parser.Operation{
		Path:      "/secure",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	headers := gen.buildSecurityHeaders(op)

	if headers["Authorization"] != "Bearer {{token}}" {
		t.Errorf("expected Bearer token for OAuth2, got: %s", headers["Authorization"])
	}
}

func TestBuildSecurityHeaders_OpenIDConnect(t *testing.T) {
	spec := &openapi3.T{
		Components: &openapi3.Components{
			SecuritySchemes: openapi3.SecuritySchemes{
				"openId": &openapi3.SecuritySchemeRef{
					Value: &openapi3.SecurityScheme{
						Type:             "openIdConnect",
						OpenIdConnectUrl: "https://example.com/.well-known/openid-configuration",
					},
				},
			},
		},
	}

	pathItem := &openapi3.PathItem{
		Get: &openapi3.Operation{
			Security: &openapi3.SecurityRequirements{
				{"openId": []string{}},
			},
		},
	}

	op := parser.Operation{
		Path:      "/secure",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	headers := gen.buildSecurityHeaders(op)

	if headers["Authorization"] != "Bearer {{token}}" {
		t.Errorf("expected Bearer token for OpenID Connect, got: %s", headers["Authorization"])
	}
}

func TestBuildSecurityHeaders_GlobalSecurity(t *testing.T) {
	spec := &openapi3.T{
		Components: &openapi3.Components{
			SecuritySchemes: openapi3.SecuritySchemes{
				"globalAuth": &openapi3.SecuritySchemeRef{
					Value: &openapi3.SecurityScheme{
						Type:   "http",
						Scheme: "bearer",
					},
				},
			},
		},
		Security: openapi3.SecurityRequirements{
			{"globalAuth": []string{}},
		},
	}

	pathItem := &openapi3.PathItem{
		Get: &openapi3.Operation{
			// No operation-level security, should use global
		},
	}

	op := parser.Operation{
		Path:      "/secure",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	headers := gen.buildSecurityHeaders(op)

	if headers["Authorization"] != "Bearer {{token}}" {
		t.Errorf("expected Bearer token from global security, got: %s", headers["Authorization"])
	}
}

func TestBuildSecurityHeaders_OperationOverridesGlobal(t *testing.T) {
	spec := &openapi3.T{
		Components: &openapi3.Components{
			SecuritySchemes: openapi3.SecuritySchemes{
				"globalAuth": &openapi3.SecuritySchemeRef{
					Value: &openapi3.SecurityScheme{
						Type:   "http",
						Scheme: "bearer",
					},
				},
				"opAuth": &openapi3.SecuritySchemeRef{
					Value: &openapi3.SecurityScheme{
						Type:   "http",
						Scheme: "basic",
					},
				},
			},
		},
		Security: openapi3.SecurityRequirements{
			{"globalAuth": []string{}},
		},
	}

	pathItem := &openapi3.PathItem{
		Get: &openapi3.Operation{
			Security: &openapi3.SecurityRequirements{
				{"opAuth": []string{}},
			},
		},
	}

	op := parser.Operation{
		Path:      "/secure",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	headers := gen.buildSecurityHeaders(op)

	if headers["Authorization"] != "Basic {{credentials}}" {
		t.Errorf("expected operation-level Basic auth to override global, got: %s", headers["Authorization"])
	}
}

func TestBuildSecurityHeaders_NoSecurity(t *testing.T) {
	spec := &openapi3.T{
		Components: &openapi3.Components{
			SecuritySchemes: openapi3.SecuritySchemes{},
		},
	}

	pathItem := &openapi3.PathItem{
		Get: &openapi3.Operation{
			// No security requirements
		},
	}

	op := parser.Operation{
		Path:      "/public",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	headers := gen.buildSecurityHeaders(op)

	if len(headers) != 0 {
		t.Errorf("expected no security headers, got: %v", headers)
	}
}

func TestBuildSecurityHeaders_EmptySecurityRequirement(t *testing.T) {
	spec := &openapi3.T{
		Components: &openapi3.Components{
			SecuritySchemes: openapi3.SecuritySchemes{
				"bearerAuth": &openapi3.SecuritySchemeRef{
					Value: &openapi3.SecurityScheme{
						Type:   "http",
						Scheme: "bearer",
					},
				},
			},
		},
		Security: openapi3.SecurityRequirements{
			{"bearerAuth": []string{}},
		},
	}

	pathItem := &openapi3.PathItem{
		Get: &openapi3.Operation{
			// Empty security requirement means no auth required (overrides global)
			Security: &openapi3.SecurityRequirements{},
		},
	}

	op := parser.Operation{
		Path:      "/public",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	headers := gen.buildSecurityHeaders(op)

	if len(headers) != 0 {
		t.Errorf("expected no security headers for empty security requirement, got: %v", headers)
	}
}

func TestBuildHeaders_WithSecurity(t *testing.T) {
	spec := &openapi3.T{
		Servers: []*openapi3.Server{
			{URL: "https://api.example.com"},
		},
		Components: &openapi3.Components{
			SecuritySchemes: openapi3.SecuritySchemes{
				"bearerAuth": &openapi3.SecuritySchemeRef{
					Value: &openapi3.SecurityScheme{
						Type:   "http",
						Scheme: "bearer",
					},
				},
			},
		},
	}

	pathItem := &openapi3.PathItem{
		Post: &openapi3.Operation{
			Security: &openapi3.SecurityRequirements{
				{"bearerAuth": []string{}},
			},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{},
					},
				},
			},
		},
	}

	op := parser.Operation{
		Path:      "/secure",
		Method:    "POST",
		Operation: pathItem.Post,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	headers := gen.buildHeaders(op)

	// Should have both Content-Type and Authorization
	if headers["Content-Type"] != "application/json" {
		t.Errorf("expected Content-Type header, got: %s", headers["Content-Type"])
	}

	if headers["Authorization"] != "Bearer {{token}}" {
		t.Errorf("expected Authorization header, got: %s", headers["Authorization"])
	}
}

func TestBuildHTTPRequest_WithBearerAuth(t *testing.T) {
	spec := &openapi3.T{
		Servers: []*openapi3.Server{
			{URL: "https://api.example.com"},
		},
		Components: &openapi3.Components{
			SecuritySchemes: openapi3.SecuritySchemes{
				"bearerAuth": &openapi3.SecuritySchemeRef{
					Value: &openapi3.SecurityScheme{
						Type:   "http",
						Scheme: "bearer",
					},
				},
			},
		},
	}

	pathItem := &openapi3.PathItem{
		Get: &openapi3.Operation{
			OperationID: "getSecureData",
			Summary:     "Get secure data",
			Security: &openapi3.SecurityRequirements{
				{"bearerAuth": []string{}},
			},
		},
	}

	op := parser.Operation{
		Path:      "/secure",
		Method:    "GET",
		Operation: pathItem.Get,
		PathItem:  pathItem,
	}

	gen := NewGenerator(spec)
	result, err := gen.BuildHTTPRequest(op)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Authorization: Bearer {{token}}") {
		t.Errorf("expected Authorization header in output, got:\n%s", result)
	}
}
