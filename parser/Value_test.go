package parser

import (
	"net/url"
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

	value, err = NewValue((*url.Values)(nil))
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

func Test_ValueFromSlice(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	value, err := NewValue([]bool{true, false})
	require.NoError(err)
	require.Equal(TypeArray, value.Type)

	value, err = NewValue([]int{1, 2, 3})
	require.NoError(err)
	require.Equal(TypeArray, value.Type)

	value, err = NewValue([]float32{3.14})
	require.NoError(err)
	require.Equal(TypeArray, value.Type)

	value, err = NewValue([]string{"text"})
	require.NoError(err)
	require.Equal(TypeArray, value.Type)
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
	if prop := value.Map["key1"]; assert.NotNil(prop) {
		require.Equal(TypeBoolean, prop.Type)
		require.True(prop.Boolean)
	}
	if prop := value.Map["key2"]; assert.NotNil(prop) {
		require.Equal(TypeInteger, prop.Type)
		require.EqualValues(9527, prop.Integer)
	}
	if prop := value.Map["key3"]; assert.NotNil(prop) {
		require.Equal(TypeNumber, prop.Type)
		require.EqualValues(3.14, prop.Number)
	}
	if prop := value.Map["key4"]; assert.NotNil(prop) {
		require.Equal(TypeString, prop.Type)
		require.EqualValues("test", prop.String)
	}
	if propRoot := value.Map["child"]; assert.NotNil(propRoot) {
		require.Equal(TypeObject, propRoot.Type)
		if propChild := propRoot.Map["childKey"]; assert.NotNil(propChild) {
			require.Equal(TypeString, propChild.Type)
			require.EqualValues("child value", propChild.String)
		}
	}
}

func Test_ValueFromURLValues(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	value, err := NewValue(url.Values{
		"key1": []string{"single"},
		"key2": []string{"foo", "bar"},
	})
	require.NoError(err)
	require.Equal(TypeObject, value.Type)
	if prop := value.Map["key1"]; assert.NotNil(prop) {
		require.Equal(TypeString, prop.Type)
		require.Equal("single", prop.String)
	}
	if prop := value.Map["key2"]; assert.NotNil(prop) {
		require.Equal(TypeArray, prop.Type)
		require.Len(prop.Array, 2)
		if elem := prop.Array[0]; assert.NotNil(elem) {
			require.Equal(TypeString, elem.Type)
			require.Equal("foo", elem.String)
		}
		if elem := prop.Array[1]; assert.NotNil(elem) {
			require.Equal(TypeString, elem.Type)
			require.Equal("bar", elem.String)
		}
	}
}
