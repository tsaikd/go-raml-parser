package parser

// AnnotationTypes map of AnnotationType
type AnnotationTypes map[string]*AnnotationType

// IsEmpty return true if it is empty
func (t AnnotationTypes) IsEmpty() bool {
	for _, elem := range t {
		if elem != nil {
			if !elem.IsEmpty() {
				return false
			}
		}
	}
	return true
}

var _ fixEmptyAnnotation = AnnotationTypes{}

func (t AnnotationTypes) fixEmptyAnnotation() (err error) {
	if t == nil {
		return
	}
	for name, elem := range t {
		if elem == nil {
			elem = &AnnotationType{}
			elem.Type = "string"
			t[name] = elem
		}
	}
	return
}

// AnnotationType wrap types defined in spec
type AnnotationType struct {
	APIType

	// The locations to which annotations are restricted. If this node
	// is specified, annotations of this type may be applied only on
	// a node corresponding to one of the locations.
	// The value MUST be one or more of the options described in the
	// Target Locations.
	AllowedTargets Unimplement `yaml:"allowedTargets" json:"allowedTargets,omitempty"`
}

// IsEmpty return true if it is empty
func (t AnnotationType) IsEmpty() bool {
	return t.APIType.IsEmpty() &&
		t.AllowedTargets.IsEmpty()
}
