package decompose

import (
	"github.com/rs/zerolog/log"
	"github.com/yeencloud/ServiceCore/src/adapters/reflect/tags"
	"github.com/yeencloud/ServiceCore/src/domain/serviceError"
	"github.com/yeencloud/ServiceCore/src/helpers"
	"reflect"
)

type MethodValue struct {
	Type       any        `required:"true"`
	Validation []tags.Tag `json:",omitempty"`
}

type MethodOutput map[string]MethodValue
type MethodInput map[string]MethodValue

type Method struct {
	Name   string       `required:"true"`
	Input  MethodInput  `required:"true"`
	Output MethodOutput `required:"true"`
}

func valueForBuiltInType(typ reflect.Type) string {
	return typ.Kind().String()
}

func valueForSliceType(typ reflect.Type) string {
	return typ.String()
}

func valueForMapType(typ reflect.Type) string {
	return typ.String()
}

func valueForStructType(typ reflect.Type) map[string]MethodValue {
	structContent := map[string]MethodValue{}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name
		fieldType := field.Type

		if !IsExported(fieldName) {
			continue
		}

		if !isExportedOrBuiltinType(fieldType) {
			log.Warn().Str("field", fieldName).Msg("field type not exported")
			continue
		}

		value := MethodValue{}

		knownType := true
		if isTypeStruct(fieldType) {
			value.Type = valueForStructType(fieldType)
		} else if isTypeSlice(fieldType) {
			value.Type = valueForSliceType(fieldType)
		} else if isTypeMap(fieldType) {
			value.Type = valueForMapType(fieldType)
		} else if isBuiltin(fieldType) {
			value.Type = valueForBuiltInType(fieldType)
		} else {
			knownType = false
		}

		if knownType && value.Type != nil {
			structContent[fieldName] = value
		}
	}

	return helpers.MapOrNil(structContent)
}

func fillValues(typ reflect.Type, input bool) map[string]MethodValue {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if !isSupportedForJson(typ) {
		return nil
	}

	if isBuiltin(typ) {
		output := map[string]MethodValue{}

		name := "output"
		if input {
			name = "input"
		}

		val := MethodValue{}
		val.Type = valueForBuiltInType(typ)
		output[name] = val

		return output
	}

	if isTypeStruct(typ) {
		output := map[string]MethodValue{}

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fieldName := field.Name
			fieldType := field.Type

			tagsForValidation := tags.GetTags(string(field.Tag))

			if !IsExported(fieldName) {
				log.Warn().Str("field", fieldName).Msg("field name must be exported")
				continue
			}

			if !isExportedOrBuiltinType(fieldType) {
				log.Warn().Str("field", fieldName).Msg("field type not exported")
				continue
			}

			fieldValue := MethodValue{
				Validation: tagsForValidation,
			}

			if isTypeStruct(fieldType) {
				v := valueForStructType(fieldType)

				if v != nil {
					fieldValue.Type = v
				}
			} else if isTypeSlice(fieldType) {
				fieldValue.Type = valueForSliceType(fieldType)
			} else if isTypeMap(fieldType) {
				fieldValue.Type = valueForMapType(fieldType)
			} else if isBuiltin(fieldType) {
				fieldValue.Type = valueForBuiltInType(fieldType)
			}

			if fieldValue.Type != nil {
				output[fieldName] = fieldValue
			}
		}
		return helpers.MapOrNil(output)
	}
	return nil
}

func decomposeMethodsOfModule(typ reflect.Type) []Method {
	methods := []Method{}

	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		methodName := method.Name
		methodType := method.Type

		if !IsExported(methodName) {
			continue
		}

		if methodType.NumIn() != 2 {
			log.Warn().Str("method", methodName).Msg("method must have one input argument (excluding the receiver)")
			continue
		}

		if methodType.NumOut() != 2 {
			log.Warn().Str("method", methodName).Msg("method must have 2 outputs: output, serviceerror")
			continue
		}

		// arg type can be exported or a built-in type
		argType := methodType.In(1)
		if !isExportedOrBuiltinType(argType) {
			log.Warn().Str("method", methodName).Msg("argument type not exported")
			continue
		}

		if argType.Kind() == reflect.Ptr {
			log.Warn().Str("method", methodName).Msg("argument type should not be a pointer")
			continue
		}

		// output type needs to be exported
		replyType := methodType.Out(0)
		if returnType := methodType.Out(1); returnType != reflect.TypeOf(&serviceError.Error{}) {
			log.Warn().Str("method", methodName).Msg("The method should return a pointer to serviceError.Error as its second return type")
			continue
		}

		if replyType.Kind() == reflect.Ptr {
			replyType = replyType.Elem()
		}

		var input MethodInput
		input = fillValues(argType, true)

		var output MethodOutput
		output = fillValues(replyType, false)

		methods = append(methods, Method{
			Name:   methodName,
			Input:  input,
			Output: output,
		})
	}

	return helpers.ArrayOrNil(methods)
}
