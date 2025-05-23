package utils

import (
	"encoding/json"
	"net/http"
	"os"
	"reflect"
	"slices"
	"strings"
	"time"

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

func DeleteCookie(w http.ResponseWriter, cookie *http.Cookie) {
	http.SetCookie(w, &http.Cookie{
		Name:        cookie.Name,
		Value:       cookie.Value,
		Quoted:      cookie.Quoted,
		Path:        cookie.Path,
		Domain:      cookie.Domain,
		Secure:      cookie.Secure,
		HttpOnly:    cookie.HttpOnly,
		SameSite:    cookie.SameSite,
		Partitioned: cookie.Partitioned,
		Raw:         cookie.Raw,
		Unparsed:    cookie.Unparsed,
		Expires:     time.Now().Add(-7 * 24 * time.Hour),
		MaxAge:      -1,
	})
}

func FilterStruct(input any, exposures map[string]bool) map[string]any {
	res := make(map[string]any)
	v := reflect.ValueOf(input)
	t := reflect.TypeOf(input)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
		t = t.Elem()
	}

	for i := range v.NumField() {
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

		isStruct := slices.Contains(structPaths, pkgPath)

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

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
