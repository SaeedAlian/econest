package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"

	db_manager "github.com/SaeedAlian/econest/api/db/manager"
	"github.com/SaeedAlian/econest/api/types"
	testutils "github.com/SaeedAlian/econest/api/utils/tests"
)

func TestUserService(t *testing.T) {
	db := testutils.SetupTestDB(t)
	manager := db_manager.NewManager(db)
	handler := NewHandler(manager)

	t.Run("should register successfully", func(t *testing.T) {
		payload := types.CreateUserPayload{
			FullName:  "Test User 1",
			Username:  strings.ToLower("testuser1"),
			Email:     strings.ToLower("testuser1@gmail.com"),
			BirthDate: time.Date(1999, 2, 25, 0, 0, 0, 0, time.UTC),
			Password:  "password123",
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/register/customer", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register/customer", handler.registerCustomer).Methods("POST")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected code %d, received %d", http.StatusCreated, rr.Code)
		}

		created, err := handler.db.GetUserByUsername(strings.ToLower("testuser1"))
		if err != nil {
			t.Fatal(err)
		}

		if created == nil {
			t.Error("Expected for the created user to exist, but it's not")
		}
	})

	t.Run("should fail to register due to duplicated username", func(t *testing.T) {
		payload := types.CreateUserPayload{
			FullName:  "Test User 1",
			Username:  strings.ToLower("testuser1"),
			Email:     strings.ToLower("SecondEmailShouldFail@gmail.com"),
			BirthDate: time.Date(1999, 2, 25, 0, 0, 0, 0, time.UTC),
			Password:  "password123",
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/register/customer", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register/customer", handler.registerCustomer).Methods("POST")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected code %d, received %d", http.StatusBadRequest, rr.Code)
		}

		created, err := handler.db.GetUserByEmail(
			strings.ToLower("SecondEmailShouldFail@gmail.com"),
		)
		if created != nil {
			t.Error("Expected for the created user to not be found, but it has been found")
		}
	})

	t.Run("should fail to register due to duplicated username", func(t *testing.T) {
		payload := types.CreateUserPayload{
			FullName:  "Test User 1",
			Username:  strings.ToLower("testuserfailed"),
			Email:     strings.ToLower("testuser1@gmail.com"),
			BirthDate: time.Date(1999, 2, 25, 0, 0, 0, 0, time.UTC),
			Password:  "password123",
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/register/customer", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register/customer", handler.registerCustomer).Methods("POST")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected code %d, received %d", http.StatusBadRequest, rr.Code)
		}

		created, err := handler.db.GetUserByUsername(
			strings.ToLower("testuserfailed"),
		)
		if created != nil {
			t.Error("Expected for the created user to not be found, but it has been found")
		}
	})

	t.Run("should fail to register because of invalid data", func(t *testing.T) {
		payload := types.CreateUserPayload{
			FullName:  "Test User 1",
			Username:  strings.ToLower("newErrorUser"),
			Email:     strings.ToLower("gmail"),
			BirthDate: time.Date(1999, 2, 25, 0, 0, 0, 0, time.UTC),
			Password:  "password123",
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/register/customer", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register/customer", handler.registerCustomer).Methods("POST")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected code %d, received %d", http.StatusBadRequest, rr.Code)
		}

		created, err := handler.db.GetUserByUsername(strings.ToLower("newErrorUser"))
		if created != nil {
			t.Error("Expected for the created user to not be found, but it has been found")
		}
	})

	t.Run("should fail to register because of short password", func(t *testing.T) {
		payload := types.CreateUserPayload{
			FullName:  "Test User 1",
			Username:  strings.ToLower("newErrorUser"),
			Email:     strings.ToLower("erroruser@gmail.com"),
			BirthDate: time.Date(1999, 2, 25, 0, 0, 0, 0, time.UTC),
			Password:  "1",
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/register/customer", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register/customer", handler.registerCustomer).Methods("POST")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected code %d, received %d", http.StatusBadRequest, rr.Code)
		}

		created, err := handler.db.GetUserByUsername(strings.ToLower("newErrorUser"))
		if created != nil {
			t.Error("Expected for the created user to not be found, but it has been found")
		}
	})

	t.Run("should fail to register because of long password", func(t *testing.T) {
		payload := types.CreateUserPayload{
			FullName:  "Test User 1",
			Username:  strings.ToLower("newErrorUser"),
			Email:     strings.ToLower("erroruser@gmail.com"),
			BirthDate: time.Date(1999, 2, 25, 0, 0, 0, 0, time.UTC),
			Password:  strings.Repeat("1", 5000000),
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/register/customer", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register/customer", handler.registerCustomer).Methods("POST")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected code %d, received %d", http.StatusBadRequest, rr.Code)
		}

		created, err := handler.db.GetUserByUsername(strings.ToLower("newErrorUser"))
		if created != nil {
			t.Error("Expected for the created user to not be found, but it has been found")
		}
	})
}
