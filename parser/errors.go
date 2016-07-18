package parser

import (
	"strings"

	"github.com/tsaikd/KDGoLib/errutil"
)

// errors
var (
	ErrorCacheNotFound                 = errutil.NewFactory("cache not found")
	ErrorLoadExternalLibrary1          = errutil.NewFactory("load external library failed: %q")
	ErrorUnsupportedParserConfig1      = errutil.NewFactory("unsupported parser config: %q")
	ErrorUnsupportedValueType1         = errutil.NewFactory("unsupported value type: %T")
	ErrorUnsupportedIncludeType1       = errutil.NewFactory("unsupported include for type: %q")
	ErrorInvaludParserConfigValueType2 = errutil.NewFactory(`value type should be "%T" but got "%T"`)
	ErrorInvaludParserConfigValueType3 = errutil.NewFactory(`value type of parser config %q should be "%T" but got "%T"`)
	ErrorUnexpectedRAMLVersion2        = errutil.NewFactory("RAML version should be %q but got %q")
	ErrorEmptyRootDocumentMediaType    = errutil.NewFactory("body without MIME-type and root document do not provide default MediaType")
	ErrorTypeUndefined1                = errutil.NewFactory("Type %q can not find in RAML")
	ErrorTypo2                         = errutil.NewFactory("Detect typo error on %q: %v")
	ErrorArrayElementTypeMismatch3     = errutil.NewFactory("Array element %d type mismatch, expected %q but got %q")
	ErrorPropertyTypeMismatch1         = errutil.NewFactory("Property %q type mismatch")
	ErrorPropertyTypeMismatch2         = errutil.NewFactory("Property type mismatch, expected %q but got %q")
	ErrorPropertyTypeMismatch3         = errutil.NewFactory("Property %q type mismatch, expected %q but got %q")
	ErrorPropertyUndefined2            = errutil.NewFactory("Property %q can not find in APIType %q")
	ErrorRequiredProperty2             = errutil.NewFactory("Property %q is required but not found in %q")
	ErrorUnusedTrait1                  = errutil.NewFactory("Trait %q is unused")
	ErrorUnusedAnnotation1             = errutil.NewFactory("Annotation %q is unused")
)

func isErrorYAMLIntoBool(err error) bool {
	return strings.Contains(err.Error(), "into bool")
}

func isErrorYAMLIntoInt64(err error) bool {
	return strings.Contains(err.Error(), "into int64")
}

func isErrorYAMLIntoFloat64(err error) bool {
	return strings.Contains(err.Error(), "into float64")
}

func isErrorYAMLIntoString(err error) bool {
	return strings.Contains(err.Error(), "into string")
}
