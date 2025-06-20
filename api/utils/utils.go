package utils

import (
	"encoding/json"
	"errors"
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
	"github.com/lib/pq"

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

func WriteErrorInResponse(w http.ResponseWriter, status int, err error) error {
	var formattedErr error = err

	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			{
				formattedErr = formatUniqueViolation(pgErr)
			}
		case "23503":
			{
				formattedErr = formatForeignKeyViolation(pgErr)
			}
		case "P0001":
			{
				formattedErr = errors.New(pgErr.Message)
			}
		case "22P02":
			{
				formattedErr = formatEnumViolation(pgErr)
			}
		default:
			{
				formattedErr = errors.New("database error: " + pgErr.Message)
			}
		}
	}

	res := types.HTTPError{
		Message: formattedErr.Error(),
	}

	return WriteJSONInResponse(w, status, res, nil)
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

		WriteJSONInResponse(w, http.StatusOK, types.FileUploadResponse{
			FileName: filename,
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

func ParseIntURLParam(name string, vars map[string]string) (int, error) {
	param, ok := vars[name]
	if !ok {
		return -1, types.ErrInvalidParamValue(name)
	}

	parsed, err := strconv.Atoi(param)
	if err != nil {
		return -1, types.ErrInvalidParamValue(name)
	}

	return parsed, nil
}

func ParseStringURLParam(name string, vars map[string]string) (string, error) {
	param, ok := vars[name]
	if !ok {
		return "", types.ErrInvalidParamValue(name)
	}

	if param == "" {
		return "", types.ErrInvalidParamValue(name)
	}

	return param, nil
}

func ParseRequestPayload[T any](r *http.Request, payload *T) error {
	if err := ParseJSONFromRequest(r, payload); err != nil {
		return types.ErrInvalidPayload
	}

	if err := Validator.Struct(*payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return types.ErrInvalidPayloadField(errors[0])
	}

	return nil
}

func ParseURLQuery(mapping map[string]any, values url.Values) error {
	for key, ptr := range mapping {
		v := reflect.ValueOf(ptr)
		if v.Kind() != reflect.Ptr || v.IsNil() {
			return types.ErrQueryMappingNilValueReceived(key)
		}
		v = v.Elem()
		vType := v.Type().Elem()
		vKind := vType.Kind()

		vals, ok := values[key]
		if !ok || len(vals) == 0 {
			continue
		}
		rawValue := vals[0]

		if vType == reflect.TypeOf(time.Time{}) {
			parsed, err := time.Parse(time.RFC3339, rawValue)
			if err != nil {
				return types.ErrInvalidQueryValue(key)
			}
			v.Set(reflect.ValueOf(&parsed))
			continue
		}

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

func CopyFileIntoResponse(dir string, filename string, w http.ResponseWriter) {
	filePath := fmt.Sprintf("%s/%s", dir, filename)
	file, err := os.Open(filePath)
	if err != nil {
		WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrCouldNotOpenFile)
		return
	}

	mime, err := mimetype.DetectReader(file)
	if err != nil {
		WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrCouldNotGetFileMimeType)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrCouldNotResetFileReader)
		return
	}

	stats, err := file.Stat()
	if err != nil {
		WriteErrorInResponse(w, http.StatusInternalServerError, types.ErrCouldNotGetFileStats)
		return
	}

	w.Header().Set("Content-Type", mime.String())
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stats.Size()))

	_, err = io.Copy(w, file)
	if err != nil {
		WriteErrorInResponse(
			w,
			http.StatusInternalServerError,
			types.ErrCouldNotCopyFileIntoResponse,
		)
		return
	}
}

func formatUniqueViolation(e *pq.Error) error {
	if e.Column != "" {
		return types.ErrUniqueConstraintViolationForColumn(e.Column)
	}

	switch e.Constraint {
	case "roles_name_key":
		return types.ErrDuplicateRoleName

	case "permission_groups_name_key":
		return types.ErrDuplicatePermissionGroupName

	case "stores_name_key":
		return types.ErrDuplicateStoreName

	case "phonenumbers_number_key":
		return types.ErrDuplicatePhoneNumber

	case "product_categories_image_name_key":
		return types.ErrDuplicateProductCategoryImageName

	case "product_images_image_name_key":
		return types.ErrDuplicateProductImageName

	case "users_email_key":
		return types.ErrDuplicateUserEmail

	case "users_username_key":
		return types.ErrDuplicateUsername

	case "products_slug_key":
		return types.ErrDuplicateProductSlug

	default:
		return types.ErrUniqueConstraintViolation
	}
}

func formatForeignKeyViolation(e *pq.Error) error {
	if e.Constraint != "" {
		switch e.Constraint {

		case "users_role_id_fkey":
			{
				return types.ErrRoleNotFound
			}

		case "users_settings_user_id_fkey":
			{
				return types.ErrUserNotFound
			}

		case "phonenumbers_user_id_fkey":
			{
				return types.ErrUserNotFound
			}

		case "phonenumbers_store_id_fkey":
			{
				return types.ErrStoreNotFound
			}

		case "addresses_user_id_fkey":
			{
				return types.ErrUserNotFound
			}

		case "addresses_store_id_fkey":
			{
				return types.ErrStoreNotFound
			}

		case "wallets_user_id_fkey":
			{
				return types.ErrUserNotFound
			}

		case "products_subcategory_id_fkey":
			{
				return types.ErrProductCategoryNotFound
			}

		case "product_offers_product_id_fkey":
			{
				return types.ErrProductNotFound
			}

		case "product_images_product_id_fkey":
			{
				return types.ErrProductNotFound
			}

		case "product_specs_product_id_fkey":
			{
				return types.ErrProductNotFound
			}

		case "product_tag_assignments_product_id_fkey":
			{
				return types.ErrProductNotFound
			}

		case "product_tag_assignments_tag_id_fkey":
			{
				return types.ErrProductTagNotFound
			}

		case "product_attribute_options_attribute_id_fkey":
			{
				return types.ErrProductAttributeNotFound
			}

		case "product_variants_product_id_fkey":
			{
				return types.ErrProductNotFound
			}

		case "product_variant_attribute_options_variant_id_fkey":
			{
				return types.ErrProductVariantNotFound
			}

		case "product_variant_attribute_options_attribute_id_fkey":
			{
				return types.ErrProductAttributeNotFound
			}

		case "product_variant_attribute_options_option_id_fkey":
			{
				return types.ErrProductAttributeOptionNotFound
			}

		case "product_comments_product_id_fkey":
			{
				return types.ErrProductNotFound
			}

		case "product_comments_user_id_fkey":
			{
				return types.ErrUserNotFound
			}

		case "role_group_assignments_role_id_fkey":
			{
				return types.ErrRoleNotFound
			}

		case "role_group_assignments_permission_group_id_fkey":
			{
				return types.ErrPermissionGroupNotFound
			}

		case "group_resource_permissions_group_id_fkey":
			{
				return types.ErrPermissionGroupNotFound
			}

		case "group_action_permissions_group_id_fkey":
			{
				return types.ErrPermissionGroupNotFound
			}

		case "wallet_transactions_wallet_id_fkey":
			{
				return types.ErrWalletNotFound
			}

		case "stores_owner_id_fkey":
			{
				return types.ErrStoreOwnerNotFound
			}

		case "stores_settings_store_id_fkey":
			{
				return types.ErrStoreNotFound
			}

		case "store_owned_products_store_id_fkey":
			{
				return types.ErrStoreNotFound
			}

		case "store_owned_products_product_id_fkey":
			{
				return types.ErrProductNotFound
			}

		case "orders_user_id_fkey":
			{
				return types.ErrUserNotFound
			}

		case "order_payments_order_id_fkey":
			{
				return types.ErrOrderNotFound
			}

		case "order_shipments_order_id_fkey":
			{
				return types.ErrOrderNotFound
			}

		case "order_shipments_receiver_address_id_fkey":
			{
				return types.ErrShipmentAddressNotFound
			}

		case "order_product_variants_order_id_fkey":
			{
				return types.ErrOrderNotFound
			}

		case "order_product_variants_variant_id_fkey":
			{
				return types.ErrProductVariantNotFound
			}

		default:
			return types.ErrForeignKeyViolationForColumn
		}
	}

	return types.ErrForeignKeyViolationForColumn
}

func formatEnumViolation(e *pq.Error) error {
	msg := e.Message
	switch {
	case strings.Contains(msg, `"actions"`):
		return types.ErrInvalidActionEnum

	case strings.Contains(msg, `"resources"`):
		return types.ErrInvalidResourceEnum

	case strings.Contains(msg, `"transaction_types"`):
		return types.ErrInvalidTransactionTypeEnum

	case strings.Contains(msg, `"transaction_statuses"`):
		return types.ErrInvalidTransactionStatusEnum

	case strings.Contains(msg, `"order_payment_statuses"`):
		return types.ErrInvalidOrderPaymentStatusEnum

	case strings.Contains(msg, `"order_shipment_statuses"`):
		return types.ErrInvalidOrderShipmentStatusEnum

	default:
		return types.ErrInvalidInputFormat
	}
}
