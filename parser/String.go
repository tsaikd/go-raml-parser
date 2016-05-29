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
	MaxLength int64 `yaml:"maxLength" json:"maxLength,omitempty"`
}

// UnmarshalYAML implement yaml unmarshaler
func (t *String) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	t.MaxLength = 2147483647

	buf := struct {
		Pattern   *string `yaml:"pattern" json:"pattern,omitempty"`
		MinLength *int64  `yaml:"minLength" json:"minLength,omitempty"`
		MaxLength *int64  `yaml:"maxLength" json:"maxLength,omitempty"`
	}{}
	if err = unmarshaler(&buf); err != nil {
		return
	}

	if buf.Pattern != nil {
		t.Pattern = *buf.Pattern
	}
	if buf.MinLength != nil {
		t.MinLength = *buf.MinLength
	}
	if buf.MaxLength != nil {
		t.MaxLength = *buf.MaxLength
	}

	return
}

// PostProcess for fill some field from RootDocument default config
func (t *String) PostProcess(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}
	return
}

// IsEmpty return true if String is empty
func (t *String) IsEmpty() bool {
	return t.Pattern == "" &&
		t.MinLength == 0 &&
		t.MaxLength == 2147483647
}
