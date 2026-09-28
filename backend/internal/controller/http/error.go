package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/korjeek/pinger/backend/pkg/apperr"
	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler(logger *slog.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		status, code, message, details := classify(err, logger, c)

		resp := ErrorResponse{
			Code:      code,
			Message:   message,
			Details:   details,
			Timestamp: new(time.Now().UTC()),
		}
		if traceID := TraceIDFromEcho(c); traceID != "" {
			resp.TraceId = &traceID
		}

		if c.Request().Method == http.MethodHead {
			_ = c.NoContent(status)
			return
		}
		_ = c.JSON(status, resp)
	}
}

func classify(err error, logger *slog.Logger, c echo.Context) (int, string, string, *map[string]any) {
	if appErr, ok := apperr.As(err); ok {
		var details *map[string]any
		if len(appErr.Details) > 0 {
			details = &appErr.Details
		}
		if appErr.Cause != nil && appErr.Status >= 500 {
			logger.Error("request failed",
				"err", appErr.Cause,
				"code", appErr.Code,
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"trace_id", TraceIDFromEcho(c),
			)
		}
		return appErr.Status, string(appErr.Code), appErr.Message, details
	}

	if httpErr, ok := errors.AsType[*echo.HTTPError](err); ok {
		status := httpErr.Code
		code := http.StatusText(status)
		if code == "" {
			code = fmt.Sprintf("HTTP_%d", status)
		}
		msg, _ := httpErr.Message.(string)
		if msg == "" {
			msg = code
		}
		if status >= 500 {
			logger.Error("http error", "status", status, "err", err,
				"trace_id", TraceIDFromEcho(c))
		}
		return status, code, msg, nil
	}

	logger.Error("unhandled error", "err", err,
		"method", c.Request().Method,
		"path", c.Request().URL.Path,
		"trace_id", TraceIDFromEcho(c),
	)
	return http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error", nil
}
