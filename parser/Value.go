package parser

import "encoding/json"

// Value multiple types value
type Value struct {
	Type    string
	Boolean bool
	Integer int64
	String  string
	Array   []*Value
	Map     map[string]*Value
}

// MarshalJSON marshal to json
func (t Value) MarshalJSON() ([]byte, error) {
	switch t.Type {
	case TypeBoolean:
		return json.Marshal(t.Boolean)
	case TypeInteger:
		return json.Marshal(t.Integer)
	case TypeString:
		return json.Marshal(t.String)
	case TypeArray:
		return json.Marshal(t.Array)
	case TypeObject:
		return json.Marshal(t.Map)
	default:
		return json.Marshal(nil)
	}
}

// UnmarshalYAML unmarshal from YAML
func (t *Value) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	if err = unmarshaler(&t.Boolean); err == nil {
		t.Type = TypeBoolean
		return
	}
	if !isErrorYAMLIntoBool(err) {
		return
	}

	if err = unmarshaler(&t.Integer); err == nil {
		t.Type = TypeInteger
		return
	}
	if !isErrorYAMLIntoInt64(err) {
		return
	}

	if err = unmarshaler(&t.String); err == nil {
		t.Type = TypeString
		return
	}
	if !isErrorYAMLIntoString(err) {
		return
	}

	if err = unmarshaler(&t.Array); err == nil {
		t.Type = TypeArray
		return
	}

	if err = unmarshaler(&t.Map); err != nil {
		return
	}
	t.Type = TypeObject
	return
}

// PostProcess for fill some field from RootDocument default config
func (t *Value) PostProcess(conf PostProcessConfig) (err error) {
	return
}

// IsEmpty return true if it is empty
func (t Value) IsEmpty() bool {
	return t.Type == "" &&
		t.Boolean == false &&
		t.Integer == 0 &&
		t.String == "" &&
		len(t.Array) < 1 &&
		len(t.Map) < 1
}
