package v1

import (
	"net/http"
	"testing"
	
	"github.com/takama/bit"
	"github.com/gromnsk/money.io/pkg/config"
	"github.com/gromnsk/money.io/pkg/logger"
	"github.com/gromnsk/money.io/pkg/logger/standard"
)

func TestHealth(t *testing.T) {
	h := New(standard.New(&logger.Config{}), new(config.Config))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(h.Health)(bit.NewControl(w, r))
	})

	testHandler(t, handler, http.StatusOK, http.StatusText(http.StatusOK))
}
