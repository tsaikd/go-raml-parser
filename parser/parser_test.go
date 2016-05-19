package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tsaikd/go-raml-parser/parser/parserConfig"
)

func Test_Parse(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	err := parser.Config(0, nil)
	require.Error(err)
	require.True(ErrorUnsupportedParserConfig1.Match(err))

	err = parser.Config(parserConfig.CheckRAMLVersion, nil)
	require.Error(err)
	require.True(ErrorInvaludParserConfigValueType3.Match(err))

	err = parser.Config(parserConfig.CheckRAMLVersion, true)
	require.NoError(err)

	_, err = parser.ParseData([]byte("#%RAML 0.8\n"))
	require.Error(err)
	require.True(ErrorUnexpectedRAMLVersion2.Match(err))
}

func Test_ParseDefiningExamples(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	rootdoc, err := parser.ParseFile("./raml-examples/defining-examples/organisation-api.raml")
	require.NoError(err)
	require.NotZero(rootdoc)

	require.Equal(rootdoc.Title, "API with Examples")
	if assert.Contains(rootdoc.Types, "User") {
		typ := rootdoc.Types["User"]
		require.Equal(typeObject, typ.Type)
		if assert.Contains(typ.Properties, "name") {
			property := typ.Properties["name"]
			require.Equal(typeString, property.Type)
		}
		if assert.Contains(typ.Properties, "lastname") {
			property := typ.Properties["lastname"]
			require.Equal(typeString, property.Type)
		}
		require.False(typ.Example.Value.IsEmpty())
		require.Equal("Bob", typ.Example.Value.Map["name"])
		require.Equal("Marley", typ.Example.Value.Map["lastname"])
	}
	if assert.Contains(rootdoc.Types, "Org") {
		typ := rootdoc.Types["Org"]
		require.Equal(typeObject, typ.Type)
		if assert.Contains(typ.Properties, "name") {
			property := typ.Properties["name"]
			require.Equal(typeString, property.Type)
			require.True(property.Required)
		}
		if assert.Contains(typ.Properties, "address") {
			property := typ.Properties["address"]
			require.Equal(typeString, property.Type)
			require.False(property.Required)
		}
		if assert.Contains(typ.Properties, "value") {
			property := typ.Properties["value"]
			require.Equal(typeString, property.Type)
			require.False(property.Required)
		}
	}
	if assert.Contains(rootdoc.Resources, "/organisation") {
		resource := rootdoc.Resources["/organisation"]
		if assert.Contains(resource.Post.Headers, "UserID") {
			header := resource.Post.Headers["UserID"]
			require.Equal("the identifier for the user that posts a new organisation", header.Description)
			require.Equal(typeString, header.Type)
			require.Equal("SWED-123", header.Example.Value.String)
		}
		if assert.Contains(resource.Post.Bodies.ForMIMEType, "application/json") {
			body := resource.Post.Bodies.ForMIMEType["application/json"]
			require.Equal("Org", body.Type)
			require.Equal("Doe Enterprise", body.Example.Value.Map["name"])
			require.Equal("Silver", body.Example.Value.Map["value"])
		}
		require.Equal("Returns an organisation entity.", resource.Get.Description)
		if assert.Contains(resource.Get.Responses, HTTPCode(201)) {
			response := resource.Get.Responses[201]
			if assert.Contains(response.Bodies.ForMIMEType, "application/json") {
				body := response.Bodies.ForMIMEType["application/json"]
				require.Equal("Org", body.Type)
				if assert.Contains(body.Examples, "acme") {
					example := body.Examples["acme"]
					require.Equal("Acme", example.Value.Map["name"])
				}
				if assert.Contains(body.Examples, "softwareCorp") {
					example := body.Examples["softwareCorp"]
					require.Equal("Software Corp", example.Value.Map["name"])
					require.Equal("35 Central Street", example.Value.Map["address"])
					require.Equal("Gold", example.Value.Map["value"])
				}
			}
		}
	}
}

func Test_ParseHelloworld(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	rootdoc, err := parser.ParseFile("./raml-examples/helloworld/helloworld.raml")
	require.NoError(err)
	require.NotZero(rootdoc)

	require.Equal("Hello world", rootdoc.Title)
	if assert.Contains(rootdoc.Resources, "/helloworld") {
		resource := rootdoc.Resources["/helloworld"]
		if assert.Contains(resource.Get.Responses, HTTPCode(200)) {
			response := resource.Get.Responses[200]
			if assert.Contains(response.Bodies.ForMIMEType, "application/json") {
				body := response.Bodies.ForMIMEType["application/json"]
				require.NotEmpty(body.Type)
				require.NotEmpty(body.Example)
			}
		}
	}
}

func Test_ParseTypesystemSimple(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	rootdoc, err := parser.ParseFile("./raml-examples/typesystem/simple.raml")
	require.NoError(err)
	require.NotZero(rootdoc)

	require.Equal("API with Types", rootdoc.Title)
	if assert.Contains(rootdoc.Types, "User") {
		typ := rootdoc.Types["User"]
		require.Equal(typeObject, typ.Type)
		if assert.Contains(typ.Properties, "age") {
			property := typ.Properties["age"]
			require.True(property.Required)
			require.Equal(typeNumber, property.Type)
		}
		if assert.Contains(typ.Properties, "firstName") {
			property := typ.Properties["firstName"]
			require.True(property.Required)
			require.Equal(typeString, property.Type)
		}
		if assert.Contains(typ.Properties, "lastName") {
			property := typ.Properties["lastName"]
			require.True(property.Required)
			require.Equal(typeString, property.Type)
		}
	}
	if assert.Contains(rootdoc.Resources, "/users/{id}") {
		resource := rootdoc.Resources["/users/{id}"]
		if assert.Contains(resource.Get.Responses, HTTPCode(200)) {
			response := resource.Get.Responses[200]
			if assert.Contains(response.Bodies.ForMIMEType, "application/json") {
				body := response.Bodies.ForMIMEType["application/json"]
				require.Equal(body.Type, "User")
			}
		}
	}
}
