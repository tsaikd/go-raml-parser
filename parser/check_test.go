package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tsaikd/KDGoLib/jsonex"
	"github.com/tsaikd/yaml"
)

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
	err = yaml.Unmarshal([]byte(strings.TrimSpace(`
type: object
properties:
    text: string
    int:  integer
    num:  number
	`)), &apiType)
	require.NoError(err)

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
	require.Error(err)

	err = testCheckValueAPIType(apiType, map[string]interface{}{
		"text": "",
		"int":  0,
		"num":  0.1,
	}, CheckValueOptionAllowRequiredPropertyToBeEmpty(true))
	require.NoError(err)

	err = testCheckValueAPIType(apiType, map[string]interface{}{
		"text": "test string",
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
		"text": "test string",
		"int":  0.1,
		"num":  0.1,
	})
	require.Error(err)

	err = testCheckValueAPIType(apiType, map[string]interface{}{
		"text": "test string",
		"int":  0,
		"num":  0,
	})
	require.Error(err)

	err = testCheckValueAPIType(apiType, map[string]interface{}{
		"text": "test string",
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
		"text": "test string",
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

func Test_CheckValueAPIType_ArrayInObject(t *testing.T) {
	var err error
	require := require.New(t)
	require.NotNil(require)

	apiType := APIType{}
	err = yaml.Unmarshal([]byte(strings.TrimSpace(`
type: object
properties:
    strs: string[]
	`)), &apiType)
	require.NoError(err)

	err = testCheckValueAPIType(apiType, nil)
	require.Error(err)
	err = testCheckValueAPIType(apiType, nil, CheckValueOptionAllowArrayToBeNull(true))
	require.NoError(err)

	strsObj := map[string]interface{}{
		"strs": nil,
	}
	err = testCheckValueAPIType(apiType, strsObj)
	require.Error(err)
	err = testCheckValueAPIType(apiType, strsObj, CheckValueOptionAllowArrayToBeNull(true))
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
