package parser

// FileType The ​file​ type can constrain the content to send through forms.
// When this type is used in the context of web forms it SHOULD be represented
// as a valid file upload in JSON format. File content SHOULD be a
// base64-encoded string.
type FileType struct {
	// A list of valid content-type strings for the file.
	// The file type */* MUST be a valid value.
	FileTypes []string `yaml:"fileTypes" json:"fileTypes,omitempty"`

	// Specifies the minimum number of bytes for a parameter value.
	// The value MUST be equal to or greater than 0.
	// Default: 0
	MinLength int64 `yaml:"minLength" json:"minLength,omitempty"`

	// Specifies the maximum number of bytes for a parameter value.
	// The value MUST be equal to or greater than 0.
	// Default: 2147483647
	MaxLength int64 `yaml:"maxLength" json:"maxLength,omitdefault" default:"2147483647"`
}

// BeforeUnmarshalYAML implement yaml Initiator
func (t *FileType) BeforeUnmarshalYAML() (err error) {
	t.MaxLength = 2147483647
	return
}

// IsEmpty return true if it is empty
func (t *FileType) IsEmpty() bool {
	return len(t.FileTypes) < 1 &&
		t.MinLength == 0 &&
		t.MaxLength == 2147483647
}
