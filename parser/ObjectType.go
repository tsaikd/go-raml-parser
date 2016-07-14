package parser

// ObjectType All types that have the built-in object type at the root of
// their inheritance tree can use the following facets in their type
// declarations:
type ObjectType struct {
	// The properties that instances of this type can or must have.
	Properties Properties `yaml:"properties" json:"properties,omitempty"`

	// The minimum number of properties allowed for instances of this type.
	MinProperties Unimplement `yaml:"minProperties" json:"minProperties,omitempty"`

	// The maximum number of properties allowed for instances of this type.
	MaxProperties Unimplement `yaml:"maxProperties" json:"maxProperties,omitempty"`

	// A Boolean that indicates if an object instance has additional properties.
	// Default: true
	AdditionalProperties bool `yaml:"additionalProperties" json:"additionalProperties,omitdefault" default:"true"`

	// Determines the concrete type of an individual object at runtime when,
	// for example, payloads contain ambiguous types due to unions or
	// inheritance. The value must match the name of one of the declared
	// properties of a type. Unsupported practices are inline type declarations
	// and using discriminator with non-scalar properties.
	Discriminator Unimplement `yaml:"discriminator" json:"discriminator,omitempty"`

	// Identifies the declaring type. Requires including a discriminator facet
	// in the type declaration. A valid value is an actual value that might
	// identify the type of an individual object and is unique in the
	// hierarchy of the type. Inline type declarations are not supported.
	// Default: The name of the type
	DiscriminatorValue Unimplement `yaml:"discriminatorValue" json:"discriminatorValue,omitempty"`
}

// BeforeUnmarshalYAML implement yaml Initiator
func (t *ObjectType) BeforeUnmarshalYAML() (err error) {
	t.AdditionalProperties = true
	return
}

// IsEmpty return true if it is empty
func (t ObjectType) IsEmpty() bool {
	return t.Properties.IsEmpty() &&
		t.MinProperties.IsEmpty() &&
		t.MaxProperties.IsEmpty() &&
		t.AdditionalProperties == true &&
		t.Discriminator.IsEmpty() &&
		t.DiscriminatorValue.IsEmpty()
}
