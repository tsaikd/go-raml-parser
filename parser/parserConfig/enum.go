package parserConfig

import "github.com/tsaikd/KDGoLib/enumutil"

// Enum main type
type Enum int8

// List all valid enum
const (
	CheckRAMLVersion Enum = 1 + iota
)

var factory = enumutil.NewEnumFactory().
	Add(CheckRAMLVersion, "CheckRAMLVersion").
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
