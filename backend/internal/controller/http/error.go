package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/korjeek/pinger/backend/internal/domain"
	"github.com/labstack/echo/v4"
)

func httpErrorHandler(err error, c echo.Context) {
	response := api.ErrorResponse{
		Error:     &errCode,
		Message:   &clientMsg,
		Timestamp: &now,
		TraceId:   &traceID,
		Details:   nil,
	}
}
