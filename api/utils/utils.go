package utils

import (
	"reflect"
	"strings"
)

func FilterStruct(input interface{}, exposures map[string]bool) map[string]interface{} {
	res := make(map[string]interface{})
	v := reflect.ValueOf(input)
	t := reflect.TypeOf(input)

	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("json")
		ex := strings.Split(f.Tag.Get("exposures"), ",")

		if tag == "-" {
			continue
		}

		if tag == "" {
			tag = strings.ToLower(f.Name)
		} else {
			tag = strings.Split(tag, ",")[0]
		}

		for i, e := range ex {
			if exposures[e] {
				res[tag] = v.Field(i).Interface()
				break
			}
		}
	}

	return res
}
