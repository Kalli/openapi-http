package parser
import (
	"context" 
	"fmt"
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
)

// Loads an open api spec based either on a filepath or URL
func LoadSpec(path string) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	// handle both file paths and URLs
	var doc *openapi3.T
	var err error

	if u, parseErr := url.Parse(path); parseErr == nil && (u.Scheme == "http" || u.Scheme == "https") {
		doc, err = loader.LoadFromURI(u)
	} else {
		doc, err = loader.LoadFromFile(path)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load spec: %w", err)
	}

	// validate
	if err := doc.Validate(context.Background()); err != nil {
		return nil, fmt.Errorf("spec validation failed: %w", err)
	}

	return doc, nil
}
