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
	case typeBoolean:
		return json.Marshal(t.Boolean)
	case typeInteger:
		return json.Marshal(t.Integer)
	case typeString:
		return json.Marshal(t.String)
	case typeArray:
		return json.Marshal(t.Array)
	case typeObject:
		return json.Marshal(t.Map)
	default:
		return json.Marshal(nil)
	}
}

// UnmarshalYAML unmarshal an Example which MIGHT be a simple string or a
// map[string]interface{}
func (t *Value) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	if err = unmarshaler(&t.Boolean); err == nil {
		t.Type = typeBoolean
		return
	}
	if !isErrorYAMLIntoBool(err) {
		return
	}

	if err = unmarshaler(&t.Integer); err == nil {
		t.Type = typeInteger
		return
	}
	if !isErrorYAMLIntoInt64(err) {
		return
	}

	if err = unmarshaler(&t.String); err == nil {
		t.Type = typeString
		return
	}
	if !isErrorYAMLIntoString(err) {
		return
	}

	if err = unmarshaler(&t.Array); err == nil {
		t.Type = typeArray
		return
	}

	if err = unmarshaler(&t.Map); err != nil {
		return
	}
	t.Type = typeObject
	return
}

// PostProcess for fill some field from RootDocument default config
func (t *Value) PostProcess(conf PostProcessConfig) (err error) {
	return
}

// IsEmpty return true if it is empty
func (t Value) IsEmpty() bool {
	return t.Boolean == false &&
		t.Integer == 0 &&
		t.String == "" &&
		len(t.Array) < 1 &&
		len(t.Map) < 1
}
