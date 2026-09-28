package api

import "github.com/labstack/echo/v4"

func NewServer() {
	e := echo.New()
	e.HTTPErrorHandler = httpErrorHandler
}
