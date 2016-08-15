package parser

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/tsaikd/KDGoLib/jsonex"
)

var jsonNull, _ = jsonex.Marshal(nil)

// LoadRAMLFromDir load RAML data from directory, concat *.raml
func LoadRAMLFromDir(dirPath string) (ramlData []byte, err error) {
	var filenames []string
	if filenames, err = filepath.Glob(filepath.Join(dirPath, "*.raml")); err != nil {
		return
	}
	sort.Strings(filenames)

	buffer := &bytes.Buffer{}
	for _, filename := range filenames {
		var filedata []byte
		if filedata, err = ioutil.ReadFile(filename); err != nil {
			return
		}
		if _, err = buffer.Write(filedata); err != nil {
			return
		}
		if _, err = buffer.WriteRune('\n'); err != nil {
			return
		}
	}

	return buffer.Bytes(), nil
}

// IsArrayType check name has suffix []
func IsArrayType(name string) (originName string, isArray bool) {
	isArray = strings.HasSuffix(name, "[]")
	if isArray {
		return name[:len(name)-2], isArray
	}
	return name, isArray
}

// ParseYAMLError return the error detail info if it's an YAML parse error,
// yaml parser return error without export error type,
// so using regexp to check
func ParseYAMLError(yamlErr error) (line int64, reason string, ok bool) {
	if yamlErr == nil {
		return 0, "", false
	}

	regYAMLError := regexp.MustCompile(`^yaml: line (\d+): (.+)$`)
	res := regYAMLError.FindAllStringSubmatch(yamlErr.Error(), -1)
	if res == nil || len(res) < 1 {
		return 0, "", false
	}

	var err error
	if line, err = strconv.ParseInt(res[0][1], 0, 64); err != nil {
		return 0, "", false
	}
	reason = res[0][2]
	return line, reason, true
}

// GetLinesInRange return text from (line - distance) to (line + distance)
func GetLinesInRange(data string, sep string, line int64, distance int64) string {
	lines := strings.Split(data, sep)
	minline := line - distance - 1
	if minline < 0 {
		minline = 0
	}
	maxline := line + distance
	if maxline >= int64(len(lines)) {
		maxline = int64(len(lines))
	}
	return strings.Join(lines[minline:maxline], sep)
}
