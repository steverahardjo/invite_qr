package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockAuthenticator struct {
	token string
	err   error
}

func (m *mockAuthenticator) LoginAdmin(_ context.Context, _, _ string) (string, error) {
	return m.token, m.err
}

func TestLoginAdmin_Success(t *testing.T) {
	mock := &mockAuthenticator{token: "test-jwt-token"}
	h := NewJwtHandler(mock)

	body := `{"username":"admin","password":"secret"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.LoginAdmin().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["token"] != "test-jwt-token" {
		t.Fatalf("expected token 'test-jwt-token', got '%s'", resp["token"])
	}
}

func TestLoginAdmin_InvalidJSON(t *testing.T) {
	h := NewJwtHandler(&mockAuthenticator{})

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{bad json}`))
	rec := httptest.NewRecorder()

	h.LoginAdmin().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestLoginAdmin_ServiceError(t *testing.T) {
	mock := &mockAuthenticator{err: errors.New("invalid credentials")}
	h := NewJwtHandler(mock)

	body := `{"username":"admin","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.LoginAdmin().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}
