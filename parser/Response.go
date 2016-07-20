package parser

import "github.com/tsaikd/KDGoLib/jsonex"

// Responses map of Response
type Responses map[HTTPCode]*Response

// MarshalJSON marshal to json
func (t Responses) MarshalJSON() ([]byte, error) {
	if t.IsEmpty() {
		return jsonNull, nil
	}

	data := map[string]interface{}{}
	for k, v := range t {
		data[k.String()] = v
	}
	return jsonex.Marshal(data)
}

// IsEmpty return true if it is empty
func (t Responses) IsEmpty() bool {
	for _, elem := range t {
		if elem != nil {
			if !elem.IsEmpty() {
				return false
			}
		}
	}
	return true
}

// Response The value of a response declaration is a map that can contain any
// of the following key-value pairs:
type Response struct {
	// A substantial, human-friendly description of a response. Its value is
	// a string and MAY be formatted using markdown.
	Description string `yaml:"description" json:"description,omitempty"`

	// Annotations to be applied to this API. An annotation is a map having
	// a key that begins with "(" and ends with ")" where the text enclosed
	// in parentheses is the annotation name, and the value is an instance of
	// that annotation.
	Annotations Annotations `yaml:",regexp:\\(.*\\)" json:"annotations,omitempty"`

	// Detailed information about any response headers returned by this method
	Headers Headers `yaml:"headers" json:"headers,omitempty"`

	// The body of the response
	Bodies Bodies `yaml:"body" json:"body,omitempty"`
}

// IsEmpty return true if it is empty
func (t Response) IsEmpty() bool {
	return t.Description == "" &&
		t.Annotations.IsEmpty() &&
		t.Headers.IsEmpty() &&
		t.Bodies.IsEmpty()
}

var _ checkAnnotation = Response{}

func (t Response) checkAnnotation(conf PostProcessConfig) (err error) {
	if err = t.Annotations.checkAnnotationTargetLocation(TargetLocationResponse); err != nil {
		return
	}
	if err = t.Bodies.checkAnnotationTargetLocation(TargetLocationResponseBody); err != nil {
		if err = t.Bodies.checkAnnotationTargetLocation(TargetLocationTypeDeclaration); err != nil {
			return
		}
	}
	return nil
}
