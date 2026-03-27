package route

import (
	"daily-hello-service/config"
	"daily-hello-service/internal/diregistry"

	echo "github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/thoas/go-funk"
)

func RegisterRoutes(httpServer *echo.Echo) {
	cfg := diregistry.GetDependency(diregistry.ConfigDIName).(*config.Config)
	if funk.Contains(config.NonProductionEnvironments, cfg.Env) {
		// use echoSwagger middleware to serve the API docs
		httpServer.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// Init route
	// Init API v1
	// e := httpServer.Group("/api")

}

func registerPublicRoutes(g *echo.Group) {

}

func registerAuthRoutes(g *echo.Group) {

}
