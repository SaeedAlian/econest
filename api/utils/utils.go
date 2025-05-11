package utils

import (
	"reflect"
	"strings"
)

func FilterStruct(input interface{}, exposures map[string]bool) map[string]interface{} {
	res := make(map[string]interface{})
	v := reflect.ValueOf(input)
	t := reflect.TypeOf(input)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		fieldVal := v.Field(i)

		if !fieldVal.CanInterface() {
			continue
		}

		tag := f.Tag.Get("json")
		ex := strings.Split(f.Tag.Get("exposure"), ",")

		if tag == "-" {
			continue
		}

		if tag == "" {
			tag = strings.ToLower(f.Name)
		} else {
			tag = strings.Split(tag, ",")[0]
		}

		include := false
		for _, e := range ex {
			if exposures[e] {
				include = true
				break
			}
		}
		if !include {
			continue
		}

		typ := fieldVal.Type()
		if typ.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				res[tag] = nil
				continue
			}
			typ = typ.Elem()
		}

		if typ.Kind() == reflect.Struct &&
			typ.PkgPath() == "github.com/SaeedAlian/econest/api/types" {
			res[tag] = FilterStruct(fieldVal.Interface(), exposures)
		} else {
			res[tag] = fieldVal.Interface()
		}
	}

	return res
}

func Ptr[T any](v T) *T {
	return &v
}
