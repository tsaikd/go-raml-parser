package parser

import "strings"

// Properties map of Property
type Properties map[string]*Property

// Property of a object type
type Property struct {
	TypeDeclaration

	// Specifies that the property is required or not.
	// Default: true.
	Required bool `yaml:"required" json:"required"`
}

// UnmarshalYAML implement yaml unmarshaler
// a Property which MIGHT be a simple string or a map[string]interface{}
func (t *Property) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	defer func() {
		if err == nil {
			t.Required = !strings.HasSuffix(t.Type, "?")
			t.Type = strings.TrimRight(t.Type, "?")
		}
	}()

	err = unmarshaler(&t.Type)
	if err == nil {
		return
	}
	if !strings.Contains(err.Error(), "into string") {
		return
	}

	if err = unmarshaler(&t.TypeDeclaration); err != nil {
		return
	}

	return
}

// PostProcess for fill some field from RootDocument default config
func (t *Property) PostProcess(rootdoc RootDocument) (err error) {
	if err = t.TypeDeclaration.PostProcess(rootdoc); err != nil {
		return
	}
	return
}
