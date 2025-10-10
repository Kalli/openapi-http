package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kalli/openapi-http/internal/generator"
	"github.com/kalli/openapi-http/internal/parser"
)

func main() {
	var operationID string
	var path string
	var outputFile string

	flag.StringVar(&operationID, "operation-id", "", "operation id to generate")
	flag.StringVar(&path, "path", "", "path to generate (e.g. /pet)")
	flag.StringVar(&outputFile, "output", "", "output file (default: stdout)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "usage: openapi-http [flags] <spec-file>\n")
		fmt.Fprintf(os.Stderr, "  flags must come before spec-file\n")
		fmt.Fprintf(os.Stderr, "\nflags:\n")
		flag.PrintDefaults()
		os.Exit(0)
	}

	specPath := flag.Arg(0)
	spec, err := parser.LoadSpec(specPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading spec: %v\n", err)
		os.Exit(1)
	}

	// no filters → just list operations
	if operationID == "" && path == "" {
		parser.ListOperations(spec)
		return
	}

	ops := parser.FindOperations(spec, operationID, path)
	if len(ops) == 0 {
		fmt.Fprintf(os.Stderr, "no operations found\n")
		os.Exit(1)
	}

	var output *os.File
	if outputFile != "" {
		output, err = os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}

	gen := generator.NewGenerator(spec)
	for i, op := range ops {
		if i > 0 {
			fmt.Fprintln(output, "")
		}
		req, err := gen.BuildHTTPRequest(op)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error generating request: %v\n", err)
			continue
		}
		fmt.Fprint(output, req)
	}
}
