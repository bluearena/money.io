package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/takama/bit"
	"github.com/gromnsk/money.io/pkg/errors"
)

// UnmarshalRequestBody reads export context body to provided structure
func (h *Handler) UnmarshalRequestBody(ctrl bit.Control, object interface{}) *errors.ResponseError {
	defer func() {
		_ = ctrl.Request().Body.Close()
	}()

	requestBody, err := ioutil.ReadAll(ctrl.Request().Body)
	if err != nil {
		return errors.NewResponseError(errors.ErrorCodeUnknown, "Request body error: %"+err.Error())
	}

	err = json.Unmarshal(requestBody, object)
	if err != nil {
		h.logger.Errorf("Error: Can't parse request (%s)", err)
		return errors.NewResponseError(errors.ErrorCodeUnknown, "Request body error: %"+err.Error())
	}

	return nil
}

// ReturnSuccess provides JSON API response in case of success.
func (h *Handler) ReturnSuccess(ctrl bit.Control, code int, response interface{}) {
	var b []byte
	var err error

	resp := Response{
		Data: response,
	}

	b, err = json.Marshal(resp)
	if err != nil {
		h.logger.Error(err)
		return
	}

	ctrl.Header().Set(HeaderContentType, CTApplicationJSON)
	ctrl.WriteHeader(code)

	_, err = ctrl.Write(b)
	if err != nil {
		h.logger.Error(err)
	}
}

// ReturnError provides JSON API method in case of user or application error.
func (h *Handler) ReturnError(ctrl bit.Control, code int, response *errors.ResponseError) {
	var b []byte
	var err error

	resp := Response{
		Error: response,
	}

	b, err = json.Marshal(resp)
	if err != nil {
		h.logger.Error(err)
		return
	}

	ctrl.Header().Set(HeaderContentType, CTApplicationJSON)
	ctrl.WriteHeader(code)

	_, err = ctrl.Write(b)
	if err != nil {
		h.logger.Error(err)
	}
}

func (h *Handler) prepareBadRequestError(c bit.Control, text string) {
	extended := &errors.ResponseError{
		Code:    http.StatusBadRequest,
		Message: text,
	}
	h.ReturnError(c, http.StatusBadRequest, extended)
}

func (h *Handler) prepareServerError(c bit.Control, text string) {
	extended := &errors.ResponseError{
		Code:    http.StatusServiceUnavailable,
		Message: text,
	}
	h.ReturnError(c, http.StatusServiceUnavailable, extended)
}