package parser

import "strings"

// APITypes map of APIType
type APITypes map[string]*APIType

// PostProcess for fill some field from RootDocument default config
func (t *APITypes) PostProcess(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}
	for _, apitype := range *t {
		if err = apitype.PostProcess(conf); err != nil {
			return
		}
	}
	return
}

// IsEmpty return true if it is empty
func (t APITypes) IsEmpty() bool {
	for _, elem := range t {
		if elem != nil {
			if !elem.IsEmpty() {
				return false
			}
		}
	}
	return true
}

// APIType wrap types defined in spec
type APIType struct {
	TypeDeclaration
	ObjectType
	ScalarType
	String
	ArrayType
}

// BeforeUnmarshalYAML implement yaml Initiator
func (t *APIType) BeforeUnmarshalYAML() (err error) {
	if err = t.ObjectType.BeforeUnmarshalYAML(); err != nil {
		return
	}
	if err = t.String.BeforeUnmarshalYAML(); err != nil {
		return
	}
	if err = t.ArrayType.BeforeUnmarshalYAML(); err != nil {
		return
	}
	return
}

// UnmarshalYAML implement yaml unmarshaler
func (t *APIType) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	if err = unmarshaler(&t.TypeDeclaration); err != nil {
		return
	}
	if strings.HasSuffix(t.TypeDeclaration.Type, "[]") {
		if err = unmarshaler(&t.ArrayType); err != nil {
			return
		}
	}
	if err = unmarshaler(&t.ObjectType); err != nil {
		return
	}
	if !t.ObjectType.IsEmpty() {
		t.Type = typeObject
		return nil
	}
	if err = unmarshaler(&t.ScalarType); err != nil {
		return
	}
	if err = unmarshaler(&t.String); err != nil {
		return
	}
	return nil
}

// MarshalJSON marshal to json
func (t APIType) MarshalJSON() ([]byte, error) {
	return MarshalJSONWithoutEmptyStruct(t)
}

// PostProcess for fill some field from RootDocument default config
func (t *APIType) PostProcess(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}
	if err = t.ObjectType.PostProcess(conf); err != nil {
		return
	}
	if err = t.ScalarType.PostProcess(conf); err != nil {
		return
	}
	if err = t.String.PostProcess(conf); err != nil {
		return
	}
	if err = t.ArrayType.PostProcess(conf); err != nil {
		return
	}

	// TypeDeclaration should go after other basic proprtyies done
	if err = t.TypeDeclaration.PostProcess(conf, *t); err != nil {
		return
	}
	return
}

// IsEmpty return true if it is empty
func (t APIType) IsEmpty() bool {
	return t.TypeDeclaration.IsEmpty() &&
		t.ObjectType.IsEmpty() &&
		t.ScalarType.IsEmpty() &&
		t.String.IsEmpty()
}
