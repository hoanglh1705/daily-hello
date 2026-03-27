package route

import (
	"daily-hello-service/config"
	"daily-hello-service/internal/diregistry"
	"daily-hello-service/internal/handlers"

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
	apiGroup := httpServer.Group("/api")
	v1 := apiGroup.Group("/v1")
	registerAuthRoutes(v1)
}

func registerPublicRoutes(g *echo.Group) {

}

func registerAuthRoutes(g *echo.Group) {
	authGroup := g.Group("/auth")

	// Branch routes
	branchHandler := diregistry.GetDependency(diregistry.BranchAPIDIName).(*handlers.BranchHandler)
	branchGroup := authGroup.Group("/branches")
	branchGroup.POST("", branchHandler.Create)
	branchGroup.GET("/:id", branchHandler.GetByID)
	branchGroup.PUT("/:id", branchHandler.Update)
	branchGroup.DELETE("/:id", branchHandler.Delete)
	branchGroup.GET("", branchHandler.List)

	// Branch Wifi routes
	branchWifiHandler := diregistry.GetDependency(diregistry.BranchWifiAPIDIName).(*handlers.BranchWifiHandler)
	branchWifiGroup := authGroup.Group("/branch-wifi")
	branchWifiGroup.POST("", branchWifiHandler.Create)
	branchWifiGroup.GET("/:id", branchWifiHandler.GetByID)
	branchWifiGroup.GET("/branch/:branch_id", branchWifiHandler.GetByBranchID)
	branchWifiGroup.PUT("/:id", branchWifiHandler.Update)
	branchWifiGroup.DELETE("/:id", branchWifiHandler.Delete)
}
