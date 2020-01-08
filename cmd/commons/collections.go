package commons

import (
	"reflect"
)

func ConvertToSliceOfInterfaces(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("passed argument must be a slice type")
	}

	sliceOfInterfaces := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		sliceOfInterfaces[i] = s.Index(i).Interface()
	}

	return sliceOfInterfaces
}
