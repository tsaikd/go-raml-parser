package parser

// Examples The OPTIONAL examples facet can be used to attach multiple examples
// to a type declaration. Its value is a map of key-value pairs, where each key
// represents a unique identifier for an example and the value is a single example.
type Examples map[string]*Example

// PostProcess for fill some field from RootDocument default config
func (t *Examples) PostProcess(conf PostProcessConfig, apiType APIType) (err error) {
	if t == nil {
		return
	}
	for _, example := range *t {
		if err = example.PostProcess(conf, apiType); err != nil {
			return
		}
	}
	return
}

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
}

// UnmarshalYAML unmarshal an Example which MIGHT be a simple string or a
// map[string]interface{}
func (t *Example) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	if err = unmarshaler(&t.SingleExample); err == nil && !t.SingleExample.IsEmpty() {
		return
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

// PostProcess for fill default example by type if not set
func (t *Example) PostProcess(conf PostProcessConfig, apiType APIType) (err error) {
	if t == nil || t.IsEmpty() {
		return
	}

	typeName, _ := GetAPITypeName(apiType)
	switch typeName {
	case TypeBoolean, TypeInteger, TypeNumber, TypeString:
		// no type check for RAML built-in type
		return
	case TypeObject:
		return CheckValueAPIType(apiType, t.Value)
	default:
		if isInlineAPIType(apiType) {
			// no type check if declared by JSON
			return
		}

		var typ *APIType
		var exist bool
		if typ, exist = conf.Library().Types[typeName]; !exist {
			return ErrorTypeUndefined1.New(nil, apiType.Type)
		}

		return CheckValueAPIType(*typ, t.Value)
	}
}
