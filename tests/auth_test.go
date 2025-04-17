package main

import (
	"api-gw/jwt"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJWTAuth(t *testing.T) {

	handler := jwt.Auth(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}, true)

	t.Run("No Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected: %d | Received: %v", http.StatusUnauthorized, rr.Code)
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer test")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusForbidden {
			t.Errorf("Expected: %d | Received: %v", http.StatusForbidden, rr.Code)
		}
	})

	t.Run("Valid Token", func(t *testing.T) {

		token, _ := jwt.CreateToken()
		req := httptest.NewRequest("GET", "/api/metrics", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected: %d | Received: %v", http.StatusCreated, rr.Code)
		}
	})
}
