package parser

// Bodies map of Body
type Bodies map[string]*Body

// Bodies contains Body types, necessary because of technical reasons.
// type Bodies struct {
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
// ForMIMEType map[string]*Body `yaml:",regexp:.*"`
// }

// UnmarshalYAML unmarshal from YAML
func (t *Bodies) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	mimetype := map[string]*Body{}
	if err = unmarshaler(&mimetype); err == nil {
		*t = mimetype
		return
	}

	body := Body{}
	if err = unmarshaler(&body); err == nil {
		(*t)["DEFAULT"] = &body
		return
	}

	return
}

// PostProcess for fill default example by type if not set
func (t *Bodies) PostProcess(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}

	if body, exist := (*t)["DEFAULT"]; exist {
		if conf.RootDocument().MediaType == "" {
			return ErrorEmptyRootDocumentMediaType.New(nil)
		}
		delete(*t, "DEFAULT")
		(*t)[conf.RootDocument().MediaType] = body
	}

	for _, body := range *t {
		if err = body.PostProcess(conf); err != nil {
			return
		}
	}
	return
}

// IsEmpty return true if it is empty
func (t Bodies) IsEmpty() bool {
	for _, elem := range t {
		if elem != nil {
			if !elem.IsEmpty() {
				return false
			}
		}
	}
	return true
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
func (t *Body) PostProcess(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}
	if err = t.APIType.PostProcess(conf); err != nil {
		return
	}
	return
}
