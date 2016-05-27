package parser

import "strings"

// Properties map of Property
type Properties map[string]*Property

// PostProcess for fill some field from RootDocument default config
func (t *Properties) PostProcess(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}
	for name, property := range *t {
		if err = property.PostProcess(conf); err != nil {
			return
		}
		if strings.HasSuffix(name, "?") {
			property.Required = false
			trimName := strings.TrimRight(name, "?")
			delete(*t, name)
			(*t)[trimName] = property
		}
	}
	return
}

// Property of a object type
type Property struct {
	TypeDeclaration
	PropertyExtra
}

// UnmarshalYAML implement yaml unmarshaler
// a Property which MIGHT be a simple string or a map[string]interface{}
func (t *Property) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	t.Required = true

	if err = unmarshaler(&t.Type); err == nil {
		return
	}
	if !isErrorYAMLIntoString(err) {
		return
	}

	if err = unmarshaler(&t.TypeDeclaration); err != nil {
		return
	}
	if err = unmarshaler(&t.PropertyExtra); err != nil {
		return
	}

	return
}

// PostProcess for fill some field from RootDocument default config
func (t *Property) PostProcess(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}
	if err = t.TypeDeclaration.PostProcess(conf); err != nil {
		return
	}
	if err = t.PropertyExtra.PostProcess(conf); err != nil {
		return
	}
	return
}

// PropertyExtra contain fields no in TypeDeclaration
type PropertyExtra struct {
	// Specifies that the property is required or not.
	// Default: true.
	Required bool `yaml:"required" json:"required"`
}

// PostProcess for fill some field from RootDocument default config
func (t *PropertyExtra) PostProcess(conf PostProcessConfig) (err error) {
	return
}
