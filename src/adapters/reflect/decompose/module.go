package decompose

import (
	"errors"
	"reflect"
)

type Module struct {
	Name    string   `required:"true" minLength:"1"`
	Methods []Method `required:"true" minLength:"1"`
}

func DecomposeModule(moduleToDecompose any, moduleName string) (*Module, error) {
	module := Module{}

	moduleValue := reflect.ValueOf(moduleToDecompose)
	if moduleName == "" {
		moduleName = reflect.Indirect(moduleValue).Type().Name()
	}
	if !IsExported(moduleName) {
		return nil, errors.New("module name must be exported")
	}
	module.Name = moduleName

	moduleType := reflect.TypeOf(moduleToDecompose)
	module.Methods = decomposeMethodsOfModule(moduleType)

	if len(module.Methods) == 0 {
		return nil, errors.New("module has no exported methods")
	}

	return &module, nil
}
