package decompose

import (
	"reflect"
	"unicode"
	"unicode/utf8"
)

func IsExported(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(ch)
}

func isBuiltin(t reflect.Type) bool {
	return t.PkgPath() == ""
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return IsExported(t.Name()) || isBuiltin(t)
}

func isTypeInt(t reflect.Type) bool {
	return t.Kind() == reflect.Int ||
		t.Kind() == reflect.Int8 ||
		t.Kind() == reflect.Int16 ||
		t.Kind() == reflect.Int32 ||
		t.Kind() == reflect.Int64
}

func isTypeUint(t reflect.Type) bool {
	return t.Kind() == reflect.Uint ||
		t.Kind() == reflect.Uint8 ||
		t.Kind() == reflect.Uint16 ||
		t.Kind() == reflect.Uint32 ||
		t.Kind() == reflect.Uint64
}

func isTypeFloat(t reflect.Type) bool {
	return t.Kind() == reflect.Float32 ||
		t.Kind() == reflect.Float64
}

func isTypeString(t reflect.Type) bool {
	return t.Kind() == reflect.String
}

func isTypeBool(t reflect.Type) bool {
	return t.Kind() == reflect.Bool
}

func isTypeStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct
}

func isTypeSlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice
}

func isTypeMap(t reflect.Type) bool {
	return t.Kind() == reflect.Map
}

func isSupportedForJson(t reflect.Type) bool {
	return isTypeInt(t) ||
		isTypeUint(t) ||
		isTypeFloat(t) ||
		isTypeString(t) ||
		isTypeBool(t) ||
		isTypeStruct(t) ||
		isTypeSlice(t) ||
		isTypeMap(t)
}