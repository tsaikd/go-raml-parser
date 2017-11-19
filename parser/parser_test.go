package parser

import (
	"bytes"
	"encoding/gob"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tsaikd/KDGoLib/jsonex"
	"github.com/tsaikd/KDGoLib/testutil/requireutil"
	"github.com/tsaikd/go-raml-parser/parser/parserConfig"
)

func Test_ParseError(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	err := parser.Config(0, nil)
	require.Error(err)
	require.True(ErrorUnsupportedParserConfig1.Match(err))

	err = parser.Config(parserConfig.CheckRAMLVersion, nil)
	require.Error(err)
	require.True(ErrorInvalidParserConfigValueType3.Match(err))

	err = parser.Config(parserConfig.CheckRAMLVersion, true)
	require.NoError(err)

	err = parser.Config(parserConfig.CheckValueOptions, nil)
	require.Error(err)

	err = parser.Config(parserConfig.CheckValueOptions, []CheckValueOption{CheckValueOptionAllowIntegerToBeNumber(true)})
	require.NoError(err)

	err = parser.Config(parserConfig.CheckValueOptions, "error")
	require.Error(err)

	_, err = parser.ParseData([]byte("#%RAML 0.8\n"), ".")
	require.Error(err)
	require.True(ErrorUnexpectedRAMLVersion2.Match(err))

	_, err = parser.ParseData([]byte(strings.TrimSpace(`
#%RAML 1.0

/get/error:
    get:
        response:
            200:
                body:
                    application/json:
                        type: string

	`)), ".")
	require.Error(err)
	require.True(ErrorTypo2.Match(err))
}

func Test_ParseAnnotationsAnnotationTargets(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	rootdoc, err := parser.ParseFile("./raml-examples/annotations/annotation-targets.raml")
	require.NoError(err)
	require.NotZero(rootdoc)

	require.Equal("Illustrating allowed targets", rootdoc.Title)
	require.Equal("application/json", rootdoc.MediaType)
	if annotationType, ok := rootdoc.AnnotationTypes["meta-resource-method"]; assert.True(ok) {
		if assert.Len(annotationType.AllowedTargets, 2) {
			require.Equal(TargetLocationResource, annotationType.AllowedTargets[0])
			require.Equal(TargetLocationMethod, annotationType.AllowedTargets[1])
		}
	}
	if annotationType, ok := rootdoc.AnnotationTypes["meta-data"]; assert.True(ok) {
		if assert.Len(annotationType.AllowedTargets, 1) {
			require.Equal(TargetLocationTypeDeclaration, annotationType.AllowedTargets[0])
		}
	}
	if apiType, ok := rootdoc.Types["User"]; assert.True(ok) {
		require.Equal(TypeObject, apiType.Type)
		if annotation, ok := apiType.Annotations["meta-data"]; assert.True(ok) {
			require.Equal("on an object; on a data type declaration", annotation.String)
		}
		if property, ok := apiType.Properties.Map()["name"]; assert.True(ok) {
			require.Equal(TypeString, property.Type)
			if annotation, ok := property.Annotations["meta-data"]; assert.True(ok) {
				require.Equal("on a string property", annotation.String)
			}
		}
	}
	if resource, ok := rootdoc.Resources["/users"]; assert.True(ok) {
		if annotation, ok := resource.Annotations["meta-resource-method"]; assert.True(ok) {
			require.Equal("on a resource", annotation.String)
		}
		if method, ok := resource.Methods["get"]; assert.True(ok) {
			if annotation, ok := method.Annotations["meta-resource-method"]; assert.True(ok) {
				require.Equal("on a method", annotation.String)
			}
			if response, ok := method.Responses[HTTPCode(200)]; assert.True(ok) {
				if body, ok := response.Bodies["application/json"]; assert.True(ok) {
					require.Equal("User[]", body.Type)
					if annotation, ok := body.Annotations["meta-data"]; assert.True(ok) {
						require.Equal("on a body", annotation.String)
					}
				}
			}
		}
	}
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
	if annotationType, ok := rootdoc.AnnotationTypes["testHarness"]; assert.True(ok) {
		require.Equal(TypeString, annotationType.Type)
	}
	if annotationType, ok := rootdoc.AnnotationTypes["badge"]; assert.True(ok) {
		require.Equal(TypeString, annotationType.Type)
	}
	if annotationType, ok := rootdoc.AnnotationTypes["clearanceLevel"]; assert.True(ok) {
		require.Equal(TypeObject, annotationType.Type)
		if property, ok := annotationType.Properties.Map()["level"]; assert.True(ok) {
			require.Len(property.Enum, 3)
			require.True(property.Required)
		}
		if property, ok := annotationType.Properties.Map()["signature"]; assert.True(ok) {
			require.Equal("\\d{3}-\\w{12}", property.Pattern)
			require.True(property.Required)
		}
	}
	if resource, ok := rootdoc.Resources["/users"]; assert.True(ok) {
		if annotation, ok := resource.Annotations["testHarness"]; assert.True(ok) {
			require.Equal("usersTest", annotation.String)
		}
		if annotation, ok := resource.Annotations["badge"]; assert.True(ok) {
			require.Equal("tested.gif", annotation.String)
		}
		if annotation, ok := resource.Annotations["clearanceLevel"]; assert.True(ok) {
			if value, ok := annotation.Map["level"]; assert.True(ok) {
				require.Equal("high", value.String)
			}
			if value, ok := annotation.Map["signature"]; assert.True(ok) {
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

	require.Equal("API with Examples", rootdoc.Title)
	if typ, ok := rootdoc.Types["User"]; assert.True(ok) {
		require.Equal(TypeObject, typ.Type)
		if property, ok := typ.Properties.Map()["name"]; assert.True(ok) {
			require.Equal(TypeString, property.Type)
		}
		if property, ok := typ.Properties.Map()["lastname"]; assert.True(ok) {
			require.Equal(TypeString, property.Type)
		}
		require.False(typ.Example.Value.IsEmpty())
		require.Equal("Bob", typ.Example.Value.Map["name"].String)
		require.Equal("Marley", typ.Example.Value.Map["lastname"].String)
	}
	if typ, ok := rootdoc.Types["Org"]; assert.True(ok) {
		require.Equal(TypeObject, typ.Type)
		if property, ok := typ.Properties.Map()["name"]; assert.True(ok) {
			require.Equal(TypeString, property.Type)
			require.True(property.Required)
		}
		if property, ok := typ.Properties.Map()["address"]; assert.True(ok) {
			require.Equal(TypeString, property.Type)
			require.False(property.Required)
		}
		if property, ok := typ.Properties.Map()["value"]; assert.True(ok) {
			require.Equal(TypeString, property.Type)
			require.False(property.Required)
		}
	}
	if resource, ok := rootdoc.Resources["/organisation"]; assert.True(ok) {
		if method, ok := resource.Methods["post"]; assert.True(ok) {
			if header, ok := method.Headers.Map()["UserID"]; assert.True(ok) {
				require.Equal("the identifier for the user that posts a new organisation", header.Description)
				require.Equal(TypeString, header.Type)
				require.Equal("SWED-123", header.Example.Value.String)
			}
			if body, ok := method.Bodies["application/json"]; assert.True(ok) {
				require.Equal("Org", body.Type)
				if name, ok := body.Example.Value.Map["name"]; assert.True(ok) {
					require.Equal("Doe Enterprise", name.String)
				}
				if value, ok := body.Example.Value.Map["value"]; assert.True(ok) {
					require.Equal("Silver", value.String)
				}
			}
		}
		if method, ok := resource.Methods["get"]; assert.True(ok) {
			require.Equal("Returns an organisation entity.", method.Description)
			if response, ok := method.Responses[201]; assert.True(ok) {
				if body, ok := response.Bodies["application/json"]; assert.True(ok) {
					require.Equal("Org", body.Type)
					if example, ok := body.Examples["acme"]; assert.True(ok) {
						require.Equal("Acme", example.Value.Map["name"].String)
					}
					if example, ok := body.Examples["softwareCorp"]; assert.True(ok) {
						require.Equal("Software Corp", example.Value.Map["name"].String)
						require.Equal("35 Central Street", example.Value.Map["address"].String)
						require.Equal("Gold", example.Value.Map["value"].String)
					}
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
	if resource, ok := rootdoc.Resources["/helloworld"]; assert.True(ok) {
		if method, ok := resource.Methods["get"]; assert.True(ok) {
			if response, ok := method.Responses[200]; assert.True(ok) {
				if body, ok := response.Bodies["application/json"]; assert.True(ok) {
					require.NotEmpty(body.Type)
					require.NotEmpty(body.Example)
				}
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
	if use, ok := rootdoc.Uses["assets"]; assert.True(ok) {
		if typ, ok := use.Types["ProductItem"]; assert.True(ok) {
			require.Equal(TypeObject, typ.Type)
			if property, ok := typ.Properties.Map()["product_id"]; assert.True(ok) {
				require.Equal(TypeString, property.Type)
			}
			if property, ok := typ.Properties.Map()["quantity"]; assert.True(ok) {
				require.Equal(TypeInteger, property.Type)
			}
		}
		if typ, ok := use.Types["Order"]; assert.True(ok) {
			require.Equal(TypeObject, typ.Type)
			if property, ok := typ.Properties.Map()["order_id"]; assert.True(ok) {
				require.Equal(TypeString, property.Type)
			}
			if property, ok := typ.Properties.Map()["creation_date"]; assert.True(ok) {
				require.Equal(TypeString, property.Type)
			}
			if property, ok := typ.Properties.Map()["items"]; assert.True(ok) {
				require.Equal("ProductItem[]", property.Type)
			}
		}
		if typ, ok := use.Types["Orders"]; assert.True(ok) {
			require.Equal(TypeObject, typ.Type)
			if property, ok := typ.Properties.Map()["orders"]; assert.True(ok) {
				require.Equal("Order[]", property.Type)
			}
		}
		if trait, ok := use.Traits["paging"]; assert.True(ok) {
			if qp, ok := trait.QueryParameters.Map()["size"]; assert.True(ok) {
				require.Equal("the amount of elements of each result page", qp.Description)
				require.Equal(TypeInteger, qp.Type)
				require.False(qp.Required)
				require.Equal(TypeInteger, qp.Example.Value.Type)
				require.EqualValues(10, qp.Example.Value.Integer)
			}
			if qp, ok := trait.QueryParameters.Map()["page"]; assert.True(ok) {
				require.Equal("the page number", qp.Description)
				require.Equal(TypeInteger, qp.Type)
				require.False(qp.Required)
				require.Equal(TypeInteger, qp.Example.Value.Type)
				require.EqualValues(0, qp.Example.Value.Integer)
			}
		}
	}
	if resource, ok := rootdoc.Resources["/orders"]; assert.True(ok) {
		require.Equal("Orders", resource.DisplayName)
		require.Equal("Orders collection resource used to create new orders.", resource.Description)
		if method, ok := resource.Methods["get"]; assert.True(ok) {
			if assert.Len(method.Is, 1) {
				is := method.Is[0]
				require.Equal("assets.paging", is.String)
			}
			require.Equal("lists all orders of a specific user", method.Description)
			if qp, ok := method.QueryParameters.Map()["userId"]; assert.True(ok) {
				require.Equal("string", qp.Type)
				require.Equal("use to query all orders of a user", qp.Description)
				require.True(qp.Required)
				require.Equal("1964401a-a8b3-40c1-b86e-d8b9f75b5842", qp.Example.Value.String)
			}
			if response, ok := method.Responses[200]; assert.True(ok) {
				if body, ok := response.Bodies["application/json"]; assert.True(ok) {
					require.Equal("assets.Orders", body.Type)
					if example, ok := body.Examples["single-order"]; assert.True(ok) {
						if orders, ok := example.Value.Map["orders"]; assert.True(ok) {
							if assert.Len(orders.Array, 1) {
								order := orders.Array[0]
								if orderID, ok := order.Map["order_id"]; assert.True(ok) {
									require.Equal("ORDER-437563756", orderID.String)
								}
								if creationDate, ok := order.Map["creation_date"]; assert.True(ok) {
									require.Equal("2016-03-30", creationDate.String)
								}
								if items, ok := order.Map["items"]; assert.True(ok) {
									if assert.Len(items.Array, 2) {
										item := items.Array[1]
										if productID, ok := item.Map["product_id"]; assert.True(ok) {
											require.Equal("PRODUCT-2", productID.String)
										}
										if quantity, ok := item.Map["quantity"]; assert.True(ok) {
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
	if typ, ok := rootdoc.Types["User"]; assert.True(ok) {
		require.Equal(TypeObject, typ.Type)
		if property, ok := typ.Properties.Map()["age"]; assert.True(ok) {
			require.True(property.Required)
			require.Equal(TypeNumber, property.Type)
		}
		if property, ok := typ.Properties.Map()["firstName"]; assert.True(ok) {
			require.True(property.Required)
			require.Equal(TypeString, property.Type)
		}
		if property, ok := typ.Properties.Map()["lastName"]; assert.True(ok) {
			require.True(property.Required)
			require.Equal(TypeString, property.Type)
		}
	}
	if resource, ok := rootdoc.Resources["/users/{id}"]; assert.True(ok) {
		require.Contains(resource.URIParameters, "id")
		if method, ok := resource.Methods["get"]; assert.True(ok) {
			if response, ok := method.Responses[200]; assert.True(ok) {
				if body, ok := response.Bodies["application/json"]; assert.True(ok) {
					require.Equal("User", body.Type)
				}
			}
		}
	}
}

func Test_ParseAnnotationOnType(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	rootdoc, err := parser.ParseFile("./test-examples/annotation-on-type.raml")
	require.NoError(err)

	if annotationType, ok := rootdoc.AnnotationTypes["AnnotationOnType"]; assert.True(ok) {
		require.Equal("annotation on type", annotationType.Description)
		require.Len(annotationType.AllowedTargets, 1)
		require.Equal(TargetLocationTypeDeclaration, annotationType.AllowedTargets[0])
		require.Equal(TypeString, annotationType.Type)
	}
	if apiType, ok := rootdoc.Types["User"]; assert.True(ok) {
		if annotation, ok := apiType.Annotations["AnnotationOnType"]; assert.True(ok) {
			require.Equal("something on annotation", annotation.String)
			annotationType := annotation.AnnotationType
			require.Equal("annotation on type", annotationType.Description)
			require.Len(annotationType.AllowedTargets, 1)
			require.Equal(TargetLocationTypeDeclaration, annotationType.AllowedTargets[0])
			require.Equal(TypeString, annotationType.Type)
		}
	}
	if resource, ok := rootdoc.Resources["/user"]; assert.True(ok) {
		if method, ok := resource.Methods["get"]; assert.True(ok) {
			if body, ok := method.Bodies["application/json"]; assert.True(ok) {
				require.Equal("User", body.Type)
				if annotation, ok := body.Annotations["AnnotationOnType"]; assert.True(ok) {
					require.Equal("something on annotation", annotation.String)
					annotationType := annotation.AnnotationType
					require.Equal("annotation on type", annotationType.Description)
					require.Len(annotationType.AllowedTargets, 1)
					require.Equal(TargetLocationTypeDeclaration, annotationType.AllowedTargets[0])
					require.Equal(TypeString, annotationType.Type)
				}
			}
			if response, ok := method.Responses[HTTPCode(200)]; assert.True(ok) {
				if body, ok := response.Bodies["application/json"]; assert.True(ok) {
					require.Equal("User", body.Type)
					if annotation, ok := body.Annotations["AnnotationOnType"]; assert.True(ok) {
						require.Equal("something on annotation", annotation.String)
						annotationType := annotation.AnnotationType
						require.Equal("annotation on type", annotationType.Description)
						require.Len(annotationType.AllowedTargets, 1)
						require.Equal(TargetLocationTypeDeclaration, annotationType.AllowedTargets[0])
						require.Equal(TypeString, annotationType.Type)
					}
				}
			}
		}
	}
}

func Test_ParseBaseURIParameters(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	rootdoc, err := parser.ParseFile("./test-examples/base-uri-parameters.raml")
	require.NoError(err)

	require.Equal("Amazon S3 REST API", rootdoc.Title)
	require.Equal("1", rootdoc.Version)
	require.Equal("https://{bucketName}.s3.amazonaws.com", rootdoc.BaseURI)
	if assert.NotNil(rootdoc.BaseURIParameters) {
		if uriParam := rootdoc.BaseURIParameters["bucketName"]; assert.NotNil(uriParam) {
			require.Equal("The name of the bucket", uriParam.Description)
		}
	}
}

func Test_ParseCheckUnusedAnnotation(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	_, err := parser.ParseFile("./test-examples/check-unused-annotation.raml")
	require.Error(err)

	err = parser.Config(parserConfig.IgnoreUnusedAnnotation, true)
	require.NoError(err)

	rootdoc, err := parser.ParseFile("./test-examples/check-unused-annotation.raml")
	require.NoError(err)

	if resource, ok := rootdoc.Resources["/get"]; assert.True(ok) {
		if annotation, ok := resource.Annotations["UsedAnnotation"]; assert.True(ok) {
			require.Equal("used annotation", annotation.AnnotationType.Description)
		}
	}
}

func Test_ParseCheckUnusedTrait(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	_, err := parser.ParseFile("./test-examples/check-unused-trait.raml")
	require.Error(err)

	err = parser.Config(parserConfig.IgnoreUnusedTrait, true)
	require.NoError(err)

	_, err = parser.ParseFile("./test-examples/check-unused-trait.raml")
	require.NoError(err)
}

func Test_ParseDefaultMediaType(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	rootdoc, err := parser.ParseFile("./test-examples/default-mediaType.raml")
	require.NoError(err)
	if resource, ok := rootdoc.Resources["/user"]; assert.True(ok) {
		if method, ok := resource.Methods["get"]; assert.True(ok) {
			if response, ok := method.Responses[HTTPCode(200)]; assert.True(ok) {
				if body, ok := response.Bodies["application/json"]; assert.True(ok) {
					require.Equal(TypeObject, body.Example.Value.Type)
					if name, ok := body.Example.Value.Map["name"]; assert.True(ok) {
						require.Equal(TypeString, name.Type)
						require.Equal("Alice", name.String)
					}
				}
			}
		}
	}
}

func Test_ParseExampleFromType(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	rootdoc, err := parser.ParseFile("./test-examples/example-from-type.raml")
	require.NoError(err)
	require.NotZero(rootdoc)

	require.Equal("Example from type", rootdoc.Title)
	if typ, ok := rootdoc.Types["User"]; assert.True(ok) {
		require.Equal(TypeObject, typ.Type)
		if property, ok := typ.Properties.Map()["name"]; assert.True(ok) {
			require.True(property.Required)
			require.Equal(TypeString, property.Type)
		}
		if property, ok := typ.Properties.Map()["email"]; assert.True(ok) {
			require.True(property.Required)
			require.Equal(TypeString, property.Type)
		}
		if example, ok := typ.Examples["user1"]; assert.True(ok) {
			if value, ok := example.Value.Map["name"]; assert.True(ok) {
				require.Equal("Alice", value.String)
			}
			if value, ok := example.Value.Map["email"]; assert.True(ok) {
				require.Equal("alice@example.com", value.String)
			}
		}
		if example, ok := typ.Examples["user2"]; assert.True(ok) {
			if value, ok := example.Value.Map["name"]; assert.True(ok) {
				require.Equal("Bob", value.String)
			}
			if value, ok := example.Value.Map["email"]; assert.True(ok) {
				require.Equal("bob@example.com", value.String)
			}
		}
	}
	if resource, ok := rootdoc.Resources["/user"]; assert.True(ok) {
		if method, ok := resource.Methods["get"]; assert.True(ok) {
			if response, ok := method.Responses[200]; assert.True(ok) {
				if body, ok := response.Bodies["application/json"]; assert.True(ok) {
					require.Equal("User", body.Type)
					if value, ok := body.Example.Value.Map["name"]; assert.True(ok) {
						require.NotEmpty(value.String)
					}
				}
			}
		}
	}
	if resource, ok := rootdoc.Resources["/user/wrap"]; assert.True(ok) {
		if method, ok := resource.Methods["get"]; assert.True(ok) {
			if response, ok := method.Responses[200]; assert.True(ok) {
				if body, ok := response.Bodies["application/json"]; assert.True(ok) {
					require.Equal(TypeObject, body.Type)
					if property, ok := body.Properties.Map()["user"]; assert.True(ok) {
						require.Equal("User", property.Type)
					}
					if user, ok := body.Example.Value.Map["user"]; assert.True(ok) {
						if value, ok := user.Map["name"]; assert.True(ok) {
							require.Equal(TypeString, value.Type)
							require.NotEmpty(value.String)
						}
						if value, ok := user.Map["email"]; assert.True(ok) {
							require.Equal(TypeString, value.Type)
							require.NotEmpty(value.String)
						}
					}
					if example, ok := body.Examples["autoGenerated"]; assert.True(ok) {
						if user, ok := example.Value.Map["user"]; assert.True(ok) {
							if value, ok := user.Map["name"]; assert.True(ok) {
								require.Equal(TypeString, value.Type)
								require.NotEmpty(value.String)
							}
							if value, ok := user.Map["email"]; assert.True(ok) {
								require.Equal(TypeString, value.Type)
								require.NotEmpty(value.String)
							}
						}
					}
				}
			}
		}
	}
	if resource, ok := rootdoc.Resources["/users"]; assert.True(ok) {
		if method, ok := resource.Methods["get"]; assert.True(ok) {
			if response, ok := method.Responses[200]; assert.True(ok) {
				if body, ok := response.Bodies["application/json"]; assert.True(ok) {
					require.Equal("User[]", body.Type)
					require.Len(body.Example.Value.Array, 2)
					for _, user := range body.Example.Value.Array {
						require.NotNil(user)
						if value, ok := user.Map["name"]; assert.True(ok) {
							require.Equal(TypeString, value.Type)
							require.NotEmpty(value.String)
						}
						if value, ok := user.Map["email"]; assert.True(ok) {
							require.Equal(TypeString, value.Type)
							require.NotEmpty(value.String)
						}
					}
					if example, ok := body.Examples["autoGenerated"]; assert.True(ok) {
						require.Len(example.Value.Array, 2)
						for _, user := range example.Value.Array {
							require.NotNil(user)
							if value, ok := user.Map["name"]; assert.True(ok) {
								require.Equal(TypeString, value.Type)
								require.NotEmpty(value.String)
							}
							if value, ok := user.Map["email"]; assert.True(ok) {
								require.Equal(TypeString, value.Type)
								require.NotEmpty(value.String)
							}
						}
					}
				}
			}
		}
	}
	if resource, ok := rootdoc.Resources["/users/wrap"]; assert.True(ok) {
		if method, ok := resource.Methods["get"]; assert.True(ok) {
			if response, ok := method.Responses[200]; assert.True(ok) {
				if body, ok := response.Bodies["application/json"]; assert.True(ok) {
					require.Equal(TypeObject, body.Type)
					if property, ok := body.Properties.Map()["users"]; assert.True(ok) {
						require.Equal("User[]", property.Type)
					}
					if users, ok := body.Example.Value.Map["users"]; assert.True(ok) {
						require.Len(users.Array, 2)
						for _, user := range users.Array {
							require.NotNil(user)
							if value, ok := user.Map["name"]; assert.True(ok) {
								require.Equal(TypeString, value.Type)
								require.NotEmpty(value.String)
							}
							if value, ok := user.Map["email"]; assert.True(ok) {
								require.Equal(TypeString, value.Type)
								require.NotEmpty(value.String)
							}
						}
					}
				}
			}
		}
	}
}

func Test_ParseExampleIncludeBinaryFile(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	rootdoc, err := parser.ParseFile("./test-examples/example-include-binary-file.raml")
	require.NoError(err)
	require.NotZero(rootdoc)

	require.Equal("Example include binary file", rootdoc.Title)
	if resource, ok := rootdoc.Resources["/binary"]; assert.True(ok) {
		if method, ok := resource.Methods["get"]; assert.True(ok) {
			if response, ok := method.Responses[200]; assert.True(ok) {
				if body, ok := response.Bodies["image/png"]; assert.True(ok) {
					require.Equal(TypeFile, body.Type)
					require.Len(body.FileTypes, 1)
					require.Equal("*/*", body.FileTypes[0])
					require.Equal(TypeBinary, body.Example.Value.Type)
					require.Len(body.Example.Value.Binary, 14865)
				}
			}
		}
	}
}

func Test_ParseObjectArray(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	rootdoc, err := parser.ParseFile("./test-examples/object-array.raml")
	require.NoError(err)
	require.NotZero(rootdoc)

	if apiType, ok := rootdoc.Types["UserList"]; assert.True(ok) {
		require.Equal("object[]", apiType.Type)
		if property, ok := apiType.Properties.Map()["name"]; assert.True(ok) {
			require.Equal(TypeString, property.Type)
		}
		require.Len(apiType.Example.Value.Array, 2)
		if value, ok := apiType.Example.Value.Array[0].Map["name"]; assert.True(ok) {
			require.Equal("Alice", value.String)
		}
		if value, ok := apiType.Example.Value.Array[1].Map["name"]; assert.True(ok) {
			require.Equal("Bob", value.String)
		}
	}
}

func Test_ParseTrait(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	rootdoc, err := parser.ParseFile("./test-examples/trait.raml")
	require.NoError(err)
	require.NotZero(rootdoc)

	if trait, ok := rootdoc.Traits["RequireLogin"]; assert.True(ok) {
		if header, ok := trait.Headers.Map()["Authorization"]; assert.True(ok) {
			require.Equal(TypeString, header.Type)
		}
	}
	if resource, ok := rootdoc.Resources["/user"]; assert.True(ok) {
		if method, ok := resource.Methods["get"]; assert.True(ok) {
			if assert.Len(method.Is, 1) {
				trait := method.Is[0]
				require.Equal(trait.String, "RequireLogin")
				if header, ok := trait.Headers.Map()["Authorization"]; assert.True(ok) {
					require.Equal(TypeString, header.Type)
				}
			}
		}
	}
}

func Test_GobEncodeDecode(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	parser := NewParser()
	require.NotNil(parser)

	rootdoc, err := parser.ParseFile("./raml-examples/others/mobile-order-api/api.raml")
	require.NoError(err)
	require.NotZero(rootdoc)

	buffer := &bytes.Buffer{}
	enc := gob.NewEncoder(buffer)
	dec := gob.NewDecoder(buffer)

	err = enc.Encode(rootdoc)
	require.NoError(err)
	var decodedDoc RootDocument
	err = dec.Decode(&decodedDoc)
	require.NoError(err)

	srcjson, err := jsonex.MarshalIndent(rootdoc, "", "  ")
	require.NoError(err)
	dstjson, err := jsonex.MarshalIndent(decodedDoc, "", "  ")
	require.NoError(err)
	requireutil.RequireText(t, string(srcjson), string(dstjson))
}
