package parser

import "encoding/json"

// Responses map of Response
type Responses map[HTTPCode]*Response

// MarshalJSON marshal to json
func (t Responses) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{}
	for k, v := range t {
		data[k.String()] = v
	}
	return json.Marshal(data)
}

// PostProcess for fill some field from RootDocument default config
func (t *Responses) PostProcess(rootdoc RootDocument) (err error) {
	for _, response := range *t {
		if err = response.PostProcess(rootdoc); err != nil {
			return
		}
	}
	return
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
	Annotations map[string]Unimplement `yaml:",regexp:\\(.*\\)" json:"annotations,omitempty"`

	// Detailed information about any response headers returned by this method
	Headers Unimplement `yaml:"headers" json:"headers,omitempty"`

	// The body of the response
	Bodies Bodies `yaml:"body"`
}

// PostProcess for fill some field from RootDocument default config
func (t *Response) PostProcess(rootdoc RootDocument) (err error) {
	if err = t.Bodies.PostProcess(rootdoc); err != nil {
		return
	}
	return
}
