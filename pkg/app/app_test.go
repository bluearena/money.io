package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/takama/bit"
	"github.com/gromnsk/money.io/pkg/config"
	"github.com/gromnsk/money.io/pkg/handlers/v1"
	pkglog "github.com/gromnsk/money.io/pkg/logger"
	stdlog "github.com/gromnsk/money.io/pkg/logger/standard"
)

func TestNewApp(t *testing.T) {
	cfg := new(config.Config)
	err := cfg.Load(config.SERVICENAME)
	if err != nil {
		t.Error("Expected loading of environment vars, got", err)
	}

	logger := stdlog.New(&pkglog.Config{
		Level: cfg.LogLevel,
		Time:  true,
		UTC:   true,
	})

	// Create router
	router := bit.NewRouter()

	// Create handlers
	h := v1.New(logger, cfg)

	// Create app
	a := NewApp(logger, h, router)

	a.SetupDefaultRoutes()
	a.SetupAppRoutes()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(notFound)(bit.NewControl(w, r))
	})

	req, err := http.NewRequest("GET", "/notfound", nil)
	if err != nil {
		t.Error(err)
	}

	trw := httptest.NewRecorder()
	handler.ServeHTTP(trw, req)

	if trw.Code != http.StatusNotFound {
		t.Error("Expected status:", http.StatusNotFound, "got", trw.Code)
	}
}
