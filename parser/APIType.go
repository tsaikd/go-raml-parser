package parser

// APITypes map of APIType
type APITypes map[string]*APIType

// PostProcess for fill some field from RootDocument default config
func (t *APITypes) PostProcess(rootdoc RootDocument) (err error) {
	for _, apitype := range *t {
		if err = apitype.PostProcess(rootdoc); err != nil {
			return
		}
	}
	return
}

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

// PostProcess for fill some field from RootDocument default config
func (t *APIType) PostProcess(rootdoc RootDocument) (err error) {
	if err = t.TypeDeclaration.PostProcess(rootdoc); err != nil {
		return
	}
	if err = t.ObjectType.PostProcess(rootdoc); err != nil {
		return
	}
	return
}
