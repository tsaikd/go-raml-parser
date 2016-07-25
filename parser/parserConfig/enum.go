package parserConfig

import "github.com/tsaikd/KDGoLib/enumutil"

// Enum main type
type Enum int8

// List all valid enum
const (
	// RAML parser cache directory, type: string, default: ""
	CacheDirectory Enum = 1 + iota
	// RAML parser should check RAML version or not, type: bool, default: false
	CheckRAMLVersion
	// options pass to CheckValueAPIType, type: []CheckValueOption, default: []CheckValueOption{}
	CheckValueOptions
	// show RAML data when error occur, set < 0 to disable, type: int64, default: 4
	ErrorTraceDistance
	// RAML parser should ignore unused annotations, type: bool, default: false
	IgnoreUnusedAnnotation
	// RAML parser should ignore unused traits, type: bool, default: false
	IgnoreUnusedTrait
)

var factory = enumutil.NewEnumFactory().
	Add(CacheDirectory, "CacheDirectory").
	Add(CheckRAMLVersion, "CheckRAMLVersion").
	Add(CheckValueOptions, "CheckValueOptions").
	Add(ErrorTraceDistance, "ErrorTraceDistance").
	Add(IgnoreUnusedAnnotation, "IgnoreUnusedAnnotation").
	Add(IgnoreUnusedTrait, "IgnoreUnusedTrait").
	Build()

func (t Enum) String() string {
	return factory.String(t)
}

// MarshalJSON return jsonfy []byte of enum
func (t Enum) MarshalJSON() ([]byte, error) {
	return factory.MarshalJSON(t)
}

// UnmarshalJSON decode json data to enum
func (t *Enum) UnmarshalJSON(b []byte) (err error) {
	return factory.UnmarshalJSON(t, b)
}

// Is check string is valid enum
func Is(s string) bool {
	return factory.IsEnumString(s)
}

// Parse string to enum
func Parse(s string) Enum {
	enum, err := factory.ParseString(s)
	if err != nil {
		return 0
	}
	return enum.(Enum)
}
