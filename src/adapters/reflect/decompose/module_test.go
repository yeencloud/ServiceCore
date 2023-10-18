package decompose

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ValidModuleExample struct {
}

type ModuleWithNoExportedMethod struct{}

type Person struct {
	Name      string `json:"regularString"`
	Age       int    `greaterThan:"0" lesserThan:"100"`
	Height    float64
	Employed  bool
	Children  map[string]Person
	Socials   map[string]string
	Parents   []Person
	Nicknames []string
}

var validModule = ValidModuleExample{}

func (mod ValidModuleExample) SomeFunction(inputData Person) (Person, error) {
	return inputData, nil
}

func TestDecomposeShouldntErrorOnValidModule(t *testing.T) {
	moduleName := "ModuleExample"

	_, err := DecomposeModule(validModule, moduleName)

	if err != nil {
		t.Error(err)
	}
}

func TestDecomposeWithoutCustomName(t *testing.T) {
	module, _ := DecomposeModule(validModule, "")

	msh, _ := json.Marshal(module)

	println(string(msh))

	assert.Equal(t, module.Name, "ValidModuleExample")
}

func TestDecomposeWithCustomName(t *testing.T) {
	name := "Module"

	module, _ := DecomposeModule(validModule, name)

	assert.Equal(t, module.Name, name)
}

func TestDecomposeShouldErrorOnInvalidModule(t *testing.T) {
	_, err := DecomposeModule(1, "")

	if err == nil {
		t.Error("Should have errored on invalid validModule")
	}
}

func TestDecomposeShouldErrorOnUnexportedModule(t *testing.T) {
	type unexportedModule struct {
	}

	_, err := DecomposeModule(unexportedModule{}, "")

	if err == nil {
		t.Error("Should have errored on unexported validModule")
	}
}

func TestDecomposeShouldErrorOnModuleWithNoExportedMethods(t *testing.T) {
	_, err := DecomposeModule(ModuleWithNoExportedMethod{}, "")

	if err == nil {
		t.Error("Should have errored on validModule with no exported methods")
	}
}

func TestDecomposeShouldReturnAModuleWithAMethodNamedSomeFunction(t *testing.T) {
	module, _ := DecomposeModule(validModule, "")

	success := assert.Len(t, module.Methods, 1)
	if !success {
		return
	}
	assert.Equal(t, module.Methods[0].Name, "SomeFunction")
}

func TestSomeFunctionExportedMethodShouldHaveAStringParameter(t *testing.T) {
	module, _ := DecomposeModule(validModule, "")

	success := assert.Len(t, module.Methods[0].Input, 1)
	if !success {
		return
	}

	val, ok := module.Methods[0].Input["RegularString"]

	assert.True(t, ok)
	assert.Equal(t, val.Type, "string")
}