package parser

import "github.com/tsaikd/KDGoLib/enumutil"

// TargetLocations slice of TargetLocation
type TargetLocations []TargetLocation

// UnmarshalYAML implement yaml unmarshaler
func (t *TargetLocations) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	// allowedTargets could be a string or string array
	*t = TargetLocations{}
	tarlocs := []string{}

	bufstr := ""
	if err = unmarshaler(&bufstr); err != nil {
		if !isErrorYAMLIntoString(err) {
			return
		}
		if err = unmarshaler(&tarlocs); err != nil {
			return
		}
	} else {
		tarlocs = []string{bufstr}
	}

	for _, tarloc := range tarlocs {
		enum := ParseTargetLocation(tarloc)
		if enum < 1 {
			return ErrorInvalidTargetLocation1.New(nil, tarloc)
		}
		*t = append(*t, enum)
	}

	return nil
}

// IsEmpty return true if it is empty
func (t TargetLocations) IsEmpty() bool {
	return len(t) < 1
}

// IsAllowed return true if targetLocation is allowed
func (t TargetLocations) IsAllowed(targetLocation TargetLocation) bool {
	if t.IsEmpty() {
		// not limited for target location
		return true
	}
	for _, tarloc := range t {
		if targetLocation == tarloc {
			return true
		}
	}
	return false
}

// TargetLocation the location within an API specification where annotations
// can be applied MUST be one of the target locations in the following
// Target Locations table. The targets are the locations themselves,
// not sub-properties within the locations; for example, the Method target
// refers to the method node, not to the method display name, description,
// and so on.
type TargetLocation int8

// List all valid enum
const (
	// The root of a RAML document
	TargetLocationAPI TargetLocation = 1 + iota
	// An item in the collection of items that is the value of the root-level documentation node
	TargetLocationDocumentationItem
	// A resource (relative URI) node, root-level or nested
	TargetLocationResource
	// A method node
	TargetLocationMethod
	// A declaration of the responses node, whose key is an HTTP status code
	TargetLocationResponse
	// The body node of a method
	TargetLocationRequestBody
	// The body node of a response
	TargetLocationResponseBody
	// A data type declaration (inline or in a global types collection),
	// header declaration, query parameter declaration, URI parameter declaration,
	// or a property within any of these declarations, where the type property
	// can be used
	TargetLocationTypeDeclaration
	// Either an example or examples node
	TargetLocationExample
	// A resource type node
	TargetLocationResourceType
	// A trait node
	TargetLocationTrait
	// A security scheme declaration
	TargetLocationSecurityScheme
	// The settings node of a security scheme declaration
	TargetLocationSecuritySchemeSettings
	// A declaration of the annotationTypes node, whose key is a name of an annotation type and whose value describes the annotation
	TargetLocationAnnotationType
	// The root of a library
	TargetLocationLibrary
	// The root of an overlay
	TargetLocationOverlay
	// The root of an extension
	TargetLocationExtension
)

var factoryTargetLocation = enumutil.NewEnumFactory().
	Add(TargetLocationAPI, "API").
	Add(TargetLocationDocumentationItem, "DocumentationItem").
	Add(TargetLocationResource, "Resource").
	Add(TargetLocationMethod, "Method").
	Add(TargetLocationResponse, "Response").
	Add(TargetLocationRequestBody, "RequestBody").
	Add(TargetLocationResponseBody, "ResponseBody").
	Add(TargetLocationTypeDeclaration, "TypeDeclaration").
	Add(TargetLocationExample, "Example").
	Add(TargetLocationResourceType, "ResourceType").
	Add(TargetLocationTrait, "Trait").
	Add(TargetLocationSecurityScheme, "SecurityScheme").
	Add(TargetLocationSecuritySchemeSettings, "SecuritySchemeSettings").
	Add(TargetLocationAnnotationType, "AnnotationType").
	Add(TargetLocationLibrary, "Library").
	Add(TargetLocationOverlay, "Overlay").
	Add(TargetLocationExtension, "Extension").
	Build()

func (t TargetLocation) String() string {
	return factoryTargetLocation.String(t)
}

// MarshalJSON return jsonfy []byte of enum
func (t TargetLocation) MarshalJSON() ([]byte, error) {
	return factoryTargetLocation.MarshalJSON(t)
}

// UnmarshalJSON decode json data to enum
func (t *TargetLocation) UnmarshalJSON(b []byte) (err error) {
	return factoryTargetLocation.UnmarshalJSON(t, b)
}

// ParseTargetLocation string to enum
func ParseTargetLocation(s string) TargetLocation {
	enum, err := factoryTargetLocation.ParseString(s)
	if err != nil {
		return 0
	}
	return enum.(TargetLocation)
}
