package parser

import (
	"encoding/json"
	"reflect"
	"strconv"
)

// NewValue cast src to value
func NewValue(src interface{}) (Value, error) {
	switch src.(type) {
	case Value:
		return src.(Value), nil
	case []*Value:
		return Value{
			Type:  TypeArray,
			Array: src.([]*Value),
		}, nil
	case []Value:
		value := src.([]Value)
		result := make([]*Value, len(value))
		for i, v := range value {
			result[i] = &v
		}
		return Value{
			Type:  TypeArray,
			Array: result,
		}, nil
	case nil:
		return Value{
			Type: TypeNull,
		}, nil
	case bool:
		return Value{
			Type:    TypeBoolean,
			Boolean: src.(bool),
		}, nil
	case int:
		return Value{
			Type:    TypeInteger,
			Integer: int64(src.(int)),
		}, nil
	case int8:
		return Value{
			Type:    TypeInteger,
			Integer: int64(src.(int8)),
		}, nil
	case int16:
		return Value{
			Type:    TypeInteger,
			Integer: int64(src.(int16)),
		}, nil
	case int32:
		return Value{
			Type:    TypeInteger,
			Integer: int64(src.(int32)),
		}, nil
	case int64:
		return Value{
			Type:    TypeInteger,
			Integer: int64(src.(int64)),
		}, nil
	case float32:
		return Value{
			Type:   TypeNumber,
			Number: float64(src.(float32)),
		}, nil
	case float64:
		return Value{
			Type:   TypeNumber,
			Number: float64(src.(float64)),
		}, nil
	case string:
		return Value{
			Type:   TypeString,
			String: src.(string),
		}, nil
	case []byte:
		return Value{
			Type:   TypeBinary,
			Binary: src.([]byte),
		}, nil
	case []interface{}:
		srcs := src.([]interface{})
		result := make([]*Value, len(srcs))
		for i, elem := range srcs {
			value, err := NewValue(elem)
			if err != nil {
				return value, err
			}
			result[i] = &value
		}
		return Value{
			Type:  TypeArray,
			Array: result,
		}, nil
	case map[string]interface{}:
		result := Value{
			Type: TypeObject,
			Map:  map[string]*Value{},
		}
		srcMap := src.(map[string]interface{})
		for k, v := range srcMap {
			newval, err := NewValue(v)
			if err != nil {
				return Value{}, err
			}
			result.Map[k] = &newval
		}
		return result, nil
	}

	refval := reflect.ValueOf(src)
	switch refval.Kind() {
	case reflect.Ptr:
		return NewValue(refval.Elem().Interface())
	}

	return Value{}, ErrorUnsupportedValueType1.New(nil, src)
}

// Value multiple types value
type Value struct {
	Type    string
	Boolean bool
	Integer int64
	Number  float64
	String  string
	Array   []*Value
	Map     map[string]*Value
	Binary  []byte
}

// MarshalJSON marshal to json
func (t Value) MarshalJSON() ([]byte, error) {
	switch t.Type {
	case TypeBoolean:
		return json.Marshal(t.Boolean)
	case TypeInteger:
		return json.Marshal(t.Integer)
	case TypeNumber:
		return json.Marshal(t.Number)
	case TypeString:
		return json.Marshal(t.String)
	case TypeArray:
		return json.Marshal(t.Array)
	case TypeObject:
		return json.Marshal(t.Map)
	case TypeBinary:
		return json.Marshal(t.Binary)
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

	// integer value should be equal to original string
	var matchstr string
	if err = unmarshaler(&matchstr); err == nil {
		if err = unmarshaler(&t.Integer); err == nil {
			if matchstr == strconv.FormatInt(t.Integer, 10) {
				t.Type = TypeInteger
				return
			}
		} else if !isErrorYAMLIntoInt64(err) {
			return
		}
	}

	if err = unmarshaler(&t.Number); err == nil {
		t.Type = TypeNumber
		return
	}
	if !isErrorYAMLIntoFloat64(err) {
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
		t.Number == 0 &&
		t.String == "" &&
		len(t.Array) < 1 &&
		len(t.Map) < 1 &&
		len(t.Binary) < 1
}
