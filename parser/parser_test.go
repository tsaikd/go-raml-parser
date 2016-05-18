package parser

import (
	"testing"

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

func Test_ParseHelloworld(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	rootdoc, err := parser.ParseFile("./raml-examples/helloworld/helloworld.raml")
	require.NoError(err)
	require.NotZero(rootdoc)

	require.Equal(rootdoc.Title, "Hello world")
	require.Contains(rootdoc.Resources, "/helloworld")
	resource := rootdoc.Resources["/helloworld"]
	require.Contains(resource.Get.Responses, HTTPCode(200))
	resourceCode := resource.Get.Responses[200]
	require.Contains(resourceCode.Bodies.ForMIMEType, "application/json")
	body := resourceCode.Bodies.ForMIMEType["application/json"]
	require.NotEmpty(body.Type)
	require.NotEmpty(body.Example)
}

func Test_ParseTypesystemSimple(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	rootdoc, err := parser.ParseFile("./raml-examples/typesystem/simple.raml")
	require.NoError(err)
	require.NotZero(rootdoc)

	require.Equal(rootdoc.Title, "API with Types")
	require.Contains(rootdoc.Types, "User")
	typeUser := rootdoc.Types["User"]
	require.Equal(typeUser.Type, "object")
	require.Contains(typeUser.Properties, "age")
	propertyAge := typeUser.Properties["age"]
	require.True(propertyAge.Required)
	require.Equal(propertyAge.Type, "number")
	require.Contains(typeUser.Properties, "firstName")
	propertyFirstName := typeUser.Properties["firstName"]
	require.True(propertyFirstName.Required)
	require.Equal(propertyFirstName.Type, "string")
	require.Contains(typeUser.Properties, "lastName")
	propertyLastName := typeUser.Properties["lastName"]
	require.True(propertyLastName.Required)
	require.Equal(propertyLastName.Type, "string")
	require.Contains(rootdoc.Resources, "/users/{id}")
	resourceUsersID := rootdoc.Resources["/users/{id}"]
	require.Contains(resourceUsersID.Get.Responses, HTTPCode(200))
	resourceUsersIDBody := resourceUsersID.Get.Responses[200]
	require.Contains(resourceUsersIDBody.Bodies.ForMIMEType, "application/json")
	body := resourceUsersIDBody.Bodies.ForMIMEType["application/json"]
	require.Equal(body.Type, "User")
}
