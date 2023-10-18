package decompose

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

//TestValueForBuiltInType

func TestValueForBuiltInTypeUsingInteger(t *testing.T) {
	var i int
	IntType := reflect.TypeOf(i)

	value := valueForBuiltInType(IntType)

	assert.Equal(t, "int", value)
}

func TestValueForBuiltInTypeUsingString(t *testing.T) {
	var s string
	StringType := reflect.TypeOf(s)

	value := valueForBuiltInType(StringType)

	assert.Equal(t, "string", value)
}

func TestValueForBuiltInTypeUsingFloat(t *testing.T) {
	var f float64
	FloatType := reflect.TypeOf(f)

	value := valueForBuiltInType(FloatType)

	assert.Equal(t, "float64", value)
}

func TestValueForBuiltInTypeUsingBool(t *testing.T) {
	var b bool
	BoolType := reflect.TypeOf(b)

	value := valueForBuiltInType(BoolType)

	assert.Equal(t, "bool", value)
}

func TestValueForBuiltInTypeUsingByte(t *testing.T) {
	var b byte
	ByteType := reflect.TypeOf(b)

	value := valueForBuiltInType(ByteType)

	assert.Equal(t, "uint8", value)
}

func TestValueForBuiltInTypeUsingStringMap(t *testing.T) {
	var m map[string]string
	MapType := reflect.TypeOf(m)

	value := valueForMapType(MapType)

	assert.Equal(t, "map[string]string", value)
}

func TestValueForBuiltInTypeUsingIntMap(t *testing.T) {
	var m map[string]int
	MapType := reflect.TypeOf(m)

	value := valueForMapType(MapType)

	assert.Equal(t, "map[string]int", value)
}

func TestValueForBuiltInTypeUsingStructMap(t *testing.T) {
	type s struct {
		Name string
	}
	var m map[string]s
	MapType := reflect.TypeOf(m)

	value := valueForMapType(MapType)

	assert.Equal(t, "map[string]string", value)
}

func TestValueForBuiltInTypeUsingSlice(t *testing.T) {
	var s []string
	SliceType := reflect.TypeOf(s)

	value := valueForSliceType(SliceType)

	assert.Equal(t, "[]string", value)
}

func TestValueForBuiltInTypeUsingStruct(t *testing.T) {
	var s struct{}
	StructType := reflect.TypeOf(s)

	value := valueForStructType(StructType)

	assert.Equal(t, "struct", value)
}