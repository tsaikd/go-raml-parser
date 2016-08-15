package parser

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_LoadRAMLFromDir(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	expected := `#%RAML 1.0
title: Load raml from directory

#%RAML 1.0
types:
  User:
    type: object
    properties:
      name:  string
      email: string
    examples:
      user1:
        name:  Alice
        email: alice@example.com
      user2:
        name:  Bob
        email: bob@example.com

#%RAML 1.0
/user:
  get:
    responses:
      200:
        body:
          application/json:
            type: User

`

	data, err := LoadRAMLFromDir("test-examples/raml-from-dir")
	require.NoError(err)
	require.Equal(expected, string(data))
}

func ExampleParseYAMLError() {
	err := errors.New("yaml: line 596: did not find expected key")
	line, reason, ok := ParseYAMLError(err)
	fmt.Printf("line: %d\nreason: %q\nok: %v\n", line, reason, ok)
	// Output:
	// line: 596
	// reason: "did not find expected key"
	// ok: true
}

func ExampleGetLinesInRange() {
	data := strings.TrimSpace(`
#%RAML 1.0
types:
    User:
        type: object
        properties:
            name:  string
            email: string
        examples:
            user1:
                name:  Alice
                email: alice@example.com
            user2:
                name:  Bob
                email: bob@example.com
	`)
	fmt.Printf("line 5, distance 1\n%s\n", GetLinesInRange(data, "\n", 5, 1))
	// Output:
	// line 5, distance 1
	//         type: object
	//         properties:
	//             name:  string
}

func Test_ParseYAMLError(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	err := errors.New("yaml: line 596: did not find expected key")
	line, reason, ok := ParseYAMLError(err)
	require.EqualValues(596, line)
	require.Equal("did not find expected key", reason)
	require.True(ok)
}
