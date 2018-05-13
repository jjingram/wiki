package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerGet(t *testing.T) {
	req, err := http.NewRequest("GET", "/foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := Handler{db: nil}

	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "foo"
	if rec.Body.String() != expected {
		t.Errorf("handler return unexpected body: got %v want %v", rec.Body.String(), expected)
	}
}

func TestHandlerPost(t *testing.T) {
	req, err := http.NewRequest("POST", "/foo", strings.NewReader("foo"))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := Handler{db: nil}

	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "foo"
	if rec.Body.String() != expected {
		t.Errorf("handler return unexpected body: got %v want %v", rec.Body.String(), expected)
	}
}

func TestHandlerMethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest("PUT", "/foo", strings.NewReader("foo"))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := Handler{db: nil}

	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := http.StatusText(http.StatusMethodNotAllowed)
	if rec.Body.String() != expected {
		t.Errorf("handler return unexpected body: got %v want %v", rec.Body.String(), expected)
	}
}
