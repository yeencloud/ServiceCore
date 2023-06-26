package decompose

import (
	"errors"
	"reflect"
)

type Module struct {
	Name    string
	Methods []Method
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

	return &module, nil
}