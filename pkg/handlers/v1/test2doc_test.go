package v1

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/adams-sarah/test2doc/doc/parse"
	"github.com/adams-sarah/test2doc/test"
	"github.com/takama/bit"
	"github.com/gromnsk/money.io/pkg/assets"
	"github.com/gromnsk/money.io/pkg/config"
	"github.com/gromnsk/money.io/pkg/logger"
	"github.com/gromnsk/money.io/pkg/logger/standard"
)

var server *test.Server

func TestMain(m *testing.M) {
	cfg := new(config.Config)

	// Setup logger
	log := standard.New(&logger.Config{
		Level: cfg.LogLevel,
		Time:  true,
		UTC:   true,
	})

	// Define handlers
	h := New(log, cfg)

	// Register new router
	r := bit.NewRouter()

	// Asset manager
	assetManager := assets.NewManagerGenerated()

	// Configure router
	r.SetupMiddleware(h.Base)
	r.GET("/", h.Root)
	r.GET("/healthz", h.Health)
	r.GET("/readyz", h.Ready)
	r.GET("/info", h.Info)
	r.GET("/metrics", h.MetricsFunc())
	r.GET("/blueprints/:name", h.GetBlueprint(assetManager))
	r.GET("/blueprints", h.Blueprints(assetManager))

	test.RegisterURLVarExtractor(BitURLVarExtractor(r))

	var err error
	server, err = test.NewServer(r)
	if err != nil {
		panic(err.Error())
	}

	exitCode := m.Run()
	server.Finish()
	os.Exit(exitCode)
}

// BitURLVarExtractor returns a func which extracts
// url vars from a request for test2doc documentation generation.
func BitURLVarExtractor(router bit.Router) parse.URLVarExtractor {
	return func(req *http.Request) map[string]string {
		// router Lookup func needs a trailing slash on path
		path := req.URL.Path

		if !strings.HasSuffix(path, "/") {
			path += "/"
		}

		_, params, ok := router.Lookup(req.Method, path)
		if !ok {
			return nil
		}

		paramsMap := make(map[string]string, len(params))
		for _, p := range params {
			paramsMap[p.Key] = p.Value
		}
		return paramsMap
	}
}

// TestHandlersCalls makes GET request to each handler http path for test2doc documentation generation
func TestHandlersCalls(t *testing.T) {
	tr := &http.Transport{
		// Because test2doc does not understand correctly gzip encoding
		DisableCompression: true,
	}
	c := &http.Client{Transport: tr}

	ps := []string{"/", "/healthz", "/readyz", "/version", "/info", "/blueprints", "/blueprints/test"}
	for _, p := range ps {

		req, err := http.NewRequest("GET", server.URL+p, nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := c.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		resp.Body.Close()
	}
}
