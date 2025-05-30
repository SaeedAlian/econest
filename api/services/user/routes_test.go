package user

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"

	"github.com/SaeedAlian/econest/api/config"
	db_manager "github.com/SaeedAlian/econest/api/db/manager"
	"github.com/SaeedAlian/econest/api/services/auth"
	"github.com/SaeedAlian/econest/api/services/smtp"
	"github.com/SaeedAlian/econest/api/types"
	"github.com/SaeedAlian/econest/api/utils"
	testutils "github.com/SaeedAlian/econest/api/utils/tests"
)

var testUsers = [10]types.CreateUserPayload{
	{
		Username:  "alice123",
		Email:     "alice@example.com",
		Password:  "password1",
		FullName:  "Alice Johnson",
		BirthDate: time.Date(1990, 5, 20, 0, 0, 0, 0, time.UTC),
	},
	{
		Username:  "bobsmith",
		Email:     "bob.smith@example.com",
		Password:  "securePass9",
		FullName:  "Bob Smith",
		BirthDate: time.Date(1985, 8, 15, 0, 0, 0, 0, time.UTC),
	},
	{
		Username:  "charlie77",
		Email:     "charlie77@example.net",
		Password:  "myPass1234",
		FullName:  "Charlie Evans",
		BirthDate: time.Date(2000, 12, 1, 0, 0, 0, 0, time.UTC),
	},
	{
		Username:  "danielleX",
		Email:     "danielle@example.org",
		Password:  "passWord99",
		FullName:  "Danielle Brooks",
		BirthDate: time.Date(1995, 3, 10, 0, 0, 0, 0, time.UTC),
	},
	{
		Username:  "evan_j",
		Email:     "evan_j@example.com",
		Password:  "evanPass12",
		FullName:  "Evan Johnson",
		BirthDate: time.Date(1992, 7, 7, 0, 0, 0, 0, time.UTC),
	},
	{
		Username:  "fiona88",
		Email:     "fiona88@example.com",
		Password:  "fionaPass8",
		FullName:  "Fiona Green",
		BirthDate: time.Date(1988, 10, 25, 0, 0, 0, 0, time.UTC),
	},
	{
		Username:  "gregT",
		Email:     "greg.t@example.net",
		Password:  "gregPass99",
		FullName:  "Greg Thompson",
		BirthDate: time.Date(1993, 1, 30, 0, 0, 0, 0, time.UTC),
	},
	{
		Username:  "hannah_b",
		Email:     "hannah.b@example.com",
		Password:  "Hannah1234",
		FullName:  "Hannah Brown",
		BirthDate: time.Date(1999, 11, 12, 0, 0, 0, 0, time.UTC),
	},
	{
		Username:  "ian_miller",
		Email:     "ian.miller@example.org",
		Password:  "IanPass321",
		FullName:  "Ian Miller",
		BirthDate: time.Date(1980, 9, 18, 0, 0, 0, 0, time.UTC),
	},
	{
		Username:  "julia_k",
		Email:     "julia.k@example.com",
		Password:  "JuliaPass5",
		FullName:  "Julia King",
		BirthDate: time.Date(1997, 4, 5, 0, 0, 0, 0, time.UTC),
	},
}

func TestUserService(t *testing.T) {
	if config.Env.Env != "test" {
		log.Panic("environment is not on test!!")
		os.Exit(1)
	}

	db := testutils.SetupTestDB(t)
	manager := db_manager.NewManager(db)

	for _, u := range testUsers {
		customerRole, err := manager.GetRoleByName(types.DefaultRoleCustomer.String())
		if err != nil {
			t.Fatal(err)
		}
		_, err = manager.CreateUser(types.CreateUserPayload{
			Username:  u.Username,
			Email:     u.Email,
			Password:  u.Password,
			FullName:  u.FullName,
			BirthDate: u.BirthDate,
			RoleId:    customerRole.Id,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	ksCache := redis.NewClient(&redis.Options{
		Addr: config.Env.KeyServerRedisAddr,
	})

	authHandlerCache := redis.NewClient(&redis.Options{
		Addr: config.Env.AuthRedisAddr,
	})

	keyServer := auth.NewKeyServer(ksCache)

	authHandler := auth.NewAuthHandler(authHandlerCache, keyServer)

	smtpServer := smtp.NewSMTPServer(
		config.Env.SMTPHost,
		config.Env.SMTPPort,
		config.Env.SMTPEmail,
		config.Env.SMTPPassword,
	)

	kid1 := "kid1"

	handler := NewHandler(manager, authHandler, smtpServer)
	router := mux.NewRouter()

	router.HandleFunc("/login", handler.login).Methods("POST")
	router.HandleFunc("/register/customer", handler.register("Customer")).Methods("POST")
	router.HandleFunc("/register/vendor", handler.register("Vendor")).Methods("POST")
	router.HandleFunc("/refresh", handler.refresh).Methods("POST")

	withAuthRouter := router.Methods("GET", "POST", "PATCH", "DELETE").Subrouter()
	withAuthRouter.HandleFunc("/logout", handler.logout).Methods("POST")
	withAuthRouter.HandleFunc("/me", handler.getMe).Methods("GET")
	withAuthRouter.HandleFunc("/user", handler.getUsers).Methods("GET")
	withAuthRouter.HandleFunc("/user", handler.updateProfile).Methods("PATCH")

	withAuthRouter.Use(authHandler.WithJWTAuth(manager))
	withAuthRouter.Use(authHandler.WithCSRFToken())
	withAuthRouter.Use(authHandler.WithUnbannedProfile(manager))

	t.Run("should initialize key server successfully", func(t *testing.T) {
		err := keyServer.RotateKeys(kid1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should register customer successfully", func(t *testing.T) {
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

		testutils.ExecuteRequest(t, router, req, http.StatusCreated)

		created, err := handler.db.GetUserByUsername(strings.ToLower("testuser1"))
		if err != nil {
			t.Fatal(err)
		}

		if created == nil {
			t.Error("Expected for the created user to exist, but it's not")
		}
	})

	t.Run("should register vendor successfully", func(t *testing.T) {
		payload := types.CreateUserPayload{
			FullName:  "Test Vendor 1",
			Username:  strings.ToLower("testvendor1"),
			Email:     strings.ToLower("testvendor1@gmail.com"),
			BirthDate: time.Date(1989, 2, 25, 0, 0, 0, 0, time.UTC),
			Password:  "password123",
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/register/vendor", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		testutils.ExecuteRequest(t, router, req, http.StatusCreated)

		created, err := handler.db.GetUserByUsername(strings.ToLower("testvendor1"))
		if err != nil {
			t.Fatal(err)
		}

		if created == nil {
			t.Error("Expected for the created user to exist, but it's not")
		}
	})

	t.Run("should fail to register customer due to duplicated username", func(t *testing.T) {
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

		testutils.ExecuteRequest(t, router, req, http.StatusBadRequest)

		created, err := handler.db.GetUserByEmail(
			strings.ToLower("SecondEmailShouldFail@gmail.com"),
		)
		if created != nil {
			t.Error("Expected for the created user to not be found, but it has been found")
		}
	})

	t.Run("should fail to register customer due to duplicated username", func(t *testing.T) {
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

		testutils.ExecuteRequest(t, router, req, http.StatusBadRequest)

		created, err := handler.db.GetUserByUsername(
			strings.ToLower("testuserfailed"),
		)
		if created != nil {
			t.Error("Expected for the created user to not be found, but it has been found")
		}
	})

	t.Run("should fail to register customer because of invalid data", func(t *testing.T) {
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

		testutils.ExecuteRequest(t, router, req, http.StatusBadRequest)

		created, err := handler.db.GetUserByUsername(strings.ToLower("newErrorUser"))
		if created != nil {
			t.Error("Expected for the created user to not be found, but it has been found")
		}
	})

	t.Run("should fail to register customer because of short password", func(t *testing.T) {
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

		testutils.ExecuteRequest(t, router, req, http.StatusBadRequest)

		created, err := handler.db.GetUserByUsername(strings.ToLower("newErrorUser"))
		if created != nil {
			t.Error("Expected for the created user to not be found, but it has been found")
		}
	})

	t.Run("should fail to register customer because of long password", func(t *testing.T) {
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

		testutils.ExecuteRequest(t, router, req, http.StatusBadRequest)

		created, err := handler.db.GetUserByUsername(strings.ToLower("newErrorUser"))
		if created != nil {
			t.Error("Expected for the created user to not be found, but it has been found")
		}
	})

	t.Run("should login successfully", func(t *testing.T) {
		payload := types.LoginUserPayload{
			Username: strings.ToLower("testuser1"),
			Password: "password123",
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		testutils.ExecuteRequest(t, router, req, http.StatusOK)
	})

	t.Run("should fail login because of invalid credentials", func(t *testing.T) {
		payload := types.LoginUserPayload{
			Username: strings.ToLower("testuser1"),
			Password: "wrongpass",
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		testutils.ExecuteRequest(t, router, req, http.StatusBadRequest)
	})

	t.Run("should login and then refresh successfully", func(t *testing.T) {
		authRes := testutils.LoginUser(t, router, strings.ToLower("testuser1"), "password123")

		refreshReq, err := http.NewRequest("POST", "/refresh", nil)
		if err != nil {
			t.Fatal(err)
		}
		for _, cookie := range authRes.Cookies {
			if cookie.Name == "refresh_token" {
				refreshReq.AddCookie(cookie)
			}
		}

		testutils.ExecuteRequest(t, router, refreshReq, http.StatusOK)
	})

	t.Run("should login and then logout successfully", func(t *testing.T) {
		authRes := testutils.LoginUser(t, router, strings.ToLower("testuser1"), "password123")
		logoutReq := testutils.CreateAuthenticatedRequest(t, "POST", "/logout", nil, authRes)
		testutils.ExecuteRequest(t, router, logoutReq, http.StatusOK)
	})

	t.Run("should get users", func(t *testing.T) {
		authRes := testutils.LoginUser(t, router, strings.ToLower("testuser1"), "password123")
		req := testutils.CreateAuthenticatedRequest(t, "GET", "/user", nil, authRes)
		testutils.ExecuteRequest(t, router, req, http.StatusOK)
	})

	t.Run("should update user profile successfully", func(t *testing.T) {
		updatePayload := types.UpdateUserPayload{
			FullName: utils.Ptr("updated name"),
		}

		authRes := testutils.LoginUser(t, router, strings.ToLower("testuser1"), "password123")
		req := testutils.CreateAuthenticatedRequest(t, "PATCH", "/user", updatePayload, authRes)
		testutils.ExecuteRequest(t, router, req, http.StatusOK)
	})
}
