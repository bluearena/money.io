package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/takama/bit"
	"github.com/gromnsk/money.io/pkg/assets"
	"github.com/gromnsk/money.io/pkg/errors"
)

// Blueprints returns list of available API blueprints.
func (h *Handler) Blueprints(manager assets.AssetManager) func(bit.Control) {
	return func(c bit.Control) {
		assetNames := manager.AssetNames()
		blueprints := make([]string, 0)
		for _, a := range assetNames {
			if strings.HasSuffix(a, ".apib") {
				blueprints = append(blueprints, strings.Replace(a, "/", "-", -1))
			}
		}
		if len(blueprints) == 0 {
			respErr := errors.NewResponseError(http.StatusNotFound, "Blueprints not found")
			h.ReturnError(c, http.StatusNotFound, respErr)
			return
		}
		h.ReturnSuccess(c, http.StatusOK, blueprints)
	}
}

// GetBlueprint returns API blueprint based on its name.
func (h *Handler) GetBlueprint(manager assets.AssetManager) func(bit.Control) {
	return func(c bit.Control) {
		name := strings.Replace(c.Query(":name"), "-", "/", -1)
		content, err := manager.Asset(fmt.Sprintf("%s", name))
		if err != nil {
			c.WriteHeader(http.StatusNotFound)
			return
		}

		c.Header().Set(HeaderContentType, CTApiB)
		c.WriteHeader(http.StatusOK)
		_, err = c.Write(content)
		if err != nil {
			h.logger.Error(err)
		}
	}
}
