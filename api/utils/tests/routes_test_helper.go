package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SaeedAlian/econest/api/types"
)

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
	Cookies     []*http.Cookie
}

func LoginUser(t *testing.T, handler http.Handler, username string, password string) AuthResponse {
	t.Helper()

	loginPayload := types.LoginUserPayload{
		Username: username,
		Password: password,
	}

	loginBody, err := json.Marshal(loginPayload)
	if err != nil {
		t.Fatalf("Failed to marshal login payload: %v", err)
	}

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Login failed with status %d: %s", rr.Code, rr.Body.String())
	}

	var authResp AuthResponse
	if err := json.NewDecoder(rr.Body).Decode(&authResp); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	authResp.Cookies = rr.Result().Cookies()
	return authResp
}

func CreateAuthenticatedRequest(
	t *testing.T,
	method string,
	path string,
	body any,
	authResp AuthResponse,
) *http.Request {
	t.Helper()

	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(bodyBytes))

	var csrfToken string
	for _, cookie := range authResp.Cookies {
		req.AddCookie(cookie)
		if cookie.Name == "csrf_token" {
			csrfToken = cookie.Value
		}
	}

	req.Header.Set("Authorization", "Bearer "+authResp.AccessToken)
	if csrfToken != "" {
		req.Header.Set("X-CSRF-Token", csrfToken)
	}

	return req
}

func ExecuteRequest(
	t *testing.T,
	handler http.Handler,
	req *http.Request,
	expectedStatus int,
) *httptest.ResponseRecorder {
	t.Helper()

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != expectedStatus {
		t.Errorf("Expected status %d, got %d. Response: %s",
			expectedStatus, rr.Code, rr.Body.String())
	}

	return rr
}
