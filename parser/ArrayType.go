package parser

// ArrayType Array types are declared by using either the array qualifier []
// at the end of a type expression or array as the value of a type facet.
// If you are defining a top-level array type, such as the Emails in the
// examples below, you can declare the following facets in addition to
// those previously described to further restrict the behavior of the array type.
type ArrayType struct {
	// Boolean value that indicates if items in the array MUST be unique.
	UniqueItems bool `yaml:"uniqueItems" json:"uniqueItems,omitempty"`

	// Indicates the type all items in the array are inherited from.
	// Can be a reference to an existing type or an inline type declaration.
	Items Unimplement `yaml:"items" json:"items,omitempty"`

	// Minimum amount of items in array. Value MUST be equal to or greater than 0.
	// Default: 0.
	MinItems int64 `yaml:"minItems" json:"minItems,omitempty"`

	// Maximum amount of items in array. Value MUST be equal to or greater than 0.
	// Default: 2147483647.
	MaxItems int64 `yaml:"maxItems" json:"maxItems,omitdefault" default:"2147483647"`
}

// BeforeUnmarshalYAML implement yaml Initiator
func (t *ArrayType) BeforeUnmarshalYAML() (err error) {
	t.MaxItems = 2147483647
	return
}

// IsEmpty return true if it is empty
func (t *ArrayType) IsEmpty() bool {
	return t.UniqueItems == false &&
		t.Items.IsEmpty() &&
		t.MinItems == 0 &&
		t.MaxItems == 2147483647
}
