package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/takama/bit"
	"github.com/gromnsk/money.io/pkg/assets"
	"github.com/gromnsk/money.io/pkg/config"
	"github.com/gromnsk/money.io/pkg/errors"
	"github.com/gromnsk/money.io/pkg/logger"
	"github.com/gromnsk/money.io/pkg/logger/standard"
)

func TestBlueprints(t *testing.T) {
	cases := []struct {
		manager  assets.AssetManager
		status   int
		response Response
	}{
		{
			manager: assets.AssetManagerMock{
				OnAssetNames: func() []string {
					return []string{"pkg/v1/v1.apib", "pkg/v2/v2.apib"}
				},
			},
			status: http.StatusOK,
			response: Response{
				Data: []string{"pkg-v1-v1.apib", "pkg-v2-v2.apib"},
			},
		},
		{
			manager: assets.AssetManagerMock{
				OnAssetNames: func() []string {
					return nil
				},
			},
			status: http.StatusNotFound,
			response: Response{
				Error: &errors.ResponseError{
					Code:    http.StatusNotFound,
					Message: "Blueprints not found",
				},
			},
		},
	}

	for _, c := range cases {
		h := New(standard.New(&logger.Config{}), new(config.Config))
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.Base(h.Blueprints(c.manager))(bit.NewControl(w, r))
		})

		b, err := json.Marshal(c.response)
		if err != nil {
			t.Error(err)
		}
		testHandler(t, handler, c.status, string(b))
	}
}

func TestGetBlueprint(t *testing.T) {
	cases := []struct {
		manager   assets.AssetManager
		blueprint string
		status    int
	}{
		{
			manager: assets.AssetManagerMock{
				OnAsset: func(name string) ([]byte, error) {
					return nil, fmt.Errorf("Asset %s not found", name)
				},
			},
			blueprint: "pkg-v1-v1.apib",
			status:    http.StatusNotFound,
		},
		{
			manager: assets.AssetManagerMock{
				OnAsset: func(name string) ([]byte, error) {
					if name == "pkg/v1/v1.apib" {
						return []byte{}, nil
					}
					return nil, fmt.Errorf("Asset %s not found", name)
				},
			},
			blueprint: "pkg-v1-v1.apib",
			status:    http.StatusOK,
		},
	}

	for _, c := range cases {
		h := New(standard.New(&logger.Config{}), new(config.Config))
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			control := bit.NewControl(w, r)
			control.Param(":name", c.blueprint)
			h.Base(h.GetBlueprint(c.manager))(control)
		})

		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Error(err)
		}

		trw := httptest.NewRecorder()
		handler.ServeHTTP(trw, req)

		if trw.Code != c.status {
			t.Error("Expected status code:", c.status, "got", trw.Code)
		}
	}
}
