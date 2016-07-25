package parser

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tsaikd/KDGoLib/jsonex"
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

func testCheckValueAPIType(apiType APIType, value interface{}, options ...CheckValueOption) (err error) {
	var v Value
	if v, err = NewValue(value); err != nil {
		return
	}
	return CheckValueAPIType(apiType, v, options...)
}

func Test_CheckValueAPIType_Boolean(t *testing.T) {
	var err error
	require := require.New(t)
	require.NotNil(require)

	apiType := APIType{}
	apiType.setType(TypeBoolean)

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
	apiType.setType(TypeInteger)

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
	apiType.setType(TypeNumber)

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
	apiType.setType(TypeString)

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
	apiType.setType(TypeObject)
	apiType.Properties = Properties{}

	property := &Property{}
	property.Name = "text"
	property.setType(TypeString)
	addProperty(&apiType.Properties, property)

	property = &Property{}
	property.Name = "int"
	property.setType(TypeInteger)
	addProperty(&apiType.Properties, property)

	property = &Property{}
	property.Name = "num"
	property.setType(TypeNumber)
	addProperty(&apiType.Properties, property)

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
	err = jsonex.Unmarshal([]byte(`{
		"text": "",
		"int": 0,
		"num": 0
	}`), &valmap)
	require.NoError(err)
	err = testCheckValueAPIType(apiType, valmap, CheckValueOptionAllowIntegerToBeNumber(true))
	require.NoError(err)
}

func Test_CheckValueAPIType_Array(t *testing.T) {
	var err error
	require := require.New(t)
	require.NotNil(require)

	apiType := APIType{}
	apiType.setType("string[]")

	err = testCheckValueAPIType(apiType, []string{"text"})
	require.NoError(err)

	err = testCheckValueAPIType(apiType, nil)
	require.Error(err)

	err = testCheckValueAPIType(apiType, nil, CheckValueOptionAllowArrayToBeNull(true))
	require.NoError(err)
}

func Test_CheckValueAPIType_ObjectArray(t *testing.T) {
	var err error
	require := require.New(t)
	require.NotNil(require)

	apiType := APIType{}
	apiType.setType("object[]")
	apiType.Properties = Properties{}

	property := &Property{}
	property.Name = "text"
	property.setType(TypeString)
	addProperty(&apiType.Properties, property)

	err = testCheckValueAPIType(apiType, []interface{}{
		map[string]interface{}{
			"text": "",
		},
	})
	require.NoError(err)

	err = testCheckValueAPIType(apiType, []interface{}{
		map[string]interface{}{
			"text": 0,
		},
	})
	require.Error(err)
}

func Test_CheckValueAPIType_CustomType(t *testing.T) {
	var err error
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	_, err = parser.ParseData([]byte(strings.TrimSpace(`
#%RAML 1.0
types:
    Username:
        type: string
    User:
        type: object
        properties:
            name: Username
        example:
            name: Alice
	`)), ".")
	require.NoError(err)
}

func Test_CheckExampleAPIType(t *testing.T) {
	var err error
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	_, err = parser.ParseData([]byte(strings.TrimSpace(`
#%RAML 1.0
types:
    User:
        type: object
        properties:
            name:  string
            email: string
        example:
            name:  Alice
	`)), ".")
	require.Error(err)
	require.True(ErrorRequiredProperty2.Match(err))
}

func addProperty(properties *Properties, property *Property) {
	if properties.mapdata == nil {
		properties.mapdata = map[string]*Property{}
	}
	if properties.propertiesSliceData == nil {
		properties.propertiesSliceData = []*Property{}
	}
	if _, exist := properties.mapdata[property.Name]; exist {
		return
	}
	properties.mapdata[property.Name] = property
	properties.propertiesSliceData = append(properties.propertiesSliceData, property)
}
