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

// Returns the first index of the target string t, or -1 if no match is found.
func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// Returns true if the target string t is in the slice.
func Contains(vs []string, t string) bool {
	return Index(vs, t) >= 0
}
