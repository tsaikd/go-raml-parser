package parser

// String A JSON string with the following additional facets:
type String struct {
	// Regular expression that this string should match.
	Pattern string `yaml:"pattern" json:"pattern,omitempty"`

	// Minimum length of the string. Value MUST be equal to or greater than 0.
	// Default: 0
	MinLength Unimplement `yaml:"minLength" json:"minLength,omitempty"`

	// Maximum length of the string. Value MUST be equal to or greater than 0.
	// Default: 2147483647
	MaxLength Unimplement `yaml:"maxLength" json:"maxLength,omitempty"`
}

// PostProcess for fill some field from RootDocument default config
func (t *String) PostProcess(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}
	if err = t.MinLength.PostProcess(conf); err != nil {
		return
	}
	if err = t.MaxLength.PostProcess(conf); err != nil {
		return
	}
	return
}

// IsEmpty return true if String is empty
func (t *String) IsEmpty() bool {
	return t.Pattern == "" &&
		t.MinLength.IsEmpty() &&
		t.MaxLength.IsEmpty()
}
