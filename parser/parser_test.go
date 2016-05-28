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

	_, err = parser.ParseData([]byte("#%RAML 0.8\n"), ".")
	require.Error(err)
	require.True(ErrorUnexpectedRAMLVersion2.Match(err))
}

func Test_ParseAnnotationsSimpleAnnotations(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	rootdoc, err := parser.ParseFile("./raml-examples/annotations/simple-annotations.raml")
	require.NoError(err)
	require.NotZero(rootdoc)

	require.Equal("Illustrating annotations", rootdoc.Title)
	require.Equal("application/json", rootdoc.MediaType)
	if assert.Contains(rootdoc.AnnotationTypes, "testHarness") {
		annotationType := rootdoc.AnnotationTypes["testHarness"]
		require.Equal(typeString, annotationType.Type)
	}
	if assert.Contains(rootdoc.AnnotationTypes, "badge") {
		annotationType := rootdoc.AnnotationTypes["badge"]
		require.Equal(typeString, annotationType.Type)
	}
	if assert.Contains(rootdoc.AnnotationTypes, "clearanceLevel") {
		annotationType := rootdoc.AnnotationTypes["clearanceLevel"]
		require.Equal(typeObject, annotationType.Type)
		if assert.Contains(annotationType.Properties, "level") {
			property := annotationType.Properties["level"]
			require.Len(property.Enum, 3)
			require.True(property.Required)
		}
		if assert.Contains(annotationType.Properties, "signature") {
			property := annotationType.Properties["signature"]
			require.Equal("\\d{3}-\\w{12}", property.Pattern)
			require.True(property.Required)
		}
	}
	if assert.Contains(rootdoc.Resources, "/users") {
		resource := rootdoc.Resources["/users"]
		if assert.Contains(resource.Annotations, "(testHarness)") {
			annotation := resource.Annotations["(testHarness)"]
			require.Equal("usersTest", annotation.String)
		}
		if assert.Contains(resource.Annotations, "(badge)") {
			annotation := resource.Annotations["(badge)"]
			require.Equal("tested.gif", annotation.String)
		}
		if assert.Contains(resource.Annotations, "(clearanceLevel)") {
			annotation := resource.Annotations["(clearanceLevel)"]
			if assert.Contains(annotation.Map, "level") {
				value := annotation.Map["level"]
				require.Equal("high", value.String)
			}
			if assert.Contains(annotation.Map, "signature") {
				value := annotation.Map["signature"]
				require.Equal("230-ghtwvfrs1itr", value.String)
			}
		}
	}
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
		require.Equal("Bob", typ.Example.Value.Map["name"].String)
		require.Equal("Marley", typ.Example.Value.Map["lastname"].String)
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
			if assert.Contains(body.Example.Value.Map, "name") {
				name := body.Example.Value.Map["name"]
				require.Equal("Doe Enterprise", name.String)
			}
			if assert.Contains(body.Example.Value.Map, "value") {
				value := body.Example.Value.Map["value"]
				require.Equal("Silver", value.String)
			}
		}
		require.Equal("Returns an organisation entity.", resource.Get.Description)
		if assert.Contains(resource.Get.Responses, HTTPCode(201)) {
			response := resource.Get.Responses[201]
			if assert.Contains(response.Bodies.ForMIMEType, "application/json") {
				body := response.Bodies.ForMIMEType["application/json"]
				require.Equal("Org", body.Type)
				if assert.Contains(body.Examples, "acme") {
					example := body.Examples["acme"]
					require.Equal("Acme", example.Value.Map["name"].String)
				}
				if assert.Contains(body.Examples, "softwareCorp") {
					example := body.Examples["softwareCorp"]
					require.Equal("Software Corp", example.Value.Map["name"].String)
					require.Equal("35 Central Street", example.Value.Map["address"].String)
					require.Equal("Gold", example.Value.Map["value"].String)
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

func Test_ParseOthersMobileOrderApi(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	rootdoc, err := parser.ParseFile("./raml-examples/others/mobile-order-api/api.raml")
	require.NoError(err)
	require.NotZero(rootdoc)

	require.Equal("Mobile Order API", rootdoc.Title)
	require.Equal("1.0", rootdoc.Version)
	require.Equal("http://localhost:8081/api", rootdoc.BaseURI)
	if assert.Contains(rootdoc.Uses, "assets") {
		use := rootdoc.Uses["assets"]
		if assert.Contains(use.Types, "ProductItem") {
			typ := use.Types["ProductItem"]
			require.Equal(typeObject, typ.Type)
			if assert.Contains(typ.Properties, "product_id") {
				property := typ.Properties["product_id"]
				require.Equal(typeString, property.Type)
			}
			if assert.Contains(typ.Properties, "quantity") {
				property := typ.Properties["quantity"]
				require.Equal(typeInteger, property.Type)
			}
		}
		if assert.Contains(use.Types, "Order") {
			typ := use.Types["Order"]
			require.Equal(typeObject, typ.Type)
			if assert.Contains(typ.Properties, "order_id") {
				property := typ.Properties["order_id"]
				require.Equal(typeString, property.Type)
			}
			if assert.Contains(typ.Properties, "creation_date") {
				property := typ.Properties["creation_date"]
				require.Equal(typeString, property.Type)
			}
			if assert.Contains(typ.Properties, "items") {
				property := typ.Properties["items"]
				require.Equal("ProductItem[]", property.Type)
			}
		}
		if assert.Contains(use.Types, "Orders") {
			typ := use.Types["Orders"]
			require.Equal(typeObject, typ.Type)
			if assert.Contains(typ.Properties, "orders") {
				property := typ.Properties["orders"]
				require.Equal("Order[]", property.Type)
			}
		}
		if assert.Contains(use.Traits, "paging") {
			trait := use.Traits["paging"]
			if assert.Contains(trait.QueryParameters, "size") {
				qp := trait.QueryParameters["size"]
				require.Equal("the amount of elements of each result page", qp.Description)
				require.Equal(typeInteger, qp.Type)
				require.False(qp.Required)
				require.Equal(typeInteger, qp.Example.Value.Type)
				require.EqualValues(10, qp.Example.Value.Integer)
			}
			if assert.Contains(trait.QueryParameters, "page") {
				qp := trait.QueryParameters["page"]
				require.Equal("the page number", qp.Description)
				require.Equal(typeInteger, qp.Type)
				require.False(qp.Required)
				require.Equal(typeInteger, qp.Example.Value.Type)
				require.EqualValues(0, qp.Example.Value.Integer)
			}
		}
	}
	if assert.Contains(rootdoc.Resources, "/orders") {
		resource := rootdoc.Resources["/orders"]
		require.Equal("Orders", resource.DisplayName)
		require.Equal("Orders collection resource used to create new orders.", resource.Description)
		if assert.Len(resource.Get.Is, 1) {
			is := resource.Get.Is[0]
			require.Equal("assets.paging", is.String)
		}
		require.Equal("lists all orders of a specific user", resource.Get.Description)
		if assert.Contains(resource.Get.QueryParameters, "userId") {
			qp := resource.Get.QueryParameters["userId"]
			require.Equal("string", qp.Type)
			require.Equal("use to query all orders of a user", qp.Description)
			require.True(qp.Required)
			require.Equal("1964401a-a8b3-40c1-b86e-d8b9f75b5842", qp.Example.Value.String)
		}
		if assert.Contains(resource.Get.Responses, HTTPCode(200)) {
			response := resource.Get.Responses[200]
			if assert.Contains(response.Bodies.ForMIMEType, "application/json") {
				body := response.Bodies.ForMIMEType["application/json"]
				require.Equal("assets.Orders", body.Type)
				if assert.Contains(body.Examples, "single-order") {
					example := body.Examples["single-order"]
					if assert.Contains(example.Value.Map, "orders") {
						orders := example.Value.Map["orders"]
						if assert.Len(orders.Array, 1) {
							order := orders.Array[0]
							if assert.Contains(order.Map, "order_id") {
								orderID := order.Map["order_id"]
								require.Equal("ORDER-437563756", orderID.String)
							}
							if assert.Contains(order.Map, "creation_date") {
								creationDate := order.Map["creation_date"]
								require.Equal("2016-03-30", creationDate.String)
							}
							if assert.Contains(order.Map, "items") {
								items := order.Map["items"]
								if assert.Len(items.Array, 2) {
									item := items.Array[1]
									if assert.Contains(item.Map, "product_id") {
										productID := item.Map["product_id"]
										require.Equal("PRODUCT-2", productID.String)
									}
									if assert.Contains(item.Map, "quantity") {
										quantity := item.Map["quantity"]
										require.EqualValues(2, quantity.Integer)
									}
								}
							}
						}
					}
				}
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
