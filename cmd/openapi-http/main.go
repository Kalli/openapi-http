package main

import (
	"fmt"
	"os"

	"github.com/kalli/openapi-http/internal/generator"
	"github.com/kalli/openapi-http/internal/parser"
	flag "github.com/spf13/pflag"
)

func main() {
	var operationID string
	var path string
	var tag string
	var outputFile string
	var all bool
	
	// Set custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "openapi-http - Generate HTTP requests from OpenAPI specs\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  openapi-http [flags] <spec-file>\n")
		fmt.Fprintf(os.Stderr, "  openapi-http <spec-file> [operation-id] [path]\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  openapi-http spec.yaml                    # list all operations\n")
		fmt.Fprintf(os.Stderr, "  openapi-http -i getPet spec.yaml          # generate request for operation\n")
		fmt.Fprintf(os.Stderr, "  openapi-http -p /pet spec.yaml            # generate requests for path\n")
		fmt.Fprintf(os.Stderr, "  openapi-http -t pet spec.yaml             # generate requests for tag\n")
		fmt.Fprintf(os.Stderr, "  openapi-http -a spec.yaml                 # generate all requests\n")
		fmt.Fprintf(os.Stderr, "  openapi-http spec.yaml getPet             # positional arguments\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.StringVarP(&operationID, "operation-id", "i", "", "operation id to generate request for")
	flag.StringVarP(&path, "path", "p", "", "path to generate requests for (e.g. /pet)")
	flag.StringVarP(&tag, "tag", "t", "", "tag to filter operations by (e.g. pet)")
	flag.StringVarP(&outputFile, "output", "o", "", "output file (default: stdout)")
	flag.BoolVarP(&all, "all", "a", false, "generate requests for all operations")
	flag.Parse()
	
	

	// Handle positional arguments
	args := flag.Args()

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

	// If operation-id not set via flag, try positional arg
	if operationID == "" && len(args) > 1 {
		operationID = args[1]
	}

	// if --all flag is set, generate all requests
	if all {
		ops := parser.FindOperations(spec, "", "", "")
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
		return
	}

	// no filters → just list operations
	if operationID == "" && path == "" && tag == "" {
		parser.ListOperations(spec)
		return
	}

	ops := parser.FindOperations(spec, operationID, path, tag)
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

func helpText(){
    
}
