package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthGet(t *testing.T) {
	handler := newServer()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	if contentType := response.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "text/plain") {
		t.Fatalf("got Content-Type %q, want text/plain", contentType)
	}
	if got := response.Body.String(); got != "ok\n" {
		t.Fatalf("got body %q, want %q", got, "ok\n")
	}

}

func TestHealthPOST(t *testing.T) {
	handler := newServer()
	request := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("got status %d, want %d", response.Code, http.StatusMethodNotAllowed)
	}

	if got := response.Header().Get("Allow"); got != http.MethodGet {
		t.Fatalf("got Allow %q, want %q", got, http.MethodGet)
	}
	if got := response.Body.String(); got != "Method not allowed\n" {
		t.Fatalf("got body %q, want %q", got, "Method not allowed\n")
	}
}
