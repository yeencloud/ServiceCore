package ServiceCore

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/yeencloud/ServiceCore/decompose"
	"net/http"
	"reflect"
)

func (shs *ServiceHttpServer) typeForField(fieldName string, inputs decompose.MethodInput) string {
	for k, v := range inputs {
		if k == fieldName {
			vv := reflect.ValueOf(v.Type)
			if vv.Kind() == reflect.Map {
				println("Struct: " + spew.Sdump(vv))
				return vv.String()
			}
			if vv.Kind() == reflect.String {
				println("String: " + spew.Sdump(vv))
				return vv.String()
			}
			return "unsupported"
		}
	}
	return ""
}

func (shs *ServiceHttpServer) valueAdjustedForType(value interface{}, wantedType string) interface{} {
	defer func() {
		if r := recover(); r != nil {
			wantedType = "unsupported"
		}
	}()
	if wantedType == "int" {
		if value == nil {
			return nil
		}
		return int(value.(float64))
	}

	if wantedType == "unsupported" {
		return nil
	}

	return value
}

// A JSON Object have limited types. Here we convert them to reflect the correct type
// For example json numbers are float64, but if we need an int we'll need to convert them smartly
func (shs *ServiceHttpServer) valueForField(fieldName string, clientRequest map[string]interface{}, inputs decompose.MethodInput) reflect.Value {
	wantedType := shs.typeForField(fieldName, inputs)

	value := shs.valueAdjustedForType(clientRequest[fieldName], wantedType)
	if wantedType == "int" && value == nil {
		if clientRequest[fieldName] == nil {
			return reflect.ValueOf(nil)
		}
	}
	return reflect.ValueOf(value)
}

func (shs *ServiceHttpServer) fillParameterArray(parameter *reflect.Value, input []interface{}) []string {
	wantedType := parameter.Type().Elem().String()

	for _, v := range input {
		value := shs.valueAdjustedForType(v, wantedType)

		if value == nil {
			return []string{"array of type '" + wantedType + "' cannot contain unsupported values"}
		}

		typeOfValue := reflect.TypeOf(value)
		if typeOfValue.String() != wantedType {
			return []string{"array of '" + wantedType + "' cannot contain '" + typeOfValue.String() + "'"}
		}
		if wantedType == "int" && value == nil {
			return []string{"array of int cannot contain nil values"}
		}
		parameter.Set(reflect.Append(*parameter, reflect.ValueOf(value)))
	}

	return []string{}
}

func (shs *ServiceHttpServer) isValueRequired(value decompose.MethodValue) bool {
	for _, v := range value.Validation {
		if v.Name == "required" && len(v.Value) > 0 {
			return v.Value[0] == "true"
		}
	}

	return false
}

func (shs *ServiceHttpServer) fillParameterStruct(parameter *reflect.Value, input map[string]decompose.MethodValue, request map[string]interface{}) []string {
	validationErrors := []string{}

	if parameter.Kind() == reflect.Ptr {
		val := parameter.Elem()
		parameter = &val
	}
	println("Parameter type: " + parameter.Kind().String())

	for k, _ := range input {
		value := shs.valueForField(k, request, input)

		if !value.IsValid() && shs.isValueRequired(input[k]) {
			validationErrors = append(validationErrors, "field '"+k+"' is required")
			continue
		}

		var valueType reflect.Type
		if value == reflect.ValueOf(nil) {
			valueType = parameter.FieldByName(k).Type()
		} else {
			valueType = reflect.TypeOf(value.Interface())
		}

		valueIsStruct := valueType.String() == "map[string]interface {}" && parameter.FieldByName(k).Kind().String() == "struct"
		valueIsArray := valueType.String() == "[]interface {}" && parameter.FieldByName(k).Kind().String() == "slice"
		if valueIsStruct {
			field := parameter.FieldByName(k)
			verrors := shs.fillParameterStruct(&field, input[k].Type.(map[string]decompose.MethodValue), request[k].(map[string]interface{}))
			validationErrors = append(validationErrors, verrors...)
		} else if valueIsArray {
			field := parameter.FieldByName(k)
			validationErrors = append(validationErrors, shs.fillParameterArray(&field, value.Interface().([]interface{}))...)
		} else if valueType != parameter.FieldByName(k).Type() {
			validationErrors = append(validationErrors, "field '"+k+"' of type '"+valueType.String()+"' cannot be assigned to "+parameter.FieldByName(k).Kind().String())
		} else {
			if value == reflect.ValueOf(nil) {
				value = reflect.Zero(parameter.FieldByName(k).Type())
			}
			parameter.FieldByName(k).Set(value)
		}
	}
	return validationErrors
}

func (shs *ServiceHttpServer) getParameterStructFromBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		request, err := shs.getRequestStruct(c)
		if err != nil {
			return
		}
		serviceInstance := reflect.ValueOf(shs.service)
		if !serviceInstance.IsValid() {
			return
		}
		methodType, _ := reflect.TypeOf(shs.service).MethodByName(request.Method)
		methodToCall := serviceInstance.MethodByName(request.Method)
		if !methodToCall.IsValid() {
			return
		}

		r := reflect.New(methodType.Type.In(1))

		uncastinput, exist := c.Get("requiredInput")
		if !exist {
			return
		}
		input := uncastinput.(decompose.MethodInput)

		validationErrors := shs.fillParameterStruct(&r, input, request.Request)

		if len(validationErrors) > 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, replyWithError(ErrValidationFailed, validationErrors))
			return
		}

		c.Set("parameter", r.Elem().Interface())
	}
}

func (shs *ServiceHttpServer) callServiceMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		request, err := shs.getRequestStruct(c)
		if err != nil {
			return
		}

		serviceInstance := reflect.ValueOf(shs.service)
		if !serviceInstance.IsValid() {
			c.AbortWithStatusJSON(http.StatusNotImplemented, replyWithError(ErrServiceIsInvalid, nil))
			return
		}

		methodType, _ := reflect.TypeOf(shs.service).MethodByName(request.Method)
		methodToCall := serviceInstance.MethodByName(request.Method)
		_ = methodType
		if !methodToCall.IsValid() {
			c.AbortWithStatusJSON(http.StatusNotImplemented, replyWithError(ErrMethodIsInvalid, nil))
			return
		}

		inpass, found := c.Get("parameter")
		if !found {
			return
		}
		spew.Dump(reflect.ValueOf(inpass).Interface())
		results := methodToCall.Call([]reflect.Value{reflect.ValueOf(inpass)})

		if len(results) == 2 {
			data := results[0]
			err := results[1]

			reply := ServiceReply{
				Module:  request.Service,
				Service: request.Method,
				Version: 1,
			}

			if err.Type() == decompose.TypeOfError && !err.IsNil() {
				err := err.Interface().(error)
				log.Err(err).Msg("service method has errored")
				reply.Error = err.Error()

				c.IndentedJSON(http.StatusInternalServerError, reply)
			} else {
				data := data.Interface()
				b, _ := json.Marshal(&data)
				var m map[string]interface{}
				_ = json.Unmarshal(b, &m)
				reply.Data = m
				c.IndentedJSON(http.StatusOK, reply)
			}
		}
	}
}