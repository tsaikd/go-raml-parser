package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tsaikd/KDGoLib/jsonex"
)

func testCheckValueAPIType(apiType APIType, value interface{}, options ...CheckValueOption) (err error) {
	var v Value
	if v, err = NewValue(value); err != nil {
		return
	}
	return CheckValueAPIType(apiType, v, options...)
}

func Test_CheckValueAPIType_Boolean(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	apiType := getAPITypeFromString(`
type: boolean
	`)
	require.NoError(testCheckValueAPIType(apiType, true))
	require.Error(testCheckValueAPIType(apiType, 0))
	require.Error(testCheckValueAPIType(apiType, 0.1))
	require.Error(testCheckValueAPIType(apiType, ""))
	require.Error(testCheckValueAPIType(apiType, map[string]interface{}{}))
}

func Test_CheckValueAPIType_Integer(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	apiType := getAPITypeFromString(`
type: integer
	`)
	require.Error(testCheckValueAPIType(apiType, true))
	require.NoError(testCheckValueAPIType(apiType, 0))
	require.Error(testCheckValueAPIType(apiType, 0.1))
	require.Error(testCheckValueAPIType(apiType, ""))
	require.Error(testCheckValueAPIType(apiType, map[string]interface{}{}))
}

func Test_CheckValueAPIType_Number(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	apiType := getAPITypeFromString(`
type: number
	`)
	require.Error(testCheckValueAPIType(apiType, true))
	require.Error(testCheckValueAPIType(apiType, 0))
	require.NoError(testCheckValueAPIType(apiType, 0, CheckValueOptionAllowIntegerToBeNumber(true)))
	require.NoError(testCheckValueAPIType(apiType, 0.1))
	require.Error(testCheckValueAPIType(apiType, ""))
	require.Error(testCheckValueAPIType(apiType, map[string]interface{}{}))
}

func Test_CheckValueAPIType_String(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	apiType := getAPITypeFromString(`
type: string
	`)
	require.Error(testCheckValueAPIType(apiType, true))
	require.Error(testCheckValueAPIType(apiType, 0))
	require.Error(testCheckValueAPIType(apiType, 0.1))
	require.NoError(testCheckValueAPIType(apiType, ""))
	require.Error(testCheckValueAPIType(apiType, map[string]interface{}{}))
}

func Test_CheckValueAPIType_Object(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	apiType := getAPITypeFromString(`
type: object
properties:
    text: string
    int:  integer
    num:  number
	`)
	require.Error(testCheckValueAPIType(apiType, true))
	require.Error(testCheckValueAPIType(apiType, 0))
	require.Error(testCheckValueAPIType(apiType, 0.1))
	require.Error(testCheckValueAPIType(apiType, ""))
	require.Error(testCheckValueAPIType(apiType,
		map[string]interface{}{
			"text": "",
			"int":  0,
			"num":  0.1,
		},
	))
	require.NoError(testCheckValueAPIType(apiType,
		map[string]interface{}{
			"text": "",
			"int":  0,
			"num":  0.1,
		},
		CheckValueOptionAllowRequiredPropertyToBeEmpty(true),
	))
	require.NoError(testCheckValueAPIType(apiType,
		map[string]interface{}{
			"text": "test string",
			"int":  0,
			"num":  0.1,
		},
	))
	require.Error(testCheckValueAPIType(apiType,
		map[string]interface{}{
			"text": 0,
			"int":  0,
			"num":  0.1,
		},
	))
	require.Error(testCheckValueAPIType(apiType,
		map[string]interface{}{
			"text": "test string",
			"int":  0.1,
			"num":  0.1,
		},
	))
	require.Error(testCheckValueAPIType(apiType,
		map[string]interface{}{
			"text": "test string",
			"int":  0,
			"num":  0,
		},
	))
	require.NoError(testCheckValueAPIType(apiType,
		map[string]interface{}{
			"text": "test string",
			"int":  float64(0),
			"num":  int64(0),
		},
		CheckValueOptionAllowIntegerToBeNumber(true),
	))

	num := float64(123)
	require.True(float64(int64(num)) == num)
	num = float64(1.23)
	require.False(float64(int64(num)) == num)

	valmap := map[string]interface{}{}
	err := jsonex.Unmarshal([]byte(`{
		"text": "test string",
		"int": 0,
		"num": 0
	}`), &valmap)
	require.NoError(err)
	require.NoError(testCheckValueAPIType(apiType, valmap, CheckValueOptionAllowIntegerToBeNumber(true)))
}

func Test_CheckValueAPIType_Array(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	apiType := getAPITypeFromString(`
type: string[]
	`)
	require.NoError(testCheckValueAPIType(apiType, []string{"text"}))
	require.Error(testCheckValueAPIType(apiType, nil))
	require.NoError(testCheckValueAPIType(apiType, nil, CheckValueOptionAllowArrayToBeNull(true)))
}

func Test_CheckValueAPIType_ArrayInObject(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	apiType := getAPITypeFromString(`
type: object
properties:
    strs: string[]
	`)
	require.Error(testCheckValueAPIType(apiType, nil))
	require.NoError(testCheckValueAPIType(apiType, nil, CheckValueOptionAllowArrayToBeNull(true)))

	strsObj := map[string]interface{}{
		"strs": nil,
	}
	require.Error(testCheckValueAPIType(apiType, strsObj))
	require.NoError(testCheckValueAPIType(apiType, strsObj, CheckValueOptionAllowArrayToBeNull(true)))
}

func Test_CheckValueAPIType_ObjectArray(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	apiType := getAPITypeFromString(`
type: object[]
properties:
    text: string
	`)
	require.Error(testCheckValueAPIType(apiType,
		[]interface{}{
			map[string]interface{}{
				"text": "",
			},
		},
	))
	require.NoError(testCheckValueAPIType(apiType,
		[]interface{}{
			map[string]interface{}{
				"text": "",
			},
		},
		CheckValueOptionAllowRequiredPropertyToBeEmpty(true),
	))
	require.NoError(testCheckValueAPIType(apiType,
		[]interface{}{
			map[string]interface{}{
				"text": "test string",
			},
		},
	))
	require.Error(testCheckValueAPIType(apiType,
		[]interface{}{
			map[string]interface{}{
				"text": 0,
			},
		},
	))
}

func Test_CheckValueAPIType_ObjectDeep(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	apiType := getAPITypeFromString(`
type: object
properties:
    level1:
        type: object
        properties:
            level2:
                type: object
                properties:
                    text: string
	`)
	require.NoError(testCheckValueAPIType(apiType,
		map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"text": "test string",
				},
			},
		},
	))
	require.Error(testCheckValueAPIType(apiType,
		map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"text": "",
				},
			},
		},
	))
	require.Error(testCheckValueAPIType(apiType,
		map[string]interface{}{
			"level1": map[string]interface{}{},
		},
	))

	apiType = getAPITypeFromString(`
type: object
properties:
    level1:
        type: object
        properties:
            lv2sibling: string
            level2?:
                type: object
                properties:
                    text: string
	`)
	require.NoError(testCheckValueAPIType(apiType,
		map[string]interface{}{
			"level1": map[string]interface{}{
				"lv2sibling": "test string",
				"level2": map[string]interface{}{
					"text": "",
				},
			},
		},
	))
	require.NoError(testCheckValueAPIType(apiType,
		map[string]interface{}{
			"level1": map[string]interface{}{
				"lv2sibling": "test string",
			},
		},
	))
}

func Test_CheckValueAPIType_CustomType(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	_, err := parser.ParseData([]byte(strings.TrimSpace(`
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
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	_, err := parser.ParseData([]byte(strings.TrimSpace(`
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
