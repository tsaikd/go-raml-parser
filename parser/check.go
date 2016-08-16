package parser

import (
	"regexp"

	"github.com/tsaikd/KDGoLib/errutil"
)

// CheckValueOption for changing CheckValueAPIType behavior
type CheckValueOption interface{}

// CheckValueOptionAllowIntegerToBeNumber allow type integer to be type number,
// e.g. APIType need a integer, but value is a number
// default: false
type CheckValueOptionAllowIntegerToBeNumber bool

// CheckValueOptionAllowArrayToBeNull allow array type to be null,
// default: false
type CheckValueOptionAllowArrayToBeNull bool

// CheckValueOptionAllowRequiredPropertyToBeEmpty allow required property to be empty value, but still should be existed
// only check the following types: TypeString, TypeArray, TypeObject, TypeBinary
// default: false
type CheckValueOptionAllowRequiredPropertyToBeEmpty bool

// CheckValueAPIType check value is valid for apiType
func CheckValueAPIType(apiType APIType, value Value, options ...CheckValueOption) (err error) {
	allowIntegerToBeNumber := CheckValueOptionAllowIntegerToBeNumber(false)
	allowArrayToBeNull := CheckValueOptionAllowArrayToBeNull(false)
	allowRequiredPropertyToBeEmpty := CheckValueOptionAllowRequiredPropertyToBeEmpty(false)

	for _, option := range options {
		switch optval := option.(type) {
		case CheckValueOptionAllowIntegerToBeNumber:
			allowIntegerToBeNumber = optval
		case CheckValueOptionAllowArrayToBeNull:
			allowArrayToBeNull = optval
		case CheckValueOptionAllowRequiredPropertyToBeEmpty:
			allowRequiredPropertyToBeEmpty = optval
		}
	}

	return checkValueAPIType(
		apiType,
		value,
		allowIntegerToBeNumber,
		allowArrayToBeNull,
		allowRequiredPropertyToBeEmpty,
	)
}

func checkValueAPIType(
	apiType APIType,
	value Value,
	allowIntegerToBeNumber CheckValueOptionAllowIntegerToBeNumber,
	allowArrayToBeNull CheckValueOptionAllowArrayToBeNull,
	allowRequiredPropertyToBeEmpty CheckValueOptionAllowRequiredPropertyToBeEmpty,
) (err error) {
	if value.IsEmpty() {
		// no need to check if value is empty
		return
	}

	if apiType.IsArray {
		if value.Type != TypeArray {
			if !allowArrayToBeNull || value.Type != TypeNull {
				return ErrorPropertyTypeMismatch2.New(nil, apiType.Type, value.Type)
			}
		}

		elemType := apiType
		elemType.IsArray = false
		for i, elemValue := range value.Array {
			if err = checkValueAPIType(
				elemType,
				*elemValue,
				allowIntegerToBeNumber,
				allowArrayToBeNull,
				allowRequiredPropertyToBeEmpty,
			); err != nil {
				switch errutil.FactoryOf(err) {
				case ErrorPropertyTypeMismatch2:
					return ErrorArrayElementTypeMismatch3.New(nil, i, elemType.Type, elemValue.Type)
				}
				return
			}
		}
		return
	}

	switch apiType.NativeType {
	case TypeBoolean, TypeString:
		if apiType.NativeType == value.Type {
			return nil
		}
		return ErrorPropertyTypeMismatch2.New(nil, apiType.Type, value.Type)
	case TypeInteger:
		if apiType.NativeType == value.Type {
			return nil
		}
		if allowIntegerToBeNumber {
			switch value.Type {
			case TypeNumber:
				if value.Number == float64(int64(value.Number)) {
					return nil
				}
			}
		}
		return ErrorPropertyTypeMismatch2.New(nil, apiType.Type, value.Type)
	case TypeNumber:
		if apiType.NativeType == value.Type {
			return nil
		}
		if allowIntegerToBeNumber {
			switch value.Type {
			case TypeInteger:
				return nil
			}
		}
		return ErrorPropertyTypeMismatch2.New(nil, apiType.Type, value.Type)
	case TypeFile:
		// no type check for file type
		return nil
	default:
		if isInlineAPIType(apiType) {
			// no type check if declared by JSON
			return nil
		}

		switch value.Type {
		case TypeObject, TypeNull:
		default:
			return ErrorPropertyTypeMismatch2.New(nil, apiType.Type, value.Type)
		}

		for _, property := range apiType.Properties.Slice() {
			if err = checkPropertyRequired(
				*property,
				value,
				allowArrayToBeNull,
				allowRequiredPropertyToBeEmpty,
				apiType,
			); err != nil {
				return err
			}

			if err = checkPropertyValue(
				*property,
				value,
				allowIntegerToBeNumber,
				allowArrayToBeNull,
				allowRequiredPropertyToBeEmpty,
			); err != nil {
				return err
			}
		}
	}

	return nil
}

func checkPropertyRequired(
	property Property,
	parent Value,
	allowArrayToBeNull CheckValueOptionAllowArrayToBeNull,
	allowRequiredPropertyToBeEmpty CheckValueOptionAllowRequiredPropertyToBeEmpty,
	apiType APIType, // only used for error message
) (err error) {
	if !property.Required {
		return nil
	}
	if property.IsArray && bool(allowArrayToBeNull) {
		return nil
	}

	switch parent.Type {
	case TypeNull:
		return ErrorRequiredProperty2.New(nil, property.Name, apiType.Type)
	case TypeObject:
		value := parent.Map[property.Name]
		if value == nil {
			return ErrorRequiredProperty2.New(nil, property.Name, apiType.Type)
		}
		if !bool(allowRequiredPropertyToBeEmpty) {
			switch value.Type {
			case TypeString, TypeArray, TypeObject, TypeBinary:
				if value.IsZero() {
					return ErrorRequiredProperty2.New(nil, property.Name, apiType.Type)
				}
			}
		}
		return nil
	default:
		panic("check property required with wrong parent type: " + parent.Type)
	}
}

func checkPropertyValue(
	property Property,
	parent Value,
	allowIntegerToBeNumber CheckValueOptionAllowIntegerToBeNumber,
	allowArrayToBeNull CheckValueOptionAllowArrayToBeNull,
	allowRequiredPropertyToBeEmpty CheckValueOptionAllowRequiredPropertyToBeEmpty,
) (err error) {
	value := parent.Map[property.Name]
	if value == nil {
		return nil
	}
	// no need to check recursive if property is not required
	if !property.Required && value.IsZero() {
		return nil
	}

	if err = checkValueAPIType(
		property.APIType,
		*value,
		allowIntegerToBeNumber,
		allowArrayToBeNull,
		allowRequiredPropertyToBeEmpty,
	); err != nil {
		switch errutil.FactoryOf(err) {
		case ErrorPropertyTypeMismatch2:
			return ErrorPropertyTypeMismatch3.New(nil, property.Name, property.Type, value.Type)
		case ErrorArrayElementTypeMismatch3:
			return ErrorPropertyTypeMismatch1.New(err, property.Name)
		}
		return err
	}

	return nil
}

func isInlineAPIType(apiType APIType) bool {
	regValidType := regexp.MustCompile(`^[\w]+(\[\])?$`)
	return !regValidType.MatchString(apiType.Type)
}
