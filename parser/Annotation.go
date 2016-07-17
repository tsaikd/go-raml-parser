package parser

import "strings"

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

// Annotation wrap types defined in spec
type Annotation struct {
	Value
}
