package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ValueError(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	_, err := NewValue(struct{}{})
	require.Error(err)
	require.True(ErrorUnsupportedValueType1.Match(err))
}

func Test_ValueFromValue(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	src, err := NewValue("test")
	require.NoError(err)
	require.Equal(TypeString, src.Type)
	require.Equal("test", src.String)

	value, err := NewValue(src)
	require.NoError(err)
	require.Equal(src.Type, value.Type)
	require.Equal(src.String, value.String)

	value, err = NewValue(&src)
	require.NoError(err)
	require.Equal(src.Type, value.Type)
	require.Equal(src.String, value.String)

	value, err = NewValue([]*Value{&src, &src})
	require.NoError(err)
	require.Equal(TypeArray, value.Type)

	value, err = NewValue([]Value{src, src})
	require.NoError(err)
	require.Equal(TypeArray, value.Type)
}

func Test_ValueFromNull(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	value, err := NewValue(nil)
	require.NoError(err)
	require.Equal(TypeNull, value.Type)
}

func Test_ValueFromBool(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	value, err := NewValue(true)
	require.NoError(err)
	require.Equal(TypeBoolean, value.Type)
	require.True(value.Boolean)
}

func Test_ValueFromInt(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	value, err := NewValue(int(9527))
	require.NoError(err)
	require.Equal(TypeInteger, value.Type)
	require.EqualValues(9527, value.Integer)

	value, err = NewValue(int8(97))
	require.NoError(err)
	require.Equal(TypeInteger, value.Type)
	require.EqualValues(97, value.Integer)

	value, err = NewValue(int16(9527))
	require.NoError(err)
	require.Equal(TypeInteger, value.Type)
	require.EqualValues(9527, value.Integer)

	value, err = NewValue(int32(9527))
	require.NoError(err)
	require.Equal(TypeInteger, value.Type)
	require.EqualValues(9527, value.Integer)

	value, err = NewValue(int64(9527))
	require.NoError(err)
	require.Equal(TypeInteger, value.Type)
	require.EqualValues(9527, value.Integer)
}

func Test_ValueFromFloat(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	value, err := NewValue(float32(3.14))
	require.NoError(err)
	require.Equal(TypeNumber, value.Type)
	require.EqualValues(float32(3.14), value.Number)

	value, err = NewValue(float64(3.14))
	require.NoError(err)
	require.Equal(TypeNumber, value.Type)
	require.EqualValues(float64(3.14), value.Number)
}

func Test_ValueFromBinary(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	value, err := NewValue([]byte("abc"))
	require.NoError(err)
	require.Equal(TypeBinary, value.Type)
	require.Len(value.Binary, 3)
}

func Test_ValueFromInterfaceSlice(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	value, err := NewValue([]interface{}{
		true,
		int(9527),
		float64(3.14),
		"test",
	})
	require.NoError(err)
	require.Equal(TypeArray, value.Type)
	require.Len(value.Array, 4)
	require.NotNil(value.Array[0])
	require.Equal(TypeBoolean, value.Array[0].Type)
	require.True(value.Array[0].Boolean)
	require.NotNil(value.Array[1])
	require.Equal(TypeInteger, value.Array[1].Type)
	require.Equal(int64(9527), value.Array[1].Integer)
	require.NotNil(value.Array[2])
	require.Equal(TypeNumber, value.Array[2].Type)
	require.Equal(float64(3.14), value.Array[2].Number)
	require.NotNil(value.Array[3])
	require.Equal(TypeString, value.Array[3].Type)
	require.Equal("test", value.Array[3].String)
}

func Test_ValueFromMap(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	value, err := NewValue(map[string]interface{}{
		"key1": true,
		"key2": 9527,
		"key3": 3.14,
		"key4": "test",
		"child": map[string]interface{}{
			"childKey": "child value",
		},
	})
	require.NoError(err)
	require.Equal(TypeObject, value.Type)
	if assert.Contains(value.Map, "key1") {
		rootval := value.Map["key1"]
		require.NotNil(rootval)
		require.Equal(TypeBoolean, rootval.Type)
		require.True(rootval.Boolean)
	}
	if assert.Contains(value.Map, "key2") {
		rootval := value.Map["key2"]
		require.NotNil(rootval)
		require.Equal(TypeInteger, rootval.Type)
		require.EqualValues(9527, rootval.Integer)
	}
	if assert.Contains(value.Map, "key3") {
		rootval := value.Map["key3"]
		require.NotNil(rootval)
		require.Equal(TypeNumber, rootval.Type)
		require.EqualValues(3.14, rootval.Number)
	}
	if assert.Contains(value.Map, "key4") {
		rootval := value.Map["key4"]
		require.NotNil(rootval)
		require.Equal(TypeString, rootval.Type)
		require.EqualValues("test", rootval.String)
	}
	if assert.Contains(value.Map, "child") {
		rootval := value.Map["child"]
		require.NotNil(rootval)
		require.Equal(TypeObject, rootval.Type)
		if assert.Contains(rootval.Map, "childKey") {
			childval := rootval.Map["childKey"]
			require.NotNil(childval)
			require.Equal(TypeString, childval.Type)
			require.EqualValues("child value", childval.String)
		}
	}
}
