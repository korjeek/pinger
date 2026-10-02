package http

import (
	"log/slog"

	"github.com/korjeek/pinger/backend/internal/controller/http/api"
	mdl "github.com/korjeek/pinger/backend/internal/controller/http/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const baseUrl = "/api"

type ServerConfig struct {
	Logger *slog.Logger

	AccessTokenParser mdl.AccessTokenParser

	CookiePath   string
	CookieDomain string
	CookieSecure bool
}

func NewServer(cfg ServerConfig, impl api.StrictServerInterface) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLogger())

	auth := mdl.JWTAuth(cfg.AccessTokenParser)
	opMW := map[string][]echo.MiddlewareFunc{
		"updateUser":              {auth},
		"getMonitors":             {auth},
		"createMonitor":           {auth},
		"deleteMonitor":           {auth},
		"updateMonitor":           {auth},
		"getSnapshotsByMonitorID": {auth},
	}

	strictMW := []api.StrictMiddlewareFunc{
		mdl.StrictErrorMiddleware,
	}

	handler := api.NewStrictHandler(impl, strictMW)

	api.RegisterHandlersWithOptions(e, handler, api.RegisterHandlersOptions{
		BaseURL:              baseUrl,
		OperationMiddlewares: opMW,
	})

	return e
}
