# External Package Guide: github.com/santhosh-tekuri/jsonschema/v6

**Date**: 2026-02-19
**Docs**: https://pkg.go.dev/github.com/santhosh-tekuri/jsonschema/v6
**Task**: 5-4 (Go Schema Validation Package)

## Overview

Pure Go JSON Schema validator supporting drafts 4, 6, 7, 2019-09, and 2020-12. No CGO dependencies.

## Key API

### Compile a schema from in-memory data

```go
import "github.com/santhosh-tekuri/jsonschema/v6"

// 1. Unmarshal schema JSON into an any value
schemaDoc, err := jsonschema.UnmarshalJSON(strings.NewReader(schemaJSON))

// 2. Create compiler, add schema as a resource
compiler := jsonschema.NewCompiler()
compiler.AddResource("schema.json", schemaDoc)

// 3. Compile
schema, err := compiler.Compile("schema.json")
```

### Validate data

```go
// Unmarshal the input data
data, err := jsonschema.UnmarshalJSON(strings.NewReader(inputJSON))

// Validate — returns nil on success, *ValidationError on failure
err = schema.Validate(data)
```

### Error handling

```go
if err != nil {
    validErr := err.(*jsonschema.ValidationError)

    // Basic output gives a flat list of errors with instance locations
    output := validErr.BasicOutput()
    for _, e := range output.Errors {
        fmt.Printf("  %s: %s\n", e.InstanceLocation, e.Error.Kind)
    }

    // Or use DetailedOutput() for nested structure
    detailed := validErr.DetailedOutput()
}
```

### Output types

- `FlagOutput()` — just valid/invalid boolean
- `BasicOutput()` — flat list of `OutputUnit` with `InstanceLocation` and `Error`
- `DetailedOutput()` — nested tree of `OutputUnit`

### OutputUnit fields

```go
type OutputUnit struct {
    Valid                   bool
    KeywordLocation         string  // schema path
    AbsoluteKeywordLocation string
    InstanceLocation        string  // JSON pointer to failing value (e.g., "/modules/0/lessons/2")
    Error                   *OutputError
    Errors                  []*OutputUnit
}
```

### Enabling format validation

```go
compiler := jsonschema.NewCompiler()
compiler.AssertFormat() // makes "format" keyword actually validate (e.g., date-time, uri)
```

## Usage pattern for Apollo

```go
//go:embed schemas/curriculum.json
var schemaBytes []byte

var (
    compiledSchema *jsonschema.Schema
    compileOnce    sync.Once
    compileErr     error
)

func getSchema() (*jsonschema.Schema, error) {
    compileOnce.Do(func() {
        doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(schemaBytes))
        if err != nil {
            compileErr = fmt.Errorf("unmarshal schema: %w", err)
            return
        }
        c := jsonschema.NewCompiler()
        if err := c.AddResource("curriculum.json", doc); err != nil {
            compileErr = fmt.Errorf("add resource: %w", err)
            return
        }
        compiledSchema, compileErr = c.Compile("curriculum.json")
    })
    return compiledSchema, compileErr
}

func Validate(jsonData []byte) error {
    sch, err := getSchema()
    if err != nil {
        return err
    }
    data, err := jsonschema.UnmarshalJSON(bytes.NewReader(jsonData))
    if err != nil {
        return fmt.Errorf("invalid JSON: %w", err)
    }
    return sch.Validate(data)
}
```
