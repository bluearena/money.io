package v1

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/takama/bit"
	"github.com/gromnsk/money.io/pkg/config"
	"github.com/gromnsk/money.io/pkg/logger"
	"github.com/gromnsk/money.io/pkg/logger/standard"
)

func TestMetrics(t *testing.T) {
	h := New(standard.New(&logger.Config{}), new(config.Config))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(h.Root)(bit.NewControl(w, r))
	})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	trw := httptest.NewRecorder()
	handler.ServeHTTP(trw, req)

	metricsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(h.MetricsFunc())(bit.NewControl(w, r))
	})

	req, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	trw = httptest.NewRecorder()
	metricsHandler.ServeHTTP(trw, req)

	metrics := trw.Body.String()

	if !strings.Contains(metrics, "http_request_duration_seconds") {
		t.Fatalf("Cannot find metrics of request durations for service")
	}

	if !strings.Contains(metrics, "http_requests_total") {
		t.Fatalf("Cannot find metrics of response statuses for service")
	}
}
