package parser

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/tsaikd/go-raml-parser/parser/parserConfig"
)

// APITypes map of APIType
type APITypes map[string]*APIType

// IsEmpty return true if it is empty
func (t APITypes) IsEmpty() bool {
	for _, elem := range t {
		if elem != nil {
			if !elem.IsEmpty() {
				return false
			}
		}
	}
	return true
}

// APIType wrap types defined in spec
type APIType struct {
	TypeDeclaration
	ObjectType
	ScalarType
	String
	ArrayType
	FileType
}

// BeforeUnmarshalYAML implement yaml Initiator
func (t *APIType) BeforeUnmarshalYAML() (err error) {
	if err = t.ObjectType.BeforeUnmarshalYAML(); err != nil {
		return
	}
	if err = t.String.BeforeUnmarshalYAML(); err != nil {
		return
	}
	if err = t.ArrayType.BeforeUnmarshalYAML(); err != nil {
		return
	}
	if err = t.FileType.BeforeUnmarshalYAML(); err != nil {
		return
	}
	return
}

// UnmarshalYAML implement yaml unmarshaler
func (t *APIType) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	if err = unmarshaler(&t.TypeDeclaration); err != nil {
		return
	}
	if strings.HasSuffix(t.TypeDeclaration.Type, "[]") {
		if err = unmarshaler(&t.ArrayType); err != nil {
			return
		}
	}
	if err = unmarshaler(&t.ObjectType); err != nil {
		return
	}
	if !t.ObjectType.IsEmpty() {
		t.Type = TypeObject
		return nil
	}
	if err = unmarshaler(&t.ScalarType); err != nil {
		return
	}
	if err = unmarshaler(&t.String); err != nil {
		return
	}
	if err = unmarshaler(&t.FileType); err != nil {
		return
	}
	return nil
}

// MarshalJSON marshal to json
func (t APIType) MarshalJSON() ([]byte, error) {
	return MarshalJSONWithoutEmptyStruct(t)
}

// IsEmpty return true if it is empty
func (t APIType) IsEmpty() bool {
	return t.TypeDeclaration.IsEmpty() &&
		t.ObjectType.IsEmpty() &&
		t.ScalarType.IsEmpty() &&
		t.String.IsEmpty() &&
		t.ArrayType.IsEmpty() &&
		t.FileType.IsEmpty()
}

var _ fillProperties = &APIType{}

func (t *APIType) fillProperties(library Library) (err error) {
	if t == nil {
		return
	}
	if !t.Properties.IsEmpty() {
		return
	}

	// fill Properties if possible
	typeName, _ := GetAPITypeName(*t)
	switch typeName {
	case "", TypeBoolean, TypeInteger, TypeNumber, TypeString, TypeObject, TypeFile:
		// no more action for RAML built-in type
		return
	default:
		if isInlineAPIType(*t) {
			// no more action if declared by JSON
			return
		}

		var typ *APIType
		var exist bool
		if typ, exist = library.Types[typeName]; !exist {
			return ErrorTypeUndefined1.New(nil, t.Type)
		}
		t.Properties = typ.Properties
	}

	return
}

var _ fillExample = &APIType{}

func (t *APIType) fillExample(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}

	if t.Example.IsEmpty() {
		if t.Example, err = generateExample(conf.Library(), *t, false); err != nil {
			return
		}
	}
	if t.Examples.IsEmpty() {
		if t.Examples, err = generateExamples(conf.Library(), *t, false); err != nil {
			return
		}
	}

	if err = fillExampleAPIType(&t.Example, conf, *t); err != nil {
		return
	}
	for _, example := range t.Examples {
		if err = fillExampleAPIType(example, conf, *t); err != nil {
			return
		}
	}

	return
}

func fillExampleAPIType(example *Example, conf PostProcessConfig, apiType APIType) (err error) {
	if example == nil || example.IsEmpty() {
		return
	}

	if example.includeTag && TypeString == example.Value.Type {
		fpath := filepath.Join(conf.RootDocument().WorkingDirectory, example.Value.String)
		var fdata []byte
		if fdata, err = ioutil.ReadFile(fpath); err != nil {
			return
		}
		switch apiType.Type {
		case TypeFile:
			if example.Value, err = NewValue(fdata); err != nil {
				return
			}
		default:
			return ErrorUnsupportedIncludeType1.New(nil, apiType.Type)
		}
	}

	if err = fillValueFromAPIType(&example.Value, conf.Library(), apiType); err != nil {
		return
	}

	return
}

func fillValueFromAPIType(value *Value, library Library, apiType APIType) (err error) {
	if value == nil {
		return nil
	}

	if value.IsEmpty() {
		if *value, err = generateExampleValue(library, apiType, false); err != nil {
			return
		}
	}

	// not support fill value from inline APIType
	if isInlineAPIType(apiType) {
		return
	}

	for name, v := range value.Map {
		if v == nil {
			v = &Value{}
			value.Map[name] = v
		}
		property := apiType.Properties[name]
		if err = fillValueFromAPIType(v, library, property.APIType); err != nil {
			return
		}
	}
	return nil
}

var _ checkExample = &APIType{}

func (t *APIType) checkExample(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}

	options := []CheckValueOption{}
	confOptions, err := conf.Parser().Get(parserConfig.CheckValueOptions)
	if err == nil {
		if opts, ok := confOptions.([]CheckValueOption); ok {
			options = opts
		}
	}

	typeName, _ := GetAPITypeName(*t)
	switch typeName {
	case TypeBoolean, TypeInteger, TypeNumber, TypeString, TypeFile:
		// no type check for RAML built-in type
		return
	case TypeObject:
		if err = CheckValueAPIType(*t, t.Example.Value, options...); err != nil {
			return
		}
		for _, example := range t.Examples {
			if err = CheckValueAPIType(*t, example.Value, options...); err != nil {
				return
			}
		}
	}

	return
}
