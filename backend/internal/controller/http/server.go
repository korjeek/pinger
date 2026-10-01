package http

import (
	"log/slog"
	"net/http"

	api2 "github.com/korjeek/pinger/backend/internal/controller/http/api"
	mdl "github.com/korjeek/pinger/backend/internal/controller/http/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ServerConfig struct {
	Logger *slog.Logger

	AccessTokenParser AccessTokenParser

	CookiePath   string
	CookieDomain string
	CookieSecure bool
}

func NewServer(cfg ServerConfig, impl api2.StrictServerInterface) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(TraceIDMiddleware())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:  true,
		LogURI:     true,
		LogStatus:  true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			cfg.Logger.Info("request",
				"method", v.Method,
				"uri", v.URI,
				"status", v.Status,
				"latency_ms", v.Latency.Milliseconds(),
				"trace_id", TraceIDFromEcho(c),
			)
			return nil
		},
	}))
	e.Use(RefreshTokenCookieMiddleware())

	auth := AuthMiddleware(cfg.AccessTokenParser)

	opMW := map[string][]echo.MiddlewareFunc{
		"updateUser":              {auth},
		"getMonitors":             {auth},
		"createMonitor":           {auth},
		"deleteMonitor":           {auth},
		"updateMonitor":           {auth},
		"getSnapshotsByMonitorID": {auth},
	}

	strictMW := []api2.StrictMiddlewareFunc{
		mdl.StrictErrorMiddleware,
		RefreshCookieWriterMiddleware(cookieConfig{
			Path:     cfg.CookiePath,
			Domain:   cfg.CookieDomain,
			Secure:   cfg.CookieSecure,
			SameSite: http.SameSiteStrictMode,
		}),
	}

	handler := api2.NewStrictHandler(impl, strictMW)

	api2.RegisterHandlersWithOptions(e, handler, api2.RegisterHandlersOptions{
		BaseURL:              "/api",
		OperationMiddlewares: opMW,
	})

	return e
}
