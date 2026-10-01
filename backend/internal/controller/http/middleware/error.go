package middleware

import (
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/korjeek/pinger/backend/internal/controller/http/api"
	"github.com/korjeek/pinger/backend/internal/controller/http/utils"
	"github.com/korjeek/pinger/backend/pkg/apperr"
	"github.com/labstack/echo/v4"
)

func StrictErrorMiddleware(f api.StrictHandlerFunc, operationID string) api.StrictHandlerFunc {
	return func(ctx echo.Context, request any) (any, error) {
		response, err := f(ctx, request)
		if err == nil {
			return response, nil
		}

		if appErr, ok := errors.AsType[apperr.Error](err); ok {
			traceID := uuid.New().String()

			slog.Error("domain error occurred",
				"operation_id", operationID,
				"code", appErr.Code.String(),
				"message", appErr.Message,
				"cause", appErr.Cause,
				"details", appErr.Details,
				"trace_id", traceID,
			)

			var outDetails *map[string]any
			if len(appErr.Details) > 0 {
				outDetails = &appErr.Details
			}

			apiError := api.ErrorResponse{
				Code:      appErr.Code.String(),
				Message:   appErr.Message,
				Details:   outDetails,
				Timestamp: new(time.Now()),
				TraceId:   &traceID,
			}

			return nil, echo.NewHTTPError(utils.CodeInt(appErr.Code), apiError)
		}

		return nil, err
	}
}
