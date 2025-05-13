package utils

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/SaeedAlian/econest/api/types"
)

var structPaths []string = []string{
	"github.com/SaeedAlian/econest/api/types",
}

var Validator = validator.New()

func ParseJSONFromRequest(r *http.Request, payload any) error {
	body := r.Body

	if body == nil {
		return types.ErrReqBodyNotFound
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSONInResponse(
	w http.ResponseWriter,
	status int,
	payload any,
	headers *map[string]string,
) error {
	if headers != nil {
		for k, v := range *headers {
			w.Header().Add(k, v)
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(payload)
}

func WriteErrorInResponse(w http.ResponseWriter, status int, message error) error {
	return WriteJSONInResponse(w, status, map[string]string{"message": message.Error()}, nil)
}

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

		pkgPath := typ.PkgPath()

		isStruct := false
		for _, p := range structPaths {
			if p == pkgPath {
				isStruct = true
				break
			}
		}

		if typ.Kind() == reflect.Struct && isStruct {
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
