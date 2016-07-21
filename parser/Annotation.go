package parser

import (
	"strings"

	"github.com/tsaikd/go-raml-parser/parser/parserConfig"
)

// Annotations map of Annotation
type Annotations map[string]*Annotation

// IsEmpty return true if it is empty
func (t Annotations) IsEmpty() bool {
	for _, elem := range t {
		if elem != nil {
			if !elem.IsEmpty() {
				return false
			}
		}
	}
	return true
}

func (t Annotations) checkAnnotationTargetLocation(targetLocation TargetLocation) (err error) {
	for _, annotation := range t {
		if !annotation.AnnotationType.AllowedTargets.IsAllowed(targetLocation) {
			return ErrorInvalidAnnotationTargetLocation2.New(nil, annotation.Name, targetLocation)
		}
	}
	return nil
}

var _ fixEmptyAnnotation = Annotations{}

func (t Annotations) fixEmptyAnnotation() (err error) {
	for name, elem := range t {
		if elem != nil {
			continue
		}
		elem = &Annotation{}
		elem.Type = TypeNull
		t[name] = elem
	}
	return
}

var _ fixAnnotationBracket = Annotations{}

func (t Annotations) fixAnnotationBracket() (err error) {
	for name, annotation := range t {
		fixedName := name
		fixedName = strings.TrimPrefix(fixedName, "(")
		fixedName = strings.TrimSuffix(fixedName, ")")
		annotation.Name = fixedName
		if fixedName == name {
			continue
		}
		delete(t, name)
		t[fixedName] = annotation
	}
	return
}

var _ checkUnusedAnnotation = Annotations{}

func (t Annotations) checkUnusedAnnotation(conf PostProcessConfig) (err error) {
	annotationUsage := conf.AnnotationUsage()
	for name := range t {
		delete(annotationUsage, conf.Library().Prefix()+name)
	}
	return
}

// Annotation wrap types defined in spec
type Annotation struct {
	Value

	// fill by fixAnnotationBracket
	Name string `yaml:"-"`
	// fill by fillAnnotation
	AnnotationType AnnotationType `yaml:"-"`
}

var _ fillAnnotation = &Annotation{}

func (t *Annotation) fillAnnotation(library Library) (err error) {
	name := t.Name
	annotype := library.AnnotationTypes[name]
	if annotype == nil {
		return ErrorAnnotationTypeUndefined1.New(nil, name)
	}
	t.AnnotationType = *annotype
	return
}

var _ checkAnnotation = Annotation{}

func (t Annotation) checkAnnotation(conf PostProcessConfig) (err error) {
	options := []CheckValueOption{}
	confOptions, err := conf.Parser().Get(parserConfig.CheckValueOptions)
	if err == nil {
		if opts, ok := confOptions.([]CheckValueOption); ok {
			options = opts
		}
	}

	if err = CheckValueAPIType(t.AnnotationType.APIType, t.Value, options...); err != nil {
		return
	}

	return nil
}
