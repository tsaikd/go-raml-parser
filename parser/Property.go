package parser

import (
	"strings"

	"github.com/tsaikd/KDGoLib/jsonex"
	"github.com/tsaikd/yaml"
)

// Properties contain multiple Property
type Properties struct {
	propertiesSliceData
	mapdata map[string]*Property
}

// UnmarshalYAML implement yaml unmarshaler
func (t *Properties) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	mapdata := map[string]*Property{}
	if err = unmarshaler(mapdata); err != nil {
		return
	}

	order := yaml.MapSlice{}
	if err = unmarshaler(&order); err != nil {
		return
	}

	slicedata := []*Property{}
	for _, item := range order {
		name := item.Key.(string)
		elem := mapdata[name]
		elem.Name = name
		slicedata = append(slicedata, elem)
	}

	t.propertiesSliceData = slicedata
	t.mapdata = mapdata

	return
}

// MarshalJSON marshal to json
func (t Properties) MarshalJSON() ([]byte, error) {
	return jsonex.Marshal(t.mapdata)
}

// IsEmpty return true if it is empty
func (t Properties) IsEmpty() bool {
	for _, elem := range t.propertiesSliceData {
		if elem != nil {
			if !elem.IsEmpty() {
				return false
			}
		}
	}
	return true
}

// Map return properties map
func (t Properties) Map() map[string]*Property {
	return t.mapdata
}

// Slice return properties slice
func (t Properties) Slice() []*Property {
	return t.propertiesSliceData
}

var _ fixRequiredBySyntax = &Properties{}

func (t Properties) fixRequiredBySyntax() (err error) {
	for name, property := range t.mapdata {
		if strings.HasSuffix(name, "?") {
			property.Required = false
			trimName := strings.TrimSuffix(name, "?")
			property.Name = trimName
			delete(t.mapdata, name)
			t.mapdata[trimName] = property
		}
	}
	return
}

type propertiesSliceData []*Property

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

	// Property Name, filled by Properties.UnmarshalYAML()
	Name string `yaml:"-" json:"name,omitempty"`
}

// BeforeUnmarshalYAML implement yaml Initiator
func (t *PropertyExtra) BeforeUnmarshalYAML() (err error) {
	t.Required = true
	return
}

// IsEmpty return true if it is empty
func (t *PropertyExtra) IsEmpty() bool {
	return t.Required == true
}
