package parser

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

func getFieldJSONName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" {
		return field.Name
	}

	tagName := strings.Split(tag, ",")[0]
	tagName = strings.TrimSpace(tagName)
	switch tagName {
	case "-":
		return ""
	case "":
		return field.Name
	default:
		return tagName
	}
}

func marshalStructToMap(output map[string]interface{}, refv reflect.Value) {
	reft := refv.Type()
	for i, max := 0, reft.NumField(); i < max; i++ {
		marshalFieldToMap(output, reft.Field(i), refv.Field(i))
	}
}

func marshalFieldToMap(output map[string]interface{}, field reflect.StructField, refv reflect.Value) {
	tag := field.Tag.Get("json")
	if refv.CanInterface() {
		if emptier, ok := refv.Interface().(SupportToCheckEmpty); ok {
			if emptier.IsEmpty() && strings.Contains(tag, "omitempty") {
				return
			}
		}
	}

	if field.Anonymous {
		marshalStructToMap(output, refv)
		return
	}

	switch field.Type.Kind() {
	case reflect.String:
		if refv.String() == "" && strings.Contains(tag, "omitempty") {
			return
		}
	case reflect.Slice:
		if refv.Len() < 1 && strings.Contains(tag, "omitempty") {
			return
		}
	case reflect.Bool:
		defval, err := strconv.ParseBool(field.Tag.Get("default"))
		if err == nil {
			if refv.Bool() == defval && strings.Contains(tag, "omitdefault") {
				return
			}
		} else {
			if refv.Bool() == false && strings.Contains(tag, "omitempty") {
				return
			}
		}
	case reflect.Int64:
		defval, err := strconv.ParseInt(field.Tag.Get("default"), 0, 64)
		if err == nil {
			if refv.Int() == defval && strings.Contains(tag, "omitdefault") {
				return
			}
		} else {
			if refv.Int() == 0 && strings.Contains(tag, "omitempty") {
				return
			}
		}
	}

	if name := getFieldJSONName(field); name != "" {
		if refv.CanInterface() {
			output[name] = refv.Interface()
		}
	}
}

// MarshalJSONWithoutEmptyStruct do not marshal empty struct
func MarshalJSONWithoutEmptyStruct(v interface{}) ([]byte, error) {
	if v == nil {
		return json.Marshal(nil)
	}

	refv := reflect.ValueOf(v)
	switch refv.Kind() {
	case reflect.Ptr:
		if refv.Elem().Kind() == reflect.Struct {
			return MarshalJSONWithoutEmptyStruct(refv.Elem().Interface())
		}
	case reflect.Struct:
		if emptier, ok := refv.Interface().(SupportToCheckEmpty); ok {
			if emptier.IsEmpty() {
				return json.Marshal(nil)
			}
		}

		bufmap := map[string]interface{}{}
		marshalStructToMap(bufmap, refv)
		if len(bufmap) > 0 {
			return json.Marshal(bufmap)
		}

		return json.Marshal(nil)
	}

	return json.Marshal(v)
}
