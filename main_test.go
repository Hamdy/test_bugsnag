package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	recoverer(http.HandlerFunc(status)).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusOK || strings.TrimSpace(rec.Body.String()) != `{"status":"ok"}` {
		t.Errorf("GET / = %d %s", rec.Code, rec.Body)
	}
}

func TestAverageLatency(t *testing.T) {
	if got := averageLatency([]int{10, 20, 30}); got != 20 {
		t.Errorf("averageLatency = %d, want 20", got)
	}
}

// A panic in a handler is answered with a 500, not a dropped connection.
func TestRecovererAnswers500(t *testing.T) {
	rec := httptest.NewRecorder()
	recoverer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("boom") })).
		ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))
	if rec.Code != http.StatusInternalServerError || !strings.Contains(rec.Body.String(), "boom") {
		t.Errorf("a panic = %d %s", rec.Code, rec.Body)
	}
}
