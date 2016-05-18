package parser

import "github.com/tsaikd/KDGoLib/errutil"

// errors
var (
	ErrorUnsupportedParserConfig1      = errutil.NewFactory("unsupported parser config: %q")
	ErrorInvaludParserConfigValueType3 = errutil.NewFactory(`value type of parser config %q should be "%T" but got "%T"`)
	ErrorUnexpectedRAMLVersion2        = errutil.NewFactory("RAML version should be %q but got %q")
)
