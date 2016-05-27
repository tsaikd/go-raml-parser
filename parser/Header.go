package parser

// Headers map of Header
type Headers map[string]*Header

// PostProcess for fill some field from RootDocument default config
func (t *Headers) PostProcess(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}
	for _, header := range *t {
		if err = header.PostProcess(conf); err != nil {
			return
		}
	}
	return
}

// IsEmpty return true if Headers is empty
func (t Headers) IsEmpty() bool {
	return len(t) < 1
}

// Header An API's methods can support or require various HTTP headers.
// The OPTIONAL headers node is used to explicitly specify those headers.
// The value of the headers node is a map, specifically a properties
// declaration, as is the value of the properties object of a type declaration.
// Each property in this declaration object is a header declaration.
// Each property name specifies an allowed header name. Each property value
// specifies the header value type as a type name or an inline type declaration.
type Header struct {
	Property
}
