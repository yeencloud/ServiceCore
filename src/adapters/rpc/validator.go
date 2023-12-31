package rpc

import (
	"fmt"
	"github.com/yeencloud/ServiceCore/src/adapters/reflect/tags"
	"github.com/yeencloud/ServiceCore/src/helpers"
	"reflect"
)

type ValidationFunc func(name string, val any, rules []string, errors *[]string)

type Validator struct {
	functions map[string]ValidationFunc
}

func NewValidator() Validator {
	validator := Validator{}

	validator.functions = map[string]ValidationFunc{
		"format":    validator.format,
		"minLength": validator.minimumLength,
		"minValue":  validator.minimumValue,
		"maxValue":  validator.maximumValue,
	}

	return validator
}

func (v *Validator) Validate(name string, value any, tags tags.Tags) []string {
	println("VALIDATE")
	println(reflect.TypeOf(value).Name())

	var errors []string

	for _, t := range tags {
		for n, f := range v.functions {
			if t.Name == n {
				f(name, value, t.Rules, &errors)
			}
		}
	}

	return helpers.ArrayOrNil(errors)
}

func (v *Validator) validationError(name string, error string) string {
	return fmt.Sprintf("%s: %s", name, error)
}

func (v *Validator) validationFieldError(name string, error string) string {
	return fmt.Sprintf("%s: validation field has errored: %s (this is an implementation serviceerror not a validation based one)", name, error)
}
