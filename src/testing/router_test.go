package testing

import (
	"net/http"
	"net/http/httptest"
	"src/service"
	"testing"
)

func TestGet(t *testing.T) {
	service.EngineConfig()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/get/1233232", nil)
	service.App.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Error()
	}
}

func TestError(t *testing.T) {
	service.EngineConfig()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/err", nil)
	service.App.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Error()
	}
}

func TestPanic(t *testing.T) {
	service.EngineConfig()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)
	service.App.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Error()
	}
}
