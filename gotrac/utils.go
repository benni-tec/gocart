package gotrac

import "reflect"

func P[T any](v T) *T {
	return &v
}

func structToMap(obj any) map[string]any {
	value := reflect.ValueOf(obj)
	value = reflect.Indirect(value)

	typ := value.Type()

	maped := make(map[string]any, typ.NumField())
	for i := range typ.NumField() {
		maped[typ.Field(i).Name] = value.Field(i).Interface()
	}

	return maped
}
