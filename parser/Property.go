package parser

import "strings"

// Properties map of Property
type Properties map[string]*Property

// IsEmpty return true if it is empty
func (t Properties) IsEmpty() bool {
	if t == nil {
		return true
	}
	for _, elem := range t {
		if elem != nil {
			if !elem.IsEmpty() {
				return false
			}
		}
	}
	return true
}

var _ fixRequiredBySyntax = Properties{}

func (t Properties) fixRequiredBySyntax() (err error) {
	if t == nil {
		return
	}
	for name, property := range t {
		if strings.HasSuffix(name, "?") {
			property.Required = false
			trimName := strings.TrimSuffix(name, "?")
			delete(t, name)
			t[trimName] = property
		}
	}
	return
}

// Property of a object type
type Property struct {
	APIType
	PropertyExtra
}

// BeforeUnmarshalYAML implement yaml Initiator
func (t *Property) BeforeUnmarshalYAML() (err error) {
	if err = t.APIType.BeforeUnmarshalYAML(); err != nil {
		return
	}
	if err = t.PropertyExtra.BeforeUnmarshalYAML(); err != nil {
		return
	}
	return
}

// UnmarshalYAML implement yaml unmarshaler
// a Property which MIGHT be a simple string or a map[string]interface{}
func (t *Property) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	if err = unmarshaler(&t.Type); err == nil {
		return
	}
	if !isErrorYAMLIntoString(err) {
		return
	}

	if err = unmarshaler(&t.APIType); err != nil {
		return
	}
	if err = unmarshaler(&t.PropertyExtra); err != nil {
		return
	}

	return
}

// MarshalJSON marshal to json
func (t Property) MarshalJSON() ([]byte, error) {
	return MarshalJSONWithoutEmptyStruct(t)
}

// IsEmpty return true if it is empty
func (t *Property) IsEmpty() bool {
	return t.APIType.IsEmpty() &&
		t.PropertyExtra.IsEmpty()
}

// PropertyExtra contain fields no in APIType
type PropertyExtra struct {
	// Specifies that the property is required or not.
	// Default: true.
	Required bool `yaml:"required" json:"required,omitdefault" default:"true"`
}

// BeforeUnmarshalYAML implement yaml Initiator
func (t *PropertyExtra) BeforeUnmarshalYAML() (err error) {
	t.Required = true
	return
}

// MarshalJSON marshal to json
func (t PropertyExtra) MarshalJSON() ([]byte, error) {
	return MarshalJSONWithoutEmptyStruct(t)
}

// IsEmpty return true if it is empty
func (t *PropertyExtra) IsEmpty() bool {
	return t.Required == true
}
