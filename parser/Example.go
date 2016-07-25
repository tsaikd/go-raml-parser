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
	Value Value `yaml:"value" json:"value,omitempty"`

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

var _ checkAnnotation = SingleExample{}

func (t SingleExample) checkAnnotation(conf PostProcessConfig) (err error) {
	return t.Annotations.checkAnnotationTargetLocation(TargetLocationExample)
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

func generateExampleValue(library Library, apiType APIType, preferArray bool) (value Value, err error) {
	if apiType.IsArray {
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

	switch apiType.BaseType {
	case TypeObject:
		valmap := map[string]interface{}{}
		for _, property := range apiType.Properties.Slice() {
			if valmap[property.Name], err = generateExampleValue(library, property.APIType, false); err != nil {
				return
			}
		}
		if apiType.IsArray || preferArray {
			return NewValue([]interface{}{valmap})
		}
		return NewValue(valmap)
	default:
		if typ, exist := library.Types[apiType.BaseType]; exist {
			return generateExampleValue(library, *typ, apiType.IsArray || preferArray)
		}
		return Value{}, nil
	}
}

func generateExample(library Library, apiType APIType, preferArray bool) (result Example, err error) {
	if !apiType.Example.IsEmpty() {
		if !apiType.IsArray && preferArray {
			example := apiType.Example
			if example.Value, err = generateExampleValue(library, apiType, preferArray); err != nil {
				return Example{}, err
			}
			return example, nil
		}
		return apiType.Example, nil
	}
	if !apiType.Examples.IsEmpty() {
		if !apiType.IsArray && preferArray {
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

	if typ, exist := library.Types[apiType.BaseType]; exist {
		return generateExample(library, *typ, apiType.IsArray || preferArray)
	}

	example := Example{}
	if example.Value, err = generateExampleValue(library, apiType, apiType.IsArray || preferArray); err != nil {
		return Example{}, err
	}
	return example, nil
}

func generateExamples(library Library, apiType APIType, preferArray bool) (result Examples, err error) {
	if !apiType.Examples.IsEmpty() {
		if !apiType.IsArray && preferArray {
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
		if !apiType.IsArray && preferArray {
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

	if typ, exist := library.Types[apiType.BaseType]; exist {
		return generateExamples(library, *typ, apiType.IsArray || preferArray)
	}

	example := Example{}
	if example.Value, err = generateExampleValue(library, apiType, apiType.IsArray); err != nil {
		return Examples{}, nil
	}
	if !example.IsEmpty() {
		return Examples{
			"autoGenerated": &example,
		}, nil
	}

	return Examples{}, nil
}
