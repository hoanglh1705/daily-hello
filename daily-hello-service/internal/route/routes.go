package route

import (
	"daily-hello-service/config"
	"daily-hello-service/internal/diregistry"
	"daily-hello-service/internal/handlers"

	auth_middleware "go-libs/http_middlewares/auth"

	echo "github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/thoas/go-funk"
)

func RegisterRoutes(httpServer *echo.Echo) {
	cfg := diregistry.GetDependency(diregistry.ConfigDIName).(*config.Config)
	if funk.Contains(config.NonProductionEnvironments, cfg.Env) {
		httpServer.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	apiGroup := httpServer.Group("/api")
	v1 := apiGroup.Group("/v1")

	registerPublicRoutes(v1)
	registerProtectedRoutes(v1)
}

func registerPublicRoutes(g *echo.Group) {
	authHandler := diregistry.GetDependency(diregistry.AuthAPIDIName).(*handlers.AuthHandler)

	// POST /api/v1/auth/login
	// POST /api/v1/auth/refresh-token
	authGroup := g.Group("/auth")
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/logout", authHandler.Logout)
	authGroup.POST("/refresh-token", authHandler.RefreshToken)
}

func registerProtectedRoutes(g *echo.Group) {
	jwtMiddleware := diregistry.GetDependency(diregistry.JWTMiddlewareDIName).(*auth_middleware.Service)
	g.Use(jwtMiddleware.MWFunc())

	// Auth routes (requires JWT)
	// POST /api/v1/auth/logout
	authHandler := diregistry.GetDependency(diregistry.AuthAPIDIName).(*handlers.AuthHandler)
	authGroup := g.Group("/auth")
	authGroup.POST("/logout", authHandler.Logout)

	// User routes
	userHandler := diregistry.GetDependency(diregistry.UserAPIDIName).(*handlers.UserHandler)
	userGroup := g.Group("/users")
	userGroup.POST("", userHandler.Register)
	userGroup.GET("", userHandler.List)
	userGroup.GET("/me", userHandler.GetMe)
	userGroup.GET("/:id", userHandler.GetByID)
	userGroup.PUT("/:id", userHandler.Update)

	// Branch routes
	branchHandler := diregistry.GetDependency(diregistry.BranchAPIDIName).(*handlers.BranchHandler)
	branchGroup := g.Group("/branches")
	branchGroup.POST("", branchHandler.Create)
	branchGroup.GET("/:id", branchHandler.GetByID)
	branchGroup.PUT("/:id", branchHandler.Update)
	branchGroup.DELETE("/:id", branchHandler.Delete)
	branchGroup.GET("", branchHandler.List)

	// Branch Wifi routes
	branchWifiHandler := diregistry.GetDependency(diregistry.BranchWifiAPIDIName).(*handlers.BranchWifiHandler)
	branchWifiGroup := g.Group("/branch-wifi")
	branchWifiGroup.POST("", branchWifiHandler.Create)
	branchWifiGroup.GET("/:id", branchWifiHandler.GetByID)
	branchWifiGroup.GET("/branch/:branch_id", branchWifiHandler.GetByBranchID)
	branchWifiGroup.PUT("/:id", branchWifiHandler.Update)
	branchWifiGroup.DELETE("/:id", branchWifiHandler.Delete)

	// Attendance routes
	attendanceHandler := diregistry.GetDependency(diregistry.AttendanceAPIDIName).(*handlers.AttendanceHandler)
	attendanceGroup := g.Group("/attendance")
	attendanceGroup.POST("/check-in", attendanceHandler.CheckIn)
	attendanceGroup.POST("/check-out", attendanceHandler.CheckOut)
	attendanceGroup.GET("/history", attendanceHandler.GetHistory)
	attendanceGroup.GET("/today", attendanceHandler.GetToday)
	attendanceGroup.GET("/:id", attendanceHandler.GetByID)

	// Device routes (user)
	deviceHandler := diregistry.GetDependency(diregistry.DeviceAPIDIName).(*handlers.DeviceHandler)
	deviceGroup := g.Group("/devices")
	deviceGroup.POST("/register", deviceHandler.Register)
	deviceGroup.GET("/status", deviceHandler.GetStatus)
	deviceGroup.GET("", deviceHandler.ListMyDevices)

	// Device admin routes
	adminGroup := g.Group("/admin")
	adminGroup.GET("/devices", deviceHandler.AdminList)
	adminGroup.PUT("/devices/:id/approve", deviceHandler.Approve)
	adminGroup.PUT("/devices/:id/reject", deviceHandler.Reject)
}
