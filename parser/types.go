package parser

import "strconv"

// Any type, for our convenience
type Any interface{}

// Unimplement For extra clarity
type Unimplement interface{}

// HTTPCode For extra clarity
type HTTPCode int // e.g. 200

func (t HTTPCode) String() string {
	return strconv.Itoa(int(t))
}
