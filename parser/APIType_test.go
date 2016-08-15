package parser

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tsaikd/yaml"
)

func Test_NewValueWithAPIType(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	var apiType APIType
	var value Value
	var err error

	apiType = APIType{}
	apiType.setType(TypeNull)
	value, err = NewValueWithAPIType(apiType, nil)
	require.NoError(err)
	require.Equal(TypeNull, value.Type)

	apiType = APIType{}
	apiType.setType(TypeBoolean)
	value, err = NewValueWithAPIType(apiType, true)
	require.NoError(err)
	require.Equal(TypeBoolean, value.Type)
	require.Equal(true, value.Boolean)

	apiType = APIType{}
	apiType.setType(TypeInteger)
	value, err = NewValueWithAPIType(apiType, 9527)
	require.NoError(err)
	require.Equal(TypeInteger, value.Type)
	require.EqualValues(9527, value.Integer)

	apiType = APIType{}
	apiType.setType(TypeNumber)
	value, err = NewValueWithAPIType(apiType, 3.14)
	require.NoError(err)
	require.Equal(TypeNumber, value.Type)
	require.EqualValues(3.14, value.Number)

	apiType = APIType{}
	apiType.setType(TypeString)
	value, err = NewValueWithAPIType(apiType, "test string value")
	require.NoError(err)
	require.Equal(TypeString, value.Type)
	require.Equal("test string value", value.String)

	apiType = APIType{}
	apiType.setType(TypeString + "[]")
	value, err = NewValueWithAPIType(apiType, nil)
	require.NoError(err)
	require.Equal(TypeArray, value.Type)
	require.NotNil(value.Array)
	require.Empty(value.Array)

	apiType = APIType{}
	err = yaml.Unmarshal([]byte(strings.TrimSpace(`
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
	`)), &apiType)
	require.NoError(err)
	value, err = NewValueWithAPIType(apiType, url.Values{
		"bool":  []string{"true"},
		"bools": []string{"false", "true"},
		"int":   []string{"9527"},
		"ints":  []string{"5566", "9527"},
		"num":   []string{"3.14"},
		"nums":  []string{"6.02", "3.14"},
		"str":   []string{"test string value"},
		"strs":  []string{"str1", "str2"},
	})
	require.NoError(err)
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
