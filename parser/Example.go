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
