package parser

import (
	"encoding/json"
	"strings"
)

// Example for mock
type Example struct {
	Type   string
	String string
	Map    map[string]interface{}
}

// MarshalJSON marshal to json
func (t Example) MarshalJSON() ([]byte, error) {
	switch t.Type {
	case "string":
		return json.Marshal(t.String)
	case "object":
		return json.Marshal(t.Map)
	default:
		return json.Marshal(nil)
	}
}

// UnmarshalYAML unmarshal an Example which MIGHT be a simple string or a
// map[string]interface{}
func (t *Example) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	err = unmarshaler(&t.String)
	if err == nil {
		t.Type = "string"
		return
	}
	if !strings.Contains(err.Error(), "into string") {
		return
	}

	if err = unmarshaler(&t.Map); err != nil {
		return
	}
	t.Type = "object"
	return
}

// IsEmpty return true if Example is empty
func (t Example) IsEmpty() bool {
	if t.String != "" {
		return false
	}
	if len(t.Map) > 0 {
		return false
	}
	return true
}

// PostProcess for fill default example by type if not set
func (t *Example) PostProcess(rootdoc RootDocument, exampleType string) (err error) {
	if !t.IsEmpty() {
		return
	}

	if rootType, exist := rootdoc.Types[exampleType]; exist {
		*t = rootType.Example
	}

	return
}
