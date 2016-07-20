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
		if fixedName == name {
			continue
		}
		delete(t, name)
		t[fixedName] = annotation
	}
	return
}

var _ checkUnusedAnnotation = Annotations{}

func (t Annotations) checkUnusedAnnotation(annotationUsage map[string]bool) (err error) {
	for name := range t {
		delete(annotationUsage, name)
	}
	return
}

var _ checkExample = Annotations{}

func (t Annotations) checkExample(conf PostProcessConfig) (err error) {
	options := []CheckValueOption{}
	confOptions, err := conf.Parser().Get(parserConfig.CheckValueOptions)
	if err == nil {
		if opts, ok := confOptions.([]CheckValueOption); ok {
			options = opts
		}
	}

	for name, anno := range t {
		annotype := conf.Library().AnnotationTypes[name]
		if annotype == nil {
			return ErrorAnnotationTypeUndefined1.New(nil, name)
		}
		if err = CheckValueAPIType(annotype.APIType, anno.Value, options...); err != nil {
			return
		}
	}

	return nil
}

// Annotation wrap types defined in spec
type Annotation struct {
	Value
}
