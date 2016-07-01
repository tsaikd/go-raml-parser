package parser

// Examples The OPTIONAL examples facet can be used to attach multiple examples
// to a type declaration. Its value is a map of key-value pairs, where each key
// represents a unique identifier for an example and the value is a single example.
type Examples map[string]*Example

// IsEmpty return true if it is empty
func (t Examples) IsEmpty() bool {
	for _, elem := range t {
		if elem != nil {
			if !elem.IsEmpty() {
				return false
			}
		}
	}
	return true
}

// SingleExample The OPTIONAL example facet can be used to attach an example of
// a type instance to the type declaration. There are two ways to represent the
// example facet value: as an explicit description of a specific type instance
// and as a map that contains additional facets.
type SingleExample struct {
	// An alternate, human-friendly name for the example. If the example is
	// part of an examples node, the default value is the unique identifier
	// that is defined for this example.
	DisplayName string `yaml:"displayName" json:"displayName,omitempty"`

	// A substantial, human-friendly description for an example. Its value is
	// a string and MAY be formatted using markdown.
	Description string `yaml:"description" json:"description,omitempty"`

	// Annotations to be applied to this API. An annotation is a map having a
	// key that begins with "(" and ends with ")" where the text enclosed in
	// parentheses is the annotation name, and the value is an instance of
	// that annotation.
	Annotations Annotations `yaml:",regexp:\\(.*\\)" json:"annotations,omitempty"`

	// The actual example of a type instance.
	Value Value `yaml:"value" json:"value"`

	// Validates the example against any type declaration (the default), or not.
	// Set this to false avoid validation.
	Strict bool `yaml:"strict" json:"strict,omitempty"`
}

// IsEmpty return true if Example is empty
func (t SingleExample) IsEmpty() bool {
	return t.DisplayName == "" &&
		t.Description == "" &&
		t.Annotations.IsEmpty() &&
		t.Value.IsEmpty()
}

// Example wrap SingleExample for unmarshal YAML
type Example struct {
	SingleExample

	// is include tag set
	includeTag bool
}

// UnmarshalYAMLTag unmarshal an Example which MIGHT be a simple string or a
// map[string]interface{}
func (t *Example) UnmarshalYAMLTag(unmarshaler func(interface{}) error, tag string) (err error) {
	if err = unmarshaler(&t.SingleExample); err == nil && !t.SingleExample.IsEmpty() {
		return
	}

	if tag == "!include" {
		t.includeTag = true
	}

	if err = unmarshaler(&t.Value); err == nil && !t.Value.IsEmpty() {
		return
	}

	return
}

// MarshalJSON marshal to json
func (t Example) MarshalJSON() ([]byte, error) {
	return MarshalJSONWithoutEmptyStruct(t)
}

func generateExampleValue(library Library, apiType APIType, preferArray bool) (value Value, err error) {
	typeName, isArray := GetAPITypeName(apiType)

	if isArray {
		result := []interface{}{}
		for _, example := range apiType.Examples {
			for _, v := range example.Value.Array {
				if !v.IsEmpty() {
					result = append(result, v)
				}
			}
		}
		if len(result) > 0 {
			return NewValue(result)
		}
		if !apiType.Example.Value.IsEmpty() {
			return apiType.Example.Value, nil
		}
	} else if preferArray {
		result := []interface{}{}
		for _, example := range apiType.Examples {
			if !example.Value.IsEmpty() {
				result = append(result, example.Value)
			}
		}
		if len(result) > 0 {
			return NewValue(result)
		}
		if !apiType.Example.Value.IsEmpty() {
			return NewValue([]interface{}{apiType.Example.Value})
		}
	}

	if !apiType.Example.Value.IsEmpty() {
		return apiType.Example.Value, nil
	}
	if !apiType.Examples.IsEmpty() {
		for _, example := range apiType.Examples {
			if !example.Value.IsEmpty() {
				return example.Value, nil
			}
		}
	}

	switch typeName {
	case TypeBoolean, TypeInteger, TypeNumber, TypeString, TypeFile:
		// no value for RAML built-in type
		return Value{}, nil
	case TypeObject:
		valmap := map[string]interface{}{}
		for name, property := range apiType.Properties {
			if valmap[name], err = generateExampleValue(library, property.APIType, false); err != nil {
				return
			}
		}
		if isArray || preferArray {
			return NewValue([]interface{}{valmap})
		}
		return NewValue(valmap)
	default:
		if typ, exist := library.Types[typeName]; exist {
			return generateExampleValue(library, *typ, isArray || preferArray)
		}
	}
	return Value{}, nil
}

func generateExample(library Library, apiType APIType, preferArray bool) (result Example, err error) {
	typeName, isArray := GetAPITypeName(apiType)

	if !apiType.Example.IsEmpty() {
		if !isArray && preferArray {
			example := apiType.Example
			if example.Value, err = generateExampleValue(library, apiType, preferArray); err != nil {
				return Example{}, err
			}
			return example, nil
		}
		return apiType.Example, nil
	}
	if !apiType.Examples.IsEmpty() {
		if !isArray && preferArray {
			example := apiType.Example
			if example.Value, err = generateExampleValue(library, apiType, preferArray); err != nil {
				return Example{}, err
			}
			return example, nil
		}
		for _, example := range apiType.Examples {
			if !example.IsEmpty() {
				return *example, nil
			}
		}
	}

	switch typeName {
	case TypeBoolean, TypeInteger, TypeNumber, TypeString, TypeFile:
		// no value for RAML built-in type
		return Example{}, nil
	default:
		if typ, exist := library.Types[typeName]; exist {
			return generateExample(library, *typ, isArray || preferArray)
		}

		example := Example{}
		if example.Value, err = generateExampleValue(library, apiType, isArray || preferArray); err != nil {
			return Example{}, err
		}
		return example, nil
	}
}

func generateExamples(library Library, apiType APIType, preferArray bool) (result Examples, err error) {
	typeName, isArray := GetAPITypeName(apiType)

	if !apiType.Examples.IsEmpty() {
		if !isArray && preferArray {
			example := apiType.Example
			if example.Value, err = generateExampleValue(library, apiType, preferArray); err != nil {
				return Examples{}, err
			}
			return Examples{
				"autoGenerated": &example,
			}, nil
		}
		return apiType.Examples, nil
	}
	if !apiType.Example.IsEmpty() {
		if !isArray && preferArray {
			example := apiType.Example
			if example.Value, err = generateExampleValue(library, apiType, preferArray); err != nil {
				return Examples{}, err
			}
			return Examples{
				"autoGenerated": &example,
			}, nil
		}
		return Examples{
			"autoGenerated": &apiType.Example,
		}, nil
	}

	switch typeName {
	case TypeBoolean, TypeInteger, TypeNumber, TypeString, TypeFile:
		// no value for RAML built-in type
		return Examples{}, nil
	default:
		if typ, exist := library.Types[typeName]; exist {
			return generateExamples(library, *typ, isArray || preferArray)
		}

		example := Example{}
		if example.Value, err = generateExampleValue(library, apiType, isArray); err != nil {
			return Examples{}, nil
		}
		if !example.IsEmpty() {
			return Examples{
				"autoGenerated": &example,
			}, nil
		}
	}

	return Examples{}, nil
}
