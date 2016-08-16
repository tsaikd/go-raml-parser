package parser

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tsaikd/yaml"
)

func getAPITypeFromString(ramlStr string) APIType {
	apiType := APIType{}
	if err := yaml.Unmarshal([]byte(strings.TrimSpace(ramlStr)), &apiType); err != nil {
		panic(err)
	}

	conf := newPostProcessConfig(nil, nil, nil, nil, nil)
	if err := postProcess(&apiType, conf); err != nil {
		panic(err)
	}

	return apiType
}

func Test_NewValueWithAPIType(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	if value, err := NewValueWithAPIType(getAPITypeFromString(`
type: null
	`), nil); assert.NoError(err) {
		require.Equal(TypeNull, value.Type)
	}

	if value, err := NewValueWithAPIType(getAPITypeFromString(`
type: boolean
	`), true); assert.NoError(err) {
		require.Equal(TypeBoolean, value.Type)
		require.Equal(true, value.Boolean)
	}

	if value, err := NewValueWithAPIType(getAPITypeFromString(`
type: integer
	`), 9527); assert.NoError(err) {
		require.Equal(TypeInteger, value.Type)
		require.EqualValues(9527, value.Integer)
	}

	if value, err := NewValueWithAPIType(getAPITypeFromString(`
type: number
	`), 3.14); assert.NoError(err) {
		require.Equal(TypeNumber, value.Type)
		require.EqualValues(3.14, value.Number)
	}

	if value, err := NewValueWithAPIType(getAPITypeFromString(`
type: string
	`), "test string value"); assert.NoError(err) {
		require.Equal(TypeString, value.Type)
		require.Equal("test string value", value.String)
	}

	if value, err := NewValueWithAPIType(getAPITypeFromString(`
type: string[]
	`), nil); assert.NoError(err) {
		require.Equal(TypeArray, value.Type)
		require.NotNil(value.Array)
		require.Empty(value.Array)
	}

	if value, err := NewValueWithAPIType(
		getAPITypeFromString(`
type: object
properties:
    bool:  boolean
    bools: boolean[]
    int:   integer
    ints:  integer[]
    num:   number
    nums:  number[]
    str:   string
    strs:  string[]
		`),
		url.Values{
			"bool":  []string{"true"},
			"bools": []string{"false", "true"},
			"int":   []string{"9527"},
			"ints":  []string{"5566", "9527"},
			"num":   []string{"3.14"},
			"nums":  []string{"6.02", "3.14"},
			"str":   []string{"test string value"},
			"strs":  []string{"str1", "str2"},
		},
	); assert.NoError(err) {
		require.Equal(TypeObject, value.Type)
		if prop := value.Map["bool"]; assert.NotNil(prop) {
			require.Equal(true, prop.Boolean)
		}
		if prop := value.Map["bools"]; assert.NotNil(prop) {
			require.Len(prop.Array, 2)
			require.Equal(false, prop.Array[0].Boolean)
			require.Equal(true, prop.Array[1].Boolean)
		}
		if prop := value.Map["int"]; assert.NotNil(prop) {
			require.EqualValues(9527, prop.Integer)
		}
		if prop := value.Map["ints"]; assert.NotNil(prop) {
			require.Len(prop.Array, 2)
			require.EqualValues(5566, prop.Array[0].Integer)
			require.EqualValues(9527, prop.Array[1].Integer)
		}
		if prop := value.Map["num"]; assert.NotNil(prop) {
			require.EqualValues(3.14, prop.Number)
		}
		if prop := value.Map["nums"]; assert.NotNil(prop) {
			require.Len(prop.Array, 2)
			require.EqualValues(6.02, prop.Array[0].Number)
			require.EqualValues(3.14, prop.Array[1].Number)
		}
		if prop := value.Map["str"]; assert.NotNil(prop) {
			require.Equal("test string value", prop.String)
		}
		if prop := value.Map["strs"]; assert.NotNil(prop) {
			require.Len(prop.Array, 2)
			require.Equal("str1", prop.Array[0].String)
			require.Equal("str2", prop.Array[1].String)
		}
	}
}
