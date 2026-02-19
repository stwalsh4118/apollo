package schema

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

const (
	curriculumSchemaFile  = "curriculum.json"
	poolSummarySchemaFile = "knowledge_pool_summary.json"
)

// schemaCache holds a lazily-compiled schema.
type schemaCache struct {
	once   sync.Once
	schema *jsonschema.Schema
	err    error
}

var (
	curriculumCache  schemaCache
	poolSummaryCache schemaCache
)

func (sc *schemaCache) get(fileName string) (*jsonschema.Schema, error) {
	sc.once.Do(func() {
		raw, err := schemaFS.ReadFile(fileName)
		if err != nil {
			sc.err = fmt.Errorf("read embedded schema %s: %w", fileName, err)
			return
		}

		doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
		if err != nil {
			sc.err = fmt.Errorf("unmarshal schema %s: %w", fileName, err)
			return
		}

		c := jsonschema.NewCompiler()
		if err := c.AddResource(fileName, doc); err != nil {
			sc.err = fmt.Errorf("add schema resource %s: %w", fileName, err)
			return
		}

		sc.schema, sc.err = c.Compile(fileName)
	})

	return sc.schema, sc.err
}

// Validate checks jsonData against the embedded curriculum JSON schema.
// It returns nil when the data is valid. On failure it returns an error
// describing the first validation issue with a field path.
func Validate(jsonData []byte) error {
	return validate(&curriculumCache, curriculumSchemaFile, jsonData)
}

// ValidatePoolSummary checks jsonData against the knowledge pool summary schema.
func ValidatePoolSummary(jsonData []byte) error {
	return validate(&poolSummaryCache, poolSummarySchemaFile, jsonData)
}

func validate(cache *schemaCache, fileName string, jsonData []byte) error {
	sch, err := cache.get(fileName)
	if err != nil {
		return fmt.Errorf("load schema: %w", err)
	}

	data, err := jsonschema.UnmarshalJSON(bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("invalid JSON input: %w", err)
	}

	if err := sch.Validate(data); err != nil {
		vErr, ok := err.(*jsonschema.ValidationError)
		if !ok {
			return fmt.Errorf("validation failed: %w", err)
		}
		return formatValidationError(vErr)
	}

	return nil
}

// formatValidationError walks the BasicOutput to produce a human-readable
// error with field paths.
func formatValidationError(vErr *jsonschema.ValidationError) error {
	output := vErr.BasicOutput()
	if output == nil {
		return fmt.Errorf("validation failed: %s", vErr.Error())
	}

	// Collect leaf errors (those with an Error field and no sub-errors).
	var msgs []string
	collectErrors(output, &msgs)

	if len(msgs) == 0 {
		return fmt.Errorf("validation failed: %s", vErr.Error())
	}

	// Return first error for clarity; full list available via DetailedOutput.
	if len(msgs) == 1 {
		return fmt.Errorf("schema validation failed: %s", msgs[0])
	}
	return fmt.Errorf("schema validation failed (%d errors, first): %s", len(msgs), msgs[0])
}

func collectErrors(unit *jsonschema.OutputUnit, msgs *[]string) {
	if unit.Error != nil && unit.InstanceLocation != "" {
		*msgs = append(*msgs, fmt.Sprintf("%s: %s", unit.InstanceLocation, unit.Error.Kind))
	}
	for i := range unit.Errors {
		collectErrors(&unit.Errors[i], msgs)
	}
}
