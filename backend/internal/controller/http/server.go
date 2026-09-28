package api

import (
	"log/slog"
	"net/http"

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

func NewServer(cfg ServerConfig, impl StrictServerInterface) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.HTTPErrorHandler = HTTPErrorHandler(cfg.Logger)

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

	strictMW := []StrictMiddlewareFunc{
		RefreshCookieWriterMiddleware(cookieConfig{
			Path:     cfg.CookiePath,
			Domain:   cfg.CookieDomain,
			Secure:   cfg.CookieSecure,
			SameSite: http.SameSiteStrictMode,
		}),
	}

	handler := NewStrictHandler(impl, strictMW)

	RegisterHandlersWithOptions(e, handler, RegisterHandlersOptions{
		BaseURL:              "/api",
		OperationMiddlewares: opMW,
	})

	return e
}
