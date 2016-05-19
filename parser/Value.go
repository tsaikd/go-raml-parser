package parser

import (
	"encoding/json"
	"strings"
)

// Value multiple types value
type Value struct {
	Type   string
	String string
	Map    map[string]interface{}
}

// MarshalJSON marshal to json
func (t Value) MarshalJSON() ([]byte, error) {
	switch t.Type {
	case typeString:
		return json.Marshal(t.String)
	case typeObject:
		return json.Marshal(t.Map)
	default:
		return json.Marshal(nil)
	}
}

// UnmarshalYAML unmarshal an Example which MIGHT be a simple string or a
// map[string]interface{}
func (t *Value) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	err = unmarshaler(&t.String)
	if err == nil {
		t.Type = typeString
		return
	}
	if !strings.Contains(err.Error(), "into string") {
		return
	}

	if err = unmarshaler(&t.Map); err != nil {
		return
	}
	t.Type = typeObject
	return
}

// IsEmpty return true if Example is empty
func (t Value) IsEmpty() bool {
	if t.String != "" {
		return false
	}
	if len(t.Map) > 0 {
		return false
	}
	return true
}
