package parser

// TypeDeclaration A type declaration references another type, or wraps or extends
// another type by adding functional facets (e.g. properties) or non-functional
// facets (e.g. a description), or is a type expression that uses other types.
// Here are the facets that all type declarations can have; certain type
// declarations might have other facets:
type TypeDeclaration struct {
	// A default value for a type. When an API request is completely missing
	// the instance of a type, for example when a query parameter described by
	// a type is entirely missing from the request, then the API must act as
	// if the API client had sent an instance of that type with the instance
	// value being the value in the default facet. Similarly, when the API
	// response is completely missing the instance of a type, the client must
	// act as if the API server had returned an instance of that type with
	// the instance value being the value in the default facet. A special case
	// is made for URI parameters: for these, the client MUST substitute the
	// value in the default facet if no instance of the URI parameter was given.
	Default Unimplement `yaml:"default" json:"default,omitempty"`

	// An alias for the equivalent "type" facet for compatibility with RAML
	// 0.8. Deprecated - API definitions should use the "type" facet because
	// the "schema" alias for that facet name might be removed in a future RAML
	// version. The "type" facet supports XML and JSON schemas.
	Schema Unimplement `yaml:"schema" json:"schema,omitempty"`

	// A base type which the current type extends or just wraps. The value of
	// a type node MUST be either a) the name of a user-defined type or b) the
	// name of a built-in RAML data type (object, array, or one of the scalar
	// types) or c) an inline type declaration.
	Type string `yaml:"type" json:"type,omitempty"`

	// An example of an instance of this type that can be used, for example,
	// by documentation generators to generate sample values for an object of
	// this type. The "example" facet MUST not be available when the "examples"
	// facet is already defined. See section Examples for more information.
	Example Example `yaml:"example" json:"example,omitempty"`

	// Examples of instances of this type. This can be used, for example, by
	// documentation generators to generate sample values for an object of this
	// type. The "examples" facet MUST not be available when the "example"
	// facet is already defined. See section Examples for more information.
	Examples Examples `yaml:"examples" json:"examples,omitempty"`

	// An alternate, human-friendly name for the type
	DisplayName string `yaml:"displayName" json:"displayName,omitempty"`

	// A substantial, human-friendly description of the type. Its value is a
	// string and MAY be formatted using markdown.
	Description string `yaml:"description" json:"description,omitempty"`

	// Annotations to be applied to this API. An annotation is a map having a
	// key that begins with "(" and ends with ")" where the text enclosed in
	// parentheses is the annotation name, and the value is an instance of
	// that annotation.
	Annotations Annotations `yaml:",regexp:\\(.*\\)" json:"annotations,omitempty"`

	// A map of additional, user-defined restrictions that will be inherited
	// and applied by any extending subtype. See section User-defined Facets
	// for more information.
	Facets Unimplement `yaml:"facets" json:"facets,omitempty"`

	// The capability to configure XML serialization of this type instance.
	XML Unimplement `yaml:"xml" json:"xml,omitempty"`
}

// MarshalJSON marshal to json
func (t TypeDeclaration) MarshalJSON() ([]byte, error) {
	return MarshalJSONWithoutEmptyStruct(t)
}

func generateExample(conf PostProcessConfig, apiType APIType) Example {
	if !apiType.Example.IsEmpty() {
		return apiType.Example
	}

	typeName, _ := GetAPITypeName(apiType)
	if typ, exist := conf.Library().Types[typeName]; exist {
		if !typ.Example.IsEmpty() {
			return typ.Example
		} else if !typ.Examples.IsEmpty() {
			for _, example := range typ.Examples {
				if !example.IsEmpty() {
					return *example
				}
			}
		}
	}

	switch typeName {
	case TypeObject:
		example := Example{
			SingleExample: SingleExample{
				Value: Value{
					Map: map[string]*Value{},
				},
			},
		}
		for name, property := range apiType.Properties {
			_, propIsArray := GetAPITypeName(property.APIType)
			if propIsArray {
				propexs := generateExamples(conf, property.APIType)
				if !propexs.IsEmpty() {
					values := Value{}
					for _, propex := range propexs {
						values.Array = append(values.Array, &propex.Value)
					}
					if !values.IsEmpty() {
						values.Type = TypeArray
					}
					example.Value.Map[name] = &values
				}
			} else {
				propex := generateExample(conf, property.APIType)
				if !propex.IsEmpty() {
					example.Value.Map[name] = &propex.Value
				}
			}
		}
		if !example.IsEmpty() {
			example.Value.Type = TypeObject
		}
		return example
	}

	return Example{}
}

func generateExamples(conf PostProcessConfig, apiType APIType) Examples {
	if !apiType.Examples.IsEmpty() {
		return apiType.Examples
	}

	typeName, _ := GetAPITypeName(apiType)
	if typ, exist := conf.Library().Types[typeName]; exist {
		if !typ.Examples.IsEmpty() {
			return typ.Examples
		} else if !typ.Example.IsEmpty() {
			example := typ.Example
			return Examples{
				"autoGenerated": &example,
			}
		}
	}

	switch typeName {
	case TypeObject:
		example := generateExample(conf, apiType)
		return Examples{
			"autoGenerated": &example,
		}
	}

	return Examples{}
}

// PostProcess for fill default example by type if not set
func (t *TypeDeclaration) PostProcess(conf PostProcessConfig, apiType APIType) (err error) {
	if t == nil {
		return
	}

	typeName, isArrayType := GetAPITypeName(apiType)
	if isArrayType {
		if t.Examples.IsEmpty() {
			if typ, exist := conf.Library().Types[typeName]; exist {
				if !typ.Examples.IsEmpty() {
					t.Examples = typ.Examples
				} else if !typ.Example.IsEmpty() {
					example := typ.Example
					t.Examples = Examples{
						"autoGenerated": &example,
					}
				}
			}
		}
	} else {
		if t.Example.IsEmpty() {
			t.Example = generateExample(conf, apiType)
		}
	}

	if err = t.Default.PostProcess(conf); err != nil {
		return
	}
	if err = t.Schema.PostProcess(conf); err != nil {
		return
	}
	if err = t.Example.PostProcess(conf, apiType); err != nil {
		return
	}
	if err = t.Examples.PostProcess(conf, apiType); err != nil {
		return
	}
	if err = t.Annotations.PostProcess(conf); err != nil {
		return
	}
	if err = t.Facets.PostProcess(conf); err != nil {
		return
	}
	if err = t.XML.PostProcess(conf); err != nil {
		return
	}
	return
}

// IsEmpty return true if it is empty
func (t *TypeDeclaration) IsEmpty() bool {
	if t == nil {
		return true
	}
	return t.Default.IsEmpty() &&
		t.Schema.IsEmpty() &&
		t.Type == "" &&
		t.Example.IsEmpty() &&
		t.Examples.IsEmpty() &&
		t.DisplayName == "" &&
		t.Description == "" &&
		t.Annotations.IsEmpty() &&
		t.Facets.IsEmpty() &&
		t.XML.IsEmpty()
}
