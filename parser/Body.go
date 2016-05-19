package parser

// Bodies contains Body types, necessary because of technical reasons.
type Bodies struct {
	// Instead of using a simple map[HTTPHeader]Body for the body
	// property of the Response and Method, we use the Bodies struct. Why?
	// Because some RAML APIs don't use the MIMEType part, instead relying
	// on the mediaType property in the APIDefinition.
	// So, you might see:
	//
	// responses:
	//   200:
	//     body:
	//       example: "some_example" : "123"
	//
	// and also:
	//
	// responses:
	//   200:
	//     body:
	//       application/json:
	//         example: |
	//           {
	//             "some_example" : "123"
	//           }

	// Resources CAN have alternate representations. For example, an API
	// might support both JSON and XML representations. This is the map
	// between MIME-type and the body definition related to it.
	ForMIMEType map[string]*Body `yaml:",regexp:.*"`
}

// PostProcess for fill default example by type if not set
func (t *Bodies) PostProcess(rootdoc RootDocument) (err error) {
	for _, body := range t.ForMIMEType {
		if err = body.PostProcess(rootdoc); err != nil {
			return
		}
	}
	return
}

// IsEmpty return true if Bodies is empty
func (t Bodies) IsEmpty() bool {
	return len(t.ForMIMEType) < 1
}

// Body used for Bodies.
// Some method verbs expect the resource to be sent as a request body.
// For example, to create a resource, the request must include the details of
// the resource to create.
// Resources CAN have alternate representations. For example, an API might
// support both JSON and XML representations.
type Body struct {
	APIType
}

// PostProcess for fill some field from RootDocument default config
func (t *Body) PostProcess(rootdoc RootDocument) (err error) {
	if err = t.APIType.PostProcess(rootdoc); err != nil {
		return
	}
	return
}
