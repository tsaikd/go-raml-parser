package parser

// QueryParameters map of QueryParameter
type QueryParameters struct {
	Properties
}

// QueryParameter The queryParameters node specifies the set of query
// parameters from which the query string is composed. When applying the
// restrictions defined by the API, processors MUST regard the query string
// as a set of query parameters according to the URL encoding format.
// The value of the queryParameters node is a properties declaration object,
// as is the value of the properties object of a type declaration.
// Each property in this declaration object is referred to as a query parameter declaration.
// Each property name specifies an allowed query parameter name.
// Each property value specifies the query parameter value type as the name
// of a type or an inline type declaration.
type QueryParameter struct {
	Property
}
