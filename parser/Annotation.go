package parser

// Annotations map of Annotation
type Annotations map[string]*Annotation

// PostProcess for fill some field from RootDocument default config
func (t *Annotations) PostProcess(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}
	for _, annotation := range *t {
		if err = annotation.PostProcess(conf); err != nil {
			return
		}
	}
	return
}

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

// Annotation wrap types defined in spec
type Annotation struct {
	Value
}
