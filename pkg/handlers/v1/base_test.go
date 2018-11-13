package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/takama/bit"
	"github.com/gromnsk/money.io/pkg/config"
	"github.com/gromnsk/money.io/pkg/errors"
	"github.com/gromnsk/money.io/pkg/logger"
	"github.com/gromnsk/money.io/pkg/logger/standard"
)

func TestReturnSuccess(t *testing.T) {
	h := New(standard.New(&logger.Config{}), new(config.Config))

	body := "test body"
	trw := httptest.NewRecorder()
	control := bit.NewControl(trw, nil)
	h.ReturnSuccess(control, http.StatusOK, body)

	if trw.Code != http.StatusOK {
		t.Error("Expected status code:", http.StatusOK, "got", trw.Code)
	}
	response := Response{
		Data: body,
	}
	b, _ := json.Marshal(response)
	if trw.Body.String() != string(b) {
		t.Error("Expected body", string(b), "got", trw.Body.String())
	}

	// error case
	trw = httptest.NewRecorder()
	control = bit.NewControl(trw, nil)
	h.ReturnSuccess(control, http.StatusOK, make(chan struct{}))
}

func TestReturnError(t *testing.T) {
	h := New(standard.New(&logger.Config{}), new(config.Config))

	responseErr := &errors.ResponseError{Code: http.StatusNotFound, Message: "Not found"}
	trw := httptest.NewRecorder()
	control := bit.NewControl(trw, nil)
	h.ReturnError(control, http.StatusNotFound, responseErr)

	if trw.Code != http.StatusNotFound {
		t.Error("Expected status code:", http.StatusOK, "got", trw.Code)
	}
	response := Response{
		Error: responseErr,
	}
	b, _ := json.Marshal(response)
	if trw.Body.String() != string(b) {
		t.Error("Expected body", string(b), "got", trw.Body.String())
	}
}

func TestPrepareBadRequestError(t *testing.T) {
	h := New(standard.New(&logger.Config{}), new(config.Config))

	responseErr := &errors.ResponseError{Code: http.StatusBadRequest, Message: "Not found"}
	trw := httptest.NewRecorder()
	control := bit.NewControl(trw, nil)
	h.prepareBadRequestError(control, responseErr.Message)

	if trw.Code != http.StatusBadRequest {
		t.Error("Expected status code:", http.StatusBadRequest, "got", trw.Code)
	}
	response := Response{
		Error: responseErr,
	}
	b, _ := json.Marshal(response)
	if trw.Body.String() != string(b) {
		t.Error("Expected body", string(b), "got", trw.Body.String())
	}
}

func TestPrepareServerError(t *testing.T) {
	h := New(standard.New(&logger.Config{}), new(config.Config))

	responseErr := &errors.ResponseError{Code: http.StatusServiceUnavailable, Message: "Not found"}
	trw := httptest.NewRecorder()
	control := bit.NewControl(trw, nil)
	h.prepareServerError(control, responseErr.Message)

	if trw.Code != http.StatusServiceUnavailable {
		t.Error("Expected status code:", http.StatusServiceUnavailable, "got", trw.Code)
	}
	response := Response{
		Error: responseErr,
	}
	b, _ := json.Marshal(response)
	if trw.Body.String() != string(b) {
		t.Error("Expected body", string(b), "got", trw.Body.String())
	}
}

func TestUnmarshalRequestBody(t *testing.T) {
	h := New(standard.New(&logger.Config{}), new(config.Config))
	trw := httptest.NewRecorder()

	type test struct {
		A string
		B int
	}

	// without error
	toRequest := test{"123", 123}
	b, _ := json.Marshal(toRequest)
	req, err := http.NewRequest("GET", "/", bytes.NewBuffer(b))
	if err != nil {
		t.Error(err)
	}
	control := bit.NewControl(trw, req)
	var fromRequest test
	responseErr := h.UnmarshalRequestBody(control, &fromRequest)
	if responseErr != nil {
		t.Error("Expected error", nil, "got", responseErr)
	}
	if fromRequest.A != toRequest.A || fromRequest.B != toRequest.B {
		t.Error("Expected objects to be equal:", fromRequest, toRequest)
	}

	// unmarshal error
	b, _ = json.Marshal(toRequest)
	req, err = http.NewRequest("GET", "/", bytes.NewBufferString("{\"A\": 123}"))
	if err != nil {
		t.Error(err)
	}
	control = bit.NewControl(trw, req)
	responseErr = h.UnmarshalRequestBody(control, &fromRequest)
	if responseErr == nil || responseErr.Code != errors.ErrorCodeUnknown {
		t.Error("Expected error, got", responseErr)
	}
}
