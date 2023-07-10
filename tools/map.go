package tools

import "encoding/json"

func MapOrNil[T any](m map[string]T) map[string]T {
	if len(m) > 0 {
		return m
	}
	return nil
}

func AnyToMap(value any) map[string]interface{} {
	b, _ := json.Marshal(value)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	return m
}