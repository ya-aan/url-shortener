package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestBasicAuth_Unauthorized(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatal("failed to load .env")
	}

	username := os.Getenv("HTTP_SERVER_USER")
	password := os.Getenv("HTTP_SERVER_PASSWORD")

	handler := BasicAuth(username, password)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}

func TestBasicAuth_Success(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatal("failed to load .env")
	}

	username := os.Getenv("HTTP_SERVER_USER")
	password := os.Getenv("HTTP_SERVER_PASSWORD")

	handler := BasicAuth(username, password)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.SetBasicAuth(username, password)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}
}
