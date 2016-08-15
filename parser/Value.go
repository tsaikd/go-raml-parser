package parser

import (
	"net/url"
	"reflect"
	"strconv"

	"github.com/tsaikd/KDGoLib/jsonex"
)

// NewValue cast src to value
func NewValue(src interface{}) (Value, error) {
	var refval reflect.Value

	switch srcval := src.(type) {
	case Value:
		return srcval, nil
	case []*Value:
		return Value{
			Type:  TypeArray,
			Array: srcval,
		}, nil
	case []Value:
		result := make([]*Value, len(srcval))
		for i, v := range srcval {
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
			Boolean: srcval,
		}, nil
	case int:
		return Value{
			Type:    TypeInteger,
			Integer: int64(srcval),
		}, nil
	case int8:
		return Value{
			Type:    TypeInteger,
			Integer: int64(srcval),
		}, nil
	case int16:
		return Value{
			Type:    TypeInteger,
			Integer: int64(srcval),
		}, nil
	case int32:
		return Value{
			Type:    TypeInteger,
			Integer: int64(srcval),
		}, nil
	case int64:
		return Value{
			Type:    TypeInteger,
			Integer: int64(srcval),
		}, nil
	case float32:
		return Value{
			Type:   TypeNumber,
			Number: float64(srcval),
		}, nil
	case float64:
		return Value{
			Type:   TypeNumber,
			Number: float64(srcval),
		}, nil
	case string:
		return Value{
			Type:   TypeString,
			String: srcval,
		}, nil
	case []byte:
		return Value{
			Type:   TypeBinary,
			Binary: srcval,
		}, nil
	case []interface{}:
		result := make([]*Value, len(srcval))
		for i, elem := range srcval {
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
		for k, v := range srcval {
			newval, err := NewValue(v)
			if err != nil {
				return result, err
			}
			result.Map[k] = &newval
		}
		return result, nil
	case url.Values:
		result := Value{
			Type: TypeObject,
			Map:  map[string]*Value{},
		}
		for k, v := range srcval {
			if len(v) == 1 {
				newval, err := NewValue(v[0])
				if err != nil {
					return result, err
				}
				result.Map[k] = &newval
			} else {
				newval, err := NewValue(v)
				if err != nil {
					return result, err
				}
				result.Map[k] = &newval
			}
		}
		return result, nil
	case reflect.Value:
		return NewValue(srcval.Interface())
	default:
		refval = reflect.ValueOf(src)
	}

	switch refval.Kind() {
	case reflect.Ptr:
		elem := refval.Elem()
		if !elem.IsValid() {
			return Value{
				Type: TypeNull,
			}, nil
		}
		return NewValue(elem.Interface())
	case reflect.Slice:
		length := refval.Len()
		result := make([]*Value, length)
		for i := 0; i < length; i++ {
			value, err := NewValue(refval.Index(i))
			if err != nil {
				return Value{}, err
			}
			result[i] = &value
		}
		return Value{
			Type:  TypeArray,
			Array: result,
		}, nil
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
		return jsonex.Marshal(t.Boolean)
	case TypeInteger:
		return jsonex.Marshal(t.Integer)
	case TypeNumber:
		return jsonex.Marshal(t.Number)
	case TypeString:
		return jsonex.Marshal(t.String)
	case TypeArray:
		return jsonex.Marshal(t.Array)
	case TypeObject:
		return jsonex.Marshal(t.Map)
	case TypeBinary:
		return jsonex.Marshal(t.Binary)
	default:
		return jsonex.Marshal(nil)
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

// IsZero return true if Value is empty or contains only type name
func (t Value) IsZero() bool {
	switch t.Type {
	case TypeBoolean:
		return t.Boolean == false
	case TypeInteger:
		return t.Integer == 0
	case TypeNumber:
		return t.Number == 0
	case TypeString:
		return t.String == ""
	case TypeArray:
		if len(t.Array) < 1 {
			return true
		}
		for _, elem := range t.Array {
			if elem != nil && !elem.IsZero() {
				return false
			}
		}
		return true
	case TypeObject:
		if len(t.Map) < 1 {
			return true
		}
		for _, elem := range t.Map {
			if elem != nil && !elem.IsZero() {
				return false
			}
		}
		return true
	case TypeBinary:
		return len(t.Binary) < 1
	default:
		return true
	}
}
