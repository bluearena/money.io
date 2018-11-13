package app

import (
	"net/http"

	"github.com/takama/bit"
	"github.com/gromnsk/money.io/pkg/assets"
	"github.com/gromnsk/money.io/pkg/handlers/v1"
	"github.com/gromnsk/money.io/pkg/logger"
)

type App struct {
	logger   logger.Logger
	handlers *v1.Handler
	router   bit.Router
}

func NewApp(logger logger.Logger, handlers *v1.Handler, router bit.Router) *App {
	return &App{
		logger:   logger,
		handlers: handlers,
		router:   router,
	}
}

// SetupDefaultRoutes configures the default routes of the service
func (a *App) SetupDefaultRoutes() {
	// Response for undefined methods
	a.router.SetupNotFoundHandler(a.handlers.Base(notFound))

	// Configure router
	a.router.SetupMiddleware(a.handlers.Base)
	a.router.GET("/", a.handlers.Root)
	a.router.GET("/healthz", a.handlers.Health)
	a.router.GET("/readyz", a.handlers.Ready)
	a.router.GET("/info", a.handlers.Info)
	a.router.GET("/metrics", a.handlers.MetricsFunc())
}

// SetupBlueprints configures the blueprints routes of the service for documentation
func (a *App) SetupBlueprints(assetManager assets.AssetManager) {
	a.router.GET("/blueprints/:name", a.handlers.GetBlueprint(assetManager))
	a.router.GET("/blueprints", a.handlers.Blueprints(assetManager))
}

// SetupAppRoutes configures the app related routes of the service
func (a *App) SetupAppRoutes() {
	// Add your app routes here.
	// Feel free to add arguments (bundle, storage, etc. for ex.).
}

// Response for undefined methods
func notFound(c bit.Control) {
	c.Code(http.StatusNotFound)
	c.Body("Method not found for " + c.Request().URL.Path)
}
