package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"

	"github.com/vektah/graphql-go/schema"
)

var output = flag.String("out", "-", "the file to write to, - for stdout")
var schemaFilename = flag.String("schema", "schema.graphql", "the graphql schema to generate types from")
var typemap = flag.String("typemap", "types.json", "a json map going from graphql to golang types")
var packageName = flag.String("package", "graphql", "the package name")
var help = flag.Bool("h", false, "this usage text")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s schema.graphql\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	schema := schema.New()
	schemaRaw, err := ioutil.ReadFile(*schemaFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to open schema: "+err.Error())
		os.Exit(1)
	}

	if err = schema.Parse(string(schemaRaw)); err != nil {
		fmt.Fprintln(os.Stderr, "unable to parse schema: "+err.Error())
		os.Exit(1)
	}

	e := extractor{
		PackageName: *packageName,
		goTypeMap:   loadTypeMap(),
		schemaRaw:   string(schemaRaw),
		Imports: map[string]string{
			"strconv": "strconv",
			"fmt":     "fmt",
			"context": "context",
			"query":   "github.com/vektah/graphql-go/query",
			"schema":  "github.com/vektah/graphql-go/schema",
		},
	}
	e.extract(schema)

	// Poke a few magic methods into query
	q := e.GetObject("Query")
	q.Fields = append(q.Fields, Field{
		Type:        e.getType("__Schema").Ptr(),
		GraphQLName: "__schema",
		NoErr:       true,
		MethodName:  "ec.introspectSchema",
	})
	q.Fields = append(q.Fields, Field{
		Type:        e.getType("__Type").Ptr(),
		GraphQLName: "__type",
		NoErr:       true,
		MethodName:  "ec.introspectType",
		Args:        []Arg{{Name: "name", Type: Type{Basic: true, Name: "string"}}},
	})

	if len(e.Errors) != 0 {
		for _, err := range e.Errors {
			fmt.Println(os.Stderr, "err: "+err)
		}
		os.Exit(1)
	}

	if err := e.introspect(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	buf := &bytes.Buffer{}
	write(e, buf)

	if *output == "-" {
		fmt.Println(string(gofmt(buf.Bytes())))
	} else {
		err := ioutil.WriteFile(*output, gofmt(buf.Bytes()), 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to write output: ", err.Error())
			os.Exit(1)
		}
	}
}

func gofmt(b []byte) []byte {
	out, err := format.Source(b)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to gofmt: "+*output+":"+err.Error())
		return b
	}
	return out
}

func loadTypeMap() map[string]string {
	goTypes := map[string]string{
		"__Directive":  "github.com/vektah/graphql-go/introspection.Directive",
		"__Type":       "github.com/vektah/graphql-go/introspection.Type",
		"__Field":      "github.com/vektah/graphql-go/introspection.Field",
		"__EnumValue":  "github.com/vektah/graphql-go/introspection.EnumValue",
		"__InputValue": "github.com/vektah/graphql-go/introspection.InputValue",
		"__Schema":     "github.com/vektah/graphql-go/introspection.Schema",
		"Query":        "interface{}",
		"Mutation":     "interface{}",
	}
	b, err := ioutil.ReadFile(*typemap)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to open typemap: "+err.Error()+" creating it.")
		return goTypes
	}
	if err = json.Unmarshal(b, &goTypes); err != nil {
		fmt.Fprintln(os.Stderr, "unable to parse typemap: "+err.Error())
		os.Exit(1)
	}

	return goTypes
}
