package parser

import (
	"encoding/json"
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

func testCheckValueAPIType(apiType APIType, value interface{}, options ...CheckValueOption) (err error) {
	v, err := NewValue(value)
	if err != nil {
		return
	}

	return CheckValueAPIType(apiType, v, options...)
}

func Test_CheckValueAPIType_Boolean(t *testing.T) {
	var err error
	require := require.New(t)
	require.NotNil(require)

	apiType := APIType{}
	apiType.Type = TypeBoolean

	err = testCheckValueAPIType(apiType, true)
	require.NoError(err)

	err = testCheckValueAPIType(apiType, 0)
	require.Error(err)

	err = testCheckValueAPIType(apiType, 0.1)
	require.Error(err)

	err = testCheckValueAPIType(apiType, "")
	require.Error(err)

	err = testCheckValueAPIType(apiType, map[string]interface{}{})
	require.Error(err)
}

func Test_CheckValueAPIType_Integer(t *testing.T) {
	var err error
	require := require.New(t)
	require.NotNil(require)

	apiType := APIType{}
	apiType.Type = TypeInteger

	err = testCheckValueAPIType(apiType, true)
	require.Error(err)

	err = testCheckValueAPIType(apiType, 0)
	require.NoError(err)

	err = testCheckValueAPIType(apiType, 0.1)
	require.Error(err)

	err = testCheckValueAPIType(apiType, "")
	require.Error(err)

	err = testCheckValueAPIType(apiType, map[string]interface{}{})
	require.Error(err)
}

func Test_CheckValueAPIType_Number(t *testing.T) {
	var err error
	require := require.New(t)
	require.NotNil(require)

	apiType := APIType{}
	apiType.Type = TypeNumber

	err = testCheckValueAPIType(apiType, true)
	require.Error(err)

	err = testCheckValueAPIType(apiType, 0)
	require.Error(err)

	err = testCheckValueAPIType(apiType, 0, CheckValueOptionAllowIntegerToBeNumber(true))
	require.NoError(err)

	err = testCheckValueAPIType(apiType, 0.1)
	require.NoError(err)

	err = testCheckValueAPIType(apiType, "")
	require.Error(err)

	err = testCheckValueAPIType(apiType, map[string]interface{}{})
	require.Error(err)
}

func Test_CheckValueAPIType_String(t *testing.T) {
	var err error
	require := require.New(t)
	require.NotNil(require)

	apiType := APIType{}
	apiType.Type = TypeString

	err = testCheckValueAPIType(apiType, true)
	require.Error(err)

	err = testCheckValueAPIType(apiType, 0)
	require.Error(err)

	err = testCheckValueAPIType(apiType, 0.1)
	require.Error(err)

	err = testCheckValueAPIType(apiType, "")
	require.NoError(err)

	err = testCheckValueAPIType(apiType, map[string]interface{}{})
	require.Error(err)
}

func Test_CheckValueAPIType_Object(t *testing.T) {
	var err error
	require := require.New(t)
	require.NotNil(require)

	apiType := APIType{}
	apiType.Type = "obj"
	apiType.Properties = Properties{}

	property := &Property{}
	property.Type = TypeString
	apiType.Properties["text"] = property

	property = &Property{}
	property.Type = TypeInteger
	apiType.Properties["int"] = property

	property = &Property{}
	property.Type = TypeNumber
	apiType.Properties["num"] = property

	err = testCheckValueAPIType(apiType, true)
	require.Error(err)

	err = testCheckValueAPIType(apiType, 0)
	require.Error(err)

	err = testCheckValueAPIType(apiType, 0.1)
	require.Error(err)

	err = testCheckValueAPIType(apiType, "")
	require.Error(err)

	err = testCheckValueAPIType(apiType, map[string]interface{}{
		"text": "",
		"int":  0,
		"num":  0.1,
	})
	require.NoError(err)

	err = testCheckValueAPIType(apiType, map[string]interface{}{
		"text": 0,
		"int":  0,
		"num":  0.1,
	})
	require.Error(err)

	err = testCheckValueAPIType(apiType, map[string]interface{}{
		"text": "",
		"int":  0.1,
		"num":  0.1,
	})
	require.Error(err)

	err = testCheckValueAPIType(apiType, map[string]interface{}{
		"text": "",
		"int":  0,
		"num":  0,
	})
	require.Error(err)

	err = testCheckValueAPIType(apiType, map[string]interface{}{
		"text": "",
		"int":  float64(0),
		"num":  int64(0),
	}, CheckValueOptionAllowIntegerToBeNumber(true))
	require.NoError(err)

	num := float64(123)
	require.True(float64(int64(num)) == num)
	num = float64(1.23)
	require.False(float64(int64(num)) == num)

	valmap := map[string]interface{}{}
	err = json.Unmarshal([]byte(`{
		"text": "",
		"int": 0,
		"num": 0
	}`), &valmap)
	require.NoError(err)
	err = testCheckValueAPIType(apiType, valmap, CheckValueOptionAllowIntegerToBeNumber(true))
	require.NoError(err)
}
