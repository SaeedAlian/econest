package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
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

func CreateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")

	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return -1
	}, slug)

	return slug
}

func FileUploadHandler(
	field string,
	maxSizeInMB int64,
	mimeTypes []string,
	directory string,
) http.HandlerFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		maxSizeInBytes := maxSizeInMB * 1024 * 1024

		r.Body = http.MaxBytesReader(w, r.Body, maxSizeInBytes)
		if err := r.ParseMultipartForm(maxSizeInBytes); err != nil {
			WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrUploadSizeTooBig(int(maxSizeInMB)),
			)
			return
		}

		file, handler, err := r.FormFile(field)
		if err != nil {
			WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrCannotRetrieveFile(err),
			)
			return
		}
		defer file.Close()

		buf := make([]byte, 512)
		_, err = file.Read(buf)
		if err != nil {
			WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrFileUpload(err),
			)
			return
		}

		mimeTypeFromHandler := handler.Header.Get("Content-Type")
		mimeTypeFromMTLib := mimetype.Detect(buf).String()

		typeFound := false

		for i := range mimeTypes {
			m := mimeTypes[i]

			if mimeTypeFromHandler == m || mimeTypeFromMTLib == m {
				typeFound = true
			}
		}

		if !typeFound {
			allowedMimeTypesString := strings.Join(mimeTypes, " , ")

			WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrNotAllowedFileType(allowedMimeTypesString),
			)
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrFileUpload(err),
			)
			return
		}

		err = os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrFileUpload(err),
			)
			return
		}

		filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), filepath.Base(handler.Filename))
		fullpath := fmt.Sprintf("%s/%s", directory, filename)

		dest, err := os.Create(fullpath)
		if err != nil {
			WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrFileUpload(err),
			)
			return
		}
		defer dest.Close()

		_, err = io.Copy(dest, file)
		if err != nil {
			WriteErrorInResponse(
				w,
				http.StatusBadRequest,
				types.ErrFileUpload(err),
			)
			return
		}

		WriteJSONInResponse(w, http.StatusOK, map[string]string{
			"fileName": filename,
		}, nil)
	}
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

func GetPageCount(totalEntities int64, pageLimit int64) int32 {
	pages := math.Floor(float64(totalEntities / pageLimit))
	remainder := totalEntities % pageLimit

	if remainder > 0 {
		pages++
	}

	return int32(pages)
}

func ParseURLQuery(mapping map[string]any, values url.Values) error {
	for key, ptr := range mapping {
		v := reflect.ValueOf(ptr)
		if v.Kind() != reflect.Ptr || v.IsNil() {
			return types.ErrQueryMappingNilValueReceived(key)
		}
		v = v.Elem()
		vKind := v.Type().Elem().Kind()

		vals, ok := values[key]
		if !ok || len(vals) == 0 {
			continue
		}
		rawValue := vals[0]

		switch vKind {
		case reflect.Bool:
			{
				var boolVal bool

				if rawValue == "1" {
					boolVal = true
				} else if rawValue == "0" {
					boolVal = false
				} else {
					return types.ErrInvalidQueryValue(key)
				}

				v.Set(reflect.ValueOf(&boolVal))
			}

		case reflect.Float32:
			{
				parsed, err := strconv.ParseFloat(rawValue, 32)
				res := float32(parsed)
				if err != nil {
					return types.ErrInvalidQueryValue(key)
				}

				v.Set(reflect.ValueOf(&res))
			}

		case reflect.Float64:
			{
				parsed, err := strconv.ParseFloat(rawValue, 64)
				if err != nil {
					return types.ErrInvalidQueryValue(key)
				}

				v.Set(reflect.ValueOf(&parsed))
			}

		case reflect.Int:
			{
				parsed, err := strconv.Atoi(rawValue)
				if err != nil {
					return types.ErrInvalidQueryValue(key)
				}

				v.Set(reflect.ValueOf(&parsed))
			}

		case reflect.Int16:
			{
				parsed, err := strconv.ParseInt(rawValue, 10, 16)
				res := int16(parsed)
				if err != nil {
					return types.ErrInvalidQueryValue(key)
				}

				v.Set(reflect.ValueOf(&res))
			}

		case reflect.Int32:
			{
				parsed, err := strconv.ParseInt(rawValue, 10, 32)
				res := int32(parsed)
				if err != nil {
					return types.ErrInvalidQueryValue(key)
				}

				v.Set(reflect.ValueOf(&res))
			}

		case reflect.Int64:
			{
				parsed, err := strconv.ParseInt(rawValue, 10, 64)
				if err != nil {
					return types.ErrInvalidQueryValue(key)
				}

				v.Set(reflect.ValueOf(&parsed))
			}

		case reflect.String:
			{
				v.Set(reflect.ValueOf(&rawValue))
			}
		}

	}

	return nil
}
