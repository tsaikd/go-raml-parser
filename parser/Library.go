package parser

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/tsaikd/yaml"
)

// Libraries map of Library
type Libraries map[string]*Library

// IsEmpty return true if it is empty
func (t Libraries) IsEmpty() bool {
	for _, elem := range t {
		if elem != nil {
			if !elem.IsEmpty() {
				return false
			}
		}
	}
	return true
}

var _ loadExternalUse = Libraries{}

func (t Libraries) loadExternalUse(conf PostProcessConfig) (err error) {
	for name, library := range t {
		filePath := filepath.Join(conf.RootDocument().WorkingDirectory, library.Name)

		fileData, err := ioutil.ReadFile(filePath)
		if err != nil {
			return ErrorLoadExternalLibrary1.New(err, filePath)
		}

		if err = yaml.Unmarshal(fileData, library); err != nil {
			return ErrorLoadExternalLibrary1.New(err, filePath)
		}

		library.Name = name
	}
	return
}

// Library wrap LibraryRAML because LibraryRAML may be a string for external library file
type Library struct {
	Name string `json:",omitempty"`

	LibraryRAML
}

// UnmarshalYAML unmarshal Library from YAML
func (t *Library) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	if err = unmarshaler(&t.Name); err == nil {
		return
	}
	if !isErrorYAMLIntoString(err) {
		return
	}

	if err = unmarshaler(&t.LibraryRAML); err != nil {
		return
	}
	return
}

// IsEmpty return true if it is empty
func (t Library) IsEmpty() bool {
	return t.Name == "" &&
		t.LibraryRAML.IsEmpty()
}

// GetTrait return trait if found
func (t Library) GetTrait(name string) (result Trait, err error) {
	if splits := strings.Split(name, "."); len(splits) == 2 {
		useName, traitName := splits[0], splits[1]
		use, ok := t.Uses[useName]
		if !ok || use == nil {
			err = ErrorUseNotFound1.New(nil, useName)
			return
		}
		return use.GetTrait(traitName)
	}

	trait, ok := t.Traits[name]
	if !ok || trait == nil {
		err = ErrorTraitNotFound1.New(nil, name)
		return
	}

	return *trait, nil
}

// Prefix return "" if Library is not external used
func (t Library) Prefix() string {
	if t.Name == "" {
		return ""
	}
	return t.Name + "."
}

var _ checkUnusedAnnotation = Library{}

func (t Library) checkUnusedAnnotation(conf PostProcessConfig) (err error) {
	annotationUsage := conf.AnnotationUsage()
	for name := range t.AnnotationTypes {
		annotationUsage[name] = true
	}
	return
}

var _ checkUnusedTrait = Library{}

func (t Library) checkUnusedTrait(conf PostProcessConfig) (err error) {
	prefix := t.Prefix()
	traitUsage := conf.TraitUsage()
	for name := range t.Traits {
		traitUsage[prefix+name] = true
	}
	return
}

var _ checkAnnotation = Library{}

func (t Library) checkAnnotation(conf PostProcessConfig) (err error) {
	return t.Annotations.checkAnnotationTargetLocation(TargetLocationLibrary)
}

// LibraryRAML RAML libraries are used to combine any collection of data type
// declarations, resource type declarations, trait declarations, and security
// scheme declarations into modular, externalized, reusable groups.
// While libraries are intended to define common declarations in external
// documents, which are then included where needed, libraries can also
// be defined inline.
type LibraryRAML struct {
	// Describes the content or purpose of a specific library. The value is
	// a string and MAY be formatted using markdown.
	Usage string `yaml:"usage" json:"usage,omitempty"`

	// An alias for the equivalent "types" node for compatibility with
	// RAML 0.8. Deprecated - API definitions should use the "types" node
	// because a future RAML version might remove the "schemas" alias with
	// that node. The "types" node supports XML and JSON schemas.
	Schemas Unimplement `yaml:"schemas" json:"schemas,omitempty"`

	// Declarations of (data) types for use within the API.
	Types APITypes `yaml:"types" json:"types,omitempty"`

	// Declarations of traits for use within the API.
	Traits Traits `yaml:"traits" json:"traits,omitempty"`

	// Declarations of resource types for use within the API.
	ResourceTypes Unimplement `yaml:"resourceTypes" json:"resourceTypes,omitempty"`

	// Declarations of annotation types for use by annotations.
	AnnotationTypes AnnotationTypes `yaml:"annotationTypes" json:"annotationTypes,omitempty"`

	// Annotations to be applied to this API. An annotation is a map having
	// a key that begins with "(" and ends with ")" where the text enclosed
	// in parentheses is the annotation name, and the value is an instance of
	// that annotation.
	Annotations Annotations `yaml:",regexp:\\(.*\\)" json:"annotations,omitempty"`

	// Declarations of security schemes for use within the API.
	SecuritySchemes Unimplement `yaml:"securitySchemes" json:"securitySchemes,omitempty"`

	// Imported external libraries for use within the API.
	Uses Libraries `yaml:"uses" json:"uses,omitempty"`
}

// IsEmpty return true if it is empty
func (t LibraryRAML) IsEmpty() bool {
	return t.Usage == "" &&
		t.Schemas.IsEmpty() &&
		t.Types.IsEmpty() &&
		t.Traits.IsEmpty() &&
		t.ResourceTypes.IsEmpty() &&
		t.AnnotationTypes.IsEmpty() &&
		t.Annotations.IsEmpty() &&
		t.SecuritySchemes.IsEmpty() &&
		t.Uses.IsEmpty()
}
