package parser

// ScalarType RAML defines a set of built-in scalar types, each of which has
// a predefined set of restrictions. All types, except the file type,
// can have an additional enum facet.
type ScalarType struct {
	// Enumeration of possible values for this built-in scalar type.
	// The value is an array containing representations of possible values,
	// or a single value if there is only one possible value.
	Enum []Value `yaml:"enum" json:"enum,omitempty"`
}

// MarshalJSON marshal to json
func (t ScalarType) MarshalJSON() ([]byte, error) {
	return MarshalJSONWithoutEmptyStruct(t)
}

// PostProcess for fill some field from RootDocument default config
func (t *ScalarType) PostProcess(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}
	return
}

// IsEmpty return true if it is empty
func (t *ScalarType) IsEmpty() bool {
	return len(t.Enum) < 1
}
