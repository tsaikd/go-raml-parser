package parser

// AnnotationTypes map of AnnotationType
type AnnotationTypes map[string]*AnnotationType

// PostProcess for fill some field from RootDocument default config
func (t *AnnotationTypes) PostProcess(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}
	for name, elem := range *t {
		if elem == nil {
			elem = &AnnotationType{}
			elem.Type = "string"
			(*t)[name] = elem
		}
		if err = elem.PostProcess(conf); err != nil {
			return
		}
	}
	return
}

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

// PostProcess for fill some field from RootDocument default config
func (t *AnnotationType) PostProcess(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}
	if err = t.APIType.PostProcess(conf); err != nil {
		return
	}
	if err = t.AllowedTargets.PostProcess(conf); err != nil {
		return
	}
	return
}

// IsEmpty return true if it is empty
func (t AnnotationType) IsEmpty() bool {
	return t.APIType.IsEmpty() &&
		t.AllowedTargets.IsEmpty()
}
