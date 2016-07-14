package parser

// String A JSON string with the following additional facets:
type String struct {
	// Regular expression that this string should match.
	Pattern string `yaml:"pattern" json:"pattern,omitempty"`

	// Minimum length of the string. Value MUST be equal to or greater than 0.
	// Default: 0
	MinLength int64 `yaml:"minLength" json:"minLength,omitempty"`

	// Maximum length of the string. Value MUST be equal to or greater than 0.
	// Default: 2147483647
	MaxLength int64 `yaml:"maxLength" json:"maxLength,omitdefault" default:"2147483647"`
}

// BeforeUnmarshalYAML implement yaml Initiator
func (t *String) BeforeUnmarshalYAML() (err error) {
	t.MaxLength = 2147483647
	return
}

// IsEmpty return true if String is empty
func (t *String) IsEmpty() bool {
	return t.Pattern == "" &&
		t.MinLength == 0 &&
		t.MaxLength == 2147483647
}
