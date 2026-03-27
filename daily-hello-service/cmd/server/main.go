package main

import (
	"context"
	"daily-hello-service/cmd/migration"
	"daily-hello-service/config"
	"daily-hello-service/internal/diregistry"
	"daily-hello-service/internal/route"
	"errors"
	"fmt"
	"go-libs/binder"
	"go-libs/errorhelper"
	"go-libs/validatehelper"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/labstack/echo-contrib/pprof"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	_ "daily-hello-service/swagger-docs"
	"go-libs/loghelper"

	echo "github.com/labstack/echo/v4"
)

var (
	cfg *config.Config

	// Build-time variables injected via -ldflags
	Version   string = "unknown"
	Branch    string = "unknown"
	BuildTime string = "unknown"
)

// @title Swagger Example API
// @version 1.0
// @description This is HDBank Agent Management (HAM) server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.hdbank.com.vn/
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8282
// @BasePath /api/
// @securityDefinitions.apikey BearerTokenAuth
// @in header
// @name Authorization

// @Security BearerTokenAuth
func main() {
	// Getting configuration base on environment
	diregistry.BuildDIContainer()
	cfg = diregistry.GetDependency(diregistry.ConfigDIName).(*config.Config)

	_, err := loghelper.InitZapLogger(&loghelper.LoggerOptions{
		AppName:       cfg.App,
		MaskingFields: cfg.SensitiveFields,
		LogLevel:      loghelper.LogLevel(cfg.LogLevel),
	})
	if err != nil {
		loghelper.Logger.Panic("Can't init zap logger", loghelper.Error(err))
	}

	loghelper.Logger.Info("*****STARTING*****")

	migration.StartMigrate(cfg)

	httpServer := initAndStartHttpServer()
	loghelper.Logger.Info(fmt.Sprintf("Gateway is started on port %d", cfg.HttpAddress))

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	<-ctx.Done()
	loghelper.Logger.Info("*****GRACEFUL SHUTTING DOWN*****")
	switch cfg.Env {
	case "prd":
		shutdown(httpServer, 30*time.Second)
	case "dev":
		shutdown(httpServer, 30*time.Second)
	}
	loghelper.Logger.Info("*****SHUTDOWN*****")
}

func initAndStartHttpServer() *echo.Echo {
	httpServer := echo.New()

	httpServer.GET("/liveness", liveness)
	httpServer.GET("/readiness", readiness)

	httpServer.Validator = validatehelper.NewValidator()
	httpServer.HTTPErrorHandler = errorhelper.NewEchoErrorHandler(httpServer).Handle
	httpServer.Binder = binder.NewBinder()
	allowOrigins := strings.Split(cfg.AllowOrigins, ",")
	allowHeaders := []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-Timestamp", "X-Request-Signature"}
	httpServer.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{LogLevel: log.ERROR}), middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "PATCH", "HEAD"},
		AllowHeaders:     allowHeaders,
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           86400,
	}))
	httpServer.Use(echoprometheus.NewMiddleware("myapp")) // adds middleware to gather metrics
	httpServer.Use(loghelper.Logger.CorrelationIDMiddleware())
	httpServer.Use(loghelper.Logger.NewEchoLoggerMiddleware())
	httpServer.Use(loghelper.Logger.BodyDump())

	httpServer.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics
	pprof.Register(httpServer)

	httpServer.Server.Addr = fmt.Sprintf(":%d", cfg.HttpAddress)
	httpServer.Server.ReadTimeout = 10 * time.Minute
	httpServer.Server.WriteTimeout = 5 * time.Minute

	route.RegisterRoutes(httpServer)
	startHttpServer(httpServer)

	return httpServer
}

func startHttpServer(httpServer *echo.Echo) {
	go func() {
		if err := httpServer.StartServer(httpServer.Server); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				httpServer.Logger.Info("shutting down the server")
			} else {
				httpServer.Logger.Errorf("error shutting down the server: ", err)
			}
		}
	}()
}

func shutdown(e *echo.Echo, gracefulShutdownTime time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTime)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func liveness(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func readiness(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":     "OK",
		"version":    Version,
		"branch":     Branch,
		"build_time": BuildTime,
	})
}
