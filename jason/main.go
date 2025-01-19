package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/dave/jennifer/jen"
)

// https://json-schema.org/specification
// https://json-schema.org/draft/2020-12/json-schema-core

type SchemaEntry struct {
	Schema  Schema
	BaseURL string
}

type SchemaRegistry struct {
	Enties map[string]*SchemaEntry
}

func (cache *SchemaRegistry) Lookup(reference, base string) (*SchemaEntry, error) {
	ref, err := url.Parse(reference)
	if err != nil {
		return nil, err
	}

	if base != "" {
		b, err := url.Parse(base)
		if err != nil {
			return nil, err
		}
		ref = b.ResolveReference(ref)
	}

	resolvedReference := ref.String()

	if entry, found := cache.Enties[resolvedReference]; found {
		return entry, nil
	}

	return nil, nil
}

func Load(path string) (*SchemaEntry, error) {
	entry := SchemaEntry{}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &entry.Schema)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func GenerateType(schema *Schema) *jen.Statement {
	var types []string
	switch t := schema.Type.(type) {
	case string:
		types = []string{t}
	case []any:
		for _, x := range t {
			switch x := x.(type) {
			case string:
				types = append(types, string(x))
			}
		}
	default:
		return jen.Any()
	}

	decl := jen.Add()
	if len(types) == 2 {
		if types[1] == "null" {
			return decl.Op("*")
		}
	}

	switch types[0] {
	case "string":
		decl.String()
	case "integer":
		decl.Int64()
	case "number":
		decl.Float64()
	case "boolean":
		decl.Bool()
	case "array":
		decl.Index().Add(GenerateType(schema.Items))
	case "object":
		decl.StructFunc(func(g *jen.Group) {
			for field, def := range schema.Properties {
				g.Id(field).Add(GenerateType(&def))
			}
		})
	case "null":
		decl.Any()
	}

	return decl
}

func Generate(entry *SchemaEntry) (string, error) {
	f := jen.NewFile("xxx")
	f.Type().Id("Schema").Add(GenerateType(&entry.Schema))
	buf := bytes.Buffer{}
	err := f.Render(&buf)
	return buf.String(), err
}

func main() {
	s, err := Load("json-schema-spec/schema.json")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%+v\n\n", s.Schema)
	code, err := Generate(s)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", code)
}
