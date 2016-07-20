package parser

import (
	"path/filepath"
	"strings"

	"github.com/tsaikd/KDGoLib/errutil"
)

// errors
var (
	ErrorTraitNotFound1 = errutil.NewFactory("trait %q not found")
	ErrorUseNotFound1   = errutil.NewFactory("use %q not found")
)

// Libraries map of LibraryWrap
type Libraries map[string]*LibraryWrap

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

// LibraryWrap wrap Library because Library may be a string for external library file
type LibraryWrap struct {
	String string `json:",omitempty"`
	Library
}

// UnmarshalYAML unmarshal LibraryWrap from YAML
func (t *LibraryWrap) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	if err = unmarshaler(&t.String); err == nil {
		return
	}
	if !isErrorYAMLIntoString(err) {
		return
	}

	if err = unmarshaler(&t.Library); err != nil {
		return
	}
	return
}

// IsEmpty return true if it is empty
func (t LibraryWrap) IsEmpty() bool {
	return t.String == "" &&
		t.Library.IsEmpty()
}

var _ loadExternalUse = &LibraryWrap{}

func (t *LibraryWrap) loadExternalUse(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}
	if t.String == "" {
		return
	}

	filePath := filepath.Join(conf.RootDocument().WorkingDirectory, t.String)
	if t.Library, err = conf.Parser().ParseLibraryFile(filePath, conf); err != nil {
		return ErrorLoadExternalLibrary1.New(err, filePath)
	}
	t.String = ""

	return
}

// Library RAML libraries are used to combine any collection of data type
// declarations, resource type declarations, trait declarations, and security
// scheme declarations into modular, externalized, reusable groups.
// While libraries are intended to define common declarations in external
// documents, which are then included where needed, libraries can also
// be defined inline.
type Library struct {
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
func (t Library) IsEmpty() bool {
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

var _ checkUnusedAnnotation = Library{}

func (t Library) checkUnusedAnnotation(annotationUsage map[string]bool) (err error) {
	for name := range t.AnnotationTypes {
		annotationUsage[name] = true
	}
	return
}

var _ checkUnusedTrait = Library{}

func (t Library) checkUnusedTrait(traitUsage map[string]bool) (err error) {
	for name := range t.Traits {
		traitUsage[name] = true
	}
	return
}

var _ checkAnnotation = Library{}

func (t Library) checkAnnotation(conf PostProcessConfig) (err error) {
	return t.Annotations.checkAnnotationTargetLocation(TargetLocationLibrary)
}
