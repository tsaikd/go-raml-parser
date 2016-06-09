package parser

import (
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
