package rpc

import (
	"fmt"
	"github.com/yeencloud/ServiceCore/src/adapters/reflect/tags"
	"reflect"
	"strings"
)

// getInnerValue if parameter is a pointer, we retrieve it's value otherwise we just return it
// that helps us being sure that we are always working with a value
func getInnerValue(parameter *reflect.Value) *reflect.Value {
	if parameter.Kind() == reflect.Ptr {
		val := parameter.Elem()
		parameter = &val
	}
	return parameter
}

// getTypeForValue will return the type of the value if it's not nil, otherwise it will return the type of the required field
func getTypeForValue(currentField reflect.Value, value reflect.Value) reflect.Type {
	if value == reflect.ValueOf(nil) {
		return currentField.Type()
	} else {
		return reflect.TypeOf(value.Interface())
	}
}

func createArrayParameter(name string, typeInfo reflect.Type, request []interface{}) (any, []string) {
	arr := reflect.MakeSlice(typeInfo, len(request), len(request))
	var errors []string

	for k, v := range request {
		fName := fieldName(name, fmt.Sprintf("[%d]", k))

		println(fName + " is a " + typeInfo.Elem().Kind().String())

		va, ers := fillAppropriateType(fName, typeInfo.Elem(), v)

		if ers != nil && len(ers) > 0 {
			errors = append(errors, ers...)
			continue
		}

		arr.Index(k).Set(reflect.ValueOf(va))
	}

	if len(errors) > 0 {
		return nil, errors
	}

	return arr.Interface(), nil
}

func fieldName(parent string, fieldname string) string {
	if parent != "" {
		if strings.HasPrefix(fieldname, "[") {
			return parent + fieldname
		}
		return parent + "." + fieldname
	} else {
		return fieldname
	}
}

func createMapParameter(name string, typeInfo reflect.Type, request map[string]interface{}) (any, []string) {
	mapp := reflect.MakeMap(typeInfo)
	var errors []string

	for k, v := range request {
		fName := fieldName(name, k)

		println(fName + " is a " + typeInfo.Elem().Kind().String())

		va, ers := fillAppropriateType(fName, typeInfo.Elem(), v)

		if ers != nil && len(ers) > 0 {
			errors = append(errors, ers...)
			continue
		}

		mapp.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(va))
	}
	if errors != nil && len(errors) > 0 {
		return nil, errors
	}

	return mapp.Interface(), nil
}

func fillAppropriateType(name string, currentField reflect.Type, data interface{}) (any, []string) {
	println(name + " is a " + currentField.Kind().String())

	var errors []string

	var nestedErrors []string
	var nestedValue any

	if currentField.Kind() == reflect.Struct {
		nestedValue = createStructParameter(name, currentField, data.(map[string]interface{}), &nestedErrors)
	} else if currentField.Kind() == reflect.Slice {
		nestedValue, nestedErrors = createArrayParameter(name, currentField, data.([]interface{}))
	} else if currentField.Kind() == reflect.Map {
		nestedValue, nestedErrors = createMapParameter(name, currentField, data.(map[string]interface{}))
	} else {
		return data, nil
	}

	if nestedErrors != nil && len(nestedErrors) > 0 {
		errors = append(errors, nestedErrors...)
		return nil, errors
	}

	return nestedValue, nil
}

func createStructParameter(name string, typeInfo reflect.Type, request map[string]interface{}, errors *[]string) any {
	builtStruct := reflect.New(typeInfo)

	for k := 0; k < typeInfo.NumField(); k++ {
		currentField := typeInfo.Field(k)
		fName := fieldName(name, currentField.Name)
		subValue, valueExists := request[currentField.Name]

		defer func() {
			if r := recover(); r != nil {
				*errors = append(*errors, "cannot convert field '"+fName+"' to "+currentField.Type.Name())
				println("Failed " + fName)
				return
			}
		}()

		fieldTag := tags.GetTags(string(currentField.Tag))

		isRequired := fieldTag.Required()

		if !valueExists {
			if isRequired {
				*errors = append(*errors, "field '"+fName+"' is required")
			}
			continue
		} else {
			//Empty checks
		}

		filledValue, err := fillAppropriateType(fName, currentField.Type, subValue)

		if err != nil && len(err) > 0 {
			*errors = append(*errors, err...)
		}

		if filledValue != nil {
			validator := NewValidator()

			getInnerValue(&builtStruct).Field(k).Set(reflect.ValueOf(filledValue).
				Convert(getInnerValue(&builtStruct).Field(k).Type())) //The convert function call should be useful when working with numbers since json numbers are considered float64 this will be used to convert them to the correct value (int for example)

			*errors = append(*errors, validator.Validate(fName, getInnerValue(&builtStruct).Field(k).Interface(), fieldTag)...)

		}

	}

	if len(*errors) > 0 {
		return nil
	}

	return builtStruct.Elem().Interface()
}

func CreateMethodParameter(typeInfo reflect.Type, request map[string]interface{}) (any, []string) {
	var errors []string

	value := createStructParameter("", typeInfo, request, &errors)

	if len(errors) > 0 {
		return nil, errors
	}

	return value, nil
}
