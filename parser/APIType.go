package parser

// APITypes map of APIType
type APITypes map[string]APIType

// APIType wrap types defined in spec
type APIType struct {
	TypeDeclaration
	ObjectType
}

// UnmarshalYAML implement yaml unmarshaler
func (t *APIType) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	t.AdditionalProperties = true
	if err = unmarshaler(&t.TypeDeclaration); err != nil {
		return
	}
	if err = unmarshaler(&t.ObjectType); err != nil {
		return
	}
	return nil
}
