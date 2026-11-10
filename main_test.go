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

// No samples yet averages to 0 instead of dividing by zero.
func TestAverageLatencyNoSamples(t *testing.T) {
	if got := averageLatency(nil); got != 0 {
		t.Errorf("averageLatency(nil) = %d, want 0", got)
	}
}

// GET /test answers 200 with an average of 0 while nothing has been recorded.
func TestLatencyReport(t *testing.T) {
	rec := httptest.NewRecorder()
	recoverer(http.HandlerFunc(latencyReport)).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/test", nil))
	if rec.Code != http.StatusOK || strings.TrimSpace(rec.Body.String()) != `{"average_ms":0}` {
		t.Errorf("GET /test = %d %s", rec.Code, rec.Body)
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
